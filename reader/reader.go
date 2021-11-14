package reader

import (
	"bufio"
	"os"
	"strings"
)

func Read(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	interfaceLines := make([]string, 0)
	isInterfaceDefinition := false

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if isPackage(line) {
			interfaceLines = append(interfaceLines, line)
		}

		if isInterfaceDefinitionStart(line) {
			isInterfaceDefinition = true
		}

		if isInterfaceDefinition {
			interfaceLines = append(interfaceLines, line)
		}

		if isInterfaceDefinition && isInterfaceDefinitionEnd(line) {
			isInterfaceDefinition = false
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return interfaceLines, nil
}

func isPackage(line string) bool {
	return strings.HasPrefix(line, "package")
}

func isInterfaceDefinitionStart(line string) bool {
	return strings.Contains(line, "type") && strings.Contains(line, "interface") && strings.Contains(line, "{")
}

func isInterfaceDefinitionEnd(line string) bool {
	return line == "}"
}
