package main

import (
	"swapnil-ex/handlers"
	"swapnil-ex/models/db"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"

)

func main() {

	defer db.Close()






	e := echo.New()

	// e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200", "https://labstack.net", "*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.POST("/register", handlers.Register, handlers.OnlySwapnil())
	e.POST("/login", handlers.Login)
	e.PUT("/updateUser", handlers.UpdateUser, handlers.IsLoggedIn)
	e.DELETE("/logout", handlers.Logout, handlers.IsLoggedIn)
	e.GET("/students", handlers.GetStudents, handlers.IsLoggedIn)
	e.GET("/students/:id", handlers.GetStudent, handlers.IsLoggedIn)
	e.GET("/students/:id/hostel", handlers.GetStudentHostel, handlers.IsLoggedIn)
	e.POST("/students", handlers.CreateStudent, handlers.IsLoggedIn)
	e.PUT("/students/:id", handlers.UpdateStudent, handlers.IsLoggedIn)
	e.DELETE("/students/:id", handlers.DeleteStudent, handlers.IsLoggedIn)
	e.GET("/standards", handlers.GetStandards, handlers.IsLoggedIn)
	e.GET("/standards/:id", handlers.GetStandard, handlers.IsLoggedIn)
	e.POST("/standards", handlers.CreateStandard, handlers.IsLoggedIn)
	e.PUT("/standards/:id", handlers.UpdateStandard, handlers.IsLoggedIn)
	e.DELETE("/standards/:id", handlers.DeleteStandard, handlers.IsLoggedIn)
	e.GET("/batchs", handlers.GetBatchs, handlers.IsLoggedIn)
	e.GET("/batchs/:id", handlers.GetBatch, handlers.IsLoggedIn)
	e.POST("/batchs", handlers.CreateBatch, handlers.IsLoggedIn)
	e.PUT("/batchs/:id", handlers.UpdateBatch, handlers.IsLoggedIn)
	e.DELETE("/batchs/:id", handlers.DeleteBatch, handlers.IsLoggedIn)


	e.Logger.Fatal(e.Start(":8080"))
}
