package semver

import (
	"errors"
	"strconv"
	"strings"
)

// Returns true if v1 is less than or equal to v2
func CompareSemver(v1 string, v2 string, separator string) bool {
	v1Parts := strings.Split(v1, separator)
	v2Parts := strings.Split(v2, separator)
	partsCount := max(len(v1Parts), len(v2Parts))
	for i := range partsCount {
		v1Part, err1 := getVerPart(&v1Parts, i)
		v2Part, err2 := getVerPart(&v2Parts, i)
		if err1 != nil || err2 != nil {
			return v1Parts[i] < v2Parts[i]
		} else {
			if v1Part != v2Part {
				return v1Part < v2Part
			}
		}
	}
	return true
}

func getVerPart(vParts *[]string, i int) (int, error) {
	var vPart int
	if i < len(*vParts) {
		if ver, err := strconv.Atoi((*vParts)[i]); err != nil {
			return 0, errors.New("not a number")
		} else {
			vPart = ver
		}
	} else {
		vPart = -1
	}
	return vPart, nil
}
