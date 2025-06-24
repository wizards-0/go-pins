#!/bin/sh
go test ./... -coverprofile=out/coverage-all
cat out/coverage-all | grep -v "mocks" > out/coverage
go tool cover -html=out/coverage