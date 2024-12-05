# gfx

```bash
go run app/filter/main.go filters
```

```
Available filters
-----------------
gray
invert
pastelize
sepia
hue::amount=0.5
sat::amount=1.0
lum::amount=1.0
hue-contrast::amount=1.0
sat-contrast::amount=1.0
lum-contrast::amount=1.0
color-shift::hue=180.0::sat=0.1::lum=0.7
brightness::amount=1.0
contrast::amount=1.0
gamma::amount=1.0
vibrance::amount=1.0
enhance::amount=1.0
sharpen::amount=1.0
blur::amount=1.0
edge-detect::amount=1.0
emboss::amount=1.0
threshold::amount=1.0
alpha-map::source=s*l::lower=0.1::upper=0.7
extract::hue=180.0::hue-tolerance=90.0::hue-feather=90.0::sat=0.50::sat-tolerance=0.25::sat-feather=0.25::lum=0.50::lum-tolerance=0.25::lum-feather=0.25
convolution::matrix=1.0,1.0,1.0,1.0,8.0,1.0,1.0,1.0,1.0
```