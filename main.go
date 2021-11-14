package main

import (
	"flag"
	"fmt"
	"mock-generator/parser"
	"mock-generator/reader"
	"mock-generator/writer"
	"os"
)

var (
	mockTemplateFile = "./mock.tmpl"
	inFile           string
	outFile          string
)

func init() {
	flag.StringVar(&inFile, "i", "", "input file")
	flag.StringVar(&outFile, "o", "<input_file>_mock.go", "output file")
}

func main() {

	flag.Parse()

	if len(inFile) == 0 {
		fmt.Println("Usage: mock-generator -i /path/to/input.go [-o /path/to/output.go]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if outFile == "<input_file>_mock.go" {
		outFile = fmt.Sprintf("%s_mock.go", inFile[:len(inFile)-3])
	}

	fmt.Printf("Generating mock for %s\n", inFile)

	interfaceLines, err := reader.Read(inFile)
	if err != nil {
		panic(err)
	}

	mockDefinition := parser.Parse(interfaceLines)

	writer.Write(outFile, mockTemplateFile, mockDefinition)

	fmt.Printf("Generated mock file %s\n", outFile)
}
