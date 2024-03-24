package main

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"strconv"
	"time"

	"stdlib/internal/taskstore"
)

type taskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *taskServer {
	store := taskstore.New()
	return &taskServer{store}
}

func renderJSON(w http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling create task at %s\n", r.URL.Path)

	type TaskStructure struct {
		Text string
		Tags []string
		Due  time.Time
	}

	type Response struct {
		Id int
	}

	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var taskStruct TaskStructure
	if err := decoder.Decode(&taskStruct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateTask(taskStruct.Text, taskStruct.Tags, taskStruct.Due)
	renderJSON(w, Response{id})
}

func (ts *taskServer) getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get all tasks at %s\n", r.URL.Path)

	tasks := ts.store.GetAllTasks()

	renderJSON(w, tasks)
}

func (ts *taskServer) deleteAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling delete all tasks at %s\n", r.URL.Path)

	if err := ts.store.DeleteAllTasks(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ts *taskServer) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get task at %s\n", r.URL.Path)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}

func (ts *taskServer) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling delete task at %s\n", r.URL.Path)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err = ts.store.DeleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ts *taskServer) tagHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get tasks by tag at %s\n", r.URL.Path)

	tag := r.PathValue("tag")
	tasks := ts.store.GetTasksByTag(tag)

	renderJSON(w, tasks)
}

func (ts *taskServer) dueHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get tasks by due date at %s\n", r.URL.Path)

	year, yearErr := strconv.Atoi(r.PathValue("year"))
	month, monthErr := strconv.Atoi(r.PathValue("month"))
	day, dayErr := strconv.Atoi(r.PathValue("day"))
	if yearErr != nil || monthErr != nil || dayErr != nil {
		http.Error(w, "expect /due/<year>/<month>/<day>/", http.StatusBadRequest)
		return
	}
	tasks := ts.store.GetByDueDate(year, time.Month(month), day)

	renderJSON(w, tasks)
}

func main() {
	mux := http.NewServeMux()
	server := NewTaskServer()
	mux.HandleFunc("POST /task/", server.createTaskHandler)
	mux.HandleFunc("GET /task/", server.getAllTasksHandler)
	mux.HandleFunc("DELETE /task/", server.deleteAllTasksHandler)
	mux.HandleFunc("GET /task/{id}", server.getTaskHandler)
	mux.HandleFunc("DELETE /task/{id}", server.deleteTaskHandler)
	mux.HandleFunc("GET /tag/{tag}", server.tagHandler)
	mux.HandleFunc("GET /due/{year}/{month}/{day}", server.dueHandler)

	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}
