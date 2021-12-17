package utils

import (
	"bufio"
	"os"
	"sort"
	"strconv"
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

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	MaybePanic(err)
	return i
}

type Reader func(scanner *bufio.Scanner)

func ReadInput(fileName string, read Reader) {
	file, err := os.Open(fileName)
	MaybePanic(err)
	defer func(file *os.File) {
		err := file.Close()
		MaybePanic(err)
	}(file)

	scanner := bufio.NewScanner(file)
	read(scanner)
}
