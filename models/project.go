package models

import (
	"database/sql"
	"strings"
	"time"
)

type Project struct {
	ID        int
	Name      string
	Alias     sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FindProjectByNameOrAlias(name string) (*Project, error) {
	selectPart := "id, name, alias, created_at, updated_at"
	downcasedName := strings.ToLower(name)

	rows, err := DB.Query("SELECT "+selectPart+" FROM projects WHERE lower(name) = $1 OR lower(alias) = $2", downcasedName, downcasedName)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	project := Project{}

	if rows.Next() {
		err = rows.Scan(&project.ID, &project.Name, &project.Alias, &project.CreatedAt, &project.UpdatedAt)

		if err != nil {
			return nil, err
		}
	} else {
		return nil, NotFoundError{}
	}

	return &project, nil
}

func GetAllProjects() ([]*Project, error) {
	selectPart := "id, name, alias, created_at, updated_at"

	rows, err := DB.Query("SELECT "+selectPart+" FROM projects ORDER BY name")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	out := []*Project{}

	for rows.Next() {
		p := &Project{}

		err = rows.Scan(&p.ID, &p.Name, &p.Alias, &p.CreatedAt, &p.UpdatedAt)

		if err != nil {
			return nil, err
		}

		out = append(out, p)
	}

	return out, nil
}

func (p *Project) Create() error {
	transaction, err := DB.Begin()

	if err != nil {
		return err
	}

	now := time.Now()

	p.CreatedAt = now
	p.UpdatedAt = now

	_, err = transaction.Exec("INSERT INTO projects (name, created_at, updated_at) VALUES ($1, $2, $3)", p.Name, p.CreatedAt, p.UpdatedAt)

	if err != nil {
		return err
	}

	err = transaction.Commit()

	if err != nil {
		return err
	}

	return nil
}
