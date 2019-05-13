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

It takes the color at point p, calculates the intensity (normalized to 0 .. 1) and multiplies the color `#4080ff` with the intensity. The result of this multiplication is a color with all color channels (r, g, b) multiplied by the intensity value.

![monochrome](https://raw.githubusercontent.com/smackem/ylang/master/doc/monochrome.png "Image")

### Saturate

```
for p in Bounds {
    @p = @p | hsv($) | hsv($.h, $.s * 1.5, $.v) | rgb($)
}
```

Like bash, powershell or F#, ylang supports pipelining.
`n := 1 | $ + 5 | $ * 2` reads like "take the value 1, add 5, multiply with 2 and finally store the result in the newly declared variable n.

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
// create a 5x5 kernel with all elements set to 1
Median := kernel(5, 5, 1)
Center := 2;2

// iterate over the image:
// fetch the red channel values of the 5x5 area around p, sort the resulting kernel
// and get the median value (center of kernel).
// then do the same for the green and blue channel.
for p in Bounds {
    r := fetchRed(p, Median) | sort($)[Center]
    g := fetchGreen(p, Median) | sort($)[Center]
    b := fetchBlue(p, Median) | sort($)[Center]
    @p = rgb(r, g, b)
}

// swap target and source image, committing all changes
flip()

SobelX := |-1  0  1
           -2  0  2
           -1  0  1|
SobelY := | 1  2  1
            0  0  0
           -1 -2 -1|

// iterate over the image:
// apply the sobel operator in both directions and calculate the gradient's magnitude.
// hypot is the same as writing sqrt(gx*gx + gy*gy).
// then calculate the gradient's angle, normalize the magnitude to 0..1 and the angle
// to 0..360 and write the resulting hsv color to the target image.
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

ylang is a dynamic script language featuring built-in types like points, kernels and colors - all of which are needed for image processing.

ylang's syntax is inspired by Go, Javascript, Bash and F#.

### Basics - Primitive Types

* All numbers in ylang are in 32 bit floating point format. The basic mathematical operations are supported on numbers. Constants can be written as `123` or `51.25`.
* Booleans have either the value `true` or `false`. All comparison operators like `==` or `<` return a boolean value.
* Strings are usually only used as hash map keys or for logging. String constants are written like this: `"Hello, world!"`.

### Basics - Variables and Constants

* Variables are declared with the `:=` operator, like in Go:
  ```
  num := 1
  str := "hello"
  ```
* Variables can be mutated with the `=` operator:
  ```
  num := 1
  log(num) // prints 1
  num = 2
  log(num) // prints 2
  ```
  Because of the dynamic type system, variables can also change type:
  ```
  v := 100 // initially a number
  log(v)
  v = "hello" // now a string
  log(v)
  ```
* Identifiers that start with a capital letter can be assigned only once:
  ```
  Ratio := 0.5
  Ratio = 1 // compilation error - Ratio is constant
  ```

### Basics - Control Flow

* ylang's `if else` statement is used like this:
  ```
  if @p.alpha > 100 {
      @p = #ff0000
  } else {
      @p = #00ff00
  }
  ```
* A shorthand for the previous `if else` statement is the ternary operator:
  ```
  @p = @p.alpha > 100 ? #ff0000 : #00ff00
  ```
* the `if else` statement can include any number of cases:
  ```
  a := @p.alpha
  if a > 200 {
      @p = #ff0000
  } else if a > 150 {
      @p = #800000
  } else if a > 100 {
      @p = #00ff00
  } else if a > 50 {
      @p = #008000
  } else {
      @p = #0000ff
  }
  ```
* The for loop iterates over ranges of numbers:
  ```
  // draw the top line blue
  for x in 0 .. Bounds.width {
      @(x;0) = #0000ff
  }
  ```
  or:
  ```
  // draw every second pixel in the top line blue
  for x in 0 .. 2 .. Bounds.width {
      @(x;0) = #0000ff
  }
  ```
* As seen before, the for loop can also iterate over iterable objects, like lists, kernels or geometrical shapes:
  ```
  for n in ["c", 0, 0, "l"] {
      log(n)
  }
  ```
* Another, much less common kind of loop is the while loop:
  ```
  x := 100
  y := 1
  while x > y {
      x = x / 2
      y = y * 2
  }
  ```
* To exit from a script or to return a value from a function, you can use the `return` statement:
  ```
  for p in Bounds {
      if p.y > 100 {
          return nil
      }
      @p = #00ff00
  }
  ```
  When breaking execution of the script, you can only return the empty value `nil`.

### Colors

Color literals are written as in HTML: `#ff0000` is the color red. You can append the alpha channel like this: `#00ff00:80` is half-transparent green.

You can also create color values with the functions `rgb`, `rgba`, `rgb01` or `rgba01`:
```
gold := rgb(255, 190, 0)
grey := rgb(128)
halfOpaqueBlue := rgba(0, 0, 255, 127)
white := rgb01(1, 1, 1)
halfOpaqueRed := rgba01(1, 0, 0, 0.5)
```

The color type defines these properties:
```
color := rgba(255, 128, 64, 32)
red := color.red // or color.r -- red is 255
green := color.green // or color.g -- green is 128
blue := color.blue // or color.b -- blue is 64
alpha := color.alpha // or color.a -- alpha is 32
r01 := color.red01 // or color.r01 -- r01 is 1
g01 := color.green01 // or color.g01 -- g01 is 0.5
b01 := color.blue01 // or color.b01 -- b01 is 0.25
a01 := color.alpha01 // or color.a01 -- a01 is 0.125
i := color.intensity // or color.i -- the intensity (brightness) of the color
i01 := color.intensity // or color.i01 -- the intensity normalized to 0..1
```

The color type supports basic arithmetic operations, which are applied per channel:
```
rgb(100, 110, 120) + rgb(10, 20, 30) // = rgb(110, 130, 150)
#ff0080 * 0.5 // = #800040
rgb(100, 200, 300) * 2 // = rgb(200, 400, 600)
```

Colors can take any value and are only clamped to 0..255 when necessary, e.g. when writing the color to the target image.

### Point

`x;y` denotes a point. `x` and `y` are implicitly converted to integer values.

Points have the following properties:
```
p := 100;120
x := p.x // = 100
y := p.y // = 120
mag := p.mag // magnitude of the point interpreted as a vector
```

### Kernel

Kernels can be cresated as literals:
```
k := |0 1 0
      1 2 1
      0 1 0|
```

... or by using the `kernel` function:

```
k := kernel(3, 3, 1) // 3x3 with all elements set to 1
k := kernel(4, 2, fn(x, y) -> x + y)
// = |0 1 2 3
//    1 2 3 4|
```

Kernels can be indexed with numbers or points:
```
k := |0 1 0
      1 2 1
      0 1 0|
n := k[4] // = 2
m := k[1;1] // = 2
```

Kernels have these properties:
```
k := kernel(2, 3, 0)
width := k.width // = 2
height := k.height // = 3
count := k.count // = 6 -- number of elements
```

### Rectangle
`rect(100, 100, 20, 50)`

### Line
`line(100;100, 200;250)`

### Polygon
`polygon(100;100, 300;200, 150;300)`

### Circle
`circle(100;100, 50)`

You can iterate over geometrical shapes.
Plotting a red circle:
```
for p in circle(100;100, 50) {
    @p = #ff0000
}
```

### Working with images

* ...

### Math Functions

* ...

### Alpha Channel

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
