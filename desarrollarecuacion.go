package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Operation struct {
	s string
	o int
	t int
}

type Section struct {
	a []Part
	s int
	e int
}

type Part struct {
	v string
	t int
}

var operations []Operation

func main() {
	var text string
	operations = []Operation{
		Operation{"1", 0, 1},
		Operation{"(", 0, 2},
		Operation{")", 0, 3},
		Operation{"*", 4, 4},
		Operation{"/", 3, 5},
		Operation{"^", 2, 6},
		Operation{"%", 1, 7},
		Operation{"+", 5, 8},
		Operation{"-", 6, 9},
		Operation{">", 0, 10},
		Operation{"<", 0, 11},
		Operation{">=", 0, 12},
		Operation{"<=", 0, 13},
		Operation{"==", 0, 14},
		Operation{"!=", 0, 15},
	}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Escribe la funcion")
	scanner.Scan()
	text = scanner.Text()
	fmt.Println(text)
	str, t := findExpresion(text)
	//fmt.Println(str, types(t))
	if t != "" {
		result := op(convToFloat(txtToMath(str[0])), convToFloat(txtToMath(str[1])), types(t))
		if result == 1 {
			fmt.Println("         ---->", text, " = Verdadero")
		} else {
			fmt.Println("         ---->", text, " = Falso")
		}

	} else {
		result := txtToMath(text)
		fmt.Println("         ---->", text, "=", result)
	}
}

func findExpresion(input string) ([]string, string) {
	exp := []string{">", "<", ">=", "<=", "!=", "=="}
	f := -1
	t := -1
	for index := 0; index < len(exp); index++ {
		i := strings.Index(input, exp[index])
		if i != -1 {
			f = i
			t = index
		}
	}
	if f != -1 && t != -1 {
		return []string{input[:f], input[f+len(exp[t]):]}, exp[t]
	}
	return []string{}, ""
}

func txtToMath(input string) string {
	var stack []string
	input = strings.ReplaceAll(input, " ", "")
	stack = strings.Split(input, "")
	equation := group(order(stack))
	result := resolve(equation)[0].v
	fmt.Println("	|   --->", joinPart(equation, " "), "=", result)
	fmt.Println("	|   Resultado Final")
	return result
}

func group(stack []Part) []Part {
	var sec []Section
	result := stack
	for {
		res := parent(result, operations[1].t, operations[2].t)
		result = res.a
		if res.s == 0 && res.e == 0 {
			break
		}
		sec = append(sec, res)
	}
	return resolveEquation(invertSection(sec), stack)
}

func resolveEquation(section []Section, stack []Part) []Part {
	for index := 0; index < len(section); index++ {
		part := resolve(section[index].a)
		fmt.Println("	|   Resultado")
		fmt.Println("	|   --->", joinPart(section[index].a, " "), "=", part[0].v)
		if index+1 < len(section) {
			section[index+1].a = cutPart(section[index+1].a, part[0].v, section[index].s-2, section[index].e)
		} else {
			stack = cutPart(stack, joinPart(resolve(section[index].a), ""), section[index].s-2, section[index].e)
		}
	}
	return stack
}

func resolve(part []Part) []Part {
	fmt.Println("	|   Seccion", joinPart(part, " "))
	for _, n := range []int{7, 6, 5, 4, 8, 9} {
		for {
			pos := findIndex(part, n)
			if pos == -1 {
				break
			}
			a := convToFloat(part[pos-1 : pos][0].v)
			b := convToFloat(part[pos+1 : pos+2][0].v)
			res := op(a, b, n)
			part = cutPart(part, convToString(res), pos-2, pos+1)
			fmt.Println("	|  |   ", a, typesReverse(n), b, "=", res)
			fmt.Println("	|  |      ", joinPart(part, " "))
		}
	}
	return part
}

func joinPart(arr []Part, concat string) string {
	res := ""
	for _, n := range arr {
		res = res + n.v + concat
	}
	return res
}

