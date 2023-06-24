package pb

import (
	"fmt"
	"strings"
)

const (
	SingleChatConversationSeparator = "&"
)

func GetSingleChatConversationId(id1 string, id2 string) string {
	if id1 < id2 {
		return fmt.Sprintf("%s%s%s", id1, SingleChatConversationSeparator, id2)
	} else {
		return fmt.Sprintf("%s%s%s", id2, SingleChatConversationSeparator, id1)
	}
}

func ParseSingleChatConversationId(convId string) (id1 string, id2 string) {
	split := strings.Split(convId, SingleChatConversationSeparator)
	if len(split) != 2 {
		return
	}
	id1 = split[0]
	id2 = split[1]
	return
}

func GetSingleChatOtherId(convId string, id1 string) (id2 string) {
	split := strings.Split(convId, SingleChatConversationSeparator)
	if len(split) != 2 {
		return
	}
	if split[0] == id1 {
		id2 = split[1]
	} else {
		id2 = split[0]
	}
	return
}
