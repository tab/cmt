package logger

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
)

// bufferedWriter intercepts zerolog JSON output and sends it to a buffer
type bufferedWriter struct {
	buffer        *LogBuffer
	consoleWriter zerolog.ConsoleWriter
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// Write implements io.Writer interface
func (w *bufferedWriter) Write(p []byte) (n int, err error) {
	var buf bytes.Buffer
	tempWriter := w.consoleWriter
	tempWriter.Out = &buf

	if _, err := tempWriter.Write(p); err != nil {
		return 0, err
	}

	formattedWithANSI := strings.TrimSpace(buf.String())
	plainText := stripANSI(formattedWithANSI)

	var entry map[string]interface{}
	level := "info"
	if err := json.Unmarshal(p, &entry); err == nil {
		if levelVal, ok := entry["level"].(string); ok {
			level = levelVal
		}
	}

	w.buffer.AddFormatted(level, plainText, formattedWithANSI)
	return len(p), nil
}

// stripANSI removes ANSI escape codes from a string
func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}
