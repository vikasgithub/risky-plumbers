package entity

type Risk struct {
	ID          string `json:"id"`
	State       string `json:"state"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
