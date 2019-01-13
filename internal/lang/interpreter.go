package lang

type interpreter struct {
}

func (i interpreter) visitProgram(stmts []statement) error {
	return nil
}

func (i interpreter) visitStmt(stmt statement) error {
	switch stmt.(type) {
	case declStmt:
	case assignStmt:
	case pixelAssignStmt:
	case ifStmt:
	case forStmt:
	case forRangeStmt:
	case yieldStmt:
	case logStmt:
	case bltStmt:
	}
	return nil
}
