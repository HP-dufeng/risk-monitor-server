package main

import (
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
		Address:  "rethinkdb-rma-dc-test.apps.cefcfco.com:30018", // endpoint without http
		Database: DbName,
	})
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool, 1)

	go Subscribe(session, TableName_SubscribeTunnelRealFund)
	go Subscribe(session, TableName_SubscribeCorpHoldMon)
	go Subscribe(session, TableName_SubscribeQuoteMon)
	go Subscribe(session, TableName_SubscribeCustRisk)
	go Subscribe(session, TableName_SubscribeCustHold)
	go Subscribe(session, TableName_SubscribeCustGroupHold)
	go Subscribe(session, TableName_SubscribeProuctGroupRisk)
	go Subscribe(session, TableName_SubscribeNearDediveHold)

	<-done
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
