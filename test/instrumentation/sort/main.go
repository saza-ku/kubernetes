/*
Copyright 2023 The Kubernetes Authors.

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

package main

import (
	"fmt"
	"os"
	"sort"

	flag "github.com/spf13/pflag"
	yaml "go.yaml.in/yaml/v2"
	"k8s.io/component-base/metrics"
)

func main() {
	var sortFile string
	flag.StringVar(&sortFile, "sort-file", "", "file of metrics to sort")
	flag.Parse()
	dat, err := os.ReadFile(sortFile)
	if err == nil {
		var parsedMetrics []metric
		err = yaml.Unmarshal(dat, &parsedMetrics)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		sort.Sort(byFQName(parsedMetrics))
		data, err := yaml.Marshal(parsedMetrics)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		fmt.Print(string(data))
	}
}

type metric struct {
	Name              string              `yaml:"name" json:"name"`
	Subsystem         string              `yaml:"subsystem,omitempty" json:"subsystem,omitempty"`
	Namespace         string              `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	Help              string              `yaml:"help,omitempty" json:"help,omitempty"`
	Type              string              `yaml:"type,omitempty" json:"type,omitempty"`
	DeprecatedVersion string              `yaml:"deprecatedVersion,omitempty" json:"deprecatedVersion,omitempty"`
	StabilityLevel    string              `yaml:"stabilityLevel,omitempty" json:"stabilityLevel,omitempty"`
	Labels            []string            `yaml:"labels,omitempty" json:"labels,omitempty"`
	Buckets           []float64           `yaml:"buckets,omitempty" json:"buckets,omitempty"`
	Objectives        map[float64]float64 `yaml:"objectives,omitempty" json:"objectives,omitempty"`
	AgeBuckets        uint32              `yaml:"ageBuckets,omitempty" json:"ageBuckets,omitempty"`
	BufCap            uint32              `yaml:"bufCap,omitempty" json:"bufCap,omitempty"`
	MaxAge            int64               `yaml:"maxAge,omitempty" json:"maxAge,omitempty"`
	ConstLabels       map[string]string   `yaml:"constLabels,omitempty" json:"constLabels,omitempty"`
}

func (m metric) BuildFQName() string {
	return metrics.BuildFQName(m.Namespace, m.Subsystem, m.Name)
}

type byFQName []metric

func (ms byFQName) Len() int { return len(ms) }
func (ms byFQName) Less(i, j int) bool {
	if ms[i].StabilityLevel < ms[j].StabilityLevel {
		return true
	} else if ms[i].StabilityLevel > ms[j].StabilityLevel {
		return false
	}
	return ms[i].BuildFQName() < ms[j].BuildFQName()
}
func (ms byFQName) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}
