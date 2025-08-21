package models

import "time"

type ExtendedOrder struct {
	Order
	Delivery Delivery `json:"delivery"`
	Payment  Payment  `json:"payment"`
	Items    []*Item  `json:"items"`
}

type Delivery struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	ID           int       `json:"id"`
	Transaction  string    `json:"transaction_id"`
	RequestID    string    `json:"request_id"`
	Currency     string    `json:"currency"`
	Provider     string    `json:"provider"`
	Amount       float64   `json:"amount"`
	PaymentDate  time.Time `json:"payment_date"`
	Bank         string    `json:"bank"`
	DeliveryCost float64   `json:"delivery_cost"`
	GoodsTotal   float64   `json:"goods_total"`
	CustomFee    float32   `json:"custom_fee"`
}

type Order struct {
	ID                int       `json:"id"`
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	DeliveryID        int       `json:"delivery_id"`
	PaymentID         int       `json:"payment_id"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shard_key"`
	SMID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OOFShard          string    `json:"oof_shard"`
}

type Item struct {
	ID          int     `json:"id"`
	OrderID     int     `json:"order_id"`
	ChrtID      int     `json:"chrt_id"`
	TrackNumber string  `json:"track_number"`
	Price       float64 `json:"price"`
	RID         string  `json:"rid"`
	Name        string  `json:"name"`
	Sale        int     `json:"sale"`
	Size        string  `json:"size"`
	TotalPrice  float64 `json:"total_price"`
	NMID        int     `json:"nm_id"`
	Brand       string  `json:"brand"`
	Status      int     `json:"status"`
}
