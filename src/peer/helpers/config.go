package helpers

import (
	"strings"

	"github.com/spf13/viper"
)

type NodeConfig struct {
	Name string
	Url  string
}

func NewNodeConfig(name string, url string) *NodeConfig {
	config := &NodeConfig{
		Name: name,
		Url:  url,
	}
	return config
}

func generateURL(host string, port string) string {
	return host + ":" + port
}

// Read configuration environment
func LoadConfig() *NodeConfig {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("node")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.BindEnv("name")
	v.BindEnv("host")
	v.BindEnv("port")
	// get values
	name := v.GetString("name")
	port := v.GetString("port")
	host := v.GetString("host")
	var config = NewNodeConfig(name, generateURL(host, port))
	config.LogConfig()
	return config
}

func (config *NodeConfig) LogConfig() {
	Log.Debugf("Name: %v | Url: %v",
		config.Name,
		config.Url,
	)
}
