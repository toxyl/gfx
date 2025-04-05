#!/bin/bash
go generate build.go
code --install-extension docs/gfxs-syntax.vsix
go run -C app/docs/ . $@
