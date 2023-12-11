package postgresdb

import (
	"context"
	"errors"
	"gostorage3081/pkg/storage/interface"

	"github.com/jackc/pgx/v4/pgxpool"
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

func (s *Storage) EditTask(t storage.Task) (int, error) {
	if t.ID == 0 {
		return 0, errors.New("ID is not set")
	}
	//update not empty and not 0 fields
	res, err := s.db.Exec(context.Background(), `
		UPDATE tasks SET
		opened=COALESCE(NULLIF($2,0), opened),
		closed=COALESCE(NULLIF($3,0), closed),
		author_id=COALESCE(NULLIF($4,0), author_id),
		assigned_id=COALESCE(NULLIF($5,0), assigned_id),
		title=COALESCE(NULLIF($6,''), title),
		content=COALESCE(NULLIF($7,''), content)
		WHERE id=$1 
		`,
		t.ID,
		t.Opened,
		t.Closed,
		t.AuthorID,
		t.AssignedID,
		t.Title,
		t.Content,
	)
	return int(res.RowsAffected()), err
}

func (s *Storage) DeleteTask(taskID int) (int, error) {
	if taskID == 0 {
		return 0, errors.New("taskID is not set")
	}
	res, err := s.db.Exec(context.Background(), `
	DELETE FROM tasks
	WHERE id=$1
	`,
		taskID,
	)
	return int(res.RowsAffected()), err
}
