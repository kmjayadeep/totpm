package types

type OtpInput struct {
	Name       string
	Secret     string
	OtpAuthUri string
}

type Otp struct {
	ID   uint
	Name string
}
