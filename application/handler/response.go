package handler

type (
	response struct {
		Status string `json:"status"`
	}

	errorResponse struct {
		Error string `json:"error"`
	}
)

func newResponse(s string) *response {
	return &response{
		Status: s,
	}
}

func newErrorResponse(err string) *errorResponse {
	return &errorResponse{
		Error: err,
	}
}
