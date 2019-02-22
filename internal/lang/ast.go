package lang

type astNode interface {
}

type statement interface {
	astNode
	getToken() token
}

type expression interface {
	astNode
}

//////////////////////////////////////////////// statements

type stmtBase struct {
	tok token
}

func (stmt stmtBase) getToken() token {
	return stmt.tok
}

type declStmt struct {
	stmtBase
	ident string
	rhs   expression
}

type assignStmt struct {
	stmtBase
	ident string
	rhs   expression
}

type indexedAssignStmt struct {
	stmtBase
	ident string
	index expression
	rhs   expression
}

type pixelAssignStmt struct {
	stmtBase
	lhs expression
	rhs expression
}

type invocationStmt struct {
	stmtBase
	invocation expression
}

type ifStmt struct {
	stmtBase
	cond       expression
	trueStmts  []statement
	falseStmts []statement // nillable
}

type forStmt struct {
	stmtBase
	ident      string
	collection expression
	stmts      []statement
}

type forRangeStmt struct {
	stmtBase
	ident string
	lower expression
	upper expression
	step  expression
	stmts []statement
}

type yieldStmt struct {
	stmtBase
	result expression
}

type logStmt struct {
	stmtBase
	parameters []expression
}

type bltStmt struct {
	stmtBase
	rect expression
}

type commitStmt struct {
	stmtBase
	rect expression // nillable
}

type returnStmt struct {
	stmtBase
	result expression
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

type indexExpr struct {
	recvr expression
	index expression
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

type functionExpr struct {
	parameterNames []string
	body           []statement
}
