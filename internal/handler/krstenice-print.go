package handler

import (
	"context"
	"fmt"
	"krstenica/internal/dto"
	"krstenica/internal/errorx"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// *************************************************************Krstenica Print*************************************
func (h *httpHandler) getKrstenicePrint() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cx := context.Background()

		krstenica, err := h.service.GetKrstenicaByID(cx, int64(id))
		if err != nil {
			if err == errorx.ErrKrstenicaNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Using a temporary directory
		targetDir, err := os.MkdirTemp("", "krstenica")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create temp directory"})
			return
		}
		defer os.RemoveAll(targetDir) // Cleanup after function execution
		targetFile := filepath.Join(targetDir, "krstenica.xlsx")

		// Generate Excel file
		err = fillKrstenicaExcelFile(ctx, krstenica, targetFile)
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

func fillKrstenicaExcelFile(ctx context.Context, krstenica *dto.Krstenica, targetFile string) error {

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
	xlsxEx.SetCellValue(sheetName, "A3", "Knjiga")
	xlsxEx.SetCellValue(sheetName, "B3", krstenica.Book)
	xlsxEx.SetCellValue(sheetName, "A4", "Strana knjige")
	xlsxEx.SetCellValue(sheetName, "B4", krstenica.Page)
	xlsxEx.SetCellValue(sheetName, "A5", "Tekuci broj")
	xlsxEx.SetCellValue(sheetName, "B5", krstenica.CurrentNumber)

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
