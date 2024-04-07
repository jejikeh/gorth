package main

import (
	"fmt"
	"os"
	"strings"
)

func genQBE(path string, program []*Instruction) {
	stack := make([]string, 0)

	bufAsm := strings.Builder{}

	bufAsm.WriteString("export function w $main() {\n")
	bufAsm.WriteString("@start\n")

	for i, inst := range program {
		switch inst.Type {
		case Push:
			stackValue := "%" + fmt.Sprintf(".st%d", len(stack))
			stack = append(stack, stackValue)

			bufAsm.WriteString(fmt.Sprintf("	%s =w copy %d\n", stackValue, inst.NumberValue))

		case Plus:
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stackValue := "%" + fmt.Sprintf(".st%d", len(stack))
			stack = append(stack, stackValue)

			bufAsm.WriteString(fmt.Sprintf("	%s =w add %s, %s\n", stackValue, b, a))

		case Minus:
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stackValue := "%" + fmt.Sprintf(".st%d", len(stack))
			stack = append(stack, stackValue)

			bufAsm.WriteString(fmt.Sprintf("	%s =w sub %s, %s\n", stackValue, b, a))

		case Multiply:
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stackValue := "%" + fmt.Sprintf(".st%d", len(stack))
			stack = append(stack, stackValue)

			bufAsm.WriteString(fmt.Sprintf("	%s =w mul %s, %s\n", stackValue, b, a))

		case Divide:
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stackValue := "%" + fmt.Sprintf(".st%d", len(stack))
			stack = append(stack, stackValue)

			bufAsm.WriteString(fmt.Sprintf("	%s =w div %s, %s\n", stackValue, b, a))

		case Equal:
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]

			stack = stack[:len(stack)-2]

			stackValue := "%" + fmt.Sprintf(".st%d", len(stack))
			stack = append(stack, stackValue)

			bufAsm.WriteString(fmt.Sprintf("	%s =w ceqw %s, %s\n", stackValue, b, a))

		case Assert:
			stack = stack[:len(stack)-1]
			fmt.Println("Assert is not implemented in codegen")

		case Dump:
			a := stack[len(stack)-1]

			bufAsm.WriteString(fmt.Sprintf("	call $printf(l $dump, ...,w %d , w %s)\n", i, a))

		default:
			panic(fmt.Sprintf("unknow instruction: %v", inst))
		}
	}

	bufAsm.WriteString("	ret 0\n")
	bufAsm.WriteString("}\n")
	bufAsm.WriteString("data $dump = {b \"%d: Dump = [%d]\\n\", b 0 }\n")
	bufAsm.WriteString("data $dump = {b \"%d: Dump = [%d]\\n\", b 0 }\n")

	asm := bufAsm.String()

	err := os.WriteFile(path, []byte(asm), 0644)

	if err != nil {
		panic(err)
	}

	fmt.Println(asm)
}
