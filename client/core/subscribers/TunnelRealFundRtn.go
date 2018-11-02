package subscribers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fengdu/risk-monitor-server/client/core"
	pb "github.com/fengdu/risk-monitor-server/pb"
	log "github.com/sirupsen/logrus"
)

type TunnelRealFundRtnDto struct {
	ActionKey          string  `json:"ActionKey"`
	ActionFlag         uint32  `json:"ActionFlag"`
	MonitorNo          string  `json:"MonitorNo"`
	TunnelCode         string  `json:"TunnelCode"`
	CurrencyCode       string  `json:"CurrencyCode"`
	AvailMarginBalance float64 `json:"AvailMarginBalance"`
	TodayMarginBalance float64 `json:"TodayMarginBalance"`
	OccupiedMargin     float64 `json:"OccupiedMargin"`
	MarginOccupiedRate float64 `json:"MarginOccupiedRate"`
}

func (s *tunnelRealFundSubscriber) Convert(done <-chan struct{}, in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for n := range in {
			rtn := n.(*pb.TunnelRealFundRtn)
			dto := &TunnelRealFundRtnDto{
				ActionKey:          string(rtn.MonitorNo) + "#" + string(rtn.TunnelCode) + "#" + string(rtn.CurrencyCode),
				ActionFlag:         rtn.ActionFlag,
				MonitorNo:          string(rtn.MonitorNo),
				TunnelCode:         string(rtn.TunnelCode),
				CurrencyCode:       string(rtn.CurrencyCode),
				AvailMarginBalance: rtn.AvailMarginBalance,
				TodayMarginBalance: rtn.TodayMarginBalance,
				OccupiedMargin:     rtn.OccupiedMargin,
				MarginOccupiedRate: rtn.MarginOccupiedRate,
			}

			select {
			case out <- dto:
			case <-done:
				return
			}
		}
	}()

	return out
}

type tunnelRealFundSubscriber struct {
	client pb.RiskMonitorServerClient
}

//NewTunnelRealFundSubscriber create new instance of TunnelRealFund Subscriber
func NewTunnelRealFundSubscriber(client pb.RiskMonitorServerClient) core.Subscriber {
	return &tunnelRealFundSubscriber{
		client,
	}
}

func (s *tunnelRealFundSubscriber) Subscribe(done <-chan struct{}) (<-chan interface{}, <-chan error) {
	out := make(chan interface{})
	errc := make(chan error, 1)

	log.Info("Looking for TunnelRealFund within ")

	go func() {
		defer close(out)

		clientDeadline := time.Now().Add(7 * time.Duration(24) * time.Hour)
		ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
		defer cancel()
		stream, err := s.client.SubscribeTunnelRealFund(ctx, &pb.SubscribeReq{})
		if err != nil {
			log.Error("%v.TunnelRealFund(_) = _, %v", s.client, err)
			errc <- err
			return
		}

		for {
			item, err := stream.Recv()
			if err != nil {
				errc <- err
				return
			}

			select {
			case out <- item:
			case <-done:
				return
			}
		}
	}()

	return out, errc
}

func (s *tunnelRealFundSubscriber) Enqueue(done <-chan struct{}, in <-chan interface{}) (<-chan interface{}, <-chan error) {
	out := make(chan interface{})
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		for n := range in {
			dto := n.(*TunnelRealFundRtnDto)
			// if dto.actionFlag != 2 {
			// 	var item map[string]interface{}
			// 	j, _ := json.Marshal(m.message)
			// 	json.Unmarshal(j, &item)
			// 	// r.Table(m.tableName).Get(m.actionKey).Replace(item).RunWrite(m.session)
			// } else {
			// 	// r.Table(m.tableName).Get(m.actionKey).Replace(nil).RunWrite(m.session)
			// }

			var m map[string]interface{}
			bs, err := json.Marshal(dto)
			if err != nil {
				errc <- err
				return
			}

			err = json.Unmarshal(bs, &m)
			if err != nil {
				errc <- err
				return
			}

			select {
			case out <- m:
			case <-done:
				return
			}
		}

	}()

	return out, errc
}
