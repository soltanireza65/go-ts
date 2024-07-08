package store

type Todo struct {
	Title string `json:"title" validate:"required"`
	Done  bool   `json:"done" validate:"required"`
}
