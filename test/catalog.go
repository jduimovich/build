// Copyright The Shipwright Contributors
//
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"context"
	"strconv"
	"time"

	. "github.com/onsi/gomega"
	build "github.com/shipwright-io/build/pkg/apis/build/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/apis"
	knativev1beta1 "knative.dev/pkg/apis/duck/v1beta1"
	crc "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

// Catalog allows you to access helper functions
type Catalog struct{}

// BuildWithClusterBuildStrategy gives you an specific Build CRD
func (c *Catalog) BuildWithClusterBuildStrategy(name string, ns string, strategyName string, secretName string) *build.Build {
	buildStrategy := build.ClusterBuildStrategyKind
	return &build.Build{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: build.BuildSpec{
			Source: build.GitSource{
				URL: "foobar",
			},
			StrategyRef: &build.StrategyRef{
				Name: strategyName,
				Kind: &buildStrategy,
			},
			Output: build.Image{
				ImageURL: "foobar",
				SecretRef: &corev1.LocalObjectReference{
					Name: secretName,
				},
			},
		},
	}
}

// BuildWithBuildStrategy gives you an specific Build CRD
func (c *Catalog) BuildWithBuildStrategy(name string, ns string, strategyName string) *build.Build {
	buildStrategy := build.NamespacedBuildStrategyKind
	return &build.Build{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: build.BuildSpec{
			Source: build.GitSource{
				URL: "foobar",
			},
			StrategyRef: &build.StrategyRef{
				Name: strategyName,
				Kind: &buildStrategy,
			},
		},
	}
}

// BuildWithNilBuildStrategyKind gives you an Build CRD with nil build strategy kind
func (c *Catalog) BuildWithNilBuildStrategyKind(name string, ns string, strategyName string) *build.Build {
	return &build.Build{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: build.BuildSpec{
			Source: build.GitSource{
				URL: "foobar",
			},
			StrategyRef: &build.StrategyRef{
				Name: strategyName,
			},
		},
	}
}

// BuildWithCustomAnnotationAndFinalizer provides a Build CRD with a customize annotation
// and optional finalizer
func (c *Catalog) BuildWithCustomAnnotationAndFinalizer(
	name string,
	ns string,
	strategyName string,
	a map[string]string,
	f []string,
) *build.Build {
	buildStrategy := build.ClusterBuildStrategyKind
	return &build.Build{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   ns,
			Annotations: a,
			Finalizers:  f,
		},
		Spec: build.BuildSpec{
			Source: build.GitSource{
				URL: "foobar",
			},
			StrategyRef: &build.StrategyRef{
				Name: strategyName,
				Kind: &buildStrategy,
			},
		},
	}
}

// ClusterBuildStrategyList to support tests
func (c *Catalog) ClusterBuildStrategyList(name string) *build.ClusterBuildStrategyList {
	return &build.ClusterBuildStrategyList{
		Items: []build.ClusterBuildStrategy{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: "build-examples",
				},
			},
		},
	}
}

// FakeClusterBuildStrategyList to support tests
func (c *Catalog) FakeClusterBuildStrategyList() *build.ClusterBuildStrategyList {
	return &build.ClusterBuildStrategyList{
		Items: []build.ClusterBuildStrategy{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foobar",
					Namespace: "build-examples",
				},
			},
		},
	}
}

// FakeNoClusterBuildStrategyList to support tests
func (c *Catalog) FakeNoClusterBuildStrategyList() *build.ClusterBuildStrategyList {
	return &build.ClusterBuildStrategyList{
		Items: []build.ClusterBuildStrategy{},
	}
}

// BuildStrategyList to support tests
func (c *Catalog) BuildStrategyList(name string, ns string) *build.BuildStrategyList {
	return &build.BuildStrategyList{
		Items: []build.BuildStrategy{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: ns,
				},
			},
		},
	}
}

// FakeBuildStrategyList to support tests
func (c *Catalog) FakeBuildStrategyList() *build.BuildStrategyList {
	return &build.BuildStrategyList{
		Items: []build.BuildStrategy{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foobar",
				},
			},
		},
	}
}

