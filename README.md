# ylang
An image manipulation language

ylang must be able to implement these algorithms:
* Convolution
* Median
* Canny edge detection
  https://www.codeproject.com/kb/cs/canny_edge_detection.aspx
* Hough transform
* Color distribution

## Roadmap
* Web interface with monaco as editor
* Drawing language for creation of graphical images
* Compile to intermediate language -> ByteCode
* Use WASM? JavaScript first :)

## Usage
```
./ylang -code script.ylang -image image.jpg -out out.png
```
* `script.ylang` contains the ylang code to execute against the passed iamge
* `image.jpg` is the input image
* `out.png` is the output image

## Samples

### Greyscale Image

```
for point in Bounds {
    color := @point
    @point = rgb(color.i)
}
```
This is a complete ylang script. It iterates over all points contained in the rectangle `Bounds`, where `Bounds` is a constant that holds the dimensions of the input image.

`color := @point` assigns the color at point `point` to the new variable `color`.

`@point = rgb(color.i)` creates an rgb color with all three channels set to the intensity of the read color, then sets the pixel at position `point` to this color.