package notify

import "github.com/spf13/pflag"

type Options struct {
	WebhookOpt  *WebhookOptions `json:"webhook" yaml:"webhookOpt"`
	MailOptions *EmailOptions   `json:"mail" yaml:"mailOptions"`
}

func NewNotifyOptions() *Options {
	return &Options{
		WebhookOpt:  NewWebhookOptions(),
		MailOptions: NewEmailOptions(),
	}
}

func (o *Options) Validate() []error {
	var err []error
	err = append(err, o.MailOptions.Validate()...)
	err = append(err, o.WebhookOpt.Validate()...)
	return err
}

func (o *Options) ApplyTo(opt *Options) {

}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.MailOptions.AddFlags(fs)
	o.WebhookOpt.AddFlags(fs)
}