// FakeNoBuildStrategyList to support tests
func (c *Catalog) FakeNoBuildStrategyList() *build.BuildStrategyList {
	return &build.BuildStrategyList{
		Items: []build.BuildStrategy{},
	}
}

// FakeSecretList to support tests
func (c *Catalog) FakeSecretList() corev1.SecretList {
	return corev1.SecretList{
		Items: []corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foobar",
				},
			},
		},
	}
}

// FakeSecretListInNamespace to support test
func (c *Catalog) FakeNoSecretListInNamespace() corev1.SecretList {
	return corev1.SecretList{
		Items: []corev1.Secret{},
	}
}

// SecretList to support tests
func (c *Catalog) SecretList(name string) corev1.SecretList {
	return corev1.SecretList{
		Items: []corev1.Secret{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
				},
			},
		},
	}
}

// StubFunc is used to simulate the status of the Build
// after a .Status().Update() call in the controller. This
// receives a parameter to return an specific status state
func (c *Catalog) StubFunc(status corev1.ConditionStatus, reason string) func(context context.Context, object runtime.Object, _ ...crc.UpdateOption) error {
	return func(context context.Context, object runtime.Object, _ ...crc.UpdateOption) error {
		switch object := object.(type) {
		case *build.Build:
			Expect(object.Status.Registered).To(Equal(status))
			Expect(object.Status.Reason).To(ContainSubstring(reason))
		}
		return nil
	}
}

// StubBuildUpdateWithFinalizers is used to simulate the state of the Build
// finalizers after a client Update on the Object happened.
func (c *Catalog) StubBuildUpdateWithFinalizers(finalizer string) func(context context.Context, object runtime.Object, _ ...crc.UpdateOption) error {
	return func(context context.Context, object runtime.Object, _ ...crc.UpdateOption) error {
		switch object := object.(type) {
		case *build.Build:
			Expect(object.Finalizers).To(ContainElement(finalizer))
		}
		return nil
	}
}

// StubBuildUpdateWithoutFinalizers is used to simulate the state of the Build
// finalizers after a client Update on the Object happened.
func (c *Catalog) StubBuildUpdateWithoutFinalizers() func(context context.Context, object runtime.Object, _ ...crc.UpdateOption) error {
	return func(context context.Context, object runtime.Object, _ ...crc.UpdateOption) error {
		switch object := object.(type) {
		case *build.Build:
			Expect(len(object.Finalizers)).To(BeZero())
		}
		return nil
	}
}

// StubBuildRunAndTaskRun is used to simulate the existence of a BuildRun
// only when there is a client GET on this object type
func (c *Catalog) StubBuildRun(
	b *build.BuildRun,
) func(context context.Context, nn types.NamespacedName, object runtime.Object) error {
	return func(context context.Context, nn types.NamespacedName, object runtime.Object) error {
		switch object := object.(type) {
		case *build.BuildRun:
			b.DeepCopyInto(object)
			return nil
		}
		return errors.NewNotFound(schema.GroupResource{}, nn.Name)
	}
}

// StubBuildRunAndTaskRun is used to simulate the existence of a BuildRun
// and a TaskRun when there is a client GET on this two objects
func (c *Catalog) StubBuildRunAndTaskRun(
	b *build.BuildRun,
	tr *v1beta1.TaskRun,
) func(context context.Context, nn types.NamespacedName, object runtime.Object) error {
	return func(context context.Context, nn types.NamespacedName, object runtime.Object) error {
		switch object := object.(type) {
		case *build.BuildRun:
			b.DeepCopyInto(object)
			return nil
		case *v1beta1.TaskRun:
			tr.DeepCopyInto(object)
			return nil
		}
		return errors.NewNotFound(schema.GroupResource{}, nn.Name)
	}
}

