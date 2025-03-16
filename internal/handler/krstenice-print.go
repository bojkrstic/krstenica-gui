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

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

var invoiceXlsxTemplateFile = "/home/krle/develop/horisen/Krstenica-new/Krstenica-Tane/krstenica/doc/template_files/krstenica-template-empty.xlsx"
var invoiceXlsxTemplateFilePreview = "/home/krle/develop/horisen/Krstenica-new/Krstenica-Tane/krstenica/doc/template_files/krstenica-template.xlsx"

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
		err = fillKrstenicaExcelFile(krstenica, targetFile)
		if err != nil {
			log.Println("Error generating Excel file:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate Excel file"})
			return
		}

		// Dodavanje slike
		imagePath := "/home/krle/develop/horisen/Krstenica-new/Krstenica-Tane/krstenica/krstenica_obrada.jpg" // Podesite putanju do slike
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			log.Println("Slika ne postoji:", imagePath)
			return
		}

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

func fillKrstenicaExcelFile(krstenica *dto.Krstenica, targetFile string) error {

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

	dateFormat := "yyyy, mmmm, dd, hh:mm"
	dataStyle, err := xlsxEx.NewStyle(&excelize.Style{
		CustomNumFmt: &dateFormat,
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Dodavanje Sheet-a
	sheetName := "Krstenica"
	index, err := xlsxEx.NewSheet(sheetName)
	if err != nil {
		log.Println("Ne može da kreira sheetName:", err)
		return err
	}
	xlsxEx.SetActiveSheet(index)

	// Naslovi kolona
	// headers := []string{"ID", "Ime", "Prezime", "Datum Rođenja", "Mesto Rođenja"}
	// for i, header := range headers {
	// 	cell := fmt.Sprintf("%s1", string(rune('A'+i))) // Generiše A1, B1, C1...
	// 	xlsxEx.SetCellValue(sheetName, cell, header)
	// }
	// xlsxEx.SetCellValue(sheetName, "A3", "Knjiga")
	xlsxEx.SetCellValue(sheetName, "C2", krstenica.Book)
	// xlsxEx.SetCellValue(sheetName, "A4", "Strana knjige")
	xlsxEx.SetCellValue(sheetName, "C3", krstenica.Page)
	// xlsxEx.SetCellValue(sheetName, "A5", "Tekuci broj")
	xlsxEx.SetCellValue(sheetName, "C4", krstenica.CurrentNumber)
	xlsxEx.SetCellValue(sheetName, "C9", krstenica.EparhijaName)
	xlsxEx.SetCellValue(sheetName, "C11", krstenica.TampleName)
	xlsxEx.SetCellValue(sheetName, "H11", krstenica.TampleCity)
	xlsxEx.SetCellValue(sheetName, "K11", krstenica.Baptism.Year())
	// xlsxEx.SetCellValue(sheetName, "F14", krstenica.BirthDate.Format("2006-01-02"))
	// xlsxEx.SetCellStyle(sheetName, "F14", "F14", dataStyle)
	// column 1
	xlsxEx.SetCellValue(sheetName, "F14", krstenica.BirthDate)
	xlsxEx.SetCellStyle(sheetName, "F14", "F14", dataStyle)
	// column 2
	xlsxEx.SetCellValue(sheetName, "E17", krstenica.PlaceOfBirthday)
	xlsxEx.SetCellValue(sheetName, "G17", krstenica.MunicipalityOfBirthday)
	// column 3
	xlsxEx.SetCellValue(sheetName, "E20", krstenica.Baptism)
	xlsxEx.SetCellStyle(sheetName, "E20", "E20", dataStyle)

	xlsxEx.SetCellValue(sheetName, "F24", krstenica.TampleCity)
	xlsxEx.SetCellValue(sheetName, "H24", krstenica.TampleName)
	xlsxEx.SetCellValue(sheetName, "D27", krstenica.FirstName)
	xlsxEx.SetCellValue(sheetName, "F27", krstenica.LastName)
	xlsxEx.SetCellValue(sheetName, "H27", krstenica.Gender)
	xlsxEx.SetCellValue(sheetName, "E30", krstenica.ParentFirstName)
	xlsxEx.SetCellValue(sheetName, "G30", krstenica.ParentLastName)
	xlsxEx.SetCellValue(sheetName, "I30", krstenica.ParentOccupation)
	xlsxEx.SetCellValue(sheetName, "E31", krstenica.ParentCity)
	xlsxEx.SetCellValue(sheetName, "G31", krstenica.ParentReligion)
	//G32 narodnost da se doda
	xlsxEx.SetCellValue(sheetName, "G35", krstenica.NumberOfCertificate)
	xlsxEx.SetCellValue(sheetName, "E38", krstenica.IsChurchMarried)
	xlsxEx.SetCellValue(sheetName, "E41", krstenica.IsTwin)
	xlsxEx.SetCellValue(sheetName, "G44", krstenica.HasPhysicalDisability)
	xlsxEx.SetCellValue(sheetName, "F47", krstenica.PriestFirstName)
	xlsxEx.SetCellValue(sheetName, "H47", krstenica.PriestLastName)
	xlsxEx.SetCellValue(sheetName, "E51", krstenica.ParohFirstName)
	xlsxEx.SetCellValue(sheetName, "G51", krstenica.ParohLastName)
	xlsxEx.SetCellValue(sheetName, "E52", krstenica.ParentOccupation)
	xlsxEx.SetCellValue(sheetName, "E55", krstenica.Anagrafa)
	xlsxEx.SetCellValue(sheetName, "D58", krstenica.Status)

	// Snimanje fajla
	if err := xlsxEx.SaveAs(targetFile); err != nil {
		log.Println("Greška pri čuvanju fajla:", err)
		return err
	}

	return nil
}

func addBackgroundImageToExcel(targetFile, imagePath string) error {
	xlsxEx, err := excelize.OpenFile(targetFile)
	if err != nil {
		log.Print("Greška pri otvaranju fajla:", err)
		return err
	}
	defer xlsxEx.Close()

	sheetName := "Krstenica" // Postarajte se da ovaj sheet postoji
	printObject := false

	absPath, err := filepath.Abs(imagePath)
	if err != nil {
		log.Println("Greška pri dobijanju apsolutne putanje slike:", err)
	} else {
		log.Println("Apsolutna putanja slike:", absPath)
	}

	// imagePath1 := filepath.Join("/", "home", "krle", "develop", "horisen", "Krstenica-new", "Krstenica-Tane", "krstenica", "krstenica_obrada.jpg")
	// log.Println("Slika ne postoji na putanji:", imagePath1)
	// if _, err := os.Stat(imagePath1); os.IsNotExist(err) {
	// 	log.Println("Slika ne postoji na putanji:", imagePath1)
	// 	return fmt.Errorf("slika nije pronađena: %s", imagePath1)
	// }

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Println("Slika ne postoji na putanji:", imagePath)
		return fmt.Errorf("slika nije pronađena: %s", imagePath)
	}

	// Dodavanje slike
	err = xlsxEx.AddPicture(sheetName, "A1", imagePath,
		&excelize.GraphicOptions{
			OffsetX: 0, OffsetY: 0, // Postavite preciznu lokaciju
			ScaleX: 2.5, ScaleY: 2.5, // Povećajte sliku da pokrije više ćelija
			PrintObject: &printObject, // Neće se videti u štampi
		})
	if err != nil {
		log.Println("Greška pri dodavanju slike:", err)
		return err
	}

	// Snimanje fajla sa dodatom slikom
	if err := xlsxEx.SaveAs(targetFile); err != nil {
		log.Println("Greška pri čuvanju fajla:", err)
		return err
	}

	return nil
}
