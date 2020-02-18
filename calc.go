package main

import (
	"errors"
	"fmt"
	"go/scanner"
	"go/token"
	"os"
	"regexp"
	"strconv"
)

var priorities = map[string]int{
	"low":     1,
	"middle":  2,
	"highest": 3,
}

type Operator struct {
	Priority  int
	isUnary   bool
	Operation func(args []float64) float64
}

var OperandsMap = map[string]*Operator{
	"+": {
		Priority:  priorities["low"],
		isUnary:   false,
		Operation: func(args []float64) float64 {
			return args[0] + args[1]
		},
	},
	"-": {
		Priority:  priorities["low"],
		isUnary:   false,
		Operation: func(args []float64) float64 {
			return args[0] - args[1]
		},
	},
	"*": {
		Priority:  priorities["middle"],
		isUnary:   false,
		Operation: func(args []float64) float64 {
			return args[0] * args[1]
		},
	},
	"/": {
		Priority:  priorities["middle"],
		isUnary:   false,
		Operation: func(args []float64) float64 {
			return args[0] / args[1]
		},
	},
	"unaryMinus": {
		Priority:  priorities["highest"],
		isUnary:   true,
		Operation: func(args []float64) float64 {
			return 0 - args[0]
		},
	},
	"unaryPlus": {
		Priority:  priorities["highest"],
		isUnary:   true,
		Operation: func(args []float64) float64 {
			return args[0]
		},
	},
}

type StringStack struct {
	stack     []string
	stackSize int
}
type FloatStack struct {
	stack     []float64
	stackSize int
}

func (s *StringStack) Push(a string) {
	s.stackSize++
	if s.stackSize < len(s.stack) {
		s.stack[s.stackSize] = a
	} else {
		s.stack = append(s.stack, a)
	}
}
func (s *StringStack) Pop() string {
	ret := s.Peek()
	s.stackSize--
	return ret
}

func (s *StringStack) Peek() string {
	if s.stackSize < 0 {
		return ""
	}
	return s.stack[s.stackSize]
}

func (s *FloatStack) Push(a float64) {
	s.stackSize++
	if s.stackSize < len(s.stack) {
		s.stack[s.stackSize] = a
	} else {
		s.stack = append(s.stack, a)
	}
}

func (s *FloatStack) Pop() float64 {
	ret := s.Peek()
	s.stackSize--
	return ret
}

func (s *FloatStack) Peek() float64 {
	if s.stackSize < 0 {
		return 0
	}
	return s.stack[s.stackSize]
}

func Solve(in string) (float64, error) {
	floats := &FloatStack{stackSize: -1, stack: []float64{}}
	ops := &StringStack{stackSize: -1, stack: []string{}}
	var s scanner.Scanner
	src := []byte(in)
	fileset := token.NewFileSet()
	file := fileset.AddFile("", fileset.Base(), len(src))
	s.Init(file, src, nil, 0)
	var prev token.Token

Loop:
	for {
		_, tok, lit := s.Scan()
		switch {
		case tok == token.EOF || tok == token.SEMICOLON:
			break Loop
		case isOperand(tok):
			val, err := strconv.ParseFloat(lit, 64)
			if err != nil {
				return 0, err
			}
			floats.Push(val)
			if prev == token.RPAREN { // Пример: (1+2+3+4)5 <=> 10 * 5
				err := solveExpression("*", ops, floats)
				if err != nil {
					return 0, err
				}
			}
		case isOperator(tok.String()):
			op := tok.String()
			if isUnary(tok, prev) { // чекаем унарные + и -
				if tok == token.SUB {
					op = "unaryMinus"
				} else {
					op = "unaryPlus"
				}
			}
			err := solveExpression(op, ops, floats)
			if err != nil {
				return 0, err
			}
		case tok == token.LPAREN: //Открытая скобка
			if isOperand(prev) { // Пример: 5(1+2+3+4) <=> 5 * 10
				err := solveExpression("*", ops, floats)
				if err != nil {
					return 0, err
				}
			}
			ops.Push(tok.String())
		case tok == token.RPAREN: //Закрытая скобка
			for ops.stackSize >= 0 && ops.Peek() != "(" {
				solveOperation(*OperandsMap[ops.Pop()], floats)
			}
			ops.Pop()
			//чекаем ++ -- ** //
		case tok == prev && (tok == token.SUB || tok == token.ADD || tok == token.MUL || tok == token.QUO):
			return 0, errors.New("duplicate operator token")
		default:
			fmt.Println("wrong input")
			return 0, errors.New("wrong token")
		}
		prev = tok
	}

	for ops.stackSize >= 0 {
		topOperand := ops.Pop()
		solveOperation(*OperandsMap[topOperand], floats)
	}
	res := floats.Peek()
	return res, nil
}