// StubBuildAndTaskRun is used to simulate the existence of a Build
// and a TaskRun when there is a client GET on this two objects
func (c *Catalog) StubBuildAndTaskRun(
	b *build.Build,
	tr *v1beta1.TaskRun,
) func(context context.Context, nn types.NamespacedName, object runtime.Object) error {
	return func(context context.Context, nn types.NamespacedName, object runtime.Object) error {
		switch object := object.(type) {
		case *build.Build:
			b.DeepCopyInto(object)
			return nil
		case *v1beta1.TaskRun:
			tr.DeepCopyInto(object)
			return nil
		}
		return errors.NewNotFound(schema.GroupResource{}, nn.Name)
	}
}

// StubBuildRunStatus asserts Status fields on a BuildRun resource
func (c *Catalog) StubBuildRunStatus(reason string, name *string, status corev1.ConditionStatus, buildSpec build.BuildSpec, tolerateEmptyStatus bool) func(context context.Context, object runtime.Object, _ ...crc.UpdateOption) error {
	return func(context context.Context, object runtime.Object, _ ...crc.UpdateOption) error {
		switch object := object.(type) {
		case *build.BuildRun:
			if !tolerateEmptyStatus || object.Status.Succeeded != "" {
				Expect(object.Status.Succeeded).To(Equal(status))
				Expect(object.Status.Reason).To(Equal(reason))
				Expect(object.Status.LatestTaskRunRef).To(Equal(name))
			}
			if object.Status.BuildSpec != nil {
				Expect(*object.Status.BuildSpec).To(Equal(buildSpec))
			}
		}
		return nil
	}
}

// StubBuildRunLabel asserts Label fields on a BuildRun resource
func (c *Catalog) StubBuildRunLabel(buildSample *build.Build) func(context context.Context, object runtime.Object, _ ...crc.UpdateOption) error {
	return func(context context.Context, object runtime.Object, _ ...crc.UpdateOption) error {
		switch object := object.(type) {
		case *build.BuildRun:
			Expect(object.Labels[build.LabelBuild]).To(Equal(buildSample.Name))
			Expect(object.Labels[build.LabelBuildGeneration]).To(Equal(strconv.FormatInt(buildSample.Generation, 10)))
		}
		return nil
	}
}

// StubBuildRunGetWithoutSA simulates the output of client GET calls
// for the BuildRun unit tests
func (c *Catalog) StubBuildRunGetWithoutSA(
	b *build.Build,
	br *build.BuildRun,
) func(context context.Context, nn types.NamespacedName, object runtime.Object) error {
	return func(context context.Context, nn types.NamespacedName, object runtime.Object) error {
		switch object := object.(type) {
		case *build.Build:
			b.DeepCopyInto(object)
			return nil
		case *build.BuildRun:
			br.DeepCopyInto(object)
			return nil
		}
		return errors.NewNotFound(schema.GroupResource{}, nn.Name)
	}
}

// StubBuildRunGetWithTaskRunAndSA returns fake object for different
// client calls
func (c *Catalog) StubBuildRunGetWithTaskRunAndSA(
	b *build.Build,
	br *build.BuildRun,
	tr *v1beta1.TaskRun,
	sa *corev1.ServiceAccount,
) func(context context.Context, nn types.NamespacedName, object runtime.Object) error {
	return func(context context.Context, nn types.NamespacedName, object runtime.Object) error {
		switch object := object.(type) {
		case *build.Build:
			b.DeepCopyInto(object)
			return nil
		case *build.BuildRun:
			br.DeepCopyInto(object)
			return nil
		case *v1beta1.TaskRun:
			tr.DeepCopyInto(object)
			return nil
		case *corev1.ServiceAccount:
			sa.DeepCopyInto(object)
			return nil
		}
		return errors.NewNotFound(schema.GroupResource{}, nn.Name)
	}
}

