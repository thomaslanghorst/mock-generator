package writer

import (
	"os"
	"text/template"
)

func Write(outFile, mockFile string, definition MockDefinition) {
	t, err := template.ParseFiles(mockFile)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = t.Execute(file, definition)
	if err != nil {
		panic(err)
	}

}
