# Filters
## alpha-map
Creates an alpha channel based on a source channel


```gfxs
alpha-map(source="l" lower=0.000000 upper=1.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| source| string| | `"l"`| | | The source channel (`l` = luminance, `s` = saturation, `s*l` = saturation * luminance)|
| lower| float| %| `0.000000`| `0`| `1`| The lower threshold, everything below will be treated as fully transparent|
| upper| float| %| `1.000000`| `0`| `1`| The upper threshold, everything above will be treated as fully opaque|


## blur
Applies a blur using a convolution matrix.


```gfxs
blur(amount=1.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| amount| float| px| `1.000000`| `-10`| `10`| The amount to blur in pixels|


## brightness
Adjusts the brightness of the image.


```gfxs
brightness(adjustment=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| adjustment| float| %| `0.000000`| `-1`| `1`| Adjustment amount where 0 means no adjustment|


## color-shift
Combines hue-shit, sat-shift and lum-shift into a single filter.


```gfxs
color-shift(hue=0.000000 sat=0.000000 lum=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| hue| float| °| `0.000000`| `-360`| `360`| |
| sat| float| %| `0.000000`| `-1`| `1`| |
| lum| float| %| `0.000000`| `-1`| `1`| |


## contrast
Adjusts the image contrast.


```gfxs
contrast(adjustment=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| adjustment| float| %| `0.000000`| `-1`| `1`| |


## convolution
Applies a 3x3 convolution matrix to the image.


```gfxs
convolution(amount=1.000000 bias=0.000000 factor=1.000000 matrix=[1 1 1] [1 8 1] [1 1 1])
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| amount| float| %| `1.000000`| `-Inf`| `+Inf`| |
| bias| float| %| `0.000000`| `-Inf`| `+Inf`| |
| factor| float| %| `1.000000`| `-Inf`| `+Inf`| |
| matrix| matrix| | `[1 1 1] [1 8 1] [1 1 1]`| | | |


## crop
Crops the layer.


```gfxs
crop(left=0.000000 right=0.000000 top=0.000000 bottom=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| left| float| %| `0.000000`| `0`| `1`| |
| right| float| %| `0.000000`| `0`| `1`| |
| top| float| %| `0.000000`| `0`| `1`| |
| bottom| float| %| `0.000000`| `0`| `1`| |


## crop-circle
Crops the layer using a circle.


```gfxs
crop-circle(radius=0.500000 offset-x=0.000000 offset-y=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| radius| float| %| `0.500000`| `0`| `+Inf`| |
| offset-x| float| %| `0.000000`| `-1`| `1`| |
| offset-y| float| %| `0.000000`| `-1`| `1`| |


## edge-detect
Detects edges.


```gfxs
edge-detect(amount=1.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| amount| float| %| `1.000000`| `-Inf`| `+Inf`| |


## emboss
Embosses the image.


```gfxs
emboss(amount=1.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| amount| float| %| `1.000000`| `-Inf`| `+Inf`| |


## enhance
Enhances the image.


```gfxs
enhance(amount=1.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| amount| float| %| `1.000000`| `-Inf`| `+Inf`| |


## extract
Extracts a color from the image.


```gfxs
extract(hue=0.000000 hue-tolerance=180.000000 hue-feather=0.000000 sat=0.500000 sat-tolerance=0.500000 sat-feather=0.000000 lum=0.500000 lum-tolerance=0.500000 lum-feather=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| hue| float| °| `0.000000`| `0`| `360`| |
| hue-tolerance| float| °| `180.000000`| `0`| `360`| |
| hue-feather| float| °| `0.000000`| `0`| `360`| |
| sat| float| %| `0.500000`| `0`| `1`| |
| sat-tolerance| float| %| `0.500000`| `0`| `1`| |
| sat-feather| float| %| `0.000000`| `0`| `1`| |
| lum| float| %| `0.500000`| `0`| `1`| |
| lum-tolerance| float| %| `0.500000`| `0`| `1`| |
| lum-feather| float| %| `0.000000`| `0`| `1`| |


## flip-h
Flips the image horizontally.


```gfxs
flip-h()
```


## flip-v
Flips the image vertically.


```gfxs
flip-v()
```


## gamma
Adjusts the image gamma.


```gfxs
gamma(adjustment=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| adjustment| float| %| `0.000000`| `-1`| `1`| |


## gray
Converts the image to grayscale


```gfxs
gray()
```


## hue
Shifts the hues of the image.


```gfxs
hue(shift=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| shift| float| °| `0.000000`| `-360`| `360`| |


## hue-contrast
Adjusts the contrast of hues in the image.


```gfxs
hue-contrast(adjustment=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| adjustment| float| %| `0.000000`| `-1`| `1`| |


## invert
Inverts all colors.


```gfxs
invert()
```


## lum
Shifts the luminosities of the image.


```gfxs
lum(shift=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| shift| float| %| `0.000000`| `-1`| `1`| |


## lum-contrast
Adjusts the contrast of luminosities in the image.


```gfxs
lum-contrast(adjustment=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| adjustment| float| %| `0.000000`| `-1`| `1`| |


## pastelize
Converts the image to pastel colors


```gfxs
pastelize()
```


## project-as-equirectangular
Converts the image to a equirectangular projection.


```gfxs
project-as-equirectangular()
```


## project-as-mercator
Converts the image to a mercator projection.


```gfxs
project-as-mercator()
```


## project-as-polar-to-rect
Converts the image to a polar-to-rect projection.


```gfxs
project-as-polar-to-rect()
```


## project-as-rect-to-polar
Converts the image to a rect-to-polar projection.


```gfxs
project-as-rect-to-polar()
```


## project-as-sinusoidal
Converts the image to a sinusoidal projection.


```gfxs
project-as-sinusoidal()
```


## project-as-stereographic
Converts the image to a stereographic projection.


```gfxs
project-as-stereographic()
```


## remove
Removes a color from the image.


```gfxs
remove(hue=0.000000 range=15.000000 feather=15.000000 amount=1.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| hue| float| °| `0.000000`| `0`| `360`| |
| range| float| °| `15.000000`| `0`| `360`| |
| feather| float| °| `15.000000`| `0`| `360`| |
| amount| float| %| `1.000000`| `0`| `1`| |


## rotate
Rotates the image.


```gfxs
rotate(angle=0.000000 offset-x=0.000000 offset-y=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| angle| float| °| `0.000000`| `-360`| `360`| |
| offset-x| float| %| `0.000000`| `-1`| `1`| Offset on the x-axis relative to the center.|
| offset-y| float| %| `0.000000`| `-1`| `1`| Offset on the y-axis relative to the center.|


## sat
Shifts the saturations of the image.


```gfxs
sat(shift=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| shift| float| %| `0.000000`| `-1`| `1`| |


## sat-contrast
Adjusts the contrast of saturations in the image.


```gfxs
sat-contrast(adjustment=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| adjustment| float| %| `0.000000`| `-1`| `1`| |


## scale
Scales the image.


```gfxs
scale(factor=0.000000 offset-x=0.000000 offset-y=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| factor| float| %| `0.000000`| `-Inf`| `+Inf`| Scale factor where 0 means no scaling, >0 upscaling and <0 downscaling|
| offset-x| float| %| `0.000000`| `-1`| `1`| Offset on the x-axis relative to the center.|
| offset-y| float| %| `0.000000`| `-1`| `1`| Offset on the y-axis relative to the center.|


## sepia
Converts the image to sepia colors.


```gfxs
sepia()
```


## sharpen
Sharpens the image.


```gfxs
sharpen(amount=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| amount| float| %| `0.000000`| `-1`| `1`| |


## threshold
Colors pixels above the threshold white and pixels below the threshold black.


```gfxs
threshold(threshold=0.500000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| threshold| float| %| `0.500000`| `0`| `1`| |


## to-polar
Converts rectangular coordinates to polar coordinates. For example useful to transform a 360 degrees camera image to an all-sky camera.


```gfxs
to-polar(angle-start=0.000000 angle-end=360.000000 rotation=0.000000 fisheye=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| angle-start| float| °| `0.000000`| `0`| `360`| |
| angle-end| float| °| `360.000000`| `0`| `360`| |
| rotation| float| °| `0.000000`| `-360`| `360`| |
| fisheye| float| %| `0.000000`| `-1`| `1`| |


## transform
Combines translate, rotate and scale into a single filter.


```gfxs
transform(translate-x=0.000000 translate-y=0.000000 rotate=0.000000 scale=0.000000 offset-x=0.000000 offset-y=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| translate-x| float| %| `0.000000`| `-1`| `1`| Movement on the x-axis|
| translate-y| float| %| `0.000000`| `-1`| `1`| Movement on the y-axis|
| rotate| float| °| `0.000000`| `-360`| `360`| Rotation|
| scale| float| %| `0.000000`| `-Inf`| `+Inf`| Scale factor where 0 means no scaling, >0 upscaling and <0 downscaling|
| offset-x| float| %| `0.000000`| `-1`| `1`| Offset on the x-axis, relative to the center.|
| offset-y| float| %| `0.000000`| `-1`| `1`| Offset on the y-axis, relative to the center.|


## translate
Moves the layer.


```gfxs
translate(x=0.000000 y=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| x| float| %| `0.000000`| `-1`| `1`| Movement on the x-axis|
| y| float| %| `0.000000`| `-1`| `1`| Movement on the y-axis|


## translate-wrap
Moves the layer and wraps pixels around.


```gfxs
translate-wrap(x=0.000000 y=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| x| float| %| `0.000000`| `-1`| `1`| Movement on the x-axis|
| y| float| %| `0.000000`| `-1`| `1`| Movement on the y-axis|


## vibrance
Adjusts the vibrance of the image.


```gfxs
vibrance(adjustment=0.000000)
```
| Argument | Type | Unit | Default | Min | Max | Description |
| --- | --- | --- | --- | --- | --- | --- |
| adjustment| float| %| `0.000000`| `-1`| `1`| |


