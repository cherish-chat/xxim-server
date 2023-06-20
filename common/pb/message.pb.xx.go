package pb

import "fmt"

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
