package lang

import (
	"bytes"
	"fmt"
	"image"
	"math"
)

type interpreter struct {
	idents map[string]value
	bitmap Bitmap
}

func (ir interpreter) visitProgram(stmts []statement) error {
	return ir.visitStmtList(stmts)
}

func (ir interpreter) visitStmtList(stmts []statement) error {
	for _, s := range stmts {
		if err := ir.visitStmt(s); err != nil {
			return err
		}
	}
	return nil
}

const lastRectIdent string = "@<lastRect>" // this is safe because its not a valid ident

func (ir interpreter) visitStmt(stmt statement) error {
	switch s := stmt.(type) {
	case declStmt:
		v, err := ir.visitExpr(s.rhs)
		if err != nil {
			return err
		}
		ir.idents[s.ident] = v

	case assignStmt:
		v, err := ir.visitExpr(s.rhs)
		if err != nil {
			return err
		}
		ir.idents[s.ident] = v

	case pixelAssignStmt:
		left, err := ir.visitExpr(s.lhs)
		if err != nil {
			return err
		}
		pos, ok := left.(Position)
		if !ok {
			return fmt.Errorf("type mismatch: expected @position = color")
		}
		right, err := ir.visitExpr(s.rhs)
		if err != nil {
			return err
		}
		color, ok := right.(Color)
		if !ok {
			return fmt.Errorf("type mismatch: expected @position = color")
		}
		ir.bitmap.SetPixel(pos.X, pos.Y, color)

	case ifStmt:
		condVal, err := ir.visitExpr(s.cond)
		if err != nil {
			return err
		}
		b, ok := condVal.(Bool)
		if !ok {
			return fmt.Errorf("type mismatch: expected if(bool)")
		}
		if b {
			return ir.visitStmtList(s.trueStmts)
		}
		if s.falseStmts != nil {
			return ir.visitStmtList(s.falseStmts)
		}

	case forStmt:
		collVal, err := ir.visitExpr(s.collection)
		if err != nil {
			return err
		}
		rect, ok := collVal.(Rect)
		if !ok {
			return fmt.Errorf("type mismatch: expected for ident in rect")
		}
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				ir.idents[s.ident] = Position{x, y}
				if err := ir.visitStmtList(s.stmts); err != nil {
					return err
				}
			}
		}

	case forRangeStmt:
		lowerVal, err := ir.visitExpr(s.lower)
		if err != nil {
			return err
		}
		lowerN, ok := lowerVal.(Number)
		if !ok {
			return fmt.Errorf("type mismatch: expected lower number")
		}
		upperVal, err := ir.visitExpr(s.upper)
		if err != nil {
			return err
		}
		upperN, ok := upperVal.(Number)
		if !ok {
			return fmt.Errorf("type mismatch: expected upper number")
		}
		for n := lowerN; n < upperN; n++ {
			ir.idents[s.ident] = n
			if err := ir.visitStmtList(s.stmts); err != nil {
				return err
			}
		}

	case yieldStmt:
		return fmt.Errorf("yield not yet implemented")

	case logStmt:
		buf := bytes.NewBuffer(nil)
		for _, expr := range s.parameters {
			v, err := ir.visitExpr(expr)
			if err != nil {
				return err
			}
			buf.WriteString(v.printStr())
		}
		fmt.Printf("%s", buf.String())

	case bltStmt:
		var expr expression
		if s.rect != nil {
			var err error
			expr, err = ir.visitExpr(s.rect)
			if err != nil {
				return err
			}
		} else {
			expr = ir.idents[lastRectIdent]
		}
		rect, ok := expr.(Rect)
		if !ok {
			return fmt.Errorf("type mismatch: blt expects rect")
		}
		ir.bitmap.Blt(rect)
	}

	return nil
}

type binaryExprVisitor func(left value, right value) (value, error)

func (ir interpreter) visitBinaryExpr(left expression, right expression, visitor binaryExprVisitor) (value, error) {
	leftVal, err := ir.visitExpr(left)
	if err != nil {
		return nil, err
	}
	rightVal, err := ir.visitExpr(right)
	if err != nil {
		return nil, err
	}
	return visitor(leftVal, rightVal)
}

