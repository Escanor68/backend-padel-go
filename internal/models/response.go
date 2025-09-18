package models

// Respuestas estándar de la API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// Funciones helper para crear respuestas
func NewSuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: "Success",
		Data:    data,
	}
}

func NewErrorResponse(message, error string) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
		Error:   error,
	}
}

func NewPaginatedResponse(data interface{}, pagination Pagination) PaginatedResponse {
	return PaginatedResponse{
		Success:    true,
		Message:    "Success",
		Data:       data,
		Pagination: pagination,
	}
}
