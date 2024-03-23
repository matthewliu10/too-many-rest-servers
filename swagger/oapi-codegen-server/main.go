package main

import (
	"oapi-codegen-server/internal/task"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	ts := task.NewTaskServer()

	task.RegisterHandlers(e, ts)
	e.Logger.Fatal(e.Start("localhost:8080"))
}
