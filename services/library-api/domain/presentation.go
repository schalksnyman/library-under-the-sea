package domain

type FindByTitleRequest struct {
	Title string `json:"title"`
}

type FindBookRequest struct {
	Id string `json:"id"`
}

type UpdateTitleRequest struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type DeleteBookRequest struct {
	Id string `json:"id"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
