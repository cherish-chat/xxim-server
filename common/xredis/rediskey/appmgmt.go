package rediskey

func AllAppMgmtConfig() string {
	return "s:appmgmt:config"
}

func AppAddressBook() string {
	return "h:appmgmt:address_book"
}

func AppAddressBookUrl() string {
	return "s:appmgmt:address_book_url"
}
