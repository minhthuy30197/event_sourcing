package controller

import (
	"time"
	"errors"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/minhthuy30197/event_sourcing/model"
)

// @Tags admin
// @Description Lấy danh sách User
// @Success 200 {string} string
// @Failure 500 {object} model.HTTPError
// @Router /auth/get-users [get]
func (c *Controller) AddTeacherToClass(ctx *gin.Context) {
	var setTeacher model.SetTeacher
	err := ctx.ShouldBindJSON(&setTeacher)
	if err != nil {
		model.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	// Lấy version 
	var version int32
	_, err = c.DB.Query(&version, `SELECT version from course.class where course_id = ?`, setTeacher.CourseID)
	if err != nil {
		model.NewError(ctx, http.StatusBadRequest, errors.New("Không lay duoc thong tin giang vien."))
		return
	} 

	// Lay thong tin teacher can them
	var newTeacher model.TeacherInfo
	_, err = c.DB.Query(&newTeacher, `select * from course.teacher where id = ?`, setTeacher.TeacherID)
	if err != nil {
		model.NewError(ctx, http.StatusBadRequest, errors.New("Không lay duoc thong tin giang vien."))
		return
	}

	// Tao event
	var addTeacherEvent model.AddTeacherEvent
	addTeacherEvent.Teacher = newTeacher
	addTeacherEvent.CourseID = setTeacher.CourseID
	aggregateID := "ClassTeacher_" + setTeacher.CourseID
	baseEvent := BuildBaseEvent(aggregateID, "", "TeacherAdded", addTeacherEvent, (version + 1))
	err = c.SaveEvent(baseEvent)
	if err != nil {
		log.Println(err)
		model.NewError(ctx, http.StatusBadRequest, errors.New("Không them duoc event."))
		return
	}

	ctx.String(http.StatusOK, "Them thành công")
}

// @Tags admin
// @Description Tạo User
// @Param user body model.CreateUser true "Thông tin tạo User"
// @Success 200 {string} string
// @Failure 500 {object} model.HTTPError
// @Router /auth/create-user [post]
func (c *Controller) Playback(ctx *gin.Context) {
	var courseID = ctx.Param("id")
	var aggregateID = "ClassTeacher_" + courseID
	var startTime time.Time
	rs, err := c.Events(aggregateID, startTime, time.Now())
	if err != nil {
		log.Println(err)
		model.NewError(ctx, http.StatusBadRequest, errors.New("Không thể lấy lich su."))
		return
	}

	class := &model.ClassTeacherAggregate{}
	for _, event := range rs {
		class.Apply(event)
		class.Version ++
	}
	ctx.JSON(http.StatusOK, class)
}

// @Tags admin
// @Description Lấy thông tin User theo ID
// @Success 200 {string} string
// @Failure 404 {object} model.HTTPError
// @Failure 500 {object} model.HTTPError
// @Router /auth/get-user/{id} [get]
func (c *Controller) RemoveTeacherFromClass(ctx *gin.Context) {
	var setTeacher model.SetTeacher
	err := ctx.ShouldBindJSON(&setTeacher)
	if err != nil {
		model.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	// Lấy version 
	var version int32
	_, err = c.DB.Query(&version, `SELECT version from course.class where course_id = ?`, setTeacher.CourseID)
	if err != nil {
		model.NewError(ctx, http.StatusBadRequest, errors.New("Không lay duoc thong tin giang vien."))
		return
	} 

	// Lay thong tin teacher can them
	var newTeacher model.TeacherInfo
	_, err = c.DB.Query(&newTeacher, `select * from course.teacher where id = ?`, setTeacher.TeacherID)
	if err != nil {
		model.NewError(ctx, http.StatusBadRequest, errors.New("Không lay duoc thong tin giang vien."))
		return
	}

	// Tao event
	var removeTeacherEvent model.RemoveTeacherEvent
	removeTeacherEvent.Teacher = newTeacher
	removeTeacherEvent.CourseID = setTeacher.CourseID
	aggregateID := "ClassTeacher_" + setTeacher.CourseID
	baseEvent := BuildBaseEvent(aggregateID, "", "TeacherRemoved", removeTeacherEvent, (version + 1))
	err = c.SaveEvent(baseEvent)
	if err != nil {
		log.Println(err)
		model.NewError(ctx, http.StatusBadRequest, errors.New("Không them duoc event."))
		return
	}

	ctx.String(http.StatusOK, "Xoa Thành công")
}

// @Tags admin
// @Description Cập nhật User theo ID
// @Success 200 {string} string
// @Failure 404 {object} model.HTTPError
// @Failure 500 {object} model.HTTPError
// @Router /auth/update-user/{id} [put]
func (c *Controller) GetTeachersOfClass(ctx *gin.Context) {
	var courseID = ctx.Param("id")

	var response model.GetTeacher
	response.CourseID = courseID

	_, err := c.DB.Query(response.Teachers, `
		SELECT teacher.id as id, teacher.name as name
		FROM course.teacher as teacher
		INNER JOIN course.class as class
		ON teacher.id = ANY(class.teacher_ids)
		WHERE class.course_id = ?`, courseID)
	if err != nil {
		model.NewError(ctx, http.StatusBadRequest, errors.New("Không thể lấy danh sách giảng viên của khoá học."))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) GetHistory(ctx *gin.Context) {
	var courseID = ctx.Param("id")
	var getHistoryRequest model.GetHistoryRequest
	err := ctx.ShouldBindJSON(&getHistoryRequest)
	if err != nil {
		model.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	var aggregateID = "ClassTeacher_" + courseID
	startTime, _ := time.Parse(time.RFC3339Nano, getHistoryRequest.StartTime)
	endTime, err := time.Parse(time.RFC3339Nano, getHistoryRequest.EndTime)
	if err != nil {
		endTime = time.Now()
	}
	log.Println(endTime)
	log.Println(startTime)
	rs, err := c.Events(aggregateID, startTime, endTime)
	if err != nil {
		log.Println(err)
		model.NewError(ctx, http.StatusBadRequest, errors.New("Không thể lấy lich su."))
		return
	}

	ctx.JSON(http.StatusOK, rs)
}