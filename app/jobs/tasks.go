package jobs

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	TypeCreateProductOrderDigiflazz = "product-order:digiflazz:create"
	TypeEmailVerification           = "email:verification"
	TypeEmailResetPassword          = "email:reset-password"
	TypeProductOrderSendWAInvoice   = "product-order:wa:invoice"
)

type EmailPayload struct {
	To      string
	Subject string
	Body    string
}

type EmailVerificationPayload struct {
	To    string
	Token string
}

type EmailResetPasswordPayload struct {
	To    string
	Token string
}

type ProductOrderSendWAInvoicePayload struct {
	OrderID string
}

func NewEmailResetPasswordTask(to string, token string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailResetPasswordPayload{To: to, Token: token})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeEmailResetPassword, payload, asynq.MaxRetry(3), asynq.Queue("high")), nil
}

func NewEmailVerificationTask(to string, token string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailVerificationPayload{To: to, Token: token})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeEmailVerification, payload, asynq.MaxRetry(3), asynq.Queue("high")), nil
}

type ProductOrderPayload struct {
	ProductOrderID *string
	ProductSKU     string
	ReferenceID    *string
	CustomerNo     *string
}

func NewProductOrderDigiflazzTask(productOrderID *string, productSKU string, referenceID *string, customerNo *string) (*asynq.Task, error) {
	payload, err := json.Marshal(ProductOrderPayload{ProductOrderID: productOrderID, ProductSKU: productSKU, ReferenceID: referenceID, CustomerNo: customerNo})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeCreateProductOrderDigiflazz, payload, asynq.MaxRetry(3), asynq.Queue("critical")), nil
}

func NewProductOrderSendWAInvoiceTask(orderID string) (*asynq.Task, error) {
	payload, err := json.Marshal(ProductOrderSendWAInvoicePayload{OrderID: orderID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeProductOrderSendWAInvoice, payload, asynq.MaxRetry(3), asynq.Queue("high")), nil
}
