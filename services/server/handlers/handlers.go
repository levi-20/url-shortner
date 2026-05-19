package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShortenRequest struct {
	URLs []string `json:"urls" binding:"required,min=1,max=20,dive,url"`
}

type ShortenedURL struct {
	Original string `json:"original"`
	Short    string `json:"short"`
}

func ShortenUrls(ctx *gin.Context, app *pgxpool.Pool) ([]string, error) {

	var req ShortenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	slog.Info("reached the method", "body", req)
	return []string{}, nil
}
