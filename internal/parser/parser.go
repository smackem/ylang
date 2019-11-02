package parser

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"github.com/smackem/ylang/internal/lexer"
	"math"
	"unicode"
)

func Parse(input []lexer.Token, omitTokenInfo bool) (Program, error) {
	parser := parser{input: input, index: 0, omitTokenInfo: omitTokenInfo}
	program, err := parser.parseProgram()
	if err != nil {
		tok := parser.current()
		return Program{nil}, fmt.Errorf("at line %d, near '%s': %s", tok.LineNumber, tok.Lexeme, err)
	}
	return Program{program}, nil
}

type parser struct {
	input         []lexer.Token
	index         int
	omitTokenInfo bool
}

func (p parser) current() lexer.Token {
	if p.index >= len(p.input) {
		return lexer.EmptyToken
	}
	return p.input[p.index]
}

func (p *parser) next() lexer.Token {
	tok := p.current()
	p.index++
	return tok
}

func (p *parser) assert(tt lexer.TokenType) (lexer.Token, error) {
	tok := p.current()
	if tok.Type != tt {
		return lexer.EmptyToken, fmt.Errorf("expected %v, found %v", lexer.TokenTypeName(tt), tok.Lexeme)
	}
	return tok, nil
}

func (p *parser) expect(tt lexer.TokenType) (lexer.Token, error) {
	tok := p.next()
	if tok.Type != tt {
		return lexer.EmptyToken, fmt.Errorf("expected %v, found %v", lexer.TokenTypeName(tt), tok.Lexeme)
	}
	return tok, nil
}

func (p *parser) makeStmtBase() StmtBase {
	if p.omitTokenInfo {
		return StmtBase{}
	}
	return StmtBase{p.current()}
}

func (p *parser) parseProgram() (stmts []Statement, err error) {
	if stmts, err = p.parseStmtList(lexer.TTEOF); err != nil {
		return
	}
	_, err = p.assert(lexer.TTEOF)
	return
}

func (p *parser) parseStmtList(terminator lexer.TokenType) (stmts []Statement, err error) {
	for {
		switch p.current().Type {
		case terminator:
			return
		case lexer.TTEOF:
			return nil, fmt.Errorf("unclosed statement list")
		}
		var stmt Statement
		if stmt, err = p.parseStmt(); err != nil {
			return
		}
		stmts = append(stmts, stmt)
	}
}

func (p *parser) parseStmt() (stmt Statement, err error) {
	tok := p.next()
	switch tok.Type {
	case lexer.TTIdent:
		stmt, err = p.parseIdentStmt(tok.Lexeme)
	case lexer.TTAt:
		stmt, err = p.parsePixelAssign()
	case lexer.TTIf:
		stmt, err = p.parseIf()
	case lexer.TTFor:
		stmt, err = p.parseFor()
	case lexer.TTWhile:
		stmt, err = p.parseWhile()
	case lexer.TTYield:
		stmt, err = p.parseYield()
	case lexer.TTLog:
		stmt, err = p.parseLog()
	case lexer.TTReturn:
		stmt, err = p.parseReturn()
	default:
		stmt, err = nil, fmt.Errorf("unexpected Token at statement begin: '%s'", tok)
	}
	return
}

func (p *parser) parseIdentStmt(ident string) (Statement, error) {
	tok := p.next()
	switch tok.Type {
	case lexer.TTColonEq:
		return p.parseDeclaration(ident)
	case lexer.TTEq:
		return p.parseAssign(ident)
	case lexer.TTLBracket:
		return p.parseIndexedAssign(ident)
	case lexer.TTLParen:
		return p.parseInvocation(ident)
	}
	return nil, fmt.Errorf("unexpected Token '%s' - expected %s or %s", tok, lexer.TokenTypeName(lexer.TTColonEq), lexer.TokenTypeName(lexer.TTEq))
}

func (p *parser) parseInvocation(ident string) (Statement, error) {
	invocation, err := p.parseInvocationAtom(ident)
	if err != nil {
		return nil, err
	}
	return InvocationStmt{p.makeStmtBase(), invocation}, nil
}

