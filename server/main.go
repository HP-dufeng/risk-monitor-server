package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	pb "github.com/fengdu/risk-monitor-server/pb"
	"golang.org/x/text/encoding/simplifiedchinese"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRiskMonitorServerServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}

var (
	csvDBFile = flag.String("csv_db_file", "testdata/appleStore_description.csv", "A csv file containing a list of contracts")
	port      = flag.Int("port", 8080, "The server port")

	encoder = simplifiedchinese.GBK.NewEncoder()
)

type server struct {
	records [][]string
}

func (s *server) HeartBeatReq(ctx context.Context, in *pb.HeartBeat) (*pb.HeartBeat, error) {
	return &pb.HeartBeat{}, nil
}

func (s *server) SubscribeTunnelRealFund(req *pb.SubscribeReq, stream pb.RiskMonitorServer_SubscribeTunnelRealFundServer) error {
	log.Println("send: SubscribeTunnelRealFund")
	// tick := time.Tick(1 * time.Nanosecond)
	for {
		// select {
		// case <-tick:
		quote := &pb.TunnelRealFundRtn{
			ActionFlag:         randomActionFlags(),
			MonitorNo:          []byte(strconv.Itoa(randomIndex(10))),
			TunnelCode:         []byte(strconv.Itoa(randomIndex(100))),
			CurrencyCode:       []byte(strconv.Itoa(randomIndex(2))),
			AvailMarginBalance: randomFloat(10),
			TodayMarginBalance: randomFloat(10),
			OccupiedMargin:     randomFloat(10),
			MarginOccupiedRate: randomFloat(10),
		}
		if err := stream.Send(quote); err != nil {
			return err
		}
		// }

	}

}

func (s *server) SubscribeCorpHoldMon(req *pb.SubscribeReq, stream pb.RiskMonitorServer_SubscribeCorpHoldMonServer) error {
	log.Println("send: SubscribeCorpHoldMon")
	// tick := time.Tick(1 * time.Nanosecond)
	for {
		// select {
		// case <-tick:
		quote := &pb.CorpHoldMonRtn{
			ActionFlag:    randomActionFlags(),
			MonitorNo:     []byte(strconv.Itoa(randomIndex(10))),
			ContractCode:  []byte(strconv.Itoa(randomIndex(100))),
			SmarketCode:   []byte("smarketCode"),
			SecCode:       []byte("secCode"),
			SecName:       []byte("SecName"),
			DirectionType: []byte("DirectionType"),
			HoldQty:       uint32(randomIndex(100)),
			HoldRate:      randomFloat(10),
		}
		if err := stream.Send(quote); err != nil {
			return err
		}
		// }

	}

}

func (s *server) SubscribeQuoteMon(req *pb.SubscribeReq, stream pb.RiskMonitorServer_SubscribeQuoteMonServer) error {
	log.Println("send: SubscribeQuoteMon")
	// tick := time.Tick(1 * time.Nanosecond)
	for {
		// select {
		// case <-tick:
		quote := &pb.QuoteMonRtn{
			ActionFlag:         randomActionFlags(),
			MonitorNo:          []byte(strconv.Itoa(randomIndex(10))),
			ContractCode:       []byte(strconv.Itoa(randomIndex(100))),
			SmarketCode:        []byte("smarketCode"),
			SecCode:            []byte("secCode"),
			TradeType:          []byte("1"),
			LastPrice:          randomFloat(10),
			PreSettlementPrice: randomFloat(10),
			SettlementPrice:    randomFloat(10),
			RfLimitPrice:       randomFloat(10),
			Chg:                randomFloat(10),
			RiskLevel:          int32(randomIndex(100)),
			OpenInterest:       uint32(randomIndex(100)),
			MinMargin:          randomFloat(10),
			MarginRatio:        randomFloat(10),
			MarginDiff:         randomFloat(10),
			TickPriceNum:       uint32(randomIndex(100)),
		}
		if err := stream.Send(quote); err != nil {
			return err
		}
		// }

	}

}

