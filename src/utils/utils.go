package utils

import (
	"sort"
	"strings"
)

func MaybePanic(err interface{}) {
	if err != nil {
		panic(err)
	}
}

func StrSort(s string) string {
	x := strings.Split(s, "")
	sort.Strings(x)
	return strings.Join(x, "")
}

func StrCat(stringsArr ...string) string {
	sb := strings.Builder{}
	for _, s := range stringsArr {
		sb.WriteString(s)
	}
	return sb.String()
}

func InBounds(i int, length int) bool {
	return i >= 0 && i < length
}
