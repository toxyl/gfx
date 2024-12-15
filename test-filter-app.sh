#!/bin/bash
go run app/filter/main.go -list
go run app/filter/main.go -in test_data/test1.png -out test_data/filter_app/test1.png -f 'enhance()' -f 'hue(shift=90)' -f 'emboss(amount=1)' 
go run app/filter/main.go -in test_data/test1.png -out test_data/filter_app/test2.png -f 'enhance()' -f 'hue(shift=270)' -f 'blur(amount=1)'
go run app/filter/main.go -in test_data/test1.png -out test_data/filter_app/test3.png -chain test_data/filter_app/test1.gfxs
