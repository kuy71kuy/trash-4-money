package web

type CreatePayoutResponse struct {
	Payouts []struct {
		Status      string `json:"status"`
		ReferenceNo string `json:"reference_no"`
	} `json:"payouts"`
}

type ApprovePayoutResponse struct {
	Status string `json:"status"`
}
