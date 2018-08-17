package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	pb "github.com/fengdu/risk-monitor-server/pb"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("10.1.7.127:%d", *port))
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
)

type server struct {
	savedContracts []*pb.QryContractRsp
	savedCustomers []*pb.QryCustSimpleInfoRsp
	savedVarieties []*pb.QryExchVariRsp
	savedExchCodes []string
}

func (s *server) QryContract(req *pb.QryContractReq, stream pb.RiskMonitorServer_QryContractServer) error {
	for _, contract := range s.savedContracts {
		if err := stream.Send(contract); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) QryCustSimpleInfo(req *pb.QryCustSimpleInfoReq, stream pb.RiskMonitorServer_QryCustSimpleInfoServer) error {
	for _, customer := range s.savedCustomers {
		if err := stream.Send(customer); err != nil {
			return err
		}
	}

	return nil
}

func (s *server) SubscribeQuoteMon(req *pb.SubscribeReq, stream pb.RiskMonitorServer_SubscribeQuoteMonServer) error {
	tick := time.Tick(1 * time.Nanosecond)
	for {
		select {
		case <-tick:
			quote := &pb.QuoteMonRtn{
				ActionFlag:   randomActionFlags(),
				MonitorNo:    []byte(strconv.Itoa(randomIndex(10))),
				ContractCode: []byte(strconv.Itoa(randomIndex(100))),
			}
			if err := stream.Send(quote); err != nil {
				return err
			}
		}

	}
	// for _, v := range []int{1, 2, 3, 4} {
	// 	contractCode := strconv.FormatInt(time.Now().UnixNano(), 10)
	// 	quote := &pb.QuoteMonRtn{
	// 		ActionFlag:   randomActionFlags(),
	// 		MonitorNo:    []byte(strconv.Itoa(v)),
	// 		ContractCode: []byte(contractCode),
	// 	}
	// 	if err := stream.Send(quote); err != nil {
	// 		return err
	// 	}
	// }
	// return nil

}

func (s *server) SubscribeCustRisk(req *pb.SubscribeReq, stream pb.RiskMonitorServer_SubscribeCustRiskServer) error {
	tick := time.Tick(1 * time.Nanosecond)
	for {
		select {
		case <-tick:

			customer := &pb.CustRiskRtn{
				ActionFlag: randomActionFlags(),
				MonitorNo:  []byte(strconv.Itoa(randomIndex(5))),
				CustNo:     []byte(strconv.Itoa(randomIndex(100))),
			}
			if err := stream.Send(customer); err != nil {
				return err
			}
		}

	}

}

func (s *server) QryExchVari(req *pb.QryExchVariReq, stream pb.RiskMonitorServer_QryExchVariServer) error {
	for _, variety := range s.savedVarieties {
		if err := stream.Send(variety); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) SetRiskContrLevel(stream pb.RiskMonitorServer_SetRiskContrLevelServer) error {
	for {
		setting, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.RspInfo{
				Errid:  0,
				Errmsg: []byte(""),
			})
		}
		if err != nil {
			return err
		}

		log.Println("SetRiskContrLevel", setting)
	}

}

func (s *server) SetCustRiskMonitor(stream pb.RiskMonitorServer_SetCustRiskMonitorServer) error {
	for {
		setting, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.RspInfo{
				Errid:  0,
				Errmsg: []byte(""),
			})
		}
		if err != nil {
			return err
		}

		log.Println("SetCustRiskMonitor", setting)
	}
}

func (s *server) SetQuoteMonitor(stream pb.RiskMonitorServer_SetQuoteMonitorServer) error {
	for {
		setting, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.RspInfo{
				Errid:  0,
				Errmsg: []byte(""),
			})
		}
		if err != nil {
			return err
		}

		log.Println("SetQuoteMonitor", setting)
	}
}

// func (s *server) SetContrGroup(stream pb.RiskMonitorServer_SetContrGroupServer) error {
// 	for {
// 		setting, err := stream.Recv()
// 		if err == io.EOF {
// 			return stream.SendAndClose(&pb.RspInfo{
// 				Errid:  0,
// 				Errmsg: []byte(""),
// 			})
// 		}
// 		if err != nil {
// 			return err
// 		}

// 		log.Println(setting)
// 	}
// }

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

	for i, v := range records {
		if i == 0 {
			continue
		}
		if i > 2000 {
			break
		}

		s.savedContracts = append(s.savedContracts, &pb.QryContractRsp{
			ContractCode:      []byte(v[0]),
			ContractShortName: []byte(v[1]),
		})

		s.savedCustomers = append(s.savedCustomers, &pb.QryCustSimpleInfoRsp{
			CustNo:   []byte(v[0]),
			CustName: []byte(v[1]),
		})

		if i < 80 {
			s.savedVarieties = append(s.savedVarieties, &pb.QryExchVariRsp{
				VariCode:  []byte(v[0]),
				VariName:  []byte(v[1]),
				ExchCode:  []byte(s.savedExchCodes[randomIndex(len(s.savedExchCodes))]),
				TradeType: []byte(strconv.Itoa(randomIndex(i+1) % 2)),
			})
		}

	}

}

func randomIndex(length int) int {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	return r.Intn(length - 1)
}

func randomActionFlags() uint32 {
	actionFlags := []uint32{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2}

	return actionFlags[randomIndex(len(actionFlags))]

}

func newServer() *server {
	s := &server{savedExchCodes: []string{"A", "B", "C", "G", "I", "S", "Z"}}
	s.loadDatas(*csvDBFile)
	return s
}
