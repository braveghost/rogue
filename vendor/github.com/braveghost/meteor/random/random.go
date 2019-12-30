package random

import (
	logging "github.com/braveghost/joker"
	"github.com/braveghost/meteor/icrypto"
	"github.com/braveghost/meteor/pool"
	"github.com/zheng-ji/goSnowFlake"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var (
	idWorker *goSnowFlake.IdWorker
)

func init() {
	var err error
	if idWorker, err = goSnowFlake.NewIdWorker(1); err != nil {
		logging.Errorw("SnowFlakeWorker.Init.Error", "err", err)
		os.Exit(1)
	}
}

type snowFlakeValue int64

func (s snowFlakeValue) String() string {
	return strconv.FormatInt(int64(s), 10)
}

func SnowFlake() snowFlakeValue {
	id, _ := idWorker.NextId()
	return snowFlakeValue(id)
}

func SnowFlakeMd5() string {
	id, _ := idWorker.NextId()
	return icrypto.StrToMd5(strconv.FormatInt(id, 10))
}

var (
	randSource = rand.New(rand.NewSource(time.Now().UnixNano()))
	source     = []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
)

var randBufferDict = make(map[int]*pool.BufferPool)

type codeChars struct {
	Chars string
	Len   int
}

func NewCodeChars(str string) *codeChars {
	return &codeChars{str, len(str)}
}

var (
	CharsAllLetterNumbers = NewCodeChars(charsLetterNumbers)
	CharsUpperLetter      = NewCodeChars(charsUpperLetter)
	CharsLowerLetter      = NewCodeChars(charsLowerLetter)
	CharsNumbers          = NewCodeChars(charsNumbers)
)

const (
	idxBits            = 6
	charsLetterNumbers = "0123456789AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	charsUpperLetter   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsLowerLetter   = "abcdefghijklmnopqrstuvwxyz"
	charsNumbers       = "1234567890"
	mask               = 1<<idxBits - 1
)

var src = rand.NewSource(time.Now().UnixNano())

// 生成随机字符串
func GenRandStr(ln int, model *codeChars) string {
	dataPool, ok := randBufferDict[ln]
	if ! ok {
		dataPool = pool.NewBytesPool(ln, ln)
		randBufferDict[ln] = dataPool
	}
	tmp := dataPool.Get()
	defer tmp.Free()
	buf := tmp.Buf.Bytes()

	chars := model.Chars
	charsLen := model.Len
	for idx, cache, remain := ln-1, src.Int63(), 10; idx >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), 10
		}
		buf[idx] = chars[int(cache&mask)%charsLen]
		cache >>= idxBits
		remain--
		idx--
	}
	return string(buf)
}

// 数字验证码
func NumCaptcha(ln int) string {
	return GenRandStr(ln, CharsNumbers)
}
