package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
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
		return
	}

	reader := createInputFile()
	defer reader.Close()
	writer := createOutputFile()
	defer writer.Close()

	dec := json.NewDecoder(reader)
	logBuilder := &strings.Builder{}
	for {
		var v map[string]interface{}
		if err := dec.Decode(&v); err != nil {
			if err != io.EOF {
				log.Fatalln("Failed to process json", err)
			}
			return
		}
		if matchesFilter(v) {
			formatLog(logBuilder, v)
			fmt.Fprintln(writer, logBuilder.String())
			logBuilder.Reset()
		}
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
	severity := parseSeverity(json["level"].(string))
	return severity == unknown || severity.weight >= minSeverity.weight
}

func formatLog(builder *strings.Builder, v map[string]interface{}) {
	fmt.Fprint(builder, v["timestamp"], " [", v["thread"], "] ", v["level"], " ", v["logger"], " ", v["message"])
	if ex, ok := v["exception"]; ok {
		fmt.Fprint(builder, " ", ex)
	}
}
