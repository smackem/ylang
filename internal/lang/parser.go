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
	return declStmt{ident, rhs}, nil
}

func (p *parser) parseAssign(ident string) (statement, error) {
	rhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return assignStmt{ident, rhs}, nil
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
	return pixelAssignStmt{lhs, rhs}, nil
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
	if p.current().Type == ttElse {
		p.next()
		if _, err := p.expect(ttLBrace); err != nil {
			return nil, err
		}
		if p.current().Type != ttRBrace {
			var err error
			falseStmts, err = p.parseStmtList(ttRBrace)
			if err != nil {
				return nil, err
			}
		}
		if _, err := p.expect(ttRBrace); err != nil {
			return nil, err
		}
	}
	return ifStmt{cond, trueStmts, falseStmts}, nil
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
	return forStmt{identTok.Lexeme, collection, stmts}, nil
}

func (p *parser) parseYield() (statement, error) {
	result, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return yieldStmt{result}, nil
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
	return logStmt{parameters}, nil
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
		var err error
		rect, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(ttRParen); err != nil {
			return nil, err
		}
	}
	return bltStmt{rect}, nil
}

func (p *parser) parseExpr() (expression, error) {
	cond, err := p.parseOrExpr()
	if err != nil {
		return nil, err
	}

	if p.current().Type == ttQMark {
		p.next()
		trueResult, err := p.parseOrExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(ttColon); err != nil {
			return nil, err
		}
		falseResult, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return ternaryExpr{cond, trueResult, falseResult}, nil
	}

	return cond, nil
}

func (p *parser) parseOrExpr() (expression, error) {
	left, err := p.parseAndExpr()
	if err != nil {
		return nil, err
	}

	if p.current().Type == ttOr {
		p.next()
		right, err := p.parseOrExpr()
		if err != nil {
			return nil, err
		}
		return orExpr{left, right}, nil
	}
	return left, nil
}

func (p *parser) parseAndExpr() (expression, error) {
	left, err := p.parseCondExpr()
	if err != nil {
		return nil, err
	}

	if p.current().Type == ttAnd {
		p.next()
		right, err := p.parseAndExpr()
		if err != nil {
			return nil, err
		}
		return andExpr{left, right}, nil
	}
	return left, nil
}

func (p *parser) parseCondExpr() (expression, error) {
	left, err := p.parseTermExpr()
	if err != nil {
		return nil, err
	}

	switch p.current().Type {
	case ttEqEq:
		p.next()
		right, err := p.parseTermExpr()
		if err != nil {
			return nil, err
		}
		return eqExpr{left, right}, nil
	case ttNeq:
		p.next()
		right, err := p.parseTermExpr()
		if err != nil {
			return nil, err
		}
		return neqExpr{left, right}, nil
	case ttGt:
		p.next()
		right, err := p.parseTermExpr()
		if err != nil {
			return nil, err
		}
		return gtExpr{left, right}, nil
	case ttGe:
		p.next()
		right, err := p.parseTermExpr()
		if err != nil {
			return nil, err
		}
		return geExpr{left, right}, nil
	case ttLt:
		p.next()
		right, err := p.parseTermExpr()
		if err != nil {
			return nil, err
		}
		return ltExpr{left, right}, nil
	case ttLe:
		p.next()
		right, err := p.parseTermExpr()
		if err != nil {
			return nil, err
		}
		return leExpr{left, right}, nil
	}
	return left, nil
}

func (p *parser) parseTermExpr() (expression, error) {
	left, err := p.parseProductExpr()
	if err != nil {
		return nil, err
	}

	switch p.current().Type {
	case ttPlus:
		p.next()
		right, err := p.parseTermExpr()
		if err != nil {
			return nil, err
		}
		return addExpr{left, right}, nil
	case ttMinus:
		p.next()
		right, err := p.parseTermExpr()
		if err != nil {
			return nil, err
		}
		return subExpr{left, right}, nil
	case ttIn:
		p.next()
		right, err := p.parseTermExpr()
		if err != nil {
			return nil, err
		}
		return inExpr{left, right}, nil
	}
	return left, nil
}

func (p *parser) parseProductExpr() (expression, error) {
	left, err := p.parseMoleculeExpr()
	if err != nil {
		return nil, err
	}

	switch p.current().Type {
	case ttStar:
		p.next()
		right, err := p.parseProductExpr()
		if err != nil {
			return nil, err
		}
		return mulExpr{left, right}, nil
	case ttSlash:
		p.next()
		right, err := p.parseProductExpr()
		if err != nil {
			return nil, err
		}
		return divExpr{left, right}, nil
	case ttPercent:
		p.next()
		right, err := p.parseProductExpr()
		if err != nil {
			return nil, err
		}
		return modExpr{left, right}, nil
	}
	return left, nil
}

func (p *parser) parseMoleculeExpr() (expression, error) {
	switch p.current().Type {
	case ttMinus:
		p.next()
		inner, err := p.parseAtom()
		if err != nil {
			return nil, err
		}
		return negExpr{inner}, nil
	case ttNot:
		p.next()
		inner, err := p.parseAtom()
		if err != nil {
			return nil, err
		}
		return notExpr{inner}, nil
	}

	atom, err := p.parseAtom()
	if err != nil {
		return nil, err
	}

	if p.current().Type == ttSemicolon {
		p.next()
		y, err := p.parseAtom()
		if err != nil {
			return nil, err
		}
		return posExpr{x: atom, y: y}, nil
	}
	if p.current().Type == ttDot {
		p.next()
		memberTok, err := p.expect(ttIdent)
		if err != nil {
			return nil, err
		}
		return memberExpr{recvr: atom, member: memberTok.Lexeme}, nil
	}
	return atom, nil
}

func (p *parser) parseAtom() (expression, error) {
	tok := p.next()
	switch tok.Type {
	case ttLParen:
		return p.parseParenAtom()
	case ttAt:
		return p.parseAtAtom()
	case ttIdent:
		return p.parseIdentAtom(tok.Lexeme)
	case ttNumber:
		return tok.parseNumber(), nil
	case ttString:
		return stringExpr(tok.Lexeme), nil
	case ttTrue:
		return boolExpr(true), nil
	case ttFalse:
		return boolExpr(false), nil
	case ttLBracket:
		return p.parseKernelAtom()
	}
	return nil, fmt.Errorf("line %d: unexpected symbol %s", tok.LineNumber, tok.Lexeme)
}

func (p *parser) parseParenAtom() (expression, error) {
	inner, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(ttRParen); err != nil {
		return nil, err
	}
	return inner, nil
}

func (p *parser) parseAtAtom() (expression, error) {
	inner, err := p.parseAtom()
	if err != nil {
		return nil, err
	}
	return atExpr{inner}, nil
}

func (p *parser) parseIdentAtom(ident string) (expression, error) {
	if p.current().Type == ttLParen {
		parameters, err := p.parseParameterList()
		if err != nil {
			return nil, err
		}
		return invokeExpr{funcName: ident, parameters: parameters}, nil
	}
	return identExpr(ident), nil
}

func (p *parser) parseKernelAtom() (expression, error) {
	elements := []expression{}
	for {
		switch p.current().Type {
		case ttRBracket:
			p.next()
			return kernelExpr{elements}, nil
		case ttEOF:
			return nil, fmt.Errorf("unclosed kernel element list")
		}
		element, err := p.parseMoleculeExpr()
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)
	}
}
