package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levi-20/url-shortner/db"
)

type ShortenRequest struct {
	URLs []string `json:"urls" binding:"required,min=1,max=20,dive,url"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func getUrlHash(URL string) string {

	sum := sha256.Sum256([]byte(URL))
	return hex.EncodeToString(sum[:])

}

type NewAndExisting struct {
	New     []db.ShortWithOriginal `json:"new"`
	Exising []db.ShortWithOriginal `json:"existing"`
}

func formatResponse(existingUrls []db.ShortWithOriginal, newShortCodes []string, newOriginalURLs []string) *NewAndExisting {

	base_url := os.Getenv("BASE_URL")

	var response NewAndExisting
	response.Exising = make([]db.ShortWithOriginal, 0, len(existingUrls))
	response.New = make([]db.ShortWithOriginal, 0, len(newShortCodes))
	if len(existingUrls) > 0 {

		for _, urlObj := range existingUrls {
			response.Exising = append(response.Exising, db.ShortWithOriginal{
				Short: fmt.Sprintf("%s/%s", base_url, urlObj.Short),
				Url:   urlObj.Url,
			})
		}
	}

	if len(newShortCodes) > 0 {

		for i, code := range newShortCodes {
			response.New = append(response.New, db.ShortWithOriginal{
				Short: fmt.Sprintf("%s/%s", base_url, code),
				Url:   newOriginalURLs[i],
			})
		}
	}
	return &response
}

func GenerateShortURL() string {

	bytes := make([]byte, 7)
	rand.Read(bytes)

	// Use URLEncoding to ensure the string is URL-safe
	return base64.URLEncoding.EncodeToString(bytes)[:7]
}

func CheckAndSaveURLs(ctx *gin.Context, pool *pgxpool.Pool, urls []string) (*NewAndExisting, error) {

	hashesToCheck := make([]string, 0, len(urls))
	for _, url := range urls {
		hashesToCheck = append(hashesToCheck, getUrlHash(url))
	}

	slog.Info("Existing URL hashes", "hashes", hashesToCheck)
	existingURLs, err := db.CheckExistingURLs(ctx, pool, hashesToCheck)
	if err != nil {
		slog.Error("Error while checking for existing URLs", "error", err.Error())
		return nil, err
	}

	urlCount := len(urls)
	realURL := make([]string, 0, urlCount)
	shortURL := make([]string, 0, urlCount)
	urlHashes := make([]string, 0, urlCount)

	for _, url := range urls {

		hash := getUrlHash(url)
		if _, ok := existingURLs.Hashes[hash]; ok {
			continue
		}
		realURL = append(realURL, url)
		urlHashes = append(urlHashes, hash)
		shortURL = append(shortURL, GenerateShortURL())
	}

	slog.Info("URLs", "realURL", realURL, "shortURL", shortURL, "hashes", urlHashes)

	if len(shortURL) > 0 {
		err = db.SaveURLs(ctx, pool, realURL, shortURL, urlHashes)
		if err != nil {
			slog.Error("error while saving the urls", "err", err.Error())
			return nil, err
		}
	}

	response := formatResponse(existingURLs.Urls, shortURL, realURL)
	return response, nil
}

func ShortenUrls(ctx *gin.Context, pool *pgxpool.Pool) {

	var URLRequest ShortenRequest

	if err := ctx.ShouldBindJSON(&URLRequest); err != nil {
		slog.Error("error while body parse", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: "Error while pasrsing URLs from request body"})
		return
	}

	response, err := CheckAndSaveURLs(ctx, pool, URLRequest.URLs)
	if err != nil {
		slog.Error("Error while Saving", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error while saving the short URLs"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func RedirectUrl(ctx *gin.Context, pool *pgxpool.Pool) {

	code := ctx.Param("slug")
	url := db.GetUrlWithCode(ctx, pool, code)
	if url == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "no url found with short code"})
		return
	}
	db.RegisterClick(ctx, pool, code)
	ctx.Redirect(http.StatusFound, *url)
}
