/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package plugin

import (
	"errors"
	"fmt"
	"testing"
)

func mockPluginFilter(*Registration) bool {
	return false
}

// Plugin types commonly used by containerd
const (
	InternalPlugin         Type = "io.containerd.internal.v1"
	RuntimePlugin          Type = "io.containerd.runtime.v1"
	RuntimePluginV2        Type = "io.containerd.runtime.v2"
	ServicePlugin          Type = "io.containerd.service.v1"
	GRPCPlugin             Type = "io.containerd.grpc.v1"
	SnapshotPlugin         Type = "io.containerd.snapshotter.v1"
	TaskMonitorPlugin      Type = "io.containerd.monitor.v1"
	DiffPlugin             Type = "io.containerd.differ.v1"
	MetadataPlugin         Type = "io.containerd.metadata.v1"
	ContentPlugin          Type = "io.containerd.content.v1"
	GCPlugin               Type = "io.containerd.gc.v1"
	LeasePlugin            Type = "io.containerd.lease.v1"
	TracingProcessorPlugin Type = "io.containerd.tracing.processor.v1"
)

func testRegistry() Registry {
	var register Registry
	return register.Register(&Registration{
		Type: TaskMonitorPlugin,
		ID:   "cgroups",
	}).Register(&Registration{
		Type: ServicePlugin,
		ID:   "tasks-service",
		Requires: []Type{
			RuntimePlugin,
			RuntimePluginV2,
			MetadataPlugin,
			TaskMonitorPlugin,
		},
	}).Register(&Registration{
		Type: ServicePlugin,
		ID:   "introspection-service",
	}).Register(&Registration{
		Type: ServicePlugin,
		ID:   "namespaces-service",
		Requires: []Type{
			MetadataPlugin,
		},
	}).Register(&Registration{
		Type: GRPCPlugin,
		ID:   "namespaces",
		Requires: []Type{
			ServicePlugin,
		},
	}).Register(&Registration{
		Type: GRPCPlugin,
		ID:   "content",
		Requires: []Type{
			ServicePlugin,
		},
	}).Register(&Registration{
		Type: GRPCPlugin,
		ID:   "containers",
		Requires: []Type{
			ServicePlugin,
		},
	}).Register(&Registration{
		Type: ServicePlugin,
		ID:   "containers-service",
		Requires: []Type{
			MetadataPlugin,
		},
	}).Register(&Registration{
		Type: GRPCPlugin,
		ID:   "events",
	}).Register(&Registration{
		Type: GRPCPlugin,
		ID:   "leases",
		Requires: []Type{
			LeasePlugin,
		},
	}).Register(&Registration{
		Type: LeasePlugin,
		ID:   "manager",
		Requires: []Type{
			MetadataPlugin,
		},
	}).Register(&Registration{
		Type: GRPCPlugin,
		ID:   "diff",
		Requires: []Type{
			ServicePlugin,
		},
	}).Register(&Registration{
		Type: ServicePlugin,
		ID:   "diff-service",
		Requires: []Type{
			DiffPlugin,
		},
	}).Register(&Registration{
		Type: ServicePlugin,
		ID:   "snapshots-service",
		Requires: []Type{
			MetadataPlugin,
		},
	}).Register(&Registration{
		Type: GRPCPlugin,
		ID:   "snapshots",
		Requires: []Type{
			ServicePlugin,
		},
	}).Register(&Registration{
		Type: GRPCPlugin,
		ID:   "version",
	}).Register(&Registration{
		Type: GRPCPlugin,
		ID:   "images",
		Requires: []Type{
			ServicePlugin,
		},
	}).Register(&Registration{
		Type: GCPlugin,
		ID:   "scheduler",
		Requires: []Type{
			MetadataPlugin,
		},
	}).Register(&Registration{
		Type: RuntimePluginV2,
		ID:   "task",
		Requires: []Type{
			MetadataPlugin,
		},
	}).Register(&Registration{
		Type: GRPCPlugin,
		ID:   "tasks",
		Requires: []Type{
			ServicePlugin,
		},
	}).Register(&Registration{
		Type:     GRPCPlugin,
		ID:       "introspection",
		Requires: []Type{"*"},
	}).Register(&Registration{
		Type: ServicePlugin,
		ID:   "content-service",
		Requires: []Type{
			MetadataPlugin,
		},
	}).Register(&Registration{
		Type: GRPCPlugin,
		ID:   "healthcheck",
	}).Register(&Registration{
		Type: InternalPlugin,
		ID:   "opt",
	}).Register(&Registration{
		Type: GRPCPlugin,
		ID:   "cri",
		Requires: []Type{
			ServicePlugin,
		},
	}).Register(&Registration{
		Type: RuntimePlugin,
		ID:   "linux",
		Requires: []Type{
			MetadataPlugin,
		},
	}).Register(&Registration{
		Type: InternalPlugin,
		Requires: []Type{
			ServicePlugin,
		},
		ID: "restart",
	}).Register(&Registration{
		Type: DiffPlugin,
		ID:   "walking",
		Requires: []Type{
			MetadataPlugin,
		},
	}).Register(&Registration{
		Type: SnapshotPlugin,
		ID:   "native",
	}).Register(&Registration{
		Type: SnapshotPlugin,
		ID:   "overlayfs",
	}).Register(&Registration{
		Type: ContentPlugin,
		ID:   "content",
	}).Register(&Registration{
		Type: MetadataPlugin,
		ID:   "bolt",
		Requires: []Type{
			ContentPlugin,
			SnapshotPlugin,
		},
	}).Register(&Registration{
		Type: TracingProcessorPlugin,
		ID:   "otlp",
	}).Register(&Registration{
		Type: InternalPlugin,
		ID:   "tracing",
		Requires: []Type{
			TracingProcessorPlugin,
		},
	})
}

