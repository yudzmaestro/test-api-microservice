package config

import (
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	"github.com/go-playground/validator/v10"
)

const (
	URL_IDM_AUTHORIZE = "authorize"
	URL_NOTIFICATION_SEND_NOTIF = "notification_send_notif"
)

type Config struct {
	Server 			*ServerConfig `yaml:"server""`
	DB 				*DBConfig `yaml:"db" validate:"required"`
	Auth         	*AuthConfig        `yaml:"auth" validate:"required"`
	Integrations 	*IntegrationConfig `yaml:"integrations" validate:"required"`
}

type ServerConfig struct {
	Protocol string `yaml:"protocol"`
	ListenAddress string `yaml:"listen_address"`
	LogLevel      int    `yaml:"log_level"`
	EnableCache   bool   `yaml:"enable_cache"`
}

type DBConfig struct {
	Host		string	`yaml:"host" validate:"required"`
	Port		int		`yaml:"port" validate:"required"`
	DbName		string	`yaml:"dbname" validate:"required"`
	User		string	`yaml:"user" validate:"required"`
	Password	string	`yaml:"password" validate:"required"`
}

// for hist serrvice
type AuthConfig struct {
	Token     string `yaml:"token" validate:"required"`
	Datetime  string `yaml:"datetime" validate:"required"`
	Signature string `yaml:"signature" validate:"required"`
}

type IntegrationConfig struct {
	HttpDialTimeoutSeconds int `yaml:"http_dial_timeout_seconds"`
	HttpRequestTimeoutSeconds int `yaml:"http_request_timeout_seconds"`
	UTKey string
	Externals *ExternalServices `yaml:"externals" validate:"required"`
}

type ExternalServices struct {
	Http map[string] *ExternalHttpServiceConfig `yaml:"http" validate:"required"`
	Grpc map[string] *ExternalGrpcServiceConfig `yaml:"grpc"`
}

type ExternalHttpServiceConfig struct {
	Scheme string `yaml:"scheme" validate:"required"`
	Host string `yaml:"host" validate:"required"`
	Endpoints map[string]string `yaml:"endpoints" validate:"required"`
}

type ExternalGrpcServiceConfig struct {
	Host string
	Port int
}

func NewConfigFromYAML(src io.Reader) (*Config, error) {
	var conf Config
	buf, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %s", err)
	}
	if err := yaml.Unmarshal(buf, &conf); err != nil {
		return nil, err
	}

	if err := checkConfig(&conf); err != nil {
		return nil, err
	}

	if err := checkConfig(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func checkConfig(config *Config) error {
	v := validator.New()
	if err := v.Struct(*config); err != nil {
		return fmt.Errorf("config missing required fields: %s", err)
	}

	if _, ok := config.Integrations.Externals.Http["idm"]; !ok {
		return fmt.Errorf("required idm integration config in config file")
	} else {
		eplist := []string{URL_IDM_AUTHORIZE}
		for _, l := range eplist {
			found := false
			for ep, _ := range config.Integrations.Externals.Http["idm"].Endpoints {
				if ep == l {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("cannot find '%s' endpoint in idm integration config", l)
			}
		}
	}

	//if _, ok := config.Integrations.Externals.Http["repayment"]; !ok {
	//	return fmt.Errorf("required repayment integration config in config file")
	//} else {
	//	eplist := []string{URL_REPAYMENT_REGISTER_REPAYMENT}
	//	for _, l := range eplist {
	//		found := false
	//		for ep, _ := range config.Integrations.Externals.Http["repayment"].Endpoints {
	//			if ep == l {
	//				found = true
	//				break
	//			}
	//		}
	//		if !found {
	//			return fmt.Errorf("cannot find '%s' endpoint in repayment integration config", l)
	//		}
	//	}
	//}

	return nil
}
