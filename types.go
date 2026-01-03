package haozpay

import "time"

type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp int64       `json:"timestamp,omitempty"`
}

type HaozPayRequest struct {
	MerchantNo string `json:"merchantNo"`
	Timestamp  int64  `json:"timestamp"`
	BizBody    string `json:"bizBody"`
	Sign       string `json:"sign"`
}

type CreatePaymentOrderRequest struct {
	OrderTitle        string  `json:"orderTitle"`
	OrderAmount       float64 `json:"orderAmount"`
	PayType           int     `json:"payType"`
	UseHaozPayCashier bool    `json:"useHaozPayCashier"`
	NotifyUrl         string  `json:"notifyUrl"`
	RedirectUrl       string  `json:"redirectUrl,omitempty"`
}

type PaymentOrderResponse struct {
	MerchantNo      string  `json:"merchantNo"`
	ChannelType     string  `json:"channelType"`
	SeqId           string  `json:"seqId"`
	PayType         int     `json:"payType"`
	OrderTitle      string  `json:"orderTitle"`
	OrderAmount     float64 `json:"orderAmount"`
	PayInfo         string  `json:"payInfo"`
	MerchantOrderNo string  `json:"merchantOrderNo"`
}

type CancelPaymentOrderRequest struct {
	OrderNo      string `json:"orderNo"`
	CancelReason string `json:"cancelReason,omitempty"`
}

type CreateRefundRequest struct {
	OrderNo      string  `json:"orderNo,omitempty"`
	ReqSeqId     string  `json:"reqSeqId,omitempty"`
	RefundAmount float64 `json:"refundAmount"`
	RefundReason string  `json:"refundReason,omitempty"`
	Remark       string  `json:"remark,omitempty"`
	NotifyUrl    string  `json:"notifyUrl,omitempty"`
}

type RefundResponse struct {
	MerchantNo        string    `json:"merchantNo"`
	OrderNo           string    `json:"orderNo"`
	SeqId             string    `json:"seqId"`
	ReqDate           string    `json:"reqDate"`
	PaySeqId          string    `json:"paySeqId"`
	PayReqDate        string    `json:"payReqDate"`
	PayUniqueId       string    `json:"payUniqueId"`
	RefundStartDate   string    `json:"refundStartDate"`
	RefundStartTime   time.Time `json:"refundStartTime"`
	RefundFinishTime  time.Time `json:"refundFinishTime"`
	RefundStatus      int       `json:"refundStatus"`
	RefundAmount      float64   `json:"refundAmount"`
	RealRefundAmount  float64   `json:"realRefundAmount"`
	TotalRefAmount    string    `json:"totalRefAmount"`
	TotalRefFeeAmount string    `json:"totalRefFeeAmount"`
	RefCount          string    `json:"refCount"`
}

type QueryRefundRequest struct {
	OrderNo     string `json:"orderNo"`
	RefundSeqId string `json:"refundSeqId,omitempty"`
}

type QueryRefundResponse struct {
	MerchantNo         string  `json:"merchantNo"`
	OrderNo            string  `json:"orderNo"`
	RefundSeqId        string  `json:"refundSeqId"`
	PaySeqId           string  `json:"paySeqId"`
	PayReqDate         string  `json:"payReqDate"`
	RefundAmount       float64 `json:"refundAmount"`
	ActualRefundAmount float64 `json:"actualRefundAmount"`
	RefundStatus       int     `json:"refundStatus"`
	RefundStatusDesc   string  `json:"refundStatusDesc"`
	TransFinishTime    string  `json:"transFinishTime"`
	FeeAmount          float64 `json:"feeAmount"`
	AcctSplitBunch     string  `json:"acctSplitBunch"`
	UnconfirmAmount    float64 `json:"unconfirmAmount"`
	ConfirmedAmount    float64 `json:"confirmedAmount"`
	PayChannel         string  `json:"payChannel"`
	Remark             string  `json:"remark"`
}

type CreateWithdrawRequest struct {
	PayChannel     string  `json:"payChannel"`
	WithdrawAmount float64 `json:"withdrawAmount"`
	ReqSeqId       string  `json:"reqSeqId"`
	Remark         string  `json:"remark,omitempty"`
	NotifyUrl      string  `json:"notifyUrl,omitempty"`
}
