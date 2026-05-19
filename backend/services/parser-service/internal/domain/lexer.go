package domain

import "fmt"

// Token represents a lexical token
type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

// TokenType represents token classifications
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenNumber
	TokenWord
	TokenMultiplier // x in "2x"
	TokenColon      // :
	TokenMinus      // - in range
	TokenUnit       // kg, lbs, etc
	TokenParenOpen  // (
	TokenParenClose // )
	TokenComma      // ,
	TokenWhitespace
	TokenNewline
	TokenInvalid
)

func (tt TokenType) String() string {
	switch tt {
	case TokenEOF:
		return "EOF"
	case TokenNumber:
		return "NUMBER"
	case TokenWord:
		return "WORD"
	case TokenMultiplier:
		return "MULTIPLIER"
	case TokenColon:
		return "COLON"
	case TokenMinus:
		return "MINUS"
	case TokenUnit:
		return "UNIT"
	case TokenParenOpen:
		return "PAREN_OPEN"
	case TokenParenClose:
		return "PAREN_CLOSE"
	case TokenComma:
		return "COMMA"
	case TokenWhitespace:
		return "WHITESPACE"
	case TokenNewline:
		return "NEWLINE"
	default:
		return "INVALID"
	}
}

// SetType represents the type of set
type SetType string

const (
	SetTypeWarmUp SetType = "WARM_UP"
	SetTypeFeeder SetType = "FEEDER"
	SetTypeWork   SetType = "WORK"
	SetTypeTopSet SetType = "TOP_SET"
)

// RepRange represents the target repetition range
type RepRange struct {
	Min int
	Max int
}

// SetGroup represents a parsed set in the DSL
type SetGroup struct {
	SetType      SetType
	PlannedSets  int
	TargetRepMin int
	TargetRepMax int
	Weight       float64
	Unit         string
	ExecutedReps []int
}

// ParsedWorkout represents the complete parsed DSL
type ParsedWorkout struct {
	ExerciseName string
	SetGroups    []SetGroup
	RawInput     string
}

// Lexer tokenizes input text
type Lexer struct {
	input  string
	pos    int
	line   int
	col    int
	tokens []Token
}

// NewLexer creates a new lexer
func NewLexer(input string) *Lexer {
	return &Lexer{
		input: input,
		pos:   0,
		line:  1,
		col:   1,
	}
}

// Tokenize performs lexical analysis
func (l *Lexer) Tokenize() ([]Token, error) {
	for l.pos < len(l.input) {
		l.skipWhitespace()
		if l.pos >= len(l.input) {
			break
		}

		ch := l.input[l.pos]

		switch ch {
		case ':':
			l.tokens = append(l.tokens, Token{Type: TokenColon, Value: ":", Line: l.line, Col: l.col})
			l.pos++
			l.col++
		case '-':
			l.tokens = append(l.tokens, Token{Type: TokenMinus, Value: "-", Line: l.line, Col: l.col})
			l.pos++
			l.col++
		case '(':
			l.tokens = append(l.tokens, Token{Type: TokenParenOpen, Value: "(", Line: l.line, Col: l.col})
			l.pos++
			l.col++
		case ')':
			l.tokens = append(l.tokens, Token{Type: TokenParenClose, Value: ")", Line: l.line, Col: l.col})
			l.pos++
			l.col++
		case ',':
			l.tokens = append(l.tokens, Token{Type: TokenComma, Value: ",", Line: l.line, Col: l.col})
			l.pos++
			l.col++
		case '\n':
			l.tokens = append(l.tokens, Token{Type: TokenNewline, Value: "\n", Line: l.line, Col: l.col})
			l.pos++
			l.line++
			l.col = 1
		case 'x':
			if l.pos > 0 && isDigit(l.input[l.pos-1]) {
				l.tokens = append(l.tokens, Token{Type: TokenMultiplier, Value: "x", Line: l.line, Col: l.col})
				l.pos++
				l.col++
			} else {
				l.readWord()
			}
		default:
			if isDigit(ch) {
				l.readNumber()
			} else if isLetter(ch) {
				l.readWord()
			} else {
				l.tokens = append(l.tokens, Token{Type: TokenInvalid, Value: string(ch), Line: l.line, Col: l.col})
				l.pos++
				l.col++
			}
		}
	}

	l.tokens = append(l.tokens, Token{Type: TokenEOF})
	return l.tokens, nil
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && (l.input[l.pos] == ' ' || l.input[l.pos] == '\t') {
		l.pos++
		l.col++
	}
}

func (l *Lexer) readNumber() {
	start := l.pos
	for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
		l.pos++
		l.col++
	}
	if l.pos < len(l.input) && l.input[l.pos] == '.' {
		l.pos++
		l.col++
		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			l.pos++
			l.col++
		}
	}
	value := l.input[start:l.pos]
	l.tokens = append(l.tokens, Token{Type: TokenNumber, Value: value, Line: l.line, Col: l.col - len(value)})
}

func (l *Lexer) readWord() {
	start := l.pos
	for l.pos < len(l.input) && (isLetter(l.input[l.pos]) || isDigit(l.input[l.pos]) || l.input[l.pos] == '_' || l.input[l.pos] == '-') {
		l.pos++
		l.col++
	}
	value := l.input[start:l.pos]

	// Classify word as unit or generic word
	if isUnit(value) {
		l.tokens = append(l.tokens, Token{Type: TokenUnit, Value: value, Line: l.line, Col: l.col - len(value)})
	} else {
		l.tokens = append(l.tokens, Token{Type: TokenWord, Value: value, Line: l.line, Col: l.col - len(value)})
	}
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isUnit(s string) bool {
	units := map[string]bool{
		"kg":  true,
		"lbs": true,
		"g":   true,
		"lb":  true,
	}
	return units[s]
}

// ParseError represents a parsing error
type ParseError struct {
	Message string
	Line    int
	Col     int
}

func (pe ParseError) Error() string {
	return fmt.Sprintf("parse error at line %d, col %d: %s", pe.Line, pe.Col, pe.Message)
}
