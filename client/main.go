package main

import (
	"flag"
	"time"

	"github.com/fengdu/risk-monitor-server/client/core/subscribers"
	log "github.com/sirupsen/logrus"

	pb "github.com/fengdu/risk-monitor-server/pb"
	"google.golang.org/grpc"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.Stamp,
	})

	flag.Parse()
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewRiskMonitorServerClient(conn)

	// r := core.NewRethinkDB(rethinkdbAddr)
	// session, err := r.Init()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// CreateDb(session)
	// session.Use(DbName)

	// CreateTable(session, TableName_SubscribeTunnelRealFund)
	// CreateTable(session, TableName_SubscribeCorpHoldMon)
	// CreateTable(session, TableName_SubscribeQuoteMon)
	// CreateTable(session, TableName_SubscribeCustRisk)
	// CreateTable(session, TableName_SubscribeCustHold)
	// CreateTable(session, TableName_SubscribeCustGroupHold)
	// CreateTable(session, TableName_SubscribeProuctGroupRisk)
	// CreateTable(session, TableName_SubscribeNearDediveHold)

	done := make(chan struct{})
	defer close(done)

	t := subscribers.NewTunnelRealFundSubscriber(client)
	s, errSubscribe := t.Subscribe(done)
	c := t.Convert(done, s)
	out, errEnqueue := t.Enqueue(done, c)
	// go SubscribeQuoteMon(client, session)
	// go SubscribeTunnelRealFund(client, session)
	// go SubscribeCorpHoldMon(client, session)
	// go SubscribeCustRisk(client, session)
	// go SubscribeCustHold(client, session)
	// go SubscribeCustGroupHold(client, session)
	// go SubscribeNearDediveHold(client, session)
	// go SubscribeProuctGroupRisk(client, session)

complete:
	for {
		select {
		case item := <-out:
			log.Infof("%+v", item)
		case err := <-errSubscribe:
			log.Error(err)
			break complete
		case err := <-errEnqueue:
			log.Error(err)
			break complete
		}
	}

	log.Info("completed")

}

var (
	serverAddr = flag.String("server_addr", "10.1.7.127:8081", "The server address in the format of host:port")

	rethinkdbAddr = "rethinkdb-rma-dc-test.apps.cefcfco.com:30018"
	deadline      = 7 * time.Duration(24) * time.Hour

	// DbName                             = "rma_subscribes"
	// TableName_SubscribeTunnelRealFund  = "SubscribeTunnelRealFund"
	// TableName_SubscribeCorpHoldMon     = "SubscribeCorpHoldMon"
	// TableName_SubscribeQuoteMon        = "SubscribeQuoteMon"
	// TableName_SubscribeCustRisk        = "SubscribeCustRisk"
	// TableName_SubscribeCustHold        = "SubscribeCustHold"
	// TableName_SubscribeCustGroupHold   = "SubscribeCustGroupHold"
	// TableName_SubscribeProuctGroupRisk = "SubscribeProuctGroupRisk"
	// TableName_SubscribeNearDediveHold  = "SubscribeNearDediveHold"
)

// // func CreateDb(session *r.Session) {
// // 	res, err := r.DBList().Run(session)
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}
// // 	defer res.Close()

// // 	var rows []string
// // 	err = res.All(&rows)
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}

// // 	existed := false
// // 	for _, d := range rows {
// // 		if d == DbName {
// // 			existed = true
// // 			break
// // 		}
// // 	}
// // 	if !existed {
// // 		_, err = r.DBCreate(DbName).Run(session)
// // 		if err != nil {
// // 			log.Fatal(err)
// // 		}
// // 		log.Infof("DB %s created.", DbName)

// // 	} else {
// // 		log.Infof("DB %s existed.", DbName)
// // 	}

// // }

// // func CreateTable(session *r.Session, tableName string) {
// // 	res, err := r.TableList().Run(session)
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}
// // 	defer res.Close()

// // 	var rows []string
// // 	err = res.All(&rows)
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}

