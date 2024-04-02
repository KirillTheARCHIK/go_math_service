package math_function

import (
	"errors"
	"fmt"
	"go_math_service/src/extensions/mymath"
	"go_math_service/src/extensions/myregexp"
	"go_math_service/src/extensions/slice"
	"maps"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type MathFunction struct {
	debug     bool
	input     string
	variables []string
}

func Constructor(input string, debug bool) MathFunction {
	this := MathFunction{}
	this.debug = debug
	this.input = input

	this.variables = getVariables(input)
	if debug {
		fmt.Println("variables: ", this.variables)
	}
	return this
}

func getVariables(input string) []string {
	variables := []string{}
	r := regexp.MustCompile(`(?<variable>[A-Z_]+)`)
	regex := myregexp.MyRegexp{Regexp: *r}
	matches := regex.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		var variable = regex.ValueByGroupName(match, "variable")
		if variable == "E" || variable == "PI" {
			continue
		}
		variables = append(variables, variable)
	}
	sort.Strings(variables)
	variables = slice.RemoveDuplicates(variables)
	return variables
}

var floatNumber = `-?(\d+[\.\,])?\d+`
var numberChars = `-?` + floatNumber + `|\(` + floatNumber + `\)`

func binaryArithmeticOperation(a, b float64, operation string) float64 {
	switch operation {
	case "+":
		{
			return a + b
		}
	case "-":
		{
			return a - b
		}
	case "*":
		{
			return a * b
		}
	case "/":
		{
			return a / b
		}

	case "^":
		{
			return math.Pow(a, b)
		}
	default:
		panic("Несуществующий бинарный оператор")
	}
}

func arithmeticFunction(a float64, fName string) float64 {
	switch fName {
	case "sin":
		{
			return math.Sin(a)
		}
	case "cos":
		{
			return math.Cos(a)
		}
	case "fact":
		{
			return float64(mymath.Fact(int(a)))
		}
	case "exp":
		{
			return math.Exp(a)
		}
	default:
		panic("Несуществующая математическая функция")
	}
}

func replaceBrackets(expression *string, variableValues map[string]float64, debug bool) bool {
	//
	// r := regexp.MustCompile(`\W+\((?<inner>[^\(\)]+?)\)`)
	r := regexp.MustCompile(`\(\((?<inner>` + floatNumber + `)\)\)`)
	regex := myregexp.MyRegexp{Regexp: *r}
	matches := regex.FindAllStringSubmatch(*expression, -1)
	if len(matches) > 0 {
		var match = matches[0]
		var inner = regex.ValueByGroupName(match, "inner")
		// if debug {
		// 	fmt.Println("Inner", inner)
		// }
		// var value, err = resolveExpression(inner, variableValues, debug)
		// if err != nil {
		// 	return false
		// }
		// var old string = "(" + inner + ")"
		// if debug {
		// 	fmt.Println("Value", value)
		// 	fmt.Println("Old", old)
		// }
		// *expression = strings.ReplaceAll(*expression, old, fmt.Sprint(value))
		*expression = strings.ReplaceAll(*expression, "(("+inner+"))", "("+inner+")")
		return true
	}

	r = regexp.MustCompile(`\d+-(?<right>(\d+[\.\,])?\d+)`)
	regex = myregexp.MyRegexp{Regexp: *r}
	matches = regex.FindAllStringSubmatch(*expression, -1)
	if len(matches) > 0 {
		var match = matches[0]
		var right = regex.ValueByGroupName(match, "right")
		*expression = strings.ReplaceAll(*expression, right, "("+right+")")
		return true
	}
	return false
}

func replaceBinary(expression *string, operators []string) bool {
	r := regexp.MustCompile("(?<valueA>" + numberChars + ")(?<operator>[" + strings.Join(operators, "") + "])(?<valueB>" + numberChars + ")")
	regex := myregexp.MyRegexp{Regexp: *r}
	matches := regex.FindAllStringSubmatch(*expression, -1)
	if len(matches) > 0 {
		match := matches[0]
		aStr := regex.ValueByGroupName(match, "valueA")
		aStr = strings.ReplaceAll(aStr, "(", "")
		aStr = strings.ReplaceAll(aStr, ")", "")
		a, _ := strconv.ParseFloat(aStr, 64)
		bStr := regex.ValueByGroupName(match, "valueB")
		bStr = strings.ReplaceAll(bStr, "(", "")
		bStr = strings.ReplaceAll(bStr, ")", "")
		b, _ := strconv.ParseFloat(bStr, 64)

		// fmt.Println(match[0], aStr, bStr)

		operator := regex.ValueByGroupName(match, "operator")
		*expression = strings.ReplaceAll(*expression, match[0], fmt.Sprint(binaryArithmeticOperation(a, b, operator)))
		return true
	}
	return false
}

func replaceFunction(expression *string, debug bool) bool {
	r := regexp.MustCompile(`(?<fname>[a-z]+)\((?<valueA>` + numberChars + `)\)`)
	regex := myregexp.MyRegexp{Regexp: *r}
	matches := regex.FindAllStringSubmatch(*expression, -1)
	if len(matches) > 0 {
		match := matches[0]
		a, _ := strconv.ParseFloat(regex.ValueByGroupName(match, "valueA"), 64)
		fname := regex.ValueByGroupName(match, "fname")
		*expression = strings.ReplaceAll(*expression, match[0], fmt.Sprint(arithmeticFunction(a, fname)))
		return true
	}
	return false
}

