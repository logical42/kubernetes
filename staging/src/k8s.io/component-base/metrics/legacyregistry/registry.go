/*
Copyright 2019 The Kubernetes Authors.

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

package legacyregistry

import (
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/component-base/metrics"
)

var (
	defaultRegistry = metrics.NewKubeRegistry()
)

// Register registers a collectable metric but uses the global registry
func Register(c metrics.Registerable) error {
	return defaultRegistry.Register(c)
}

// MustRegister registers registerable metrics but uses the global registry.
func MustRegister(cs ...metrics.Registerable) {
	defaultRegistry.MustRegister(cs...)
}

// RawMustRegister registers prometheus collectors but uses the global registry, this
// bypasses the metric stability framework
//
// Deprecated
func RawMustRegister(cs ...prometheus.Collector) {
	defaultRegistry.RawMustRegister(cs...)
}

// RawRegister registers a prometheus collector but uses the global registry, this
// bypasses the metric stability framework
//
// Deprecated
func RawRegister(c prometheus.Collector) error {
	return defaultRegistry.RawRegister(c)
}