// // 	existed := false
// // 	for _, t := range rows {
// // 		if t == tableName {
// // 			existed = true
// // 			break
// // 		}
// // 	}
// // 	if !existed {
// // 		_, err = r.TableCreate(tableName, r.TableCreateOpts{PrimaryKey: "ActionKey"}).Run(session)
// // 		if err != nil {
// // 			log.Fatal(err)
// // 		}
// // 		log.Infof("Table %s created.", tableName)

// // 	} else {
// // 		log.Infof("Table %s existed.", tableName)
// // 	}
// // }

// // SubscribeQuoteMon lists all the Quotes.
// type RtnDto struct {
// 	TableName  string
// 	session    *r.Session
// 	ActionFlag uint32
// 	ActionKey  string
// 	Dto        interface{}
// }

// func (item *RtnDto) Replace() {

// 	if item.ActionFlag != 2 {
// 		var m map[string]interface{}
// 		j, _ := json.Marshal(item.Dto)
// 		json.Unmarshal(j, &m)
// 		r.Table(item.TableName).Get(item.ActionKey).Replace(m).RunWrite(item.session)
// 	} else {
// 		r.Table(item.TableName).Get(item.ActionKey).Replace(nil).RunWrite(item.session)
// 	}

// }

// type TunnelRealFundRtnDto struct {
// 	ActionKey string `json:"ActionKey,omitempty"`
// 	pb.TunnelRealFundRtn
// }

// func SubscribeTunnelRealFund(client pb.RiskMonitorServerClient, session *r.Session) {
// 	log.Info("Looking for TunnelRealFund within ")

// 	clientDeadline := time.Now().Add(deadline)
// 	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
// 	defer cancel()
// 	stream, err := client.SubscribeTunnelRealFund(ctx, &pb.SubscribeReq{})
// 	if err != nil {
// 		log.Fatal("%v.TunnelRealFund(_) = _, %v", client, err)
// 	}
// 	for {
// 		item, err := stream.Recv()
// 		if err == io.EOF {
// 			log.Info("EOF")
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal("%v.TunnelRealFund(_) = _, %v", client, err)
// 		}

// 		d := &TunnelRealFundRtnDto{
// 			ActionKey:         string(item.MonitorNo) + "#" + string(item.TunnelCode) + "#" + string(item.CurrencyCode),
// 			TunnelRealFundRtn: *item,
// 		}

// 		rethink := &RtnDto{
// 			TableName:  TableName_SubscribeTunnelRealFund,
// 			ActionFlag: d.TunnelRealFundRtn.ActionFlag,
// 			ActionKey:  d.ActionKey,
// 			session:    session,
// 			Dto:        d,
// 		}
// 		rethink.Replace()
// 	}

// }

// type CorpHoldMonRtnDto struct {
// 	ActionKey string `json:"ActionKey,omitempty"`
// 	pb.CorpHoldMonRtn
// }

// func SubscribeCorpHoldMon(client pb.RiskMonitorServerClient, session *r.Session) {
// 	log.Info("Looking for CorpHoldMon within ")

// 	clientDeadline := time.Now().Add(deadline)
// 	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
// 	defer cancel()
// 	stream, err := client.SubscribeCorpHoldMon(ctx, &pb.SubscribeReq{})
// 	if err != nil {
// 		log.Fatal("%v.CorpHoldMon(_) = _, %v", client, err)
// 	}
// 	for {
// 		item, err := stream.Recv()
// 		if err == io.EOF {
// 			log.Info("EOF")
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal("%v.CorpHoldMon(_) = _, %v", client, err)
// 		}

// 		d := &CorpHoldMonRtnDto{
// 			ActionKey:      string(item.MonitorNo) + "#" + string(item.ContractCode),
// 			CorpHoldMonRtn: *item,
// 		}

// 		rethink := &RtnDto{
// 			TableName:  TableName_SubscribeCorpHoldMon,
// 			ActionFlag: d.CorpHoldMonRtn.ActionFlag,
// 			ActionKey:  d.ActionKey,
// 			session:    session,
// 			Dto:        d,
// 		}
// 		rethink.Replace()
// 	}

// }

// type QuoteMonRtnDto struct {
// 	ActionKey string `json:"ActionKey,omitempty"`
// 	pb.QuoteMonRtn
// }

// func SubscribeQuoteMon(client pb.RiskMonitorServerClient, session *r.Session) {
// 	log.Info("Looking for Quotes within ")

