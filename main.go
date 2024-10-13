package main

import (
	"fmt"
	"strconv"
)

func Reverse(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func Interpret(tests []string) {
	stack := Stack{}
	operations := map[string]Node{"+": Add{}, "-": Subtract{}, "*": Multiply{}, "/": Divide{}, "^": Pow{}, "cos": Cos{}, "sin": Sin{}, "tan": Tan{}, "ln": Log{}, "exp": Exp{}}
	for _, exp := range tests {
		for i := len(exp) - 1; i >= 0; i-- {
			if val, ok := operations[string(exp[i])]; ok {
				stack.Parse(val)
			} else if exp[i] > 'A' && exp[i] < 'z' {
				str := ""
				// Collects string
				for exp[i] > 'A' && exp[i] < 'z' {
					if i == 0 {
						break
					}
					str += string(exp[i])
					i--
				}
				// Checks if str is a keyword
				if val, ok := operations[Reverse(str)]; ok {
					stack.Parse(val)
				} else {
					stack.Push(Identifier{str})
				}
			} else if exp[i] >= '0' && exp[i] <= '9' {
				num := ""
				for exp[i] >= '0' && exp[i] <= '9' {
					num += string(exp[i])
					i--
					if i < 0 {
						break
					}
				}
				intNum, err := strconv.Atoi(Reverse(num))
				if err == nil {
					stack.Push(Literal{float64(intNum)})
				} else {
					panic("Unable to convert string to float.")
				}
			}
		}
		diffed := stack.data[0].Diff()
		simplified := diffed.Simplify()
		fmt.Println(simplified.PrefixString())
		stack.Clear()
	}
}

func main() {
	tests := []string{"(ln x)"}
	Interpret(tests)
}
