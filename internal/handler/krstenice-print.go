package handler

import (
	"fmt"
	"io"
	"krstenica/internal/dto"
	"krstenica/internal/errorx"
	"krstenica/pkg"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// var invoiceXlsxTemplateFile = "/home/krle/develop/horisen/Krstenica-new/Krstenica-Tane/krstenica/doc/template_files/krstenica-template-empty.xlsx"
// var invoiceXlsxTemplateFilePreview = "/home/krle/develop/horisen/Krstenica-new/Krstenica-Tane/krstenica/doc/template_files/krstenica-template.xlsx"

var templateFileRelative = "doc/template_files/krstenica-template-empty.xlsx"
var templateEmptyFileRelative = "doc/template_files/krstenica-template.xlsx"

// *************************************************************Krstenica Print*************************************
func (h *httpHandler) getKrstenicePrint() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cx := ctx.Request.Context()
		filters := pkg.ParseUrlQuery(ctx)
		log.Println("filters", filters)

		krstenica, err := h.service.GetKrstenicaByID(cx, int64(id))
		if err != nil {
			if err == errorx.ErrKrstenicaNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var file string
		templateDir := resolveDir("doc/template_files")
		invoiceXlsxTemplateFilePreview := filepath.Join(templateDir, filepath.Base(templateFileRelative))
		invoiceXlsxTemplateFile := filepath.Join(templateDir, filepath.Base(templateEmptyFileRelative))

		fmt.Println("Template file preview:", invoiceXlsxTemplateFilePreview)
		fmt.Println("Empty template file:", invoiceXlsxTemplateFile)

		if v, ok := filters.Filters[pkg.FilterKey{Property: "preview", Operator: "eq"}]; ok && len(v) > 0 && v[0] == "true" {
			file = invoiceXlsxTemplateFilePreview
		} else {
			file = invoiceXlsxTemplateFile
		}

		outputFormat := "xlsx"
		if v, ok := filters.Filters[pkg.FilterKey{Property: "format", Operator: "eq"}]; ok && len(v) > 0 {
			outputFormat = strings.ToLower(strings.TrimSpace(v[0]))
		}
		if outputFormat == "" {
			outputFormat = "xlsx"
		}

		targetDir, err := os.MkdirTemp("", "krstenica")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create temp directory"})
			return
		}
		defer os.RemoveAll(targetDir)

		backgroundImage := resolveFile("krstenica_obrada.jpg")
		backgroundFullBleed := true
		fontKey := ""
		if v, ok := filters.Filters[pkg.FilterKey{Property: "template_version", Operator: "eq"}]; ok && len(v) > 0 {
			version := strings.TrimSpace(strings.ToLower(v[0]))
			switch version {
			case "2", "v2", "verzija2", "version2":
				backgroundImage = ""
				backgroundFullBleed = false
			}
		}
		if v, ok := filters.Filters[pkg.FilterKey{Property: "font", Operator: "eq"}]; ok && len(v) > 0 {
			fontKey = strings.TrimSpace(v[0])
		}

		var (
			targetFile   string
			contentType  string
			downloadName string
		)

		switch outputFormat {
		case "pdf":
			targetFile = filepath.Join(targetDir, "krstenica.pdf")
			if err := fillKrstenicaPDFFile(krstenica, file, targetFile, backgroundImage, backgroundFullBleed, fontKey); err != nil {
				log.Println("Error generating PDF file:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate PDF file: %v", err)})
				return
			}
			contentType = "application/pdf"
			downloadName = "krstenica.pdf"
		default:
			from, err := os.Open(file)
			if err != nil {
				log.Println("Can't open Excel template file:", err)
				return
			}
			defer from.Close()

			targetFile = filepath.Join(targetDir, "krstenica.xlsx")
			to, err := os.OpenFile(targetFile, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				log.Print(err)
				return
			}
			defer to.Close()

			if _, err = io.Copy(to, from); err != nil {
				log.Print(err)
				return
			}

			if err := fillKrstenicaExcelFile(krstenica, targetFile, backgroundImage, backgroundFullBleed); err != nil {
				log.Println("Error generating Excel file:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate Excel file: %v", err)})
				return
			}

			contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
			downloadName = "krstenica.xlsx"
		}

		fi, err := os.Stat(targetFile)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "file not found"})
			return
		}

		size := fi.Size()

		ctx.Writer.Header().Set("Content-Type", contentType)
		ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", size))
		ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", downloadName))
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")

		b, err := os.ReadFile(targetFile)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
			return
		}

		n, err := ctx.Writer.Write(b)
		if err != nil {
			log.Println("Error while writing file to response:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send file"})
			return
		}

		if int64(n) != size {
			log.Println("Incomplete file transfer:", n, "bytes written, expected:", size)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "file transfer incomplete"})
			return
		}

	}
}

