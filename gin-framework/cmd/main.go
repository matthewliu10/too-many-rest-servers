package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"gin-framework/internal/taskstore"

	"github.com/gin-gonic/gin"
)

type taskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *taskServer {
	store := taskstore.New()
	return &taskServer{store}
}

func (ts *taskServer) createTaskHandler(c *gin.Context) {
	type TaskStructure struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	var task TaskStructure
	if err := c.ShouldBindJSON(&task); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	id := ts.store.CreateTask(task.Text, task.Tags, task.Due)
	c.JSON(http.StatusOK, gin.H{"Id": id})
}

func (ts *taskServer) getAllTasksHandler(c *gin.Context) {
	tasks := ts.store.GetAllTasks()
	c.JSON(http.StatusOK, tasks)
}

func (ts *taskServer) deleteAllTasksHandler(c *gin.Context) {
	if err := ts.store.DeleteAllTasks(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (ts *taskServer) getTaskHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id")) // PRINT WHAT ID IS -------------------------------
	if err != nil {
		c.String(http.StatusBadRequest, "invalid id")
		return
	}

	task, err := ts.store.GetTask(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, task)
}

func (ts *taskServer) deleteTaskHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid id")
	}

	if err = ts.store.DeleteTask(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.Status(http.StatusNoContent)
}

func (ts *taskServer) tagHandler(c *gin.Context) {
	tag := c.Param("tag")
	tasks := ts.store.GetTasksByTag(tag)
	c.JSON(http.StatusOK, tasks)
}

func (ts *taskServer) dueHandler(c *gin.Context) {
	year, yearErr := strconv.Atoi(c.Param("year"))
	month, monthErr := strconv.Atoi(c.Param("month"))
	day, dayErr := strconv.Atoi(c.Param("day"))
	if yearErr != nil || monthErr != nil || dayErr != nil {
		c.String(http.StatusBadRequest, "expect /due/<year>/<month>/<day>/")
		return
	}

	tasks := ts.store.GetByDueDate(year, time.Month(month), day)
	c.JSON(http.StatusOK, tasks)
}

func main() {
	router := gin.Default()
	server := NewTaskServer()

	router.POST("/task/", server.createTaskHandler)
	router.GET("/task/", server.getAllTasksHandler)
	router.DELETE("/task/", server.deleteAllTasksHandler)
	router.GET("/task/:id", server.getTaskHandler)
	router.DELETE("/task/:id", server.deleteTaskHandler)
	router.GET("/tag/:tag/", server.tagHandler)
	router.GET("/due/:year/:month/:day", server.dueHandler)

	router.Run("localhost:" + os.Getenv("SERVERPORT"))
}
