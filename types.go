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
	ApplySplit        int     `json:"applySplit,omitempty"`
	ExpireTime        string  `json:"expireTime,omitempty"`
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
	OrderNo      string  `json:"orderNo"`
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
	OrderNo string `json:"orderNo"`
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
	PayChannel       string  `json:"payChannel"`
	WithdrawAmount   float64 `json:"withdrawAmount"`
	IntoAcctDateType string  `json:"intoAcctDateType"`
	Remark           string  `json:"remark,omitempty"`
}

// ApplySplitRequest 商户申请分账请求
type ApplySplitRequest struct {
	OrderNo          string               `json:"orderNo"`
	SplitAmount      float64              `json:"splitAmount"`
	AcctSplitBunches []AcctSplitBunchItem `json:"acctSplitBunches"`
}

// AcctSplitBunchItem 分账配置子项
type AcctSplitBunchItem struct {
	MerchantNo string  `json:"merchantNo"`
	DivAmt     float64 `json:"divAmt"`
}

// ReverseSplitRequest 商户撤销分账请求
type ReverseSplitRequest struct {
	OrderNo      string `json:"orderNo"`
	CancelReason string `json:"cancelReason,omitempty"`
}

// ReverseSplitResponse 撤销分账响应
type ReverseSplitResponse struct {
	OrderNo               string                 `json:"orderNo"`
	TotalRefundAmount     float64                `json:"totalRefundAmount"`
	MerchantRefundDetails []MerchantRefundDetail `json:"merchantRefundDetails"`
}

// MerchantRefundDetail 退款商户明细
type MerchantRefundDetail struct {
	DivAmt     string `json:"divAmt"`
	MerchantNo string `json:"merchantNo"`
}

// QueryHostingOrderRequest 挂靠订单列表查询请求
type QueryHostingOrderRequest struct {
	OrderNo     string `json:"orderNo,omitempty"`
	StartTime   string `json:"startTime,omitempty"`
	EndTime     string `json:"endTime,omitempty"`
	OrderStatus int    `json:"orderStatus,omitempty"`
	PageNum     int    `json:"pageNum,omitempty"`
	PageSize    int    `json:"pageSize,omitempty"`
}

// HostingOrderPageResponse 挂靠订单分页响应
type HostingOrderPageResponse struct {
	Records  []HostingOrderBO `json:"records"`
	Total    int64            `json:"total"`
	PageNum  int              `json:"pageNum"`
	PageSize int              `json:"pageSize"`
}

// HostingOrderBO 挂靠订单业务对象
type HostingOrderBO struct {
	OrderNo     string  `json:"orderNo"`
	MerchantNo  string  `json:"merchantNo"`
	OrderAmount float64 `json:"orderAmount"`
	TransAmt    float64 `json:"transAmt"`
	OrderStatus int     `json:"orderStatus"`
	CreateTime  string  `json:"createTime"`
	UpdateTime  string  `json:"updateTime"`
	ReqSeqId    string  `json:"reqSeqId"`
	TradeType   string  `json:"tradeType"`
	TransStat   string  `json:"transStat"`
	PayAmt      float64 `json:"payAmt"`
	EndTime     string  `json:"endTime"`
}

// QueryOrderDetailRequest 挂靠订单详情查询请求
type QueryOrderDetailRequest struct {
	OrderNo string `json:"orderNo"`
}

// OrderDetailResponse 订单详情响应
type OrderDetailResponse struct {
	HostOrderInfo       *HostOrderInfo        `json:"hostOrderInfo"`
	PayInfo             *PayInfoDetail        `json:"payInfo"`
	TradeConfirmRecord  *TradeConfirmRecord   `json:"tradeConfirmRecord"`
	RefundInfo          *RefundInfoDetail     `json:"refundInfo"`
	RefundDetails       []RefundDetailItem    `json:"refundDetails"`
	SplitSettlementList []SplitSettlementItem `json:"splitSettlementList"`
}

// HostOrderInfo 挂靠订单详情
type HostOrderInfo struct {
	ID              int64   `json:"id"`
	OrderNo         string  `json:"orderNo"`
	MerchantNo      string  `json:"merchantNo"`
	OrderAmount     float64 `json:"orderAmount"`
	OrderStatus     int     `json:"orderStatus"`
	OrderStatusDesc string  `json:"orderStatusDesc"`
	OrderTitle      string  `json:"orderTitle"`
	PayAmount       float64 `json:"payAmount"`
	FinishTime      string  `json:"finishTime"`
	PayChannel      string  `json:"payChannel"`
	PayChannelDesc  string  `json:"payChannelDesc"`
	PayType         int     `json:"payType"`
	PayTypeDesc     string  `json:"payTypeDesc"`
	CreateTime      string  `json:"createTime"`
	UpdateTime      string  `json:"updateTime"`
	RedirectUrl     string  `json:"redirectUrl"`
}

