package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Handler string

const (
	HandlerJSON Handler = "json"
	HandlerText Handler = "text"
)

type (
	Option        func(*LoggerOptions)
	LoggerOptions struct {
		level  slog.Level
		format Handler
		output *os.File
		source bool
	}
)

func WithLevel(level slog.Level) Option {
	return func(opts *LoggerOptions) {
		opts.level = level
	}
}

func WithFormat(format Handler) Option {
	return func(opts *LoggerOptions) {
		opts.format = format
	}
}

func WithSource(s bool) Option {
	return func(opts *LoggerOptions) {
		opts.source = s
	}
}

func WithOutput(out *os.File) Option {
	return func(opts *LoggerOptions) {
		opts.output = out
	}
}

func New(options ...Option) (*slog.Logger, error) {
	opts := LoggerOptions{
		level: slog.LevelInfo,
		// format: HandlerText,
		format: HandlerJSON,
		source: false,
	}

	for _, opt := range options {
		opt(&opts)
	}

	if opts.output == nil {
		var err error

		opts.output, err = newLogFile(opts.format)
		if err != nil {
			return nil, err
		}
	}

	baseOpts := &slog.HandlerOptions{
		AddSource: opts.source,
		Level:     opts.level,
	}

	var slogHandler slog.Handler
	if opts.format == HandlerText {
		slogHandler = slog.NewTextHandler(opts.output, baseOpts)
	} else {
		slogHandler = slog.NewJSONHandler(opts.output, baseOpts)
	}
	logger := slog.New(slogHandler)
	slog.SetDefault(logger)

	return logger, nil
}

func newLogFile(format Handler) (*os.File, error) {
	defaultFolder := "~/.config/gotex"
	defaultFile := "out.log"

	defaultFolder, err := ReplaceHomeDirChar(defaultFolder)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(defaultFolder, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	if format == HandlerJSON {
		defaultFile = "log.json"
	}

	logFilePath := filepath.Join(defaultFolder, defaultFile)

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("Failed to open log file: %v", err)
	}

	return file, nil
}

func ReplaceHomeDirChar(fp string) (string, error) {
	if !strings.Contains(fp, "~") {
		return fp, nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error getting home directory: %w", err)
	}
	// Replace ~ with the home directory path
	if fp[:2] == "~/" {
		fp = filepath.Join(homeDir, fp[2:])
	}
	return fp, nil
}
