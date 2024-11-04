package request

type SaveUser struct {
	Login    string `json:"login,string"`
	Password string `json:"password,string"`
	Name     string `json:"name,string"`
}
