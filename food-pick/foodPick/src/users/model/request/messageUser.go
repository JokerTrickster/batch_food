package request

type ReqMessageUser struct {
	Role    string `json:"role"`
	UserID  int    `json:"userId"`
	Title   string `json:"title"`
	Message string `json:"message"`
}
