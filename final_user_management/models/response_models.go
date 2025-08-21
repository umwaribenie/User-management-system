package models

type PaginatedResponse struct {
	CurrentPage  int         `json:"currentPage"`
	LastPage     int         `json:"lastPage"`
	List         interface{} `json:"list"`
	NextPage     *int        `json:"nextPage,omitempty"`
	PreviousPage *int        `json:"previousPage,omitempty"`
	Status       string      `json:"status"`
	Total        int64       `json:"total"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}