// PayInfoDetail 支付订单详情
type PayInfoDetail struct {
	OrderNo              string      `json:"orderNo"`
	ChannelUniqueOrderNo string      `json:"channelUniqueOrderNo"`
	PartyOrderId         string      `json:"partyOrderId"`
	MerchantNo           string      `json:"merchantNo"`
	OrderAmount          float64     `json:"orderAmount"`
	PayStatus            int         `json:"payStatus"`
	PayStatusDesc        string      `json:"payStatusDesc"`
	OrderTitle           string      `json:"orderTitle"`
	PayAmount            float64     `json:"payAmount"`
	FinishTime           string      `json:"finishTime"`
	PayChannel           string      `json:"payChannel"`
	PayChannelDesc       string      `json:"payChannelDesc"`
	PayType              int         `json:"payType"`
	PayTypeDesc          string      `json:"payTypeDesc"`
	CreateTime           string      `json:"createTime"`
	UpdateTime           string      `json:"updateTime"`
	RedirectUrl          string      `json:"redirectUrl"`
	RefundedAmount       float64     `json:"refundedAmount"`
	RefundableAmount     float64     `json:"refundableAmount"`
	RefundTime           string      `json:"refundTime"`
	SettleInfo           *SettleInfo `json:"settleInfo"`
}

// SettleInfo 结算信息
type SettleInfo struct {
	EstimatedSettlementTime string  `json:"estimatedSettlementTime"`
	RealSettlementTime      string  `json:"realSettlementTime"`
	SettlementStatus        int     `json:"settlementStatus"`
	SettlementStatusDesc    string  `json:"settlementStatusDesc"`
	MerchantIncome          float64 `json:"merchantIncome"`
	PayFee                  float64 `json:"payFee"`
	SettleAmount            float64 `json:"settleAmount"`
}

// TradeConfirmRecord 交易确认记录
type TradeConfirmRecord struct {
	SplitAmount      float64 `json:"splitAmount"`
	PayAmount        float64 `json:"payAmount"`
	RefundAmount     float64 `json:"refundAmount"`
	RefundableAmount float64 `json:"refundableAmount"`
	Settled          bool    `json:"settled"`
	SettlementTime   string  `json:"settlementTime"`
	PayFee           float64 `json:"payFee"`
}

// RefundInfoDetail 退款主单详情
type RefundInfoDetail struct {
	MerchantNo        string  `json:"merchantNo"`
	OrderNo           string  `json:"orderNo"`
	SeqId             string  `json:"seqId"`
	ReqDate           string  `json:"reqDate"`
	PaySeqId          string  `json:"paySeqId"`
	PayReqDate        string  `json:"payReqDate"`
	PayUniqueId       string  `json:"payUniqueId"`
	RefundStartDate   string  `json:"refundStartDate"`
	RefundStartTime   string  `json:"refundStartTime"`
	RefundFinishTime  string  `json:"refundFinishTime"`
	RefundStatus      int     `json:"refundStatus"`
	RefundAmount      float64 `json:"refundAmount"`
	RealRefundAmount  float64 `json:"realRefundAmount"`
	TotalRefAmount    string  `json:"totalRefAmount"`
	TotalRefFeeAmount string  `json:"totalRefFeeAmount"`
	RefCount          string  `json:"refCount"`
}

// RefundDetailItem 退款子单明细
type RefundDetailItem struct {
	MerchantNo      string  `json:"merchantNo"`
	OrderNo         string  `json:"orderNo"`
	ReqDate         string  `json:"reqDate"`
	ReqSeqId        string  `json:"reqSeqId"`
	OrgReqDate      string  `json:"orgReqDate"`
	OrgReqId        string  `json:"orgReqId"`
	TransDate       string  `json:"transDate"`
	TransTime       string  `json:"transTime"`
	TransFinishTime string  `json:"transFinishTime"`
	TransStat       string  `json:"transStat"`
	OrdAmt          float64 `json:"ordAmt"`
	ActualRefAmt    float64 `json:"actualRefAmt"`
	UnconfirmAmt    float64 `json:"unconfirmAmt"`
	FundFreezeStat  string  `json:"fundFreezeStat"`
	PayChannel      string  `json:"payChannel"`
	TotalRefAmt     float64 `json:"totalRefAmt"`
	TotalRefFeeAmt  float64 `json:"totalRefFeeAmt"`
	RefCut          string  `json:"refCut"`
	NotifyUrl       string  `json:"notifyUrl"`
	CofRfdReqId     string  `json:"cofRfdReqId"`
}

// SplitSettlementItem 分账结算子项
type SplitSettlementItem struct {
	MerchantNo string  `json:"merchantNo"`
	DivAmount  float64 `json:"divAmount"`
}
