package models

import "time"

type TimeEntry struct {
	ID        int
	UserId    int
	Date      time.Time
	Time      string
	Minutes   int
	Details   string
	CreatedAt time.Time
	UpdatedAt time.Time
	ProjectId int
}

func (t *TimeEntry) Create() error {
	transaction, err := DB.Begin()

	if err != nil {
		return err
	}

	_, err = transaction.Exec("INSERT INTO time_entries (user_id, date, time, minutes, details, project_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", t.UserId, t.Date, t.Time, t.Minutes, t.Details, t.ProjectId, t.CreatedAt, t.UpdatedAt)

	if err != nil {
		return err
	}

	err = transaction.Commit()

	if err != nil {
		return err
	}

	return nil
}
