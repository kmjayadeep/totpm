package types

type Account struct {
	ID        string `json:"id"`
	Service   string `json:"service"`
	Account   string `json:"account"`
	Icon      string `json:"icon"`
	OtpType   string `json:"otp_type"`
	Digits    int    `json:"digits"`
	Algorithm string `json:"algorithm"`
	Period    int    `json:"period"`
	Counter   int    `json:"counter"`
	Secret    string `json:"secret"`
}
