/*
 * @Author: lwnmengjing
 * @Date: 2021/10/29 11:21 下午
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2021/10/29 11:21 下午
 */

package chart

import (
	"os"
	"strconv"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	"github.com/mss-boot-io/cd-template/imports/k8s"
	"github.com/mss-boot-io/cd-template/pkg/config"
)

func NewWorkloadChart(scope constructs.Construct, id string, props *cdk8s.ChartProps) cdk8s.Chart {
	chart := cdk8s.NewChart(scope, jsii.String(id), props)
	ports := make([]*k8s.ContainerPort, 0)
	//port
	for i := range config.Cfg.Ports {
		ports = append(ports, &k8s.ContainerPort{
			ContainerPort: jsii.Number(float64(config.Cfg.Ports[i].Port)),
		})
	}
	//env
	env := make([]*k8s.EnvVar, 0)
	for i := range config.Cfg.ImportEnvNames {
		if config.Cfg.ImportEnvNames[i] == "" {
			continue
		}
		v := os.Getenv(config.Cfg.ImportEnvNames[i])
		env = append(env, &k8s.EnvVar{
			Name:  &config.Cfg.ImportEnvNames[i],
			Value: &v,
		})
	}
	env = append(env, &k8s.EnvVar{
		Name: jsii.String("NODE_NAME"),
		ValueFrom: &k8s.EnvVarSource{
			FieldRef: &k8s.ObjectFieldSelector{
				FieldPath: jsii.String("metadata.name"),
			},
		},
	}, &k8s.EnvVar{
		Name: jsii.String("STAGE"),
		ValueFrom: &k8s.EnvVarSource{
			FieldRef: &k8s.ObjectFieldSelector{
				FieldPath: jsii.String("metadata.namespace"),
			},
		},
	})
	//container
	containers := make([]*k8s.Container, 0)
	if len(config.Cfg.Containers) > 0 {
		for i := range config.Cfg.Containers {
			containerPorts := make([]*k8s.ContainerPort, 0)
			for j := range config.Cfg.Containers[i].Ports {
				port := k8s.ContainerPort{
					Name:          &config.Cfg.Containers[i].Ports[j].Name,
					HostIp:        &config.Cfg.Containers[i].Ports[j].HostIp,
					HostPort:      &config.Cfg.Containers[i].Ports[j].HostPort,
					ContainerPort: &config.Cfg.Containers[i].Ports[j].HostPort,
					Protocol:      &config.Cfg.Containers[i].Ports[j].Protocol,
				}
				containerPorts = append(containerPorts, &port)
			}
			containers = append(containers, &k8s.Container{
				Name:  &config.Cfg.Containers[i].Name,
				Image: &config.Cfg.Containers[i].Image,
				Ports: &containerPorts,
				Env:   &env,
			})
		}
	}
	//config
	volumeMounts := make([]*k8s.VolumeMount, 0)
	volumes := make([]*k8s.Volume, 0)
	if len(config.Cfg.Config) > 0 {
		readOnly := true
		for i := range config.Cfg.Config {
			volumes = append(volumes, &k8s.Volume{
				Name: &config.Cfg.Config[i].Name,
				ConfigMap: &k8s.ConfigMapVolumeSource{
					Name: &config.Cfg.Config[i].Name,
				},
			})

			volumeMounts = append(volumeMounts, &k8s.VolumeMount{
				MountPath: &config.Cfg.Config[i].Path,
				Name:      &config.Cfg.Config[i].Name,
				ReadOnly:  &readOnly,
			})
		}
	}
	for i := range config.Cfg.Storages {
		volumes = append(volumes, &k8s.Volume{
			Name: &config.Cfg.Storages[i].Name,
			PersistentVolumeClaim: &k8s.PersistentVolumeClaimVolumeSource{
				ClaimName: &config.Cfg.Storages[i].Name,
			},
		})
		vm := &k8s.VolumeMount{
			MountPath: &config.Cfg.Storages[i].Path,
			Name:      &config.Cfg.Storages[i].Name,
		}
		if config.Cfg.Storages[i].SubPath != "" {
			vm.SubPath = &config.Cfg.Storages[i].SubPath
		}
		volumeMounts = append(volumeMounts, vm)
	}

	var serviceAccountName *string
	if config.Cfg.ServiceAccount {
		serviceAccountName = jsii.String(config.Cfg.GetName())
	}
	if config.Cfg.ServiceAccountName != "" {
		serviceAccountName = jsii.String(config.Cfg.ServiceAccountName)
	}
	var command *[]*string
	if len(config.Cfg.Command) > 0 {

		command = &config.Cfg.Command
	}
	var args *[]*string
	if len(config.Cfg.Args) > 0 {
		args = &config.Cfg.Args
	}

	var resources k8s.ResourceRequirements
	if len(config.Cfg.Resources) > 0 {
		for k, r := range config.Cfg.Resources {
			switch k {
			case "limits":
				resources.Limits = &map[string]k8s.Quantity{
					"cpu":    k8s.Quantity_FromString(&r.CPU),
					"memory": k8s.Quantity_FromString(&r.Memory),
				}
			case "requests":
				resources.Requests = &map[string]k8s.Quantity{
					"cpu":    k8s.Quantity_FromString(&r.CPU),
					"memory": k8s.Quantity_FromString(&r.Memory),
				}
			}
		}
	}
	annotations := make(map[string]*string)
	if config.Cfg.Metrics.Scrape {
		annotations["prometheus.io/scrape"] = jsii.String("true")
		annotations["prometheus.io/port"] = jsii.String(strconv.Itoa(int(config.Cfg.Metrics.Port)))
		annotations["prometheus.io/path"] = jsii.String(config.Cfg.Metrics.Path)
	}
	var replicas *float64
	if !config.Cfg.Hpa && config.Cfg.Replicas > 0 {
		replicas = jsii.Number(float64(config.Cfg.Replicas))
	}
	imagePullSecrets := make([]*k8s.LocalObjectReference, len(config.Cfg.Image.Secrets))
	for i := range config.Cfg.Image.Secrets {
		imagePullSecrets[i] = &k8s.LocalObjectReference{Name: &config.Cfg.Image.Secrets[i]}
	}
	// default container
	containers = append(containers, &k8s.Container{
		Name:         jsii.String(config.Cfg.Service),
		Image:        jsii.String(config.Cfg.Image.String()),
		Ports:        &ports,
		Env:          &env,
		VolumeMounts: &volumeMounts,
		Command:      command,
		Args:         args,
		Resources:    &resources,
	})

	switch config.Cfg.WorkloadType {
	case "statefulset":
		k8s.NewKubeStatefulSet(chart, jsii.String("statefulset"), &k8s.KubeStatefulSetProps{
			Metadata: &k8s.ObjectMeta{
				Name:   &config.Cfg.Service,
				Labels: props.Labels,
			},
			Spec: &k8s.StatefulSetSpec{
				ServiceName: &config.Cfg.Service,
				Replicas:    replicas,
				Selector: &k8s.LabelSelector{
					MatchLabels: props.Labels,
				},
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels:      props.Labels,
						Annotations: &annotations,
					},
					Spec: &k8s.PodSpec{
						ServiceAccountName: serviceAccountName,
						Containers:         &containers,
						Volumes:            &volumes,
						ImagePullSecrets:   &imagePullSecrets,
					},
				},
			},
		})
	default:
		k8s.NewKubeDeployment(chart, jsii.String("deployment"), &k8s.KubeDeploymentProps{
			Metadata: &k8s.ObjectMeta{
				Name:   &config.Cfg.Service,
				Labels: props.Labels,
			},
			Spec: &k8s.DeploymentSpec{
				Replicas: replicas,
				Selector: &k8s.LabelSelector{
					MatchLabels: props.Labels,
				},
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels:      props.Labels,
						Annotations: &annotations,
					},
					Spec: &k8s.PodSpec{
						ServiceAccountName: serviceAccountName,
						Containers:         &containers,
						Volumes:            &volumes,
						ImagePullSecrets:   &imagePullSecrets,
					},
				},
			},
		})
	}
	return chart
}
