package migrator

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/wizards-0/go-pins/logger"
	"github.com/wizards-0/go-pins/migrator/types"
	"github.com/wizards-0/go-pins/semver"
)

func parseDirectory(path string) ([]types.Migration, error) {

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	verMigrationMap := map[string]types.Migration{}

	if err := addDirToMap(path, verMigrationMap); err != nil {
		return nil, fmt.Errorf("error while processing dir with path '%v'\n%w", path, err)
	}

	mArr := slices.Collect(maps.Values(verMigrationMap))
	sort.Slice(mArr, func(i1, i2 int) bool {
		return semver.CompareSemver(mArr[i1].Version, mArr[i2].Version, types.VERSION_SEPARATOR)
	})

	if err := validateMigrations(mArr); err != nil {
		return nil, fmt.Errorf("error while validating migrations\n%w", err)
	}
	return mArr, nil
}

func addDirToMap(path string, verMigrationMap map[string]types.Migration) error {

	entries, dirReadErr := os.ReadDir(path)
	if dirReadErr != nil {
		return logger.WrapAndLogError(dirReadErr, "error in reading directory - "+path)
	}

	for _, entry := range entries {
		if entry.Type().IsDir() {
			addDirToMap(path+entry.Name()+"/", verMigrationMap)
		} else {
			fileProcessErr := addFileToMap(path+entry.Name(), entry.Name(), verMigrationMap)
			if fileProcessErr != nil {
				return logger.WrapAndLogError(fileProcessErr, "error in processing file "+entry.Name())
			}
		}
	}

	return nil
}

func addFileToMap(filePath string, fileName string, verMigrationMap map[string]types.Migration) error {
	qBytes, fileReadErr := os.ReadFile(filePath)
	if fileReadErr != nil {
		return logger.WrapAndLogError(fileReadErr, "error in reading file "+filePath)
	}
	query := string(qBytes)

	ver, name, isQuery, fileNameErr := parseFileName(fileName)
	if fileNameErr != nil {
		return fmt.Errorf("error in parsing filename '%v'\n%w", fileName, fileNameErr)
	}

	m, versionExists := verMigrationMap[ver]
	if !versionExists {
		m = types.Migration{
			Version: ver,
			Name:    name,
		}
	} else {
		if name != m.Name {
			return nameMismatchError(m.Version+"-"+m.Name, ver+"-"+name)
		}
	}
	if isQuery {
		m.Query = query
	} else {
		m.Rollback = query
	}
	verMigrationMap[ver] = m
	return nil
}

func parseFileName(fileName string) (string, string, bool, error) {
	fileNameParts := strings.Split(fileName, ".")

	if len(fileNameParts) != 4 {
		return "", "", false, fileNameError(fileName)
	}

	if fileNameParts[3] != "sql" {
		return "", "", false, fileTypeError(fileName)
	}

	ver := fileNameParts[0]
	name := fileNameParts[1]
	var isQuery bool
	switch fileNameParts[2] {
	case "query":
		isQuery = true
	case "rollback":
		isQuery = false
	default:
		return "", "", false, fileNameError(fileName)
	}
	return ver, name, isQuery, nil
}

func validateMigrations(mArr []types.Migration) error {
	for _, m := range mArr {
		if len(m.Query) == 0 {
			return missingQuery(m)
		}
		if len(m.Rollback) == 0 {
			return missingRollback(m)
		}
	}
	return nil
}

func fileNameError(fileName string) error {
	err := fmt.Errorf("invalid filename - %v . File name has to be of format 'ver.name.query|rollback.sql'. E.g. 1-1.user-setup.query.sql, 1-1.user-setup.rollback.sql", fileName)
	return logger.LogError(err)
}

func fileTypeError(fileName string) error {
	err := fmt.Errorf("invalid filename - %v . Only files with sql extensions are supported. E.g. 1-1.user-setup.query.sql, 1-1.user-setup.rollback.sql", fileName)
	return logger.LogError(err)
}

func nameMismatchError(n1 string, n2 string) error {
	err := fmt.Errorf("invalid filenames - %v %v . Version and Name must be same for query and rollback files", n1, n2)
	return logger.LogError(err)
}

func missingQuery(m types.Migration) error {
	err := fmt.Errorf("missing query file for Version: %v, Name: %v ", m.Version, m.Name)
	return logger.LogError(err)
}

func missingRollback(m types.Migration) error {
	err := fmt.Errorf("missing rollback file for Version: %v, Name: %v ", m.Version, m.Name)
	return logger.LogError(err)
}
