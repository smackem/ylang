package parser

import "github.com/smackem/ylang/internal/lexer"

type astNode interface {
}

type Statement interface {
	astNode
	Token() lexer.Token
}

type Expression interface {
	astNode
}

// Program is the complete, parsed ylang program.
type Program struct {
	Stmts []Statement
}

//////////////////////////////////////////////// statements

type StmtBase struct {
	tok lexer.Token
}

func (stmt StmtBase) Token() lexer.Token {
	return stmt.tok
}

type DeclStmt struct {
	StmtBase
	Ident string
	Rhs   Expression
}

type AssignStmt struct {
	StmtBase
	Ident string
	Rhs   Expression
}

type IndexedAssignStmt struct {
	StmtBase
	Ident string
	Index Expression
	Rhs   Expression
}

type PixelAssignStmt struct {
	StmtBase
	Lhs Expression
	Rhs Expression
}

type InvocationStmt struct {
	StmtBase
	Invocation Expression
}

type IfStmt struct {
	StmtBase
	Cond       Expression
	TrueStmts  []Statement
	FalseStmts []Statement // nillable
}

type ForStmt struct {
	StmtBase
	Ident      string
	Collection Expression
	Stmts      []Statement
}

type ForRangeStmt struct {
	StmtBase
	Ident string
	Lower Expression
	Upper Expression
	Step  Expression
	Stmts []Statement
}

type WhileStmt struct {
	StmtBase
	Cond  Expression
	Stmts []Statement
}

type YieldStmt struct {
	StmtBase
	Result Expression
}

type LogStmt struct {
	StmtBase
	Args []Expression
}

type ReturnStmt struct {
	StmtBase
	Result Expression
}

//////////////////////////////////////////////// expressions

type TernaryExpr struct {
	Cond        Expression
	TrueResult  Expression
	FalseResult Expression
}

type BinaryExpr struct {
	Left  Expression
	Right Expression
}

type PipelineExpr BinaryExpr
type OrExpr BinaryExpr
type AndExpr BinaryExpr
type EqExpr BinaryExpr
type NeqExpr BinaryExpr
type GtExpr BinaryExpr
type GeExpr BinaryExpr
type LtExpr BinaryExpr
type LeExpr BinaryExpr

type ConcatExpr BinaryExpr

type AddExpr BinaryExpr
type SubExpr BinaryExpr
type MulExpr BinaryExpr
type DivExpr BinaryExpr
type ModExpr BinaryExpr
type InExpr BinaryExpr

type UnaryExpr struct {
	Inner Expression
}

type NegExpr UnaryExpr
type NotExpr UnaryExpr

type PosExpr struct {
	X Expression
	Y Expression
}

type MemberExpr struct {
	Recvr  Expression
	Member string
}

type IndexExpr struct {
	Recvr Expression
	Index Expression
}

type IndexRangeExpr struct {
	Recvr Expression
	Lower Expression
	Upper Expression
}

type IdentExpr string
type AtExpr UnaryExpr

type InvokeExpr struct {
	FuncName string
	Args     []Expression
}

type KernelExpr struct {
	Elements []Expression
}

type FunctionExpr struct {
	ParameterNames []string
	Body           []Statement
}

type HashMapExpr struct {
	Entries []HashEntryExpr
}

type HashEntryExpr struct {
	Key   Expression
	Value Expression
}

type ListExpr struct {
	Elements []Expression
}
