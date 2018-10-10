package model

import (
	"fmt"
	"log"
	"time"

	"github.com/minhthuy30197/event_sourcing/config"
	"github.com/minhthuy30197/event_sourcing/helper"
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
	AggregateId string      `json:"aggregate_id"`
	Time        time.Time   `json:"time" sql:"default:now()"`
	Version     int32       `json:"version"`
	Data        interface{} `json:"data"`
	EventType   string      `json:"event_type"`
	UserID      string      `json:"user_id"`
	Revision    int32       `json:"revision"`
}

func (event Event) Apply(c config.Config) error {
	err := event.Data.(EventInterface).Apply(event, c)
	if err != nil {
		return err
	}
	return nil
}

type EventInterface interface {
	Apply(Event, config.Config) error
}

type AddTeacherEvent struct {
	CourseID string      `json:"course_id"`
	Teacher  TeacherInfo `json:"teacher"`
}

func (event AddTeacherEvent) Apply(ev Event, config config.Config) error {
	dbConfig := config.Database
	db := ConnectDb(dbConfig.User, dbConfig.Password, dbConfig.Database, dbConfig.Address)
	defer db.Close()

	// Update read databse
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

type RemoveTeacherEvent struct {
	CourseID string      `json:"course_id"`
	Teacher  TeacherInfo `json:"teacher"`
}

func (event RemoveTeacherEvent) Apply(ev Event, config config.Config) error {
	dbConfig := config.Database
	db := ConnectDb(dbConfig.User, dbConfig.Password, dbConfig.Database, dbConfig.Address)
	defer db.Close()

	// Update read databse
	_, err := db.Exec(`UPDATE course.class SET teacher_ids = array_remove(teacher_ids, ?) WHERE course_id = ?`,
		event.Teacher.Id, event.CourseID)
	if err != nil {
		return err
	}

	return nil
}

func (class *Class) Transition(event interface{}) {
	switch e := event.(type) {
	case RemoveTeacherEvent:
		log.Println("RemoveTeacherEvent")
		pos := helper.GetPosStringElementInSlice(class.TeacherIDS, e.Teacher.Id)
		if pos != -1 {
			class.CourseID = e.CourseID
			copy(class.TeacherIDS[pos:], class.TeacherIDS[pos+1:])
			class.TeacherIDS[len(class.TeacherIDS)-1] = ""
			class.TeacherIDS = class.TeacherIDS[:len(class.TeacherIDS)-1]
		}
	case AddTeacherEvent:
		log.Println("AddTeacherEvent")
		class.CourseID = e.CourseID
		class.TeacherIDS = append(class.TeacherIDS, e.Teacher.Id)
	default:
		fmt.Printf("%T\n", e)
	}
}