func (ir interpreter) visitExpr(expr expression) (value, error) {
	switch e := expr.(type) {
	case ternaryExpr:
		condVal, err := ir.visitExpr(e.cond)
		if err != nil {
			return nil, err
		}
		b, ok := condVal.(Bool)
		if !ok {
			return nil, fmt.Errorf("type mismatch: expected bool?x:x")
		}
		if b {
			return ir.visitExpr(e.trueResult)
		}
		return ir.visitExpr(e.falseResult)

	case orExpr:
		leftVal, err := ir.visitExpr(e.left)
		if err != nil {
			return nil, err
		}
		b, ok := leftVal.(Bool)
		if !ok {
			return nil, fmt.Errorf("type mismatch: expected bool")
		}
		if b {
			return Bool(true), nil
		}
		rightVal, err := ir.visitExpr(e.right)
		if err != nil {
			return nil, err
		}
		b, ok = rightVal.(Bool)
		if !ok {
			return nil, fmt.Errorf("type mismatch: expected bool")
		}
		return Bool(b), nil

	case andExpr:
		leftVal, err := ir.visitExpr(e.left)
		if err != nil {
			return nil, err
		}
		b, ok := leftVal.(Bool)
		if !ok {
			return nil, fmt.Errorf("type mismatch: expected bool")
		}
		if !b {
			return Bool(false), nil
		}
		rightVal, err := ir.visitExpr(e.right)
		if err != nil {
			return nil, err
		}
		b, ok = rightVal.(Bool)
		if !ok {
			return nil, fmt.Errorf("type mismatch: expected bool")
		}
		return Bool(b), nil

	case eqExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			return left.equals(right)
		})

	case neqExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			result, err := left.equals(right)
			if err != nil {
				return nil, err
			}
			return Bool(!result.(Bool)), nil
		})

	case gtExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			return left.greaterThan(right)
		})

	case geExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			return left.greaterThanOrEqual(right)
		})

	case ltExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			return left.lessThan(right)
		})

	case leExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			return left.lessThanOrEqual(right)
		})

	case addExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			return left.add(right)
		})

	case subExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			return left.sub(right)
		})

	case mulExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			return left.mul(right)
		})

	case divExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			return left.div(right)
		})

	case modExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			return left.mod(right)
		})

	case inExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			return left.in(right)
		})

	case negExpr:
		leftVal, err := ir.visitExpr(e.inner)
		if err != nil {
			return nil, err
		}
		return leftVal.neg()

	case notExpr:
		leftVal, err := ir.visitExpr(e.inner)
		if err != nil {
			return nil, err
		}
		return leftVal.not()

	case posExpr:
		return ir.visitBinaryExpr(e.x, e.y, func(xVal value, yVal value) (value, error) {
			x, ok := xVal.(Number)
			if !ok {
				return nil, fmt.Errorf("type mismatch: expected pos(Number, Number)")
			}
			y, ok := yVal.(Number)
			if !ok {
				return nil, fmt.Errorf("type mismatch: expected pos(Number, Number)")
			}
			return Position{int(x + 0.5), int(y + 0.5)}, nil
		})

	case memberExpr:
		recvrVal, err := ir.visitExpr(e.recvr)
		if err != nil {
			return nil, err
		}
		return recvrVal.property(e.member)

	case stringExpr:
		return String(e), nil

	case boolExpr:
		return Bool(e), nil

	case identExpr:
		return ir.idents[string(e)], nil

	case atExpr:
		val, err := ir.visitExpr(e.inner)
		if err != nil {
			return nil, err
		}
		pos, ok := val.(Position)
		if !ok {
			return nil, fmt.Errorf("")
		}
		return ir.bitmap.GetPixel(pos.X, pos.Y), nil

	case invokeExpr:
		parameterVals := []value{}
		for _, parameter := range e.parameters {
			parameterVal, err := ir.visitExpr(parameter)
			if err != nil {
				return nil, err
			}
			parameterVals = append(parameterVals, parameterVal)
		}
		return ir.invokeFunc(e.funcName, parameterVals)

	case kernelExpr:
		elementNumbers := []Number{}
		for _, element := range e.elements {
			elementVal, err := ir.visitExpr(element)
			if err != nil {
				return nil, err
			}
			n, ok := elementVal.(Number)
			if !ok {
				return nil, fmt.Errorf("type mismatch: kernel expr expects number elements")
			}
			elementNumbers = append(elementNumbers, n)
		}
		rootOfLen := int(math.Sqrt(float64(len(elementNumbers))))
		return kernel{
			values: elementNumbers,
			length: rootOfLen,
			radius: rootOfLen / 2,
		}, nil
	}

	return nil, nil
}

