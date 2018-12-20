package core

import (
	log "github.com/sirupsen/logrus"
)

type Subscriber interface {
	Subscribe(done <-chan struct{}) (<-chan interface{}, <-chan error)
	Convert(done <-chan struct{}, in <-chan interface{}) <-chan interface{}
	Enqueue(done <-chan struct{}, in <-chan interface{}) (<-chan interface{}, <-chan error)
}

func Subscribe(s Subscriber) {
	done := make(chan struct{})
	defer close(done)

	sourceC, errSubscribe := s.Subscribe(done)
	convertC := s.Convert(done, sourceC)
	out, errEnqueue := s.Enqueue(done, convertC)

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
