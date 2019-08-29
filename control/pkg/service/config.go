package service

import (
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

type FileConfig struct {
	Hosts string `mapstructure:"hosts"`
}

type DbConfig struct {
	Hosts              string `mapstructure:"hosts"`
	Database           string `mapstructure:"database"`
	Username           string `mapstructure:"username"`
	Password           string `mapstructure:"password"`
	UserCollection     string `mapstructure:"user_collection"`
	HospitalCollection string `mapstructure:"hospital_collection"`
	ConfigCollection   string `mapstructure:"config_collection"`
	TopicCollection    string `mapstructure:"topic_collection"`
}
type RedisConfig struct {
	Hosts            string `mapstructure:"hosts"`
	WebSocketChannel string `mapstructure:"websocket_channel"`
}
type HostsConfig struct {
	HttpHosts string `mapstructure:"http_hosts"`
	GrpcHosts string `mapstructure:"grpc_hosts"`
	Hosts     string `mapstructure:"hosts"`
	File      bool   `mapstructure:"without_file_server"`
}
type ServerConfig struct {
	Control HostsConfig `mapstructure:"control"`
}
type ClientConfig struct {
	Control   HostsConfig `mapstructure:"control"`
	Image     HostsConfig `mapstructure:"image"`
	Patient   HostsConfig `mapstructure:"patient"`
	Date      HostsConfig `mapstructure:"date"`
	Diagnosis HostsConfig `mapstructure:"diagnosis"`
	File      HostsConfig `mapstructure:"file"`
	Master    HostsConfig `mapstructure:"master"`
	Websocket HostsConfig `mapstructure:"websocket"`
}
type TimeConfig struct {
	Format string `mapstructure:"format"`
}
type ListConfig struct {
	Default string   `mapstructure:"default"`
	List    []string `mapstructure:"list"`
}
type ExternalConfig struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
}
type RuleConfig struct {
	Name []string `mapstructure:"name"`
}

var serverHostConfig ServerConfig
var clientHostConfig ClientConfig
var mgoConfig DbConfig
var redisConfig RedisConfig
var envOs string
var timeConfig TimeConfig
var ethnicityConfig ListConfig
var countryConfig ListConfig
var skinConfig ListConfig
var diseaseConfig ListConfig
var locationConfig ListConfig
var genderConfig ListConfig
var tagConfig ListConfig
var systemType string
var externalConfig ExternalConfig
var ruleConfig RuleConfig
var withoutFileServer bool

