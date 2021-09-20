package logger

import (
	"io"
	"os"
	"path"

	"github.com/rs/zerolog"
	//"github.com/rs/zerolog/log"
	"github.com/codemaestro64/skeleton/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*zerolog.Logger
}

func New(cfg *config.LoggerConfig, env config.Environment) *Logger {
	var writers []io.Writer

	filename := path.Join(cfg.Directory, "log.txt")
	writers = append(writers, newRollingFile(cfg, filename))
	if env == config.DevelopmentEnv {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	mw := io.MultiWriter(writers...)

	// zerolog.SetGlobalLevel(zerolog.DebugLevel)
	logger := zerolog.New(mw).With().Timestamp().Logger()

	return &Logger{
		Logger: &logger,
	}
}

func newRollingFile(cfg *config.LoggerConfig, filename string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxBackups: cfg.MaxBackups, // files
		MaxSize:    cfg.MaxSize,    // megabytes
		MaxAge:     cfg.MaxAge,     // days
	}
}
