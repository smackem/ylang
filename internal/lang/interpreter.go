package lang

import (
	"fmt"
	"image"
	"math"
	"reflect"
	"strings"
)

func interpret(program Program, bitmap BitmapContext) error {
	ir := newInterpreter(bitmap)
	if err := ir.visitStmtList(program.stmts); err != nil {
		if _, ok := err.(returnSignal); !ok { // return statement encountered
			return err
		}
	}
	return nil
}

type scope map[string]value
type functionScope struct {
	retval value
}

type interpreter struct {
	idents         []scope
	bitmap         BitmapContext
	functionScopes []functionScope
}

type returnSignal string

func (rs returnSignal) Error() string {
	return string(rs)
}

const returnSig returnSignal = returnSignal("RET")
const initialScopeCount int = 2

func newInterpreter(bitmap BitmapContext) *interpreter {
	ir := &interpreter{
		idents: []scope{make(scope)},
		bitmap: bitmap,
	}
	ir.newIdent(lastRectIdent, rect{})
	ir.newIdent("Black", NewRgba(0, 0, 0, 255))
	ir.newIdent("White", NewRgba(255, 255, 255, 255))
	ir.newIdent("Transparent", NewRgba(255, 255, 255, 0))
	ir.newIdent("Pi", Number(math.Pi))
	ir.newIdent("Rad2Deg", Number(180/math.Pi))
	ir.newIdent("Deg2Rad", Number(math.Pi/180))
	if bitmap != nil {
		ir.assignBounds(true)
		ir.newIdent("W", Number(bitmap.SourceWidth()))
		ir.newIdent("H", Number(bitmap.SourceHeight()))
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

func (ir *interpreter) pushFunctionScope(closure []scope) int {
	ir.functionScopes = append(ir.functionScopes, functionScope{})
	prevScopeCount := len(ir.idents)
	ir.idents = append(ir.idents, closure...)
	ir.pushScope()
	return prevScopeCount
}

func (ir *interpreter) popFunctionScope(scopeCount int) {
	ir.functionScopes = ir.functionScopes[:len(ir.functionScopes)-1]
	ir.idents = ir.idents[:scopeCount]
}

func (ir *interpreter) getReturnValue() value {
	if len(ir.functionScopes) <= 0 {
		return nil
	}
	return ir.functionScopes[len(ir.functionScopes)-1].retval
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

func (ir *interpreter) removeIdent(ident string) {
	delete(ir.idents[len(ir.idents)-1], ident)
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

func (ir *interpreter) assignBounds(new boolean) error {
	bounds := rect{image.Point{0, 0}, image.Point{ir.bitmap.SourceWidth(), ir.bitmap.SourceHeight()}}
	if new {
		return ir.newIdent("Bounds", bounds)
	}
	return ir.assignIdent("Bounds", bounds)
}

func (ir *interpreter) visitStmtList(stmts []statement) error {
	for _, s := range stmts {
		if err := ir.visitStmt(s); err != nil {
			if _, ok := err.(returnSignal); ok { // return statement encountered
				return err
			}

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

	case indexedAssignStmt:
		lval, ok := ir.findIdent(s.ident)
		if !ok {
			return fmt.Errorf("unkown identifier '%s'", s.ident)
		}
		ival, err := ir.visitExpr(s.index)
		if err != nil {
			return err
		}
		rval, err := ir.visitExpr(s.rhs)
		if err != nil {
			return err
		}
		if err := lval.indexAssign(ival, rval); err != nil {
			return err
		}

	case pixelAssignStmt:
		left, err := ir.visitExpr(s.lhs)
		if err != nil {
			return err
		}
		pos, ok := left.(point)
		if !ok {
			return fmt.Errorf("type mismatch: expected @point = color")
		}
		right, err := ir.visitExpr(s.rhs)
		if err != nil {
			return err
		}
		color, ok := right.(Color)
		if !ok {
			return fmt.Errorf("type mismatch: expected @point = color")
		}
		ir.bitmap.SetPixel(pos.X, pos.Y, color)

	case invocationStmt:
		if _, err := ir.visitExpr(s.invocation); err != nil {
			return err
		}

	case ifStmt:
		condVal, err := ir.visitExpr(s.cond)
		if err != nil {
			return err
		}
		b, ok := condVal.(boolean)
		if !ok {
			return fmt.Errorf("type mismatch: expected if(boolean)")
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
		rect, ok := collVal.(rect)
		if ok {
			ir.assignIdent(lastRectIdent, rect)
		}
		ir.pushScope()
		defer ir.popScope()
		ir.newIdent(s.ident, nil)
		return collVal.iterate(func(val value) error {
			ir.assignIdent(s.ident, val)
			if err := ir.visitStmtList(s.stmts); err != nil {
				return err
			}
			return nil
		})

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
		stepVal, err := ir.visitExpr(s.step)
		if err != nil {
			return err
		}
		stepN, ok := stepVal.(Number)
		if !ok {
			return fmt.Errorf("type mismatch: expected upper number")
		}
		ir.pushScope()
		defer ir.popScope()
		ir.newIdent(s.ident, nil)
		for n := lowerN; n < upperN; n += stepN {
			ir.assignIdent(s.ident, n)
			if err := ir.visitStmtList(s.stmts); err != nil {
				return err
			}
		}

	case whileStmt:
		for {
			condVal, err := ir.visitExpr(s.cond)
			if err != nil {
				return err
			}
			b, ok := condVal.(boolean)
			if !ok {
				return fmt.Errorf("type mismatch: expected while(boolean)")
			}
			if !b {
				break
			}
			err = ir.visitStmtList(s.stmts)
			if err != nil {
				return err
			}
		}

	case yieldStmt:
		return fmt.Errorf("yield not yet implemented")

	case logStmt:
		buf := strings.Builder{}
		for _, expr := range s.args {
			v, err := ir.visitExpr(expr)
			if err != nil {
				return err
			}
			buf.WriteString(formatValue(v, "", false))
		}
		fmt.Println(buf.String())

	case returnStmt:
		result, err := ir.visitExpr(s.result)
		if err != nil {
			return err
		}
		if len(ir.functionScopes) > 0 {
			ir.functionScopes[len(ir.functionScopes)-1].retval = result
			return returnSig
		}
		if _, isNil := result.(nilval); isNil {
			return returnSig
		}
		return fmt.Errorf("A script can only return 'nil' from root level")
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
	v, err := ir.visitExprInner(expr)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nilval{}, nil
	}
	return v, nil
}

func (ir *interpreter) visitExprInner(expr expression) (value, error) {
	switch e := expr.(type) {
	case ternaryExpr:
		condVal, err := ir.visitExpr(e.cond)
		if err != nil {
			return nil, err
		}
		b, ok := condVal.(boolean)
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
		b, ok := leftVal.(boolean)
		if !ok {
			return nil, fmt.Errorf("type mismatch: expected bool")
		}
		if b {
			return boolean(true), nil
		}
		rightVal, err := ir.visitExpr(e.right)
		if err != nil {
			return nil, err
		}
		b, ok = rightVal.(boolean)
		if !ok {
			return nil, fmt.Errorf("type mismatch: expected bool")
		}
		return boolean(b), nil

	case andExpr:
		leftVal, err := ir.visitExpr(e.left)
		if err != nil {
			return nil, err
		}
		b, ok := leftVal.(boolean)
		if !ok {
			return nil, fmt.Errorf("type mismatch: expected bool")
		}
		if !b {
			return boolean(false), nil
		}
		rightVal, err := ir.visitExpr(e.right)
		if err != nil {
			return nil, err
		}
		b, ok = rightVal.(boolean)
		if !ok {
			return nil, fmt.Errorf("type mismatch: expected bool")
		}
		return boolean(b), nil

	case eqExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			cmp, _ := left.compare(right)
			if cmp == nil {
				return falseVal, nil
			}
			n, ok := cmp.(Number)
			return boolean(ok && n == 0), nil
		})

	case neqExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			cmp, _ := left.compare(right)
			if cmp == nil {
				return boolean(true), nil
			}
			n, ok := cmp.(Number)
			return boolean(!ok || n != 0), nil
		})

	case gtExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			cmp, _ := left.compare(right)
			if cmp == nil {
				return falseVal, nil
			}
			n, ok := cmp.(Number)
			return boolean(ok && n > 0), nil
		})

	case geExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			cmp, _ := left.compare(right)
			if cmp == nil {
				return falseVal, nil
			}
			n, ok := cmp.(Number)
			return boolean(ok && n >= 0), nil
		})

	case ltExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			cmp, _ := left.compare(right)
			if cmp == nil {
				return falseVal, nil
			}
			n, ok := cmp.(Number)
			return boolean(ok && n < 0), nil
		})

	case leExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			cmp, _ := left.compare(right)
			if cmp == nil {
				return falseVal, nil
			}
			n, ok := cmp.(Number)
			return boolean(ok && n <= 0), nil
		})

	case concatExpr:
		return ir.visitBinaryExpr(e.left, e.right, func(left value, right value) (value, error) {
			return left.concat(right)
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
			return point{int(x + 0.5), int(y + 0.5)}, nil
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
		index, err := ir.visitExpr(e.index)
		if err != nil {
			return nil, err
		}
		return recvr.index(index)

	case indexRangeExpr:
		recvr, err := ir.visitExpr(e.recvr)
		if err != nil {
			return nil, err
		}
		lower, err := ir.visitExpr(e.lower)
		if err != nil {
			return nil, err
		}
		upper, err := ir.visitExpr(e.upper)
		if err != nil {
			return nil, err
		}
		return recvr.indexRange(lower, upper)

	case str:
		return e, nil

	case boolean:
		return e, nil

	case Number:
		return e, nil

	case Color:
		return e, nil

	case nilval:
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
		pos, ok := val.(point)
		if !ok {
			return nil, fmt.Errorf("")
		}
		return ir.bitmap.GetPixel(pos.X, pos.Y), nil

	case invokeExpr:
		args := []value{}
		for _, arg := range e.args {
			arg, err := ir.visitExpr(arg)
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		}
		return ir.invokeFunc(e.funcName, args)

	case kernelExpr:
		elementNumbers := make([]Number, len(e.elements))
		for i, element := range e.elements {
			elementVal, err := ir.visitExpr(element)
			if err != nil {
				return nil, err
			}
			n, ok := elementVal.(Number)
			if !ok {
				return nil, fmt.Errorf("type mismatch: kernel expr expects number elements")
			}
			elementNumbers[i] = n
		}
		rootOfLen := int(math.Sqrt(float64(len(elementNumbers))))
		return kernel{
			values: elementNumbers,
			width:  rootOfLen,
			height: rootOfLen,
		}, nil

	case functionExpr:
		return function{
			parameterNames: e.parameterNames,
			body:           e.body,
			closure:        ir.idents[2:], // omit constants and top scope - they are visible in any context
		}, nil

	case hashMapExpr:
		h := make(hashMap)
		for _, entry := range e.entries {
			key, err := ir.visitExpr(entry.key)
			if err != nil {
				return nil, err
			}
			val, err := ir.visitExpr(entry.value)
			if err != nil {
				return nil, err
			}
			h[key] = val
		}
		return h, nil

	case listExpr:
		l := list{
			elements: make([]value, len(e.elements)),
		}
		for i, elem := range e.elements {
			val, err := ir.visitExpr(elem)
			if err != nil {
				return nil, err
			}
			l.elements[i] = val
		}
		return l, nil

	case pipelineExpr:
		left, err := ir.visitExpr(e.left)
		if err != nil {
			return nil, err
		}
		pipelineValueIdent := tokenTypeNames[ttDollar]
		ir.newIdent(pipelineValueIdent, left)
		right, err := ir.visitExpr(e.right)
		if err != nil {
			return nil, err
		}
		ir.removeIdent(pipelineValueIdent)
		return right, nil
	}

	return nil, fmt.Errorf("unknown expression type %s", reflect.TypeOf(expr))
}

