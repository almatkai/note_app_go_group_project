package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *pgxpool.Pool
}

func (m *SnippetModel) Insert(ctx context.Context, title string, content string, expires int, user_id int) (int, error) {
	var id int

	query := `INSERT INTO snippets (title, content, created, expires, user_id)
	VALUES($1, $2, NOW(), NOW() + INTERVAL '1 DAY' * $3, $4)
	RETURNING id`

	err := m.DB.QueryRow(ctx, query, title, content, expires, user_id).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *SnippetModel) Get(ctx context.Context, id, user_id int) ([]*Snippet, error) {
	if user_id == 0 {
		fmt.Println("some error has occurred")
	}
	query := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > NOW() and id <> $1 and user_id = $2
	ORDER BY id DESC`
	rows, err := m.DB.Query(ctx, query, id, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*Snippet{}
	////////////////////////////////////////////////////////////////
	queryId := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > NOW() and id = $1 and user_id = $2`
	row := m.DB.QueryRow(ctx, queryId, id, user_id)

	sID := &Snippet{}
	err = row.Scan(&sID.ID, &sID.Title, &sID.Content, &sID.Created, &sID.Expires)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	snippets = append(snippets, sID)
	////////////////////////////////////////////////////////////////
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil

}
func (m *SnippetModel) Latest(ctx context.Context, user_id int) ([]*Snippet, error) {
	if user_id == 0 {
		fmt.Println("some error has occurred")
	}
	query := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > NOW() and user_id = $1
	ORDER BY id DESC`
	rows, err := m.DB.Query(ctx, query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
