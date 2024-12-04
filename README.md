# gfx

```bash
go run app/filter/main.go filters
```

```
Available filters
-----------------
alpha-map::source=s*l::lower=0.1::upper=0.7
blur::amount=1.0
brightness::amount=1.0
color-shift::hue=180.0::sat=0.1::lum=0.7
contrast::amount=1.0
edge-detect::amount=1.0
emboss::amount=1.0
enhance::amount=1.0
extract::hue=180.0::hue-tolerance=90.0::hue-feather=90.0::sat=0.50::sat-tolerance=0.25::sat-feather=0.25::lum=0.50::lum-tolerance=0.25::lum-feather=0.25
gamma::amount=1.0
grayscale
hue-rotate::amount=180.0
invert
lightness-contrast::amount=1.0
pastelize
saturation::amount=1.0
sepia
sharpen::amount=1.0
threshold::amount=1.0
vibrance::amount=1.0
```