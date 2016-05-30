package server

func authenticate(r respondent, ctx *context) (e error) {
	if e = r.initializeSecureChannel(); e != nil {
		return
	}

	// var buffer []byte
	// if buffer, e = r.getMessage(); e != nil {
	// 	return
	// }

	return
}
