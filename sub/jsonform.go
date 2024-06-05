package main

type JsonForm struct {
	Order_uid    string `json:"order_uid" validate:"required"`
	Track_number string `json:"track_number" vailidate:"required"`
	Entry        string `json:"entry" validate:"required"`

	Delivery struct {
		Name    string `json:"name" validate:"required"`
		Phone   string `json:"phone" validate:"required,e164"`
		Zip     string `json:"zip" validate:"required"`
		City    string `json:"city" validate:"required"`
		Address string `json:"address" validate:"required"`
		Region  string `json:"region" validate:"required"`
		Email   string `json:"email" validate:"email"`
	} `json:"delivery"`

	Payment struct {
		Transaction   string `json:"transaction" validate:"required"`
		Request_id    string `json:"request_id"`
		Currency      string `json:"currency" validate:"required"`
		Provider      string `json:"provider" validate:"required"`
		Amount        uint64 `json:"amount" validate:"required"`
		Payment_dt    uint64 `json:"payment_dt" validate:"required"`
		Bank          string `json:"bank" validate:"required"`
		Delivery_cost uint64 `json:"delivery_cost"`
		Goods_total   uint64 `json:"goods_total" validate:"required"`
		Custom_fee    uint64 `json:"custom_fee"`
	} `json:"payment"`

	Items []struct {
		Chrt_id      uint64 `json:"chrt_id" validate:"required"`
		Track_number string `json:"track_number" validate:"required"`
		Price        uint64 `json:"price" validate:"required"`
		Rid          string `json:"rid" validate:"required"`
		Name         string `json:"name" validate:"required"`
		Sale         uint64 `json:"sale"`
		Size         string `json:"size"`
		Total_price  uint64 `json:"total_price"`
		Nm_id        uint64 `json:"nm_id" validate:"required"`
		Brand        string `json:"brand"`
		Status       uint64 `json:"status"`
	} `json:"items"`

	Locale             string `json:"locale" validate:"required"`
	Internal_signature string `json:"internal_signature"`
	Customer_id        string `json:"customer_id" vaildate:"required"`
	Delivery_service   string `json:"delivery_service"`
	Shardkey           string `json:"shardkey"`
	Sm_id              uint64 `json:"sm_id"`
	Date_created       string `json:"date_created" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Oof_shard          string `json:"oof_shard"`
}
