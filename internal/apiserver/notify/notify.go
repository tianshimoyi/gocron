package notify

type Message map[string]interface{}

type Notify interface {
	Send(msg Message)
}

type Inform struct {
	informers map[string]Notify
	queue     chan Message
}

func NewInform() *Inform {
	i := &Inform{
		informers: make(map[string]Notify),
		queue:     make(chan Message, 100),
	}
	//i.informers["webhook"] =
	return i
}

func (i *Inform) Push(msg Message) {
	i.queue <- msg
}

func (i *Inform) Run() {

}
