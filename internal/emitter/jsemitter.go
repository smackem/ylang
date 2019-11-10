package emitter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"github.com/smackem/ylang/internal/parser"
	"strings"
)

func EmitJS(program parser.Program) string {
	js := jsemitter{}
	js.pushBuffer()
	js.printf("function ylangExecute(surface) {")
	js.println()
	js.print("const Bounds = new Rect(0, 0, surface.width(), surface.height());")
	js.println()
	js.visitStmtList(program.Stmts)
	js.print("}")
	js.println()
	return js.popBuffer()
}

type jsemitter struct {
	bufs []*strings.Builder
}

func (js *jsemitter) pushBuffer() {
	js.bufs = append(js.bufs, &strings.Builder{})
}

func (js *jsemitter) popBuffer() string {
	if len(js.bufs) == 0 {
		panic("buffer stack underflow")
	}
	last := len(js.bufs) - 1
	buf := js.bufs[last]
	js.bufs = js.bufs[:last]
	return buf.String()
}

func (js *jsemitter) printf(format string, a ...interface{}) {
	js.bufs[len(js.bufs)-1].WriteString(fmt.Sprintf(format, a...))
}

func (js *jsemitter) print(s string) {
	js.bufs[len(js.bufs)-1].WriteString(s)
}

func (js *jsemitter) println() {
	js.bufs[len(js.bufs)-1].WriteRune('\n')
}

func (js *jsemitter) visitStmtList(stmts []parser.Statement) {
	for _, stmt := range stmts {
		js.visitStmt(stmt)
	}
}

func (js *jsemitter) visitStmt(stmt parser.Statement) {
	switch s := stmt.(type) {
	case parser.DeclStmt:
		js.printf("var %s = %s;", s.Ident, js.visitExpr(s.Rhs))

	case parser.AssignStmt:
		js.printf("%s = %s;", s.Ident, js.visitExpr(s.Rhs))

	case parser.IndexedAssignStmt:
		// lhs.setAt(index, rhs)
		js.printf("%s[%s] = %s;", s.Ident, js.visitExpr(s.Index), js.visitExpr(s.Rhs))

	case parser.PixelAssignStmt:
		js.printf("surface.setPixel(%s, %s);", js.visitExpr(s.Lhs), js.visitExpr(s.Rhs))

	case parser.InvocationStmt:
		js.printf("%s;", js.visitExpr(s.Invocation))

	case parser.IfStmt:
		js.printf("if (%s) {", js.visitExpr(s.Cond))
		js.println()
		js.visitStmtList(s.TrueStmts)
		js.print("}")
		if s.FalseStmts != nil {
			js.print(" else {")
			js.visitStmtList(s.FalseStmts)
			js.print("}")
		}

	case parser.ForStmt:
		js.printf("for (var %s of %s.iter()) {", s.Ident, js.visitExpr(s.Collection))
		js.println()
		js.visitStmtList(s.Stmts)
		js.print("}")

	case parser.ForRangeStmt:
		js.printf("for (var %s = %s; %s < %s; %s += %s) {",
			s.Ident, js.visitExpr(s.Lower),
			s.Ident, js.visitExpr(s.Upper),
			s.Ident, js.visitExpr(s.Step))
		js.println()
		js.visitStmtList(s.Stmts)
		js.print("}")

	case parser.WhileStmt:
		js.printf("while (%s) {", js.visitExpr(s.Cond))
		js.println()
		js.visitStmtList(s.Stmts)
		js.print("}")

	case parser.YieldStmt:
		js.print("UNSUPPORTED_STMT")

	case parser.LogStmt:
		js.print("console.log(")
		for i, expr := range s.Args {
			js.printf("%s", js.visitExpr(expr))
			if i < len(s.Args)-1 {
				js.print(", ")
			}
		}
		js.print(");")

	case parser.ReturnStmt:
		js.printf("return %s;", js.visitExpr(s.Result))
	}

	js.println()
}