// StubBuildRunGetWithSA returns fake object for different
// client calls
func (c *Catalog) StubBuildRunGetWithSA(
	b *build.Build,
	br *build.BuildRun,
	sa *corev1.ServiceAccount,
) func(context context.Context, nn types.NamespacedName, object runtime.Object) error {
	return func(context context.Context, nn types.NamespacedName, object runtime.Object) error {
		switch object := object.(type) {
		case *build.Build:
			b.DeepCopyInto(object)
			return nil
		case *build.BuildRun:
			br.DeepCopyInto(object)
			return nil
		case *corev1.ServiceAccount:
			sa.DeepCopyInto(object)
			return nil
		}
		return errors.NewNotFound(schema.GroupResource{}, nn.Name)
	}
}

// StubBuildRunGetWithSAandStrategies simulates the output of client GET
// calls for the BuildRun unit tests
func (c *Catalog) StubBuildRunGetWithSAandStrategies(
	b *build.Build,
	br *build.BuildRun,
	sa *corev1.ServiceAccount,
	cb *build.ClusterBuildStrategy,
	bs *build.BuildStrategy,
) func(context context.Context, nn types.NamespacedName, object runtime.Object) error {
	return func(context context.Context, nn types.NamespacedName, object runtime.Object) error {
		switch object := object.(type) {
		case *build.Build:
			b.DeepCopyInto(object)
			return nil
		case *build.BuildRun:
			br.DeepCopyInto(object)
			return nil
		case *corev1.ServiceAccount:
			sa.DeepCopyInto(object)
			return nil
		case *build.ClusterBuildStrategy:
			cb.DeepCopyInto(object)
			return nil
		case *build.BuildStrategy:
			bs.DeepCopyInto(object)
			return nil
		}
		return errors.NewNotFound(schema.GroupResource{}, nn.Name)
	}
}

// StubBuildCRDsPodAndTaskRun stubs different objects in case a client
// GET call is executed against them
func (c *Catalog) StubBuildCRDsPodAndTaskRun(
	b *build.Build,
	br *build.BuildRun,
	sa *corev1.ServiceAccount,
	cb *build.ClusterBuildStrategy,
	bs *build.BuildStrategy,
	tr *v1beta1.TaskRun,
	pod *corev1.Pod,
) func(context context.Context, nn types.NamespacedName, object runtime.Object) error {
	return func(context context.Context, nn types.NamespacedName, object runtime.Object) error {
		switch object := object.(type) {
		case *build.Build:
			b.DeepCopyInto(object)
			return nil
		case *build.BuildRun:
			br.DeepCopyInto(object)
			return nil
		case *corev1.ServiceAccount:
			sa.DeepCopyInto(object)
			return nil
		case *build.ClusterBuildStrategy:
			cb.DeepCopyInto(object)
			return nil
		case *build.BuildStrategy:
			bs.DeepCopyInto(object)
			return nil
		case *v1beta1.TaskRun:
			tr.DeepCopyInto(object)
			return nil
		case *corev1.Pod:
			pod.DeepCopyInto(object)
			return nil
		}
		return errors.NewNotFound(schema.GroupResource{}, nn.Name)
	}
}

// DefaultTaskRunWithStatus returns a minimal tekton TaskRun with an Status
func (c *Catalog) DefaultTaskRunWithStatus(trName string, buildRunName string, ns string, status corev1.ConditionStatus, reason string) *v1beta1.TaskRun {
	return &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      trName,
			Namespace: ns,
			Labels:    map[string]string{"buildrun.build.dev/name": buildRunName},
		},
		Spec: v1beta1.TaskRunSpec{},
		Status: v1beta1.TaskRunStatus{
			Status: knativev1beta1.Status{
				Conditions: knativev1beta1.Conditions{
					{
						Type:   apis.ConditionSucceeded,
						Reason: reason,
						Status: status,
					},
				},
			},
		},
	}
}

