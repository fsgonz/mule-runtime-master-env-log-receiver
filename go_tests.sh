#!/bin/bash

# Run Go tests
go install github.com/jstemmer/go-junit-report/v2@latest
go test -v 2>&1 ./... ./... | go-junit-report -set-exit-code > /test-reports/report.xml

echo "Go module tests run!"