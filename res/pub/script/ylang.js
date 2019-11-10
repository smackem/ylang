function Color(r, g, b, a) {
    const self = this;
    self.r = r;
    self.g = g;
    self.b = b;
    self.a = a;

    function intensity() {
        return 0.299 * self.r + 0.587 * self.g + 0.114 * self.b;
    }

    Object.defineProperty(this, "i", {
        get() {
            return intensity();
        }
    });
    Object.defineProperty(this, "intensity", {
        get() {
            return intensity();
        }
    });
    Object.defineProperty(this, "i01", {
        get() {
            return intensity() / 255.0;
        }
    });
    Object.defineProperty(this, "intensity01", {
        get() {
            return intensity() / 255.0;
        }
    });
    Object.defineProperty(this, "r01", {
        get() {
            return self.r / 255.0;
        }
    });
    Object.defineProperty(this, "g01", {
        get() {
            return self.g / 255.0;
        }
    });
    Object.defineProperty(this, "b01", {
        get() {
            return self.b / 255.0;
        }
    });
    Object.defineProperty(this, "a01", {
        get() {
            return self.a / 255.0;
        }
    });

    self.neg = function () {
        return new Color(255 - self.r, 255 - self.g, 255 - self.b, self.a);
    };
}

function color01(r01, g01, b01, a01) {
    return new Color(r01 * 255.0, g01 * 255.0, b01 * 255.0, a01 * 255.0);
}

function Point(x, y) {
    const self = this;
    self.x = x;
    self.y = y;
    self.neg = function () {
        return new Point(-self.x, -self.y);
    }
}

function Kernel(...elements) {
    const self = this;
    let rootOfLen = Math.sqrt(elements.length);
    self.width = rootOfLen;
    self.height = rootOfLen;
    self.elements = elements;
    self.at = function (index) {
        if (index instanceof Point) {
            return self.elements[index.y * self.width + self.x];
        }
        return self.elements[index];
    };
    self.iter = function () {
        return self.elements;
    };
    self.contains = function (item) {
        return self.elements.includes(item);
    };
}

function List(...elements) {
    const self = this;
    self.elements = elements;
    self.at = function (index) {
        return self.elements[index];
    };
    self.concat = function (list) {
        return new List(self.elements.concat(list));
    };
    self.slice = function (lower, upper) {
        return new List(self.elements.slice(lower, upper));
    };
    self.iter = function () {
        return self.elements;
    };
    self.contains = function (item) {
        return self.elements.includes(item);
    };
}

function HashMap(obj) {
    const self = this;
    self.obj = obj
    self.at = function (key) {
        return self.obj[key];
    };
    self.iter = function () {
        return self.obj;
    };
}

function Rect(x, y, width, height) {
    const self = this;
    self.x = x;
    self.y = y;
    self.width = width;
    self.height = height;
    self.iter = function () {
        let points = [];
        let lastY = y + self.height;
        let lastX = x + self.width;
        for (y = self.y; y < lastY; y++) {
            for (x = self.x; x < lastX; x++) {
                points.push(new Point(x, y));
            }
        }
        return points;
    };
    self.contains = function (pt) {
        return pt.x >= self.x && pt.x < self.x + self.width
            && pt.y >= self.y && pt.y < self.y + self.height;
    }
}

function Circle(center, radius) {
    const self = this;
    self.center = center;
    self.radius = radius;
    self.contains = function (pt) {
        let dx = pt.x - self.center.x;
        let dy = pt.y - self.center.y;
        let distance = Math.sqrt(dx * dx + dy * dy);
        return distance <= self.radius;
    };
    self.iter = function () {
        throw "NOT IMPLEMENTED";
    };
}

function Line(point1, point2) {
    const self = this;
    self.point1 = point1;
    self.point2 = point2;
    self.iter = function () {
        throw "NOT IMPLEMENTED";
    };
}

function Polygon(points) {
    const self = this;
    self.points = points;
    self.iter = function () {
        throw "NOT IMPLEMENTED";
    };
}

