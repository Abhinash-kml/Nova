package utils

import (
	"encoding/base64"
	"strconv"
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
