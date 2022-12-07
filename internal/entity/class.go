package entity

import "time"

type Class struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TeacherID   int64     `json:"teacher_id"`
	CreatedAt   time.Time `json:"created_at"`
	Version     int64     `json:"-"`
}
