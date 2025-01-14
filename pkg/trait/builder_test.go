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

package trait

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	v1 "github.com/apache/camel-k/v2/pkg/apis/camel/v1"
	"github.com/apache/camel-k/v2/pkg/util/camel"
	"github.com/apache/camel-k/v2/pkg/util/defaults"
	"github.com/apache/camel-k/v2/pkg/util/kubernetes"
	"github.com/apache/camel-k/v2/pkg/util/test"
)

func TestBuilderTraitNotAppliedBecauseOfNilKit(t *testing.T) {
	environments := []*Environment{
		createBuilderTestEnv(v1.IntegrationPlatformClusterOpenShift, v1.IntegrationPlatformBuildPublishStrategyS2I, v1.BuildStrategyRoutine),
		createBuilderTestEnv(v1.IntegrationPlatformClusterKubernetes, v1.IntegrationPlatformBuildPublishStrategyKaniko, v1.BuildStrategyRoutine),
	}

	for _, e := range environments {
		e := e // pin
		e.IntegrationKit = nil

		t.Run(string(e.Platform.Status.Cluster), func(t *testing.T) {
			err := NewBuilderTestCatalog().apply(e)

			assert.Nil(t, err)
			assert.NotEmpty(t, e.ExecutedTraits)
			assert.Nil(t, e.GetTrait("builder"))
			assert.Empty(t, e.Pipeline)
		})
	}
}

func TestBuilderTraitNotAppliedBecauseOfNilPhase(t *testing.T) {
	environments := []*Environment{
		createBuilderTestEnv(v1.IntegrationPlatformClusterOpenShift, v1.IntegrationPlatformBuildPublishStrategyS2I, v1.BuildStrategyRoutine),
		createBuilderTestEnv(v1.IntegrationPlatformClusterKubernetes, v1.IntegrationPlatformBuildPublishStrategyKaniko, v1.BuildStrategyRoutine),
	}

	for _, e := range environments {
		e := e // pin
		e.IntegrationKit.Status.Phase = v1.IntegrationKitPhaseInitialization

		t.Run(string(e.Platform.Status.Cluster), func(t *testing.T) {
			err := NewBuilderTestCatalog().apply(e)

			assert.Nil(t, err)
			assert.NotEmpty(t, e.ExecutedTraits)
			assert.Nil(t, e.GetTrait("builder"))
			assert.Empty(t, e.Pipeline)
		})
	}
}

func TestS2IBuilderTrait(t *testing.T) {
	env := createBuilderTestEnv(v1.IntegrationPlatformClusterOpenShift, v1.IntegrationPlatformBuildPublishStrategyS2I, v1.BuildStrategyRoutine)
	err := NewBuilderTestCatalog().apply(env)

	assert.Nil(t, err)
	assert.NotEmpty(t, env.ExecutedTraits)
	assert.NotNil(t, env.GetTrait("builder"))
	assert.NotEmpty(t, env.Pipeline)
	assert.Len(t, env.Pipeline, 2)
	assert.NotNil(t, env.Pipeline[0].Builder)
	assert.NotNil(t, env.Pipeline[1].S2i)
}

func TestKanikoBuilderTrait(t *testing.T) {
	env := createBuilderTestEnv(v1.IntegrationPlatformClusterKubernetes, v1.IntegrationPlatformBuildPublishStrategyKaniko, v1.BuildStrategyRoutine)
	err := NewBuilderTestCatalog().apply(env)

	assert.Nil(t, err)
	assert.NotEmpty(t, env.ExecutedTraits)
	assert.NotNil(t, env.GetTrait("builder"))
	assert.NotEmpty(t, env.Pipeline)
	assert.Len(t, env.Pipeline, 2)
	assert.NotNil(t, env.Pipeline[0].Builder)
	assert.NotNil(t, env.Pipeline[1].Kaniko)
}