func cutPart(arr []Part, add string, s int, e int) []Part {
	res := []Part{}
	c := 0
	for i, n := range arr {
		if i > s && i <= e {
			if c == 0 {
				res = append(res, Part{add, 1})
				c++
			}
		} else {
			res = append(res, Part{n.v, n.t})
		}
	}
	return res
}

func op(a float64, b float64, types int) float64 {
	switch types {
	case 4:
		return a * b
	case 5:
		return a / b
	case 6:
		return mult(a, b)
	case 7:
		return mod(a, b)
	case 8:
		return a + b
	case 9:
		return a - b
	case 10:
		if a > b {
			return 1
		}
		return 0
	case 11:
		if a < b {
			return 1
		}
		return 0
	case 12:
		if a >= b {
			return 1
		}
		return 0
	case 13:
		if a <= b {
			return 1
		}
		return 0
	case 14:
		if a == b {
			return 1
		}
		return 0
	case 15:
		if a != b {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func mult(a float64, b float64) float64 {
	c := a
	for index := 0; index < int(b)-1; index++ {
		c = c * a
	}
	return c
}

func mod(a float64, b float64) float64 {
	return a - b*float64(int(a/b))
}

func convToString(s float64) string {
	return strconv.FormatFloat(s, 'f', 6, 64)
}

func convToFloat(s string) float64 {
	if res, err := strconv.ParseFloat(s, 64); err == nil {
		return res
	}
	return 0
}

func findIndex(arr []Part, types int) int {
	for i, n := range arr {
		if n.t == types {
			return i
		}
	}
	return -1
}

func invertSection(arr []Section) []Section {
	invert := []Section{}
	for i := range arr {
		n := arr[len(arr)-1-i]
		invert = append(invert, n)
	}
	return invert
}

func parent(arr []Part, a int, b int) Section {
	c := 0
	s := 0
	e := 0
	for i, n := range arr {
		if n.t == a {
			if c == 0 {
				s = i + 1
			}
			c++
		} else if n.t == b {
			c--
			if c == 0 {
				e = i
				break
			}
		}
	}
	return Section{arr[s:e], s, e}
}

func order(stack []string) []Part {
	arr := []Part{}
	a := Part{"", 0}
	b := Part{"", 0}
	num := ""
	c := 0
	parent := 0
	for i := 0; i < len(stack); i++ {
		a = Part{"", 0}
		b = Part{"", 0}
		if i != 0 {
			a = Part{stack[i-1], types(stack[i-1])}
		} else {
			a = Part{"", 0}
		}
		b = Part{stack[i], types(stack[i])}
		if b.t == 0 {
			log.Fatal("Caracter inesperado es la formula")
		}
		if (a.t > 3 && b.t > 3 && b.t <= 7) || (i == 0 && b.t >= 3 && b.t <= 7) || (i+1 >= len(stack) && (b.t >= 4 || b.t == 2)) || (a.t == 1 && b.t == 2) || (a.t == 2 && b.t == 3) || (a.t >= 8 && b.t >= 8) || (b.t == 3 && a.t > 1) || (a.v == "." && b.t > 1) {
			log.Fatal("Formula mal implementada")
		}
		if b.t == 2 {
			parent++
		} else if b.t == 3 {
			parent--
		}
		if b.t != 1 || a.t != 1 {
			c++
			arr = append(arr, Part{"", 0})
			num = ""
		}
		num = num + string(stack[i])
		arr[c-1] = Part{num, b.t}
	}
	if parent != 0 {
		log.Fatal("Error de parentesis")
	}
	return arr
}

func types(s string) int {
	for _, c := range s {
		if unicode.IsDigit(c) {
			return 1
		}
	}
	for _, n := range operations {
		if n.s == s {
			return n.t
		}
	}
	return 0
}

func typesReverse(n int) string {
	for _, s := range operations {
		if s.t == n {
			return s.s
		}
	}
	return ""
}
