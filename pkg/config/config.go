/*
 * @Author: lwnmengjing
 * @Date: 2021/10/29 10:30 下午
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2021/10/29 10:30 下午
 */

package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var Cfg Config

type Config struct {
	App                string              `json:"app" yaml:"app"`
	Project            string              `json:"project" yaml:"project"`
	Service            string              `json:"service" yaml:"service"`
	Version            string              `json:"version" yaml:"version"`
	Hpa                bool                `json:"hpa" yaml:"hpa"`
	Resources          map[string]Resource `json:"resources" yaml:"resources"`
	Replicas           uint                `json:"replicas" yaml:"replicas"`
	TestReplicas       uint                `json:"testReplicas" yaml:"testReplicas"`
	MinReplicas        uint                `json:"minReplicas" yaml:"minReplicas"`
	MaxReplicas        uint                `json:"maxReplicas" yaml:"maxReplicas"`
	ServiceAccount     bool                `json:"serviceAccount" yaml:"serviceAccount"`
	ServiceAccountName string              `json:"serviceAccountName" yaml:"serviceAccountName"`
	Image              Image               `json:"image" yaml:"image"`
	Ports              []Port              `json:"ports" yaml:"ports"`
	Metrics            Metrics             `json:"metrics" yaml:"metrics"`
	ImportEnvNames     []string            `json:"importEnvNames" yaml:"importEnvNames"`
	Config             []ConfigmapData     `json:"config" yaml:"config"`
	Secret             []SecretData        `json:"secret" yaml:"secret"`
	Storages           []Storage           `json:"storages" yaml:"storages"`
	WorkloadType       string              `json:"workloadType" yaml:"workloadType"`
	Command            []*string           `json:"command" yaml:"command"`
	Args               []*string           `json:"args" yaml:"args"`
	Containers         []Container         `json:"containers" yaml:"containers"`
	NodeSelector       map[string]*string  `json:"nodeSelector" yaml:"nodeSelector"`
}

type Storage struct {
	Name         string `json:"name" yaml:"name"`
	Path         string `json:"path" yaml:"path"`
	SubPath      string `json:"subPath" yaml:"subPath"`
	Size         string `json:"size" yaml:"size"`
	StorageClass string `json:"storageClass" yaml:"storageClass"`
}

func (c *Config) GetName() string {
	a := make([]string, 0)
	if c.Project != "" {
		a = append(a, c.Project)
	}
	if c.App != "" {
		a = append(a, c.App)
	}
	if c.Service != "" {
		a = append(a, c.Service)
	}
	return strings.Join(a, "-")
}

type Container struct {
	Image string          `json:"image" yaml:"image"`
	Name  string          `json:"name" yaml:"name"`
	Ports []ContainerPort `json:"ports" yaml:"ports"`
}

type ContainerPort struct {
	Name          string  `json:"name"`
	HostIp        string  `json:"hostIp"`
	HostPort      float64 `json:"hostPort"`
	ContainerPort float64 `json:"containerPort"`
	Protocol      string  `json:"protocol"`
}

type Resource struct {
	CPU    string `json:"cpu" yaml:"cpu"`
	Memory string `json:"memory" yaml:"memory"`
}

type Port struct {
	Port       uint   `json:"port" yaml:"port"`
	TargetPort uint   `json:"targetPort" yaml:"targetPort"`
	Name       string `json:"name" yaml:"name"`
}

type Image struct {
	Path    string   `json:"path" yaml:"path"`
	Tag     string   `json:"tag" yaml:"tag"`
	Secrets []string `json:"secrets" yaml:"secrets"`
}

func (e Image) String() string {
	if e.Tag == "" {
		return e.Path
	}
	return e.Path + ":" + e.Tag
}

type Metrics struct {
	Scrape bool   `json:"scrape" yaml:"scrape"`
	Path   string `json:"path" yaml:"path"`
	Port   uint   `json:"port" yaml:"port"`
}

type ConfigmapData struct {
	Name    string            `json:"name" yaml:"name"`
	Path    string            `json:"path" yaml:"path"`
	EnvName string            `json:"envName" yaml:"envName"`
	Key     string            `json:"key" yaml:"key"`
	Data    map[string]string `json:"data" yaml:"data"`
}

type SecretData struct {
	Name    string `json:"name" yaml:"name"`
	Path    string `json:"path" yaml:"path"`
	EnvName string `json:"envName" yaml:"envName"`
	Key     string `json:"key" yaml:"key"`
}