func init() {
	if err := LoadConfig(); err != nil {
		panic(err)
	}

	viper.GetStringMap("db")
	viper.UnmarshalKey("db", &mgoConfig)
	viper.GetStringMap("redis")
	viper.UnmarshalKey("redis", &redisConfig)
	viper.GetStringMap("server")
	viper.UnmarshalKey("server", &serverHostConfig)
	viper.GetStringMap("client")
	viper.UnmarshalKey("client", &clientHostConfig)
	viper.GetStringMap("time")
	viper.UnmarshalKey("time", &timeConfig)
	viper.GetStringMap("ethnicity")
	viper.UnmarshalKey("ethnicity", &ethnicityConfig)
	viper.GetStringMap("country")
	viper.UnmarshalKey("country", &countryConfig)
	viper.GetStringMap("skin")
	viper.UnmarshalKey("skin", &skinConfig)
	viper.GetStringMap("disease")
	viper.UnmarshalKey("disease", &diseaseConfig)
	viper.GetStringMap("location")
	viper.UnmarshalKey("location", &locationConfig)
	viper.GetStringMap("gender")
	viper.UnmarshalKey("gender", &genderConfig)
	viper.GetStringMap("tag")
	viper.UnmarshalKey("tag", &tagConfig)
	viper.GetStringMap("external")
	viper.UnmarshalKey("external", &externalConfig)
	viper.GetStringMap("rule")
	viper.UnmarshalKey("rule", &ruleConfig)
	envOs = viper.GetString("platform")
	if len(envOs) == 0 {
		envOs = "linux"
	}
	systemType = viper.GetString("type")
	withoutFileServer = GetConfigWithoutFileServer()

	//viper.Debug()
}
func LoadConfig() (err error) {
	viper.SetConfigFile(os.Getenv("PIKABU_HOME") + "/config/config.json")
	if err = viper.ReadInConfig(); err != nil {
		viper.SetConfigFile(os.Getenv("PIKABU_HOME") + "/pikabu-control" + "/config/config.json")
		if err = viper.ReadInConfig(); err != nil {
			return
		}
		return
	}
	return
}
func GetConfigServerControlHttp() string {
	if strings.Contains(serverHostConfig.Control.HttpHosts, "PORT") {
		port := os.Getenv("PIKABU_CONTROL_PORT")
		hosts := strings.Replace(serverHostConfig.Control.HttpHosts, "PORT", port, 1)
		return hosts
	} else {
		return serverHostConfig.Control.HttpHosts
	}
}
func GetConfigServerControlGrpc() string {
	return serverHostConfig.Control.GrpcHosts
}
func GetConfigWithoutFileServer() bool {
	return serverHostConfig.Control.File
}
func GetConfigClientControlHttp() string {
	return clientHostConfig.Control.HttpHosts
}
func GetConfigClientControlGrpc() string {
	return clientHostConfig.Control.GrpcHosts
}
func GetConfigClientFileHttp() string {
	if strings.Contains(clientHostConfig.File.HttpHosts, "localhost") {
		ip, _ := ExternalIP()
		hosts := strings.Replace(clientHostConfig.File.HttpHosts, "localhost", ip, 1)
		return hosts
	} else if strings.Contains(clientHostConfig.File.HttpHosts, "LOCALHOST:PORT") {
		port := os.Getenv("PIKABU_CONTROL_PORT")
		ip, _ := ExternalIP()
		hosts := "http://" + ip + ":" + port
		return hosts
	} else {
		return clientHostConfig.File.HttpHosts
	}
}
func GetConfigClientMasterHttp() string {
	return clientHostConfig.Master.HttpHosts
}
func GetConfigClientImageGrpc() string {
	return clientHostConfig.Image.GrpcHosts
}
func GetConfigClientPatientGrpc() string {
	return clientHostConfig.Patient.GrpcHosts
}
func GetConfigClientDateGrpc() string {
	return clientHostConfig.Date.GrpcHosts
}
func GetConfigClientDiagnosisGrpc() string {
	return clientHostConfig.Diagnosis.GrpcHosts
}
func GetPikabuTimeFormat() string {
	return time.RFC3339
}
func GetTimeFormat() string {
	return timeConfig.Format
}
func GetDefaultEthnicity() string {
	return ethnicityConfig.Default
}
func GetEthnicityList() []string {
	return ethnicityConfig.List
}
func GetDefaultCountry() string {
	return countryConfig.Default
}
func GetCountryList() []string {
	return countryConfig.List
}
func GetDefaultSkin() string {
	return skinConfig.Default
}
func GetSkinList() []string {
	return skinConfig.List
}
func GetDefaultDisease() string {
	return diseaseConfig.Default
}
func GetDiseaseList() []string {
	return diseaseConfig.List
}
func GetDefaultLocation() string {
	return locationConfig.Default
}
func GetLocationList() []string {
	return locationConfig.List
}
func GetDefaultGender() string {
	return genderConfig.Default
}
func GetGenderList() []string {
	return genderConfig.List
}
func GetDefaultTag() string {
	return tagConfig.Default
}
func GetTagList() []string {
	return tagConfig.List
}
func GetPlatform() string {
	return envOs
}
func GetConfigSystemType() string {
	return systemType
}
func GetConfigExternalAddress() string {
	return externalConfig.Address
}
func GetConfigExternalPort() string {
	return externalConfig.Port
}
func GetConfigRuleName() []string {
	return ruleConfig.Name
}
func GetConfigWebsocketHosts() string {
	if len(clientHostConfig.Websocket.Hosts) == 0 {
		ip, _ := ExternalIPRefresh()
		s := strings.Split(GetConfigServerControlHttp(), ":")
		return "ws://" + ip + ":" + s[len(s) - 1] + "/ws"
	} else {
		return clientHostConfig.Websocket.Hosts
	}
}
