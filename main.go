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
	log     = logrus.WithFields(logrus.Fields{"app": "astrio"})
)

func main() {
	s, err := server.New(configuration.Load())
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(s.Run())
}
