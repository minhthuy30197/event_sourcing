package model

import (
	"github.com/minhthuy30197/event_sourcing/helper"
)

type AggregateMethod interface {
	Apply(Event)  
}

type Aggregate struct {
	Item interface{}
}

func (aggregate *Aggregate) Apply(event Event) {
	aggregate.Item.(AggregateMethod).Apply(event)
} 

type ClassTeacherAggregate struct {
	// Mã course (chuỗi ngẫu nhiên duy nhất)
	CourseID string `json:"course_id"`

	// Tên hiển thị
	TeacherIDS []string `json:"teacher_ids" pg:",array"`

	// Version
	Version int32 `json:"version"`
}

func (class *ClassTeacherAggregate) Apply(event Event) {
	switch event.EventType {
	case "TeacherRemoved":
		pos := helper.GetPosStringElementInSlice(class.TeacherIDS, event.Data.(RemoveTeacherEvent).Teacher.Id)
		if pos != -1 {
			class.CourseID = event.Data.(RemoveTeacherEvent).CourseID
			copy(class.TeacherIDS[pos:], class.TeacherIDS[pos+1:])
			class.TeacherIDS[len(class.TeacherIDS)-1] = ""
			class.TeacherIDS = class.TeacherIDS[:len(class.TeacherIDS)-1]
		}
	case "TeacherAdded":
		class.CourseID = event.Data.(AddTeacherEvent).CourseID
		class.TeacherIDS = append(class.TeacherIDS, event.Data.(AddTeacherEvent).Teacher.Id)
	}
}
