package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// Severity represents how important a particular log is
type Severity struct {
	description string
	weight      int
}

func (r Severity) String() string {
	return r.description
}

var (
	unknown    = Severity{"unknown", 99}
	trace      = Severity{"trace", 0}
	debug      = Severity{"debug", 1}
	info       = Severity{"info", 2}
	warn       = Severity{"warn", 3}
	err        = Severity{"error", 4}
	severities = []Severity{trace, debug, info, warn, err}
)

var (
	levelFlag = flag.String("severity", "info", "min severity level to filter on, one of: "+fmt.Sprintf("%v", severities))
	fileFlag  = flag.String("file", "", "json log file to read from - if not supplied then stdin is used")
	outFlag   = flag.String("out", "", "output file to write the parsed logs to - if not supplied then stdout is used")
)

var minSeverity = info

func main() {
	flag.Parse()

	minSeverity = parseSeverity(*levelFlag)
	if minSeverity == unknown {
		fmt.Println("Invalid level value supplied.")
		flag.Usage()
		os.Exit(1)
	}

	reader := createInputFile()
	defer reader.Close()
	writer := createOutputFile()
	defer writer.Close()

	logBuilder := &strings.Builder{}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var v map[string]interface{}
		var dec = json.NewDecoder(strings.NewReader(scanner.Text()))
		if err := dec.Decode(&v); err != nil {
			// Skip decode errors - these indicate that lines are not valid json (logs may contain noise)
			continue
		}
		if matchesFilter(v) {
			formatLog(logBuilder, v)
			fmt.Fprintln(writer, logBuilder.String())
			logBuilder.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func createInputFile() *os.File {
	if *fileFlag == "" || *fileFlag == "-" {
		return os.Stdin
	} else {
		file, err := os.Open(*fileFlag)
		if err != nil {
			fmt.Println("Failed to read input file due to err", err)
			os.Exit(1)
		}
		return file
	}
}

func createOutputFile() *os.File {
	if *outFlag == "" {
		return os.Stdout
	} else {
		file, err := os.Create(*outFlag)
		if err != nil {
			fmt.Println("Failed to create output file due to err", err)
			os.Exit(1)
		}
		return file
	}
}

func parseSeverity(description string) Severity {
	description = strings.ToLower(description)
	for _, l := range severities {
		if l.description == description {
			return l
		}
	}
	return unknown
}

func matchesFilter(json map[string]interface{}) bool {
	return looksLikeLogMessageFilter(json) && matchesSeverityFilter(json)
}

func matchesSeverityFilter(json map[string]interface{}) bool {
	severity := parseSeverity(json["level"].(string))
	return severity == unknown || severity.weight >= minSeverity.weight
}

var logMessageMandatoryFields = []string{"timestamp", "level", "message"}

// looksLikeLogMessageFilter
// This filter sanity checks that the parsed JSON looks like a log message, this is useful when logs contain raw JSON
func looksLikeLogMessageFilter(json map[string]interface{}) bool {
	for i := range logMessageMandatoryFields {
		if _, in := json[logMessageMandatoryFields[i]]; !in {
			return false
		}
	}
	return true
}

func formatLog(builder *strings.Builder, v map[string]interface{}) {
	fmt.Fprint(builder, v["timestamp"], " [", v["thread"], "] ", v["level"], " ", v["logger"], " ", v["message"])
	if ex, ok := v["exception"]; ok {
		fmt.Fprint(builder, " ", ex)
	}
}