// 	clientDeadline := time.Now().Add(deadline)
// 	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
// 	defer cancel()
// 	stream, err := client.SubscribeQuoteMon(ctx, &pb.SubscribeReq{})
// 	if err != nil {
// 		log.Fatal("%v.QryQuote(_) = _, %v", client, err)
// 	}
// 	for {
// 		item, err := stream.Recv()
// 		if err == io.EOF {
// 			log.Info("EOF")
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal("%v.QryQuote(_) = _, %v", client, err)
// 		}

// 		d := &QuoteMonRtnDto{
// 			ActionKey:   string(item.MonitorNo) + "#" + string(item.ContractCode),
// 			QuoteMonRtn: *item,
// 		}

// 		rethink := &RtnDto{
// 			TableName:  TableName_SubscribeQuoteMon,
// 			ActionFlag: d.QuoteMonRtn.ActionFlag,
// 			ActionKey:  d.ActionKey,
// 			session:    session,
// 			Dto:        d,
// 		}
// 		rethink.Replace()
// 	}

// }

// type CustRiskRtnDto struct {
// 	ActionKey string `json:"ActionKey,omitempty"`
// 	pb.CustRiskRtn
// }

// func SubscribeCustRisk(client pb.RiskMonitorServerClient, session *r.Session) {
// 	log.Info("Looking for CustRisk within ")

// 	clientDeadline := time.Now().Add(deadline)
// 	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
// 	defer cancel()
// 	stream, err := client.SubscribeCustRisk(ctx, &pb.SubscribeReq{})
// 	if err != nil {
// 		log.Fatal("%v.CustRisk(_) = _, %v", client, err)
// 	}
// 	for {
// 		item, err := stream.Recv()
// 		if err == io.EOF {
// 			log.Info("EOF")
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal("%v.CustRisk(_) = _, %v", client, err)
// 		}

// 		d := &CustRiskRtnDto{
// 			ActionKey:   string(item.MonitorNo) + "#" + string(item.CustNo) + "#" + string(item.CurrencyCode),
// 			CustRiskRtn: *item,
// 		}

// 		rethink := &RtnDto{
// 			TableName:  TableName_SubscribeCustRisk,
// 			ActionFlag: d.CustRiskRtn.ActionFlag,
// 			ActionKey:  d.ActionKey,
// 			session:    session,
// 			Dto:        d,
// 		}
// 		rethink.Replace()
// 	}

// }

// type CustHoldRtnDto struct {
// 	ActionKey string `json:"ActionKey,omitempty"`
// 	pb.CustHoldRtn
// }

// func SubscribeCustHold(client pb.RiskMonitorServerClient, session *r.Session) {
// 	log.Info("Looking for CustHold within ")

// 	clientDeadline := time.Now().Add(deadline)
// 	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
// 	defer cancel()
// 	stream, err := client.SubscribeCustHold(ctx, &pb.SubscribeReq{})
// 	if err != nil {
// 		log.Fatal("%v.CustHold(_) = _, %v", client, err)
// 	}
// 	for {
// 		item, err := stream.Recv()
// 		if err == io.EOF {
// 			log.Info("EOF")
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal("%v.CustHold(_) = _, %v", client, err)
// 		}

// 		d := &CustHoldRtnDto{
// 			ActionKey:   string(item.MonitorNo) + "#" + string(item.CustNo) + "#" + string(item.ContractCode) + "#" + strconv.Itoa(int(item.HoldType)),
// 			CustHoldRtn: *item,
// 		}

// 		rethink := &RtnDto{
// 			TableName:  TableName_SubscribeCustHold,
// 			ActionFlag: d.CustHoldRtn.ActionFlag,
// 			ActionKey:  d.ActionKey,
// 			session:    session,
// 			Dto:        d,
// 		}
// 		rethink.Replace()
// 	}

// }

// type CustGroupHoldRtnDto struct {
// 	ActionKey string `json:"ActionKey,omitempty"`
// 	pb.CustGroupHoldRtn
// }

// func SubscribeCustGroupHold(client pb.RiskMonitorServerClient, session *r.Session) {
// 	log.Info("Looking for CustGroupHold within ")

