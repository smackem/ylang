package lang

import (
	"bytes"
	"fmt"
	"image"
	"math"
	"reflect"
)

func interpret(program Program, bitmap BitmapContext) error {
	return newInterpreter(bitmap).visitStmtList(program.stmts)
}

type scope map[string]value

type interpreter struct {
	idents []scope
	bitmap BitmapContext
}

func newInterpreter(bitmap BitmapContext) *interpreter {
	ir := &interpreter{
		idents: []scope{make(scope)},
		bitmap: bitmap,
	}
	ir.newIdent(lastRectIdent, Rect{})
	ir.newIdent("Black", NewRgba(0, 0, 0, 255))
	ir.newIdent("White", NewRgba(255, 255, 255, 255))
	ir.newIdent("Transparent", NewRgba(255, 255, 255, 0))
	ir.newIdent("Pi", Number(math.Pi))
	if bitmap != nil {
		ir.newIdent("Bounds", Rect{image.Point{0, 0}, image.Point{bitmap.Width(), bitmap.Height()}})
		ir.newIdent("W", Number(bitmap.Width()))
		ir.newIdent("H", Number(bitmap.Height()))
	}
	ir.pushScope()
	return ir
}

func (ir *interpreter) pushScope() {
	ir.idents = append(ir.idents, make(scope))
}

func (ir *interpreter) popScope() {
	ir.idents = ir.idents[:len(ir.idents)-1]
}

func (ir interpreter) findIdent(ident string) (value, bool) {
	last := len(ir.idents) - 1
	for i := range ir.idents {
		scope := ir.idents[last-i]
		val, ok := scope[ident]
		if ok {
			return val, true
		}
	}
	return nil, false
}

func (ir *interpreter) newIdent(ident string, val value) error {
	ir.idents[len(ir.idents)-1][ident] = val
	return nil
}

func (ir *interpreter) assignIdent(ident string, val value) error {
	last := len(ir.idents) - 1
	for i := range ir.idents {
		scope := ir.idents[last-i]
		_, ok := scope[ident]
		if ok {
			scope[ident] = val
			return nil
		}
	}
	return fmt.Errorf("identifier '%s' not found", ident)
}

func (ir *interpreter) visitStmtList(stmts []statement) error {
	for _, s := range stmts {
		if err := ir.visitStmt(s); err != nil {
			tok := s.getToken()
			return fmt.Errorf("line %d near '%s': %s", tok.LineNumber, tok.Lexeme, err)
		}
	}
	return nil
}

const lastRectIdent string = "@:R" // this is safe because its not a valid ident

func (ir *interpreter) visitStmt(stmt statement) error {
	switch s := stmt.(type) {
	case declStmt:
		v, err := ir.visitExpr(s.rhs)
		if err != nil {
			return err
		}
		err = ir.newIdent(s.ident, v)
		if err != nil {
			return err
		}

	case assignStmt:
		v, err := ir.visitExpr(s.rhs)
		if err != nil {
			return err
		}
		err = ir.assignIdent(s.ident, v)
		if err != nil {
			return err
		}

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
			ir.pushScope()
			defer ir.popScope()
			return ir.visitStmtList(s.trueStmts)
		}
		if s.falseStmts != nil {
			ir.pushScope()
			defer ir.popScope()
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
		ir.pushScope()
		defer ir.popScope()
		ir.newIdent(s.ident, nil)
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				ir.assignIdent(s.ident, Position{x, y})
				if err := ir.visitStmtList(s.stmts); err != nil {
					return err
				}
			}
		}
		ir.assignIdent(lastRectIdent, rect)

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
		ir.pushScope()
		defer ir.popScope()
		ir.newIdent(s.ident, nil)
		for n := lowerN; n < upperN; n++ {
			ir.assignIdent(s.ident, n)
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
		fmt.Printf("%s\n", buf.String())

	case bltStmt:
		expr, err := ir.visitExpr(s.rect)
		if err != nil {
			return err
		}
		rect, ok := expr.(Rect)
		if !ok {
			return fmt.Errorf("type mismatch: blt expects rect")
		}
		ir.bitmap.BltToTarget(rect.Min.X, rect.Min.Y, rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y)

	case commitStmt:
		var expr expression
		if s.rect != nil {
			var err error
			expr, err = ir.visitExpr(s.rect)
			if err != nil {
				return err
			}
		} else {
			expr, _ = ir.findIdent(lastRectIdent)
		}
		rect, ok := expr.(Rect)
		if !ok {
			return fmt.Errorf("type mismatch: commit expects rect")
		}
		ir.bitmap.BltToSource(rect.Min.X, rect.Min.Y, rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y)
	}

	return nil
}

type binaryExprVisitor func(left value, right value) (value, error)

func (ir *interpreter) visitBinaryExpr(left expression, right expression, visitor binaryExprVisitor) (value, error) {
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

func (ir *interpreter) visitExpr(expr expression) (value, error) {
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

	case indexExpr:
		recvr, err := ir.visitExpr(e.recvr)
		if err != nil {
			return nil, err
		}
		k, ok := recvr.(kernel)
		if !ok {
			return nil, fmt.Errorf("type mismatch: expected kernel@index but found %s@index", recvr)
		}
		index, err := ir.visitExpr(e.index)
		if err != nil {
			return nil, err
		}
		switch i := index.(type) {
		case Number:
			return k.values[int(i)], nil
		case Position:
			return k.values[i.Y*k.width+i.X], nil
		}
		return nil, fmt.Errorf("type mismatch: expected kernel@number or kernel@position but found kernel@%s", index)

	case String:
		return e, nil

	case Bool:
		return e, nil

	case Number:
		return e, nil

	case Color:
		return e, nil

	case identExpr:
		val, ok := ir.findIdent(string(e))
		if !ok {
			return nil, fmt.Errorf("identifier '%s' not found", e)
		}
		return val, nil

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
			width:  rootOfLen,
			height: rootOfLen,
		}, nil
	}

	return nil, fmt.Errorf("unknown expression type %s", reflect.TypeOf(expr))
}

func (ir *interpreter) invokeFunc(name string, values []value) (value, error) {
	f, ok := functions[name]
	if !ok {
		return nil, fmt.Errorf("unkown function '%s'", name)
	}
	if len(values) != len(f.params) {
		return nil, fmt.Errorf("wrong number of arguments for function '%s': expected %d, got %d", name, len(f.params), len(values))
	}
	for i := 0; i < len(values); i++ {
		if reflect.TypeOf(values[i]) != f.params[i] {
			return nil, fmt.Errorf("argument type mismatch for function '%s' at argument %d: expected %s, got %s", name, i, f.params[i].Name(), reflect.TypeOf(values[i]).Name())
		}
	}
	return f.body(ir, values)
}
