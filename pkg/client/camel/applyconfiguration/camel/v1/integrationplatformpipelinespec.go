/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	camelv1 "github.com/apache/camel-k/v2/pkg/apis/camel/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IntegrationPlatformPipelineSpecApplyConfiguration represents an declarative configuration of the IntegrationPlatformPipelineSpec type for use
// with apply.
type IntegrationPlatformPipelineSpecApplyConfiguration struct {
	BuildConfiguration      *BuildConfigurationApplyConfiguration            `json:"buildConfiguration,omitempty"`
	PublishStrategy         *camelv1.IntegrationPlatformBuildPublishStrategy `json:"publishStrategy,omitempty"`
	RuntimeVersion          *string                                          `json:"runtimeVersion,omitempty"`
	RuntimeProvider         *camelv1.RuntimeProvider                         `json:"runtimeProvider,omitempty"`
	BaseImage               *string                                          `json:"baseImage,omitempty"`
	Registry                *RegistrySpecApplyConfiguration                  `json:"registry,omitempty"`
	BuildCatalogToolTimeout *metav1.Duration                                 `json:"buildCatalogToolTimeout,omitempty"`
	Timeout                 *metav1.Duration                                 `json:"timeout,omitempty"`
	Maven                   *MavenSpecApplyConfiguration                     `json:"maven,omitempty"`
	PublishStrategyOptions  map[string]string                                `json:"PublishStrategyOptions,omitempty"`
	MaxRunningPipelines     *int32                                           `json:"maxRunningPipelines,omitempty"`
}

// IntegrationPlatformPipelineSpecApplyConfiguration constructs an declarative configuration of the IntegrationPlatformPipelineSpec type for use with
// apply.
func IntegrationPlatformPipelineSpec() *IntegrationPlatformPipelineSpecApplyConfiguration {
	return &IntegrationPlatformPipelineSpecApplyConfiguration{}
}

// WithBuildConfiguration sets the BuildConfiguration field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the BuildConfiguration field is set to the value of the last call.
func (b *IntegrationPlatformPipelineSpecApplyConfiguration) WithBuildConfiguration(value *BuildConfigurationApplyConfiguration) *IntegrationPlatformPipelineSpecApplyConfiguration {
	b.BuildConfiguration = value
	return b
}

// WithPublishStrategy sets the PublishStrategy field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PublishStrategy field is set to the value of the last call.
func (b *IntegrationPlatformPipelineSpecApplyConfiguration) WithPublishStrategy(value camelv1.IntegrationPlatformBuildPublishStrategy) *IntegrationPlatformPipelineSpecApplyConfiguration {
	b.PublishStrategy = &value
	return b
}

// WithRuntimeVersion sets the RuntimeVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RuntimeVersion field is set to the value of the last call.
func (b *IntegrationPlatformPipelineSpecApplyConfiguration) WithRuntimeVersion(value string) *IntegrationPlatformPipelineSpecApplyConfiguration {
	b.RuntimeVersion = &value
	return b
}

// WithRuntimeProvider sets the RuntimeProvider field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RuntimeProvider field is set to the value of the last call.
func (b *IntegrationPlatformPipelineSpecApplyConfiguration) WithRuntimeProvider(value camelv1.RuntimeProvider) *IntegrationPlatformPipelineSpecApplyConfiguration {
	b.RuntimeProvider = &value
	return b
}

// WithBaseImage sets the BaseImage field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the BaseImage field is set to the value of the last call.
func (b *IntegrationPlatformPipelineSpecApplyConfiguration) WithBaseImage(value string) *IntegrationPlatformPipelineSpecApplyConfiguration {
	b.BaseImage = &value
	return b
}

// WithRegistry sets the Registry field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Registry field is set to the value of the last call.
func (b *IntegrationPlatformPipelineSpecApplyConfiguration) WithRegistry(value *RegistrySpecApplyConfiguration) *IntegrationPlatformPipelineSpecApplyConfiguration {
	b.Registry = value
	return b
}

// WithBuildCatalogToolTimeout sets the BuildCatalogToolTimeout field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the BuildCatalogToolTimeout field is set to the value of the last call.
func (b *IntegrationPlatformPipelineSpecApplyConfiguration) WithBuildCatalogToolTimeout(value metav1.Duration) *IntegrationPlatformPipelineSpecApplyConfiguration {
	b.BuildCatalogToolTimeout = &value
	return b
}

// WithTimeout sets the Timeout field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Timeout field is set to the value of the last call.
func (b *IntegrationPlatformPipelineSpecApplyConfiguration) WithTimeout(value metav1.Duration) *IntegrationPlatformPipelineSpecApplyConfiguration {
	b.Timeout = &value
	return b
}

// WithMaven sets the Maven field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Maven field is set to the value of the last call.
func (b *IntegrationPlatformPipelineSpecApplyConfiguration) WithMaven(value *MavenSpecApplyConfiguration) *IntegrationPlatformPipelineSpecApplyConfiguration {
	b.Maven = value
	return b
}

// WithPublishStrategyOptions puts the entries into the PublishStrategyOptions field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the PublishStrategyOptions field,
// overwriting an existing map entries in PublishStrategyOptions field with the same key.
func (b *IntegrationPlatformPipelineSpecApplyConfiguration) WithPublishStrategyOptions(entries map[string]string) *IntegrationPlatformPipelineSpecApplyConfiguration {
	if b.PublishStrategyOptions == nil && len(entries) > 0 {
		b.PublishStrategyOptions = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.PublishStrategyOptions[k] = v
	}
	return b
}

// WithMaxRunningPipelines sets the MaxRunningPipelines field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MaxRunningPipelines field is set to the value of the last call.
func (b *IntegrationPlatformPipelineSpecApplyConfiguration) WithMaxRunningPipelines(value int32) *IntegrationPlatformPipelineSpecApplyConfiguration {
	b.MaxRunningPipelines = &value
	return b
}
