# Overview
This tool is a command line app which parses json logs produced by a Java application and formats it as a plain text.

This tool is useful when investigating logs without requiring 3rd party tooling like (Google Cloud logs or Kibana). 
This is especially useful when trying to read/investigate exceptions as it formats the stacktraces over multiple lines.

# Installing
```bash
make install
```
This will install `json-log-parser` into the `~/go/bin/` directory.

# Command line args
```
Usage of json-log-parser:
  -file string
        json log file to read from - if not supplied then stdin is used
  -out string
        output file to write the parsed logs to - if not supplied then stdout is used
  -severity string
        min severity level to filter on, one of: [trace debug info warn error] (default "info")
```

# Example commands
These commands assume that json-log-parser has been installed and is available on the PATH.
## Format k8s logs
One scenario that this works particular well for, retrieving logs from a k8s pod using kubectl.
Example command:
```
kubectl logs POD-NAME | json-log-parser
```

# JSON log file structure
Example log object:
```json
 {
    "timestamp":"2024-06-10T08:46:25.518Z",
    "level":"ERROR",
    "thread":"thread-12",
    "logger":"com.blah.SomeClass",
    "message":"Something went wrong, see exception for details",
    "exception":"java.lang.NullPointerException: null\n\tat java.base/java.util.Objects.requireNonNull(Unknown Source)\n\tredacted..."
 }
 ```
A log file contains a log object per line and not an array of logs objects.

# Output format
Example formatted plain text log line output by the tool:
```
2024-06-10T08:46:25.518Z [thread-12] ERROR com.blah.SomeClass Something went wrong, see exception for details java.lang.NullPointerException: null
        at java.base/java.util.Objects.requireNonNull(Unknown Source)
        redacted...
```

# Future work
- support for logging java mdc (mapped diagnostic context
- support other languages / json log formats
- allow output format to be customised
- further log filtering
 