# Graphics Library

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
alpha-map(source=l lower=0 upper=0)
blur(amount=1)
brightness(adjustment=1)
color-shift(hue=0 sat=0 lum=0)
contrast(adjustment=1)
convolution(amount=1 bias=0 factor=1 matrix=[[1 1 1] [1 8 1] [1 1 1]])
crop(left=0 right=0 top=0 bottom=0)
crop-circle(radius=0 offset-x=0 offset-y=0)
edge-detect(amount=1)
emboss(amount=1)
enhance(amount=1)
extract(hue=0 hue-tolerance=180 hue-feather=0 sat=0.5 sat-tolerance=0.5 sat-feather=0 lum=0.5 lum-tolerance=0.5 lum-feather=0)
flip-h()
flip-v()
gamma(adjustment=1)
gray()
hue-contrast(adjustment=0)
hue(shift=0)
invert()
lum-contrast(adjustment=0)
lum(shift=0)
pastelize()
rotate(angle=0 offset-x=0 offset-y=0)
sat-contrast(adjustment=0)
sat(shift=0)
scale(scale=0 offset-x=0 offset-y=0)
sepia()
sharpen(amount=0)
threshold(amount=0)
transform(transform-x=0 transform-y=0 rotate=0 scale=0 offset-x=0 offset-y=0)
translate(x=0 y=0)
translate-wrap(x=0 y=0)
vibrance(adjustment=0)
```
After the test `test_data/filter_app/` must contain `test1.png`, `test2.png`, `test3.png`. 

## Test Composer app
```bash
./test-composer-app.sh
```
Once complete, `test_data/composer_app/` must contain `sun.png`, `sun_spots.png`, `lasco_c3.png`. 

# GFXScript
The `composer` and `filter` apps used above make use of the `GFXScript` language to compose images / apply filters to images.  
For example: to render a sun image as used on https://aurora-map.toxyl.nl a script similar to this is used:
```
[VARS] # defines variables that can be used in the [FILTERS] section 
b1 = 4.0
b2 = 1.5
e1 = 1.025
s0 = 0
l0 = 0
h1 = -60
h2 = -150
h3 = -190
h4 = 5
s4 = -0.05
l4 = 0.01
src = `l`

[FILTERS] # defines filters that can be used in the [COMPOSITION] and [LAYERS] sections 
filter1    { blur(amount=b1) }
filter2    { enhance(e1) }
filter3    { 
  use(filter2) # use(...) includes a previously defined filter
  blur(amount=b2) 
}
filter4    { color-shift(hue=h1 sat = s0 lum=l0) enhance() }
filter5    { color-shift(hue=h2 sat = s0 lum=l0) enhance() }
filter6    { color-shift(hue=h3 sat = s0 lum=l0) enhance() }
filter7    { color-shift(h4 s4 l4) enhance() } # args can be unnamed but then their order must match as per "Available Filters" above
compFilter { enhance() alpha-map(lower=0.02 upper=0.275 source=src) }

[COMPOSITION] # describes general properties of the composition
name   = `Sun (GOES)`
width  = 300		
height = 300
color  = hsla(240 0.5 0.25 1.0) # color of the first layer (only unnamed args allowed)
filter = compFilter # name of filter to apply after rendering all layers, must be defined in [FILTERS] section
crop   = 20 20 260 260 # x, y, w, h
resize = 300 300 # w, h

[LAYERS] # all layers of the composition, read from bottom to top, like in photoshop
#     mode  alpha  filter source
 pin-light 0.5000 filter7 ./test_data/compositions/layers/goes_18_171.png # top-most layer
   average 0.5000 filter6 ./test_data/compositions/layers/goes_16_171.png
 pin-light 0.7500 filter5 ./test_data/compositions/layers/goes_18_131.png
    darken 0.5000 filter4 ./test_data/compositions/layers/goes_16_131.png
    darken 0.7500 filter2 ./test_data/compositions/layers/goes_18_195.png
difference 0.7500 filter2 ./test_data/compositions/layers/goes_16_195.png
   average 0.5000 filter3 ./test_data/compositions/layers/goes_18_094.png
   average 0.5000 filter2 ./test_data/compositions/layers/goes_16_094.png
difference 1.0000       * ./test_data/compositions/layers/goes_18_304.png # * means 'no filter'
difference 1.0000       * ./test_data/compositions/layers/goes_16_304.png
  multiply 1.0000 filter1 ./test_data/compositions/layers/goes_18_284.png
  multiply 1.0000       * ./test_data/compositions/layers/goes_16_284.png # bottom-most layer
```

## VSCode extension for syntax highlighting
```bash
./install-syntax-highlighter-vscode.sh
```
This will generate a VSCode extension (`*.vsix` archive) with the latest version of the language rules (see `parser/vsix_builder.go` for details) and then install it in VSCode. To refresh syntax highlighting after an update you have to press `CTRL+Shift+P` and run `Developer: Reload Window`. If you want to modify tokens used by the language, have a look at `parser/config.go`.  