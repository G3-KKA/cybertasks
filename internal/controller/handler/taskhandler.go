package handler

import (
	"context"
	"cybertask/internal/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	MessageSuccessful = "successful"
)

type taskhandler struct {
	l     *logger.Logger
	ucase TaskUsecase
}

// Task handler constructor.
func NewTaskHandler(l *logger.Logger, ucase TaskUsecase) *taskhandler {
	return &taskhandler{
		l:     l,
		ucase: ucase,
	}

}

// @Summary     Update task.
// @Description Update task via replacing it with provided task in body.
// @ID          upadte
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param		request body UpdateTaskRequest true "update task"
// @Success     200 {object} UpdateTaskResponse
// @Failure     400 {object} TaskError
// @Failure     404 {object} TaskError
// @Failure     500 {object} TaskError
// @Router      /task [PUT]
func (h *taskhandler) Update(gctx *gin.Context) {
	req := UpdateTaskRequest{}

	err := gctx.BindJSON(&req)
	if err != nil {
		h.l.Err(ErrGotInvalidJSON).Err(err).Send()
		JSONE(gctx, err, http.StatusBadRequest, TaskError{
			Msg: ErrGotInvalidJSON.Error(),
		})

		return
	}

	validate := validator.New()

	err = validate.Struct(req)
	if err != nil {
		h.l.Err(ErrUnsuccessfulValidation).Err(err).Send()
		JSONE(gctx, err, http.StatusBadRequest, TaskError{
			Msg: ErrUnsuccessfulValidation.Error(),
		})

		return
	}

	err = h.ucase.UpdateTask(context.TODO(), req.Task)
	if err != nil {
		h.l.Err(err).Send()
		JSONE(gctx, err, http.StatusInternalServerError, TaskError{
			Msg: ErrInernalError.Error(),
		})

		return
	}

	gctx.JSON(http.StatusOK, UpdateTaskResponse{
		Msg: MessageSuccessful,
	})

	return

}

// @Summary     Create task.
// @Description Create provided in body.
// @ID          create
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param		request body CreateTaskRequest true "create task"
// @Success     200 {object} CreateTaskResponse
// @Failure     400 {object} TaskError
// @Failure     404 {object} TaskError
// @Failure     500 {object} TaskError
// @Router      /task [POST]
func (h *taskhandler) Create(gctx *gin.Context) {
	req := CreateTaskRequest{}

	err := gctx.BindJSON(&req)
	if err != nil {
		h.l.Err(ErrGotInvalidJSON).Err(err).Send()
		JSONE(gctx, err, http.StatusBadRequest, TaskError{
			Msg: ErrGotInvalidJSON.Error(),
		})

		return
	}

	validate := validator.New()

	err = validate.Struct(req)
	if err != nil {
		h.l.Err(ErrUnsuccessfulValidation).Err(err).Send()
		JSONE(gctx, err, http.StatusBadRequest, TaskError{
			Msg: ErrUnsuccessfulValidation.Error(),
		})

		return
	}

	err = h.ucase.CreateTask(context.TODO(), req.Task)
	if err != nil {
		h.l.Err(err).Send()
		JSONE(gctx, err, http.StatusInternalServerError, TaskError{
			Msg: ErrInernalError.Error(),
		})

		return
	}

	gctx.JSON(http.StatusOK, CreateTaskResponse{
		Msg: MessageSuccessful,
	})

	return
}

// @Summary     Delete task.
// @Description Delete task if exist with ID provided in query params.
// @ID          delete
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param		request body DeleteTaskRequest true "delete task"
// @Success     200 {object} DeleteTaskResponse
// @Failure     400 {object} TaskError
// @Failure     404 {object} TaskError
// @Failure     500 {object} TaskError
// @Router      /task/:id [DELETE]
func (h *taskhandler) Delete(gctx *gin.Context) {

	id, err := strconv.ParseUint(gctx.Params.ByName("id"), 10, 64)

	if err != nil {
		h.l.Err(ErrIncorrectID).Err(err).Send()
		JSONE(gctx, err, http.StatusBadRequest, TaskError{
			Msg: ErrIncorrectID.Error(),
		})

		return
	}

	err = h.ucase.DeleteTask(context.TODO(), id)
	if err != nil {
		h.l.Err(err).Send()
		JSONE(gctx, err, http.StatusInternalServerError, TaskError{
			Msg: ErrInernalError.Error(),
		})

		return
	}

	gctx.JSON(http.StatusOK, DeleteTaskResponse{
		Msg: MessageSuccessful,
	})

	return
}

// @Summary     Get task.
// @Description Get task if exist with ID provided in query params.
// @ID          get
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param		request body GetTaskRequest true "get task"
// @Success     200 {object} GetTaskResponse
// @Failure     400 {object} TaskError
// @Failure     404 {object} TaskError
// @Failure     500 {object} TaskError
// @Router      /task/:id [GET]
func (h *taskhandler) Get(gctx *gin.Context) {

	id, err := strconv.ParseUint(gctx.Params.ByName("id"), 10, 64)

	if err != nil {
		h.l.Err(ErrIncorrectID).Err(err).Send()
		JSONE(gctx, err, http.StatusBadRequest, TaskError{
			Msg: ErrIncorrectID.Error(),
		})

		return
	}

	task, err := h.ucase.GetTask(context.TODO(), id)
	if err != nil {
		h.l.Err(err).Send()
		JSONE(gctx, err, http.StatusInternalServerError, TaskError{
			Msg: ErrInernalError.Error(),
		})

		return
	}

	gctx.JSON(http.StatusOK, GetTaskResponse{
		Task: task,
	})

	return
}
