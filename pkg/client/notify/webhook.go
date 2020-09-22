package notify

import (
	"encoding/json"
	"github.com/spf13/pflag"
	"github.com/x893675/gocron/pkg/client/httpclient"
	"github.com/x893675/gocron/pkg/utils/reflectutils"
	"github.com/x893675/gocron/pkg/utils/stringutils"
	"html"
	"k8s.io/klog/v2"
	"net/http"
	"time"
)

type WebHookMsg struct {
	Msg Message `json:"message"`
}

type WebhookOptions struct {
	Url        string        `json:"url,omitempty" yaml:"url"`
	RetryTimes int           `json:"retry,omitempty" yaml:"retry"`
	TimeOut    time.Duration `json:"timeout,omitempty" yaml:"timeout"`
}

func NewWebhookOptions() *WebhookOptions {
	return &WebhookOptions{
		Url:        "",
		RetryTimes: 3,
		TimeOut:    30 * time.Second,
	}
}

func (w *WebhookOptions) Validate() []error {
	return nil
}

func (w *WebhookOptions) ApplyTo(opt *WebhookOptions) {
	if w.Url != "" {
		reflectutils.Override(opt, w)
	}
}

func (w *WebhookOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&w.Url, "webhook-url", w.Url, ""+
		"Notify webhook url, if left blank, all of the following mail options will "+
		"be ignored and webhook will be disabled.")

	fs.IntVar(&w.RetryTimes, "webhook-retry", w.RetryTimes, ""+
		"Webhook request retry times.")

	fs.DurationVar(&w.TimeOut, "webhook-timeout", w.TimeOut, ""+
		"Webhook request timeout, default is 30s")
}

var _ Notify = (*WebhookClient)(nil)

type WebhookClient struct {
	Url        string
	RetryTimes int
	Timeout    time.Duration
}

func NewWebhookClient(opt *WebhookOptions) *WebhookClient {
	return &WebhookClient{
		Url:        opt.Url,
		RetryTimes: opt.RetryTimes,
		Timeout:    opt.TimeOut,
	}
}

func (w *WebhookClient) Send(msg Message) {
	buf := parseWebhookTemplate(msg)
	content := html.UnescapeString(buf)
	klog.V(2).Infof("webhook content is %v", content)
	i := 0
	for i < w.RetryTimes {
		resp := httpclient.PostJson(w.Url, content, int(w.Timeout))
		if resp.StatusCode == http.StatusAccepted {
			break
		}
		i += 1
		time.Sleep(2 * time.Second)
		if i < w.RetryTimes {
			klog.Errorf("webHook#发送消息失败#%s#消息内容-%s", resp.Body, content)
		}
	}
}

func parseWebhookTemplate(msg Message) string {
	m := WebHookMsg{
		Msg: msg,
	}
	s, err := json.Marshal(m)
	if err != nil {
		klog.Errorf("marshal webhook msg error")
		return ""
	}
	return stringutils.Bytes2string(s)
}
