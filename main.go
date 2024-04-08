package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildPath := buildCmd.String("i", "", "path to project")
	outputPath := buildCmd.String("o", "", "output path")

	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	runPath := runCmd.String("i", "", "path to project")

	flag.Parse()

	switch os.Args[1] {
	case "build":
		buildCmd.Parse(os.Args[2:])

		lexer := NewLexer(*buildPath)

		program, err := lexer.loadProgramFromFile()

		if err != nil {
			panic(err)
		}

		genQBE(*outputPath, program)

	case "run":
		runCmd.Parse(os.Args[2:])

		lexer := NewLexer(*runPath)

		program, err := lexer.loadProgramFromFile()

		if err != nil {
			panic(err)
		}

		runProgram(program)

	default:
		flag.Usage()
		os.Exit(1)
	}
}

type InstructionType int

const (
	Push InstructionType = iota

	Plus
	Minus
	Multiply
	Divide

	Equal

	LeftBracket
	RightBracket

	If
	Else

	PopValue

	Dump
	Assert
)

type Instruction struct {
	Type        InstructionType
	NumberValue int
	Name        string
}

func push(x int) *Instruction {
	return &Instruction{
		Type:        Push,
		NumberValue: x,
	}
}

func plus() *Instruction {
	return &Instruction{
		Type: Plus,
	}
}

func sub() *Instruction {
	return &Instruction{
		Type: Minus,
	}
}

func mul() *Instruction {
	return &Instruction{
		Type: Multiply,
	}
}

func div() *Instruction {
	return &Instruction{
		Type: Divide,
	}
}

func equal() *Instruction {
	return &Instruction{
		Type: Equal,
	}
}

func dump() *Instruction {
	return &Instruction{
		Type: Dump,
	}
}

func rightBracket() *Instruction {
	return &Instruction{
		Type: RightBracket,
	}
}

func leftBracket() *Instruction {
	return &Instruction{
		Type: LeftBracket,
	}
}

func iff() *Instruction {
	return &Instruction{
		Type: If,
	}
}

func elsee() *Instruction {
	return &Instruction{
		Type: Else,
	}
}

func popval() *Instruction {
	return &Instruction{
		Type: PopValue,
	}
}

func assert() *Instruction {
	return &Instruction{
		Type: Assert,
	}
}

func runProgram(program []*Instruction) {
	stack := make([]int, 0)

	for i := 0; i < len(program); i++ {
		inst := program[i]

		switch inst.Type {
		case Push:
			stack = append(stack, inst.NumberValue)

		case Plus:
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stack = append(stack, a+b)

		case Minus:
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stack = append(stack, b-a)

		case Multiply:
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stack = append(stack, a*b)

		case Divide:
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stack = append(stack, b/a)

		case Equal:
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stack = append(stack, boolToInt(a == b))

		case Dump:
			a := stack[len(stack)-1]

			fmt.Printf("%d: Dump = [%d]\n", i, a)

		case Assert:
			expected := stack[len(stack)-1]
			got := stack[len(stack)-2]

			stack = stack[:len(stack)-1]

			if expected != got {
				panic(fmt.Sprintf("%d: Assert = [expected '%d' but got '%d']", i, expected, got))
			}

		case LeftBracket:

		case RightBracket:

		case If:
			a := stack[len(stack)-1]

			if a != 1 {
				i = inst.NumberValue
			}

		case Else:
			a := stack[len(stack)-1]

			if a != 0 {
				i = inst.NumberValue
			}

		case PopValue:
			stack = stack[:len(stack)-1]

		default:
			panic(fmt.Sprintf("unknow instruction: %v", inst))
		}
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}