func createBuilderTestEnv(cluster v1.IntegrationPlatformCluster, strategy v1.IntegrationPlatformBuildPublishStrategy, buildStrategy v1.BuildStrategy) *Environment {
	c, err := camel.DefaultCatalog()
	if err != nil {
		panic(err)
	}
	client, _ := test.NewFakeClient()
	res := &Environment{
		Ctx:          context.TODO(),
		CamelCatalog: c,
		Catalog:      NewCatalog(nil),
		Client:       client,
		Integration: &v1.Integration{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: "ns",
			},
			Status: v1.IntegrationStatus{
				Phase: v1.IntegrationPhaseDeploying,
			},
		},
		IntegrationKit: &v1.IntegrationKit{
			Status: v1.IntegrationKitStatus{
				Phase: v1.IntegrationKitPhaseBuildSubmitted,
			},
		},
		Platform: &v1.IntegrationPlatform{
			Spec: v1.IntegrationPlatformSpec{
				Cluster: cluster,
				Build: v1.IntegrationPlatformBuildSpec{
					PublishStrategy:        strategy,
					Registry:               v1.RegistrySpec{Address: "registry"},
					RuntimeVersion:         defaults.DefaultRuntimeVersion,
					RuntimeProvider:        v1.RuntimeProviderQuarkus,
					PublishStrategyOptions: map[string]string{},
					BuildConfiguration: v1.BuildConfiguration{
						Strategy: buildStrategy,
					},
				},
			},
			Status: v1.IntegrationPlatformStatus{
				Phase: v1.IntegrationPlatformPhaseReady,
			},
		},
		EnvVars:        make([]corev1.EnvVar, 0),
		ExecutedTraits: make([]Trait, 0),
		Resources:      kubernetes.NewCollection(),
	}

	res.Platform.ResyncStatusFullConfig()

	return res
}

func NewBuilderTestCatalog() *Catalog {
	return NewCatalog(nil)
}

func TestMavenPropertyBuilderTrait(t *testing.T) {
	env := createBuilderTestEnv(v1.IntegrationPlatformClusterKubernetes, v1.IntegrationPlatformBuildPublishStrategyKaniko, v1.BuildStrategyRoutine)
	builderTrait := createNominalBuilderTraitTest()
	builderTrait.Properties = append(builderTrait.Properties, "build-time-prop1=build-time-value1")

	err := builderTrait.Apply(env)

	assert.Nil(t, err)
	assert.Equal(t, "build-time-value1", env.Pipeline[0].Builder.Maven.Properties["build-time-prop1"])
}

func createNominalBuilderTraitTest() *builderTrait {
	builderTrait, _ := newBuilderTrait().(*builderTrait)
	builderTrait.Enabled = pointer.Bool(true)

	return builderTrait
}

func TestCustomTaskBuilderTrait(t *testing.T) {
	env := createBuilderTestEnv(v1.IntegrationPlatformClusterKubernetes, v1.IntegrationPlatformBuildPublishStrategySpectrum, v1.BuildStrategyPod)
	builderTrait := createNominalBuilderTraitTest()
	builderTrait.Tasks = append(builderTrait.Tasks, "test;alpine;ls")

	err := builderTrait.Apply(env)

	assert.Nil(t, err)
	builderTask := findCustomTaskByName(env.Pipeline, "builder")
	publisherTask := findCustomTaskByName(env.Pipeline, "spectrum")
	customTask := findCustomTaskByName(env.Pipeline, "test")
	assert.NotNil(t, customTask)
	assert.NotNil(t, builderTask)
	assert.NotNil(t, publisherTask)
	assert.Equal(t, 3, len(env.Pipeline))
	assert.Equal(t, "test", customTask.Custom.Name)
	assert.Equal(t, "alpine", customTask.Custom.ContainerImage)
	assert.Equal(t, "ls", customTask.Custom.ContainerCommand)
}