func (s *server) SubscribeCustRisk(req *pb.SubscribeReq, stream pb.RiskMonitorServer_SubscribeCustRiskServer) error {
	log.Println("send: SubscribeCustRisk")
	// tick := time.Tick(1 * time.Nanosecond)
	for {
		// 	select {
		// 	case <-tick:
		custName, _ := encoder.Bytes([]byte("客户名称"))
		item := &pb.CustRiskRtn{
			ActionFlag:          randomActionFlags(),
			MonitorNo:           []byte(strconv.Itoa(randomIndex(10))),
			CustNo:              []byte(strconv.Itoa(randomIndex(10000))),
			CustClass:           []byte("CustClass"),
			CustName:            custName,
			MobilePhone:         []byte("MobilePhone"),
			Clientmode:          []byte("Clientmode"),
			RiskLevel:           []byte("RiskLevel"),
			RiskDegree0:         randomFloat(10),
			RiskDegree1:         randomFloat(10),
			RiskDegree2:         randomFloat(10),
			RiskDegree3:         randomFloat(10),
			LastRiskLevel:       []byte("LastRiskLevel"),
			LastRemain:          randomFloat(10),
			Margin:              randomFloat(10),
			DropProfit:          randomFloat(10),
			HoldProfit:          randomFloat(10),
			TodayInout:          randomFloat(10),
			RoyaltyInout:        randomFloat(10),
			DynCapRight:         randomFloat(10),
			ExchMargin:          randomFloat(10),
			AvailFund:           randomFloat(10),
			OptionCap:           randomFloat(10),
			DynRights:           randomFloat(10),
			OptionDynMargin:     randomFloat(10),
			FrznMargin:          randomFloat(10),
			FrznRoyalty:         randomFloat(10),
			ExchFrznMargin:      randomFloat(10),
			FrznStrikeMargin:    randomFloat(10),
			OptionNowMargin:     randomFloat(10),
			ExchOptionNowMargin: randomFloat(10),
			ExchOptionDynMargin: randomFloat(10),
			RiskContractQty:     uint32(randomIndex(100)),
			CurrencyCode:        []byte(strconv.Itoa(randomIndex(2))),
			TradingNo:           []byte("TradingNo"),
			DynRatio:            randomFloat(10),
		}
		if err := stream.Send(item); err != nil {
			return err
		}

		// }

	}

}

func (s *server) SubscribeCustHold(req *pb.SubscribeReq, stream pb.RiskMonitorServer_SubscribeCustHoldServer) error {
	log.Println("send: SubscribeCustHold")
	// tick := time.Tick(1 * time.Nanosecond)
	for {
		// select {
		// case <-tick:
		quote := &pb.CustHoldRtn{
			ActionFlag:   randomActionFlags(),
			MonitorNo:    []byte(strconv.Itoa(randomIndex(10))),
			CustNo:       []byte(strconv.Itoa(randomIndex(50))),
			CustClass:    []byte("CustClass"),
			ExchCode:     []byte("ExchCode"),
			VariCode:     []byte("VariCode"),
			ContractCode: []byte(strconv.Itoa(randomIndex(50))),
			DelivDate:    []byte("DelivDate"),
			HoldSum:      uint32(randomIndex(100)),
			TradeType:    []byte("1"),
			CpFlag:       []byte("CpFlag"),
			MonthFlag:    []byte("MonthFlag"),
			LimitRatio:   randomFloat(10),
			LimitVolmue:  uint32(randomIndex(100)),
			OverVolume:   uint32(randomIndex(100)),
			LimitWrning:  uint32(randomIndex(100)),
			OverWrning:   uint32(randomIndex(100)),
		}
		if randomActionFlags() == 1 {
			quote.HoldType = pb.CustHoldRtn_Speculate
		} else {
			quote.HoldType = pb.CustHoldRtn_Total
		}
		if err := stream.Send(quote); err != nil {
			return err
		}
		// }

	}

}

func (s *server) SubscribeCustGroupHold(req *pb.SubscribeReq, stream pb.RiskMonitorServer_SubscribeCustGroupHoldServer) error {
	log.Println("send: SubscribeCustGroupHold")
	// tick := time.Tick(1 * time.Nanosecond)
	for {
		// select {
		// case <-tick:
		quote := &pb.CustGroupHoldRtn{
			ActionFlag:    randomActionFlags(),
			MonitorNo:     []byte(strconv.Itoa(randomIndex(10))),
			CustGroupNo:   []byte(strconv.Itoa(randomIndex(10))),
			CustGroupName: []byte("CustGroupName"),
			ExchCode:      []byte("ExchCode"),
			VariCode:      []byte("VariCode"),
			ContractCode:  []byte(strconv.Itoa(randomIndex(100))),
			DelivDate:     []byte("DelivDate"),
			HoldSum:       uint32(randomIndex(100)),
			TradeType:     []byte("1"),
			CpFlag:        []byte("CpFlag"),
			MonthFlag:     []byte("MonthFlag"),
			LimitRatio:    randomFloat(10),
			LimitVolmue:   uint32(randomIndex(100)),
			OverVolume:    uint32(randomIndex(100)),
			LimitWrning:   uint32(randomIndex(100)),
			OverWrning:    uint32(randomIndex(100)),
		}
		if randomActionFlags() == 1 {
			quote.HoldType = pb.CustGroupHoldRtn_Speculate
		} else {
			quote.HoldType = pb.CustGroupHoldRtn_Total
		}
		if err := stream.Send(quote); err != nil {
			return err
		}
		// }

	}

}

