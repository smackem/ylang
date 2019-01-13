package lang

import "fmt"

func parse(input []token) ([]statement, error) {
	parser := parser{input: input, index: 0}
	return parser.parseProgram()
}

type parser struct {
	input []token
	index int
}

func (p parser) current() token {
	if p.index >= len(p.input) {
		return emptyToken
	}
	return p.input[p.index]
}

func (p *parser) next() token {
	tok := p.current()
	p.index++
	return tok
}

func (p *parser) assert(tt tokenType) (token, error) {
	tok := p.current()
	if tok.Type != tt {
		return emptyToken, fmt.Errorf("line %d: expected %v, found %v", tok.LineNumber, tokenTypeNames[tt], tok.Lexeme)
	}
	return tok, nil
}

func (p *parser) expect(tt tokenType) (token, error) {
	tok := p.next()
	if tok.Type != tt {
		return emptyToken, fmt.Errorf("line %d: expected %v, found %v", tok.LineNumber, tokenTypeNames[tt], tok.Lexeme)
	}
	return tok, nil
}

func (p *parser) parseProgram() (stmts []statement, err error) {
	if stmts, err = p.parseStmtList(ttEOF); err != nil {
		return
	}
	_, err = p.assert(ttEOF)
	return
}

func (p *parser) parseStmtList(terminator tokenType) (stmts []statement, err error) {
	for {
		switch p.current().Type {
		case terminator:
			return
		case ttEOF:
			return nil, fmt.Errorf("unclosed statement list")
		}
		var stmt statement
		if stmt, err = p.parseStmt(); err != nil {
			return
		}
		stmts = append(stmts, stmt)
	}
}

func (p *parser) parseStmt() (statement, error) {
	tok := p.next()
	switch tok.Type {
	case ttIdent:
		return p.parseIdentStmt(tok.Lexeme)
	case ttAt:
		return p.parsePixelAssign()
	case ttIf:
		return p.parseIf()
	case ttFor:
		return p.parseFor()
	case ttYield:
		return p.parseYield()
	case ttLog:
		return p.parseLog()
	case ttBlt:
		return p.parseBlt()
	}
	return nil, fmt.Errorf("unexpected token at statement begin: %s", tok)
}

func (p *parser) parseIdentStmt(ident string) (statement, error) {
	tok := p.next()
	switch tok.Type {
	case ttColonEq:
		return p.parseDeclaration(ident)
	case ttEq:
		return p.parseAssign(ident)
	}
	return nil, fmt.Errorf("unexpected token %s - expected %s or %s", tok, tokenTypeNames[ttColonEq], tokenTypeNames[ttEq])
}

func (p *parser) parseDeclaration(ident string) (statement, error) {
	rhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return declStmt{ident: ident, rhs: rhs}, nil
}

func (p *parser) parseAssign(ident string) (statement, error) {
	rhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return assignStmt{ident: ident, rhs: rhs}, nil
}

func (p *parser) parsePixelAssign() (statement, error) {
	lhs, err := p.parseAtom()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(ttEq); err != nil {
		return nil, err
	}
	rhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return pixelAssignStmt{lhs: lhs, rhs: rhs}, nil
}

