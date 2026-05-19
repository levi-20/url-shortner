package server

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levi-20/url-shortner/db"
	handlers "github.com/levi-20/url-shortner/handlers"
)

type Environment string

const (
	EnvironmentDebug      Environment = "debug"
	EnvironmentProduction Environment = "production"
)

type ServerConfig struct {
	Env         Environment
	ExecBaseDir string
	ListenPort  uint16
}

type App struct {
	DatabasePool *pgxpool.Pool
}

var GlobalServerConfig ServerConfig

const LoggerKey = "l"
const RequestIdKey = "rid"

func NewApp(serverConfig ServerConfig) *App {

	app := &App{
		DatabasePool: db.CreateDatabaseConnection(),
	}

	return app
}

func StartGinServer(app *App, addr string) {

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	// @Levi TODO
	// router.Use(GinStructuredLogger(slog.Default()))
	// router.Use(CORSMiddleware())

	router.POST("/generate", func(ctx *gin.Context) {

		urls, err := handlers.ShortenUrls(ctx, app.DatabasePool)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}

		ctx.JSON(http.StatusOK, urls)
	})

	if err := router.Run(addr); err != nil {
		slog.Error("Failed to start router", "address", addr, "error", err)
		os.Exit(-2)
	}
}
