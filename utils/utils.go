package utils

import (
	"log"
	"strconv"
	"time"
)

func ConvertToInt(s string) int {
	// i, err := strconv.ParseInt(s, 10, 64)
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Print(err)
		return 0
	}
	return i
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func Elapsed(what string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %v\n", what, time.Since(start))
	}
}
