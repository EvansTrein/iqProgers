package logs

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type CustomHandler struct {
	handler slog.Handler
	output  io.Writer
	attrs   []slog.Attr
	mu      sync.Mutex
}

func NewCustomHandler(output io.Writer, opts *slog.HandlerOptions) *CustomHandler {
	return &CustomHandler{
		handler: slog.NewTextHandler(output, opts),
		output:  output,
		mu:      sync.Mutex{},
	}
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	var buf bytes.Buffer

	buf.Write([]byte("\n"))

	buf.Write([]byte("Apps log: "))

	buf.Write([]byte(r.Time.Format(time.Stamp) + "\n"))

	level := r.Level.String()
	switch r.Level {
	case slog.LevelInfo:
		level = "\033[32m" + level + "\033[0m" // green color
	case slog.LevelError:
		level = "\033[31m" + level + "\033[0m" // red color
	case slog.LevelDebug:
		level = "\033[34m" + level + "\033[0m" // blue color
	case slog.LevelWarn:
		level = "\033[33m" + level + "\033[0m" // yellow color
	}
	buf.Write([]byte("level--> " + level + "\n"))

	buf.Write([]byte("\033[4m" + "message--> " + r.Message + "\033[0m" + "\n")) // underlined text

	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		source := "file--> " + f.File +
			"\ncode_line--> " + "\033[38;5;208m" + strconv.Itoa(f.Line) + "\033[0m" + "\n" // orange color
		buf.Write([]byte(source))
	}

	for _, attr := range h.attrs {
		if attr.Key == "operation" {
			buf.Write([]byte("\033[38;5;90m" + attr.Key + "--> " + attr.Value.String() + "\033[0m" + "\n")) // purple color
		} else {
			buf.Write([]byte(attr.Key + "--> " + attr.Value.String() + "\n"))
		}
	}

	r.Attrs(func(attr slog.Attr) bool {
		if attr.Key == "error" || attr.Key == "err" {
			buf.Write([]byte("\033[31m" + attr.Key + "--> " + attr.Value.String() + "\033[0m" + "\n")) // red color
		} else {
			buf.Write([]byte(attr.Key + "--> " + attr.Value.String() + "\n"))
		}
		return true
	})

	buf.Write([]byte("\n"))

	h.output.Write(buf.Bytes())
	return nil
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomHandler{
		handler: h.handler.WithAttrs(attrs),
		output:  h.output,
		attrs:   attrs,
	}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return &CustomHandler{
		handler: h.handler.WithGroup(name),
		output:  h.output,
	}
}

func InitLog(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(NewCustomHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))
	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
