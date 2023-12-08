package simpleHttp

type ResponseBuild struct {
	response *Response
}

func (r *ResponseBuild) Status(code int) *ResponseBuild {
	r.response.SetStatusCode(code)
	return r
}

func (r *ResponseBuild) Header(key, value string) *ResponseBuild {
	r.response.AddHeader(key, value)
	return r
}

func (r *ResponseBuild) Server(server string) *ResponseBuild {
	r.response.Server = server
	return r
}

func (r *ResponseBuild) Body(body []byte) *ResponseBuild {
	return r
}

func (r *ResponseBuild) Build() *Response {
	return r.response
}

func NewResponseBuild() *ResponseBuild {
	return &ResponseBuild{response: new(Response)}
}
