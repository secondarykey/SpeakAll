package ws

type message struct {
	Type      string
	MessageId int
	UserId    int
	UserName  string
	Project   string
	ClientId  string
	Content   string
	Category  string
	Created   string
}

// Type
//
// Open = "Open client id"
// Move = "client move category"

// Category
//
// Public   = "unused login"
// Dashbord = "top page message"
// uuid     = "create category message"

func createOpenMessage(clientId string) *message {
	return &message{
		Type:     "Open",
		ClientId: clientId,
		Content:  clientId,
	}
}

func createBadgeMessage(categoryId string) *message {
	return &message{
		Type:     "Notify",
		Category: categoryId,
	}
}