package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
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

		buildProgram(*outputPath, program)

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
	Dump
)

type Instruction struct {
	Type        InstructionType
	NumberValue int
	Name        string
}

func push(x int) Instruction {
	return Instruction{
		Type:        Push,
		NumberValue: x,
	}
}

func plus() Instruction {
	return Instruction{
		Type: Plus,
	}
}

func sub() Instruction {
	return Instruction{
		Type: Minus,
	}
}

func mul() Instruction {
	return Instruction{
		Type: Multiply,
	}
}

func div() Instruction {
	return Instruction{
		Type: Divide,
	}
}

func dump() Instruction {
	return Instruction{
		Type: Dump,
	}
}

func runProgram(program []Instruction) {
	stack := make([]int, 0)

	for i, inst := range program {
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

		case Dump:
			a := stack[len(stack)-1]

			fmt.Printf("%d: Dump = %d\n", i, a)

		default:
			panic(fmt.Sprintf("unknow instruction: %v", inst))
		}
	}
}

func buildProgram(path string, program []Instruction) {
	stack := make([]string, 0)

	asmBuf := strings.Builder{}

	asmBuf.WriteString("export function w $main() {\n")
	asmBuf.WriteString("@start\n")

	for i, inst := range program {
		switch inst.Type {
		case Push:
			stackValue := "%" + fmt.Sprintf(".st%d", len(stack))
			stack = append(stack, stackValue)

			asmBuf.WriteString(fmt.Sprintf("	%s =w copy %d\n", stackValue, inst.NumberValue))

		case Plus:
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stackValue := "%" + fmt.Sprintf(".st%d", len(stack))
			stack = append(stack, stackValue)

			asmBuf.WriteString(fmt.Sprintf("	%s =w add %s, %s\n", stackValue, b, a))

		case Minus:
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stackValue := "%" + fmt.Sprintf(".st%d", len(stack))
			stack = append(stack, stackValue)

			asmBuf.WriteString(fmt.Sprintf("	%s =w sub %s, %s\n", stackValue, b, a))

		case Multiply:
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stackValue := "%" + fmt.Sprintf(".st%d", len(stack))
			stack = append(stack, stackValue)

			asmBuf.WriteString(fmt.Sprintf("	%s =w mul %s, %s\n", stackValue, b, a))

		case Divide:
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stackValue := "%" + fmt.Sprintf(".st%d", len(stack))
			stack = append(stack, stackValue)

			asmBuf.WriteString(fmt.Sprintf("	%s =w div %s, %s\n", stackValue, b, a))

		case Dump:
			a := stack[len(stack)-1]

			asmBuf.WriteString(fmt.Sprintf("	call $printf(l $dump, ...,w %d , w %s)\n", i, a))

		default:
			panic(fmt.Sprintf("unknow instruction: %v", inst))
		}
	}

	asmBuf.WriteString("	ret 0\n")
	asmBuf.WriteString("}\n")
	asmBuf.WriteString("data $dump = {b \"%d: Dump = %d\\n\", b 0 }\n")

	asm := asmBuf.String()

	err := os.WriteFile(path, []byte(asm), 0644)

	if err != nil {
		panic(err)
	}

	fmt.Println(asm)
}