func getKrstenicaCellValues(krstenica *dto.Krstenica) map[string]string {
	values := map[string]string{
		"C1":  krstenica.Book,
		"C2":  formatInt(krstenica.Page),
		"C3":  formatInt(krstenica.CurrentNumber),
		"F8":  krstenica.EparhijaName,
		"C10": krstenica.TampleName,
		"I10": krstenica.TampleCity,
		"F13": formatDateTimeComma(krstenica.BirthDate),
		// "I17": krstenica.Country,
		"F19": formatDateComma(krstenica.Baptism),
		"G24": krstenica.TampleCity,
		"I24": krstenica.TampleName,
		"F27": krstenica.FirstName,
		"I27": mapGenderToCyrillic(krstenica.Gender),
		"F30": krstenica.ParentFirstName,
		"I30": krstenica.ParentLastName,
		"K30": krstenica.ParentOccupation,
		"F31": krstenica.ParentCity,
		"I31": krstenica.ParentReligion,
		"I32": strings.TrimSpace(krstenica.BirthOrder),
		"I36": strings.TrimSpace(krstenica.IsChurchMarried),
		"I38": strings.TrimSpace(krstenica.IsTwin),
		"I41": strings.TrimSpace(krstenica.HasPhysicalDisability),
		"F43": krstenica.PriestFirstName,
		"H43": krstenica.PriestLastName,
		"K43": strings.TrimSpace(krstenica.PriestTitle),
		"E48": krstenica.GodfatherFirstName,
		"G48": krstenica.GodfatherLastName,
		"E49": krstenica.GodfatherCity,
		"G49": krstenica.GodfatherReligion,
		"E51": krstenica.Anagrafa,
		"C54": krstenica.Comment,
		"B62": strings.TrimSpace(krstenica.NumberOfCertificate),
		"B63": "",
		"C63": "",
		"B65": krstenica.TownOfCertificate,
		// "F60": strings.TrimSpace(fmt.Sprintf("%s %s", krstenica.ParohFirstName, krstenica.ParohLastName)),
		// "C62": krstenica.Status,
	}

	if values["K48"] == "" {
		values["K48"] = strings.TrimSpace(krstenica.GodfatherOccupation)
	}

	placeBirth := strings.TrimSpace(krstenica.PlaceOfBirthday)
	municipalityBirth := strings.TrimSpace(krstenica.MunicipalityOfBirthday)
	values["E16"] = ""
	values["F16"] = ""
	values["G16"] = ""

	switch {
	case placeBirth != "" && municipalityBirth != "" && strings.EqualFold(placeBirth, municipalityBirth):
		values["F16"] = municipalityBirth
	case placeBirth == "" && municipalityBirth != "":
		values["G16"] = municipalityBirth
	case municipalityBirth == "" && placeBirth != "":
		values["E16"] = placeBirth
	default:
		values["E16"] = placeBirth
		values["G16"] = municipalityBirth
	}

	if !krstenica.Certificate.IsZero() {
		dayMonth, yearSuffix := splitDateDayMonthYearSuffix(krstenica.Certificate)
		values["B63"] = dayMonth
		values["C63"] = yearSuffix
	}

	if !krstenica.Baptism.IsZero() {
		values["N10"] = formatBaptismYear(krstenica.Baptism)
	} else {
		values["N10"] = ""
	}

	return values
}

func formatBaptismYear(t time.Time) string {
	year := t.Year()
	if year < 2000 {
		return fmt.Sprintf("%04d", year)
	}
	return t.Format("06")
}

func fillKrstenicaExcelFile(krstenica *dto.Krstenica, targetFile string, backgroundImage string, fullBleed bool) error {

	// Proveriti da li fajl postoji
	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
		// Ako ne postoji, kreiraj novi Excel fajl
		xlsxEx := excelize.NewFile()

		// Snimi inicijalno prazan fajl da bi mogao kasnije da se otvori
		if err := xlsxEx.SaveAs(targetFile); err != nil {
			log.Println("Greška pri kreiranju fajla:", err)
			return err
		}
		xlsxEx.Close() // Zatvaranje inicijalnog fajla
	}

	xlsxEx, err := excelize.OpenFile(targetFile)
	if err != nil {
		log.Print("Greška pri otvaranju fajla:", err)
		return err
	}
	defer xlsxEx.Close()

	sheetName := "krstenica"
	if idx, err := xlsxEx.GetSheetIndex(sheetName); err != nil || idx == -1 {
		sheetName = xlsxEx.GetSheetName(xlsxEx.GetActiveSheetIndex())
	}

	if backgroundImage != "" {
		if _, err := os.Stat(backgroundImage); err == nil {
			if err := addBackgroundPicture(xlsxEx, sheetName, backgroundImage, fullBleed); err != nil {
				log.Println("Ne može da doda pozadinsku sliku:", err)
			}
		} else {
			log.Println("Pozadinska slika nije pronađena:", backgroundImage)
		}
	}

	set := func(cell, value string) {
		if err := xlsxEx.SetCellValue(sheetName, cell, value); err != nil {
			log.Printf("set cell %s failed: %v", cell, err)
		}
	}

	for cell, value := range getKrstenicaCellValues(krstenica) {
		set(cell, value)
	}

	setCellBold(xlsxEx, sheetName, "F27")

	// Snimanje fajla
	if err := xlsxEx.SaveAs(targetFile); err != nil {
		log.Println("Greška pri čuvanju fajla:", err)
		return err
	}

	return nil
}

