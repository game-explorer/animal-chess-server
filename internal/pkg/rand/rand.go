package util

import (
	"math/rand"
	"sync/atomic"
	"time"
)

// RandomString 显示随机字符串,默认8位,如需自定义长度请传入int型的长度单位,1-32位之间
func RandomString(length ...int) string {
	var (
		seedData = []byte(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890`)
		res      []byte
	)
	if len(length) == 0 {
		length = make([]int, 1)
		length[0] = 8
	}
	if length[0] < 1 || length[0] > 32 {
		length[0] = 8
	}
	res = make([]byte, 0, length[0])

	for i := 0; i < length[0]; i++ {
		res = append(res, seedData[RandomInt64n(62)])
	}

	return string(res)
}

// 生成随机Int, 在min和max之间, 有可能为min和max
func RandomInt(min, max int) int {
	return min + int(RandomInt64n(int64(max-min+1)))
}

var added int64 = 0

func RandomInt64n(n int64) int64 {
	v := atomic.AddInt64(&added, 1)
	rand.Seed(time.Now().UnixNano() + v)
	return rand.Int63n(n)
}

// 生成随机数字串 主要用于4位数生成短信验证码
func RandomInt4(n int) string {
	const alphanum = "0123456789"
	var bytes = make([]byte, n)
	v := atomic.AddInt64(&added, 1)
	rand.Seed(time.Now().UnixNano() + v)
	for i, _ := range bytes {
		bytes[i] = alphanum[rand.Intn(len(alphanum))] // 返回一个取值范围在[0,n)的伪随机int值
	}
	return string(bytes)
}
