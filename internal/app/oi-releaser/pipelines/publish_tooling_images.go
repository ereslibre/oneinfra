/**
 * Copyright 2020 Rafael Fernández López <ereslibre@ereslibre.es>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 **/

package pipelines

import (
	"fmt"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/oneinfra/oneinfra/internal/app/oi-releaser/pipelines/azure"
	"github.com/oneinfra/oneinfra/internal/pkg/constants"
)

// AzurePublishToolingImages builds the Azure publish tooling images pipeline
func AzurePublishToolingImages() error {
	pipeline := azure.Pipeline{
		Variables: map[string]string{
			"CI": "1",
		},
		Jobs: []azure.Job{},
		Trigger: &azure.Trigger{
			Branches: &azure.BranchesTrigger{
				Include: []string{"master"},
			},
			Paths: &azure.PathsTrigger{
				Include: []string{
					"RELEASE",
					".azure-pipelines/publish-tooling-images.yml",
				},
			},
			PR: &azure.PRTrigger{
				Branches: &azure.BranchesTrigger{
					Exclude: []string{"*"},
				},
			},
		},
	}
	for _, kubernetesVersion := range constants.ReleaseData.KubernetesVersions {
		pipeline.Jobs = append(
			pipeline.Jobs,
			publishContainerJob(fmt.Sprintf("kubelet-installer:%s", kubernetesVersion.Version), []string{}, defaultOptions),
		)
	}
	marshaledPipeline, err := yaml.Marshal(&pipeline)
	if err != nil {
		return err
	}
	fmt.Println("# Code generated by oi-releaser. DO NOT EDIT.")
	azurifiedYAML := strings.ReplaceAll(
		string(marshaledPipeline),
		"- _",
		"- ",
	)
	fmt.Print(azurifiedYAML)
	return nil
}
