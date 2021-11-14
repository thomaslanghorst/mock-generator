package writer

type MockDefinition struct {
	Package             string
	Name                string
	FunctionDefinitions []FunctionDefinition
}

type FunctionDefinition struct {
	Name             string
	Params           string
	ReturnValues     string
	CalledLine       string
	TypeCastingLines []string
	ReturnLine       string
}
