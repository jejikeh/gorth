package main

import (
	"fmt"
	"os"
	"strings"
)

func genQBE(path string, program []Instruction) {
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

		case Equal:
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stackValue := "%" + fmt.Sprintf(".st%d", len(stack))
			stack = append(stack, stackValue)

			asmBuf.WriteString(fmt.Sprintf("	%s =w ceqw %s, %s\n", stackValue, b, a))

		case Assert:
			// @Note: Pop expected value from stack
			stack = stack[:len(stack)-1]
			fmt.Println("Assert is not implemented in codegen")

		case Dump:
			a := stack[len(stack)-1]

			asmBuf.WriteString(fmt.Sprintf("	call $printf(l $dump, ...,w %d , w %s)\n", i, a))

		default:
			panic(fmt.Sprintf("unknow instruction: %v", inst))
		}
	}

	asmBuf.WriteString("	ret 0\n")
	asmBuf.WriteString("}\n")
	asmBuf.WriteString("data $dump = {b \"%d: Dump = [%d]\\n\", b 0 }\n")

	asm := asmBuf.String()

	err := os.WriteFile(path, []byte(asm), 0644)

	if err != nil {
		panic(err)
	}

	fmt.Println(asm)
}
