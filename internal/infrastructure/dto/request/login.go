package request

type Login struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