// TaskRunWithCompletionAndStartTime provides a TaskRun object with a
// Completion and StartTime
func (c *Catalog) TaskRunWithCompletionAndStartTime(trName string, buildRunName string, ns string) *v1beta1.TaskRun {
	return &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      trName,
			Namespace: ns,
			Labels:    map[string]string{"buildrun.build.dev/name": buildRunName},
		},
		Spec: v1beta1.TaskRunSpec{},
		Status: v1beta1.TaskRunStatus{
			TaskRunStatusFields: v1beta1.TaskRunStatusFields{
				CompletionTime: &metav1.Time{
					Time: time.Now(),
				},
				StartTime: &metav1.Time{
					Time: time.Now(),
				},
				PodName: "foobar",
			},
			Status: knativev1beta1.Status{
				Conditions: knativev1beta1.Conditions{
					{
						Type:    apis.ConditionSucceeded,
						Reason:  "something bad happened",
						Status:  corev1.ConditionFalse,
						Message: "some message",
					},
				},
			},
		},
	}
}

// DefaultTaskRunWithFalseStatus returns a minimal tektont TaskRun with a FALSE status
func (c *Catalog) DefaultTaskRunWithFalseStatus(trName string, buildRunName string, ns string) *v1beta1.TaskRun {
	return &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      trName,
			Namespace: ns,
			Labels:    map[string]string{"buildrun.build.dev/name": buildRunName},
		},
		Spec: v1beta1.TaskRunSpec{},
		Status: v1beta1.TaskRunStatus{
			Status: knativev1beta1.Status{
				Conditions: knativev1beta1.Conditions{
					{
						Type:    apis.ConditionSucceeded,
						Reason:  "something bad happened",
						Status:  corev1.ConditionFalse,
						Message: "some message",
					},
				},
			},
		},
	}
}

// DefaultBuild returns a minimal Build object
func (c *Catalog) DefaultBuild(buildName string, strategyName string, strategyKind build.BuildStrategyKind) *build.Build {
	return &build.Build{
		ObjectMeta: metav1.ObjectMeta{
			Name: buildName,
		},
		Spec: build.BuildSpec{
			StrategyRef: &build.StrategyRef{
				Name: strategyName,
				Kind: &strategyKind,
			},
		},
		Status: build.BuildStatus{
			Registered: corev1.ConditionTrue,
		},
	}
}

// DefaultBuildWithFalseRegistered returns a minimal Build object with a FALSE Registered
func (c *Catalog) DefaultBuildWithFalseRegistered(buildName string, strategyName string, strategyKind build.BuildStrategyKind) *build.Build {
	return &build.Build{
		ObjectMeta: metav1.ObjectMeta{
			Name: buildName,
		},
		Spec: build.BuildSpec{
			StrategyRef: &build.StrategyRef{
				Name: strategyName,
				Kind: &strategyKind,
			},
		},
		Status: build.BuildStatus{
			Registered: corev1.ConditionFalse,
			Reason:     "something bad happened",
		},
	}
}

// DefaultBuildRun returns a minimal BuildRun object
func (c *Catalog) DefaultBuildRun(buildRunName string, buildName string) *build.BuildRun {
	return &build.BuildRun{
		ObjectMeta: metav1.ObjectMeta{
			Name: buildRunName,
		},
		Spec: build.BuildRunSpec{
			BuildRef: &build.BuildRef{
				Name: buildName,
			},
		},
	}
}

// PodWithInitContainerStatus returns a pod with a single
// entry under the Status field for InitContainer Status
func (c *Catalog) PodWithInitContainerStatus(podName string, initContainerName string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
		},
		Status: corev1.PodStatus{
			InitContainerStatuses: []corev1.ContainerStatus{
				{
					Name: initContainerName,
				},
			},
		},
	}
}

// BuildRunWithBuildSnapshot returns BuildRun Object with a populated
// BuildSpec in the Status field
func (c *Catalog) BuildRunWithBuildSnapshot(buildRunName string, buildName string) *build.BuildRun {
	return &build.BuildRun{
		ObjectMeta: metav1.ObjectMeta{
			Name: buildRunName,
			CreationTimestamp: metav1.Time{
				Time: time.Now(),
			},
		},
		Status: build.BuildRunStatus{
			BuildSpec: &build.BuildSpec{
				StrategyRef: &build.StrategyRef{
					Name: "foobar",
				},
			},
		},
		Spec: build.BuildRunSpec{
			BuildRef: &build.BuildRef{
				Name: buildName,
			},
		},
	}
}

