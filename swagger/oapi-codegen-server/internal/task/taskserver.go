package task

import (
	"net/http"
	"strconv"
	"time"

	"oapi-codegen-server/internal/taskstore"

	"github.com/labstack/echo/v4"
)

type TaskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *TaskServer {
	store := taskstore.New()
	return &TaskServer{store}
}

func (ts *TaskServer) GetByDueDate(ctx echo.Context, year string, month string, day string) error {
	yearInt, yearErr := strconv.Atoi(year)
	monthInt, monthErr := strconv.Atoi(month)
	dayInt, dayErr := strconv.Atoi(day)
	if yearErr != nil || monthErr != nil || dayErr != nil {
		return ctx.String(http.StatusBadRequest, "Invalid format for parameters")
	}

	tasks := ts.store.GetByDueDate(yearInt, time.Month(monthInt), dayInt)
	return ctx.JSON(http.StatusOK, tasks)
}

func (ts *TaskServer) GetTasksByTag(ctx echo.Context, tagname string) error {
	tasks := ts.store.GetTasksByTag(tagname)

	return ctx.JSON(http.StatusOK, tasks)
}

func (ts *TaskServer) DeleteAllTasks(ctx echo.Context) error {
	if err := ts.store.DeleteAllTasks(); err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (ts *TaskServer) GetAllTasks(ctx echo.Context) error {
	tasks := ts.store.GetAllTasks()

	return ctx.JSON(http.StatusOK, tasks)
}

func (ts *TaskServer) CreateTask(ctx echo.Context) error {
	type Response struct {
		Id int
	}

	var rt RequestTask
	if err := ctx.Bind(&rt); err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid request body")
	}

	var tags []string
	if rt.Tags != nil {
		tags = *rt.Tags
	}

	if rt.Text == nil || rt.Due == nil {
		return ctx.String(http.StatusBadRequest, "Text and Due are required fields")
	}

	id := ts.store.CreateTask(*rt.Text, tags, *rt.Due)
	return ctx.JSON(http.StatusOK, Response{id})
}

func (ts *TaskServer) DeleteTask(ctx echo.Context, id int) error {
	if err := ts.store.DeleteTask(id); err != nil {
		return ctx.String(http.StatusNotFound, err.Error())
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (ts *TaskServer) GetTask(ctx echo.Context, id int) error {
	task, err := ts.store.GetTask(id)
	if err != nil {
		return ctx.String(http.StatusNotFound, err.Error())
	}

	return ctx.JSON(http.StatusOK, task)
}
