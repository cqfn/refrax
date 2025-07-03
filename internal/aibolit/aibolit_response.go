package aibolit

type AibolitResponse struct {
	data string
}

func NewAibolitResponse(data string) *AibolitResponse {
	return &AibolitResponse{data}
}

func (r *AibolitResponse) Sanitized() string {
	return ""
}