// DefaultTaskRun returns a minimal TaskRun object
func (c *Catalog) DefaultTaskRun(taskRunName string, ns string) *v1beta1.TaskRun {
	return &v1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      taskRunName,
			Namespace: ns,
		},
	}
}

// DefaultServiceAccount returns a minimal SA object
func (c *Catalog) DefaultServiceAccount(name string) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

// DefaultClusterBuildStrategy returns a minimal ClusterBuildStrategy
// object with a inmutable name
func (c *Catalog) DefaultClusterBuildStrategy() *build.ClusterBuildStrategy {
	return &build.ClusterBuildStrategy{
		ObjectMeta: metav1.ObjectMeta{
			Name: "foobar",
		},
	}
}

// DefaultNamespacedBuildStrategy returns a minimal BuildStrategy
// object with a inmutable name
func (c *Catalog) DefaultNamespacedBuildStrategy() *build.BuildStrategy {
	return &build.BuildStrategy{
		ObjectMeta: metav1.ObjectMeta{
			Name: "foobar",
		},
	}
}

// BuildRunWithSA returns a customized BuildRun object that defines a
// service account
func (c *Catalog) BuildRunWithSA(buildRunName string, buildName string, saName string) *build.BuildRun {
	return &build.BuildRun{
		ObjectMeta: metav1.ObjectMeta{
			Name: buildRunName,
		},
		Spec: build.BuildRunSpec{
			BuildRef: &build.BuildRef{
				Name: buildName,
			},
			ServiceAccount: &build.ServiceAccount{
				Name:     &saName,
				Generate: false,
			},
		},
	}
}

// BuildRunWithoutSA returns a buildrun without serviceAccountName and generate serviceAccount is false
func (c *Catalog) BuildRunWithoutSA(buildRunName string, buildName string) *build.BuildRun {
	return &build.BuildRun{
		ObjectMeta: metav1.ObjectMeta{
			Name: buildRunName,
		},
		Spec: build.BuildRunSpec{
			BuildRef: &build.BuildRef{
				Name: buildName,
			},
			ServiceAccount: &build.ServiceAccount{
				Generate: false,
			},
		},
	}
}

// BuildRunWithSAGenerate returns a customized BuildRun object that defines a
// service account
func (c *Catalog) BuildRunWithSAGenerate(buildRunName string, buildName string) *build.BuildRun {
	return &build.BuildRun{
		ObjectMeta: metav1.ObjectMeta{
			Name: buildRunName,
		},
		Spec: build.BuildRunSpec{
			BuildRef: &build.BuildRef{
				Name: buildName,
			},
			ServiceAccount: &build.ServiceAccount{
				Generate: true,
			},
		},
	}
}

// LoadCustomResources returns a container set of resources based on cpu and memory
func (c *Catalog) LoadCustomResources(cpu string, mem string) corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(cpu),
			corev1.ResourceMemory: resource.MustParse(mem),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(cpu),
			corev1.ResourceMemory: resource.MustParse(mem),
		},
	}
}

// LoadBuildYAML parses YAML bytes into JSON and from JSON
// into a Build struct
func (c *Catalog) LoadBuildYAML(d []byte) (*build.Build, error) {
	b := &build.Build{}
	err := yaml.Unmarshal(d, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// LoadBuildRunYAML parses YAML bytes into JSON and from JSON
// into a BuildRun struct
func (c *Catalog) LoadBuildRunYAML(d []byte) (*build.BuildRun, error) {
	b := &build.BuildRun{}
	err := yaml.Unmarshal(d, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// LoadBuildStrategyYAML parses YAML bytes into JSON and from JSON
// into a BuildStrategy struct
func (c *Catalog) LoadBuildStrategyYAML(d []byte) (*build.BuildStrategy, error) {
	b := &build.BuildStrategy{}
	err := yaml.Unmarshal(d, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
