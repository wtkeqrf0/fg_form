package util

import (
	"math/rand"
	"strings"
)

func ParseCallBack(callback string) (string, string) {
	data := strings.SplitN(callback, "/", 2)
	switch len(data) {
	case 0:
		return "", ""
	case 1:
		return data[0], ""
	}
	return data[0], data[1]
}

func GenerateString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = rand.Int31n(26) + 97
	}
	return string(b)
}
