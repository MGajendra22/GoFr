package main

import (
	"fmt"
	"github.com/MGajendra22/GoFr/handler/task"
	"github.com/MGajendra22/GoFr/handler/user"
	"github.com/MGajendra22/GoFr/migrations"

	taskServicePkg "github.com/MGajendra22/GoFr/service/task"
	userServicePkg "github.com/MGajendra22/GoFr/service/user"
	taskStorePkg "github.com/MGajendra22/GoFr/store/task"
	userStorePkg "github.com/MGajendra22/GoFr/store/user"
	"gofr.dev/pkg/gofr"
)

func main() {

	userStore := userStorePkg.NewUserStore()
	userService := userServicePkg.NewUserService(userStore)
	userHandler := user.NewUserHandler(userService)
	// Init task dependencies
	taskStore := taskStorePkg.NewStore()
	taskService := taskServicePkg.NewService(taskStore, userService)
	taskHandler := task.NewHandler(taskService)

	app := gofr.New()

	app.Migrate(migrations.All())

	app.POST("/task", taskHandler.Create)
	app.GET("/task/{id}", taskHandler.GetTask)
	app.GET("/task", taskHandler.All)
	app.PUT("/task/{id}", taskHandler.Complete)
	app.DELETE("/task/{id}", taskHandler.Delete)
	app.GET("task/user/{id}", taskHandler.GetTasksByUserID)

	app.POST("/user", userHandler.Create)
	app.GET("/user", userHandler.All)
	app.GET("/user/{id}", userHandler.Get)
	app.DELETE("/user/{id}", userHandler.Delete)

	fmt.Println("Server running at http://localhost:8000")
	app.AddHTTPService("Task-Manager Http Service", "http://localhost:8000")

	app.Run()
}
