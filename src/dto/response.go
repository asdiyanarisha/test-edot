package dto

type (
	ErrorResponse struct {
		Error string `json:"error"`
	}

	ResponseToken struct {
		Token string `json:"token"`
	}

	Response struct {
		Message string `json:"message"`
		Data    any    `json:"data,omitempty"`
	}
)
