package main

import (
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	Port               string `envconfig:"PORT" default:"8081"`
	SharedKey          string `envconfig:"KEY" default:"0x24FEEDFACEDEADBEEFCAFE"`
	MaxRedirect        int    `envconfig:"MAX_REDIRECTS" default:"4"`
	GamoHostname       string `envconfig:"HOSTNAME" default:"unknown"`
	SocketTimeout      int    `envconfig:"SOCKET_TIMEOUT" default:"15"`
	KeepAlive          bool   `envconfig:"KEEP_ALIVE" default:"false"`
	ContentLengthLimit int64  `envconfig:"LENGTH_LIMIT" default:"5242880"`
}

var env *Env

func initEnv() {
	env = new(Env)
	envconfig.Process("GAMO", env)
}

func (e *Env) getAddr() string {
	return ":" + e.Port
}