// 	clientDeadline := time.Now().Add(deadline)
// 	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
// 	defer cancel()
// 	stream, err := client.SubscribeCustGroupHold(ctx, &pb.SubscribeReq{})
// 	if err != nil {
// 		log.Fatal("%v.CustGroupHold(_) = _, %v", client, err)
// 	}
// 	for {
// 		item, err := stream.Recv()
// 		if err == io.EOF {
// 			log.Info("EOF")
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal("%v.CustGroupHold(_) = _, %v", client, err)
// 		}

// 		d := &CustGroupHoldRtnDto{
// 			ActionKey:        string(item.MonitorNo) + "#" + string(item.CustGroupNo) + "#" + string(item.ContractCode) + "#" + strconv.Itoa(int(item.HoldType)),
// 			CustGroupHoldRtn: *item,
// 		}

// 		rethink := &RtnDto{
// 			TableName:  TableName_SubscribeCustGroupHold,
// 			ActionFlag: d.CustGroupHoldRtn.ActionFlag,
// 			ActionKey:  d.ActionKey,
// 			session:    session,
// 			Dto:        d,
// 		}
// 		rethink.Replace()
// 	}

// }

// type NearDediveHoldRtnDto struct {
// 	ActionKey string `json:"ActionKey,omitempty"`
// 	pb.NearDediveHoldRtn
// }

// func SubscribeNearDediveHold(client pb.RiskMonitorServerClient, session *r.Session) {
// 	log.Info("Looking for NearDediveHold within ")

// 	clientDeadline := time.Now().Add(deadline)
// 	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
// 	defer cancel()
// 	stream, err := client.SubscribeNearDediveHold(ctx, &pb.SubscribeReq{})
// 	if err != nil {
// 		log.Fatal("%v.NearDediveHold(_) = _, %v", client, err)
// 	}
// 	for {
// 		item, err := stream.Recv()
// 		if err == io.EOF {
// 			log.Info("EOF")
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal("%v.NearDediveHold(_) = _, %v", client, err)
// 		}

// 		d := &NearDediveHoldRtnDto{
// 			ActionKey:         string(item.MonitorNo) + "#" + string(item.CustNo) + "#" + string(item.ExchCode) + "#" + string(item.ContractCode),
// 			NearDediveHoldRtn: *item,
// 		}

// 		rethink := &RtnDto{
// 			TableName:  TableName_SubscribeNearDediveHold,
// 			ActionFlag: d.NearDediveHoldRtn.ActionFlag,
// 			ActionKey:  d.ActionKey,
// 			session:    session,
// 			Dto:        d,
// 		}
// 		rethink.Replace()
// 	}

// }

// type ProuctGroupRiskRtnDto struct {
// 	ActionKey string `json:"ActionKey,omitempty"`
// 	pb.ProuctGroupRiskRtn
// }

// func SubscribeProuctGroupRisk(client pb.RiskMonitorServerClient, session *r.Session) {
// 	log.Info("Looking for ProuctGroupRisk within ")

// 	clientDeadline := time.Now().Add(deadline)
// 	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
// 	defer cancel()
// 	stream, err := client.SubscribeProuctGroupRisk(ctx, &pb.SubscribeReq{})
// 	if err != nil {
// 		log.Fatal("%v.ProuctGroupRisk(_) = _, %v", client, err)
// 	}
// 	for {
// 		item, err := stream.Recv()
// 		if err == io.EOF {
// 			log.Info("EOF")
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal("%v.ProuctGroupRisk(_) = _, %v", client, err)
// 		}

// 		d := &ProuctGroupRiskRtnDto{
// 			ActionKey:          string(item.MonitorNo) + "#" + string(item.ProductGroupNo) + "#" + string(item.ContractCode),
// 			ProuctGroupRiskRtn: *item,
// 		}

// 		rethink := &RtnDto{
// 			TableName:  TableName_SubscribeProuctGroupRisk,
// 			ActionFlag: d.ProuctGroupRiskRtn.ActionFlag,
// 			ActionKey:  d.ActionKey,
// 			session:    session,
// 			Dto:        d,
// 		}
// 		rethink.Replace()
// 	}

// }
