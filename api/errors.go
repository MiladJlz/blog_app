package api

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}
	apiError := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiError.Code).JSON(apiError)
}

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func ErrNotResourceNotFound(err error) Error {
	return Error{
		Code: http.StatusNotFound,
		Err:  "resource not found -> " + err.Error(),
	}
}

func ErrBadRequest(err error) Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid JSON request -> " + err.Error(),
	}
}
