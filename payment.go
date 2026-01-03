package haozpay

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type PaymentService struct {
	client *resty.Client
	config *Config
}

func NewPaymentService(client *resty.Client, config *Config) *PaymentService {
	return &PaymentService{
		client: client,
		config: config,
	}
}

func (s *PaymentService) CreateOrder(ctx context.Context, req *CreatePaymentOrderRequest) (*PaymentOrderResponse, error) {
	bizBodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, &SDKError{
			Code:       ErrInvalidResponse.Code,
			Message:    fmt.Sprintf("failed to marshal request: %v", err),
			StatusCode: 0,
		}
	}

	haozReq := &HaozPayRequest{
		MerchantNo: s.config.MerchantNo,
		Timestamp:  currentTimestampMillis(),
		BizBody:    string(bizBodyBytes),
	}

	var result struct {
		Response
		Data *PaymentOrderResponse `json:"data"`
	}

	_, err = s.client.R().
		SetContext(ctx).
		SetBody(haozReq).
		SetResult(&result).
		Post("/pay-core/payment/order")

	if err != nil {
		return nil, &SDKError{
			Code:       ErrNetworkError.Code,
			Message:    fmt.Sprintf("failed to create payment order: %v", err),
			StatusCode: 0,
		}
	}

	if result.Code != 0 {
		return nil, NewSDKErrorWithRequestID(
			result.Code,
			result.Message,
			0,
			result.RequestID,
		)
	}

	return result.Data, nil
}

func (s *PaymentService) CancelOrder(ctx context.Context, req *CancelPaymentOrderRequest) error {
	bizBodyBytes, err := json.Marshal(req)
	if err != nil {
		return &SDKError{
			Code:       ErrInvalidResponse.Code,
			Message:    fmt.Sprintf("failed to marshal request: %v", err),
			StatusCode: 0,
		}
	}

	haozReq := &HaozPayRequest{
		MerchantNo: s.config.MerchantNo,
		Timestamp:  currentTimestampMillis(),
		BizBody:    string(bizBodyBytes),
	}

	var result Response

	_, err = s.client.R().
		SetContext(ctx).
		SetBody(haozReq).
		SetResult(&result).
		Post("/pay-core/payment/cancel")

	if err != nil {
		return &SDKError{
			Code:       ErrNetworkError.Code,
			Message:    fmt.Sprintf("failed to cancel payment order: %v", err),
			StatusCode: 0,
		}
	}

	if result.Code != 0 {
		return NewSDKErrorWithRequestID(
			result.Code,
			result.Message,
			0,
			result.RequestID,
		)
	}

	return nil
}

func (s *PaymentService) CreateRefund(ctx context.Context, req *CreateRefundRequest) (*RefundResponse, error) {
	// 业务校验: OrderNo 和 ReqSeqId 不能同时为空
	if req.OrderNo == "" && req.ReqSeqId == "" {
		return nil, &SDKError{
			Code:       ErrInvalidRequest.Code,
			Message:    "OrderNo and ReqSeqId cannot both be empty, at least one must be provided",
			StatusCode: 0,
		}
	}

	bizBodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, &SDKError{
			Code:       ErrInvalidResponse.Code,
			Message:    fmt.Sprintf("failed to marshal request: %v", err),
			StatusCode: 0,
		}
	}

	haozReq := &HaozPayRequest{
		MerchantNo: s.config.MerchantNo,
		Timestamp:  currentTimestampMillis(),
		BizBody:    string(bizBodyBytes),
	}

	var result struct {
		Response
		Data *RefundResponse `json:"data"`
	}

	_, err = s.client.R().
		SetContext(ctx).
		SetBody(haozReq).
		SetResult(&result).
		Post("/pay-core/payment/refund")

	if err != nil {
		return nil, &SDKError{
			Code:       ErrNetworkError.Code,
			Message:    fmt.Sprintf("failed to create refund: %v", err),
			StatusCode: 0,
		}
	}

	if result.Code != 0 {
		return nil, NewSDKErrorWithRequestID(
			result.Code,
			result.Message,
			0,
			result.RequestID,
		)
	}

	return result.Data, nil
}

func (s *PaymentService) QueryRefund(ctx context.Context, req *QueryRefundRequest) (*QueryRefundResponse, error) {
	bizBodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, &SDKError{
			Code:       ErrInvalidResponse.Code,
			Message:    fmt.Sprintf("failed to marshal request: %v", err),
			StatusCode: 0,
		}
	}

	haozReq := &HaozPayRequest{
		MerchantNo: s.config.MerchantNo,
		Timestamp:  currentTimestampMillis(),
		BizBody:    string(bizBodyBytes),
	}

	var result struct {
		Response
		Data *QueryRefundResponse `json:"data"`
	}

	_, err = s.client.R().
		SetContext(ctx).
		SetBody(haozReq).
		SetResult(&result).
		Post("/pay-core/payment/refund/query")

	if err != nil {
		return nil, &SDKError{
			Code:       ErrNetworkError.Code,
			Message:    fmt.Sprintf("failed to query refund: %v", err),
			StatusCode: 0,
		}
	}

	if result.Code != 0 {
		return nil, NewSDKErrorWithRequestID(
			result.Code,
			result.Message,
			0,
			result.RequestID,
		)
	}

	return result.Data, nil
}

func currentTimestampMillis() int64 {
	return time.Now().UnixMilli()
}