func (js *jsemitter) visitExpr(expr parser.Expression) string {
	switch e := expr.(type) {
	case parser.TernaryExpr:
		return fmt.Sprintf("%s ? %s : %s", js.visitExpr(e.Cond), js.visitExpr(e.TrueResult), js.visitExpr(e.FalseResult))

	case parser.OrExpr:
		return fmt.Sprintf("%s || %s", js.visitExpr(e.Left), js.visitExpr(e.Right))

	case parser.AndExpr:
		return fmt.Sprintf("%s && %s", js.visitExpr(e.Left), js.visitExpr(e.Right))

	case parser.EqExpr:
		return fmt.Sprintf("%s == %s", js.visitExpr(e.Left), js.visitExpr(e.Right))

	case parser.NeqExpr:
		return fmt.Sprintf("%s != %s", js.visitExpr(e.Left), js.visitExpr(e.Right))

	case parser.GtExpr:
		return fmt.Sprintf("%s > %s", js.visitExpr(e.Left), js.visitExpr(e.Right))

	case parser.GeExpr:
		return fmt.Sprintf("%s >= %s", js.visitExpr(e.Left), js.visitExpr(e.Right))

	case parser.LtExpr:
		return fmt.Sprintf("%s < %s", js.visitExpr(e.Left), js.visitExpr(e.Right))

	case parser.LeExpr:
		return fmt.Sprintf("%s < %s", js.visitExpr(e.Left), js.visitExpr(e.Right))

	case parser.ConcatExpr:
		return fmt.Sprintf("%s.concat(%s)", js.visitExpr(e.Left), js.visitExpr(e.Right))

	case parser.AddExpr:
		return fmt.Sprintf("Op.add(%s, %s)", js.visitExpr(e.Left), js.visitExpr(e.Right))

	case parser.SubExpr:
		return fmt.Sprintf("Op.sub(%s, %s)", js.visitExpr(e.Left), js.visitExpr(e.Right))

	case parser.MulExpr:
		return fmt.Sprintf("Op.mul(%s, %s)", js.visitExpr(e.Left), js.visitExpr(e.Right))

	case parser.DivExpr:
		return fmt.Sprintf("Op.div(%s, %s)", js.visitExpr(e.Left), js.visitExpr(e.Right))

	case parser.ModExpr:
		return fmt.Sprintf("Op.mod(%s, %s)", js.visitExpr(e.Left), js.visitExpr(e.Right))

	case parser.InExpr:
		return fmt.Sprintf("(%s).contains(%s)", js.visitExpr(e.Right), js.visitExpr(e.Left))

	case parser.NegExpr:
		return fmt.Sprintf("Op.neg(%s)", js.visitExpr(e.Inner))

	case parser.NotExpr:
		return fmt.Sprintf("Op.not(%s)", js.visitExpr(e.Inner))

	case parser.PosExpr:
		return fmt.Sprintf("new Point(%s, %s)", js.visitExpr(e.X), js.visitExpr(e.Y))

	case parser.MemberExpr:
		return fmt.Sprintf("(%s).%s", js.visitExpr(e.Recvr), e.Member)

	case parser.IndexExpr:
		return fmt.Sprintf("(%s).at(%s)", js.visitExpr(e.Recvr), js.visitExpr(e.Index))

	case parser.IndexRangeExpr:
		return fmt.Sprintf("(%s).slice(%s, %s)", js.visitExpr(e.Recvr), js.visitExpr(e.Lower), js.visitExpr(e.Upper))

	case lang.Str:
		return fmt.Sprintf("\"%s\"", e)

	case lang.Boolean:
		if e {
			return "true"
		} else {
			return "false"
		}

	case lang.Number:
		return fmt.Sprintf("%f", e)

	case lang.Color:
		return fmt.Sprintf("new Color(%d, %d, %d, %d)", int(e.R), int(e.G), int(e.B), int(e.A))

	case lang.Nil:
		return "null"

	case parser.IdentExpr:
		return string(e)

	case parser.AtExpr:
		return fmt.Sprintf("surface.getPixel(%s)", js.visitExpr(e.Inner))

	case parser.InvokeExpr:
		return fmt.Sprintf("%s(%s)", e.FuncName, js.joinArgs(e.Args))

	case parser.KernelExpr:
		return fmt.Sprintf("new Kernel(%s)", js.joinArgs(e.Elements))

	case parser.FunctionExpr:
		params := strings.Join(e.ParameterNames, ",")
		js.pushBuffer()
		js.visitStmtList(e.Body)
		return fmt.Sprintf("function(%s) {\n%s}", params, js.popBuffer())

	case parser.HashMapExpr:
		js.pushBuffer()
		js.print("new HashMap({")
		js.println()
		for i, entry := range e.Entries {
			js.printf("%s: %s", js.visitExpr(entry.Key), js.visitExpr(entry.Value))
			if i < len(e.Entries)-1 {
				js.print(",")
			}
			js.println()
		}
		js.print("})")
		return js.popBuffer()

	case parser.ListExpr:
		return fmt.Sprintf("new List(%s)", js.joinArgs(e.Elements))

		//case parser.PipelineExpr:
		//	left, err := ir.visitExpr(e.Left)
		//	if err != nil {
		//		return nil, err
		//	}
		//	pipelineValueIdent := lexer.TokenTypeName(lexer.TTDollar)
		//	_ = ir.newIdent(pipelineValueIdent, left)
		//	defer ir.removeIdent(pipelineValueIdent)
		//	right, err := ir.visitExpr(e.Right)
		//	if err != nil {
		//		return nil, err
		//	}
		//	return right, nil
	}

	return "UNKNOWN_EXPR"
}

func (js *jsemitter) joinArgs(args []parser.Expression) string {
	buf := strings.Builder{}
	last := len(args) - 1
	for i, arg := range args {
		buf.WriteString(js.visitExpr(arg))
		if i < last {
			buf.WriteString(", ")
		}
	}
	return buf.String()
}
