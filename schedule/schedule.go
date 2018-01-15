package schedule

import (
	"errors"
	"time"
)

// Schedule record the information of a schedule in parents day
type Schedule struct {
	ClassCode  string    `json:"classcode"`
	ClassNo    int       `json:"classno"`
	ArrivedAt  time.Time `json:"arrivedAt"`
	Order      int       `json:"order"`
	IsNotified bool      `json:"isNotified"`
	IsMeeting  bool      `json:"isMeeting"`
	IsComplete bool      `json:"isComplete"`
}

// MapSchedules is a slice of Schedule
type MapSchedules map[string][]*Schedule

// New create new MapSchedules
func New() MapSchedules {
	return make(map[string][]*Schedule)
}

// AppendSchedule append a new schedule to schedules
func (m MapSchedules) AppendSchedule(classCode string, classNo int) error {
	_, s := m.FindSchedule(classCode, classNo)
	if s != nil {
		return errors.New("schedule is already exist")
	}

	m[classCode] = append(m[classCode], &Schedule{
		ClassCode:  classCode,
		ClassNo:    classNo,
		ArrivedAt:  time.Now(),
		Order:      0,
		IsNotified: false,
		IsMeeting:  false,
		IsComplete: false,
	})
	return nil
}

// RemoveSchedule remove schedule
func (m MapSchedules) RemoveSchedule(classCode string, classNo int) error {
	i, _ := m.FindSchedule(classCode, classNo)
	if i == -1 {
		return errors.New("schedule not found")
	}

	m[classCode] = append(m[classCode][:i], m[classCode][i+1:]...)
	return nil
}

// UpdateOrder is update order
func (m MapSchedules) UpdateOrder(classCode string, classNo int, order int) error {
	i, s := m.FindSchedule(classCode, classNo)
	if i == -1 {
		return errors.New("schedule not found")
	}

	s.Order = order
	return nil
}

// ToggleIsNotified toggle IsNotified
func (m MapSchedules) ToggleIsNotified(classCode string, classNo int) error {
	i, s := m.FindSchedule(classCode, classNo)
	if i == -1 {
		return errors.New("schedule not found")
	}
	s.IsNotified = !s.IsNotified
	return nil
}

// ToggleIsMeeting toggle IsMeeting
func (m MapSchedules) ToggleIsMeeting(classCode string, classNo int) error {
	i, s := m.FindSchedule(classCode, classNo)
	if i == -1 {
		return errors.New("schedule not found")
	}
	s.IsMeeting = !s.IsMeeting
	return nil
}

// ToggleIsComplete toggle IsComplete
func (m MapSchedules) ToggleIsComplete(classCode string, classNo int) error {
	i, s := m.FindSchedule(classCode, classNo)
	if i == -1 {
		return errors.New("schedule not found")
	}
	s.IsComplete = !s.IsComplete
	return nil
}

// FindSchedule find schedule by given classCode and classNo
func (m MapSchedules) FindSchedule(classCode string, classNo int) (int, *Schedule) {
	for i, d := range m[classCode] {
		if d.ClassNo == classNo {
			return i, d
		}
	}
	return -1, nil
}
