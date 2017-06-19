package models

import (
	"time"
	"fmt"
	"database/sql"
)

type TimeEntry struct {
	ID        int
	UserId    int
	Date      time.Time
	Time      string
	Minutes   int
	Details   sql.NullString
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

func GetTimeEntriesInPeriodWithProjectAndUser(user *User, project *Project, from time.Time, to time.Time) ([]*TimeEntry, error) {
	selectPart := "id, user_id, date, time, minutes, details, created_at, updated_at, project_id"
	sqlQuery := fmt.Sprintf("SELECT %s FROM time_entries WHERE user_id = $1 AND date >= $2 and date <= $3", selectPart)

	var rows *sql.Rows
	var err error

	if project == nil {
		rows, err = DB.Query(sqlQuery, user.ID, formatDate(from), formatDate(to))
	} else {
		rows, err = DB.Query(sqlQuery + " AND project_id = $4", user.ID, formatDate(from), formatDate(to), project.ID)
	}

	if err != nil {
		return nil, err
	}

	timeEntries := []*TimeEntry{}

	for rows.Next() {
		timeEntry := TimeEntry{}
		err = rows.Scan(&timeEntry.ID, &timeEntry.UserId, &timeEntry.Date, &timeEntry.Time, &timeEntry.Minutes, &timeEntry.Details, &timeEntry.CreatedAt, &timeEntry.UpdatedAt, &timeEntry.ProjectId)

		if err != nil {
			return nil, err
		}

		timeEntries = append(timeEntries, &timeEntry)
	}

	return timeEntries, nil
}

func formatDate(datetime time.Time) string {
	return datetime.Format("2006-01-02")
}
