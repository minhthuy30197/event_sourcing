package model

type Aggregate interface {
	UpdateVersion()
}

type BaseAggregate struct {
	Version int32 `json:"version"`
}

func (agg *BaseAggregate) UpdateVersion() {
	agg.Version++
}

type ClassTeacher struct {
	// Mã course (chuỗi ngẫu nhiên duy nhất)
	CourseID string `json:"course_id"`

	// Tên hiển thị
	TeacherIDS []string `json:"teacher_ids" pg:",array"`

	// Base aggregate
	BaseAggregate `json:"base_aggregate"`
}
