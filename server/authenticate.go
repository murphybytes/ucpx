package server

func authenticate(r respondent, ctx *context) (e error) {
	if e = r.initializeSecureChannel(); e != nil {
		return
	}

	return
}
