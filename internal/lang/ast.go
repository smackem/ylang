package lang

type statement interface {
}

type expression interface {
}

//////////////////////////////////////////////// statements

type declStmt struct {
	ident string
	rhs   expression
}

type assignStmt struct {
	ident string
	rhs   expression
}

type pixelAssignStmt struct {
	lhs expression
	rhs expression
}

type ifStmt struct {
	cond       expression
	trueStmts  []statement
	falseStmts []statement // nillable
}

type forStmt struct {
	ident      string
	collection expression
	stmts      []statement
}

type forRangeStmt struct {
	ident string
	lower expression
	upper expression
	stmts []statement
}

type yieldStmt struct {
	result expression
}

type logStmt struct {
	parameters []expression
}

type bltStmt struct {
	rect expression // nillable
}

//////////////////////////////////////////////// expressions

type ternaryExpr struct {
	cond        expression
	trueResult  expression
	falseResult expression
}

type binaryExpr struct {
	left  expression
	right expression
}

type orExpr binaryExpr
type andExpr binaryExpr

type eqExpr binaryExpr
type neqExpr binaryExpr
type gtExpr binaryExpr
type geExpr binaryExpr
type ltExpr binaryExpr
type leExpr binaryExpr

type addExpr binaryExpr
type subExpr binaryExpr
type mulExpr binaryExpr
type divExpr binaryExpr
type modExpr binaryExpr
type inExpr binaryExpr

type unaryExpr struct {
	inner expression
}

type negExpr unaryExpr
type notExpr unaryExpr
type posExpr struct {
	x expression
	y expression
}
type memberExpr struct {
	recvr  expression
	member string
}
type identExpr string
type atExpr unaryExpr
type invokeExpr struct {
	funcName   string
	parameters []expression
}
type kernelExpr struct {
	elements []expression
}
