package utils

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/example"
	"github.com/midtrans/midtrans-go/iris"
)

var irisCreator iris.Client
var irisApprover iris.Client

func setupIrisGateway() {
	irisCreator.New(example.IrisCreatorKeySandbox, midtrans.Sandbox)
	irisApprover.New(example.IrisApproverKeySandbox, midtrans.Sandbox)
}

func createPayout(name string, number string, numberType string, email string, amount string) []iris.CreatePayoutDetailResponse {
	p := iris.CreatePayoutDetailReq{
		BeneficiaryName:    name,
		BeneficiaryAccount: number,
		BeneficiaryBank:    numberType,
		BeneficiaryEmail:   email,
		Amount:             amount,
		Notes:              "Congrats this is your money",
	}
	var payouts []iris.CreatePayoutDetailReq
	payouts = append(payouts, p)

	cp := iris.CreatePayoutReq{Payouts: payouts}

	payoutReps, _ := irisCreator.CreatePayout(cp)
	return payoutReps.Payouts
}

func CreateAndApprovePayout(name string, number string, numberType string, email string, amount string) (*iris.ApprovePayoutResponse, string) {
	setupIrisGateway()
	var payouts = createPayout(name, number, numberType, email, amount)

	var refNos []string
	refNos = append(refNos, payouts[0].ReferenceNo)

	ap := iris.ApprovePayoutReq{
		ReferenceNo: refNos,
		OTP:         "335163",
	}

	approveResp, _ := irisApprover.ApprovePayout(ap)
	return approveResp, payouts[0].ReferenceNo
}
