package rediskey

import "fmt"

func ConvMaxSeq(convId string) string {
	return fmt.Sprintf("im/convmaxseq/convid:%s", convId)
}

func ConvMinSeq(convId string) string {
	return fmt.Sprintf("im/convminseq/convid:%s", convId)
}
