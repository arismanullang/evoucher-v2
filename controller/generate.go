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
	numberRoutine int
	codeLength    int
	letterConsume string
	processFunc   func(string)
}

// NewGenerateCode returns a GenerateVoucher
// f for listener each process routines
// fixed configuration given length for fixed character code length
// routine number of generated code
// each of code will be using 1 routines
func NewGenerateCode(f func(str string), length, routine int) *GenerateCode {
	var defaultLetter = "abcdefghijklmnopqrstuvwxyzQWERTYUIOPASDFGHJKLZXCVBNM0123456789"
	return &GenerateCode{randomSrc: rand.NewSource(time.Now().UnixNano()),
		cacheMap:      make(map[string]interface{}),
		numberRoutine: routine,
		codeLength:    length,
		letterConsume: defaultLetter,
		processFunc:   f}
}

// make sure app not terminated before all routines complete
func (c *GenerateCode) start() *GenerateCode {
	c.cacheMapWg.Add(c.numberRoutine)
	for i := 0; i < c.numberRoutine; i++ {
		go c.process()
	}
	c.cacheMapWg.Wait()
	return c
}

// duplicate code will be generated 2 times max. may be changed next plan
func (c *GenerateCode) process() *GenerateCode {
	str := c.getUniqueString()
	c.processFunc(str)
	c.cacheMapWg.Done()
	return c
}

func (c *GenerateCode) getUniqueString() (str string) {
	for {
		str = string(randString(c.letterConsume, c.codeLength, c.randomSrc))
		if c.getCacheMapValue(str) != nil {
			// fmt.Println("WARNING!!DUPLICATE.")
		} else {
			c.setCacheMap(str, "")
			break
		}
	}
	return
}

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
