# gfx

## Run test suite
```bash
go test
```
See `main_test.go` for details as to which files will be created in `test_data/`.

## Test Filter app
```bash
./test-filter-app.sh
```

Should return a list like this (order might be different):
```
Available filters
-----------------
sharpen(amount=0)
blur(amount=1)
edge-detect(amount=1)
invert()
brightness(adjustment=1)
contrast(adjustment=1)
vibrance(adjustment=0)
color-shift(hue=0 sat=0 lum=0)
gamma(adjustment=1)
threshold(amount=0)
alpha-map(source=l lower=0 upper=0)
sepia()
lum(shift=0)
hue-contrast(adjustment=0)
sat-contrast(adjustment=0)
convolution(amount=1 bias=0 factor=1 matrix=[[1 1 1] [1 8 1] [1 1 1]])
pastelize()
hue(shift=0)
emboss(amount=1)
extract(hue=0 hue-tolerance=180 hue-feather=0 sat=0.5 sat-tolerance=0.5 sat-feather=0 lum=0.5 lum-tolerance=0.5 lum-feather=0)
gray()
sat(shift=0)
lum-contrast(adjustment=0)
enhance(amount=1)
```
After the test `test_data/filter_app/` must contain `test1.png`, `test2.png`, `test3.png`. 

## Test Composer app
```bash
./test-composer-app.sh
```
Once complete, `test_data/composer_app/` must contain `sun.png`, `sun_spots.png`, `lasco_c3.png`. 