// TestContainerdPlugin tests the logic of Graph, use the containerd's plugin
func TestContainerdPlugin(t *testing.T) {
	register := testRegistry()
	ordered := register.Graph(mockPluginFilter)
	expectedURI := []string{
		"io.containerd.monitor.v1.cgroups",
		"io.containerd.content.v1.content",
		"io.containerd.snapshotter.v1.native",
		"io.containerd.snapshotter.v1.overlayfs",
		"io.containerd.metadata.v1.bolt",
		"io.containerd.runtime.v1.linux",
		"io.containerd.runtime.v2.task",
		"io.containerd.service.v1.tasks-service",
		"io.containerd.service.v1.introspection-service",
		"io.containerd.service.v1.namespaces-service",
		"io.containerd.service.v1.containers-service",
		"io.containerd.differ.v1.walking",
		"io.containerd.service.v1.diff-service",
		"io.containerd.service.v1.snapshots-service",
		"io.containerd.service.v1.content-service",
		"io.containerd.grpc.v1.namespaces",
		"io.containerd.grpc.v1.content",
		"io.containerd.grpc.v1.containers",
		"io.containerd.grpc.v1.events",
		"io.containerd.lease.v1.manager",
		"io.containerd.grpc.v1.leases",
		"io.containerd.grpc.v1.diff",
		"io.containerd.grpc.v1.snapshots",
		"io.containerd.grpc.v1.version",
		"io.containerd.grpc.v1.images",
		"io.containerd.gc.v1.scheduler",
		"io.containerd.grpc.v1.tasks",
		"io.containerd.grpc.v1.healthcheck",
		"io.containerd.internal.v1.opt",
		"io.containerd.grpc.v1.cri",
		"io.containerd.internal.v1.restart",
		"io.containerd.tracing.processor.v1.otlp",
		"io.containerd.internal.v1.tracing",
		"io.containerd.grpc.v1.introspection",
	}
	cmpOrdered(t, ordered, expectedURI)
}

func cmpOrdered(t *testing.T, ordered []Registration, expectedURI []string) {
	if len(ordered) != len(expectedURI) {
		t.Fatalf("ordered compare failed, %d != %d", len(ordered), len(expectedURI))
	}
	for i := range ordered {
		if ordered[i].URI() != expectedURI[i] {
			t.Fatalf("graph failed, expected: %s, but return: %s", expectedURI[i], ordered[i].URI())
		}
	}
}

