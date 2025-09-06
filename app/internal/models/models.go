package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(modelsStruct interface{}) error {
	return validate.Struct(modelsStruct)
}

type ExtendedOrder struct {
	Order
	Delivery Delivery `json:"delivery" validate:"required"`
	Payment  Payment  `json:"payment" validate:"required"`
	Items    []*Item  `json:"items" validate:"required,min=1,dive,required"`
}

type Delivery struct {
	ID      int64  `json:"id"`
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required,e164"` // e164 -> +734264
	Zip     string `json:"zip" validate:"required"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

type Payment struct {
	ID           int64   `json:"id"`
	Transaction  string  `json:"transaction" validate:"required"`
	RequestID    string  `json:"request_id" validate:"omitempty"`
	Currency     string  `json:"currency" validate:"required,len=3,uppercase"`
	Provider     string  `json:"provider" validate:"required"`
	Amount       float64 `json:"amount" validate:"required,gt=0"`
	PaymentDate  int64   `json:"payment_dt" validate:"required"`
	Bank         string  `json:"bank" validate:"required"`
	DeliveryCost float64 `json:"delivery_cost" validate:"required,gt=0"`
	GoodsTotal   float64 `json:"goods_total" validate:"required,gt=0"`
	CustomFee    float64 `json:"custom_fee"`
}

type Order struct {
	ID                int64     `json:"id"`
	OrderUID          string    `json:"order_uid" validate:"required"`
	TrackNumber       string    `json:"track_number" validate:"required"`
	Entry             string    `json:"entry" validate:"required"`
	DeliveryID        int64     `json:"delivery_id"`
	PaymentID         int64     `json:"payment_id"`
	Locale            string    `json:"locale" validate:"required"`
	InternalSignature string    `json:"internal_signature" validate:"omitempty"`
	CustomerID        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" validate:"required"`
	ShardKey          string    `json:"shardkey" validate:"required"`
	SMID              int       `json:"sm_id" validate:"required"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OOFShard          string    `json:"oof_shard" validate:"required"`
}

type Item struct {
	ID          int64   `json:"id"`
	OrderID     int64   `json:"order_id"`
	ChrtID      int     `json:"chrt_id" validate:"required"`
	TrackNumber string  `json:"track_number" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	RID         string  `json:"rid" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Sale        int     `json:"sale" validate:"required,gte=0,lt=100"`
	Size        string  `json:"size" validate:"required"`
	TotalPrice  float64 `json:"total_price" validate:"required,gt=0"`
	NMID        int     `json:"nm_id" validate:"required"`
	Brand       string  `json:"brand" validate:"required"`
	Status      int     `json:"status" validate:"required"`
}
