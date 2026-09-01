package api

type APIErrorResponse struct {
	Error APIErrorDetail `json:"error"`
}

type APIErrorDetail struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
