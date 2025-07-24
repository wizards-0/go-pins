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
		propLines := strings.Split(propFile, "\n")

		for _, prop := range propLines {
			prop = strings.TrimSpace(prop)
			if !strings.HasPrefix(prop, "#") && prop != "" {
				propParts := strings.SplitN(prop, "=", 2)
				if len(propParts) != 2 {
					return nil, fmt.Errorf("invalid property %s, in file %s", prop, filePath)
				}
				key := strings.TrimSpace(propParts[0])
				value := strings.Trim(strings.TrimSpace(propParts[1]), "\"")
				props[key] = value
			}
		}
	}
	return props, nil
}
