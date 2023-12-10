package storage

// Задача.
type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

type DbInterface interface {
	Tasks (taskID, authorID int) ([]Task, error)
	NewTask(Task) (int, error)
}