func (p *parser) parseDeclaration(ident string) (Statement, error) {
	rhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return DeclStmt{p.makeStmtBase(), ident, rhs}, nil
}

func (p *parser) parseAssign(ident string) (Statement, error) {
	if unicode.IsUpper(rune(ident[0])) {
		return nil, fmt.Errorf("identifier '%s' is a constant and cannot be assigned to", ident)
	}
	rhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return AssignStmt{p.makeStmtBase(), ident, rhs}, nil
}

func (p *parser) parseIndexedAssign(ident string) (Statement, error) {
	indexExpr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.TTRBracket); err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.TTEq); err != nil {
		return nil, err
	}
	rhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return IndexedAssignStmt{p.makeStmtBase(), ident, indexExpr, rhs}, nil
}

func (p *parser) parsePixelAssign() (Statement, error) {
	lhs, err := p.parseAtom()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.TTEq); err != nil {
		return nil, err
	}
	rhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return PixelAssignStmt{p.makeStmtBase(), lhs, rhs}, nil
}

func (p *parser) parseIf() (Statement, error) {
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err = p.expect(lexer.TTLBrace); err != nil {
		return nil, err
	}

	var trueStmts []Statement
	if p.current().Type != lexer.TTRBrace {
		trueStmts, err = p.parseStmtList(lexer.TTRBrace)
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(lexer.TTRBrace); err != nil {
		return nil, err
	}

	var falseStmts []Statement
	if p.current().Type == lexer.TTElse {
		p.next()

		if p.current().Type == lexer.TTIf {
			p.next()
			stmt, err := p.parseIf()
			if err != nil {
				return nil, err
			}
			falseStmts = []Statement{stmt}
		} else {
			if _, err := p.expect(lexer.TTLBrace); err != nil {
				return nil, err
			}
			if p.current().Type != lexer.TTRBrace {
				var err error
				falseStmts, err = p.parseStmtList(lexer.TTRBrace)
				if err != nil {
					return nil, err
				}
			}
			if _, err := p.expect(lexer.TTRBrace); err != nil {
				return nil, err
			}
		}
	}
	return IfStmt{p.makeStmtBase(), cond, trueStmts, falseStmts}, nil
}

func (p *parser) parseFor() (Statement, error) {
	identTok, err := p.expect(lexer.TTIdent)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.TTIn); err != nil {
		return nil, err
	}
	collection, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	var upper Expression
	var step Expression
	if p.current().Type == lexer.TTDotDot {
		p.next()
		upper, err = p.parseExpr()
		if err != nil {
			return nil, err
		}

		if p.current().Type == lexer.TTDotDot {
			p.next()
			step = upper
			upper, err = p.parseExpr()
			if err != nil {
				return nil, err
			}
		} else {
			step = lang.Number(1)
		}
	}

	if _, err := p.expect(lexer.TTLBrace); err != nil {
		return nil, err
	}
	stmts, err := p.parseStmtList(lexer.TTRBrace)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.TTRBrace); err != nil {
		return nil, err
	}

	if upper != nil {
		return ForRangeStmt{
			Ident: identTok.Lexeme,
			Lower: collection,
			Upper: upper,
			Step:  step,
			Stmts: stmts}, nil
	}
	return ForStmt{p.makeStmtBase(), identTok.Lexeme, collection, stmts}, nil
}

func (p *parser) parseWhile() (Statement, error) {
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.TTLBrace); err != nil {
		return nil, err
	}
	stmts, err := p.parseStmtList(lexer.TTRBrace)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.TTRBrace); err != nil {
		return nil, err
	}
	return WhileStmt{p.makeStmtBase(), cond, stmts}, nil
}

func (p *parser) parseYield() (Statement, error) {
	result, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return YieldStmt{p.makeStmtBase(), result}, nil
}

func (p *parser) parseLog() (Statement, error) {
	if _, err := p.expect(lexer.TTLParen); err != nil {
		return nil, err
	}
	args, err := p.parseArgumentList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.TTRParen); err != nil {
		return nil, err
	}
	return LogStmt{p.makeStmtBase(), args}, nil
}

