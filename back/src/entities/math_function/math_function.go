package math_function

import (
	"errors"
	"fmt"
	"go_math_service/src/extensions/myregexp"
	"go_math_service/src/extensions/slice"
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
		variables = append(variables, variable)
	}
	sort.Strings(variables)
	variables = slice.RemoveDuplicates(variables)
	return variables
}

var numberChars = `-?(\d+[\.\,])?\d+`

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
	default:
		panic("Несуществующая математическая функция")
	}
}

func replaceBrackets(expression *string, variableValues map[string]float64, debug bool) bool {
	//
	r := regexp.MustCompile(`\W*\((?<inner>[^\(\)]+?)\)`)
	regex := myregexp.MyRegexp{Regexp: *r}
	matches := regex.FindAllStringSubmatch(*expression, -1)
	if len(matches) > 0 {
		var match = matches[0]
		var inner = regex.ValueByGroupName(match, "inner")
		fmt.Println("Inner ", inner)
		var value, err = resolveExpression(&inner, variableValues, debug)
		if err != nil {
			return false
		}
		*expression = strings.ReplaceAll(*expression, "("+inner+")", fmt.Sprint(value))
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
		a, _ := strconv.ParseFloat(regex.ValueByGroupName(match, "valueA"), 64)
		b, _ := strconv.ParseFloat(regex.ValueByGroupName(match, "valueB"), 64)
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

func (mathFunction *MathFunction) ResolveExpression(variableValues map[string]float64) (value float64, err error) {
	for _, variable := range mathFunction.variables {
		if variableValues[variable] == 0 {
			variableValues[variable] = 0
		}
	}
	return resolveExpression(&mathFunction.input, variableValues, mathFunction.debug)
}
func resolveExpression(expression *string, variableValues map[string]float64, debug bool) (value float64, err error) {
	for variable, value := range variableValues {
		*expression = strings.ReplaceAll(*expression, variable, fmt.Sprint(value))
	}
	for {
		_, err = strconv.ParseFloat(*expression, 64)
		if err == nil {
			break
		}

		if debug {
			fmt.Println(*expression)
		}
		if replaceBrackets(expression, variableValues, debug) {
			continue
		}
		if replaceFunction(expression, debug) {
			continue
		}
		if replaceBinary(expression, []string{"\\^"}) {
			continue
		}
		if replaceBinary(expression, []string{"\\*", "/"}) {
			continue
		}
		if replaceBinary(expression, []string{"\\+", "-"}) {
			continue
		}
		err = errors.New("выражение составлено некорректно")
		return
	}
	if debug {
		fmt.Println(*expression)
	}
	parsed, err := strconv.ParseFloat(*expression, 64)
	if err == nil {
		return parsed, nil
	}
	err = errors.New("выражение составлено некорректно")
	return
}

//     public float64 GetDerivative(map[string]float64 variableValues, char axis)
//     {
//         var delta = 1e-8;
//         return (
//             ResolveExpression(variableValues.Select(
//                 pair => pair.Key != axis ?
//                     pair :
//                      KeyValuePair<char, float64?>(pair.Key, pair.Value + delta)).ToDictionary()
//             ) -
//             ResolveExpression(variableValues.Select(
//                 pair => pair.Key != axis ?
//                     pair :
//                      KeyValuePair<char, float64?>(pair.Key, pair.Value - delta)).ToDictionary()
//             )
//          ) / (2 * delta);
//     }
//     bool replaceUnary( string expression, []string operators)
//     {
//         r := regexp.MustCompile("(?<operator>({string.Join("|", operators)}))(?<valueA>{numberChars})");
//         var match = regex.Match(expression);
//         if (match.Success)
//         {
//             expression = expression.Replace(
//                 match.Value,
//                     ArithmeticOperation(
//                         strconv.ParseFloat(regex.ValueByGroupName(match, "valueA")),
//                         regex.ValueByGroupName(match, "operator")
//                     )
//                 .ToString()
//             );
//             return true;
//         }
//         return false;
//     }