func (ir *interpreter) invokeFunc(name string, arguments []value) (value, error) {
	val, ok := ir.findIdent(name)
	if ok {
		return ir.invokeFunctionExpr(name, val, arguments)
	}
	return ir.invokeBuiltinFunction(name, arguments)
}

func (ir *interpreter) invokeFunctionExpr(name string, val value, arguments []value) (value, error) {
	fn, ok := val.(function)
	if !ok {
		return nil, fmt.Errorf("%s is invoked like a function, but refers to a %s", name, reflect.TypeOf(val))
	}
	if len(arguments) != len(fn.parameterNames) {
		return nil, fmt.Errorf("%s is invoked with %d arguments, but is declared with %d parameters", name, len(arguments), len(fn.parameterNames))
	}

	prevScopeCount := ir.pushFunctionScope(fn.closure)
	defer ir.popFunctionScope(prevScopeCount)
	for i, argument := range arguments {
		ir.newIdent(fn.parameterNames[i], argument)
	}

	if err := ir.visitStmtList(fn.body); err != nil {
		if _, ok := err.(returnSignal); !ok { // return statement encountered
			return nil, err
		}
	}
	return ir.getReturnValue(), nil
}

func (ir *interpreter) invokeBuiltinFunction(name string, arguments []value) (value, error) {
	fs, ok := functions[name]
	if !ok {
		return nil, fmt.Errorf("unkown function '%s'", name)
	}
	var err error
	for _, f := range fs {
		if err = validateArguments(arguments, f.params); err == nil {
			return f.body(ir, arguments)
		}
	}

	buffer := strings.Builder{}
	for i, f := range fs {
		if i > 0 {
			buffer.WriteString("\n")
		}
		buffer.WriteString(signature(name, f))
	}
	return nil, fmt.Errorf("no fitting overload for function '%s': %s. possible overloads:\n%s", name, err, buffer.String())
}