func (p *parser) parseArgumentList() ([]Expression, error) {
	args := []Expression{}
	for {
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if p.current().Type == lexer.TTComma {
			p.next()
		} else {
			break
		}
	}
	return args, nil
}

func (p *parser) parseReturn() (Statement, error) {
	result, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return ReturnStmt{p.makeStmtBase(), result}, nil
}

func (p *parser) parseExpr() (Expression, error) {
	cond, err := p.parsePipelineExpr()
	if err != nil {
		return nil, err
	}

	if p.current().Type == lexer.TTQMark {
		p.next()
		trueResult, err := p.parsePipelineExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.TTColon); err != nil {
			return nil, err
		}
		falseResult, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return TernaryExpr{cond, trueResult, falseResult}, nil
	}

	return cond, nil
}

func (p *parser) parsePipelineExpr() (Expression, error) {
	left, err := p.parseOrExpr()
	if err != nil {
		return nil, err
	}

	if p.current().Type == lexer.TTPipe {
		p.next()
		right, err := p.parsePipelineExpr()
		if err != nil {
			return nil, err
		}
		return PipelineExpr{left, right}, nil
	}

	return left, nil
}

func (p *parser) parseOrExpr() (Expression, error) {
	left, err := p.parseAndExpr()
	if err != nil {
		return nil, err
	}

	for {
		switch p.current().Type {
		case lexer.TTOr:
			p.next()
			right, err := p.parseAndExpr()
			if err != nil {
				return nil, err
			}
			left = OrExpr{left, right}
		default:
			return left, nil
		}
	}
}

func (p *parser) parseAndExpr() (Expression, error) {
	left, err := p.parseCondExpr()
	if err != nil {
		return nil, err
	}

	for {
		switch p.current().Type {
		case lexer.TTAnd:
			p.next()
			right, err := p.parseCondExpr()
			if err != nil {
				return nil, err
			}
			left = AndExpr{left, right}
		default:
			return left, nil
		}
	}
}

func (p *parser) parseCondExpr() (Expression, error) {
	left, err := p.parseConcatExpr()
	if err != nil {
		return nil, err
	}

	switch p.current().Type {
	case lexer.TTEqEq:
		p.next()
		right, err := p.parseConcatExpr()
		if err != nil {
			return nil, err
		}
		return EqExpr{left, right}, nil
	case lexer.TTNeq:
		p.next()
		right, err := p.parseConcatExpr()
		if err != nil {
			return nil, err
		}
		return NeqExpr{left, right}, nil
	case lexer.TTGt:
		p.next()
		right, err := p.parseConcatExpr()
		if err != nil {
			return nil, err
		}
		return GtExpr{left, right}, nil
	case lexer.TTGe:
		p.next()
		right, err := p.parseConcatExpr()
		if err != nil {
			return nil, err
		}
		return GeExpr{left, right}, nil
	case lexer.TTLt:
		p.next()
		right, err := p.parseConcatExpr()
		if err != nil {
			return nil, err
		}
		return LtExpr{left, right}, nil
	case lexer.TTLe:
		p.next()
		right, err := p.parseConcatExpr()
		if err != nil {
			return nil, err
		}
		return LeExpr{left, right}, nil
	}
	return left, nil
}

func (p *parser) parseConcatExpr() (Expression, error) {
	left, err := p.parseTupleExpr()
	if err != nil {
		return nil, err
	}

	for {
		switch p.current().Type {
		case lexer.TTColonColon:
			p.next()
			right, err := p.parseTupleExpr()
			if err != nil {
				return nil, err
			}
			left = ConcatExpr{left, right}
		default:
			return left, nil
		}
	}
}

func (p *parser) parseTupleExpr() (Expression, error) {
	left, err := p.parseTermExpr()
	if err != nil {
		return nil, err
	}

	if p.current().Type == lexer.TTSemicolon {
		p.next()
		right, err := p.parseTermExpr()
		if err != nil {
			return nil, err
		}
		return PosExpr{X: left, Y: right}, nil
	}

	return left, nil
}

