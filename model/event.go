package model

import (
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

type Snapshot struct {
	TableName   struct{}  `json:"table_name" sql:"es.snapshot"`
	AggregateId string    `json:"aggregate_id"`
	Time        time.Time `json:"time"`
	Version     int32     `json:"version"`
	Data        []string  `json:"data"`
}

type Event struct {
	AggregateId string      `json:"aggregate_id"`
	Time        time.Time   `json:"time" sql:"default:now()"`
	Version     int32       `json:"version"`
	Data        EventInterface `json:"data"`
	EventType   string      `json:"event_type"`
	UserID      string      `json:"user_id"`
	Revision    int32       `json:"revision"`
}

func (event Event) SaveReadDB(c config.Config) error {
	err := event.Data.(EventInterface).SaveReadDB(event, c)
	if err != nil {
		return err
	}
	return nil
}

func (event Event) Apply(agg Aggregate) error {
	err := event.Data.(EventInterface).Apply(agg, event)
	if err != nil {
		return err
	}

	agg.UpdateVersion()

	return nil
}

type EventInterface interface {
	SaveReadDB(Event, config.Config) error
	Apply(Aggregate, Event) error
}

type AddTeacherEvent struct {
	CourseID string      `json:"course_id"`
	Teacher  TeacherInfo `json:"teacher"`
}

func (event AddTeacherEvent) SaveReadDB(ev Event, config config.Config) error {
	dbConfig := config.Database
	db := ConnectDb(dbConfig.User, dbConfig.Password, dbConfig.Database, dbConfig.Address)
	defer db.Close()

	// Update read databse
	var teacherOfClass Class
	err := db.Model(&teacherOfClass).Where(`course_id = ?`, event.CourseID).Select()
	if err != nil {
		teacherOfClass.CourseID = event.CourseID
		teacherOfClass.Version = ev.Version
		teacherOfClass.TeacherIDS = []string{event.Teacher.Id}
		err = db.Insert(&teacherOfClass)
		if err != nil {
			return err
		}
	} else {
		_, err = db.Exec(`UPDATE course.class SET teacher_ids = array_append(teacher_ids, ?), version = ? WHERE course_id = ?`,
			event.Teacher.Id, ev.Version, event.CourseID)
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

func (event RemoveTeacherEvent) SaveReadDB(ev Event, config config.Config) error {
	dbConfig := config.Database
	db := ConnectDb(dbConfig.User, dbConfig.Password, dbConfig.Database, dbConfig.Address)
	defer db.Close()

	// Update read databse
	_, err := db.Exec(`UPDATE course.class SET teacher_ids = array_remove(teacher_ids, ?), version = ? WHERE course_id = ?`,
		event.Teacher.Id, ev.Version, event.CourseID)
	if err != nil {
		return err
	}

	return nil
}

func (event RemoveTeacherEvent) Apply(agg Aggregate, ev Event) error {
	class := agg.(*ClassTeacher)
	pos := helper.GetPosStringElementInSlice(class.TeacherIDS, event.Teacher.Id)
	if pos != -1 {
		class.CourseID = event.CourseID
		copy(class.TeacherIDS[pos:], class.TeacherIDS[pos+1:])
		class.TeacherIDS[len(class.TeacherIDS)-1] = ""
		class.TeacherIDS = class.TeacherIDS[:len(class.TeacherIDS)-1]
	}
	
	return nil
}

func (event AddTeacherEvent) Apply(agg Aggregate, ev Event) error {
	class := agg.(*ClassTeacher)
	class.CourseID = event.CourseID
	class.TeacherIDS = append(class.TeacherIDS, event.Teacher.Id)

	return nil 
}
