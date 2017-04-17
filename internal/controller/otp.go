package controller

func OTPAuth(sn, token string) bool {
	if sn != "" {
		return true
	}
	return false
}
