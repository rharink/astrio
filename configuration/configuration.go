package configuration

import (
	"time"

	"gopkg.in/fsnotify.v1"

	_ "github.com/coreos/dex/pkg/log"
	"github.com/spf13/viper"
)

var (
	//the address the server is bound on
	ServerAddress string

	//frames per second
	ServerFPS int

	//time allowed to write a message to the peer
	ServerWriteWait time.Duration

	//time allowed to read the nex pong message from the peer
	ServerPongWait time.Duration

	//send pings to peer with this period. Must be less than pongWait
	ServerPingPeriod time.Duration

	//maximum message size allowed from peer
	ServerMaxMessageSize int64

	//size of the upgrader read buffer
	ServerReadBufferSize int

	//size of the upgrader write buffer
	ServerWriteBufferSize int

	//jwt algorithm
	ServerJWTAlgorithm string

	//path to the jwt key
	ServerJWTSecret string

	//private signing key
	ServerJWTPrivate string

	//the public counterpart of the private key
	ServerJWTPublic string

	//allowed origins * means all
	ServerAllowedOrigins []string
)

//GetConfiguration retrieves the configuration object collapsing all inputs
//into one configuration object
func Load() {
	//setup default configuration file locations
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	//setup environment variables
	viper.SetEnvPrefix("astrio")
	viper.AutomaticEnv()

	//set default server configuration
	viper.SetDefault("server.address", "127.0.0.1:4000")
	viper.SetDefault("server.fps", 24)
	viper.SetDefault("server.writeWait", "1s")
	viper.SetDefault("server.pongWait", "1s")
	viper.SetDefault("server.pingPeriod", "0.5s")
	viper.SetDefault("server.maxMessageSize", 1024)
	viper.SetDefault("server.readBufferSize", 512)
	viper.SetDefault("server.writeBufferSize", 512)
	viper.SetDefault("server.jwtAlgorithm", "RS256")
	viper.SetDefault("server.jwtSecret", "secret")
	viper.SetDefault("server.jwtPrivate", "./server.key")
	viper.SetDefault("server.jwtPublic", "./server.pub")
	viper.SetDefault("server.allowedOrigins", []string{"astr.io"})

	//read configuration from file(s)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	//set variables
	ServerAddress = viper.GetString("server.address")
	ServerFPS = viper.GetInt("server.fps")
	ServerWriteWait = viper.GetDuration("server.writeWait")
	ServerPongWait = viper.GetDuration("server.pongWait")
	ServerPingPeriod = viper.GetDuration("server.pingPeriod")
	ServerMaxMessageSize = int64(viper.GetInt("server.maxMessageSize"))
	ServerReadBufferSize = viper.GetInt("server.readBufferSize")
	ServerWriteBufferSize = viper.GetInt("server.writeBufferSize")
	ServerJWTAlgorithm = viper.GetString("server.jwtAlgorithm")
	ServerJWTSecret = viper.GetString("server.jwtKey")
	ServerJWTPrivate = viper.GetString("server.jwtPrivate")
	ServerJWTPublic = viper.GetString("server.jwtPublic")
	ServerAllowedOrigins = viper.GetStringSlice("server.allowedOrigins")
}

//WatchConfiguration watch the configuration for any changes and run the
//given function when it occurs
func Watch(f func(fsnotify.Event)) {
	Load()
	viper.WatchConfig()
	viper.OnConfigChange(f)
}
