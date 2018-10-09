package model

import (
	"time"
	"git.hocngay.com/hocngay/event-sourcing/config"
)

type EventSource struct {
	TableName   struct{}  `json:"table_name" sql:"es.event_source"`
	AggregateId string    `json:"aggregate_id"`
	Time        time.Time `json:"time" sql:"default:now()"`
	Version     int32     `json:"version"`
	Data        []string  `json:"data"`
	EventType   string    `json:"event_type"`
	UserID      string    `json:"user_id"`
	Revision    int32     `json:"revision"`
} 

type Event struct {
	AggregateId string    `json:"aggregate_id"`
	Time        time.Time `json:"time" sql:"default:now()"`
	Version     int32     `json:"version"`
	Data        interface{}  `json:"data"`
	EventType   string    `json:"event_type"`
	UserID      string    `json:"user_id"`
	Revision    int32     `json:"revision"`
}

type EventInterface interface {
	Apply(Event, config.Config) error
	AggregateType() string
}

type AddTeacherEvent struct {
	CourseID string 
	Teacher TeacherInfo
}

func (event AddTeacherEvent) AggregateType() string {
	return "AddTeacherEvent"
}

func (event AddTeacherEvent) Apply(ev Event, config config.Config) error {
	dbConfig := config.Database
	db := ConnectDb(dbConfig.User, dbConfig.Password, dbConfig.Database, dbConfig.Address)
	defer db.Close()

	// Update vào read databse
	var teacherOfClass Class
	err := db.Model(&teacherOfClass).Where(`course_id = ?`, event.CourseID).Select()
	if err != nil {
		teacherOfClass.CourseID = event.CourseID
		teacherOfClass.TeacherIDS = []string{event.Teacher.Id}
		err = db.Insert(&teacherOfClass)
		if err != nil {
			return err
		}
	} else {
		_, err = db.Exec(`UPDATE course.class SET teacher_ids = array_append(teacher_ids, ?) WHERE course_id = ?`,
			event.Teacher.Id, event.CourseID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (event Event) Apply(c config.Config) error {
	err := event.Data.(EventInterface).Apply(event, c)
	if err != nil {
		return err
	}
	return nil
}

