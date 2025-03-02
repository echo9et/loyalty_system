package entities

type User struct {
	Login    int    `json:"login"`
	Password string `json:"password"`
}
