package database

import (
	"context"
	"database/sql"
	"log"

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

func (r *PostgresRepository) SetQuestion(ctx context.Context, question *models.Question) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO questions (id, test_id, question, answer) VALUES ($1, $2, $3, $4)", question.Id, question.TestId, question.Question, question.Answer)
	return err
}

func (r *PostgresRepository) SetEnrollment(ctx context.Context, enrollment *models.Enrollment) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO enrollments (student_id, test_id) VALUES ($1, $2)", enrollment.StudentId, enrollment.TestId)
	return err
}

func (r *PostgresRepository) GetStudentsPerTest(ctx context.Context, testId string) ([]*models.Student, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, age FROM students WHERE id IN (SELECT student_id FROM enrollments WHERE test_id = $1)", testId)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	students := make([]*models.Student, 0)
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.Id, &student.Name, &student.Age); err != nil {
			return nil, err
		}
		students = append(students, &student)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func (r *PostgresRepository) GetQuestionsPerTest(ctx context.Context, testId string) ([]*models.Question, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, question FROM questions WHERE test_id = $1", testId)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	questions := make([]*models.Question, 0)
	for rows.Next() {
		var question models.Question
		if err := rows.Scan(&question.Id, &question.Question); err != nil {
			return nil, err
		}
		questions = append(questions, &question)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}
