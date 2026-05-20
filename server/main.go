package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/joho/godotenv"
	server "github.com/levi-20/url-shortner/server"
	"github.com/lmittmann/tint"
)

func InitLogger(environment server.Environment) {
	var handler slog.Handler
	if environment == server.EnvironmentProduction {
		slog.SetLogLoggerLevel(slog.LevelInfo)
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})
	} else {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			AddSource:  true,
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		})
	}

	slog.SetDefault(slog.New(handler))
}

func main() {

	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic", "error", r, "stack", string(debug.Stack()))
			time.Sleep(2 * time.Second)
			panic(r)
		}
	}()

	godotenv.Load()

	server.GlobalServerConfig.Env = server.EnvironmentDebug
	server.GlobalServerConfig.ListenPort = 8900

	InitLogger(server.GlobalServerConfig.Env)

	exec, err := os.Executable()
	if err != nil {
		panic(err)
	}

	server.GlobalServerConfig.ExecBaseDir = filepath.Dir(exec)

	listenAddr := fmt.Sprintf(":%d", server.GlobalServerConfig.ListenPort)
	slog.Info(fmt.Sprintf("📦 Starting GIN server on %s", listenAddr))
	app := server.NewApp(server.GlobalServerConfig)

	server.StartGinServer(app, listenAddr)
}
