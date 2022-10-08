package stage

import (
	"path/filepath"

	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s"

	"github.com/mss-boot-io/cd-template/chart"
	"github.com/mss-boot-io/cd-template/pkg/config"
)

// Synth chart generate
func Synth(stage string) {
	app := cdk8s.NewApp(&cdk8s.AppProps{Outdir: jsii.String(filepath.Join("dist", config.Cfg.Service, stage))})
	chartProps := &cdk8s.ChartProps{
		Labels: &map[string]*string{
			"app":     &config.Cfg.Service,
			"version": &config.Cfg.Version,
		},
	}
	if &config.Cfg.Project != nil && *&config.Cfg.Project != "" {
		chartProps = &cdk8s.ChartProps{
			Labels: &map[string]*string{
				"app":     &config.Cfg.Service,
				"version": &config.Cfg.Version,
				"project": &config.Cfg.Project,
			},
		}
	}
	chart.NewServiceChart(app, config.Cfg.App+"-"+config.Cfg.Service+"-service", chartProps)
	needConfigmap := false
	if len(config.Cfg.Config) > 0 {
		for i := range config.Cfg.Config {
			if len(config.Cfg.Config[i].Data) > 0 {
				needConfigmap = true
			}
		}
	}
	if needConfigmap {
		chart.NewConfigmapChart(app, config.Cfg.App+"-"+config.Cfg.Service+"-configmap", chartProps)
	}
	chart.NewWorkloadChart(app, config.Cfg.App+"-"+config.Cfg.Service+"-workload", chartProps)
	if config.Cfg.Hpa {
		chart.NewHpaChart(app, config.Cfg.App+"-"+config.Cfg.Service+"-hpa", chartProps)
	}
	app.Synth()
}
