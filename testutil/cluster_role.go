// Copyright 2017 The prometheus-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testutil

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
)

func CreateClusterRole(kubeClient kubernetes.Interface, relativePath string) error {
	clusterRole, err := parseClusterRoleYaml(relativePath)
	if err != nil {
		return err
	}

	_, err = kubeClient.RbacV1().ClusterRoles().Get(clusterRole.Name, metav1.GetOptions{})

	if err == nil {
		// ClusterRole already exists -> Update
		_, err = kubeClient.RbacV1().ClusterRoles().Update(clusterRole)
		if err != nil {
			return err
		}

	} else {
		// ClusterRole doesn't exists -> Create
		_, err = kubeClient.RbacV1().ClusterRoles().Create(clusterRole)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteClusterRole(kubeClient kubernetes.Interface, relativePath string) error {
	clusterRole, err := parseClusterRoleYaml(relativePath)
	if err != nil {
		return err
	}

	if err := kubeClient.RbacV1().ClusterRoles().Delete(clusterRole.Name, &metav1.DeleteOptions{}); err != nil {
		return err
	}

	return nil
}

func parseClusterRoleYaml(relativePath string) (*rbacv1.ClusterRole, error) {
	var manifest *os.File
	var err error

	var clusterRole rbacv1.ClusterRole
	if manifest, err = PathToOSFile(relativPath); err != nil {
		return nil, err
	}

	decoder := yaml.NewYAMLOrJSONDecoder(manifest, 100)
	for {
		var out unstructured.Unstructured
		err = decoder.Decode(&out)
		if err != nil {
			// this would indicate it's malformed YAML.
			break
		}

		if out.GetKind() == "ClusterRole" {
			var marshaled []byte
			marshaled, err = out.MarshalJSON()
			json.Unmarshal(marshaled, &clusterRole)
			break
		}
	}

	if err != io.EOF && err != nil {
		return nil, err
	}
	return &clusterRole, nil
}
