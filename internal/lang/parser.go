package lang

import (
	"fmt"
	"math"
	"unicode"
)

func parse(input []token) (Program, error) {
	parser := parser{input: input, index: 0, symbols: make(map[astNode]token)}
	program, err := parser.parseProgram()
	if err != nil {
		tok := parser.current()
		return Program{nil, nil}, fmt.Errorf("at line %d, near '%s': %s", tok.LineNumber, tok.Lexeme, err)
	}
	return Program{program, parser.symbols}, nil
}

type parser struct {
	input   []token
	index   int
	symbols map[astNode]token
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
		return emptyToken, fmt.Errorf("expected %v, found %v", tokenTypeNames[tt], tok.Lexeme)
	}
	return tok, nil
}

func (p *parser) expect(tt tokenType) (token, error) {
	tok := p.next()
	if tok.Type != tt {
		return emptyToken, fmt.Errorf("expected %v, found %v", tokenTypeNames[tt], tok.Lexeme)
	}
	return tok, nil
}

func (p *parser) makeStmt(stmt statement) statement {
	p.symbols[&stmt] = p.current()
	return stmt
}

func (p *parser) makeExpr(expr expression) expression {
	p.symbols[&expr] = p.current()
	return expr
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
	case ttCommit:
		return p.parseCommit()
	}
	return nil, fmt.Errorf("unexpected token at statement begin: '%s'", tok)
}

func (p *parser) parseIdentStmt(ident string) (statement, error) {
	tok := p.next()
	switch tok.Type {
	case ttColonEq:
		return p.parseDeclaration(ident)
	case ttEq:
		return p.parseAssign(ident)
	}
	return nil, fmt.Errorf("unexpected token '%s' - expected %s or %s", tok, tokenTypeNames[ttColonEq], tokenTypeNames[ttEq])
}

func (p *parser) parseDeclaration(ident string) (statement, error) {
	rhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return p.makeStmt(declStmt{ident, rhs}), nil
}

func (p *parser) parseAssign(ident string) (statement, error) {
	if unicode.IsUpper(rune(ident[0])) {
		return nil, fmt.Errorf("identifier '%s' is a constant and cannot be assigned to", ident)
	}
	rhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return p.makeStmt(assignStmt{ident, rhs}), nil
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
	return p.makeStmt(pixelAssignStmt{lhs, rhs}), nil
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

		if p.current().Type == ttIf {
			p.next()
			stmt, err := p.parseIf()
			if err != nil {
				return nil, err
			}
			falseStmts = []statement{stmt}
		} else {
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
	}
	return p.makeStmt(ifStmt{cond, trueStmts, falseStmts}), nil
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
		return p.makeStmt(forRangeStmt{ident: identTok.Lexeme, lower: collection, upper: upper, stmts: stmts}), nil
	}
	return p.makeStmt(forStmt{identTok.Lexeme, collection, stmts}), nil
}

func (p *parser) parseYield() (statement, error) {
	result, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return p.makeStmt(yieldStmt{result}), nil
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
	return p.makeStmt(logStmt{parameters}), nil
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
	p.next()
	rect, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(ttRParen); err != nil {
		return nil, err
	}
	return p.makeStmt(bltStmt{rect}), nil
}