func TestCustomTaskBuilderTraitValidStrategyOverride(t *testing.T) {
	env := createBuilderTestEnv(v1.IntegrationPlatformClusterKubernetes, v1.IntegrationPlatformBuildPublishStrategySpectrum, v1.BuildStrategyRoutine)
	builderTrait := createNominalBuilderTraitTest()
	builderTrait.Tasks = append(builderTrait.Tasks, "test;alpine;ls")
	builderTrait.Strategy = "pod"

	err := builderTrait.Apply(env)

	assert.Nil(t, err)

	builderTask := findCustomTaskByName(env.Pipeline, "builder")
	publisherTask := findCustomTaskByName(env.Pipeline, "spectrum")
	customTask := findCustomTaskByName(env.Pipeline, "test")
	assert.NotNil(t, customTask)
	assert.NotNil(t, builderTask)
	assert.NotNil(t, publisherTask)
	assert.Equal(t, 3, len(env.Pipeline))
	assert.Equal(t, "test", customTask.Custom.Name)
	assert.Equal(t, "alpine", customTask.Custom.ContainerImage)
	assert.Equal(t, "ls", customTask.Custom.ContainerCommand)
}

func TestCustomTaskBuilderTraitInvalidStrategy(t *testing.T) {
	env := createBuilderTestEnv(v1.IntegrationPlatformClusterKubernetes, v1.IntegrationPlatformBuildPublishStrategySpectrum, v1.BuildStrategyRoutine)
	builderTrait := createNominalBuilderTraitTest()
	builderTrait.Tasks = append(builderTrait.Tasks, "test;alpine;ls")

	err := builderTrait.Apply(env)

	assert.NotNil(t, err)
	assert.Equal(t, env.IntegrationKit.Status.Phase, v1.IntegrationKitPhaseError)
	assert.Equal(t, env.IntegrationKit.Status.Conditions[0].Status, corev1.ConditionFalse)
	assert.Equal(t, env.IntegrationKit.Status.Conditions[0].Type, v1.IntegrationKitConditionType("IntegrationKitTasksValid"))
}

func TestCustomTaskBuilderTraitInvalidStrategyOverride(t *testing.T) {
	env := createBuilderTestEnv(v1.IntegrationPlatformClusterKubernetes, v1.IntegrationPlatformBuildPublishStrategySpectrum, v1.BuildStrategyPod)
	builderTrait := createNominalBuilderTraitTest()
	builderTrait.Tasks = append(builderTrait.Tasks, "test;alpine;ls")
	builderTrait.Strategy = "routine"

	err := builderTrait.Apply(env)

	assert.NotNil(t, err)
	assert.Equal(t, env.IntegrationKit.Status.Phase, v1.IntegrationKitPhaseError)
	assert.Equal(t, env.IntegrationKit.Status.Conditions[0].Status, corev1.ConditionFalse)
	assert.Equal(t, env.IntegrationKit.Status.Conditions[0].Type, v1.IntegrationKitConditionType("IntegrationKitTasksValid"))
}

func findCustomTaskByName(tasks []v1.Task, name string) v1.Task {
	for _, t := range tasks {
		if t.Custom != nil && t.Custom.Name == name {
			return t
		}
	}
	return v1.Task{}
}

func TestMavenProfileBuilderTrait(t *testing.T) {
	env := createBuilderTestEnv(v1.IntegrationPlatformClusterKubernetes, v1.IntegrationPlatformBuildPublishStrategyKaniko, v1.BuildStrategyRoutine)
	builderTrait := createNominalBuilderTraitTest()
	builderTrait.MavenProfile = "configmap:maven-profile/owasp-profile.xml"

	err := builderTrait.Apply(env)

	assert.Nil(t, err)

	assert.Equal(t, v1.ValueSource{
		ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: "maven-profile",
			},
			Key: "owasp-profile.xml",
		},
	}, env.Pipeline[0].Builder.Maven.MavenSpec.Profile)
}
