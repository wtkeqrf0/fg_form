package appeal

type Appeal struct {
	Division string
	Subject  string
	Text     string
	ChatID   int64
	Username string
}

type NewAppeal struct {
	AppealID    int64
	AdminChatID int64
}
