package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	stan "github.com/nats-io/stan.go"
	"jagodkiL0/rest"
	"jagodkiL0/store"
	"log"
	"os"
	"path/filepath"
	"time"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func main() {

	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)
	dbConn := "postgresql://localhost:5432/dbwb?sslmode=disable&user=postgres&password=741985"

	migrationPth := fmt.Sprintf("file:%s", filepath.Join(exPath+"/migrations"))
	db, err := store.NewStore(dbConn, migrationPth)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	h := rest.NewAPI(&db)

	subject := "test"
	sc, err := stan.Connect("test-cluster", "client-123")
	if err != nil {
		panic(err)
	}

	sub, err := sc.Subscribe(subject, func(m *stan.Msg) {
		o := store.Order{}
		data := m.Data
		err = json.Unmarshal(data, &o)
		if err != nil {
			log.Println(err)
			return
		}
		err = db.AddOrder(o)
		if err != nil {
			log.Println(err)
			return
		}
	})
	if err != nil {
		panic(err)
	}

	o := store.Order{}
o.OrderUid = "test3"
o.Items = []store.Item{{
2,
"test",
12,
"tracknum",
5.0,
"rid_test",
"item_name",
2.0,
"size",
3.0,
126,
"str",
457,
}}
o.Delivery = store.Delivery{
OrderId: "test",
Name:    "del_name",
Phone:   "del_phone",
Zip:     "125",
City:    "459",
Address: "678",
Region:  "123123",
Email:   "email",
}
o.Payment = store.Payment{
OrderId:      "test",
Transaction:  "transac",
RequestId:    "reqid",
Currency:     "cur",
Provider:     "provid",
Amount:       123,
PaymentDt:    1,
Bank:         "s",
DeliveryCost: 1.0,
GoodsTotal:   2,
CustomFee:    0.5,
}
o.CustomerId = "cusid"
o.DateCreated = time.Now()
o.DeliveryService = "or_entry"
o.Entry = "or_entry"
o.InternalSignature = "intsig"
o.Locate = "lo"
o.OofShard = "ofsh"
o.ShardKey = "shkey"
o.SmId = 2
o.TrackNumber = "tracknumbe"
b, _ := json.Marshal(o)
err = sc.Publish(subject, []byte(b))
	if err != nil {
		panic(err)
	}

	defer sub.Close()

	e := gin.Default()

	e.Use(CORSMiddleware())
	e.Use(gin.Logger())
	e.Use(gin.Recovery())

	e.Static("/frontend", "./frontend")
	api := e.Group("/api")
	api.GET("/order/:order_uid", h.GetOrder)
	e.Run(":8080")
}
