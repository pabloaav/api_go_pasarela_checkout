package linkqr

import "time"

type EstadoOperacion struct {
	DateCreated time.Time    `json:"date_created"`
	ID          string       `json:"id"`
	Action      string       `json:"action"`
	Type        string       `json:"type"`
	LiveMode    bool         `json:"live_mode"`
	APIVersion  string       `json:"api_version"`
	Data        TransferData `json:"data"`
}

type TransferData struct {
	OperationID string      `json:"operation_id"`
	Amount      Amount      `json:"amount"`
	Merchant    Merchant    `json:"merchant"`
	Payer       Payer       `json:"payer"`
	Transaction Transaction `json:"transaction"`
}

type Amount struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

type Merchant struct {
	SoftDescriptor string `json:"soft_descriptor"`
	CUIT           string `json:"cuit"`
	BranchCode     string `json:"branch_code"`
	PosCode        string `json:"pos_code"`
}

type Payer struct {
	Name     string   `json:"name"`
	Document Document `json:"document"`
	Wallet   Wallet   `json:"wallet"`
}

type Document struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}

type Wallet struct {
	Name string `json:"name"`
	CUIT string `json:"cuit"`
}

type Transaction struct {
	AuthorizationCode string `json:"authorization_code"`
	Code              string `json:"code"`
	Description       string `json:"description"`
	OnRejection       struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"on_rejection"`
	Datetime          time.Time         `json:"datetime"`
	Type              string            `json:"type"`
	GrossAmount       Amount            `json:"gross_amount"`
	NetAmount         Amount            `json:"net_amount"`
	MerchantDiscounts MerchantDiscounts `json:"merchant_discounts"`
}

type MerchantDiscounts struct {
	TotalDiscounts        Amount `json:"total_discounts"`
	WithholdingsDiscounts Amount `json:"withholdings_discounts"`
	MerchantDiscountRate  struct {
		PercentageOverGrossAmount int    `json:"percentage_over_gross_amount"`
		Amount                    Amount `json:"amount"`
		IVAPercentage             int    `json:"iva_percentage"`
		IVAAmount                 Amount `json:"iva_amount"`
		Total                     Amount `json:"total"`
	} `json:"merchant_discount_rate"`
}
