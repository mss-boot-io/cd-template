package main

import (
	"flag"
	"github.com/mss-boot-io/cd-template/pkg/config"
	"github.com/mss-boot-io/cd-template/stage"
)

var configPath = flag.String("config", "", "config path")

func main() {
	flag.Parse()
	config.NewConfig(configPath)
	stage.Synth("prod")
	config.Cfg.Hpa = false
	config.Cfg.Resources = nil
	config.Cfg.Replicas = config.Cfg.TestReplicas
	stage.Synth("test")
}
