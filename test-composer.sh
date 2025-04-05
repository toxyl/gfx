#!/bin/bash
go generate build.go
code --install-extension docs/gfxs-syntax.vsix
# go run app/composer/main.go test.gfxs test.jpg image.png
go run app/composer/main.go test.gfxs test.jpg image.jpg
# go run app/composer/main.go -in test_data/compositions/sun_spots.gfxs -out test_data/composer_app/sun_spots.png 
# go run app/composer/main.go -in test_data/compositions/lasco_c3.gfxs -out test_data/composer_app/lasco_c3.png 
