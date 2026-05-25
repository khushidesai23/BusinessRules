package evaluator

import (
	"calculationengine/service/parser"
	"errors"
	"fmt"
	"math"
	"strconv"
)

var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

func Eval(node parser.Node, env *Environment) Object {
	switch node := node.(type) {
	case *parser.ExpressionStatement:
		return Eval(node.Expression, env)

	case *parser.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *parser.IntegerLiteral:
		return &Integer{Value: int64(node.Value)}

	case *parser.FloatLiteral:
		return &Float{Value: node.Value}

	case *parser.Identifier:
		return evalIdentifier(node, env)

	case *parser.Boolean:
		return NativeBoolToBooleanObject(node.Value)

	case *parser.StringLiteral:
		return &String{Value: node.Value}

	case *parser.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		value := evalInfixExpression(node.Operator, left, right)
		return value

	case *parser.IfExpression:
		return evalIfExpression(node, env)

	case *parser.CallExpression:
		return evalCallExpression(node, env)
	}
	return newError("Unknown parse function")
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

func evalIdentifier(i *parser.Identifier, env *Environment) Object {
	value, ok := env.Get(i.Value)
	if !ok {
		return newError("Empty value in variable %s", i.Value)
	}
	return value
}

func evalIfExpression(ie *parser.IfExpression, env *Environment) Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	isTruth, err := isTruthy(condition)
	if err != nil {
		return newError("Invalid Condition: %s", ie.Condition.String())
	}
	if isTruth {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(object Object) (bool, error) {
	switch {
	case object == TRUE:
		return true, nil
	case object == FALSE:
		return false, nil
	case object.Type() == INTEGER_OBJ && object.(*Integer).Value == 0:
		return false, nil
	case object.Type() == INTEGER_OBJ && object.(*Integer).Value > 0:
		return true, nil
	default:
		return false, errors.New("Invalid Condition")
	}
}

func evalInfixExpression(operator string, left Object, right Object) Object {
	if left.Type() == STRING_OBJ {
		parsedLeft := evalStringAsNummber(left.(*String).Value)
		if parsedLeft != nil {
			left = parsedLeft
		}
	}
	if right.Type() == STRING_OBJ {
		parsedRight := evalStringAsNummber(right.(*String).Value)
		if parsedRight != nil {
			right = parsedRight
		}
	}
	switch {
	case left.Type() == FLOAT_OBJ && right.Type() == FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)

	case left.Type() == FLOAT_OBJ && right.Type() == INTEGER_OBJ:
		rightAsFloat := &Float{Value: float64(right.(*Integer).Value)}
		return evalFloatInfixExpression(operator, left, rightAsFloat)

	case left.Type() == INTEGER_OBJ && right.Type() == FLOAT_OBJ:
		leftAsFloat := &Float{Value: float64(left.(*Integer).Value)}
		return evalIntegerInfixExpression(operator, leftAsFloat, right)

	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	/*
		We’re using pointer comparison here to check for equality between booleans. That works
		because we’re always using pointers to our objects and in the case of booleans we only ever use
		two: TRUE and FALSE. So, if something has the same value as TRUE (the memory address that is)
		then it’s true.
	*/
	case operator == "=":
		if left.Type() != right.Type() {
			return FALSE // Return false by default if types mismatch
		}
		return NativeBoolToBooleanObject(left == right)

	case operator == "<>":
		if left.Type() != right.Type() {
			return TRUE // Return true by default if types mismatch
		}
		return NativeBoolToBooleanObject(left != right)

	case left.Type() != right.Type():
		return newError("Type mismatch: %s %s %s", left.Type(), operator, right.Type())

	default:
		return newError("Unknown Operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringAsNummber(s string) Object {
	if intVal, err := strconv.ParseInt(s, 10, 64); err == nil {
		return &Integer{Value: intVal}
	}
	if floatVal, err := strconv.ParseFloat(s, 64); err == nil {
		return &Float{Value: floatVal}
	}
	return nil
}

func evalIntegerInfixExpression(operator string, left Object, right Object) Object {
	leftVal := left.(*Integer).Value
	rightVal := right.(*Integer).Value
	switch operator {
	case "+":
		return &Integer{Value: (leftVal + rightVal)}
	case "-":
		return &Integer{Value: (leftVal - rightVal)}
	case "*":
		return &Integer{Value: (leftVal * rightVal)}
	case "/":
		if rightVal == 0 {
			return newError("Division by zero not allowed")
		}
		return &Integer{Value: (leftVal / rightVal)}
	case "%":
		if rightVal == 0 {
			return newError("Modulo by zero not allowed")
		}
		return &Integer{Value: (leftVal % rightVal)}
	case "^":
		if rightVal < 0 {
			return newError("Negative exponents not supported for integers")
		}
		result := int64(1)
		for i := int64(0); i < rightVal; i++ {
			result *= leftVal
		}
		return &Integer{Value: result}
	case ">":
		return NativeBoolToBooleanObject(leftVal > rightVal)
	case "<":
		return NativeBoolToBooleanObject(leftVal < rightVal)
	case "<>":
		return NativeBoolToBooleanObject(leftVal != rightVal)
	case "=":
		return NativeBoolToBooleanObject(leftVal == rightVal)
	default:
		return newError("Unknown Operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left Object, right Object) Object {
	leftVal := left.(*Float).Value
	rightVal := right.(*Float).Value
	switch operator {
	case "+":
		return &Float{Value: (leftVal + rightVal)}
	case "-":
		return &Float{Value: (leftVal - rightVal)}
	case "*":
		return &Float{Value: (leftVal * rightVal)}
	case "/":
		if rightVal == 0.0 {
			return newError("Division by zero not allowed")
		}
		return &Float{Value: (leftVal / rightVal)}
	case "%":
		if rightVal == 0.0 {
			return newError("Modulo by zero not allowed")
		}
		return &Float{Value: math.Mod(leftVal, rightVal)}
	case "^":
		return &Float{Value: math.Pow(leftVal, rightVal)}
	case ">":
		return NativeBoolToBooleanObject(leftVal > rightVal)
	case "<":
		return NativeBoolToBooleanObject(leftVal < rightVal)
	case "<>":
		return NativeBoolToBooleanObject(leftVal != rightVal)
	case "=":
		return NativeBoolToBooleanObject(leftVal == rightVal)
	default:
		return newError("Unknown Operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left Object, right Object) Object {
	leftVal := left.(*String).Value
	rightVal := right.(*String).Value
	switch operator {
	case "<>":
		return NativeBoolToBooleanObject(leftVal != rightVal)
	case "=":
		return NativeBoolToBooleanObject(leftVal == rightVal)
	default:
		return newError("Unknown Operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func NativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("Unknown Operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != INTEGER_OBJ {
		return newError("Unknown operator: -%s", right.Type())
	}
	value := right.(*Integer).Value
	return &Integer{Value: -value}
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

// evalCallExpression evaluates function calls like MIN, MAX, ROUND
func evalCallExpression(call *parser.CallExpression, env *Environment) Object {
	funcIdent, ok := call.Function.(*parser.Identifier)
	if !ok {
		return newError("Invalid function call")
	}

	args := evalExpressions(call.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	return applyFunction(funcIdent.Value, args)
}

// evalExpressions evaluates a slice of expressions
func evalExpressions(exps []parser.Expression, env *Environment) []Object {
	var result []Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

// applyFunction applies built-in functions (MIN, MAX, ROUND)
func applyFunction(name string, args []Object) Object {
	switch name {
	case "MIN":
		return builtinMin(args)
	case "MAX":
		return builtinMax(args)
	case "ROUND":
		return builtinRound(args)
	default:
		return newError("Unknown function: %s", name)
	}
}

// builtinMin returns the minimum value from arguments
func builtinMin(args []Object) Object {
	if len(args) == 0 {
		return newError("MIN requires at least 1 argument")
	}

	var minVal float64
	firstNum := true

	for _, arg := range args {
		var num float64
		switch arg.Type() {
		case INTEGER_OBJ:
			num = float64(arg.(*Integer).Value)
		case FLOAT_OBJ:
			num = arg.(*Float).Value
		default:
			return newError("MIN arguments must be numeric")
		}

		if firstNum {
			minVal = num
			firstNum = false
		} else if num < minVal {
			minVal = num
		}
	}

	if allIntegers(args) {
		return &Integer{Value: int64(minVal)}
	}
	return &Float{Value: minVal}
}

// builtinMax returns the maximum value from arguments
func builtinMax(args []Object) Object {
	if len(args) == 0 {
		return newError("MAX requires at least 1 argument")
	}

	var maxVal float64
	firstNum := true

	for _, arg := range args {
		var num float64
		switch arg.Type() {
		case INTEGER_OBJ:
			num = float64(arg.(*Integer).Value)
		case FLOAT_OBJ:
			num = arg.(*Float).Value
		default:
			return newError("MAX arguments must be numeric")
		}

		if firstNum {
			maxVal = num
			firstNum = false
		} else if num > maxVal {
			maxVal = num
		}
	}

	if allIntegers(args) {
		return &Integer{Value: int64(maxVal)}
	}
	return &Float{Value: maxVal}
}

// builtinRound rounds a number to specified decimal places
// ROUND(number, decimals) - if decimals not provided, rounds to nearest integer
func builtinRound(args []Object) Object {
	if len(args) == 0 || len(args) > 2 {
		return newError("ROUND requires 1 or 2 arguments")
	}

	var number float64
	switch args[0].Type() {
	case INTEGER_OBJ:
		number = float64(args[0].(*Integer).Value)
	case FLOAT_OBJ:
		number = args[0].(*Float).Value
	default:
		return newError("ROUND first argument must be numeric")
	}

	decimals := int64(0)
	if len(args) == 2 {
		if args[1].Type() != INTEGER_OBJ {
			return newError("ROUND second argument must be an integer")
		}
		decimals = args[1].(*Integer).Value
	}

	multiplier := math.Pow(10, float64(decimals))
	rounded := math.Round(number*multiplier) / multiplier

	if decimals == 0 && allIntegers([]Object{args[0]}) {
		return &Integer{Value: int64(rounded)}
	}
	return &Float{Value: rounded}
}

// allIntegers checks if all arguments are integers
func allIntegers(args []Object) bool {
	for _, arg := range args {
		if arg.Type() != INTEGER_OBJ {
			return false
		}
	}
	return true
}
