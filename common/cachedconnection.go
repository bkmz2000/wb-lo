package common

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/nats-io/nats.go"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "orders"
)

type CachedConnection struct {
	db    *sql.DB
	cache map[string]Order
}

func (c *CachedConnection) Connect() error {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)

	if err != nil {
		return err
	}

	err = db.Ping()

	if err != nil {
		return err
	}

	c.db = db

	c.ReCache()

	return nil
}

func (c *CachedConnection) ReCache() error {
	c.cache = make(map[string]Order)

	rows, err := c.db.Query("SELECT order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders;")

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		order, err := NewOrderFromRows(rows)

		if err != nil {
			return err
		}

		c.cache[order.OrderUID] = order
	}

	return nil
}

func (c *CachedConnection) Get(uid string) (Order, error) {
	order, ok := c.cache[uid]

	if ok {
		return order, nil
	}

	rows, err := c.db.Query("SELECT order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard WHERE uid=$1 FROM orders;", uid)

	if err != nil {
		c.checkConnection()
		return Order{}, err
	}

	defer rows.Close()

	order, err = NewOrderFromRows(rows)

	if err == nil {
		c.cache[uid] = order
	}

	return order, err // uid is unique
}

func (c *CachedConnection) Insert(jsonOrder string) error {
	var order Order

	err := json.Unmarshal([]byte(jsonOrder), &order)

	if err != nil {
		return err
	}

	q, err := c.db.Prepare(`
		INSERT INTO 
			ORDERS(
				order_uid,
				track_number, 
				entry, 
				delivery, 
				payment, 
				items, 
				locale, 
				internal_signature, 
				customer_id, 
				delivery_service, 
				shardkey, 
				sm_id, 
				date_created, 
				oof_shard)
			VALUES(
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8,
				$9,
				$10,
				$11,
				$12,
				$13,
				$14);`)

	if err != nil {
		return err
	}

	_, err = q.Exec(
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Delivery,
		order.Payment,
		order.Items,
		order.Locale,
		order.InternalSignature,
		order.CustomerId,
		order.DeliveryService,
		order.ShardKey,
		order.SMID,
		order.DateCreated,
		order.OofShard,
	)

	if err == nil {
		c.cache[order.OrderUID] = order
	}

	return err
}

func (c *CachedConnection) Close() error {
	return c.db.Close()
}

func (c *CachedConnection) HTTPGetter(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	uid := string(body)

	order, err := c.Get(uid)

	if err != nil {
		http.Error(w, "Error retrieving order", http.StatusInternalServerError)
		return
	}

	response, _ := json.MarshalIndent(order, "", "  ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (c *CachedConnection) checkConnection() {
	err := c.db.Ping()
	if err != nil {
		c.Connect()
	}
}

func (c *CachedConnection) NATSInserter(msg *nats.Msg) {
	err := c.Insert(string(msg.Data))

	if err != nil {
		c.checkConnection()
		log.Printf("Error while inserting via NATS: %v", err)
	} else {
		log.Println("Succesfully inserted!")
		log.Println(string(msg.Data))
		log.Println("\n\n\n\n")
	}
}
