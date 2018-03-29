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
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
)

func CreateServiceAccount(kubeClient kubernetes.Interface, namespace string, relativPath string) (finalizerFn, error) {
	finalizerFn := func() error { return DeleteServiceAccount(kubeClient, namespace, relativPath) }

	serviceAccount, err := parseServiceAccountYaml(relativPath)
	if err != nil {
		return finalizerFn, err
	}
	_, err = kubeClient.CoreV1().ServiceAccounts(namespace).Create(serviceAccount)
	if err != nil {
		return finalizerFn, err
	}

	return finalizerFn, nil
}

func parseServiceAccountYaml(relativPath string) (*v1.ServiceAccount, error) {
	var manifest *os.File
	var err error

	var serviceAccount v1.ServiceAccount
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

		if out.GetKind() == "ServiceAccount" {
			var marshaled []byte
			marshaled, err = out.MarshalJSON()
			json.Unmarshal(marshaled, &serviceAccount)
			break
		}
	}

	if err != io.EOF && err != nil {
		return nil, err
	}
	return &serviceAccount, nil
}

func DeleteServiceAccount(kubeClient kubernetes.Interface, namespace string, relativPath string) error {
	serviceAccount, err := parseServiceAccountYaml(relativPath)
	if err != nil {
		return err
	}

	return kubeClient.CoreV1().ServiceAccounts(namespace).Delete(serviceAccount.Name, nil)
}
