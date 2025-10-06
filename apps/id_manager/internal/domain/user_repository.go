package domain

type UserRepository interface {
	Save(user *User) error
	FindByUsername(username string) (*User, error)
	FindAll() ([]*User, error)
	Delete(username string) error
	UpdateConnectionStatus(username string, connected bool, ip string) error
	AddPendingMessage(username string, messageID string) error
	RemovePendingMessage(username string, messageID string) error
}
