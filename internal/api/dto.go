package api

type putRequest struct {
	Value string `json:"value"`
}

type valueResponse struct {
	Value string `json:"value"`
}

type errorResponse struct {
	Error string `json:"error"`
}
