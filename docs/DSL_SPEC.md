# DSL Specification - IRONLOG Workout Language

## Overview

The IRONLOG DSL is a lightweight, deterministic Domain-Specific Language for recording strength training workouts with minimal typing while capturing complete information.

## Grammar

### EBNF Notation

```ebnf
Workout      = ExerciseName ( SetLine )*
ExerciseName = ( Word )+ Newline
SetLine      = SetType Colon SetGroup Newline
SetGroup     = SetMultiplier RepRange Weight [ ExecutedReps ]
SetType      = SetTypeKeyword
SetTypeKeyword = "Warm up" | "Feeder" | "Work" | "Top set"
                | "WARM UP" | "FEEDER" | "WORK" | "TOP SET"
                | "Warm-up" | "warm up" (case-insensitive)

SetMultiplier = Number "x"
RepRange     = Number "-" Number
Weight       = Number Unit
Unit         = "kg" | "lbs" | "lb" | "g"
ExecutedReps = "(" Number ( Number )* ")"

Number       = Digit+ [ "." Digit+ ]
Word         = Letter+ [ "-" Letter+ ]*
Digit        = "0" | "1" | ... | "9"
Letter       = "a" | ... | "z" | "A" | ... | "Z"
Newline      = "\n"
```

## Lexical Tokens

| Token | Pattern | Example |
|-------|---------|---------|
| NUMBER | `[0-9]+(\.[0-9]+)?` | `20`, `10.5` |
| WORD | `[a-zA-Z][a-zA-Z0-9_-]*` | `Supino`, `Bench-Press` |
| MULTIPLIER | `x` (after number) | `3x` |
| COLON | `:` | `:` |
| MINUS | `-` | `-` |
| UNIT | `kg`\|`lbs`\|`lb`\|`g` | `kg` |
| PAREN_OPEN | `(` | `(` |
| PAREN_CLOSE | `)` | `)` |
| COMMA | `,` | `,` |
| WHITESPACE | Space, Tab | ` ` |
| NEWLINE | `\n` | Newline |

## Parsing Process

### 1. Lexical Analysis (Lexer)

Input: Raw text
```
Supino Reto
Warm up: 1x 1-20 10kg (10)
```

Output: Token stream
```
WORD("Supino")
WORD("Reto")
NEWLINE
WORD("Warm")
WORD("up")
COLON
NUMBER(1)
MULTIPLIER(x)
NUMBER(1)
MINUS
NUMBER(20)
NUMBER(10)
UNIT(kg)
PAREN_OPEN
NUMBER(10)
PAREN_CLOSE
NEWLINE
EOF
```

### 2. Syntax Analysis (Parser)

Tokens → AST (Abstract Syntax Tree)

```
ParsedWorkout {
  exercise_name: "Supino Reto",
  set_groups: [
    SetGroup {
      set_type: "WARM_UP",
      planned_sets: 1,
      target_rep_min: 1,
      target_rep_max: 20,
      weight: 10.0,
      unit: "kg",
      executed_reps: [10]
    }
  ]
}
```

## Examples

### Example 1: Complete Upper Body Session

**Input:**
```
SUPINO RETO BARRA

Warm up: 1x 1-20 10kg (10)
Feeder: 1x 08-10 15kg (10)
Work: 2x 08-10 20kg (6 8)
Top set: 1x 6-8 25kg (3)
```

**AST Output:**
```json
{
  "exercise_name": "SUPINO RETO BARRA",
  "set_groups": [
    {
      "set_type": "WARM_UP",
      "planned_sets": 1,
      "target_rep_min": 1,
      "target_rep_max": 20,
      "weight": 10.0,
      "unit": "kg",
      "executed_reps": [10]
    },
    {
      "set_type": "FEEDER",
      "planned_sets": 1,
      "target_rep_min": 8,
      "target_rep_max": 10,
      "weight": 15.0,
      "unit": "kg",
      "executed_reps": [10]
    },
    {
      "set_type": "WORK",
      "planned_sets": 2,
      "target_rep_min": 8,
      "target_rep_max": 10,
      "weight": 20.0,
      "unit": "kg",
      "executed_reps": [6, 8]
    },
    {
      "set_type": "TOP_SET",
      "planned_sets": 1,
      "target_rep_min": 6,
      "target_rep_max": 8,
      "weight": 25.0,
      "unit": "kg",
      "executed_reps": [3]
    }
  ]
}
```

### Example 2: Leg Day

**Input:**
```
AGACHAMENTO LIVRE

Warm up: 2x 1-10 40kg (10 5)
Work: 4x 5-7 80kg (7 6 6 5)
Accessory: 3x 10-12 60kg (11 11 10)
```

**Key Points:**
- Multiple sets per line handled by `planned_sets`
- Rep ranges allow flexibility (e.g., 5-7)
- All executed reps tracked in array

### Example 3: Minimal Format

**Input:**
```
Deadlift

Work: 3x 3-5 100kg (5 4 3)
```

**Valid Pattern:**
- Exercise name on first line(s)
- Sets on separate lines
- Minimal whitespace

### Example 4: Variations

**Input:**
```
Leg Press

Warm up: 1x 15-20 80kg (15)
Work: 1x 6-8 180kg (6)
Drop: 1x 8-10 120kg (10 9)
```

**Note:** Set type is flexible (any label accepted), but standard types are:
- `Warm up` (or `Warm-up`)
- `Feeder`
- `Work`
- `Top set` (or `Top-set`)

