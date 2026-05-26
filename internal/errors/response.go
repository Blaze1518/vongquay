package errors

import "time"

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
	Path 	   string      `json:"path,omitempty"`
	RequestID  string      `json:"request_id,omitempty"`
	RetryAfter *int        `json:"retry_after,omitempty"`
}

type Meta struct {
	RequestID  string    `json:"request_id,omitempty"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
	Page       int       `json:"page,omitempty"`
	PerPage    int       `json:"per_page,omitempty"`
	Total      int64     `json:"total,omitempty"`
	TotalPages int       `json:"total_pages,omitempty"`
	Links      *Links    `json:"links,omitempty"`
}

type Links struct {
	Self  string `json:"self,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
}

func Success(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}