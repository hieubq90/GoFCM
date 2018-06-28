package config

import (
	"github.com/spf13/viper"
	"strings"
	"io/ioutil"
	"bytes"
	"fmt"
)

var defaultConf = []byte(`
core:
	host: "0.0.0.0" # ip address to bind (default: any)
	port: "2701"
	cpu_num: 0 # default is runtime.NumCPU()
	http_proxy: "" # proxy for FMC (default: none)
thrift:
	enabled: true
	port: "2702"
	framed: false
	buffered: false
rpc:
	enabled: false
	port: "2703"
api:
	admin_uri: "/admin"
	app_uri: "/api/applications"
	push_uri: "/api/push"
	stat_go_uri: "/api/stat/go"
	stat_app_uri: "/api/stat/app"
	config_uri: "/api/config"
	sys_stat_uri: "/sys/stats"
	metric_uri: "/metrics"
	health_uri: "/healthz"
log:
	format: "string" # string or json
	access_log: "stdout" # stdout: output to console, or define log path like "log/access_log"
	access_level: "debug"
  	error_log: "stderr" # stderr: output to console, or define log path like "log/error_log"
  	error_level: "error"
storage:
	path: "level.db"
`)

type ConfYaml struct {
	Core    SectionCore    `yaml:"core"`
	Thrift  SectionThrift  `yaml:"thrift"`
	RPC     SectionRPC     `yaml:"rpc"`
	API     SectionAPI     `yaml:"api"`
	Log     SectionLog     `yaml:"log"`
	Storage SectionStorage `yaml:"storage"`
}

type SectionCore struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	CpuNum    int    `yaml"cpu_num"`
	HttpProxy string `yaml:"http_proxy"`
}

type SectionThrift struct {
	Enabled  bool   `yaml:"host"`
	Port     string `yaml:"port"`
	Framed   bool   `yaml:"framed"`
	Buffered bool   `yaml:"buffered"`
}

type SectionRPC struct {
	Enabled bool   `yaml:"host"`
	Port    string `yaml:"port"`
}

type SectionAPI struct {
	AdminURI   string `yaml:"admin_uri"`
	AppURI     string `yaml:"app_uri"`
	PushURI    string `yaml:"push_uri"`
	StatGoURI  string `yaml:"stat_go_uri"`
	StatAppURI string `yaml:"stat_app_uri"`
	ConfigURI  string `yaml:"config_uri"`
	SysStatURI string `yaml:"sys_stat_uri"`
	MetricURI  string `yaml:"metric_uri"`
	HealthURI  string `yaml:"health_uri"`
}

type SectionLog struct {
	Format      string `yaml:"format"`
	AccessLog   string `yaml:"access_log"`
	AccessLevel string `yaml:"access_level"`
	ErrorLog    string `yaml:"error_log"`
	ErrorLevel  string `yaml:"error_level"`
}

type SectionStorage struct {
	Path    string `yaml:"path"`
}

func LoadConf(confPath string) (ConfYaml, error) {
	var conf ConfYaml

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("gofcm")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if confPath != "" {
		content, err := ioutil.ReadFile(confPath)

		if err != nil {
			return conf, err
		}

		viper.ReadConfig(bytes.NewBuffer(content))
	} else {
		// Search config in home directory with name ".gorush" (without extension).
		viper.AddConfigPath("/etc/gofcm/")
		viper.AddConfigPath("$HOME/.gofcm")
		viper.AddConfigPath(".")
		viper.SetConfigName("go_fcm")

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			// load default config
			viper.ReadConfig(bytes.NewBuffer(defaultConf))
		}
	}

	return conf, nil
}
