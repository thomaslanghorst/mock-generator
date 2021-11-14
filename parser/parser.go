package parser

import (
	"fmt"
	"mock-generator/writer"
	"strings"
)

func Parse(interfaceLines []string) writer.MockDefinition {

	mockDefinition := writer.MockDefinition{}
	funcs := make([]writer.FunctionDefinition, 0)

	for _, line := range interfaceLines {

		if isPackage(line) {
			mockDefinition.Package = parsePackage(line)
		}

		if isInterfaceDefinitionStart(line) {
			mockDefinition.Name = parseInterfaceName(line)
		}

		if isFunctionDefinition(line) {
			name := parseFunctionName(line)
			params := parseParams(line)
			returnValues := parseReturnValues(line)
			calledLine := parseCalledLine(params, returnValues)
			typeCastringLines, returnLine := parseTypeCastingLinesAndReturnLine(returnValues)

			funcs = append(funcs, writer.FunctionDefinition{
				Name:             name,
				Params:           params,
				ReturnValues:     returnValues,
				CalledLine:       calledLine,
				TypeCastingLines: typeCastringLines,
				ReturnLine:       returnLine,
			})
		}

	}

	mockDefinition.FunctionDefinitions = funcs
	return mockDefinition
}

func isPackage(line string) bool {
	return strings.HasPrefix(line, "package")
}

func parsePackage(line string) string {
	return strings.Split(line, " ")[1]
}

func parseInterfaceName(line string) string {
	// type InterfaceName interface
	//     ^             ^
	//     |             |
	//     i1            i2
	i1 := len("type ")
	i2 := strings.Index(line, " interface")
	return strings.Trim(line[i1:i2], " ")
}

func parseFunctionName(line string) string {
	i := strings.Index(line, "(")
	return strings.TrimSpace(line[0:i])
}

func parseParams(line string) string {
	// ListUsers() ([]User, error)
	//          ^^
	//          ||
	//        i1,i2
	i1 := strings.Index(line, "(")
	i2 := strings.Index(line, ")")
	return strings.Trim(line[i1+1:i2], " ")

}

func parseReturnValues(line string) string {
	i := strings.Index(line, ")")
	return strings.TrimSpace(line[i+1:])
}

func parseCalledLine(params string, returnValues string) string {
	called := ""
	paramNames := make([]string, 0)

	if len(returnValues) > 0 {
		called = "args := "
	}

	splits := strings.Split(params, ",")
	for _, split := range splits {
		trimmed := strings.TrimSpace(split)

		if strings.Contains(trimmed, " ") {
			// if params look line: s1 string, s2 string, s3 string
			paramNames = append(paramNames, strings.TrimSpace(strings.Split(trimmed, " ")[0]))
		} else {
			// if params look line: s1, s2, s3 string
			paramNames = append(paramNames, strings.TrimSpace(trimmed))
		}
	}

	called = fmt.Sprintf("%sm.Called(%s)", called, strings.Join(paramNames, ", "))

	return called
}

func parseTypeCastingLinesAndReturnLine(returnValues string) ([]string, string) {
	typeCastingLines := make([]string, 0)
	returnLineSplits := make([]string, 0)
	returnValueTypes := make([]string, 0)

	if len(returnValues) == 0 {
		return typeCastingLines, ""
	}

	if hasMultipleReturnTypes(returnValues) {
		vals := returnValues[1 : len(returnValues)-1] // strip ( and )
		returnValueTypes = strings.Split(vals, ",")
	} else {
		returnValueTypes = append(returnValueTypes, returnValues)
	}

	for i, returnValueType := range returnValueTypes {

		trimmedReturnValueType := strings.TrimSpace(returnValueType)

		if isRetrievable, funcName := isRetrievableType(trimmedReturnValueType); !isRetrievable {

			typeCastingLines = append(typeCastingLines, fmt.Sprintf("var v%d %s", i, trimmedReturnValueType))
			typeCastingLines = append(typeCastingLines, fmt.Sprintf("if args.Get(%d) != nil {", i))
			typeCastingLines = append(typeCastingLines, fmt.Sprintf("\tv%d = args.Get(%d).(%s)", i, i, trimmedReturnValueType))
			typeCastingLines = append(typeCastingLines, "}")

			returnLineSplits = append(returnLineSplits, fmt.Sprintf("v%d", i))

		} else {
			returnLineSplits = append(returnLineSplits, fmt.Sprintf("args%s%d)", funcName, i))
		}
	}

	returnLine := fmt.Sprintf("return %s", strings.Join(returnLineSplits, ", "))
	return typeCastingLines, returnLine

}

func hasMultipleReturnTypes(returnValues string) bool {
	return strings.Contains(returnValues, "(") && strings.Contains(returnValues, ")") && strings.Contains(returnValues, ",")
}

func isRetrievableType(returnValueType string) (bool, string) {
	// testify.mock has the following functions to retrieve values from arguments
	//args.Bool()
	//args.String()
	//args.Error()
	//args.Int()
	//args.Get()
	// if returnValueType is a premitive type, or error, the corresponding args.XYZ() function should be called
	// otherwise args.Get() should be called and must be casted

	switch returnValueType {
	case "bool":
		return true, ".Bool("
	case "string":
		return true, ".String("
	case "error":
		return true, ".Error("
	case "int":
		return true, ".Int("
	default:
		return false, ""
	}
}

func isInterfaceDefinitionStart(line string) bool {
	return strings.Contains(line, "type") && strings.Contains(line, "interface") && strings.Contains(line, "{")
}

func isFunctionDefinition(line string) bool {
	return strings.Contains(line, "(") && strings.Contains(line, ")")
}
