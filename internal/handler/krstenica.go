package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *httpHandler) createKrstenica() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// req := &dto.ManualWalletUpdateReq{}

		// if err := ctx.Bind(req); err != nil {
		// 	fmt.Println("Error when parsing body", err)
		// 	ctx.JSON(400, gin.H{"error": "error when parsing data"})
		// 	return
		// }

		// var err error

		// brandID, err := middleware.GetBrandIDFromContext(ctx)
		// if err != nil || brandID == 0 {
		// 	ctx.JSON(400, gin.H{"error": err.Error()})
		// 	return
		// }

		// req.BrandId = brandID

		// adminID, err := middleware.GetUserIdFromContext(ctx)
		// if err != nil || brandID == 0 {
		// 	ctx.JSON(400, gin.H{"error": err.Error()})
		// 	return
		// }

		// err = h.service.ManualBalanceUpdate(ctx, req, adminID)

		// if err != nil {
		// 	fmt.Println("Error when maually updating balance", err)
		// 	ctx.JSON(400, gin.H{"error": "error when manually updating balance"})
		// 	return
		// }

	}
}

func (h *httpHandler) getKrstenica() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func (h *httpHandler) updateKrstenica() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func (h *httpHandler) deleteKrstenica() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
