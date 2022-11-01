/*
 * @Author: lwnmengjing
 * @Date: 2021/10/29 10:56 下午
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2021/10/29 10:56 下午
 */

package chart

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/mss-boot-io/cd-template/imports/k8s"
	"github.com/mss-boot-io/cd-template/pkg/config"
)

func NewServiceChart(scope constructs.Construct, id string, props *cdk8s.ChartProps) cdk8s.Chart {
	chart := cdk8s.NewChart(scope, jsii.String(id), props)

	// define resources here
	ports := make([]*k8s.ServicePort, len(config.Cfg.Ports))
	for i := range config.Cfg.Ports {
		ports[i] = &k8s.ServicePort{
			Name:       &config.Cfg.Ports[i].Name,
			Port:       jsii.Number(float64(config.Cfg.Ports[i].Port)),
			TargetPort: k8s.IntOrString_FromNumber(jsii.Number(float64(config.Cfg.Ports[i].TargetPort))),
		}
	}

	k8s.NewKubeService(chart, jsii.String("service"), &k8s.KubeServiceProps{
		Metadata: &k8s.ObjectMeta{
			Name:   &config.Cfg.Service,
			Labels: props.Labels,
		},
		Spec: &k8s.ServiceSpec{
			Type:  jsii.String("ClusterIP"),
			Ports: &ports,
			Selector: &map[string]*string{
				"app": &config.Cfg.Service,
			},
		},
	})

	if !config.Cfg.ServiceAccount {
		return chart
	}
	k8s.NewKubeServiceAccount(chart, jsii.String("service-account"), &k8s.KubeServiceAccountProps{
		Metadata: &k8s.ObjectMeta{
			Name:   jsii.String(config.Cfg.GetName()),
			Labels: props.Labels,
		},
	})

	return chart
}
