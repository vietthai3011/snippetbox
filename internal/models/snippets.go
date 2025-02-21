package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (s *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := s.DB.Exec(stmt, title, content, expires)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT * FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := s.DB.QueryRow(stmt, id)

	var snippet Snippet

	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	return snippet, nil
}

// This will return the 10 most recently created snippets.
func (s *SnippetModel) Laster() ([]Snippet, error) {
	stmt := `SELECT * FROM snippets
				WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := s.DB.Query(stmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var snippet Snippet

		err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, snippet)
	}

	// Kiểm tra lỗi xảy ra trong quá trình lặp nhưng không bắt được trong Scan()
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
