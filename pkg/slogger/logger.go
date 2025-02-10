package slogger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/sgrumley/gotex/pkg/path"
)

type Logger struct {
	*slog.Logger
}

type (
	Handler string
)

const (
	HandlerJSON Handler = "json"
	HandlerText Handler = "text"
)

type (
	Option        func(*LoggerOptions)
	LoggerOptions struct {
		level   slog.Level
		format  Handler
		output  *os.File
		source  bool
		handler slog.Handler
	}
)

func WithLevel(level slog.Level) Option {
	return func(opts *LoggerOptions) {
		opts.level = level
	}
}

// Replace this with WithHandler??
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

// WithHandler allows a custom handler to be set
// NOTE: not yet implemented
func WithHandler(handler slog.Handler) Option {
	return func(opts *LoggerOptions) {
		opts.handler = handler
	}
}

func New(options ...Option) (*Logger, error) {
	opts := LoggerOptions{
		level:  slog.LevelInfo,
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

	return &Logger{
		logger,
	}, nil
}

func newLogFile(format Handler) (*os.File, error) {
	defaultFolder := "~/.config/gotex"
	defaultFile := "out.log"

	defaultFolder, err := path.ReplaceHomeDirChar(defaultFolder)
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

func (l *Logger) Error(msg string, err error, args ...any) {
	l.Logger.Error(msg, append([]any{"error", err}, args...)...)
}

func (l *Logger) Fatal(msg string, err error) {
	l.Error(msg, err)
	os.Exit(1)
}

func (l *Logger) With(args ...any) *Logger {
	lw := l.Logger.With(args...)
	return &Logger{
		lw,
	}
}
