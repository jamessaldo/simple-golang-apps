package commands

type CreatePost struct {
	ID          uint64
	UserID      uint64
	Title       string
	Description string
}
