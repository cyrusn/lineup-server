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

// Schedules is a slice of Schedule
type Schedules map[string][]*Schedule

// New create new Schedules
func New() Schedules {
	return make(map[string][]*Schedule)
}

// AppendSchedule append a new schedule to schedules
func (s Schedules) AppendSchedule(classCode string, classNo int) error {
	_, d := s.FindSchedule(classCode, classNo)
	if d != nil {
		return errors.New("schedule is already exist")
	}

	s[classCode] = append(s[classCode], &Schedule{
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
func (s Schedules) RemoveSchedule(classCode string, classNo int) error {
	i, _ := s.FindSchedule(classCode, classNo)
	if i == -1 {
		return errors.New("schedule not found")
	}

	s[classCode] = append(s[classCode][:i], s[classCode][i+1:]...)
	return nil
}

// UpdateOrder is update order
func (s Schedules) UpdateOrder(classCode string, classNo int, order int) error {
	i, p := s.FindSchedule(classCode, classNo)
	if i == -1 {
		return errors.New("schedule not found")
	}

	p.Order = order
	return nil
}

// ToggleIsNotified toggle IsNotified
func (s Schedules) ToggleIsNotified(classCode string, classNo int) error {
	i, p := s.FindSchedule(classCode, classNo)
	if i == -1 {
		return errors.New("schedule not found")
	}
	p.IsNotified = !p.IsNotified
	return nil
}

// ToggleIsMeeting toggle IsMeeting
func (s Schedules) ToggleIsMeeting(classCode string, classNo int) error {
	i, p := s.FindSchedule(classCode, classNo)
	if i == -1 {
		return errors.New("schedule not found")
	}
	p.IsMeeting = !p.IsMeeting
	return nil
}

// ToggleIsComplete toggle IsComplete
func (s Schedules) ToggleIsComplete(classCode string, classNo int) error {
	i, p := s.FindSchedule(classCode, classNo)
	if i == -1 {
		return errors.New("schedule not found")
	}
	p.IsComplete = !p.IsComplete
	return nil
}

// FindSchedule find schedule by given classCode and classNo
func (s Schedules) FindSchedule(classCode string, classNo int) (int, *Schedule) {
	for i, d := range s[classCode] {
		if d.ClassNo == classNo {
			return i, d
		}
	}
	return -1, nil
}
