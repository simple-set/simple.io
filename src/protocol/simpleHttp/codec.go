package simpleHttp

// HttpDecoder http解码器
//type HttpDecoder struct{}

//func (h *HttpDecoder) Input(context *event.HandleContext, reader *bufio.Reader) (*Request, bool) {
//	request, err := NewRequestDecoder(reader).Decoder()
//	if err != nil {
//		logrus.Errorln("Decoding HTTP protocol error, ", err)
//		_ = context.Session().Close()
//		return nil, false
//	}
//
//	request.Response = NewReplyResponse(request)
//	return request, true
//}

//func NewHttpDecoder() *HttpDecoder {
//	return &HttpDecoder{}
//}
//
//// HttpEncoder http编码器
//type HttpEncoder struct{}
//
//func (h HttpEncoder) Output(_ *event.HandleContext, data interface{}) (any, bool) {
//	var response *Response
//	if value, ok := data.(*Request); ok && value.Response != nil {
//		response = value.Response
//	} else if value, ok := data.(*Response); ok {
//		response = value
//	}
//
//	if response == nil {
//		logrus.Warnln("HTTP encoder execution failed with no available response")
//		return nil, false
//	}
//
//	decode, err := NewResponseDecode(response).Decode()
//	if err != nil {
//		logrus.Errorln("HTTP encoder execution failed, ", err)
//		return nil, false
//	}
//	return decode, true
//}
//
//func NewHttpEncoder() *HttpEncoder {
//	return &HttpEncoder{}
//}
