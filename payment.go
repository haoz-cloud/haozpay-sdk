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

func (s *PaymentService) CreateWithdraw(ctx context.Context, req *CreateWithdrawRequest) error {
	var result Response

	_, err := s.client.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&result).
		Post("/pay-core/account/withdraw")

	if err != nil {
		return &SDKError{
			Code:       ErrNetworkError.Code,
			Message:    fmt.Sprintf("failed to create withdraw: %v", err),
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

func (s *PaymentService) ApplySplit(ctx context.Context, req *ApplySplitRequest) error {
	acctSplitBunchesJSON, err := json.Marshal(req.AcctSplitBunches)
	if err != nil {
		return &SDKError{
			Code:       ErrInvalidRequest.Code,
			Message:    fmt.Sprintf("failed to marshal acctSplitBunches: %v", err),
			StatusCode: 0,
		}
	}

	timestamp := currentTimestampMillis()
	paramsMap := map[string]interface{}{
		"merchantNo":       s.config.MerchantNo,
		"timestamp":        timestamp,
		"orderNo":          req.OrderNo,
		"splitAmount":      req.SplitAmount,
		"acctSplitBunches": string(acctSplitBunchesJSON),
	}

	sign, err := GenerateSign(paramsMap, s.config.PrivateKey)
	if err != nil {
		return &SDKError{
			Code:       ErrInvalidRequest.Code,
			Message:    fmt.Sprintf("failed to generate sign: %v", err),
			StatusCode: 0,
		}
	}

	body := map[string]interface{}{
		"merchantNo":       s.config.MerchantNo,
		"timestamp":        timestamp,
		"orderNo":          req.OrderNo,
		"splitAmount":      req.SplitAmount,
		"acctSplitBunches": req.AcctSplitBunches,
		"sign":             sign,
	}

	var result Response
	_, err = s.client.R().
		SetContext(ctx).
		SetBody(body).
		SetResult(&result).
		Post("/pay-core/merchant/split/acctSplitBunch")

	if err != nil {
		return &SDKError{
			Code:       ErrNetworkError.Code,
			Message:    fmt.Sprintf("failed to apply split: %v", err),
			StatusCode: 0,
		}
	}

	if result.Code != 0 {
		return NewSDKErrorWithRequestID(result.Code, result.Message, 0, result.RequestID)
	}

	return nil
}

func (s *PaymentService) ReverseSplit(ctx context.Context, req *ReverseSplitRequest) (*ReverseSplitResponse, error) {
	timestamp := currentTimestampMillis()
	paramsMap := map[string]interface{}{
		"merchantNo": s.config.MerchantNo,
		"timestamp":  timestamp,
		"orderNo":    req.OrderNo,
	}
	if req.CancelReason != "" {
		paramsMap["cancelReason"] = req.CancelReason
	}

	sign, err := GenerateSign(paramsMap, s.config.PrivateKey)
	if err != nil {
		return nil, &SDKError{
			Code:       ErrInvalidRequest.Code,
			Message:    fmt.Sprintf("failed to generate sign: %v", err),
			StatusCode: 0,
		}
	}

	body := map[string]interface{}{
		"merchantNo": s.config.MerchantNo,
		"timestamp":  timestamp,
		"orderNo":    req.OrderNo,
		"sign":       sign,
	}
	if req.CancelReason != "" {
		body["cancelReason"] = req.CancelReason
	}

	var result struct {
		Response
		Data *ReverseSplitResponse `json:"data"`
	}

	_, err = s.client.R().
		SetContext(ctx).
		SetBody(body).
		SetResult(&result).
		Post("/pay-core/merchant/split/reverseSplitRefund")

	if err != nil {
		return nil, &SDKError{
			Code:       ErrNetworkError.Code,
			Message:    fmt.Sprintf("failed to reverse split: %v", err),
			StatusCode: 0,
		}
	}

	if result.Code != 0 {
		return nil, NewSDKErrorWithRequestID(result.Code, result.Message, 0, result.RequestID)
	}

	return result.Data, nil
}

func (s *PaymentService) QueryHostingOrderPage(ctx context.Context, req *QueryHostingOrderRequest) (*HostingOrderPageResponse, error) {
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
		Data *HostingOrderPageResponse `json:"data"`
	}

	_, err = s.client.R().
		SetContext(ctx).
		SetBody(haozReq).
		SetResult(&result).
		Post("/pay-core/payment/queryHostingOrderPage")

	if err != nil {
		return nil, &SDKError{
			Code:       ErrNetworkError.Code,
			Message:    fmt.Sprintf("failed to query hosting order page: %v", err),
			StatusCode: 0,
		}
	}

	if result.Code != 0 {
		return nil, NewSDKErrorWithRequestID(result.Code, result.Message, 0, result.RequestID)
	}

	return result.Data, nil
}

func (s *PaymentService) QueryOrderDetail(ctx context.Context, req *QueryOrderDetailRequest) (*OrderDetailResponse, error) {
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
		Data *OrderDetailResponse `json:"data"`
	}

	_, err = s.client.R().
		SetContext(ctx).
		SetBody(haozReq).
		SetResult(&result).
		Post("/pay-core/payment/queryOrderDetail")

	if err != nil {
		return nil, &SDKError{
			Code:       ErrNetworkError.Code,
			Message:    fmt.Sprintf("failed to query order detail: %v", err),
			StatusCode: 0,
		}
	}

	if result.Code != 0 {
		return nil, NewSDKErrorWithRequestID(result.Code, result.Message, 0, result.RequestID)
	}

	return result.Data, nil
}

func currentTimestampMillis() int64 {
	return time.Now().UnixMilli()
}
