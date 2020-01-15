package controller

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

// GenerateCode represent configuration and sync
// for generate code
type GenerateCode struct {
	cacheMapWg    sync.WaitGroup
	cacheMapMu    sync.Mutex
	cacheMap      map[string]interface{}
	randomSrc     rand.Source
	numberRoutine int64
	codeLength    int64
	letterConsume string
	processFunc   func(string)
	listBuffer    []interface{}
}

const (
	defaultLetter = "abcdefghijklmnopqrstuvwxyzQWERTYUIOPASDFGHJKLZXCVBNM0123456789"
	urlSafeLetter = "abcdefghijklmnopqrstuvwxyzQWERTYUIOPASDFGHJKLZXCVBNM0123456789-_"
)

// NewGenerateCodeRoutine returns a GenerateVoucher
// f for listener each process routines
// fixed configuration given length for fixed character code length
// routine number of generated code
// each of code will be using 1 routines
func NewGenerateCodeRoutine(f func(str string), length, routine int64, i []interface{}) *GenerateCode {
	return &GenerateCode{randomSrc: rand.NewSource(time.Now().UnixNano()),
		cacheMap:      make(map[string]interface{}),
		numberRoutine: routine,
		codeLength:    length,
		letterConsume: defaultLetter,
		processFunc:   f,
		listBuffer:    i}
}

// NewGenerateCode basic generate random string
// return unique array of string with GetUniqueStrings()
func NewGenerateCode() *GenerateCode {
	return &GenerateCode{
		randomSrc:     rand.NewSource(time.Now().UnixNano()),
		cacheMap:      make(map[string]interface{}),
		letterConsume: defaultLetter,
	}
}

// make sure app.main() not terminated before all routines complete
func (c *GenerateCode) start() *GenerateCode {
	c.cacheMapWg.Add(int(c.numberRoutine))
	for i := 0; i < int(c.numberRoutine); i++ {
		go c.process()
	}
	c.cacheMapWg.Wait()
	return c
}

// Duplicate code will be re-generated until unique found as many times posible. May be changed next plan.
// Calculation unique string len(letterConsume)^codeLength,
// ex. letterConsume = qweQWE123, len(letterConsume) = 9, codeLength = 4
// then max Unique Code can be generated is 9^4 = 6.561
// TODO: what happen if unique string cache reached max try.
func (c *GenerateCode) process() *GenerateCode {
	c.processFunc(c.getUniqueString())
	c.cacheMapWg.Done()
	return c
}

//GetUniqueStrings : generate string of length by [size]string arrays
func (c *GenerateCode) GetUniqueStrings(length, size int) []string {
	r := make([]string, size)
	for i := 0; i < size; i++ {
		r = append(r, c.getUniqueString())
	}
	return r
}

func (c *GenerateCode) getUniqueString() (str string) {
	for {
		str = string(randString(c.letterConsume, int(c.codeLength), c.randomSrc))
		if c.getCacheMapValue(str) != nil {
			// fmt.Println("WARNING!!DUPLICATE.")
		} else {
			c.setCacheMap(str, "")
			break
		}
	}
	return
}

// setLetterTemplate config set for letter template to be generated
// must set before start() executed
func (c *GenerateCode) setLetterTemplate(letter string) *GenerateCode {
	c.letterConsume = letter
	return c
}

func (c *GenerateCode) setCacheMap(key string, val interface{}) {
	c.cacheMapMu.Lock()
	c.cacheMap[key] = val
	c.cacheMapMu.Unlock()
}

func (c *GenerateCode) getCacheMapValue(key string) interface{} {
	c.cacheMapMu.Lock()
	r := c.cacheMap[key]
	c.cacheMapMu.Unlock()
	return r
}

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randString(letters string, n int, randomSrc rand.Source) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A randomSrc.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randomSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randomSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letters) {
			sb.WriteByte(letters[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}