func (p *parser) parseTermExpr() (Expression, error) {
	left, err := p.parseProductExpr()
	if err != nil {
		return nil, err
	}

	for {
		switch p.current().Type {
		case lexer.TTPlus:
			p.next()
			right, err := p.parseProductExpr()
			if err != nil {
				return nil, err
			}
			left = AddExpr{left, right}
		case lexer.TTMinus:
			p.next()
			right, err := p.parseProductExpr()
			if err != nil {
				return nil, err
			}
			left = SubExpr{left, right}
		case lexer.TTIn:
			p.next()
			right, err := p.parseProductExpr()
			if err != nil {
				return nil, err
			}
			left = InExpr{left, right}
		default:
			return left, nil
		}
	}
}

func (p *parser) parseProductExpr() (Expression, error) {
	left, err := p.parseMoleculeExpr()
	if err != nil {
		return nil, err
	}

	for {
		switch p.current().Type {
		case lexer.TTStar:
			p.next()
			right, err := p.parseMoleculeExpr()
			if err != nil {
				return nil, err
			}
			left = MulExpr{left, right}
		case lexer.TTSlash:
			p.next()
			right, err := p.parseMoleculeExpr()
			if err != nil {
				return nil, err
			}
			left = DivExpr{left, right}
		case lexer.TTPercent:
			p.next()
			right, err := p.parseMoleculeExpr()
			if err != nil {
				return nil, err
			}
			left = ModExpr{left, right}
		default:
			return left, nil
		}
	}
}

func (p *parser) parseMoleculeExpr() (Expression, error) {
	switch p.current().Type {
	case lexer.TTMinus:
		p.next()
		inner, err := p.parseMoleculeExpr()
		if err != nil {
			return nil, err
		}
		return NegExpr{inner}, nil
	case lexer.TTNot:
		p.next()
		inner, err := p.parseMoleculeExpr()
		if err != nil {
			return nil, err
		}
		return NotExpr{inner}, nil
	}

	atom, err := p.parseAtom()
	if err != nil {
		return nil, err
	}

	for {
		switch p.current().Type {
		case lexer.TTDot:
			p.next()
			memberTok, err := p.expect(lexer.TTIdent)
			if err != nil {
				return nil, err
			}
			atom = MemberExpr{Recvr: atom, Member: memberTok.Lexeme}
		case lexer.TTLBracket:
			p.next()
			index, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if p.current().Type == lexer.TTDotDot {
				p.next()
				upper, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				atom = IndexRangeExpr{Recvr: atom, Lower: index, Upper: upper}
			} else {
				atom = IndexExpr{Recvr: atom, Index: index}
			}
			if _, err := p.expect(lexer.TTRBracket); err != nil {
				return nil, err
			}
		default:
			return atom, nil
		}
	}
}

func (p *parser) parseAtom() (Expression, error) {
	tok := p.next()
	switch tok.Type {
	case lexer.TTLParen:
		return p.parseParenAtom()
	case lexer.TTAt:
		return p.parseAtAtom()
	case lexer.TTIdent:
		return p.parseIdentAtom(tok.Lexeme)
	case lexer.TTNumber:
		return tok.ParseNumber(), nil
	case lexer.TTString:
		return lang.Str(tok.ParseString()), nil
	case lexer.TTTrue:
		return lang.TrueVal, nil
	case lexer.TTFalse:
		return lang.FalseVal, nil
	case lexer.TTNil:
		return lang.NilVal, nil
	case lexer.TTDollar:
		return IdentExpr(tok.Lexeme), nil
	case lexer.TTColor:
		return tok.ParseColor(), nil
	case lexer.TTPipe:
		return p.parseKernelAtom()
	case lexer.TTFn:
		return p.parseFunctionDef()
	case lexer.TTLBrace:
		return p.parseHashMap()
	case lexer.TTLBracket:
		return p.parseList()
	}
	return nil, fmt.Errorf("unexpected symbol '%s'", tok.Lexeme)
}

