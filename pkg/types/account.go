package types

type Account struct {
	ID        string `json:"id"`
	Service   string `json:"service"`
	Account   string `json:"account"`
	Icon      string `json:"icon"`
	OtpType   string `json:"otp_type"`
	Digits    uint   `json:"digits"`
	Algorithm string `json:"algorithm"`
	Period    uint   `json:"period"`
	Counter   uint64 `json:"counter"`
	Secret    string `json:"secret"`
}
