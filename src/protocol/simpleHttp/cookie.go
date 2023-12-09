package simpleHttp

//type Cookie struct {
//	http.Cookie
//}
//
//func NewCookie(name, value string) *Cookie {
//	name = sanitizeCookieName(name)
//	value = sanitizeCookieValue(value)
//	return &Cookie{http.Cookie{Name: name, Value: value, Path: "/", MaxAge: 86400}}
//}
//
//func sanitizeCookieName(n string) string {
//	return strings.NewReplacer("\n", "-", "\r", "-").Replace(n)
//}
//
//func sanitizeCookieValue(v string) string {
//	v = sanitizeOrWarn("Cookie.Value", validCookieValueByte, v)
//	if len(v) == 0 {
//		return v
//	}
//	if strings.ContainsAny(v, " ,") {
//		return `"` + v + `"`
//	}
//	return v
//}
//
//func sanitizeOrWarn(fieldName string, valid func(byte) bool, v string) string {
//	ok := true
//	for i := 0; i < len(v); i++ {
//		if valid(v[i]) {
//			continue
//		}
//		log.Printf("net/http: invalid byte %q in %s; dropping invalid bytes", v[i], fieldName)
//		ok = false
//		break
//	}
//	if ok {
//		return v
//	}
//	buf := make([]byte, 0, len(v))
//	for i := 0; i < len(v); i++ {
//		if b := v[i]; valid(b) {
//			buf = append(buf, b)
//		}
//	}
//	return string(buf)
//}
//
//func validCookieValueByte(b byte) bool {
//	return 0x20 <= b && b < 0x7f && b != '"' && b != ';' && b != '\\'
//}
