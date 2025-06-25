package config

import (
	"fmt"
	"os"

	"gorm.io/gorm"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

var DB *gorm.DB
var ServerConfig *Configuration

// Configuration holds the necessary information
// for our server, including the Iris one.
type Configuration struct {
	ServerName string `yaml:"ServerName" envconfig:"SERVER_NAME"`
	Env        string `yaml:"Env" envconfig:"ENV"`
	// The server's host, if empty, defaults to 0.0.0.0
	Host string `yaml:"Host" envconfig:"HOST"`
	// The server's port, e.g. 80
	Port int `yaml:"Port" envconfig:"PORT"`
	// If not empty runs under tls with this domain using letsencrypt.
	Domain string `yaml:"Domain" envconfig:"DOMAIN"`
	// Enables write response and read request compression.
	EnableCompression bool `yaml:"EnableCompression" envconfig:"ENABLE_COMPRESSION"`
	// Defines the "Access-Control-Allow-Origin" header of the CORS middleware.
	// Many can be separated by comma.
	// Defaults to "*" (allow all).
	AllowOrigin string `yaml:"AllowOrigin" envconfig:"ALLOW_ORIGIN"`
	// If not empty a request logger is registered,
	// note that this will cost a lot in performance, use it only for debug.
	EnableRequestLog bool `yaml:"EnableRequestLog" envconfig:"ENABLE_REQUEST_LOG"`
	Database         struct {
		Host               string `yaml:"Host" envconfig:"MASTER_DB_HOST"`
		Port               string `yaml:"Port" envconfig:"MASTER_DB_PORT"`
		Name               string `yaml:"Name" envconfig:"MASTER_DB_NAME"`
		Username           string `yaml:"Username" envconfig:"MASTER_DB_USER"`
		Password           string `yaml:"Password" envconfig:"MASTER_DB_PASSWORD"`
		ConnectionString   string `yaml:"ConnectionString" envconfig:"DB_CONNECTION_STRING"`
		MaxIdleConnections int    `yaml:"MaxIdleConnections" envconfig:"DB_MAX_IDLE_CONNECTIONS"`
		MaxOpenConnections int    `yaml:"MaxOpenConnections" envconfig:"DB_MAX_OPEN_CONNECTIONS"`
		SSLmode            string `yaml:"SSLmode" envconfig:"DB_SSL_MODE"`
		TablePrefix        string `yaml:"TablePrefix" envconfig:"DB_TABLE_PREFIX"`
	} `yaml:"Database"`
	Redis struct {
		URL string `yaml:"Url" envconfig:"REDIS_URL"`
	} `yaml:"Redis"`
}

func ReadConfigDotEnv(c *Configuration) {
	// DB_PASS
	c.Database.Password = os.Getenv("DB_PASS")
}

func readEnv(cfg *Configuration) {
	err := envconfig.Process("welab", cfg)
	if err != nil {
		panic(err)
	}
}

// BindFile binds the yaml file's contents to this Configuration.
func (c *Configuration) BindFile(filename string) error {
	ReadConfigDotEnv(c)
	contents, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(contents, c); err != nil {
		return err
	}
	// Read from env
	readEnv(c)
	ServerConfig = c
	fmt.Println("🚀 Server:  ", c.ServerName)
	return nil
}