const Op = {
    add: function (left, right) {
        if (left instanceof Color) {
            if (right instanceof Color) {
                return new Color(left.r + right.r, left.g + right.g, left.b + right.b, left.a);
            }
            return new Color(left.r + right, left.g + right, left.b + right, left.a);
        }
        if (left instanceof Point) {
            if (right instanceof Point) {
                return new Point(left.x + right.x, left.y + right.y);
            }
            return new Point(left.x + right, left.y + right);
        }
        if (right instanceof Color) {
            return new Color(right.r + left, right.g + left, right.b + left, right.a);
        }
        if (right instanceof Point) {
            return new Point(right.x + left, right.y + left);
        }
        return left + right;
    },
    sub: function (left, right) {
        if (left instanceof Color) {
            if (right instanceof Color) {
                return new Color(left.r - right.r, left.g - right.g, left.b - right.b, left.a);
            }
            return new Color(left.r - right, left.g - right, left.b - right, left.a);
        }
        if (left instanceof Point) {
            if (right instanceof Point) {
                return new Point(left.x - right.x, left.y - right.y);
            }
            return new Point(left.x - right, left.y - right);
        }
        if (right instanceof Color) {
            return new Color(left - right.r, left - right.g, left - right.b, right.a);
        }
        if (right instanceof Point) {
            return new Point(left - right.x, left - right.y);
        }
        return left - right;
    },
    mul: function (left, right) {
        if (left instanceof Color) {
            if (right instanceof Color) {
                return color01(left.r01 * right.r01, left.g01 * right.g01, left.b01 * right.b01, left.a01);
            }
            return new Color(left.r * right, left.g * right, left.b * right, left.a);
        }
        if (left instanceof Point) {
            if (right instanceof Point) {
                return new Point(left.x * right.x, left.y * right.y);
            }
            return new Point(left.x * right, left.y * right);
        }
        if (right instanceof Color) {
            return new Color(right.r * left, right.g * left, right.b * left, right.a);
        }
        if (right instanceof Point) {
            return new Point(right.x * left, right.y * left);
        }
        return left + right;
    },
    div: function (left, right) {
        if (left instanceof Color) {
            if (right instanceof Color) {
                return color01(left.r01 / right.r01, left.g01 / right.g01, left.b01 / right.b01, left.a01);
            }
            return new Color(left.r / right, left.g / right, left.b / right, left.a);
        }
        if (left instanceof Point) {
            if (right instanceof Point) {
                return new Point(left.x / right.x, left.y / right.y);
            }
            return new Point(left.x / right, left.y / right);
        }
        if (right instanceof Color) {
            return new Color(left / right.r, left / right.g, left / right.b, right.a);
        }
        if (right instanceof Point) {
            return new Point(left / right.x, left / right.y);
        }
        return left + right;
    },
    neg: function (val) {
        if (val.neg) {
            return val.neg();
        }
        return -val;
    },
    not: function (val) {
        return !val;
    },
};

function Surface(context, width, height) {
    const self = this;
    let imgdata = context.getImageData(0, 0, width, height);
    let pixels = imgdata.data;

    self.width = function () {
        return width;
    };
    self.height = function () {
        return height;
    };

    self.getPixel = function (pt) {
        let i = 4 * (pt.y * width + pt.x);
        return new Color(pixels[i], pixels[i + 1], pixels[i + 2], pixels[i + 3]);
    };
    self.setPixel = function (pt, color) {
        let i = 4 * (pt.y * width + pt.x);
        pixels[i + 0] = color.r;
        pixels[i + 1] = color.g;
        pixels[i + 2] = color.b;
        pixels[i + 3] = color.a;
    };

    self.getImageData = function () {
        return imgdata;
    };
}

function rgba(r, g, b, a) {
    return new Color(r, g, b, a);
}
function rgb(r, g, b) {
    return new Color(r, g, b, 255);
}
function rgba01(r01, g01, b01, a01) {
    return color01(r01, g01, b01, a01);
}
function rgb01(r01, g01, b01, a01) {
    return color01(r01, g01, b01, 1.0);
}
