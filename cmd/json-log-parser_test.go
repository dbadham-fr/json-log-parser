package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"testing"
)

var jsonLogWithException = map[string]interface{}{
	"level":     "DEBUG",
	"timestamp": "2024-06-10T08:47:35.469Z",
	"thread":    "Thread-1",
	"logger":    "Logger-1",
	"message":   "Something bad happened when processing this request - see exception for details",
	"exception": "java.lang.NullPointerException: null\\n\\tat java.base/java.util.Objects.requireNonNull(Unknown Source)\\n\\tat redacted",
}

var plainTextLogWithException = "2024-06-10T08:47:35.469Z [Thread-1] DEBUG Logger-1 " +
	"Something bad happened when processing this request - see exception for details " +
	"java.lang.NullPointerException: null\\n\\tat java.base/java.util.Objects.requireNonNull(Unknown Source)\\n" +
	"\\tat redacted"

var jsonLogWithoutException = map[string]interface{}{
	"level":     "INFO",
	"timestamp": "2024-06-10T08:47:35.444Z",
	"thread":    "Thread-4",
	"logger":    "Logger-3",
	"message":   "Something of note happened",
}

var plainTextLogWithoutException = "2024-06-10T08:47:35.444Z [Thread-4] INFO Logger-3 Something of note happened"

var builder = &strings.Builder{}

func BenchmarkLogPrintingWithBuilder(b *testing.B) {
	for n := 0; n < b.N; n++ {
		formatLog(builder, jsonLogWithException)
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
	formatLog(b, jsonLogWithException)

	assert.Equal(t, plainTextLogWithException, b.String())
}

func TestFormatLogNoException(t *testing.T) {
	var b = &strings.Builder{}
	formatLog(b, jsonLogWithoutException)

	assert.Equal(t, plainTextLogWithoutException, b.String())
}

func TestParseValidLogs(t *testing.T) {
	reader := strings.NewReader(createStringBuilderContainingJson(jsonLogWithException, jsonLogWithoutException).String())
	writer := &bytes.Buffer{}

	parseLog(reader, writer)

	assert.NotEmpty(t, writer.String())
	assert.Equal(t, plainTextLogWithException+"\n"+plainTextLogWithoutException+"\n", writer.String())
}

func createStringBuilderContainingJson(values ...map[string]interface{}) *strings.Builder {
	sb := &strings.Builder{}
	encodeJson(sb, values...)
	return sb
}

func encodeJson(sb *strings.Builder, values ...map[string]interface{}) {
	for i := range values {
		encoder := json.NewEncoder(sb)
		v := values[i]
		err := encoder.Encode(v)
		if err != nil {
			log.Fatal("Unexpected error encoding test data as json", err)
		}
	}
}

func TestParseLogsContainingNoJson(t *testing.T) {
	sb := &strings.Builder{}
	sb.WriteString(plainTextLogWithException)
	sb.WriteString("\n")
	sb.WriteString(plainTextLogWithoutException)
	sb.WriteString("\n")
	sb.WriteString("Junk Written to stdout\n")
	sb.WriteString("Junk Written to stderr\n")
	sb.WriteString("<xml><blah>some text</blah></xml>\n")

	reader := strings.NewReader(sb.String())
	writer := &bytes.Buffer{}

	parseLog(reader, writer)

	assert.Empty(t, writer.String())
}

func TestParseJsonLogsContainingNoise(t *testing.T) {
	logs := createStringBuilderContainingJson(jsonLogWithException)
	logs.WriteString("Junk data in logs\n")
	logs.WriteString("Junk data in logs\n")
	encodeJson(logs, jsonLogWithoutException)

	// Add raw json to the logs which has not been written via a logger
	encodeJson(logs, map[string]interface{}{
		"response_time": 129,
		"response_code": "200",
		"endpoint":      "/admin/v1/user",
		"method":        "POST",
	})

	reader := strings.NewReader(logs.String())
	writer := &bytes.Buffer{}

	parseLog(reader, writer)

	assert.NotEmpty(t, writer.String())
	// Verify junk data is omitted
	assert.Equal(t, plainTextLogWithException+"\n"+plainTextLogWithoutException+"\n", writer.String())
}
