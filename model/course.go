package model

type Teacher struct {
	// Tên bảng
	TableName struct{} `json:"table_name" sql:"course.teacher"`

	// Mã teacher (chuỗi ngẫu nhiên duy nhất)
	Id string `json:"id"`

	// Tên hiển thị
	Name string `json:"name"`
}

type Class struct {
	// Tên bảng
	TableName struct{} `json:"table_name" sql:"course.class"`

	// Mã course (chuỗi ngẫu nhiên duy nhất)
	CourseID string `json:"course_id"`

	// Tên hiển thị
	TeacherIDS []string `json:"teacher_ids" pg:",array"`

	// Version
	Version int32 `json:"version"`
}

type Course struct {
	// Tên bảng
	TableName struct{} `json:"table_name" sql:"course.course"`

	// Mã User (chuỗi ngẫu nhiên duy nhất)
	Id string `json:"id"`

	// Tên hiển thị
	Title string `json:"title"`
}

type SetTeacher struct {
	CourseID  string `json:"course_id"`
	TeacherID string `json:"teacher_id"`
}

type TeacherInfo struct {
	// Mã teacher (chuỗi ngẫu nhiên duy nhất)
	Id string `json:"id"`

	// Tên hiển thị
	Name string `json:"name"`
}

type GetTeacher struct {
	CourseID string        `json:"course_id"`
	Teachers []TeacherInfo `json:"teachers"`
}

type GetHistoryRequest struct {
	StartTime string `json:"start_time"`
	EndTime string `json:"end_time"`
}
