/*
Copyright 2018 The Kubernetes Authors.

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

package app

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	goruntime "runtime"

	genericapifilters "k8s.io/apiserver/pkg/endpoints/filters"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	apiserver "k8s.io/apiserver/pkg/server"
	genericfilters "k8s.io/apiserver/pkg/server/filters"
	"k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/apiserver/pkg/server/mux"
	"k8s.io/apiserver/pkg/server/routes"
	componentbaseconfig "k8s.io/component-base/config"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/util/configz"

	clientmetrics "k8s.io/kubernetes/pkg/util/prometheusclientgo" // load all the prometheus client-go plugin
	versionmetrics "k8s.io/kubernetes/pkg/version/prometheus"      // for version metric registration
)

// BuildHandlerChain builds a handler chain with a base handler and CompletedConfig.
func BuildHandlerChain(apiHandler http.Handler, authorizationInfo *apiserver.AuthorizationInfo, authenticationInfo *apiserver.AuthenticationInfo) http.Handler {
	requestInfoResolver := &apirequest.RequestInfoFactory{}
	failedHandler := genericapifilters.Unauthorized(legacyscheme.Codecs, false)

	handler := apiHandler
	if authorizationInfo != nil {
		handler = genericapifilters.WithAuthorization(apiHandler, authorizationInfo.Authorizer, legacyscheme.Codecs)
	}
	if authenticationInfo != nil {
		handler = genericapifilters.WithAuthentication(handler, authenticationInfo.Authenticator, failedHandler, nil)
	}
	handler = genericapifilters.WithRequestInfo(handler, requestInfoResolver)
	handler = genericfilters.WithPanicRecovery(handler)

	return handler
}

// NewBaseHandler takes in CompletedConfig and returns a handler.
func NewBaseHandler(c *componentbaseconfig.DebuggingConfiguration, checks ...healthz.HealthzChecker) *mux.PathRecorderMux {
	mux := mux.NewPathRecorderMux("controller-manager")
	healthz.InstallHandler(mux, checks...)
	if c.EnableProfiling {
		routes.Profiling{}.Install(mux)
		if c.EnableContentionProfiling {
			goruntime.SetBlockProfileRate(1)
		}
	}
	configz.InstallHandler(mux)

	r := prometheus.NewRegistry()
	registerMetrics(r)
	mux.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}))
	return mux
}

func registerMetrics(registerer prometheus.Registerer) {
	versionmetrics.Register(registerer)
	clientmetrics.Register(registerer)
}