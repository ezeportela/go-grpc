package database

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/ezeportela/go-grpc/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) *PostgresRepository {
	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}

	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) GetStudent(ctx context.Context, id string) (*models.Student, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name, age FROM students WHERE id = $1", id)

	var student models.Student
	if err := row.Scan(&student.Id, &student.Name, &student.Age); err != nil {
		return nil, err
	}

	return &student, nil
}

func (r *PostgresRepository) SetStudent(ctx context.Context, student *models.Student) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO students (id, name, age) VALUES ($1, $2, $3)", student.Id, student.Name, student.Age)
	return err
}

func (r *PostgresRepository) GetTest(ctx context.Context, id string) (*models.Test, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name FROM tests WHERE id = $1", id)

	var test models.Test
	if err := row.Scan(&test.Id, &test.Name); err != nil {
		return nil, err
	}

	return &test, nil
}

func (r *PostgresRepository) SetTest(ctx context.Context, test *models.Test) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO tests (id, name) VALUES ($1, $2)", test.Id, test.Name)
	return err
}