// TestPluginGraph tests the logic of Graph
func TestPluginGraph(t *testing.T) {
	for _, testcase := range []struct {
		input       []*Registration
		expectedURI []string
		filter      DisableFilter
	}{
		// test requires *
		{
			input: []*Registration{
				{
					Type: "grpc",
					ID:   "introspection",
					Requires: []Type{
						"*",
					},
				},
				{
					Type: "service",
					ID:   "container",
				},
			},
			expectedURI: []string{
				"service.container",
				"grpc.introspection",
			},
		},
		// test requires
		{
			input: []*Registration{
				{
					Type: "service",
					ID:   "container",
					Requires: []Type{
						"metadata",
					},
				},
				{
					Type: "metadata",
					ID:   "bolt",
				},
			},
			expectedURI: []string{
				"metadata.bolt",
				"service.container",
			},
		},
		{
			input: []*Registration{
				{
					Type: "metadata",
					ID:   "bolt",
					Requires: []Type{
						"content",
						"snapshotter",
					},
				},
				{
					Type: "snapshotter",
					ID:   "overlayfs",
				},
				{
					Type: "content",
					ID:   "content",
				},
			},
			expectedURI: []string{
				"content.content",
				"snapshotter.overlayfs",
				"metadata.bolt",
			},
		},
		// test disable
		{
			input: []*Registration{
				{
					Type: "content",
					ID:   "content",
				},
				{
					Type: "disable",
					ID:   "disable",
				},
			},
			expectedURI: []string{
				"content.content",
			},
			filter: func(r *Registration) bool {
				return r.Type == "disable"
			},
		},
	} {
		var register Registry
		for _, in := range testcase.input {
			register = register.Register(in)
		}
		var filter DisableFilter = mockPluginFilter
		if testcase.filter != nil {
			filter = testcase.filter
		}
		ordered := register.Graph(filter)
		cmpOrdered(t, ordered, testcase.expectedURI)
	}
}

func TestGetPlugins(t *testing.T) {
	otherError := fmt.Errorf("other error")
	plugins := NewPluginSet()
	for _, p := range []*Plugin{
		testPlugin("type1", "id1", "id1", nil),
		testPlugin("type1", "id2", "id2", ErrSkipPlugin),
		testPlugin("type2", "id3", "id3", ErrSkipPlugin),
		testPlugin("type3", "id4", "id4", nil),
		testPlugin("type4", "id5", "id5", nil),
		testPlugin("type4", "id6", "id6", nil),
		testPlugin("type5", "id7", "id7", otherError),
	} {
		plugins.Add(p)
	}

	ic := InitContext{
		plugins: plugins,
	}

	for _, tc := range []struct {
		pluginType string
		err        error
	}{
		{"type1", nil},
		{"type2", ErrPluginNotFound},
		{"type3", nil},
		{"type4", ErrPluginMultipleInstances},
		{"type5", otherError},
	} {
		t.Run("GetSingle", func(t *testing.T) {
			instance, err := ic.GetSingle(Type(tc.pluginType))
			if err != nil {
				if tc.err == nil {
					t.Fatalf("unexpected error %v", err)
				} else if !errors.Is(err, tc.err) {
					t.Fatalf("unexpected error %v, expected %v", err, tc.err)
				}
				return
			} else if tc.err != nil {
				t.Fatalf("expected error %v, got no error", tc.err)
			}
			_, ok := instance.(string)
			if !ok {
				t.Fatalf("unexpected instance value %v", instance)
			}
		})
	}

	for _, tc := range []struct {
		pluginType string
		expected   []string
		err        error
	}{
		{"type1", []string{"id1"}, nil},
		{"type2", nil, ErrPluginNotFound},
		{"type3", []string{"id4"}, nil},
		{"type4", []string{"id5", "id6"}, nil},
		{"type5", nil, otherError},
	} {
		t.Run("GetByType", func(t *testing.T) {
			m, err := ic.GetByType(Type(tc.pluginType))
			if err != nil {
				if tc.err == nil {
					t.Fatalf("unexpected error %v", err)
				} else if !errors.Is(err, tc.err) {
					t.Fatalf("unexpected error %v, expected %v", err, tc.err)
				}
				return
			} else if tc.err != nil {
				t.Fatalf("expected error %v, got no error", tc.err)
			}

			if len(m) != len(tc.expected) {
				t.Fatalf("unexpected result %v, expected %v", m, tc.expected)
			}
			for _, v := range tc.expected {
				instance, ok := m[v]
				if !ok {
					t.Errorf("missing value for %q", v)
					continue
				}
				if instance.(string) != v {
					t.Errorf("unexpected value %v, expected %v", instance, v)
				}
			}
		})
	}

	for _, tc := range []struct {
		pluginType string
		id         string
		err        error
	}{
		{"type1", "id1", nil},
		{"type1", "id2", ErrSkipPlugin},
		{"type2", "id3", ErrSkipPlugin},
		{"type3", "id4", nil},
		{"type4", "id5", nil},
		{"type4", "id6", nil},
		{"type5", "id7", otherError},
	} {
		t.Run("GetByID", func(t *testing.T) {
			instance, err := ic.GetByID(Type(tc.pluginType), tc.id)
			if err != nil {
				if tc.err == nil {
					t.Fatalf("unexpected error %v", err)
				} else if !errors.Is(err, tc.err) {
					t.Fatalf("unexpected error %v, expected %v", err, tc.err)
				}
				return
			} else if tc.err != nil {
				t.Fatalf("expected error %v, got no error", tc.err)
			}

			if instance.(string) != tc.id {
				t.Errorf("unexpected value %v, expected %v", instance, tc.id)
			}
		})
	}

}

