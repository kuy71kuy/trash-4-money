package req

import (
	"app/model"
	"app/model/web"
)

func PaymentPassBody(payment web.PaymentResponse) *model.Payment {
	return &model.Payment{
		UserId:      payment.UserId,
		PointId:     payment.PointId,
		Amount:      payment.Amount,
		Type:        payment.Type,
		Number:      payment.Number,
		Status:      payment.Status,
		ReferenceNo: payment.ReferenceNo,
	}
}
