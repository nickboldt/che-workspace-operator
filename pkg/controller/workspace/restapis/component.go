//
// Copyright (c) 2019-2020 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//

package restapis

import (
	"github.com/che-incubator/che-workspace-operator/pkg/apis/workspace/v1alpha1"
	devworkspace "github.com/devfile/kubernetes-api/pkg/apis/workspaces/v1alpha1"
	"github.com/che-incubator/che-workspace-operator/pkg/common"
	"github.com/che-incubator/che-workspace-operator/pkg/config"
	corev1 "k8s.io/api/core/v1"
)

const cheRestAPIsName = "che-rest-apis"
const cheRestApisPort = 9999

func GetCheRestApisComponent(workspaceName, workspaceId, namespace string) v1alpha1.ComponentDescription {
	container := corev1.Container{
		Image:           config.ControllerCfg.GetCheRestApisDockerImage(),
		ImagePullPolicy: corev1.PullPolicy(config.ControllerCfg.GetSidecarPullPolicy()),
		Name:            cheRestAPIsName,
		Ports: []corev1.ContainerPort{
			{
				ContainerPort: int32(cheRestApisPort),
				Protocol:      corev1.ProtocolTCP,
			},
		},
		Env: []corev1.EnvVar{
			{
				Name:  "CHE_WORKSPACE_NAME",
				Value: workspaceName,
			},
			{
				Name:  "CHE_WORKSPACE_ID",
				Value: workspaceId,
			},
			{
				Name:  "CHE_WORKSPACE_NAMESPACE",
				Value: namespace,
			},
			{
				Name:  "CHE_WORKSPACE_RUNTIME_JSON_PATH",
				Value: config.RestAPIsRuntimeVolumePath + config.RestAPIsRuntimeJSONFilename,
			},
			{
				Name:  "CHE_WORKSPACE_DEVFILE_YAML_PATH",
				Value: config.RestAPIsRuntimeVolumePath + config.RestAPIsDevfileYamlFilename,
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      config.RestAPIsVolumeName,
				ReadOnly:  true,
				MountPath: config.RestAPIsRuntimeVolumePath,
			},
		},
	}

	volumeMode := corev1.ConfigMapVolumeSourceDefaultMode
	configmapVolume := corev1.Volume{
		Name: config.RestAPIsVolumeName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: common.CheRestAPIsConfigmapName(workspaceId),
				},
				DefaultMode: &volumeMode,
			},
		},
	}

	return v1alpha1.ComponentDescription{
		Name: cheRestAPIsName,
		PodAdditions: v1alpha1.PodAdditions{
			Containers: []corev1.Container{container},
			Volumes:    []corev1.Volume{configmapVolume},
		},
		ComponentMetadata: v1alpha1.ComponentMetadata{
			Containers: map[string]v1alpha1.ContainerDescription{
				cheRestAPIsName: {
					Attributes: map[string]string{
						config.RestApisContainerSourceAttribute: config.RestApisRecipeSourceToolAttribute,
					},
					Ports: []int{cheRestApisPort},
				},
			},
			Endpoints: []devworkspace.Endpoint{
				{
					Attributes: map[string]string{
						string(v1alpha1.PUBLIC_ENDPOINT_ATTRIBUTE):   "false",
						string(v1alpha1.PROTOCOL_ENDPOINT_ATTRIBUTE): "tcp",
					},
					Name: cheRestAPIsName,
					TargetPort: cheRestApisPort,
				},
			},
		},
	}
}
