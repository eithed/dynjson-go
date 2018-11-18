package common

import (
	"math/rand"
	"strconv"
	"time"
)

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImprSrc(n int) string {

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func Itos(i interface{}) string {
	switch i.(type) {
	case bool:
		return strconv.FormatBool(i.(bool))
	case int:
		return strconv.FormatInt(i.(int64), 10)
	case int16:
		return strconv.FormatInt(i.(int64), 10)
	case int32:
		return strconv.FormatInt(i.(int64), 10)
	case int64:
		return strconv.FormatInt(i.(int64), 10)
	case float32:
		return strconv.FormatFloat(i.(float64), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(i.(float64), 'f', -1, 64)
	case string:
		return i.(string)
	}

	return ""
}
