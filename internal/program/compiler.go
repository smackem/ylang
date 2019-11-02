package program

import (
	"github.com/smackem/ylang/internal/interpreter"
	"github.com/smackem/ylang/internal/lexer"
	"github.com/smackem/ylang/internal/parser"
)

// Compile compiles the given source code into a executable Program.
func Compile(src string) (parser.Program, error) {
	tokens, err := lexer.Lex(src)
	if err != nil {
		return parser.Program{}, err
	}
	prog, err := parser.Parse(tokens, false)
	if err != nil {
		return parser.Program{}, err
	}
	return prog, nil
}

// Execute executes the Program against the specified Bitmap.
func Execute(prog parser.Program, bitmap interpreter.BitmapContext) error {
	return interpreter.Interpret(prog, bitmap)
}
