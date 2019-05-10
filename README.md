# ylang
An image manipulation language

ylang is able to easily express these algorithms:
* Convolution
* Median
* Color distribution
* Edge detection
* Hough transform

## Usage
```
./ylang -code script.ylang -image image.jpg -out out.png
```
* `script.ylang` contains the ylang code to execute against the passed image
* `image.jpg` is the input image
* `out.png` is the output image

## Samples

This is the original image:

![original](https://raw.githubusercontent.com/smackem/ylang/master/doc/fish.jpg "Original")

### Greyscale

```
for point in Bounds {
    color := @point
    @point = rgb(color.i)
}
```
This is a complete ylang script. It iterates over all points contained in the rectangle `Bounds`, where `Bounds` is a constant that holds the dimensions of the input image.

`color := @point` assigns the color at point `point` to the new variable `color`.

`@point = rgb(color.i)` creates an rgb color with all three channels set to the intensity of the read color, then sets the pixel at position `point` to this color.

![greyscale](https://raw.githubusercontent.com/smackem/ylang/master/doc/greyscale.png "Image")

### Monochrome

```
for p in Bounds {
    @p = @p.i01 * #4080ff
}
```

The single statement that is executed for each pixel in the source image is `@p = @p.i01 * #4080ff`.

`@p` takes the color at point p, calculates the intensity (normalized to 0 .. 1) and multiplies the color `#4080ff` with the intensity. The result of this multiplication is a color with all color channels (r, g, b) multiplied by the intensity value.

![monochrome](https://raw.githubusercontent.com/smackem/ylang/master/doc/monochrome.png "Image")

### Saturate

```
for p in Bounds {
    @p = @p | hsv($) | hsv($.h, $.s * 1.5, $.v) | rgb($)
}
```

Like bash, powershell or F#, ylang supports pipelining.
`n := 1 | $ + 5 | $ * 2` reads like "takes the value 1, add 5, multiply with 2 and finally store the result in the newly declared variable n.

The statement is equivalent to `n := (1 + 5) * 2`.

The sample takes the input pixel color, converts it to hsv, creates a new hsv color with 50% more saturation, converts it to rgb and stores it in the target image.

![saturate](https://raw.githubusercontent.com/smackem/ylang/master/doc/saturate.png "Image")

### Convolution

Convolution with a HPF (high-pass filter) kernel emphasizes edges:

```
Kernel := |-1 -1 -1
           -1  8 -1
           -1 -1 -1|

for p in Bounds {
    @p = convolute(p, Kernel)
}
```

ylang supports kernel (a two-dimensional array of numbers) as a data type with concise and natural syntax. The built-in function `convolute` takes a point and a kernel as input and returns the resulting color.

The convolution is applied to all color channels (r, g, b).

To find out more about convolution and kernels, visit
https://en.wikipedia.org/wiki/Kernel_(image_processing)

![edges](https://raw.githubusercontent.com/smackem/ylang/master/doc/edges.png "Image")

Running the same program with a LPF (low-pass filter) kernel blurs the image:

```
Kernel := |0  1  2  1  0
           1  2  4  2  1
           2  4  8  4  2
           1  2  4  2  1
           0  1  2  1  0|

for p in Bounds {
    @p = convolute(p, Kernel)
}
```

![blur](https://raw.githubusercontent.com/smackem/ylang/master/doc/blur.png "Image")

### Slightly More Complex Sample

```
Median := kernel(5, 5, 1)
Center := 2;2

for p in Bounds {
    r := fetchRed(p, Median) | sort($)[Center]
    g := fetchGreen(p, Median) | sort($)[Center]
    b := fetchBlue(p, Median) | sort($)[Center]
    @p = rgb(r, g, b)
}

flip() // swap target and source image, committing all changes

SobelX := |-1  0  1
           -2  0  2
           -1  0  1|
SobelY := | 1  2  1
            0  0  0
           -1 -2 -1|

for p in Bounds {
    gx := convolute(p, SobelX).i
    gy := convolute(p, SobelY).i
    mag := hypot(gx, gy)
    angle := gx == 0 ? 0 : atan(gy / gx) * Rad2Deg + 90
    mag = mag / 255
    @p := round(angle * 2) | hsv($, mag, mag) | rgb($)
}
```

This sample first applies the median filter to the source image, then applies the sobel operator to the resulting image. The sobel operator yields gradient magnitude and gradient angle of edges, both of which are encoded into a the target pixel - the angle as the hue and the magnitude as value and saturation.

![complex](https://raw.githubusercontent.com/smackem/ylang/master/doc/complex.png "Image")

## Walkthrough

### Basics - Primitive Types

* Number
* Boolean
* String

### Basics - Control Flow

* if
* ternary operator
* for
* while
* return

### More Types

* Color
* Point
* Kernel
* Rectangle `rect(100, 100, 20, 50)`
* Line `line(100;100, 200;250)`
* Polygon `polygon(100;100, 300;200, 150;300)`
* Circle `circle(100;100, 50)`

You can iterate over geometrical shapes.
Plotting a red circle:
```
for p in circle(100;100, 50) {
    @p = #ff0000
}
```

### Math Functions

* ...

### Lists

* ...

### Objects

* ...

### Functions and Lambdas

* ...

## Roadmap

* Web interface with monaco as editor
* Compile to intermediate language -> ByteCode
* Compile to JavaScript, maybe WASM?
* Implement canny
  https://www.codeproject.com/kb/cs/canny_edge_detection.aspx
* Implement Rectange detection through hough transform
