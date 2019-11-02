function Color(r, g, b, a) {
    const self = this;
    self.r = () => r;
    self.g = () => g;
    self.b = () => b;
    self.a = () => a;
}

function Surface(canvas, width, height) {
    const self = this;
    let context = canvas.getContext('2d');
    let imgdata = context.getImageData(0, 0, width, height);
    let pixels = imgdata.data;

    self.getPixel = function(x, y) {
        let i = y * width + x;
        return new Color(pixels[i], pixels[i + 1], pixels[i + 2], pixels[i + 3]);
    }

    self.setPixel = function(x, y, color) {
        let i = y * width + x;
        pixels[i + 0] = color.r()
        pixels[i + 1] = color.g()
        pixels[i + 2] = color.b()
        pixels[i + 3] = color.a()
    }
}
