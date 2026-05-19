package domain

import (
	"strconv"
	"strings"
)

// Parser parses tokenized DSL into a structured AST
type Parser struct {
	tokens []Token
	pos    int
	errors []ParseError
}

// NewParser creates a new parser
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
		errors: []ParseError{},
	}
}

// Parse performs syntax analysis and returns structured data
func (p *Parser) Parse() (*ParsedWorkout, error) {
	workout := &ParsedWorkout{
		SetGroups: []SetGroup{},
	}

	// First line is exercise name
	exerciseName := p.parseExerciseName()
	if exerciseName == "" {
		return nil, ParseError{Message: "exercise name required", Line: 1, Col: 1}
	}
	workout.ExerciseName = exerciseName

	// Parse set groups
	for p.pos < len(p.tokens) {
		if p.peekToken().Type == TokenEOF {
			break
		}

		setGroup, err := p.parseSetGroup()
		if err != nil {
			p.errors = append(p.errors, err.(ParseError))
			p.skipToNextLine()
			continue
		}

		if setGroup != nil {
			workout.SetGroups = append(workout.SetGroups, *setGroup)
		}
	}

	if len(p.errors) > 0 {
		return nil, p.errors[0]
	}

	return workout, nil
}

// parseExerciseName extracts the exercise name from the first line
func (p *Parser) parseExerciseName() string {
	var name []string

	for p.pos < len(p.tokens) && p.peekToken().Type != TokenNewline && p.peekToken().Type != TokenEOF {
		tok := p.consumeToken()
		if tok.Type == TokenWord {
			name = append(name, tok.Value)
		}
	}

	p.skipNewline()
	return strings.Join(name, " ")
}

// parseSetGroup parses a single set group like "Work: 2x 08-10 20kg (6 8)"
func (p *Parser) parseSetGroup() (*SetGroup, error) {
	sg := &SetGroup{}

	// Parse set type (Warm up, Feeder, Work, Top set)
	setTypeToken := p.consumeToken()
	if setTypeToken.Type != TokenWord {
		return nil, ParseError{Message: "expected set type", Line: setTypeToken.Line, Col: setTypeToken.Col}
	}

	setType := parseSetType(setTypeToken.Value)
	if setType == "" {
		return nil, ParseError{Message: "invalid set type: " + setTypeToken.Value, Line: setTypeToken.Line, Col: setTypeToken.Col}
	}
	sg.SetType = SetType(setType)

	// Expect colon
	if p.peekToken().Type != TokenColon {
		return nil, ParseError{Message: "expected ':'", Line: p.peekToken().Line, Col: p.peekToken().Col}
	}
	p.consumeToken()

	// Parse number of sets (e.g., "2x")
	numToken := p.consumeToken()
	if numToken.Type != TokenNumber {
		return nil, ParseError{Message: "expected number of sets", Line: numToken.Line, Col: numToken.Col}
	}
	num, _ := strconv.Atoi(numToken.Value)
	sg.PlannedSets = num

	// Expect multiplier 'x'
	if p.peekToken().Type != TokenMultiplier {
		return nil, ParseError{Message: "expected 'x' after number", Line: p.peekToken().Line, Col: p.peekToken().Col}
	}
	p.consumeToken()

	// Parse rep range (e.g., "08-10")
	repMinToken := p.consumeToken()
	if repMinToken.Type != TokenNumber {
		return nil, ParseError{Message: "expected minimum reps", Line: repMinToken.Line, Col: repMinToken.Col}
	}
	repMin, _ := strconv.Atoi(repMinToken.Value)
	sg.TargetRepMin = repMin

	// Expect minus
	if p.peekToken().Type != TokenMinus {
		return nil, ParseError{Message: "expected '-' in rep range", Line: p.peekToken().Line, Col: p.peekToken().Col}
	}
	p.consumeToken()

	// Parse max reps
	repMaxToken := p.consumeToken()
	if repMaxToken.Type != TokenNumber {
		return nil, ParseError{Message: "expected maximum reps", Line: repMaxToken.Line, Col: repMaxToken.Col}
	}
	repMax, _ := strconv.Atoi(repMaxToken.Value)
	sg.TargetRepMax = repMax

	// Parse weight with unit (e.g., "20kg")
	weightToken := p.consumeToken()
	if weightToken.Type != TokenNumber {
		return nil, ParseError{Message: "expected weight", Line: weightToken.Line, Col: weightToken.Col}
	}
	weight, _ := strconv.ParseFloat(weightToken.Value, 64)
	sg.Weight = weight

	unitToken := p.peekToken()
	if unitToken.Type == TokenUnit {
		sg.Unit = unitToken.Value
		p.consumeToken()
	} else {
		sg.Unit = "kg" // default unit
	}

	// Parse executed reps in parentheses (e.g., "(6 8)")
	if p.peekToken().Type != TokenParenOpen {
		return sg, nil // executed reps are optional
	}
	p.consumeToken()

	reps := []int{}
	for p.peekToken().Type != TokenParenClose && p.peekToken().Type != TokenEOF {
		repToken := p.consumeToken()
		if repToken.Type == TokenNumber {
			rep, _ := strconv.Atoi(repToken.Value)
			reps = append(reps, rep)
		}
	}

	if p.peekToken().Type != TokenParenClose {
		return nil, ParseError{Message: "expected ')'", Line: p.peekToken().Line, Col: p.peekToken().Col}
	}
	p.consumeToken()

	sg.ExecutedReps = reps
	p.skipNewline()

	return sg, nil
}

func parseSetType(s string) string {
	s = strings.ToUpper(s)
	switch s {
	case "WARM", "WARM-UP", "WARMUP":
		return "WARM_UP"
	case "FEEDER":
		return "FEEDER"
	case "WORK":
		return "WORK"
	case "TOP", "TOP-SET", "TOPSET":
		return "TOP_SET"
	default:
		return ""
	}
}

func (p *Parser) peekToken() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) consumeToken() Token {
	tok := p.peekToken()
	p.pos++
	return tok
}

func (p *Parser) skipNewline() {
	for p.pos < len(p.tokens) && p.peekToken().Type == TokenNewline {
		p.pos++
	}
}

func (p *Parser) skipToNextLine() {
	for p.pos < len(p.tokens) && p.peekToken().Type != TokenNewline {
		p.pos++
	}
	p.skipNewline()
}
