package main

import (
	"github.com/stretchr/testify/assert"
	"go/token"
	"testing"
)

func TestSolveExpression(t *testing.T) {
	floatsTest := &FloatStack{stackSize: -1, stack: []float64{}}
	opsTest := &StringStack{stackSize: -1, stack: []string{}}

	floatsTest.Push(1)
	floatsTest.Push(1)
	opsTest.Push("*")

	err0 := solveExpression("+", opsTest, floatsTest)
	assert.Equal(t, 1.0, floatsTest.Peek())
	assert.Nil(t, err0, "No errors")

	floatsTest.Push(2)
	floatsTest.Push(4)
	opsTest.Push("/")

	err1 := solveExpression("+", opsTest, floatsTest)
	assert.Equal(t, 1.5, floatsTest.Peek())
	assert.Nil(t, err1, "No errors")

	floatsTest.Push(3)
	floatsTest.Push(5)
	opsTest.Push("+")

	err2 := solveExpression("+", opsTest, floatsTest)
	assert.Equal(t, 9.5, floatsTest.Peek())
	assert.Nil(t, err2, "No errors")

	floatsTest.Push(4)
	err3 := solveExpression("+", opsTest, floatsTest)
	assert.Equal(t, 13.5, floatsTest.Peek())
	assert.Nil(t, err3, "No errors")
}

func TestPopNext(t *testing.T) {
	res0, err := popNext("unaryMinus", "+")
	assert.Equal(t, false, res0)
	assert.Nil(t, err, "No errors")

	res1, err := popNext("unaryPlus", "+")
	assert.Equal(t, false, res1)
	assert.Nil(t, err, "No errors")

	res2, err := popNext("+", "Wrong operation")
	assert.Equal(t, false, res2)
	assert.Nil(t, err, "No errors")

	res3, err := popNext("-", "Wrong operation")
	assert.Equal(t, false, res3)
	assert.Nil(t, err, "No errors")

	res4, err := popNext("+", "*")
	assert.Equal(t, true, res4)
	assert.Nil(t, err, "No errors")

	res5, err := popNext("-", "/")
	assert.Equal(t, true, res5)
	assert.Nil(t, err, "No errors")

	res6, err := popNext("*", "unaryPlus")
	assert.Equal(t, true, res6)
	assert.Nil(t, err, "No errors")

	res7, err := popNext("/", "unaryMinus")
	assert.Equal(t, true, res7)
	assert.Nil(t, err, "No errors")

	res8, err := popNext("+", "+")
	assert.Equal(t, true, res8)
	assert.Nil(t, err, "Error with one priority(low) operators")

	res9, err := popNext("/", "/")
	assert.Equal(t, true, res9)
	assert.Nil(t, err, "Error with one priority(middle) operators")

	res10, err := popNext("unaryMinus", "unaryMinus")
	assert.Equal(t, false, res10)
	assert.Nil(t, err, "Error with unary operators")

}

func TestSolveOperation(t *testing.T) {
	floatsTest := &FloatStack{stackSize: -1, stack: []float64{}}
	floatsTest.Push(1)
	floatsTest.Push(1)
	solveOperation(*OperandsMap["+"], floatsTest)
	assert.Equal(t, 2.0, floatsTest.Peek())

	floatsTest.Push(3)
	solveOperation(*OperandsMap["-"], floatsTest)
	assert.Equal(t, -1.0, floatsTest.Peek())

	floatsTest.Push(-2)
	solveOperation(*OperandsMap["*"], floatsTest)
	assert.Equal(t, 2.0, floatsTest.Peek())

	floatsTest.Push(2)
	solveOperation(*OperandsMap["/"], floatsTest)
	assert.Equal(t, 1.0, floatsTest.Peek())

	solveOperation(*OperandsMap["unaryMinus"], floatsTest)
	assert.Equal(t, -1.0, floatsTest.Peek())

	solveOperation(*OperandsMap["unaryPlus"], floatsTest)
	assert.Equal(t, -1.0, floatsTest.Peek())

}

func TestIsOperand(t *testing.T) {
	assert.Equal(t, true, isOperand(token.FLOAT))
	assert.Equal(t, true, isOperand(token.INT))
	assert.Equal(t, false, isOperand(token.ADD))
	assert.Equal(t, false, isOperand(token.RPAREN))
	assert.Equal(t, false, isOperand(token.ILLEGAL))
}

