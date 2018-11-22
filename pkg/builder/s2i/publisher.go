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

package s2i

import (
	"io/ioutil"
	"time"

	"github.com/apache/camel-k/pkg/builder"

	"github.com/apache/camel-k/pkg/util/kubernetes"
	"github.com/apache/camel-k/pkg/util/kubernetes/customclient"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	buildv1 "github.com/openshift/api/build/v1"
	imagev1 "github.com/openshift/api/image/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pkg/errors"
)

// Publisher --
func Publisher(ctx *builder.Context) error {
	bc := buildv1.BuildConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: buildv1.SchemeGroupVersion.String(),
			Kind:       "BuildConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "camel-k-" + ctx.Request.Meta.Name,
			Namespace: ctx.Namespace,
		},
		Spec: buildv1.BuildConfigSpec{
			CommonSpec: buildv1.CommonSpec{
				Source: buildv1.BuildSource{
					Type: buildv1.BuildSourceBinary,
				},
				Strategy: buildv1.BuildStrategy{
					SourceStrategy: &buildv1.SourceBuildStrategy{
						From: v1.ObjectReference{
							Kind: "DockerImage",
							Name: ctx.Image,
						},
					},
				},
				Output: buildv1.BuildOutput{
					To: &v1.ObjectReference{
						Kind: "ImageStreamTag",
						Name: "camel-k-" + ctx.Request.Meta.Name + ":" + ctx.Request.Meta.ResourceVersion,
					},
				},
			},
		},
	}

	sdk.Delete(&bc)
	err := sdk.Create(&bc)
	if err != nil {
		return errors.Wrap(err, "cannot create build config")
	}

	is := imagev1.ImageStream{
		TypeMeta: metav1.TypeMeta{
			APIVersion: imagev1.SchemeGroupVersion.String(),
			Kind:       "ImageStream",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "camel-k-" + ctx.Request.Meta.Name,
			Namespace: ctx.Namespace,
		},
		Spec: imagev1.ImageStreamSpec{
			LookupPolicy: imagev1.ImageLookupPolicy{
				Local: true,
			},
		},
	}

	sdk.Delete(&is)
	err = sdk.Create(&is)
	if err != nil {
		return errors.Wrap(err, "cannot create image stream")
	}

	resource, err := ioutil.ReadFile(ctx.Archive)
	if err != nil {
		return errors.Wrap(err, "cannot fully read tar file "+ctx.Archive)
	}

	restClient, err := customclient.GetClientFor("build.openshift.io", "v1")
	if err != nil {
		return err
	}

	result := restClient.Post().
		Namespace(ctx.Namespace).
		Body(resource).
		Resource("buildconfigs").
		Name("camel-k-" + ctx.Request.Meta.Name).
		SubResource("instantiatebinary").
		Do()

	if result.Error() != nil {
		return errors.Wrap(result.Error(), "cannot instantiate binary")
	}

	data, err := result.Raw()
	if err != nil {
		return errors.Wrap(err, "no raw data retrieved")
	}

	u := unstructured.Unstructured{}
	err = u.UnmarshalJSON(data)
	if err != nil {
		return errors.Wrap(err, "cannot unmarshal instantiate binary response")
	}

	ocbuild, err := k8sutil.RuntimeObjectFromUnstructured(&u)
	if err != nil {
		return err
	}

	err = kubernetes.WaitCondition(ocbuild, func(obj interface{}) (bool, error) {
		if val, ok := obj.(*buildv1.Build); ok {
			if val.Status.Phase == buildv1.BuildPhaseComplete {
				return true, nil
			} else if val.Status.Phase == buildv1.BuildPhaseCancelled ||
				val.Status.Phase == buildv1.BuildPhaseFailed ||
				val.Status.Phase == buildv1.BuildPhaseError {
				return false, errors.New("build failed")
			}
		}
		return false, nil
	}, 5*time.Minute)

	err = sdk.Get(&is)
	if err != nil {
		return err
	}

	if is.Status.DockerImageRepository == "" {
		return errors.New("dockerImageRepository not available in ImageStream")
	}

	ctx.Image = is.Status.DockerImageRepository + ":" + ctx.Request.Meta.ResourceVersion

	return nil
}
