package utils

import (
	"math/rand"
	"time"
)

type RandomUtil struct {
}

var Random = RandomUtil{}

// String 此方法生成n位随机字符串
// strList: 字符串列表
// length: 随机字符串长度
func (r RandomUtil) String(strList []string, length int) string {
	var str string
	for i := 0; i < length; i++ {
		// 从strList中随机获取一个字符串
		rand.Seed(time.Now().UnixNano())
		str += strList[rand.Intn(len(strList))]
	}
	return str
}
