package models

import (
	"database/sql"
	"fmt"
	"time"
)

type TimeEntry struct {
	ID        int
	UserId    int
	Date      Date
	Time      string
	Minutes   int
	Details   sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
	ProjectId int
	Project   *Project
}

func (t *TimeEntry) Create() error {
	transaction, err := DB.Begin()

	if err != nil {
		return err
	}

	_, err = transaction.Exec("INSERT INTO time_entries (user_id, date, time, minutes, details, project_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", t.UserId, t.Date._time, t.Time, t.Minutes, t.Details, t.ProjectId, t.CreatedAt, t.UpdatedAt)

	if err != nil {
		return err
	}

	err = transaction.Commit()

	if err != nil {
		return err
	}

	return nil
}

func (t *TimeEntry) Description() string {
	// "*#{id}: #{project.name}* - #{time} - #{details}"

	return fmt.Sprintf("*%d: %s* - %s - %s", t.ID, t.Project.Name, t.Time, t.Details.String)
}

func GetTimeEntriesInPeriodWithProjectAndUser(user *User, project *Project, from Date, to Date) ([]*TimeEntry, error) {
	selectPart := "time_entries.id, time_entries.user_id, time_entries.date, time_entries.time, time_entries.minutes, time_entries.details, time_entries.created_at, time_entries.updated_at, time_entries.project_id"
	projectSelectPart := "projects.id, projects.name, projects.alias, projects.created_at, projects.updated_at"
	joins := "INNER JOIN projects ON projects.id = time_entries.project_id"
	sqlQuery := fmt.Sprintf("SELECT %s, %s FROM time_entries %s WHERE user_id = $1 AND date >= %s and date <= %s", selectPart, projectSelectPart, joins, from.SQL(), to.SQL())

	var rows *sql.Rows
	var err error

	if project == nil {
		rows, err = DB.Query(sqlQuery, user.ID)
	} else {
		rows, err = DB.Query(sqlQuery+" AND project_id = $2", user.ID, project.ID)
	}

	if err != nil {
		return nil, err
	}

	timeEntries := []*TimeEntry{}

	for rows.Next() {
		timeEntry := TimeEntry{}
		timeEntry.Project = &Project{}

		var d time.Time

		err = rows.Scan(&timeEntry.ID, &timeEntry.UserId, &d, &timeEntry.Time, &timeEntry.Minutes,
			&timeEntry.Details, &timeEntry.CreatedAt, &timeEntry.UpdatedAt, &timeEntry.ProjectId,
			&timeEntry.Project.ID, &timeEntry.Project.Name, &timeEntry.Project.Alias, &timeEntry.Project.CreatedAt,
			&timeEntry.Project.UpdatedAt)

		if err != nil {
			return nil, err
		}

		timeEntry.Date = NewDate(d)

		timeEntries = append(timeEntries, &timeEntry)
	}

	return timeEntries, nil
}

func (t TimeEntry) String() string {
	return fmt.Sprintf("User ID: %d; Time: %s; Details: \"%s\"", t.ID, t.Time, t.Details.String)
}
