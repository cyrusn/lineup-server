package model

import (
	"database/sql"
	"fmt"
	"time"
)

// Schedule record the information of a schedule in parents day
type Schedule struct {
	ClassCode  string    `json:"classcode"`
	ClassNo    int       `json:"classno"`
	ArrivedAt  time.Time `json:"arrivedAt"`
	Priority   int       `json:"priority"`
	IsNotified bool      `json:"isNotified"`
	IsMeeting  bool      `json:"isMeeting"`
	IsComplete bool      `json:"isComplete"`
}

// ScheduleDB is *sql.DB
type ScheduleDB struct {
	*sql.DB
}

// SelectByClassCode find all schedules by classcode
func (db ScheduleDB) SelectByClassCode(classcode string) ([]*Schedule, error) {
	var schedules []*Schedule
	rows, err := db.Query(`SELECT * FROM SCHEDULE WHERE classcode = ?`, classcode)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		s := new(Schedule)
		if err := rows.Scan(
			&s.ClassCode,
			&s.ClassNo,
			&s.ArrivedAt,
			&s.Priority,
			&s.IsNotified,
			&s.IsMeeting,
			&s.IsComplete,
		); err != nil {
			return nil, err
		}

		schedules = append(schedules, s)
	}
	return schedules, nil
}

// Insert Schedule by given classCode and classNo
func (db ScheduleDB) Insert(classCode string, classNo int) error {
	_, err := db.Exec(`INSERT INTO SCHEDULE (
      classcode,
      classno,
      arrived_at,
      priority,
      is_notified,
      is_meeting,
      is_complete
    ) values (?, ?, ?, ?, ?, ?, ?)`,
		classCode, classNo, time.Now(), 0, false, false, false,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteSchedule delete schedule
func (db ScheduleDB) DeleteSchedule(classcode string, classno int) error {

	if _, err := db.Exec(
		`DELETE FROM SCHEDULE WHERE (
      classcode = ? and classno = ?
    )`,
		classcode, classno,
	); err != nil {
		return err
	}

	return nil
}

// UpdatePriority update schedule's priority
func (db ScheduleDB) UpdatePriority(classcode string, classno int, priority int) error {
	if _, err := db.Exec(
		`UPDATE SCHEDULE SET priority = ? WHERE (
      classcode = ? and classno = ?
    )`,
		priority, classcode, classno,
	); err != nil {
		return err
	}

	return nil
}

func (db ScheduleDB) toggleFactory(key string) func(string, int) error {
	return func(classCode string, classNo int) error {
		var boolVar bool
		query := fmt.Sprintf(`SELECT %s FROM SCHEDULE WHERE (
      classcode = ? and classno = ?
    )`, key)
		err := db.QueryRow(query,
			classCode, classNo).Scan(&boolVar)

		if err != nil {
			return err
		}

		exec := fmt.Sprintf(`
    UPDATE SCHEDULE SET %s = ? WHERE (
      classcode = ? and classno = ?
    )`, key)

		if _, err := db.Exec(exec,
			!boolVar, classCode, classNo,
		); err != nil {
			return err
		}
		return nil

	}
}

// ToggleIsNotified toggle IsNotified
func (db ScheduleDB) ToggleIsNotified(classCode string, classNo int) error {
	return db.toggleFactory("is_notified")(classCode, classNo)
}

// ToggleIsMeeting toggle IsMeeting
func (db ScheduleDB) ToggleIsMeeting(classCode string, classNo int) error {
	return db.toggleFactory("is_meeting")(classCode, classNo)
}

// ToggleIsComplete toggle IsComplete
func (db ScheduleDB) ToggleIsComplete(classCode string, classNo int) error {
	return db.toggleFactory("is_complete")(classCode, classNo)
}
