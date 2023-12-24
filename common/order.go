package common

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

// Payment represents the "payment" object within the JSON
type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type Order struct {
	OrderUID          string    `json:"order_uid" sql:"order_uid"`                   // SQL: VARCHAR(255) PRIMARY KEY
	TrackNumber       string    `json:"track_number" sql:"track_number"`             // SQL: VARCHAR(255)
	Entry             string    `json:"entry" sql:"entry"`                           // SQL: VARCHAR(255)
	Delivery          string    `json:"delivery" sql:"delivery"`                     // SQL: JSON
	Payment           string    `json:"payment" sql:"payment"`                       // SQL: JSON
	Items             string    `json:"items" sql:"items"`                           // SQL: JSON
	Locale            string    `json:"locale" sql:"locale"`                         // SQL: VARCHAR(10)
	InternalSignature string    `json:"internal_signature" sql:"internal_signature"` // SQL: VARCHAR(255)
	CustomerId        string    `json:"customer_id" sql:"customer_id"`               // SQL: VARCHAR(255)
	DeliveryService   string    `json:"delivery_service" sql:"delivery_service"`     // SQL: VARCHAR(255)
	ShardKey          string    `json:"shardkey" sql:"shardkey"`                     // SQL: VARCHAR(10)
	SMID              int       `json:"sm_id" sql:"sm_id"`                           // SQL: INTEGER
	DateCreated       time.Time `json:"date_created" sql:"date_created"`             // SQL: TIMESTAMP
	OofShard          string    `json:"oof_shard" sql:"oof_shard"`                   // SQL: VARCHAR(10)
}

func NewOrderFromRows(rows *sql.Rows) (Order, error) {
	var order Order
	err := rows.Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Delivery,
		&order.Payment,
		&order.Items,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerId,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SMID,
		&order.DateCreated,
		&order.OofShard,
	)

	if err != nil {
		return Order{}, err
	}

	return order, nil
}

func (o Order) String() string {
	// Convert DateCreated to string representation
	dateCreatedStr := o.DateCreated.Format("2006-01-02 15:04:05")

	// Convert struct fields to a map for JSON formatting
	orderMap := map[string]interface{}{
		"order_uid":          o.OrderUID,
		"track_number":       o.TrackNumber,
		"entry":              o.Entry,
		"delivery":           o.Delivery,
		"payment":            o.Payment,
		"items":              o.Items,
		"locale":             o.Locale,
		"internal_signature": o.InternalSignature,
		"customer":           o.CustomerId,
		"delivery_service":   o.DeliveryService,
		"shardkey":           o.ShardKey,
		"sm_id":              o.SMID,
		"date_created":       dateCreatedStr,
		"oof_shard":          o.OofShard,
	}

	jsonStr, err := json.MarshalIndent(orderMap, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error converting to JSON: %v", err)
	}

	return string(jsonStr)
}

func generateRandomTime() time.Time {
	start := time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC)
	end := time.Now()

	duration := end.Sub(start)

	dt := time.Duration(rand.Int63n(int64(duration)))

	randomTime := start.Add(dt)

	return randomTime
}

func GenerateRandomOrderJSON() string {
	delivery, _ := json.Marshal(generateRandomDelivery())
	payment, _ := json.Marshal(generateRandomPayment())
	items, _ := json.Marshal(generateRandomItems())

	orderMap := map[string]interface{}{
		"order_uid":          generateRandomString(10),
		"track_number":       generateRandomString(12),
		"entry":              generateRandomString(5),
		"delivery":           string(delivery),
		"payment":            string(payment),
		"items":              string(items),
		"locale":             generateRandomString(2),
		"internal_signature": generateRandomString(20),
		"customer_id":        generateRandomString(8),
		"delivery_service":   generateRandomString(10),
		"shardkey":           generateRandomString(3),
		"sm_id":              rand.Intn(100),
		"date_created":       time.Now().UTC(),
		"oof_shard":          generateRandomString(3),
	}

	jsonData, _ := json.MarshalIndent(orderMap, "", "  ")

	return string(jsonData)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generateRandomDelivery() map[string]interface{} {
	name, surname := generateRandomString(10), generateRandomString(10)
	fullname := fmt.Sprintf("%s %s", name, surname)

	street, building := generateRandomString(10), rand.Intn(100)
	address := fmt.Sprintf("%s %d", street, building)
	number := generateRandomNumber()
	email := fmt.Sprintf("%s@gmail.com", generateRandomString(10))

	delivery := map[string]interface{}{
		"name":    fullname,
		"phone":   number,
		"zip":     fmt.Sprint(rand.Intn(1000000)),
		"city":    generateRandomString(10),
		"address": address,
		"region":  generateRandomString(10),
		"email":   email,
	}

	return delivery
}

func generateRandomNumber() string {
	const charset = "0123456789"
	b := make([]byte, 11)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	b[0] = '+'

	return string(b)
}

func generateRandomItems() []map[string]interface{} {
	count := rand.Intn(9) + 1

	items := make([]map[string]interface{}, count)

	for i := range items {
		items[i] = generateRandomItem()
	}

	return items
}

func generateRandomItem() map[string]interface{} {
	itemMap := map[string]interface{}{
		"chrt_id":      rand.Intn(10000),
		"track_number": generateRandomString(10),
		"price":        rand.Intn(10000),
		"rid":          generateRandomString(10),
		"name":         generateRandomString(10),
		"sale":         rand.Intn(10000),
		"size":         generateRandomString(10),
		"total_price":  rand.Intn(10000),
		"nm_id":        rand.Intn(10000),
		"brand":        generateRandomString(10),
		"status":       rand.Intn(10000),
	}

	return itemMap
}

func generateRandomPayment() map[string]interface{} {
	paymentMap := map[string]interface{}{
		"transaction":   generateRandomString(10),
		"request_id":    generateRandomString(10),
		"currency":      generateRandomString(10),
		"provider":      generateRandomString(10),
		"amount":        rand.Intn(10000),
		"payment_dt":    rand.Int63n(10000),
		"bank":          generateRandomString(10),
		"delivery_cost": rand.Intn(10000),
		"goods_total":   rand.Intn(10000),
		"custom_fee":    rand.Intn(10000),
	}

	return paymentMap
}
