package server

func authenticate(r respondent) (e error) {
	if e = r.initializeSecureChannel(); e != nil {
		return
	}

	if e = r.initializeTransfer(); e != nil {
		return
	}

	// var buffer []byte
	// if buffer, e = r.getMessage(); e != nil {
	// 	return
	// }

	return
}
