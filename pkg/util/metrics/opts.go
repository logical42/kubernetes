package metrics

import (
	"github.com/blang/semver"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

type Opts struct {
	Namespace         string
	Subsystem         string
	Name              string
	Help              string
	ConstLabels       prometheus.Labels
	DeprecatedVersion *semver.Version
	deprecateOnce     sync.Once
}
