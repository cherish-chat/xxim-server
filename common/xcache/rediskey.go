package xcache

import "fmt"

const (
	IncrKeyGroupId = "group:i:group_id"
)

func UserUsernameLockKey(username string) string {
	return fmt.Sprintf("user:l:username:%s", username)
}

func UserPhoneLockKey(phone string, phoneCode string) string {
	return fmt.Sprintf("user:l:phone:%s:%s", phone, phoneCode)
}

func UserEmailLockKey(email string) string {
	return fmt.Sprintf("user:l:email:%s:%s", email)
}
