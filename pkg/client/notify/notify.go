package notify

import (
	"k8s.io/klog/v2"
	"time"
)

type Message struct {
	// sender id
	Sid string `json:"sender_id"`
	// receive ids
	Rids    []string `json:"receiver_ids"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	// category id
	Cid        string `json:"category_id"`
	NotifyType string `json:"-"`
}

type Notify interface {
	Send(msg Message)
}

type Inform struct {
	informers map[string]Notify
	queue     chan Message
}

func NewInform(opt *Options) *Inform {
	i := &Inform{
		informers: make(map[string]Notify),
		queue:     make(chan Message, 100),
	}
	if opt.WebhookOpt != nil && opt.WebhookOpt.Url != "" {
		klog.V(2).Infof("init webhook informer")
		i.informers["webhook"] = NewWebhookClient(opt.WebhookOpt)
	}
	return i
}

func (i *Inform) Push(msg Message) {
	i.queue <- msg
}

func (i *Inform) Run(stopCh <-chan struct{}) {
	for {
		select {
		case <-stopCh:
			goto EedLoop
		case msg := <-i.queue:
			klog.V(2).Infof("msg queue consumer %v", msg)
			switch msg.NotifyType {
			case "webhook":
				if h, ok := i.informers["webhook"]; ok {
					go h.Send(msg)
				} else {
					klog.Warning("webhook informer not exist")
				}
			default:
				klog.Error("unknown notify type")
			}
		default:
			time.Sleep(1 * time.Second)
		}
	}
EedLoop:
	klog.V(2).Infof("informer stop")
}
