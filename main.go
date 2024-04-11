/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/11/1 09:31:56
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/11/1 09:31:56
 */

package main

import (
	"flag"
	"log"

	"github.com/mss-boot-io/cd-template/pkg/config"
	"github.com/mss-boot-io/cd-template/stage"
)

var configPath = flag.String("config", "", "config path")

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	config.NewConfig(configPath)
	stage.Synth("prod")
	config.Cfg.Hpa.Enabled = false
	config.Cfg.Resources = nil
	config.Cfg.Replicas = config.Cfg.TestReplicas
	stage.Synth("test")
}
