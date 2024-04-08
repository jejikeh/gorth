package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	buildCommand := NewBuildCommand()
	runCommand := NewRunCommand()

	flag.Parse()

	switch os.Args[1] {
	case "build":
		buildCommand.Run()

	case "run":
		runCommand.Run()

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

	PopValueFromCondition

	Dump
	Assert
	StackPrintln
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
		Type: PopValueFromCondition,
	}
}

func assert() *Instruction {
	return &Instruction{
		Type: Assert,
	}
}

func stackprintln() *Instruction {
	return &Instruction{
		Type: StackPrintln,
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
			} else {
				stack = stack[:len(stack)-1]
			}

		case Else:
			a := stack[len(stack)-1]

			if a != 0 {
				i = inst.NumberValue
			} else {
				stack = stack[:len(stack)-1]
			}

		case PopValueFromCondition:
			stack = stack[:len(stack)-1]

		case StackPrintln:
			fmt.Println("\n__stack_println__")
			fmt.Printf("%d: Stack = [%d]\n", i, len(stack))
			for i, v := range stack {
				fmt.Printf("%d: [%d]\n", i, v)
			}

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

type BuildCommand struct {
	buildCmd *flag.FlagSet
	path     *string
	output   *string
}

func NewBuildCommand() *BuildCommand {
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)

	return &BuildCommand{
		buildCmd: buildCmd,
		path:     buildCmd.String("i", "", "path to project"),
		output:   buildCmd.String("o", "", "output path"),
	}
}

func (c *BuildCommand) Run() {
	c.buildCmd.Parse(os.Args[2:])

	lexer := NewLexer(*c.path)

	program, err := lexer.loadProgramFromFile()

	if err != nil {
		panic(err)
	}

	genQBE(*c.output, program)
}

type RunCommand struct {
	runCmd *flag.FlagSet
	path   *string
}

func NewRunCommand() *RunCommand {
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)

	return &RunCommand{
		runCmd: runCmd,
		path:   runCmd.String("i", "", "path to project"),
	}
}

func (c *RunCommand) Run() {
	c.runCmd.Parse(os.Args[2:])

	lexer := NewLexer(*c.path)

	program, err := lexer.loadProgramFromFile()

	if err != nil {
		panic(err)
	}

	runProgram(program)
}
