package main

import "github.com/Sirupsen/logrus"

//version information variables
var (
	version string
	build   string
	date    string
	log     = logrus.WithFields(logrus.Fields{"app": "astrio"})
)

func main() {
	cfg := GetConfiguration()

	s, err := NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(s.Run())
}
