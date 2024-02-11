package auth

func newTestAdmin() admin {
	return admin{PasswordHash: "test hash"}
}
