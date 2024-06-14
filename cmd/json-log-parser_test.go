package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var logWithException = map[string]interface{}{
	"level":     "DEBUG",
	"timestamp": "2024-06-10T08:47:35.469Z",
	"thread":    "Thread-1",
	"logger":    "Logger-1",
	"message":   "Something bad happened when processing this request - see exception for details",
	"exception": "java.lang.NullPointerException: null\\n\\tat java.base/java.util.Objects.requireNonNull(Unknown Source)\\n\\tat redacted",
}

var logWithoutException = map[string]interface{}{
	"level":     "INFO",
	"timestamp": "2024-06-10T08:47:35.444Z",
	"thread":    "Thread-4",
	"logger":    "Logger-3",
	"message":   "Something of note happened",
}

var builder = &strings.Builder{}

func BenchmarkLogPrintingWithBuilder(b *testing.B) {
	for n := 0; n < b.N; n++ {
		formatLog(builder, logWithException)
		builder.Reset()
	}
}

func TestParseSeverity(t *testing.T) {
	var scenarios = map[string]Severity{
		"":            unknown,
		"sdfdfsdfsfd": unknown,
		"unknown":     unknown,
		"trace":       trace,
		"TRACE":       trace,
		"debug":       debug,
		"DEBUG":       debug,
		"info":        info,
		"INFO":        info,
		"warn":        warn,
		"WARN":        warn,
		"error":       err,
		"ERROR":       err,
	}

	for desc := range scenarios {
		assert.Equal(t, scenarios[desc], parseSeverity(desc))
	}

}

func TestFormatLog(t *testing.T) {
	var b = &strings.Builder{}
	formatLog(b, logWithException)

	assert.Equal(t, "2024-06-10T08:47:35.469Z [Thread-1] DEBUG Logger-1 "+
		"Something bad happened when processing this request - see exception for details "+
		"java.lang.NullPointerException: null\\n\\tat java.base/java.util.Objects.requireNonNull(Unknown Source)\\n"+
		"\\tat redacted", b.String())
}

func TestFormatLogNoException(t *testing.T) {
	var b = &strings.Builder{}
	formatLog(b, logWithoutException)

	assert.Equal(t, "2024-06-10T08:47:35.444Z [Thread-4] INFO Logger-3 "+
		"Something of note happened", b.String())
}