func (p *parser) parseCommit() (statement, error) {
	var rect expression
	if p.current().Type == ttLParen {
		p.next()
		var err error
		rect, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(ttRParen); err != nil {
			return nil, err
		}
	}
	return p.makeStmt(commitStmt{rect}), nil
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
		return p.makeExpr(ternaryExpr{cond, trueResult, falseResult}), nil
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
		return p.makeExpr(orExpr{left, right}), nil
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
		return p.makeExpr(andExpr{left, right}), nil
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
		return p.makeExpr(eqExpr{left, right}), nil
	case ttNeq:
		p.next()
		right, err := p.parseTermExpr()
		if err != nil {
			return nil, err
		}
		return p.makeExpr(neqExpr{left, right}), nil
	case ttGt:
		p.next()
		right, err := p.parseTermExpr()
		if err != nil {
			return nil, err
		}
		return p.makeExpr(gtExpr{left, right}), nil
	case ttGe:
		p.next()
		right, err := p.parseTermExpr()
		if err != nil {
			return nil, err
		}
		return p.makeExpr(geExpr{left, right}), nil
	case ttLt:
		p.next()
		right, err := p.parseTermExpr()
		if err != nil {
			return nil, err
		}
		return p.makeExpr(ltExpr{left, right}), nil
	case ttLe:
		p.next()
		right, err := p.parseTermExpr()
		if err != nil {
			return nil, err
		}
		return p.makeExpr(leExpr{left, right}), nil
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
		return p.makeExpr(addExpr{left, right}), nil
	case ttMinus:
		p.next()
		right, err := p.parseTermExpr()
		if err != nil {
			return nil, err
		}
		return p.makeExpr(subExpr{left, right}), nil
	case ttIn:
		p.next()
		right, err := p.parseTermExpr()
		if err != nil {
			return nil, err
		}
		return p.makeExpr(inExpr{left, right}), nil
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
		return p.makeExpr(mulExpr{left, right}), nil
	case ttSlash:
		p.next()
		right, err := p.parseProductExpr()
		if err != nil {
			return nil, err
		}
		return p.makeExpr(divExpr{left, right}), nil
	case ttPercent:
		p.next()
		right, err := p.parseProductExpr()
		if err != nil {
			return nil, err
		}
		return p.makeExpr(modExpr{left, right}), nil
	}
	return left, nil
}

func (p *parser) parseMoleculeExpr() (expression, error) {
	switch p.current().Type {
	case ttMinus:
		p.next()
		inner, err := p.parseMoleculeExpr()
		if err != nil {
			return nil, err
		}
		return p.makeExpr(negExpr{inner}), nil
	case ttNot:
		p.next()
		inner, err := p.parseMoleculeExpr()
		if err != nil {
			return nil, err
		}
		return p.makeExpr(notExpr{inner}), nil
	}

	atom, err := p.parseAtom()
	if err != nil {
		return nil, err
	}

	for {
		switch p.current().Type {
		case ttSemicolon:
			p.next()
			y, err := p.parseAtom()
			if err != nil {
				return nil, err
			}
			atom = p.makeExpr(posExpr{x: atom, y: y})
		case ttDot:
			p.next()
			memberTok, err := p.expect(ttIdent)
			if err != nil {
				return nil, err
			}
			atom = p.makeExpr(memberExpr{recvr: atom, member: memberTok.Lexeme})
		case ttLBracket:
			p.next()
			index, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			atom = p.makeExpr(indexExpr{recvr: atom, index: index})
			if _, err := p.expect(ttRBracket); err != nil {
				return nil, err
			}
		default:
			return atom, nil
		}
	}
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
		return p.makeExpr(String(tok.Lexeme)), nil
	case ttTrue:
		return p.makeExpr(Bool(true)), nil
	case ttFalse:
		return p.makeExpr(Bool(false)), nil
	case ttColor:
		return tok.parseColor(), nil
	case ttLBracket:
		return p.parseKernelAtom()
	}
	return nil, fmt.Errorf("unexpected symbol '%s'", tok.Lexeme)
}

func (p *parser) parseParenAtom() (expression, error) {
	inner, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(ttRParen); err != nil {
		return nil, err
	}
	return p.makeExpr(inner), nil
}

func (p *parser) parseAtAtom() (expression, error) {
	inner, err := p.parseAtom()
	if err != nil {
		return nil, err
	}
	return p.makeExpr(atExpr{inner}), nil
}

func (p *parser) parseIdentAtom(ident string) (expression, error) {
	if p.current().Type == ttLParen {
		p.next()
		parameters, err := p.parseParameterList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(ttRParen); err != nil {
			return nil, err
		}
		return p.makeExpr(invokeExpr{funcName: ident, parameters: parameters}), nil
	}
	return identExpr(ident), nil
}

func (p *parser) parseKernelAtom() (expression, error) {
	elements := []expression{}
	for {
		switch p.current().Type {
		case ttRBracket:
			width := math.Sqrt(float64(len(elements)))
			if width-math.Trunc(width) != 0 {
				return nil, fmt.Errorf("kernel defined in kernel expression must be quadratic")
			}
			p.next()
			return p.makeExpr(kernelExpr{elements}), nil
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
