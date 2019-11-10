package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"github.com/smackem/ylang/internal/lexer"
	"github.com/smackem/ylang/internal/parser"
	"image"
	"math"
	"reflect"
	"strings"
)

// package initialization
func init() {
	initFunctions()
}

func Interpret(program parser.Program, bitmap BitmapContext) error {
	ir := newInterpreter(bitmap)
	if err := ir.visitStmtList(program.Stmts); err != nil {
		if _, ok := err.(returnSignal); !ok { // return statement encountered
			return err
		}
	}
	return nil
}

type scope map[string]Value
type functionScope struct {
	retval Value
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

//noinspection ALL
func newInterpreter(bitmap BitmapContext) *interpreter {
	ir := &interpreter{
		idents: []scope{make(scope)},
		bitmap: bitmap,
	}
	ir.newIdent(lastRectIdent, Rect{})
	ir.newIdent("Black", Color(lang.NewRgba(0, 0, 0, 255)))
	ir.newIdent("White", Color(lang.NewRgba(255, 255, 255, 255)))
	ir.newIdent("Transparent", Color(lang.NewRgba(255, 255, 255, 0)))
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

func (ir *interpreter) getReturnValue() Value {
	if len(ir.functionScopes) <= 0 {
		return nil
	}
	return ir.functionScopes[len(ir.functionScopes)-1].retval
}

func (ir interpreter) findIdent(ident string) (Value, bool) {
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

func (ir *interpreter) newIdent(ident string, val Value) error {
	ir.idents[len(ir.idents)-1][ident] = val
	return nil
}

func (ir *interpreter) removeIdent(ident string) {
	delete(ir.idents[len(ir.idents)-1], ident)
}

func (ir *interpreter) assignIdent(ident string, val Value) error {
	last := len(ir.idents) - 1
	for i := range ir.idents {
		scope := ir.idents[last-i]
		_, ok := scope[ident]
		if ok {
			scope[ident] = val
			return nil
		}
	}
	return fmt.Errorf("identifier '%s' Not found", ident)
}

func (ir *interpreter) assignBounds(new Boolean) error {
	bounds := Rect{image.Point{0, 0}, image.Point{ir.bitmap.SourceWidth(), ir.bitmap.SourceHeight()}}
	if new {
		return ir.newIdent("Bounds", bounds)
	}
	return ir.assignIdent("Bounds", bounds)
}

func (ir *interpreter) visitStmtList(stmts []parser.Statement) error {
	for _, s := range stmts {
		if err := ir.visitStmt(s); err != nil {
			if _, ok := err.(returnSignal); ok { // return statement encountered
				return err
			}

			tok := s.Token()
			return fmt.Errorf("line %d near '%s': %s", tok.LineNumber, tok.Lexeme, err)
		}
	}
	return nil
}

const lastRectIdent string = "@:R" // this is safe because its Not a valid ident

func (ir *interpreter) visitStmt(stmt parser.Statement) error {
	switch s := stmt.(type) {
	case parser.DeclStmt:
		v, err := ir.visitExpr(s.Rhs)
		if err != nil {
			return err
		}
		err = ir.newIdent(s.Ident, v)
		if err != nil {
			return err
		}

	case parser.AssignStmt:
		v, err := ir.visitExpr(s.Rhs)
		if err != nil {
			return err
		}
		err = ir.assignIdent(s.Ident, v)
		if err != nil {
			return err
		}

	case parser.IndexedAssignStmt:
		lval, ok := ir.findIdent(s.Ident)
		if !ok {
			return fmt.Errorf("unkown identifier '%s'", s.Ident)
		}
		ival, err := ir.visitExpr(s.Index)
		if err != nil {
			return err
		}
		rval, err := ir.visitExpr(s.Rhs)
		if err != nil {
			return err
		}
		if err := lval.IndexAssign(ival, rval); err != nil {
			return err
		}

	case parser.PixelAssignStmt:
		left, err := ir.visitExpr(s.Lhs)
		if err != nil {
			return err
		}
		pos, ok := left.(Point)
		if !ok {
			return fmt.Errorf("type mismatch: expected @point = color")
		}
		right, err := ir.visitExpr(s.Rhs)
		if err != nil {
			return err
		}
		color, ok := right.(Color)
		if !ok {
			return fmt.Errorf("type mismatch: expected @point = color")
		}
		ir.bitmap.SetPixel(pos.X, pos.Y, lang.Color(color))

	case parser.InvocationStmt:
		if _, err := ir.visitExpr(s.Invocation); err != nil {
			return err
		}

	case parser.IfStmt:
		condVal, err := ir.visitExpr(s.Cond)
		if err != nil {
			return err
		}
		b, ok := condVal.(Boolean)
		if !ok {
			return fmt.Errorf("type mismatch: expected if(boolean)")
		}
		if b {
			ir.pushScope()
			defer ir.popScope()
			return ir.visitStmtList(s.TrueStmts)
		}
		if s.FalseStmts != nil {
			ir.pushScope()
			defer ir.popScope()
			return ir.visitStmtList(s.FalseStmts)
		}

	case parser.ForStmt:
		collVal, err := ir.visitExpr(s.Collection)
		if err != nil {
			return err
		}
		rect, ok := collVal.(Rect)
		if ok {
			ir.assignIdent(lastRectIdent, rect)
		}
		ir.pushScope()
		defer ir.popScope()
		ir.newIdent(s.Ident, nil)
		return collVal.Iterate(func(val Value) error {
			ir.assignIdent(s.Ident, val)
			if err := ir.visitStmtList(s.Stmts); err != nil {
				return err
			}
			return nil
		})

	case parser.ForRangeStmt:
		lowerVal, err := ir.visitExpr(s.Lower)
		if err != nil {
			return err
		}
		lowerN, ok := lowerVal.(Number)
		if !ok {
			return fmt.Errorf("type mismatch: expected lower number")
		}
		upperVal, err := ir.visitExpr(s.Upper)
		if err != nil {
			return err
		}
		upperN, ok := upperVal.(Number)
		if !ok {
			return fmt.Errorf("type mismatch: expected upper number")
		}
		stepVal, err := ir.visitExpr(s.Step)
		if err != nil {
			return err
		}
		stepN, ok := stepVal.(Number)
		if !ok {
			return fmt.Errorf("type mismatch: expected upper number")
		}
		ir.pushScope()
		defer ir.popScope()
		ir.newIdent(s.Ident, nil)
		for n := lowerN; n < upperN; n += stepN {
			ir.assignIdent(s.Ident, n)
			if err := ir.visitStmtList(s.Stmts); err != nil {
				return err
			}
		}

	case parser.WhileStmt:
		for {
			condVal, err := ir.visitExpr(s.Cond)
			if err != nil {
				return err
			}
			b, ok := condVal.(Boolean)
			if !ok {
				return fmt.Errorf("type mismatch: expected while(boolean)")
			}
			if !b {
				break
			}
			err = ir.visitStmtList(s.Stmts)
			if err != nil {
				return err
			}
		}

	case parser.YieldStmt:
		return fmt.Errorf("yield Not yet implemented")

	case parser.LogStmt:
		buf := strings.Builder{}
		for _, expr := range s.Args {
			v, err := ir.visitExpr(expr)
			if err != nil {
				return err
			}
			buf.WriteString(formatValue(v, "", false))
		}
		fmt.Println(buf.String())

	case parser.ReturnStmt:
		result, err := ir.visitExpr(s.Result)
		if err != nil {
			return err
		}
		if len(ir.functionScopes) > 0 {
			ir.functionScopes[len(ir.functionScopes)-1].retval = result
			return returnSig
		}
		if _, isNil := result.(Nilval); isNil {
			return returnSig
		}
		return fmt.Errorf("A script can only return 'nil' from root level")
	}

	return nil
}

type binaryExprVisitor func(left Value, right Value) (Value, error)

func (ir *interpreter) visitBinaryExpr(left parser.Expression, right parser.Expression, visitor binaryExprVisitor) (Value, error) {
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

func (ir *interpreter) visitExpr(expr parser.Expression) (Value, error) {
	v, err := ir.visitExprInner(expr)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return Nilval{}, nil
	}
	return v, nil
}

func (ir *interpreter) visitExprInner(expr parser.Expression) (Value, error) {
	switch e := expr.(type) {
	case parser.TernaryExpr:
		condVal, err := ir.visitExpr(e.Cond)
		if err != nil {
			return nil, err
		}
		b, ok := condVal.(Boolean)
		if !ok {
			return nil, fmt.Errorf("type mismatch: expected bool?x:x")
		}
		if b {
			return ir.visitExpr(e.TrueResult)
		}
		return ir.visitExpr(e.FalseResult)

	case parser.OrExpr:
		leftVal, err := ir.visitExpr(e.Left)
		if err != nil {
			return nil, err
		}
		b, ok := leftVal.(Boolean)
		if !ok {
			return nil, fmt.Errorf("type mismatch: expected bool")
		}
		if b {
			return Boolean(true), nil
		}
		rightVal, err := ir.visitExpr(e.Right)
		if err != nil {
			return nil, err
		}
		b, ok = rightVal.(Boolean)
		if !ok {
			return nil, fmt.Errorf("type mismatch: expected bool")
		}
		return Boolean(b), nil

	case parser.AndExpr:
		leftVal, err := ir.visitExpr(e.Left)
		if err != nil {
			return nil, err
		}
		b, ok := leftVal.(Boolean)
		if !ok {
			return nil, fmt.Errorf("type mismatch: expected bool")
		}
		if !b {
			return Boolean(false), nil
		}
		rightVal, err := ir.visitExpr(e.Right)
		if err != nil {
			return nil, err
		}
		b, ok = rightVal.(Boolean)
		if !ok {
			return nil, fmt.Errorf("type mismatch: expected bool")
		}
		return Boolean(b), nil

	case parser.EqExpr:
		return ir.visitBinaryExpr(e.Left, e.Right, func(left Value, right Value) (Value, error) {
			cmp, _ := left.Compare(right)
			if cmp == nil {
				return Boolean(lang.FalseVal), nil
			}
			n, ok := cmp.(Number)
			return Boolean(ok && n == 0), nil
		})

	case parser.NeqExpr:
		return ir.visitBinaryExpr(e.Left, e.Right, func(left Value, right Value) (Value, error) {
			cmp, _ := left.Compare(right)
			if cmp == nil {
				return Boolean(lang.TrueVal), nil
			}
			n, ok := cmp.(Number)
			return Boolean(!ok || n != 0), nil
		})

	case parser.GtExpr:
		return ir.visitBinaryExpr(e.Left, e.Right, func(left Value, right Value) (Value, error) {
			cmp, _ := left.Compare(right)
			if cmp == nil {
				return Boolean(lang.FalseVal), nil
			}
			n, ok := cmp.(Number)
			return Boolean(ok && n > 0), nil
		})

	case parser.GeExpr:
		return ir.visitBinaryExpr(e.Left, e.Right, func(left Value, right Value) (Value, error) {
			cmp, _ := left.Compare(right)
			if cmp == nil {
				return Boolean(lang.FalseVal), nil
			}
			n, ok := cmp.(Number)
			return Boolean(ok && n >= 0), nil
		})

	case parser.LtExpr:
		return ir.visitBinaryExpr(e.Left, e.Right, func(left Value, right Value) (Value, error) {
			cmp, _ := left.Compare(right)
			if cmp == nil {
				return Boolean(lang.FalseVal), nil
			}
			n, ok := cmp.(Number)
			return Boolean(ok && n < 0), nil
		})

	case parser.LeExpr:
		return ir.visitBinaryExpr(e.Left, e.Right, func(left Value, right Value) (Value, error) {
			cmp, _ := left.Compare(right)
			if cmp == nil {
				return Boolean(lang.FalseVal), nil
			}
			n, ok := cmp.(Number)
			return Boolean(ok && n <= 0), nil
		})

	case parser.ConcatExpr:
		return ir.visitBinaryExpr(e.Left, e.Right, func(left Value, right Value) (Value, error) {
			return left.Concat(right)
		})

	case parser.AddExpr:
		return ir.visitBinaryExpr(e.Left, e.Right, func(left Value, right Value) (Value, error) {
			return left.Add(right)
		})

	case parser.SubExpr:
		return ir.visitBinaryExpr(e.Left, e.Right, func(left Value, right Value) (Value, error) {
			return left.Sub(right)
		})

	case parser.MulExpr:
		return ir.visitBinaryExpr(e.Left, e.Right, func(left Value, right Value) (Value, error) {
			return left.Mul(right)
		})

	case parser.DivExpr:
		return ir.visitBinaryExpr(e.Left, e.Right, func(left Value, right Value) (Value, error) {
			return left.Div(right)
		})

	case parser.ModExpr:
		return ir.visitBinaryExpr(e.Left, e.Right, func(left Value, right Value) (Value, error) {
			return left.Mod(right)
		})

	case parser.InExpr:
		return ir.visitBinaryExpr(e.Left, e.Right, func(left Value, right Value) (Value, error) {
			return left.In(right)
		})

	case parser.NegExpr:
		leftVal, err := ir.visitExpr(e.Inner)
		if err != nil {
			return nil, err
		}
		return leftVal.Neg()

	case parser.NotExpr:
		leftVal, err := ir.visitExpr(e.Inner)
		if err != nil {
			return nil, err
		}
		return leftVal.Not()

	case parser.PosExpr:
		return ir.visitBinaryExpr(e.X, e.Y, func(xVal Value, yVal Value) (Value, error) {
			x, ok := xVal.(Number)
			if !ok {
				return nil, fmt.Errorf("type mismatch: expected pos(Number, Number)")
			}
			y, ok := yVal.(Number)
			if !ok {
				return nil, fmt.Errorf("type mismatch: expected pos(Number, Number)")
			}
			return Point{int(x + 0.5), int(y + 0.5)}, nil
		})

	case parser.MemberExpr:
		recvrVal, err := ir.visitExpr(e.Recvr)
		if err != nil {
			return nil, err
		}
		return recvrVal.Property(e.Member)

	case parser.IndexExpr:
		recvr, err := ir.visitExpr(e.Recvr)
		if err != nil {
			return nil, err
		}
		index, err := ir.visitExpr(e.Index)
		if err != nil {
			return nil, err
		}
		return recvr.Index(index)

	case parser.IndexRangeExpr:
		recvr, err := ir.visitExpr(e.Recvr)
		if err != nil {
			return nil, err
		}
		lower, err := ir.visitExpr(e.Lower)
		if err != nil {
			return nil, err
		}
		upper, err := ir.visitExpr(e.Upper)
		if err != nil {
			return nil, err
		}
		return recvr.IndexRange(lower, upper)

	case lang.Str:
		return Str(e), nil

	case lang.Boolean:
		return Boolean(e), nil

	case lang.Number:
		return Number(e), nil

	case lang.Color:
		return Color(e), nil

	case lang.Nil:
		return Nilval(e), nil

	case parser.IdentExpr:
		val, ok := ir.findIdent(string(e))
		if !ok {
			return nil, fmt.Errorf("identifier '%s' Not found", e)
		}
		return val, nil

	case parser.AtExpr:
		val, err := ir.visitExpr(e.Inner)
		if err != nil {
			return nil, err
		}
		pos, ok := val.(Point)
		if !ok {
			return nil, fmt.Errorf("")
		}
		return Color(ir.bitmap.GetPixel(pos.X, pos.Y)), nil

	case parser.InvokeExpr:
		args := []Value{}
		for _, arg := range e.Args {
			arg, err := ir.visitExpr(arg)
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		}
		return ir.invokeFunc(e.FuncName, args)

	case parser.KernelExpr:
		elementNumbers := make([]lang.Number, len(e.Elements))
		for i, element := range e.Elements {
			elementVal, err := ir.visitExpr(element)
			if err != nil {
				return nil, err
			}
			n, ok := elementVal.(Number)
			if !ok {
				return nil, fmt.Errorf("type mismatch: kernel expr expects number elements")
			}
			elementNumbers[i] = lang.Number(n)
		}
		rootOfLen := int(math.Sqrt(float64(len(elementNumbers))))
		return Kernel{
			Values: elementNumbers,
			Width:  rootOfLen,
			Height: rootOfLen,
		}, nil

	case parser.FunctionExpr:
		return Function{
			ParameterNames: e.ParameterNames,
			Body:           e.Body,
			closure:        ir.idents[2:], // omit constants and top scope - they are visible In any context
		}, nil

	case parser.HashMapExpr:
		h := make(HashMap)
		for _, entry := range e.Entries {
			key, err := ir.visitExpr(entry.Key)
			if err != nil {
				return nil, err
			}
			val, err := ir.visitExpr(entry.Value)
			if err != nil {
				return nil, err
			}
			h[key] = val
		}
		return h, nil

	case parser.ListExpr:
		l := List{
			Elements: make([]Value, len(e.Elements)),
		}
		for i, elem := range e.Elements {
			val, err := ir.visitExpr(elem)
			if err != nil {
				return nil, err
			}
			l.Elements[i] = val
		}
		return l, nil

	case parser.PipelineExpr:
		left, err := ir.visitExpr(e.Left)
		if err != nil {
			return nil, err
		}
		pipelineValueIdent := lexer.TokenTypeName(lexer.TTDollar)
		_ = ir.newIdent(pipelineValueIdent, left)
		defer ir.removeIdent(pipelineValueIdent)
		right, err := ir.visitExpr(e.Right)
		if err != nil {
			return nil, err
		}
		return right, nil
	}

	return nil, fmt.Errorf("unknown expression type %s", reflect.TypeOf(expr))
}

func (ir *interpreter) invokeFunc(name string, arguments []Value) (Value, error) {
	val, ok, err := ir.invokeBuiltinFunction(name, arguments)
	if err != nil {
		return nil, err
	}
	if ok {
		return val, nil
	}
	fval, ok := ir.findIdent(name)
	if !ok {
		return nil, fmt.Errorf("unknown identifier '%s'", name)
	}
	return ir.invokeFunctionExpr(name, fval, arguments)
}

func (ir *interpreter) invokeFunctionExpr(name string, val Value, arguments []Value) (Value, error) {
	fn, ok := val.(Function)
	if !ok {
		return nil, fmt.Errorf("%s is invoked like a function, but refers to a %s", name, reflect.TypeOf(val))
	}
	if len(arguments) != len(fn.ParameterNames) {
		return nil, fmt.Errorf("%s is invoked with %d arguments, but is declared with %d parameters", name, len(arguments), len(fn.ParameterNames))
	}

	prevScopeCount := ir.pushFunctionScope(fn.closure)
	defer ir.popFunctionScope(prevScopeCount)
	for i, argument := range arguments {
		ir.newIdent(fn.ParameterNames[i], argument)
	}

	if err := ir.visitStmtList(fn.Body); err != nil {
		if _, ok := err.(returnSignal); !ok { // return statement encountered
			return nil, err
		}
	}
	return ir.getReturnValue(), nil
}

func (ir *interpreter) invokeBuiltinFunction(name string, arguments []Value) (Value, bool, error) {
	fs, ok := functions[name]
	if !ok {
		return nil, false, nil
	}
	var err error
	for _, f := range fs {
		if err = validateArguments(arguments, f.params); err == nil {
			val, err := f.body(ir, arguments)
			return val, true, err
		}
	}

	buffer := strings.Builder{}
	for i, f := range fs {
		if i > 0 {
			buffer.WriteString("\n")
		}
		buffer.WriteString(signature(name, f))
	}
	return nil, false, fmt.Errorf("no fitting overload for function '%s': %s. possible overloads:\n%s", name, err, buffer.String())
}

func validateArguments(arguments []Value, params []reflect.Type) error {
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
				return fmt.Errorf("argument type mismatch At argument %d: expected %s, got %s", i, lastParam.Elem().Name(), reflect.TypeOf(arguments[i]).Name())
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
		return fmt.Errorf("argument type mismatch At argument %d: expected %s, got %s", i, params[i].Name(), reflect.TypeOf(arguments[i]).Name())
	}

	return nil
}

func hasMatchingType(v Value, typ reflect.Type) bool {
	return typ == valueType || reflect.TypeOf(v) == typ
}
