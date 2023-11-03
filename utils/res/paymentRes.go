package res

import (
	"app/model"
	"app/model/web"
)

func PaymentConvertIndex(payments []model.Payment) []web.PaymentResponse {
	var results []web.PaymentResponse
	for _, payment := range payments {
		paymentResponse := web.PaymentResponse{
			ID:          int(payment.ID),
			UserId:      payment.UserId,
			PointId:     payment.PointId,
			Amount:      payment.Amount,
			Type:        payment.Type,
			Number:      payment.Number,
			Status:      payment.Status,
			ReferenceNo: payment.ReferenceNo,
		}
		results = append(results, paymentResponse)
	}
	return results
}

func PassPaymentBody(payment *model.Payment) web.PaymentResponse {
	return web.PaymentResponse{
		ID:          int(payment.ID),
		UserId:      payment.UserId,
		PointId:     payment.PointId,
		Amount:      payment.Amount,
		Type:        payment.Type,
		Number:      payment.Number,
		Status:      payment.Status,
		ReferenceNo: payment.ReferenceNo,
	}
}
