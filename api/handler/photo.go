package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"photo-gallery/database"
)

func PhotoHandler(db *database.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		photos, err := db.GetAll()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]string{
				"message": "unexpected error. please try again later",
			})
			return
		}

		ctx.JSON(http.StatusOK, photos)
	}
}
