package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LinksContainer struct {
	Self     string `json:"self"`
	Previous string `json:"previous"`
	Next     string `json:"next"`
	First    string `json:"first"`
	Last     string `json:"last"`
}

type MetaContainer struct {
	CurrentPage string `json:"current_page"`
	TotalPage   string `json:"total_page"`
}

type PaginatedResponse[T any] struct {
	Data  []T            `json:"data"`
	Links LinksContainer `json:"_links"`
	Meta  MetaContainer  `json:"metadata"`
}

type ProblemDetailErrors struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

type ProblemDetails struct {
	Type        string                `json:"type"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	StatusCode  int                   `json:"status_code"`
	Instance    string                `json:"instance"`
	Errors      []ProblemDetailErrors `json:"errors,omitempty"`
}

func DecodeCursor(c string) (int, error) {
	if c == "nil" {
		return -1, nil
	}
	bytes, err := base64.URLEncoding.DecodeString(c)
	if err != nil {
		return -1, err
	}
	num, err := strconv.Atoi(string(bytes))
	if err != nil {
		return -1, err
	}

	return num, nil
}

func EncodeToCursor(c int) string {
	return base64.URLEncoding.EncodeToString([]byte(strconv.Itoa(c)))
}

func Paginate[T any](data []T) PaginatedResponse[T] {
	response := PaginatedResponse[T]{
		Data:  data,
		Links: LinksContainer{},
		Meta:  MetaContainer{},
	}

	return response
}

func SendProblemDetails(c *gin.Context, defaultStatus int, err error) {
	problem := ProblemDetails{
		StatusCode: defaultStatus,
	}

	// 1. Check for validation errors
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		problem.Type = "nova.com/validation-error"
		problem.Title = "Request Field Validation Failed"
		problem.Description = "One or more parameters in your request violated structural constrainsts."

		for _, fe := range validationErrors {
			fieldName := strings.ToLower(fe.Field())
			reason := fmt.Sprintf("Rule %q failed", fe.Tag())
			if fe.Param() != "" {
				reason += "=" + fe.Param()
			}

			problem.Errors = append(problem.Errors, ProblemDetailErrors{
				Field:   fieldName,
				Message: reason,
				Code:    "400",
			})
		}

		c.Header("Content-Type", "application/problem+json")
		c.JSON(problem.StatusCode, problem)
		return
	}

	// 2. Check for Syntax/marshalling failures
	// Gin wraps unmarshalling syntax failures inside a unmarshaltypeError
	var typeError interface {
		Field() string
		Type() string
	}

	// Extract structural metrics errors safely if they exist
	if errors.As(err, &typeError) {
		problem.Type = "nova.com/malformed-parameters"
		problem.Title = "Malformed Parameter Data Type"
		problem.Description = fmt.Sprintf("The field %q couldn't be parsed. Expected a valid value matching data type %q",
			strings.ToLower(typeError.Field()), typeError.Type())

		c.Header("Content-Type", "application/problem+json")
		c.JSON(problem.StatusCode, problem)
		return
	}

	// Commom errors

	// Resource not found
	if errors.Is(err, common.ErrResourceNotFound) {
		problem.Type = "nova.com/no-resource"
		problem.Title = "Resource Not Found"
		problem.Description = "The requested resource cannot be found"

		c.Header("Content-Type", "application/problem+json")
		c.JSON(problem.StatusCode, problem)
		return
	}

	// Resource already exists
	if errors.Is(err, common.ErrResourceExists) {
		problem.Type = "nova.com/already-exists"
		problem.Title = "Resource Exists"
		problem.Description = "The requested resource already exists"

		c.Header("Content-Type", "application/problem+json")
		c.JSON(problem.StatusCode, problem)
		return
	}
}
