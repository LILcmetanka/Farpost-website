// В файле internal/lib/api/response/response.go добавьте:

package response

// Response represents a standard API response
// @Description Стандартный ответ API
type Response struct {
	// Статус операции: OK или ERROR
	Status string `json:"status" example:"OK"`
	// Сообщение об ошибке (если статус ERROR)
	Error string `json:"error,omitempty" example:""`
}

// Ok returns a successful response
func Ok() Response {
	return Response{
		Status: "OK",
	}
}

// Error returns an error response
func Error(msg string) Response {
	return Response{
		Status: "ERROR",
		Error:  msg,
	}
}