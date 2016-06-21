package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/rauwekost/astrio/configuration"
	"github.com/rauwekost/astrio/server"
)

//version information variables
var (
	version string
	build   string
	date    string
	log     = logrus.WithField("p", "main,b")
)

func main() {
	//load initial configuation
	configuration.Load()

	//run the server
	log.Fatal(server.New().Run())
}
