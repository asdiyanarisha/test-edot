package dto

type (
	ErrorResponse struct {
		Error string `json:"error"`
	}

	Response struct {
		Message string `json:"message"`
	}
)
