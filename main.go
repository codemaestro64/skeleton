package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/codemaestro64/skeleton/config"
	"github.com/codemaestro64/skeleton/lib"
	"github.com/codemaestro64/skeleton/lib/logger"
	"github.com/codemaestro64/skeleton/web"
)

const (
	configFile = "config.toml"
)

func createDataDirIfNotExists(logDir string) error {
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directories: %s", err.Error())
	}

	return nil
}

func mainCore() error {
	cfg, err := config.Load(configFile)
	if err != nil {
		return err
	}

	// get data directory path
	dataDir := lib.AppDataDir(strings.ToLower(cfg.App.Name), true)
	logDir := filepath.Join(dataDir, "logs")

	err = createDataDirIfNotExists(logDir)
	if err != nil {
		return err
	}
	cfg.Logger.Directory = logDir

	// init logger
	logger := logger.New(&cfg.Logger, cfg.App.Env)

	server, err := web.NewServer(cfg, logger)
	if err != nil {
		return err
	}

	return server.Serve()
}

func main() {
	if err := mainCore(); err != nil {
		log.Fatal(err)
	}
}
