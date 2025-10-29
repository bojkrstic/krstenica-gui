package handler

import (
	"context"
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

		cx := context.Background()
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

		var (
			targetFile   string
			contentType  string
			downloadName string
		)

		switch outputFormat {
		case "pdf":
			targetFile = filepath.Join(targetDir, "krstenica.pdf")
			if err := fillKrstenicaPDFFile(krstenica, file, targetFile, backgroundImage); err != nil {
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

			if err := fillKrstenicaExcelFile(krstenica, targetFile, backgroundImage); err != nil {
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
		"C2":  krstenica.Book,
		"C3":  formatInt(krstenica.Page),
		"C4":  formatInt(krstenica.CurrentNumber),
		"C9":  krstenica.EparhijaName,
		"C11": krstenica.TampleName,
		"H11": krstenica.TampleCity,
		"F14": formatDateTimeComma(krstenica.BirthDate),
		"E17": krstenica.PlaceOfBirthday,
		"G17": krstenica.MunicipalityOfBirthday,
		// "I17": krstenica.Country,
		"E20": formatDateTimeComma(krstenica.Baptism),
		"F24": krstenica.TampleCity,
		"H24": krstenica.TampleName,
		"D27": krstenica.FirstName,
		"F27": krstenica.LastName,
		"H27": mapGenderToCyrillic(krstenica.Gender),
		"E30": krstenica.ParentFirstName,
		"G30": krstenica.ParentLastName,
		"I30": krstenica.ParentOccupation,
		"E31": krstenica.ParentCity,
		"G31": krstenica.ParentReligion,
		"G32": strings.TrimSpace(krstenica.BirthOrder),
		"E35": strings.TrimSpace(krstenica.IsChurchMarried),
		"E37": boolToYesNoCyrillic(krstenica.IsTwin),
		"G39": boolToYesNoCyrillic(krstenica.HasPhysicalDisability),
		"F42": krstenica.PriestFirstName,
		"H42": krstenica.PriestLastName,
		"E45": krstenica.GodfatherFirstName,
		"G45": krstenica.GodfatherLastName,
		"I45": krstenica.GodfatherOccupation,
		"E46": krstenica.GodfatherCity,
		"G46": krstenica.GodfatherReligion,
		"E49": krstenica.Anagrafa,
		"C51": krstenica.Comment,
		"B58": formatInt(krstenica.NumberOfCertificate),
		"B60": "",
		"C60": "",
		"B61": krstenica.TownOfCertificate,
		// "F60": strings.TrimSpace(fmt.Sprintf("%s %s", krstenica.ParohFirstName, krstenica.ParohLastName)),
		// "C62": krstenica.Status,
	}

	if !krstenica.Certificate.IsZero() {
		dayMonth, yearSuffix := splitDateDayMonthYearSuffix(krstenica.Certificate)
		values["B60"] = dayMonth
		values["C60"] = yearSuffix
	}

	if !krstenica.Baptism.IsZero() {
		values["K11"] = fmt.Sprintf("%d", krstenica.Baptism.Year())
	} else {
		values["K11"] = ""
	}

	return values
}

func fillKrstenicaExcelFile(krstenica *dto.Krstenica, targetFile string, backgroundImage string) error {

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
			if err := addBackgroundPicture(xlsxEx, sheetName, backgroundImage); err != nil {
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

	// Snimanje fajla
	if err := xlsxEx.SaveAs(targetFile); err != nil {
		log.Println("Greška pri čuvanju fajla:", err)
		return err
	}

	return nil
}

func formatDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02.01.2006. 15:04")
}

func formatDateTimeComma(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006, 01, 02, 15:04")
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02.01.2006.")
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

func addBackgroundPicture(file *excelize.File, sheetName, imagePath string) error {
	printObject := true
	if err := file.AddPicture(sheetName, "A1", imagePath, &excelize.GraphicOptions{
		OffsetX:     0,
		OffsetY:     0,
		ScaleX:      1.6,
		ScaleY:      1.55,
		Positioning: "moveAndSize",
		PrintObject: &printObject,
	}); err != nil {
		return err
	}
	return nil
}
