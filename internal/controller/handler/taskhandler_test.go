package handler

import (
	"bytes"
	mock_handler "cybertask/internal/controller/handler/mocks"
	"cybertask/internal/logger"
	"cybertask/model"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var _happyPathTask = model.Task{
	ID:          22,
	Header:      "hdr",
	Description: "",
	CreatedAt:   time.Now(),
	Status:      false,
}

func TestGetHappyPath(t *testing.T) {
	t.Parallel()

	// Prepare gin context.
	w := httptest.NewRecorder()
	gctx, _ := gin.CreateTestContext(w)

	req := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	gctx.Request = req

	gctx.Params = append(
		gctx.Params,
		gin.Param{
			Key:   "id",
			Value: strconv.Itoa(int(_happyPathTask.ID))},
	)

	tlogger := logger.NewTesting(t)

	ucase := mock_handler.NewMockTaskUsecase(t)
	ucase.On("GetTask",
		mock.Anything,
		mock.Anything).
		Return(_happyPathTask, error(nil))

	h := NewTaskHandler(tlogger, ucase)

	h.Get(gctx)
	rsp := w.Result()
	defer rsp.Body.Close()
	bytes, err := io.ReadAll(rsp.Body)
	assert.NoError(t, err)

	compareResult := func() {

		rsp := GetTaskResponse{}

		err = sonic.Unmarshal(bytes, &rsp)

		assert.NoError(t, err)

		task := rsp.Task

		timediff := _happyPathTask.CreatedAt.Compare(task.CreatedAt)

		assert.Equal(t, 0, timediff)
		assert.Equal(t, _happyPathTask.ID, task.ID)
		assert.Equal(t, _happyPathTask.Header, task.Header)
		assert.Equal(t, _happyPathTask.Description, task.Description)
		assert.Equal(t, _happyPathTask.Status, task.Status)
	}
	compareResult()

}

func TestUpdateHappyPath(t *testing.T) {
	t.Parallel()

	_happyPathTask := model.Task{
		ID:          22,
		Header:      "hdr",
		Description: "",
		CreatedAt:   time.Now(),
		Status:      false,
	}
	testreq := UpdateTaskRequest{
		Task: _happyPathTask,
	}

	// Prepare gin context.
	w := httptest.NewRecorder()
	gctx, _ := gin.CreateTestContext(w)

	body, err := sonic.Marshal(testreq)
	require.NoError(t, err)

	buf := bytes.NewBuffer(body)

	reqbody := io.NopCloser(buf)

	req := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
		Body:   reqbody,
	}
	gctx.Request = req

	tlogger := logger.NewTesting(t)

	ucase := mock_handler.NewMockTaskUsecase(t)
	ucase.On("UpdateTask",
		mock.Anything,
		mock.Anything).
		Return(error(nil))

	h := NewTaskHandler(tlogger, ucase)

	h.Update(gctx)
	rsp := w.Result()
	defer rsp.Body.Close()
	bytes, err := io.ReadAll(rsp.Body)
	assert.NoError(t, err)

	compareResult := func() {

		rsp := UpdateTaskResponse{}

		err = sonic.Unmarshal(bytes, &rsp)

		require.NoError(t, err)

		assert.Equal(t, rsp.Msg, MessageSuccessful)
	}
	compareResult()

}
func TestCreateHappyPath(t *testing.T) {
	t.Parallel()

	_happyPathTask := model.Task{
		ID:          22,
		Header:      "hdr",
		Description: "",
		CreatedAt:   time.Now(),
		Status:      false,
	}
	testreq := CreateTaskRequest{
		Task: _happyPathTask,
	}

	// Prepare gin context.
	w := httptest.NewRecorder()
	gctx, _ := gin.CreateTestContext(w)

	body, err := sonic.Marshal(testreq)
	require.NoError(t, err)

	buf := bytes.NewBuffer(body)

	reqbody := io.NopCloser(buf)

	req := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
		Body:   reqbody,
	}
	gctx.Request = req

	tlogger := logger.NewTesting(t)

	ucase := mock_handler.NewMockTaskUsecase(t)
	ucase.On("CreateTask",
		mock.Anything,
		mock.Anything).
		Return(error(nil))

	h := NewTaskHandler(tlogger, ucase)

	h.Create(gctx)

	rsp := w.Result()
	defer rsp.Body.Close()
	bytes, err := io.ReadAll(rsp.Body)
	assert.NoError(t, err)

	compareResult := func() {

		rsp := CreateTaskResponse{}

		err = sonic.Unmarshal(bytes, &rsp)

		require.NoError(t, err)

		assert.Equal(t, rsp.Msg, MessageSuccessful)
	}
	compareResult()

}
func TestDeleteHappyPath(t *testing.T) {
	t.Parallel()

	_happyPathTask := model.Task{
		ID:          22,
		Header:      "hdr",
		Description: "",
		CreatedAt:   time.Now(),
		Status:      false,
	}

	// Prepare gin context.
	w := httptest.NewRecorder()
	gctx, _ := gin.CreateTestContext(w)

	req := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	gctx.Request = req

	gctx.Params = append(
		gctx.Params,
		gin.Param{
			Key:   "id",
			Value: strconv.Itoa(int(_happyPathTask.ID))},
	)

	tlogger := logger.NewTesting(t)

	ucase := mock_handler.NewMockTaskUsecase(t)
	ucase.On("DeleteTask",
		mock.Anything,
		mock.Anything).
		Return(error(nil))

	h := NewTaskHandler(tlogger, ucase)

	h.Delete(gctx)
	rsp := w.Result()
	defer rsp.Body.Close()
	bytes, err := io.ReadAll(rsp.Body)
	assert.NoError(t, err)

	compareResult := func() {

		rsp := DeleteTaskResponse{}

		err = sonic.Unmarshal(bytes, &rsp)

		require.NoError(t, err)

		assert.Equal(t, rsp.Msg, MessageSuccessful)
	}
	compareResult()

}
