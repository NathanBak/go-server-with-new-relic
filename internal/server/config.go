package server

import (
	"fmt"
	"time"

	"github.com/NathanBak/go-server-with-new-relic/pkg/widget"
)

// Config contains information necessary to set up a Server.
type Config struct {
	Port         int           `json:"port" envvar:"PORT,required"`
	ReadTimeout  time.Duration `json:"readTimeout" envvar:"READ_TIMEOUT,default=3s"`
	WriteTimeout time.Duration `json:"writeTimeout" envvar:"WRITE_TIMEOUT,default=3s"`

	Logger Logger `json:"-" envvar:"-"`

	IncludeStatusCodeInMessages bool `json:"-" envvar:"-"`

	Storage Storage `json:"-" envvar:"-"`
}

// The Logger interface defines the methods required by the Server for logging.
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string)
}

// The Storage interface defines the methods required to access backing Widget storage.
type Storage interface {
	Get(string) (widget.Widget, bool, error)
	Set(string, widget.Widget) error
	Delete(string) (widget.Widget, bool, error)
	Keys() ([]string, error)
}

// CfgBuildInit initializes the Logger.  It should only be called by a cfgbuild.Builder.
func (cfg *Config) CfgBuildInit() error {
	if cfg.Logger == nil {
		cfg.Logger = defaultLogger{}
	}
	return nil
}

// CfgBuildValidate checks the specified values.  It should only be called by a cfgbuild.Builder.
func (cfg *Config) CfgBuildValidate() error {
	if cfg.Port < 1 || cfg.Port > 65535 {
		return fmt.Errorf("%d is not a valid port", cfg.Port)
	}
	return nil
}
