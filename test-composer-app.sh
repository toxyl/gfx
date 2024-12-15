#!/bin/bash
go run app/composer/main.go -in test_data/compositions/sun.gfxs -out test_data/composer_app/sun.png 
go run app/composer/main.go -in test_data/compositions/sun_spots.gfxs -out test_data/composer_app/sun_spots.png 
go run app/composer/main.go -in test_data/compositions/lasco_c3.gfxs -out test_data/composer_app/lasco_c3.png 
