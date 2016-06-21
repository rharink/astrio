package configuration

import (
	"time"

	"github.com/spf13/viper"
)

var Server serverConfig
var Game gameConfig

type (
	//server specific configuration
	serverConfig struct {
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

	//game specific configuration
	gameConfig struct {
		Scramble   bool
		MaxPlayers int
	}
)

//load retrieves the configuration object collapsing all inputs
//into one configuration object
func Load() {
	//setup default configuration file locations
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	//setup environment variables
	viper.SetEnvPrefix("astrio")
	viper.AutomaticEnv()

	//read configuration from file(s)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	Server = serverConfig{
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
	}
	Game = gameConfig{
		Scramble:   viper.GetBool("game.scramble"),
		MaxPlayers: viper.GetInt("game.maxPlayers"),
	}
}
