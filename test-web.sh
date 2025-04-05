#!/bin/bash
go generate build.go
go run app/web/main.go
