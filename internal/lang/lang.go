package lang

// Compile compiles the given source code into a executable Program.
func Compile(src string) (Program, error) {
	tokens, err := lex(src)
	if err != nil {
		return Program{}, err
	}
	prog, err := parse(tokens, false)
	if err != nil {
		return Program{}, err
	}
	return prog, nil
}

// Program is the compiled, executable ylang program.
type Program struct {
	stmts []statement
}

// Execute executes the Program against the specified Bitmap.
func (prog Program) Execute(bitmap BitmapContext) error {
	return interpret(prog, bitmap)
}

// BitmapContext is the surface a ylang Program works on.
type BitmapContext interface {
	GetPixel(x int, y int) Color
	SetPixel(x int, y int, color Color)
	SourceWidth() int
	SourceHeight() int
	TargetWidth() int
	TargetHeight() int
	Convolute(x, y, width, height int, kernel []Number) Color
	MapRed(x, y, width, height int, kernel []Number) []Number
	MapGreen(x, y, width, height int, kernel []Number) []Number
	MapBlue(x, y, width, height int, kernel []Number) []Number
	MapAlpha(x, y, width, height int, kernel []Number) []Number
	Blt(x, y, width, height int)
	ResizeTarget(width, height int)
	Flip() int // return imageID for Recall()
	Recall(imageID int) error
}