func (p *parser) parseIf() (statement, error) {
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err = p.expect(ttLBrace); err != nil {
		return nil, err
	}

	var trueStmts []statement
	if p.current().Type != ttRBrace {
		trueStmts, err = p.parseStmtList(ttRBrace)
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(ttRBrace); err != nil {
		return nil, err
	}

	var falseStmts []statement
	if p.next().Type == ttElse {
		if _, err := p.expect(ttLBrace); err != nil {
			return nil, err
		}
		if p.current().Type != ttRBrace {
			falseStmts, err := p.parseStmtList(ttRBrace)
			if err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(ttRBrace); err != nil {
			return nil, err
		}
	}
	return ifStmt{cond: cond, trueStmts: trueStmts, falseStmts: falseStmts}, nil
}

func (p *parser) parseFor() (statement, error) {
	identTok, err := p.expect(ttIdent)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(ttIn); err != nil {
		return nil, err
	}
	collection, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	var upper statement
	if p.current().Type == ttDotDot {
		p.next()
		upper, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.expect(ttLBrace); err != nil {
		return nil, err
	}
	stmts, err := p.parseStmtList(ttRBrace)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(ttRBrace); err != nil {
		return nil, err
	}

	if upper != nil {
		return forRangeStmt{ident: identTok.Lexeme, lower: collection, upper: upper, stmts: stmts}, nil
	}
	return forStmt{ident: identTok.Lexeme, collection: collection, stmts: stmts}, nil
}

func (p *parser) parseYield() (statement, error) {
	result, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return yieldStmt{result: result}, nil
}

func (p *parser) parseLog() (statement, error) {
	if _, err := p.expect(ttLParen); err != nil {
		return nil, err
	}
	parameters, err := p.parseParameterList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(ttRParen); err != nil {
		return nil, err
	}
	return logStmt{parameters: parameters}, nil
}

func (p *parser) parseParameterList() ([]expression, error) {
	parameters := []expression{}
	for {
		parameter, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, parameter)
		if p.current().Type == ttComma {
			p.next()
		} else {
			break
		}
	}
	return parameters, nil
}

func (p *parser) parseBlt() (statement, error) {
	var rect expression
	if p.next().Type == ttLParen {
		rect, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(ttRParen); err != nil {
			return nil, err
		}
	}
	return bltStmt{rect: rect}, nil
}

func (p *parser) parseExpr() (expression, error) {
	if err := p.parseOrExpr(); err != nil {
		return nil, err
	}

	if p.current().Type == ttQMark {
		p.next()
		if err := p.parseOrExpr(); err != nil {
			return nil, err
		}
		if _, err := p.expect(ttColon); err != nil {
			return nil, err
		}
		falseResult, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (p *parser) parseOrExpr() error {
	if err := p.parseAndExpr(); err != nil {
		return err
	}

	if p.current().Type == ttOr {
		p.next()
		if err := p.parseOrExpr(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseAndExpr() error {
	if err := p.parseCondExpr(); err != nil {
		return err
	}

	if p.current().Type == ttAnd {
		p.next()
		if err := p.parseAndExpr(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseCondExpr() error {
	if err := p.parseTermExpr(); err != nil {
		return err
	}

	switch p.current().Type {
	case ttEqEq:
		p.next()
		if err := p.parseTermExpr(); err != nil {
			return err
		}
	case ttNeq:
		p.next()
		if err := p.parseTermExpr(); err != nil {
			return err
		}
	case ttGt:
		p.next()
		if err := p.parseTermExpr(); err != nil {
			return err
		}
	case ttGe:
		p.next()
		if err := p.parseTermExpr(); err != nil {
			return err
		}
	case ttLt:
		p.next()
		if err := p.parseTermExpr(); err != nil {
			return err
		}
	case ttLe:
		p.next()
		if err := p.parseTermExpr(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseTermExpr() error {
	if err := p.parseProductExpr(); err != nil {
		return err
	}

	switch p.current().Type {
	case ttPlus:
		p.next()
		if err := p.parseTermExpr(); err != nil {
			return err
		}
	case ttMinus:
		p.next()
		if err := p.parseTermExpr(); err != nil {
			return err
		}
	case ttIn:
		p.next()
		if err := p.parseTermExpr(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseProductExpr() error {
	if err := p.parseMoleculeExpr(); err != nil {
		return err
	}

	switch p.current().Type {
	case ttStar:
		p.next()
		if err := p.parseProductExpr(); err != nil {
			return err
		}
	case ttSlash:
		p.next()
		if err := p.parseProductExpr(); err != nil {
			return err
		}
	case ttPercent:
		p.next()
		if err := p.parseProductExpr(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseMoleculeExpr() error {
	switch p.current().Type {
	case ttMinus:
		p.next()
		return p.parseAtom()
	case ttNot:
		p.next()
		return p.parseAtom()
	}

	if err := p.parseAtom(); err != nil {
		return err
	}

	if p.current().Type == ttSemicolon {
		p.next()
		if err := p.parseAtom(); err != nil {
			return err
		}
	}
	if p.current().Type == ttDot {
		p.next()
		if _, err := p.expect(ttIdent); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) parseAtom() error {
	tok := p.next()
	switch tok.Type {
	case ttLParen:
		return p.parseParenAtom()
	case ttAt:
		return p.parseAtAtom()
	case ttIdent:
		return p.parseIdentAtom()
	case ttNumber:
		return nil
	case ttString:
		return nil
	case ttTrue:
		return nil
	case ttFalse:
		return nil
	case ttLBracket:
		return p.parseKernelAtom()
	}
	return fmt.Errorf("line %d: unexpected symbol %s", tok.LineNumber, tok.Lexeme)
}

func (p *parser) parseParenAtom() error {
	if err := p.parseExpr(); err != nil {
		return err
	}
	if _, err := p.expect(ttRParen); err != nil {
		return err
	}
	return nil
}

func (p *parser) parseAtAtom() error {
	return p.parseAtom()
}

func (p *parser) parseIdentAtom() error {
	if p.current().Type == ttLParen { // function call
		if err := p.parseParameterList(); err != nil {
			return nil
		}
	}
	return nil
}

func (p *parser) parseKernelAtom() error {
	for {
		switch p.current().Type {
		case ttRBracket:
			p.next()
			return nil
		case ttEOF:
			return fmt.Errorf("unclosed kernel element list")
		}
		if err := p.parseMoleculeExpr(); err != nil {
			return err
		}
	}
}
