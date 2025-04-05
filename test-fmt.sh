#!/bin/bash
go generate build.go
code --install-extension docs/gfxs-syntax.vsix
go run app/fmt/main.go test.gfxs test-parsed.gfxs 
# go run app/fmt/main.go test.gfxs test-parsed-1.gfxs 
# go run app/fmt/main.go test-parsed-1.gfxs test-parsed-2.gfxs 
# go run app/fmt/main.go test-parsed-2.gfxs test-parsed-3.gfxs 