func validateArguments(arguments []value, params []reflect.Type) error {
	argsCount := len(arguments)
	paramsCount := len(params)
	discreteArgsCount := argsCount

	if argsCount != paramsCount {
		lastParam := params[paramsCount-1]
		if lastParam.Kind() != reflect.Slice {
			return fmt.Errorf("wrong number of arguments: expected %d, got %d", paramsCount, argsCount)
		}
		if argsCount < paramsCount-1 {
			return fmt.Errorf("wrong number of arguments: expected %d, got %d", paramsCount, argsCount)
		}
		for i := paramsCount - 1; i < argsCount; i++ {
			if !hasMatchingType(arguments[i], lastParam.Elem()) {
				return fmt.Errorf("argument type mismatch at argument %d: expected %s, got %s", i, lastParam.Elem().Name(), reflect.TypeOf(arguments[i]).Name())
			}
		}
		discreteArgsCount = paramsCount - 1
	}

	for i := 0; i < discreteArgsCount; i++ {
		if hasMatchingType(arguments[i], params[i]) {
			// direct match
			continue
		}
		if i == discreteArgsCount-1 && params[i].Kind() == reflect.Slice && hasMatchingType(arguments[i], params[i].Elem()) {
			// trailing slice with one arg
			continue
		}
		return fmt.Errorf("argument type mismatch at argument %d: expected %s, got %s", i, params[i].Name(), reflect.TypeOf(arguments[i]).Name())
	}

	return nil
}

func hasMatchingType(v value, typ reflect.Type) bool {
	return typ == valueType || reflect.TypeOf(v) == typ
}
