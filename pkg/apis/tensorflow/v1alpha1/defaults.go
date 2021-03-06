package v1alpha1

import (
	"github.com/golang/protobuf/proto"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// SetDefaults_TFJob sets any unspecified values to defaults
func SetDefaults_TFJob(obj *TFJob) {
	c := &obj.Spec

	if c.TFImage == "" {
		c.TFImage = DefaultTFImage
	}

	// Check that each replica has a TensorFlow container.
	for _, r := range c.ReplicaSpecs {

		if r.TFPort == nil {
			r.TFPort = proto.Int32(TFPort)
		}

		if string(r.TFReplicaType) == "" {
			r.TFReplicaType = MASTER
		}

		if r.Replicas == nil {
			r.Replicas = proto.Int32(Replicas)
		}

		//Set the default configuration for a PS server if the user didn't specify a PodTemplateSpec
		if r.Template == nil && r.TFReplicaType == PS {
			setDefault_PSPodTemplateSpec(r, c.TFImage)
		}
	}
	if c.TerminationPolicy == nil {
		c.TerminationPolicy = &TerminationPolicySpec{
			Chief: &ChiefSpec{
				ReplicaName:  "MASTER",
				ReplicaIndex: 0,
			},
		}
	}

}

func setDefault_PSPodTemplateSpec(r *TFReplicaSpec, tfImage string) {
	r.IsDefaultPS = true
	r.Template = &v1.PodTemplateSpec{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				v1.Container{
					Image: tfImage,
					Name:  "tensorflow",
					VolumeMounts: []v1.VolumeMount{
						v1.VolumeMount{
							Name:      "ps-config-volume",
							MountPath: "/ps-server",
						},
					},
				},
			},
			RestartPolicy: v1.RestartPolicyOnFailure,
		},
	}
}
