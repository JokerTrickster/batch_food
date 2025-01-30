package request

type ReqAllMessageUser struct {
	Role    string `json:"role"`
	Title   string `json:"title"`
	Message string `json:"message"`
}
