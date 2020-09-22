package notify

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/x893675/gocron/pkg/utils/reflectutils"
)

type EmailOptions struct {
	Host        string `json:"host,omitempty" yaml:"host"`
	Port        int    `json:"port,omitempty" yaml:"port"`
	AccountAddr string `json:"addr,omitempty" yaml:"accountAddr"`
	Password    string `json:"password,omitempty" yaml:"password"`
	FromAddr    string `json:"from,omitempty" yaml:"fromAddr"`
	SkipVerify  bool   `json:"skip_verify,omitempty" yaml:"skipVerify"`
}

func NewEmailOptions() *EmailOptions {
	return &EmailOptions{
		Host:        "",
		Port:        587,
		AccountAddr: "",
		Password:    "",
		FromAddr:    "",
		SkipVerify:  true,
	}
}

func (e *EmailOptions) Validate() []error {
	var errors []error

	return errors
}

func (e *EmailOptions) ApplyTo(options *EmailOptions) {
	if e.Host != "" && e.AccountAddr != "" {
		reflectutils.Override(options, e)
	}
}

func (e *EmailOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&e.Host, "mail-host", e.Host, ""+
		"Mail server host, if left blank, all of the following mail options will "+
		"be ignored and mail will be disabled.")

	fs.IntVar(&e.Port, "mail-port", e.Port, ""+
		"Mail server port.")

	fs.StringVar(&e.AccountAddr, "mail-addr", e.AccountAddr, ""+
		"Mail account addr.")

	fs.StringVar(&e.Password, "mail-password", e.Password, ""+
		"Mail account password.")

	fs.StringVar(&e.FromAddr, "mail-from", e.FromAddr, ""+
		"Mail from.")

	fs.BoolVar(&e.SkipVerify, "mail-skip-verify", e.SkipVerify, ""+
		"Mail enable/disable ssl verify.")
}

func (e *EmailOptions) GetMailDSN() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}