## Error Handling

### Parse Errors

| Error | Message | Example |
|-------|---------|---------|
| Missing set type | "Expected set type" | `3x 6-8 20kg (8)` |
| Missing colon | "Expected ':'" | `Work 3x 6-8 20kg` |
| Invalid rep range | "Expected '-' in rep range" | `Work: 3x 6_8 20kg` |
| Missing weight | "Expected weight" | `Work: 3x 6-8 (8)` |
| Invalid parentheses | "Expected ')'" | `Work: 3x 6-8 20kg (8` |

### Validation Errors

```go
if planned_sets <= 0 {
    error: "Planned sets must be greater than 0"
}
if target_rep_min > target_rep_max {
    error: "Min reps cannot exceed max reps"
}
if weight <= 0 {
    error: "Weight must be positive"
}
```

## Set Type Classifications

### Warm Up
- **Purpose:** Preparation sets
- **Characteristics:** High reps, light weight
- **Example:** `Warm up: 1x 1-20 10kg (10)`

### Feeder
- **Purpose:** Building up to working sets
- **Characteristics:** Moderate reps, moderate weight
- **Example:** `Feeder: 1x 08-10 15kg (10)`

### Work
- **Purpose:** Primary training sets
- **Characteristics:** Target rep range, working weight
- **Example:** `Work: 3x 6-8 20kg (8 7 6)`

### Top Set
- **Purpose:** Maximum effort set
- **Characteristics:** Lower reps, highest weight
- **Example:** `Top set: 1x 3-5 30kg (3)`

## Volume Calculation

Volume = Weight × Total Reps

```
Example:
Executed reps: [6, 8]
Weight: 20kg
Volume = 20 × (6 + 8) = 20 × 14 = 280 kg
```

## Tokenization Algorithm

```go
func Tokenize(input string) []Token {
    tokens := []Token{}
    
    for pos < len(input) {
        ch := input[pos]
        
        switch ch {
        case ':':
            tokens.append(TokenColon)
        case '-':
            if !isPartOfNumber(input[pos-1]) {
                tokens.append(TokenMinus)
            }
        case '(':
            tokens.append(TokenParenOpen)
        case ')':
            tokens.append(TokenParenClose)
        case ' ', '\t':
            skipWhitespace()
        case '\n':
            tokens.append(TokenNewline)
        case isDigit(ch):
            tokens.append(readNumber())
        case isLetter(ch):
            word := readWord()
            if isUnit(word) {
                tokens.append(TokenUnit(word))
            } else if word == "x" && isAfterNumber() {
                tokens.append(TokenMultiplier)
            } else {
                tokens.append(TokenWord(word))
            }
        }
    }
    
    return tokens
}
```

## Parsing Algorithm

```go
func Parse(tokens []Token) (*ParsedWorkout, error) {
    workout := ParsedWorkout{}
    
    // Line 1: Exercise name
    workout.ExerciseName = parseExerciseName()
    
    // Following lines: Set groups
    for token != EOF {
        setGroup, err := parseSetGroup()
        if err != nil {
            return nil, err
        }
        workout.SetGroups.append(setGroup)
    }
    
    return &workout, nil
}

func parseSetGroup() (*SetGroup, error) {
    // Expect: SetType Colon Multiplier RepRange Weight [ ExecutedReps ]
    
    sg := SetGroup{}
    
    setType := expect(TokenWord)
    sg.SetType = normalizeSetType(setType)
    
    expect(TokenColon)
    
    num := expect(TokenNumber)
    expect(TokenMultiplier)
    sg.PlannedSets = parseInt(num)
    
    minReps := expect(TokenNumber)
    expect(TokenMinus)
    maxReps := expect(TokenNumber)
    sg.TargetRepMin = parseInt(minReps)
    sg.TargetRepMax = parseInt(maxReps)
    
    weight := expect(TokenNumber)
    unit := expect(TokenUnit)
    sg.Weight = parseFloat(weight)
    sg.Unit = unit
    
    if peek() == TokenParenOpen {
        sg.ExecutedReps = parseExecutedReps()
    }
    
    return &sg, nil
}
```

## Edge Cases

### Multi-Word Exercise Names
```
SUPINO RETO BARRA COM PAUSA
```
Parser collects all words before first line with ":"

### Decimal Weights
```
Work: 3x 6-8 22.5kg (8 7 6)
```
Correctly parsed as float: 22.5

### Missing Executed Reps
```
Work: 3x 6-8 20kg
```
Valid - `executed_reps` is optional

### Varying Rep Counts per Set
```
Work: 3x 6-8 20kg (8 7 6)
```
Different reps per set captured separately

### Short Rep Ranges
```
Top set: 1x 3 20kg (3)
```
Interpreted as 3-3 (same min and max)

## Performance Considerations

- **Lexer**: O(n) where n = input length
- **Parser**: O(m) where m = number of tokens
- **Memory**: O(m) for token buffer
- **No backtracking**: Single-pass parsing

## Future Extensions

Potential DSL enhancements:

```
# Comments (future)
# This is a warm-up set
Warm up: 1x 1-20 10kg (10)

# Rest periods (future)
Work: 3x 6-8 20kg (8 7 6) [120s]

# Rate of perceived exertion (future)
Work: 3x 6-8 20kg (8 7 6) RPE:8

# Tempo notation (future)
Work: 3x 6-8 20kg (8 7 6) [3-1-2-1]
```

