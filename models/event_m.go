package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"index" json:"title" binding:"required"`
	EventDate string         `gorm:"type:date" json:"event_date" binding:"required"`
	StartTime string         `gorm:"type:timetz" json:"start_time" binding:"required"`
	EndTime   string         `gorm:"type:timetz" json:"end_time" binding:"required"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"-" `
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Event) TableName() string {
	return "events"
}

func (e *Event) GetStartTime() time.Time {
	t, err := time.Parse("15:04:05-07", e.StartTime)
	if err != nil {
		return time.Time{}
	}
	return t
}

func (e *Event) SetStartTime(t time.Time) {
	e.StartTime = t.Format("15:04:05-07")
}

func (e *Event) GetEndTime() time.Time {
	t, err := time.Parse("15:04:05-07", e.EndTime)
	if err != nil {
		return time.Time{}
	}
	return t
}

func (e *Event) SetEndTime(t time.Time) {
	e.EndTime = t.Format("15:04:05-07")
}

func (e *Event) GetEventDate() time.Time {
	t, err := time.Parse("2023-12-05", e.EventDate)
	if err != nil {
		return time.Time{}
	}
	return t
}

func (e *Event) SetEventDate(t time.Time) {
	e.EventDate = t.Format("2023-12-05")
}
