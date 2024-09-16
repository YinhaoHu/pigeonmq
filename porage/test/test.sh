#!/bin/bash

# This script runs all tests in the e2e-test and integration-test directories.
# And then merges the coverage profiles and displays the combined coverage report.
#
# This script should be run from the `porage/test` directory.


# Run tests in e2e-test directory
echo "Running tests in e2e-test..."
cd ./e2e-test
go test -coverpkg=../../internal/... -coverprofile=./coverage/coverage.out 
cd ..

# Run tests in integration-test directory
echo "Running tests in integration-test..."
cd ./integration-test
PORAGE_LOG_LEVEL=warn go test -coverpkg=../../internal/... -coverprofile=./coverage/coverage.out
cd ..

# Merge coverage profiles
echo "Merging coverage profiles..."
gocovmerge ./e2e-test/coverage/coverage.out  ./integration-test/coverage/coverage.out  > ./coverage/coverage.out

# Show the coverage report
echo "Displaying combined coverage report..."
cd coverage
go tool cover -func=coverage.out | grep "total:" --color=never
go tool cover -html=coverage.out -o coverage.html
cd ..