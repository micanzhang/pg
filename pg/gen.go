package pg

import (
	"math/rand"
	"time"
)

var (
	letters = [26]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
	numbers = [10]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	chars   = []byte{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-', '_', '+', '=', '[', ']', '{'}
)

func GenRandPassword(len, uCaseLen, lCaseLen, numLen, puncLen int) string {
	minLen := uCaseLen + lCaseLen + numLen + puncLen
	if minLen > len {
		len = minLen
	}

	dst := make([]byte, 0)

	for i := len - minLen; i > 0; i-- {
		if i%4 == 0 {
			puncLen++
		} else if i%3 == 0 {
			numLen++
		} else if i%2 == 0 {
			lCaseLen++
		} else {
			uCaseLen++
		}
	}

	randSpecficBytes := func(n int, randFunc func() byte) []byte {
		bytes := make([]byte, n)
		for i := 0; i < n; i++ {
			bytes[i] = randFunc()
		}
		return bytes
	}

	dst = append(randSpecficBytes(uCaseLen, randUpperCaseLetter), dst...)
	dst = append(randSpecficBytes(lCaseLen, randLowerCaseLetter), dst...)
	dst = append(randSpecficBytes(numLen, randNumber), dst...)
	dst = append(randSpecficBytes(puncLen, randPunc), dst...)

	for i := 0; i < len; i++ {
		m := randInt(len)
		n := randInt(len)
		dst = exchange(dst, m, n)
	}

	return string(dst)
}

func randLowerCaseLetter() byte {
	return letters[randInt(26)]
}

func randUpperCaseLetter() byte {
	return randLowerCaseLetter() - 32
}

func randNumber() byte {
	return numbers[randInt(10)]
}

func randPunc() byte {
	return chars[randInt(len(chars))]
}

func randInt(n int) int {
	return int(rand.Int63n(time.Now().UnixNano()) % int64(n))
}

func exchange(bytes []byte, m, n int) []byte {
	c := bytes[m]
	bytes[m] = bytes[n]
	bytes[n] = c
	return bytes
}
