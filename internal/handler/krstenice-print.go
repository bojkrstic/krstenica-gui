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

		v, ok := filters.Filters[pkg.FilterKey{Property: "preview", Operator: "eq"}]
		if ok && len(v) > 0 && v[0] == "true" {
			file = invoiceXlsxTemplateFilePreview
		} else {
			file = invoiceXlsxTemplateFile

		}

		//copying template
		from, err := os.Open(file)
		if err != nil {
			log.Println("Can't open Excel template file:", err)
			return
		}
		defer from.Close()

		// Using a temporary directory
		targetDir, err := os.MkdirTemp("", "krstenica")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create temp directory"})
			return
		}
		defer os.RemoveAll(targetDir) // Cleanup after function execution
		targetFile := filepath.Join(targetDir, "krstenica.xlsx")

		to, err := os.OpenFile(targetFile, os.O_RDWR|os.O_CREATE, 0666)
		log.Printf("to %v", to)
		if err != nil {
			log.Print(err)
			return
		}
		defer to.Close()

		_, err = io.Copy(to, from)
		if err != nil {
			log.Print(err)
			return
		}

		// Generate Excel file
		backgroundImage := resolveFile("krstenica_obrada.jpg")
		err = fillKrstenicaExcelFile(krstenica, targetFile, backgroundImage)
		if err != nil {
			log.Println("Error generating Excel file:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate Excel file"})
			return
		}

		// // Dodavanje slike
		// imagePath := "/home/krle/develop/horisen/Krstenica-new/Krstenica-Tane/krstenica/krstenica_obrada.jpg" // Podesite putanju do slike
		// if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		// 	log.Println("Slika ne postoji:", imagePath)
		// 	return
		// }

		// err = addBackgroundImageToExcel(targetFile, imagePath)
		// if err != nil {
		// 	log.Println("Error adding background image:", err)
		// }

		// Get file information
		fi, err := os.Stat(targetFile)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "file not found"})
			return
		}

		// get the size
		size := fi.Size()

		ctx.Writer.Header().Add("Content-type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		ctx.Writer.Header().Add("Content-length", fmt.Sprintf("%d", size))
		ctx.Writer.Header().Add("Content-Disposition", "attachment; filename=krstenica.xlsx")
		ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")

		// Read file content
		b, err := os.ReadFile(targetFile)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
			return
		}

		// Write file content to response
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

	set("C2", krstenica.Book)
	set("C3", formatInt(krstenica.Page))
	set("C4", formatInt(krstenica.CurrentNumber))
	set("C9", krstenica.EparhijaName)
	set("C11", krstenica.TampleName)
	set("H11", krstenica.TampleCity)
	if !krstenica.Baptism.IsZero() {
		set("K11", fmt.Sprintf("%d", krstenica.Baptism.Year()))
	}
	set("F14", formatDateTime(krstenica.BirthDate))
	set("E17", krstenica.PlaceOfBirthday)
	set("G17", krstenica.MunicipalityOfBirthday)
	set("I17", krstenica.Country)
	set("E20", formatDateTime(krstenica.Baptism))
	set("F24", krstenica.TampleCity)
	set("H24", krstenica.TampleName)
	set("D27", krstenica.FirstName)
	set("F27", krstenica.LastName)
	set("H27", krstenica.Gender)
	set("E30", krstenica.ParentFirstName)
	set("G30", krstenica.ParentLastName)
	set("I30", krstenica.ParentOccupation)
	set("E31", krstenica.ParentCity)
	set("G31", krstenica.ParentReligion)
	set("G35", formatInt(krstenica.BirthOrder))
	set("K35", formatInt(krstenica.NumberOfCertificate))
	set("E38", boolToYesNo(krstenica.IsChurchMarried))
	set("E41", boolToYesNo(krstenica.IsTwin))
	set("G44", boolToYesNo(krstenica.HasPhysicalDisability))
	set("F47", krstenica.PriestFirstName)
	set("H47", krstenica.PriestLastName)
	set("E51", krstenica.GodfatherFirstName)
	set("G51", krstenica.GodfatherLastName)
	set("I51", krstenica.GodfatherOccupation)
	set("E52", krstenica.GodfatherCity)
	set("G52", krstenica.GodfatherReligion)
	set("E55", krstenica.Anagrafa)
	set("C58", krstenica.Comment)
	set("I58", krstenica.TownOfCertificate)
	set("K58", formatDate(krstenica.Certificate))
	set("F60", strings.TrimSpace(fmt.Sprintf("%s %s", krstenica.ParohFirstName, krstenica.ParohLastName)))
	set("C62", krstenica.Status)

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

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02.01.2006.")
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
