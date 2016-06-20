package configuration

import (
	"time"

	"gopkg.in/fsnotify.v1"

	_ "github.com/coreos/dex/pkg/log"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server *Server
		Game   *Game
	}

	//Server specific configuration
	Server struct {
		//the address the server is bound on
		Address string

		//frames per second
		FPS int

		//time allowed to write a message to the peer
		WriteWait time.Duration

		//time allowed to read the nex pong message from the peer
		PongWait time.Duration

		//send pings to peer with this period. Must be less than pongWait
		PingPeriod time.Duration

		//maximum message size allowed from peer
		MaxMessageSize int64

		//size of the upgrader read buffer
		ReadBufferSize int

		//size of the upgrader write buffer
		WriteBufferSize int

		//jwt algorithm
		JWTAlgorithm string

		//path to the jwt key
		JWTSecret string

		//private signing key
		JWTPrivate string

		//the public counterpart of the private key
		JWTPublic string

		//allowed origins * means all
		AllowedOrigins []string
	}

	//Game specific configuration
	Game struct {
		Scramble bool
	}
)

//GetConfiguration retrieves the configuration object collapsing all inputs
//into one configuration object
func Load() *Config {
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

	return &Config{
		Server: &Server{
			Address:         viper.GetString("server.address"),
			FPS:             viper.GetInt("server.fps"),
			WriteWait:       viper.GetDuration("server.writeWait"),
			PongWait:        viper.GetDuration("server.pongWait"),
			PingPeriod:      viper.GetDuration("server.pingPeriod"),
			MaxMessageSize:  int64(viper.GetInt("server.maxMessageSize")),
			ReadBufferSize:  viper.GetInt("server.readBufferSize"),
			WriteBufferSize: viper.GetInt("server.writeBufferSize"),
			JWTAlgorithm:    viper.GetString("server.jwtAlgorithm"),
			JWTSecret:       viper.GetString("server.jwtKey"),
			JWTPrivate:      viper.GetString("server.jwtPrivate"),
			JWTPublic:       viper.GetString("server.jwtPublic"),
			AllowedOrigins:  viper.GetStringSlice("server.allowedOrigins"),
		},
		Game: &Game{
			Scramble: viper.GetBool("game.scramble"),
		},
	}
}

//WatchConfiguration watch the configuration for any changes and run the
//given function when it occurs
func Watch(f func(fsnotify.Event)) {
	viper.WatchConfig()
	viper.OnConfigChange(f)
}
