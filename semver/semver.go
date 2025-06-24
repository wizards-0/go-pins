package semver

import (
	"strconv"
	"strings"
)

// Returns true if v1 is less than or equal to v2
func CompareSemver(v1 string, v2 string, separator string) bool {
	v1Parts := strings.Split(v1, separator)
	v2Parts := strings.Split(v2, separator)
	partsCount := max(len(v1Parts), len(v2Parts))
	for i := range partsCount {
		v1Part := getVerPart(&v1Parts, i)
		v2Part := getVerPart(&v2Parts, i)
		if v1Part != v2Part {
			return v1Part < v2Part
		}
	}
	return true
}

func getVerPart(vParts *[]string, i int) int {
	var vPart int
	if i < len(*vParts) {
		if ver, err := strconv.Atoi((*vParts)[i]); err != nil {
			vPart = 9999
		} else {
			vPart = ver
		}
	} else {
		vPart = -1
	}
	return vPart
}