func (mathFunction *MathFunction) ResolveExpression(variableValues map[string]float64, debug bool) (value float64, err error) {
	for _, variable := range mathFunction.variables {
		if variableValues[variable] == 0 {
			variableValues[variable] = 0
		}
	}
	return resolveExpression(mathFunction.input, variableValues, debug)
}
func resolveExpression(expression string, variableValues map[string]float64, debug bool) (value float64, err error) {
	for variable, value := range variableValues {
		expression = strings.ReplaceAll(expression, variable, fmt.Sprint(value))
	}
	expression = strings.ReplaceAll(expression, "E", fmt.Sprint(math.E))
	expression = strings.ReplaceAll(expression, "PI", fmt.Sprint(math.Pi))
	for {
		_, err = strconv.ParseFloat((expression)[1:], 64)
		if err == nil {
			expression = (expression)[1:]
			break
		}

		if debug {
			fmt.Println(expression)
		}
		if replaceBrackets(&expression, variableValues, debug) {
			continue
		}
		if replaceFunction(&expression, debug) {
			continue
		}
		if replaceBinary(&expression, []string{"\\^"}) {
			continue
		}
		if replaceBinary(&expression, []string{"\\*", "/"}) {
			continue
		}
		if replaceBinary(&expression, []string{"\\+", "-"}) {
			continue
		}
	}
	if debug {
		fmt.Println(expression)
	}
	parsed, err := strconv.ParseFloat(expression, 64)
	if err == nil {
		return parsed, nil
	}
	err = errors.New("выражение составлено некорректно")
	return
}

func (mathFunction *MathFunction) GetDerivative(variableValues map[string]float64, axis string, debug bool) (value float64, err error) {
	var delta = 1e-8
	var vV1 = maps.Clone(variableValues)
	vV1[axis] += delta
	var vV2 = maps.Clone(variableValues)
	vV2[axis] -= delta

	if mathFunction.debug {
		fmt.Println("vV1:", vV1, "vV2:", vV2)
	}

	f1, err := mathFunction.ResolveExpression(vV1, debug)
	if err != nil {
		return
	}
	f2, err := mathFunction.ResolveExpression(vV2, debug)
	if err != nil {
		return
	}

	if mathFunction.debug {
		fmt.Println("f1:", f1, "f2:", f2)
	}

	value = (f1 - f2) / (2 * delta)
	return
}

func (mathFunction *MathFunction) FindRootsDividing(a, b, eps float64) (value []float64, err error) {
	if len(mathFunction.variables) != 1 {
		return nil, errors.New("уравнение должно содержать только одну переменную, кроме математических констант")
	}
	for iterations := 1; ; iterations++ {
		var mid = (a + b) / 2
		if math.Abs(b-a) < eps {
			if mathFunction.debug {
				fmt.Println("Iterations:", iterations)
			}
			return []float64{mid}, nil
		}
		if mathFunction.debug {
			fmt.Println(fmt.Sprint(iterations)+")", "a =", a, "mid =", mid, "b =", b)
		}
		//f(a)
		fa, err := mathFunction.ResolveExpression(map[string]float64{mathFunction.variables[0]: a}, false)
		if err != nil {
			return nil, err
		}
		//f(mid)
		fmid, err := mathFunction.ResolveExpression(map[string]float64{mathFunction.variables[0]: mid}, false)
		if err != nil {
			return nil, err
		}
		if mathFunction.debug {
			fmt.Println("fa =", fa, "fmid =", fmid)
		}
		if fa == 0 {
			return []float64{fa}, nil
		}
		if fmid == 0 {
			return []float64{fmid}, nil
		}
		if fa*fmid < 0 {
			b = mid
			continue
		}

		//f(b)
		fb, err := mathFunction.ResolveExpression(map[string]float64{mathFunction.variables[0]: b}, false)
		if err != nil {
			return nil, err
		}
		if mathFunction.debug {
			fmt.Println("fb =", fb)
		}
		if fb == 0 {
			return []float64{fb}, nil
		}
		if fb*fmid < 0 {
			a = mid
			continue
		}

		return nil, errors.New("нет корней на данном отрезке")
	}
}

func (mathFunction *MathFunction) FindRootsSimple(eps float64) (value []float64, err error) {
	if len(mathFunction.variables) != 1 {
		return nil, errors.New("уравнение должно содержать только одну переменную, кроме математических констант")
	}
	var fx float64 = 0
	for iterations := 1; ; iterations++ {
		if mathFunction.debug {
			fmt.Println(fmt.Sprint(iterations)+")", fx)
		}
		newFx, err := mathFunction.ResolveExpression(map[string]float64{mathFunction.variables[0]: fx}, false)
		if err != nil {
			return nil, err
		}
		if math.Abs(newFx-fx) < eps {
			fmt.Println("Iterations:", iterations)
			return []float64{newFx}, nil
		}
		fx = newFx
	}
}

func (mathFunction *MathFunction) FindRootsNewton(a, b, eps float64) (value []float64, err error) {
	if len(mathFunction.variables) != 1 {
		return nil, errors.New("уравнение должно содержать только одну переменную, кроме математических констант")
	}
	return nil, nil
}
