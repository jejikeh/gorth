package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	program := []Instruction{
		push("a", 2),
		push("b", 3),
		plus("c"),
		push("r", -1),
		multiply("t"),
		dump("t"),
	}

	// buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	// buildPath := buildCmd.String("path", "", "path to project")
	// outputPath := buildCmd.String("output", "", "output path")

	// runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	// runPath := runCmd.String("path", "", "path to project")

	flag.Parse()

	switch os.Args[1] {
	case "build":
		buildProgram(program)

	case "run":
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

func push(name string, x int) Instruction {
	return Instruction{
		Type:        Push,
		NumberValue: x,
		Name:        name,
	}
}

func plus(name string) Instruction {
	return Instruction{
		Type: Plus,
		Name: name,
	}
}

func minus(name string) Instruction {
	return Instruction{
		Type: Minus,
		Name: name,
	}
}

func multiply(name string) Instruction {
	return Instruction{
		Type: Multiply,
		Name: name,
	}
}

func divide(name string) Instruction {
	return Instruction{
		Type: Divide,
		Name: name,
	}
}

func dump(name string) Instruction {
	return Instruction{
		Type: Dump,
		Name: name,
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

			stack = stack[:len(stack)-1]

			fmt.Printf("%d: Dump = %d\n", i, a)

		default:
			panic(fmt.Sprintf("unknow instruction: %v", inst))
		}
	}
}

func buildProgram(program []Instruction) {
	stack := make([]Instruction, 0)

	asmBuf := strings.Builder{}

	asmBuf.WriteString("export function w $main() {\n")
	asmBuf.WriteString("@start\n")

	for i, inst := range program {
		switch inst.Type {
		case Push:
			stack = append(stack, inst)

			asmBuf.WriteString(fmt.Sprintf("	%s =w copy %d\n", "%"+inst.Name, inst.NumberValue))

		case Plus:
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stack = append(stack, inst)

			asmBuf.WriteString(fmt.Sprintf("	%s =w add %s, %s\n", "%"+inst.Name, "%"+a.Name, "%"+b.Name))

		case Minus:
			// a := stack[len(stack)-1]
			// b := stack[len(stack)-2]

			// stack = stack[:len(stack)-2]

			// stack = append(stack, b-a)

			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stack = append(stack, inst)

			asmBuf.WriteString(fmt.Sprintf("	%s =w sub %s, %s\n", "%"+inst.Name, "%"+a.Name, "%"+b.Name))

		case Multiply:
			// a := stack[len(stack)-1]
			// b := stack[len(stack)-2]

			// stack = stack[:len(stack)-2]

			// stack = append(stack, a*b)

			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stack = append(stack, inst)

			asmBuf.WriteString(fmt.Sprintf("	%s =w mul %s, %s\n", "%"+inst.Name, "%"+a.Name, "%"+b.Name))

		case Divide:
			// 	a := stack[len(stack)-1]
			// 	b := stack[len(stack)-2]

			// 	stack = stack[:len(stack)-2]

			// 	stack = append(stack, b/a)

			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stack = append(stack, inst)

			asmBuf.WriteString(fmt.Sprintf("	%s =w div %s, %s\n", "%"+inst.Name, "%"+a.Name, "%"+b.Name))

		case Dump:
			// a := stack[len(stack)-1]

			// stack = stack[:len(stack)-1]

			// fmt.Printf("%d: Dump = %d\n", i, a)

			asmBuf.WriteString(fmt.Sprintf("	call $printf(l $dump, ...,w %d , w %s)\n", i, "%"+inst.Name))

		default:
			panic(fmt.Sprintf("unknow instruction: %v", inst))
		}
	}

	asmBuf.WriteString("	ret 0\n")
	asmBuf.WriteString("}\n")
	asmBuf.WriteString("data $dump = {b \"%d: Dump = %d\\n\", b 0 }\n")

	asm := asmBuf.String()

	err := os.WriteFile("file.ssa", []byte(asm), 0644)

	if err != nil {
		panic(err)
	}
}
