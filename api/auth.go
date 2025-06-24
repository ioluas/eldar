package api

func NewSignUpRequest(email, password string, options ...SignUpOption) *SignUpRequest {
	sr := SignUpRequest{
		Email:    email,
		Password: password,
		Options:  map[string]interface{}{},
	}
	for _, option := range options {
		sr = option(sr)
	}
	return &sr
}

func (c *DatabaseHTTPClient) SignUp(request *SignUpRequest) SignUpResponse {
	return SignUpResponse{} // @todo.
}