func setCellBold(xlsxEx *excelize.File, sheetName, cell string) {
	style := &excelize.Style{}
	if styleID, err := xlsxEx.GetCellStyle(sheetName, cell); err == nil {
		if existing, err := xlsxEx.GetStyle(styleID); err == nil && existing != nil {
			style = existing
		} else if err != nil {
			log.Printf("get style definition for %s failed: %v", cell, err)
		}
	} else {
		log.Printf("get style for %s failed: %v", cell, err)
	}

	if style.Font == nil {
		style.Font = &excelize.Font{}
	}
	if style.Font.Bold {
		return
	}
	style.Font.Bold = true

	boldStyleID, err := xlsxEx.NewStyle(style)
	if err != nil {
		log.Printf("create bold style for %s failed: %v", cell, err)
		return
	}
	if err := xlsxEx.SetCellStyle(sheetName, cell, cell, boldStyleID); err != nil {
		log.Printf("apply bold style to %s failed: %v", cell, err)
	}
}

func formatDateTime(t time.Time) string {
	return formatSerbianDateTime(t)
}

func formatDateTimeComma(t time.Time) string {
	return formatSerbianDateTime(t)
}

func formatDateComma(t time.Time) string {
	return formatSerbianDate(t)
}

func formatDate(t time.Time) string {
	return formatSerbianDate(t)
}

func splitDateDayMonthYearSuffix(t time.Time) (string, string) {
	if t.IsZero() {
		return "", ""
	}
	dayMonth := t.Format("02.01.")
	yearSuffix := t.Format("06")
	return dayMonth, yearSuffix
}

func formatInt(v int64) string {
	if v == 0 {
		return ""
	}
	return fmt.Sprintf("%d", v)
}

func boolToYesNo(b bool) string {
	if b {
		return "DA"
	}
	return "NE"
}

func boolToYesNoCyrillic(b bool) string {
	if b {
		return "Да"
	}
	return "Не"
}

func mapGenderToCyrillic(gender string) string {
	switch strings.ToLower(strings.TrimSpace(gender)) {
	case "m", "musko", "male", "muško":
		return "Мушко"
	case "z", "zensko", "female", "žensko":
		return "Женско"
	default:
		return gender
	}
}

var serbianMonths = []string{
	"",
	"јануар",
	"фебруар",
	"март",
	"април",
	"мај",
	"јун",
	"јул",
	"август",
	"септембар",
	"октобар",
	"новембар",
	"децембар",
}

func formatSerbianDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	local := t.In(time.Local)
	monthIdx := int(local.Month())
	monthName := ""
	if monthIdx >= 1 && monthIdx <= 12 {
		monthName = serbianMonths[monthIdx]
	}

	return fmt.Sprintf("%d    %s    %d", local.Year(), monthName, local.Day())
}

func formatSerbianDateTime(t time.Time) string {
	date := formatSerbianDate(t)
	if date == "" {
		return ""
	}

	local := t.In(time.Local)
	if hasClockComponent(local) {
		return fmt.Sprintf("%s    у    %02d:%02d часова", date, local.Hour(), local.Minute())
	}

	return date
}

func hasClockComponent(t time.Time) bool {
	return t.Hour() != 0 || t.Minute() != 0 || t.Second() != 0
}

func filterEmpty(lines []string) []string {
	res := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		res = append(res, trimmed)
	}
	return res
}

const yearLabelColumn = 12

func formatYearLine(label, year string) string {
	padded := padToColumn(label, yearLabelColumn)
	return fmt.Sprintf("%s%s    год.", padded, year)
}

func padToColumn(label string, column int) string {
	width := utf8.RuneCountInString(label)
	spaces := column - width
	if spaces < 1 {
		spaces = 1
	}
	return label + strings.Repeat(" ", spaces)
}

func addBackgroundPicture(file *excelize.File, sheetName, imagePath string, fullBleed bool) error {
	printObject := true
	scaleX := 1.6
	scaleY := 1.55
	if fullBleed {
		scaleX = 2.2
		scaleY = 2.2
	}
	options := &excelize.GraphicOptions{
		OffsetX:     0,
		OffsetY:     0,
		ScaleX:      scaleX,
		ScaleY:      scaleY,
		Positioning: "moveAndSize",
		PrintObject: &printObject,
	}
	if err := file.AddPicture(sheetName, "A1", imagePath, options); err != nil {
		return err
	}
	return nil
}
