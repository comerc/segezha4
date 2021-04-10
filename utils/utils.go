package utils

import (
	"log"
	"math"
	"net/url"
	"os"
	"strconv"
	"strings"
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

var timeoutFactor int

func InitTimeoutFactor() {
	timeoutFactor = ConvertToInt(os.Getenv("SEGEZHA4_TIMEOUT_FACTOR"))
	if timeoutFactor == 0 {
		timeoutFactor = 100
	}
}

func GetTimeout(average int) time.Duration {
	f := (float64(average) / 100) * float64(timeoutFactor)
	return time.Duration(math.Round(f)) * time.Second
}

func GetHost(linkURL string) string {
	u, err := url.Parse(linkURL)
	if err != nil {
		log.Println(err)
	}
	if strings.HasPrefix(u.Host, "www.") {
		return u.Host[4:]
	}
	return u.Host
}
