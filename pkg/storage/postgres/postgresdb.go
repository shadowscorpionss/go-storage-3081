package postgresdb

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"gostorage3081/pkg/storage/interface"
)

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// Tasks возвращает список задач из БД.
func (s *Storage) Tasks(taskID, authorID int) ([]storage.Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			($1 = 0 OR id = $1) AND
			($2 = 0 OR author_id = $2)
		ORDER BY id;
	`,
		taskID,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []storage.Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t storage.Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// NewTask создаёт новую задачу и возвращает её id.
func (s *Storage) NewTask(t storage.Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (title, content)
		VALUES ($1, $2) RETURNING id;
		`,
		t.Title,
		t.Content,
	).Scan(&id)
	return id, err
}