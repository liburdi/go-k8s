package go_k8s

import (
	"context"
	"errors"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var Job = &job{K8sClient: nil}

type job struct {
	K8sClient client.Client
}

func (j *job) Run(name, ns, image, cmd string, args []string) error {
	if j.K8sClient == nil {
		return errors.New("k8s client should not be nil")
	}
	containers := make([]corev1.Container, 0, 1)
	container := getContainers(name, image, cmd, args)
	containers = append(containers, container)
	err := j.K8sClient.Create(context.Background(),
		&batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: ns,
				Labels: map[string]string{
					"app.kubernetes.io/name":     "job",
					"app.kubernetes.io/instance": name,
				},
			},
			Spec: batchv1.JobSpec{
				BackoffLimit:            Int32Ptr(1),
				ActiveDeadlineSeconds:   Int64Ptr(100),
				TTLSecondsAfterFinished: Int32Ptr(0),
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app.kubernetes.io/name":     "job",
							"app.kubernetes.io/instance": name},
						Annotations: map[string]string{
							"eks.tke.cloud.tencent.com/cpu-type":                  "SA3,SA2,S5",
							"eks.tke.cloud.tencent.com/root-cbs-size":             "50Gi",
							"eks.tke.cloud.tencent.com/use-image-cache":           "auto",
							"eks.tke.cloud.tencent.com/reserve-sandbox-duration":  "20m",
							"eks.tke.cloud.tencent.com/reserve-succeeded-sandbox": "true",
							"eks.tke.cloud.tencent.com/reserve-task-shorter-than": "2m",
						},
					},
					Spec: corev1.PodSpec{
						ShareProcessNamespace:         BoolPtr(true),
						TerminationGracePeriodSeconds: Int64Ptr(30),
						ServiceAccountName:            "go-k8s-operator",
						RestartPolicy:                 corev1.RestartPolicyNever,
						NodeSelector:                  map[string]string{"node_type": "eks-job"},

						ImagePullSecrets: []corev1.LocalObjectReference{
							{
								Name: "registry-secrete",
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: "filebeat-log-volume-newlogs",
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{},
								},
							},
						},
						Tolerations: []corev1.Toleration{
							{
								Key:      "eks.tke.cloud.tencent.com/eklet",
								Operator: corev1.TolerationOpExists,
								Effect:   corev1.TaintEffectNoSchedule,
							},
						},
						Containers: containers,
					},
				},
			},
		}, &client.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func getContainers(name, image, cmd string, args []string) corev1.Container {
	return corev1.Container{
		Name:            "job",
		Image:           image,
		Command:         []string{cmd},
		Args:            args,
		ImagePullPolicy: corev1.PullAlways,
		Env: []corev1.EnvVar{
			{
				Name: "POD_NAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.name",
					},
				},
			},
			{
				Name: "POD_NAMESPACE",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.namespace",
					},
				},
			},
			{
				Name: "NODE_NAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "spec.nodeName",
					},
				},
			},
			{
				Name: "POD_IP",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "status.podIP",
					},
				},
			},
		},
		EnvFrom: []corev1.EnvFromSource{
			{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "env-config",
					},
				},
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "filebeat-log-volume-newlogs",
				MountPath: "/data/log/service",
			},
		},
		Resources: corev1.ResourceRequirements{ // TODO from tps config or metrics
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("100m"),
				corev1.ResourceMemory: resource.MustParse("150Mi"),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("2"),
				corev1.ResourceMemory: resource.MustParse("4Gi"),
			},
		},
	}
}
