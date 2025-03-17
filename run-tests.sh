#!/bin/sh
go test ./... -coverprofile=out/coverage
go tool cover -html=out/coverage