func (p *parser) parseParenAtom() (Expression, error) {
	inner, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.TTRParen); err != nil {
		return nil, err
	}
	return inner, nil
}

func (p *parser) parseAtAtom() (Expression, error) {
	inner, err := p.parseAtom()
	if err != nil {
		return nil, err
	}
	return AtExpr{inner}, nil
}
func (p *parser) parseIdentAtom(ident string) (Expression, error) {

	if p.current().Type == lexer.TTLParen {
		p.next()
		return p.parseInvocationAtom(ident)
	}
	return IdentExpr(ident), nil
}

func (p *parser) parseInvocationAtom(ident string) (Expression, error) {
	var args []Expression
	var err error
	if p.current().Type == lexer.TTRParen {
		args = nil
	} else {
		args, err = p.parseArgumentList()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(lexer.TTRParen); err != nil {
		return nil, err
	}
	return InvokeExpr{FuncName: ident, Args: args}, nil
}

func (p *parser) parseKernelAtom() (Expression, error) {
	elements := []Expression{}
	for {
		if p.current().Type == lexer.TTPipe {
			width := math.Sqrt(float64(len(elements)))
			if width-math.Trunc(width) != 0 {
				return nil, fmt.Errorf("kernel defined in kernel expression must be quadratic")
			}
			p.next()
			return KernelExpr{elements}, nil
		}
		element, err := p.parseMoleculeExpr()
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)
	}
}

func (p *parser) parseFunctionDef() (Expression, error) {
	if _, err := p.expect(lexer.TTLParen); err != nil {
		return nil, err
	}
	var parameterNames []string
	var err error
	if p.current().Type == lexer.TTRParen {
		p.next()
	} else {
		parameterNames, err = p.parseIdentList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.TTRParen); err != nil {
			return nil, err
		}
	}

	if p.current().Type == lexer.TTArrow {
		p.next()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return FunctionExpr{
			ParameterNames: parameterNames,
			Body: []Statement{
				ReturnStmt{
					StmtBase: p.makeStmtBase(),
					Result:   expr,
				},
			},
		}, nil
	}

	if _, err := p.expect(lexer.TTLBrace); err != nil {
		return nil, err
	}
	body, err := p.parseStmtList(lexer.TTRBrace)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.TTRBrace); err != nil {
		return nil, err
	}

	return FunctionExpr{
		ParameterNames: parameterNames,
		Body:           body,
	}, nil
}

func (p *parser) parseIdentList() ([]string, error) {
	idents := []string{}
	for {
		tok, err := p.expect(lexer.TTIdent)
		if err != nil {
			return nil, err
		}

		idents = append(idents, tok.Lexeme)

		if p.current().Type == lexer.TTComma {
			p.next()
		} else {
			break
		}
	}
	return idents, nil
}

func (p *parser) parseHashMap() (Expression, error) {
	hashMap := HashMapExpr{}
	for {
		var key Expression
		var err error

		if p.current().Type == lexer.TTRBrace {
			p.next()
			return hashMap, nil
		}

		tok := p.current()
		if tok.Type == lexer.TTIdent {
			key = lang.Str(tok.Lexeme)
			p.next()
		} else {
			key, err = p.parseExpr()
			if err != nil {
				return nil, err
			}
		}

		if _, err := p.expect(lexer.TTColon); err != nil {
			return nil, err
		}

		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		hashMap.Entries = append(hashMap.Entries, HashEntryExpr{key, value})

		if p.current().Type == lexer.TTComma {
			p.next()
		} else {
			break
		}
	}
	if _, err := p.expect(lexer.TTRBrace); err != nil {
		return nil, err
	}
	return hashMap, nil
}

func (p *parser) parseList() (Expression, error) {
	list := ListExpr{}
	for {
		if p.current().Type == lexer.TTRBracket {
			p.next()
			return list, nil
		}

		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		list.Elements = append(list.Elements, expr)

		if p.current().Type == lexer.TTComma {
			p.next()
		} else {
			break
		}
	}
	if _, err := p.expect(lexer.TTRBracket); err != nil {
		return nil, err
	}
	return list, nil
}
