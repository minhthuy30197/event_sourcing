package controller

import (
	"errors"
	"log"
	"net/http"

	"git.hocngay.com/hocngay/event-sourcing/model"
	"github.com/gin-gonic/gin"
)

// @Tags admin
// @Description Lấy danh sách User
// @Success 200 {string} string
// @Failure 500 {object} model.HTTPError
// @Router /auth/get-users [get]
func (c *Controller) AddTeacherToClass(ctx *gin.Context) {
	log.Println("------------------------")
	var setTeacher model.SetTeacher
	err := ctx.ShouldBindJSON(&setTeacher)
	if err != nil {
		model.NewError(ctx, http.StatusBadRequest, err)
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
	aggregateID := "ClassTeacher_" + setTeacher.CourseID
	baseEvent := BuildBaseEvent(aggregateID, "", "TeacherAdded", addTeacherEvent, 1)
	err = c.SaveEvent(baseEvent)
	if err != nil {
		log.Println(err)
		model.NewError(ctx, http.StatusBadRequest, errors.New("Không them duoc event."))
		return 
	}

	ctx.String(http.StatusOK, "Thành công")
}

// @Tags admin
// @Description Tạo User
// @Param user body model.CreateUser true "Thông tin tạo User"
// @Success 200 {string} string
// @Failure 500 {object} model.HTTPError
// @Router /auth/create-user [post]
func (c *Controller) Playback(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Playback")
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

	// TODO: Thêm sự kiện vào event_sourcing

	// Update vào read databse

	_, err = c.DB.Exec(`UPDATE course.class SET teacher_ids = array_remove((teacher_ids, ?) WHERE course_id = ?`,
		setTeacher.TeacherID, setTeacher.CourseID)
	if err != nil {
		model.NewError(ctx, http.StatusBadRequest, errors.New("Không thể rút giảng viên khỏi khoá học."))
		return
	}

	ctx.String(http.StatusOK, "Thành công")
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
	log.Println("Here")
	var courseID = ctx.Param("id")
	var aggregateID = "ClassTeacher_" + courseID
	
	rs, err := c.Events(aggregateID)
	if err != nil {
		log.Println(err)
		model.NewError(ctx, http.StatusBadRequest, errors.New("Không thể lấy lich su."))
		return
	}

	ctx.JSON(http.StatusOK, rs)
}