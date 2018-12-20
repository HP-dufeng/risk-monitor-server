package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"

	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.Stamp,
	})

	session, err := r.Connect(r.ConnectOpts{
		// Address:  "websocket-rethinkdb-rma-7x24.apps.dev-cefcfco.com:30028", // endpoint without http
		Address:  "websocket-rethinkdb:28015", // endpoint without http
		Database: DbName,
	})
	if err != nil {
		log.Fatal(err)
	}

	go Subscribe(session, TableName_SubscribeTunnelRealFund)
	go Subscribe(session, TableName_SubscribeCorpHoldMon)
	go Subscribe(session, TableName_SubscribeQuoteMon)
	go Subscribe(session, TableName_SubscribeCustRisk)
	go Subscribe(session, TableName_SubscribeCustHold)
	go Subscribe(session, TableName_SubscribeCustGroupHold)
	go Subscribe(session, TableName_SubscribeProuctGroupRisk)
	go Subscribe(session, TableName_SubscribeNearDediveHold)

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
}

var (
	DbName                             = "rma_subscribes"
	TableName_SubscribeTunnelRealFund  = "SubscribeTunnelRealFund"
	TableName_SubscribeCorpHoldMon     = "SubscribeCorpHoldMon"
	TableName_SubscribeQuoteMon        = "SubscribeQuoteMon"
	TableName_SubscribeCustRisk        = "SubscribeCustRisk"
	TableName_SubscribeCustHold        = "SubscribeCustHold"
	TableName_SubscribeCustGroupHold   = "SubscribeCustGroupHold"
	TableName_SubscribeProuctGroupRisk = "SubscribeProuctGroupRisk"
	TableName_SubscribeNearDediveHold  = "SubscribeNearDediveHold"
)

func Subscribe(session *r.Session, tableName string) {
	res, err := r.Table(tableName).Changes().Run(session)
	if err != nil {
		log.Fatalln(err)
	}

	var value interface{}

	log.Infof("%s received start", tableName)
	for res.Next(&value) {
		log.Info(value)
	}

}
