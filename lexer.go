package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Lexer struct {
	filePath                    string
	source                      []string
	currentLine, currentCollumn int
	fileContainError            bool
}

func NewLexer(path string) *Lexer {
	fileContent, err := os.ReadFile(path)

	if err != nil {
		// @Cleanup:
		panic(err)
	}

	// @Cleanup: I think it is better to use maybe custon func in FiledsFunc
	// c := strings.ReplaceAll(, ";", " ; ")
	// c = strings.ReplaceAll(string(fileContent), "}", " } ")
	// c = strings.ReplaceAll(string(fileContent), "{", " { ")

	s := strings.Split(string(fileContent), "\n")

	return &Lexer{
		filePath: path,
		source:   s,
	}
}

func (l *Lexer) loadProgramFromFile() ([]*Instruction, error) {

	instructions := make([]*Instruction, 0)

	for lineIndex, line := range l.source {
		l.currentLine = lineIndex

		// @Incomplete: The source location of words are lost because of Fields.
		// We are reporting not the collumn, but rather world index

		// @Note: If you decide what you need to look ahead some tokens
		// Make sure what you realy need it, because it is just a lexer.
		// You can always look ahead in Parser while iterating AST.
		words := strings.Fields(string(line))

	lineParsing:
		for wordIndex, w := range words {
			l.currentCollumn = wordIndex

			switch {
			case unicode.IsDigit(rune(w[0])):
				if num, err := strconv.Atoi(w); err == nil {
					instructions = append(instructions, push(num))
				} else {
					// @Incomplete: Now we have the lineIndex and wordIndex.
					// We can do a better error reporting with the line and something like cursor?

					// @Cleanup: Move error creating to separate function.
					// Something like ExpectButGotError(Number, w)
					l.reportExpectedButGotError("number", w)
				}
			case w == "+":
				instructions = append(instructions, plus())

			case w == "-":
				instructions = append(instructions, sub())

			case w == "*":
				instructions = append(instructions, mul())

			case w == "==":
				instructions = append(instructions, equal())

			case w == "{":
				instructions = append(instructions, leftBracket())

			case w == "}":
				instructions = append(instructions, rightBracket())

			case w == "/":
				instructions = append(instructions, div())

			case w == "assert":
				instructions = append(instructions, assert())

			case w == "println":
				instructions = append(instructions, dump())

			case w == "if":
				instructions = append(instructions, iff())

			case w == "else":
				instructions = append(instructions, elsee())

			case w == "};", w == "{};":
				instructions = append(instructions, popval())

			case w == "//":
				// @Incomplete: We do not handle code after comment on the same line. Sad?
				break lineParsing

			case w == "__stack_println__":
				instructions = append(instructions, stackprintln())

			default:
				// @Incomplete: If we got error on the line, should we skip this line to not
				// print too many error?
				l.reportUnexpectedTokenError(w)
				break lineParsing
			}
		}
	}

	instructions = l.crossReference(instructions)

	if l.fileContainError {
		// This way it is easier to see the errors, but i might just delete that when
		// we have better way to panic?
		fmt.Println()
		return instructions, fmt.Errorf("too many errors")
	}

	return instructions, nil
}

// @Cleanup: This doesnt seams to fit right in lexer part.
// Maybe look some jblow videos to figouraut some things?
func (l *Lexer) crossReference(instructions []*Instruction) []*Instruction {
	stack := make([]int, 0)

	for i := 0; i < len(instructions); i++ {
		instruction := instructions[i]

		switch instruction.Type {
		case If, Else:
			if i == len(instructions)-1 {
				l.reportExpectedButGotError("{", "EndOfFile")

				break
			}

			if instructions[i+1].Type != LeftBracket {
				// @Incomplete: Need to get the lexeme in here for error reporting
				l.reportExpectedButGotError("{", "NotImplemented")

				break
			}

			stack = append(stack, i)

		case RightBracket, PopValueFromCondition:
			if len(stack) == 0 {
				// @Cleanup: Better error
				// @Incomplete: currentLine and currentCursor are incorrect here
				l.reportExpectedButGotError("}", "No { condition on the stack")

				break
			}

			ifInstruction := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			instructions[ifInstruction].NumberValue = i - 1

			// case PopValue:
			// 	stack = stack[:len(stack)-1]
		}

	}

	// @Incomplete: is better to check if some Instructions was not initialized?

	return instructions
}

func (l *Lexer) reportError(errorMessage string) {
	l.fileContainError = true

	fmt.Printf("%s:%d:%d: %s\n", l.filePath, l.currentLine+1, l.currentCollumn+1, errorMessage)
}

func (l *Lexer) reportExpectedButGotError(expected, got string) {
	l.reportError(fmt.Sprintf("expected '%s', but got '%s'", expected, got))
}

func (l *Lexer) reportUnexpectedTokenError(unexpectedToken string) {
	l.reportError(fmt.Sprintf("unexpected token '%s'", unexpectedToken))
}
