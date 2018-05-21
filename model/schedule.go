package model

import (
	"database/sql"
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

// Insert ...
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

// // RemoveSchedule ...
// func RemoveSchedule(classcode string, classno int) error {

// }

// // UpdatePriority ...
// func UpdatePriority(classcode string, classno int) error {}

// // ToggleIsNotified ...
// func ToggleIsNotified(classcode string, classno int) error {}

// // ToggleIsMeeting ...
// func ToggleIsMeeting(classcode string, classno int) error {}

// // ToggleIsComplete ...
// func ToggleIsComplete(classcode string, classno int) error {}
