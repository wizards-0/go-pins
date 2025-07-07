package props

import (
	"fmt"
	"os"
	"strings"

	"github.com/wizards-0/go-pins/logger"
)

func ReadFiles(filePaths ...string) (map[string]string, error) {
	props := map[string]string{}
	for _, filePath := range filePaths {
		pBytes, fileReadErr := os.ReadFile(filePath)
		if fileReadErr != nil {
			return nil, logger.WrapAndLogError(fileReadErr, "error in reading file "+filePath)
		}
		propFile := string(pBytes)
		var propLines []string
		if strings.Contains(propFile, "\r\n") {
			propLines = strings.Split(propFile, "\r\n")
		} else {
			propLines = strings.Split(propFile, "\n")
		}

		for _, prop := range propLines {
			propParts := strings.SplitN(prop, "=", 2)
			if len(propParts) != 2 {
				return nil, fmt.Errorf("invalid property %s, in file %s", prop, filePath)
			}
			props[propParts[0]] = strings.Trim(propParts[1], "\"")
		}
	}
	return props, nil
}
