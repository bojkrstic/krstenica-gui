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
		// targetFile := "/tmp"

		// // ctx.JSON(http.StatusOK, krstenica)
		// targetFileBasename := filepath.Base(targetFile)

		// targetFolder := strings.TrimRight(targetFile, "/"+targetFileBasename)
		// defer os.RemoveAll(targetFolder)

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
		// ctx.Writer.Header().Add("Content-Disposition", "attachment; filename="+targetFileBasename)
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
	// Kreiranje novog Excel fajla
	f := excelize.NewFile()

	// Dodavanje Sheet-a
	sheetName := "Krstenica"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		log.Println("Ne moze da kreira sheetName")
	}
	f.SetActiveSheet(index)

	// Naslovi kolona
	headers := []string{"ID", "Ime", "Prezime", "Datum Rođenja", "Mesto Rođenja"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i))) // Generiše A1, B1, C1...
		f.SetCellValue(sheetName, cell, header)
	}

	return nil
}

//****************************************************end******Krstenica Print*************************************