func TestIsOperator(t *testing.T) {
	assert.Equal(t, true, isOperator("+"))
	assert.Equal(t, true, isOperator("unaryMinus"))
	assert.Equal(t, false, isOperator("1"))
	assert.Equal(t, false, isOperator("1.1"))
	assert.Equal(t, false, isOperator("#"))
}
func TestIsUnary(t *testing.T) {
	assert.Equal(t, true, isUnary(token.SUB, token.LPAREN))
	assert.Equal(t, true, isUnary(token.SUB, token.SUB))
	assert.Equal(t, true, isUnary(token.SUB, token.QUO))
	assert.Equal(t, true, isUnary(token.SUB, token.MUL))

	assert.Equal(t, true, isUnary(token.ADD, token.LPAREN))
	assert.Equal(t, true, isUnary(token.ADD, token.ADD))
	assert.Equal(t, true, isUnary(token.ADD, token.MUL))
	assert.Equal(t, true, isUnary(token.ADD, token.QUO))

	assert.Equal(t, false, isUnary(token.SUB, token.ILLEGAL))
	assert.Equal(t, false, isUnary(token.ADD, token.ILLEGAL))

	assert.Equal(t, false, isUnary(token.ILLEGAL, token.MUL))
	assert.Equal(t, false, isUnary(token.ILLEGAL, token.QUO))
	assert.Equal(t, false, isUnary(token.ILLEGAL, token.LPAREN))

}

func TestGetOperator(t *testing.T) {
	operator0, err := getOperator("+")
	assert.Nil(t, err, "No errors")
	assert.Equal(t, OperandsMap["+"], operator0)

	operator1, err := getOperator("unaryMinus")
	assert.Nil(t, err, "No errors")
	assert.Equal(t, OperandsMap["unaryMinus"], operator1)

	operator2, err := getOperator("$")
	assert.NotNil(t, err, "wrong operator")
	assert.Nil(t, operator2, "Nil operator")

	operator3, err := getOperator("SomeWrongThing")
	assert.NotNil(t, err, "wrong operator")
	assert.Nil(t, operator3, "Nil operator")
}

func TestValidation(t *testing.T) {
	assert.Equal(t, true, validation("1+2"))
	assert.Equal(t, true, validation("1+2*3"))
	assert.Equal(t, true, validation("(1+2)*3"))
	assert.Equal(t, true, validation("1+2+3*(1+(10/5))"))
	assert.Equal(t, true, validation("(1+2)*(3*(1/2))"))


	assert.Equal(t, false, validation("1+2+3*(1+(10/hello))"))
	assert.Equal(t, false, validation("1+2+3*(1+(10/5)))"))
	assert.Equal(t, false, validation("(1+2+3*(1+(10/5))"))
	assert.Equal(t, false, validation("1+2+3*(1+(10//5)))"))
	assert.Equal(t, false, validation("1++2+3*(1+(10/5)))"))
}

func TestSolve(t *testing.T) {
	res0, err := Solve("1+2")
	assert.Equal(t, 3.0, res0)
	assert.Nil(t, err, "No errors")

	res1, err := Solve("1+2*3")
	assert.Equal(t, 7.0, res1)
	assert.Nil(t, err, "No errors")

	res2, err := Solve("(1+2)*3")
	assert.Equal(t, 9.0, res2)
	assert.Nil(t, err, "No errors")

	res3, err := Solve("(1+2)*(3*(1/2))")
	assert.Equal(t, 4.5, res3)
	assert.Nil(t, err, "No errors")

	res4, err := Solve("( 1 +  2) * ( 3* (1+ 2 ))")
	assert.Equal(t, 27.0, res4)
	assert.Nil(t, err, "No errors")

	res5, err := Solve("1-(-1)")
	assert.Equal(t, 2.0, res5)
	assert.Nil(t, err, "No errors")

	res6, err := Solve("1-(-1)+(-1)*8")
	assert.Equal(t, -6.0, res6)
	assert.Nil(t, err, "No errors")


	res7, err := Solve("1+2+3*(1+(10/5))")
	assert.Equal(t, 12.0, res7)
	assert.Nil(t, err, "No errors")

}


