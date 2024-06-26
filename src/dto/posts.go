package dto

type (
	AddPost struct {
		Id      int      `json:"id,omitempty"`
		Title   string   `json:"title" binding:"required"`
		Content string   `json:"content" binding:"required"`
		Tags    []string `json:"tags" binding:"required"`
	}

	ParamGetPost struct {
		Offset int `form:"offset"`
		Limit  int `form:"limit"`
	}
)