func (s *server) SubscribeNearDediveHold(req *pb.SubscribeReq, stream pb.RiskMonitorServer_SubscribeNearDediveHoldServer) error {
	log.Println("send: SubscribeNearDediveHold")
	// tick := time.Tick(1 * time.Nanosecond)
	for {
		// select {
		// case <-tick:
		quote := &pb.NearDediveHoldRtn{
			ActionFlag:   randomActionFlags(),
			MonitorNo:    []byte(strconv.Itoa(randomIndex(10))),
			CustNo:       []byte(strconv.Itoa(randomIndex(50))),
			CustClass:    []byte("CustClass"),
			ExchCode:     []byte(randomExchCodes()),
			VariCode:     []byte("VariCode"),
			ContractCode: []byte(strconv.Itoa(randomIndex(50))),
			DelivDate:    []byte("DelivDate"),
			HoldSum:      uint32(randomIndex(100)),
			TradeType:    []byte("1"),
			CpFlag:       []byte("CpFlag"),
			ExpireDays:   uint32(randomIndex(100)),
			LimitBase:    uint32(randomIndex(100)),
			CloseDays:    uint32(randomIndex(100)),
		}
		if err := stream.Send(quote); err != nil {
			return err
		}
		// }

	}

}

func (s *server) SubscribeProuctGroupRisk(req *pb.SubscribeReq, stream pb.RiskMonitorServer_SubscribeProuctGroupRiskServer) error {
	log.Println("send: SubscribeProuctGroupRisk")
	// tick := time.Tick(1 * time.Nanosecond)
	for {
		// select {
		// case <-tick:
		quote := &pb.ProuctGroupRiskRtn{
			ActionFlag:       randomActionFlags(),
			MonitorNo:        []byte(strconv.Itoa(randomIndex(10))),
			ProductGroupNo:   []byte(strconv.Itoa(randomIndex(10))),
			ProductGroupName: []byte("ProductGroupName"),
			Count:            uint32(randomIndex(100)),
			RiskCount:        uint32(randomIndex(100)),

			RiskDegree:   randomFloat(10),
			SmarketCode:  []byte("secCode"),
			SecCode:      []byte("secCode"),
			ContractCode: []byte(strconv.Itoa(randomIndex(100))),

			TradeType: []byte("1"),

			LastPrice:          randomFloat(10),
			PreSettlementPrice: randomFloat(10),
			SettlementPrice:    randomFloat(10),
			RfLimitPrice:       randomFloat(10),
			Chg:                randomFloat(10),
			RiskLevel:          int32(randomIndex(100)),
			OpenInterest:       uint32(randomIndex(100)),
			MinMargin:          randomFloat(10),
			MarginRatio:        randomFloat(10),
			MarginDiff:         randomFloat(10),
			TickPriceNum:       uint32(randomIndex(100)),
		}
		if err := stream.Send(quote); err != nil {
			return err
		}
		// }

	}

}

// loadDatas loads contracts from a csv file.
func (s *server) loadDatas(filePath string) {
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}

	s.records = records
}

func randomIndex(length int) int {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	return r.Intn(length - 1)
}

func randomFloat(length int) float64 {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	return float64(r.Intn(length-1)) + r.Float64()
}

func randomActionFlags() uint32 {
	actionFlags := []uint32{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2}

	return actionFlags[randomIndex(len(actionFlags))]

}

func randomExchCodes() string {
	savedExchCodes := []string{"A", "B", "C", "G", "I", "S", "Z"}
	return savedExchCodes[randomIndex(len(savedExchCodes))]
}

func newServer() *server {
	s := &server{}
	// s.loadDatas(*csvDBFile)
	return s
}
