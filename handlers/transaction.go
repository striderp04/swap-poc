package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"

	"github.com/labstack/echo/v4"
)

func GetStudentTransactions(c echo.Context) error {
	// Get student
	studentId := c.Param("student_id")
	newStudentId, err := strconv.Atoi(studentId)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}
	student := &models.Student{ID: uint(newStudentId)}
	err = student.Find()
	if err != nil {
		fmt.Println("s.Find(GetStudent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	transaction := &models.Transaction{}
	transactions, err := transaction.All(newStudentId)
	if err != nil {
		fmt.Println("s.ALL(GetTransactions)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, transactions)
}

func GetStudentTransaction(c echo.Context) error {
	// Get a single user by ID
	studentId := c.Param("student_id")
	newStudentId, err := strconv.Atoi(studentId)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}
	student := &models.Student{ID: uint(newStudentId)}
	err = student.Find()
	if err != nil {
		fmt.Println("s.Find(GetStudent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, student)
}