var (
	app              = flag.String("app", "", "application")
	project          = flag.String("project", "", "project")
	service          = flag.String("service", "", "service")
	version          = flag.String("version", "v1", "service version")
	httpPort         = flag.Uint("httpPort", 0, "http server listen port")
	grpcPort         = flag.Uint("grpcPort", 0, "grpc server listen port")
	image            = flag.String("image", "", "image:tag")
	importEnvNames   = flag.String("importEnvNames", "", "import env names, split ','")
	configDataFiles  = flag.String("configDataFiles", "", "config data file path, multi split ','")
	configPath       = flag.String("configPath", "", "application config path")
	configmapName    = flag.String("configmapName", "", "exist configmap name")
	secretName       = flag.String("secretName", "", "exist secret name")
	replicas         = flag.Uint("replicas", 1, "replicas")
	workloadType     = flag.String("workloadType", "deployment", "workload type, e.g. deployment, statefulset")
	hpa              = flag.Bool("hpa", false, "enable hpa")
	metricsScrape    = flag.Bool("metricsScrape", false, "enable metrics export")
	imagePullSecrets = flag.String("imagePullSecrets", "", "image pull secrets, multi split ','")
	nodeSelector     = flag.String("nodeSelector", "", "node selector, e.g. app0=demo0,app1=demo1")
	_                = flag.String("namespace", "", "")
)

// NewConfig set config
func NewConfig(path *string) {
	viper.SetConfigType("yaml")
	f, err := os.OpenFile(*path, os.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("open config file error: %v", err)
	}

	err = viper.ReadConfig(f)
	if err != nil {
		log.Fatalf("read config file error: %v", err)
	}

	//set default config
	viper.SetDefault("app", *app)
	viper.SetDefault("service", *service)
	viper.SetDefault("project", *project)
	viper.SetDefault("image.path", *image)
	viper.SetDefault("version", *version)
	viper.SetDefault("image.secrets", strings.Split(*imagePullSecrets, ","))

	err = viper.Unmarshal(&Cfg)
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	//config.DefaultConfig, err = config.NewConfig()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//if path != nil && *path != "" {
	//	err = config.LoadFile(*path)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//	err = config.Scan(&Cfg)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//}

	//Cfg.App = config.Get("app").String(*app)
	//Cfg.Service = config.Get("service").String(*service)
	//Cfg.Project = config.Get("project").String(*project)
	if len(Cfg.Ports) == 0 && (*httpPort > 0 || *grpcPort > 0) {
		Cfg.Ports = make([]Port, 0)
		if *httpPort > 0 {
			Cfg.Ports = append(Cfg.Ports, Port{
				Port:       *httpPort,
				Name:       "http",
				TargetPort: *httpPort,
			})
		}
		if *grpcPort > 0 {
			Cfg.Ports = append(Cfg.Ports, Port{
				Port:       *grpcPort,
				Name:       "grpc",
				TargetPort: *grpcPort,
			})
		}

	}
	if *metricsScrape || Cfg.Metrics.Scrape {
		Cfg.Metrics.Scrape = true
		if Cfg.Metrics.Port == 0 {
			Cfg.Metrics.Port = 5000
		}
		if Cfg.Metrics.Path == "" {
			Cfg.Metrics.Path = "/metrics"
		}
	}
	//Cfg.Image.Path = config.Get("image", "path").String(*image)
	//Cfg.Version = config.Get("version").String(*version)
	if len(Cfg.ImportEnvNames) == 0 {
		Cfg.ImportEnvNames = strings.Split(*importEnvNames, ",")
	}
	if configDataFiles != nil && *configDataFiles != "" {
		if Cfg.Config == nil {
			Cfg.Config = make([]ConfigmapData, 0)
		}
		if Cfg.Secret == nil {
			Cfg.Secret = make([]SecretData, 0)
		}
		configData := ConfigmapData{
			Data: make(map[string]string),
		}
		var secretData SecretData
		for _, p := range strings.Split(*configDataFiles, ",") {
			rb, err := os.ReadFile(p)
			if err != nil {
				log.Fatalln(err)
			}
			configData.Data[filepath.Base(p)] = string(rb)
		}
		if configPath != nil && *configPath != "" {
			configData.Path = *configPath
		}
		if configmapName != nil && *configmapName != "" {
			configData.Name = *configmapName
		}
		if secretName != nil && *secretName != "" {
			secretData.Name = *secretName
		}
		Cfg.Config = append(Cfg.Config, configData)
		Cfg.Secret = append(Cfg.Secret, secretData)
	}
	if Cfg.Replicas < 1 {
		Cfg.Replicas = *replicas
	}
	if Cfg.MinReplicas < 1 {
		Cfg.MinReplicas = 3
	}
	if Cfg.MaxReplicas < Cfg.Replicas {
		Cfg.MaxReplicas = 5
	}
	if workloadType != nil && (Cfg.WorkloadType == "" || *workloadType != "deployment") {
		Cfg.WorkloadType = *workloadType
	}
	if hpa != nil && *hpa {
		Cfg.Hpa = *hpa
	}

	if len(Cfg.Resources) == 0 && Cfg.Hpa {
		Cfg.Resources = map[string]Resource{
			"limits": {
				CPU:    "1",
				Memory: "1Gi",
			},
			"requests": {
				CPU:    "500m",
				Memory: "800Mi",
			},
		}
	}
	if nodeSelector != nil && *nodeSelector != "" {
		Cfg.NodeSelector = stringToMap(*nodeSelector)
	}
}

func stringToMap(s string) map[string]*string {
	m := make(map[string]*string)
	for _, v := range strings.Split(s, ",") {
		vs := strings.Split(v, "=")
		if len(vs) <= 1 {
			continue
		}
		m[vs[0]] = &vs[1]
	}
	return m
}
