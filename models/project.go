package models

import (
	"errors"
	"strings"
	"time"
)

type Project struct {
	ID        int
	Name      string
	Alias     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FindProjectByNameOrAlias(name string) (*Project, error) {
	selectPart := "id, name, alias, created_at, updated_at"
	downcasedName := strings.ToLower(name)

	rows, err := DB.Query("SELECT "+selectPart+" FROM projects WHERE lower(name) = ? OR lower(alias) = ?", downcasedName, downcasedName)

	if err != nil {
		return nil, err
	}

	project := Project{}

	if rows.Next() {
		err = rows.Scan(&project.ID, &project.Name, &project.Alias, &project.CreatedAt, &project.UpdatedAt)

		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("The project was not found")
	}

	return &project, nil
}