func solveExpression(op string, operationStack *StringStack, floats *FloatStack) error { //вычисляем операции со скобками
	peekedOperand := operationStack.Peek()
	popNext, err := popNext(op, peekedOperand)
	if err != nil {
		return err
	}
	for operationStack.stackSize >= 0 && popNext {
		solveOperation(*OperandsMap[operationStack.Pop()], floats)
	}
	operationStack.Push(op)
	return nil
}

func popNext(n1 string, n2 string) (bool, error) {
	if !isOperator(n2) {
		return false, nil
	}
	if n1 == "unaryMinus" {
		return false, nil
	} else if n1 == "unaryPlus" {
		return false, nil
	}
	op1, err := getOperator(n1)
	if err != nil {
		return false, err
	}
	op2, err := getOperator(n2)
	if err != nil {
		return false, err
	}
	return op1.Priority <= op2.Priority, nil
}

func solveOperation(op Operator, floats *FloatStack) { //вычисляем базовые операции
	var argsCount int
	if op.isUnary {
		argsCount = 1
	} else {
		argsCount = 2
	}
	var args = make([]float64, argsCount)
	for i := argsCount - 1; i >= 0; i-- {
		args[i] = floats.Pop()
	}
	floats.Push(op.Operation(args))
}

func isOperand(t token.Token) bool {
	return t == token.FLOAT || t == token.INT
}
func isOperator(s string) bool {
	_, exist := OperandsMap[s]
	return exist
}
func isUnary(tok token.Token, prev token.Token) bool {
	return (tok == token.SUB || tok == token.ADD) && (isOperator(prev.String()) || prev == token.LPAREN)
}

func getOperator(str string) (*Operator, error) {
	op, exist := OperandsMap[str]
	if exist {
		return op, nil
	}
	return nil, errors.New("wrong operator")
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func validation(input string) bool {
	regexWrongSymbols := regexp.MustCompile("[^\\+\\-\\*\\/\\)\\(0-9]")
	wrongSymbols := regexWrongSymbols.FindAllString(input, -1)
	if len(wrongSymbols) > 0 {
		return false
	}
	regexDuplicateSymbols := regexp.MustCompile("[\\+\\-\\*\\/\\(][\\+\\-\\*\\/\\)]")
	duplicateSymbols := regexDuplicateSymbols.FindAllString(input, -1)
	if len(duplicateSymbols) > 0 {
		return false
	}
	regexOpenBrackets := regexp.MustCompile("\\(")
	open := regexOpenBrackets.FindAllString(input, -1)
	regexClosedBrackets := regexp.MustCompile("\\)")
	closed := regexClosedBrackets.FindAllString(input, -1)
	if len(open) != len(closed) {
		return false
	}
	return true
}

func main() {
	args := os.Args
	var input string
	for i := 1; i < len(args); i++ {
		input += args[i]
	}

	if validation(input) {
		res, err := Solve(input)
		if err != nil {
			fmt.Println("Wrong input")
			return
		}
		fmt.Println(res)
		return
	}
	fmt.Println("Wrong input")
	return
}
