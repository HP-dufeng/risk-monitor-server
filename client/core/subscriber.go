package core

type Subscriber interface {
	Subscribe(done <-chan struct{}) (<-chan interface{}, <-chan error)
	Convert(done <-chan struct{}, in <-chan interface{}) <-chan interface{}
	Enqueue(done <-chan struct{}, in <-chan interface{}) (<-chan interface{}, <-chan error)
}
