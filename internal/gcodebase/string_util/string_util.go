package string_util

import (
	"galaxy_walker/internal/gcodebase/time_util"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"encoding/json"
	"unicode/utf8"
)

func Purify(s string, dirty ...string) string {
	n := s
	for _, d := range dirty {
		n = strings.Replace(n, d, "", -1)
	}
	return n
}

func IsEmpty(s string) bool {
	return s == ""
}

// Reverse a utf8 encoded string.
func Reverse(str string) []byte {
	var size int

	tail := len(str)
	buf := make([]byte, tail)
	s := buf

	for len(str) > 0 {
		_, size = utf8.DecodeRuneInString(str)
		tail -= size
		s = append(s[:tail], []byte(str[:size])...)
		str = str[size:]
	}
	return buf
}

func StringAppendF(s *string, format string, a ...interface{}) {
	app := fmt.Sprintf(format, a...)
	*s = *s + app
}
func RandomIntn(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}
func Shuffle(slice *[]byte) {
	for i := range *slice {
		rand.Seed(time_util.GetCurrentTimeStamp())
		j := RandomIntn(i + 1)
		(*slice)[i], (*slice)[j] = (*slice)[j], (*slice)[i]
	}
}
func ShuffleInt(len int) []byte {
	sli := []byte{}
	for i := 0; i < len; i++ {
		sli = append(sli, byte(i))
	}
	Shuffle(&sli)
	return sli
}

func PrettyFormat(v interface{}) string {
	if v == nil {
		return ""
	}
	outs, _ := json.MarshalIndent(v, "", "\t")
	return string(outs)
}