func testPlugin(t Type, id string, i interface{}, err error) *Plugin {
	return &Plugin{
		Registration: Registration{
			Type: t,
			ID:   id,
		},
		instance: i,
		err:      err,
	}
}

func TestRequiresAll(t *testing.T) {
	var register Registry
	register = register.Register(&Registration{
		Type: InternalPlugin,
		ID:   "system",
	}).Register(&Registration{
		Type: ServicePlugin,
		ID:   "introspection",
		Requires: []Type{
			"*",
		},
	}).Register(&Registration{
		Type: ServicePlugin,
		ID:   "task",
		Requires: []Type{
			InternalPlugin,
		},
	}).Register(&Registration{
		Type: ServicePlugin,
		ID:   "version",
	})
	ordered := register.Graph(mockPluginFilter)
	expectedURI := []string{
		"io.containerd.internal.v1.system",
		"io.containerd.service.v1.task",
		"io.containerd.service.v1.version",
		"io.containerd.service.v1.introspection",
	}
	cmpOrdered(t, ordered, expectedURI)
}

func TestRegisterErrors(t *testing.T) {

	for _, tc := range []struct {
		name     string
		expected error
		register func(Registry) Registry
	}{
		{
			name:     "duplicate",
			expected: ErrIDRegistered,
			register: func(r Registry) Registry {
				return r.Register(&Registration{
					Type: TaskMonitorPlugin,
					ID:   "cgroups",
				}).Register(&Registration{
					Type: TaskMonitorPlugin,
					ID:   "cgroups",
				})
			},
		},
		{
			name:     "circular",
			expected: ErrPluginCircularDependency,
			register: func(r Registry) Registry {
				// Circular dependencies should not loop but order will be based on registration order
				return r.Register(&Registration{
					Type: InternalPlugin,
					ID:   "p1",
					Requires: []Type{
						RuntimePlugin,
					},
				}).Register(&Registration{
					Type: RuntimePlugin,
					ID:   "p2",
					Requires: []Type{
						InternalPlugin,
					},
				}).Register(&Registration{
					Type: InternalPlugin,
					ID:   "p3",
				})
			},
		},
		{
			name:     "self",
			expected: ErrInvalidRequires,
			register: func(r Registry) Registry {
				// Circular dependencies should not loop but order will be based on registration order
				return r.Register(&Registration{
					Type: InternalPlugin,
					ID:   "p1",
					Requires: []Type{
						InternalPlugin,
					},
				})
			},
		},
		{
			name:     "no-type",
			expected: ErrNoType,
			register: func(r Registry) Registry {
				// Circular dependencies should not loop but order will be based on registration order
				return r.Register(&Registration{
					Type: "",
					ID:   "p1",
				})
			},
		},
		{
			name:     "no-ID",
			expected: ErrNoPluginID,
			register: func(r Registry) Registry {
				// Circular dependencies should not loop but order will be based on registration order
				return r.Register(&Registration{
					Type: InternalPlugin,
					ID:   "",
				})
			},
		},
		{
			name:     "bad-requires-all",
			expected: ErrInvalidRequires,
			register: func(r Registry) Registry {
				// Circular dependencies should not loop but order will be based on registration order
				return r.Register(&Registration{
					Type: InternalPlugin,
					ID:   "p1",
					Requires: []Type{
						"*",
						InternalPlugin,
					},
				})
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var (
				r        Registry
				panicAny any
			)
			func() {
				defer func() {
					panicAny = recover()
				}()

				tc.register(r).Graph(mockPluginFilter)
			}()
			if panicAny == nil {
				t.Fatalf("expected panic with error %v", tc.expected)
			}
			err, ok := panicAny.(error)
			if !ok {
				t.Fatalf("expected panic: %v, expected error %v", panicAny, tc.expected)
			}
			if !errors.Is(err, tc.expected) {
				t.Fatalf("unexpected error type: %v, expected %v", panicAny, tc.expected)
			}
		})
	}
}

func BenchmarkGraph(b *testing.B) {
	register := testRegistry()
	b.ResetTimer()
	for range b.N {
		register.Graph(mockPluginFilter)
	}
}

func BenchmarkUnique(b *testing.B) {
	register := testRegistry()
	b.ResetTimer()
	for range b.N {
		checkUnique(register, &Registration{
			Type: InternalPlugin,
			ID:   "new",
		})
	}
}
