package chart

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	"github.com/mss-boot-io/cd-template/imports/k8s"
	"github.com/mss-boot-io/cd-template/pkg/config"
)

func NewPvcChart(scope constructs.Construct, id string, props *cdk8s.ChartProps, storage config.Storage) cdk8s.Chart {
	chart := cdk8s.NewChart(scope, jsii.String(id), props)
	pvcProps := &k8s.KubePersistentVolumeClaimProps{
		Metadata: &k8s.ObjectMeta{
			Labels: props.Labels,
			Name:   jsii.String(storage.Name),
		},
		Spec: &k8s.PersistentVolumeClaimSpec{
			AccessModes: &[]*string{
				jsii.String("ReadWriteOnce"),
			},
			Resources: &k8s.ResourceRequirements{
				Requests: &map[string]k8s.Quantity{
					"storage": k8s.Quantity_FromString(&storage.Size),
				},
			},
			VolumeMode: jsii.String("Filesystem"),
		},
	}
	if storage.StorageClass != "" {
		pvcProps.Spec.StorageClassName = &storage.StorageClass
	}
	k8s.NewKubePersistentVolumeClaim(chart, jsii.String("pvc"), pvcProps)

	return chart
}