func (ir interpreter) invokeFunc(name string, values []value) (value, error) {
	switch name {
	case "rgb":
		return ir.invokeRgb(values)
	case "srgb":
		return ir.invokeSrgb(values)
	case "rgba":
		return ir.invokeRgba(values)
	case "srgba":
		return ir.invokeSrgba(values)
	case "rect":
		return ir.invokeRect(values)
	case "convolute":
		return ir.invokeConvolute(values)
	}
	return nil, nil
}

func getNumbers(values []value, errStr string) ([]Number, error) {
	numbers := make([]Number, len(values))
	for i, v := range values {
		n, ok := v.(Number)
		if !ok {
			return nil, fmt.Errorf(errStr)
		}
		numbers[i] = n
	}
	return numbers, nil
}

func (ir interpreter) invokeRgb(values []value) (value, error) {
	const argMismatch string = "argument mismatch: rgb(number, number, number)"

	if len(values) != 3 {
		return nil, fmt.Errorf(argMismatch)
	}
	numbers, err := getNumbers(values, argMismatch)
	if err != nil {
		return nil, err
	}

	return NewRgba(numbers[0], numbers[1], numbers[2], 255), nil
}

func (ir interpreter) invokeSrgb(values []value) (value, error) {
	const argMismatch string = "argument mismatch: srgb(number, number, number)"

	if len(values) != 3 {
		return nil, fmt.Errorf(argMismatch)
	}
	numbers, err := getNumbers(values, argMismatch)
	if err != nil {
		return nil, err
	}

	return NewSrgba(numbers[0], numbers[1], numbers[2], 1.0), nil
}

func (ir interpreter) invokeRgba(values []value) (value, error) {
	const argMismatch string = "argument mismatch: rgba(number, number, number, number)"

	if len(values) != 4 {
		return nil, fmt.Errorf(argMismatch)
	}
	numbers, err := getNumbers(values, argMismatch)
	if err != nil {
		return nil, err
	}

	return NewRgba(numbers[0], numbers[1], numbers[2], numbers[3]), nil
}

func (ir interpreter) invokeSrgba(values []value) (value, error) {
	const argMismatch string = "argument mismatch: srgba(number, number, number, number)"

	if len(values) != 4 {
		return nil, fmt.Errorf(argMismatch)
	}
	numbers, err := getNumbers(values, argMismatch)
	if err != nil {
		return nil, err
	}

	return NewSrgba(numbers[0], numbers[1], numbers[2], numbers[3]), nil
}

func (ir interpreter) invokeRect(values []value) (value, error) {
	const argMismatch string = "argument mismatch: rect(x:number, y:number, w:number, h:number)"

	if len(values) != 4 {
		return nil, fmt.Errorf(argMismatch)
	}
	numbers, err := getNumbers(values, argMismatch)
	if err != nil {
		return nil, err
	}

	return Rect{
		Min: image.Point{int(numbers[0] + 0.5), int(numbers[1] + 0.5)},
		Max: image.Point{int(numbers[2] + 0.5), int(numbers[3] + 0.5)},
	}, nil
}

func (ir interpreter) invokeConvolute(values []value) (value, error) {
	const argMismatch string = "argument mismatch: convolute(position, kernel)"

	if len(values) != 2 {
		return nil, fmt.Errorf(argMismatch)
	}
	posVal, ok := values[0].(Position)
	if !ok {
		return nil, fmt.Errorf(argMismatch)
	}
	kernelVal, ok := values[1].(kernel)
	if !ok {
		return nil, fmt.Errorf(argMismatch)
	}

	return ir.bitmap.Convolute(posVal.X, posVal.Y, kernelVal.radius, kernelVal.length, kernelVal.values), nil
}
