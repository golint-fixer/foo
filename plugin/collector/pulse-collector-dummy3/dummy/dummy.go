package dummy

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/intelsdi-x/pulse/control/plugin"
	"github.com/intelsdi-x/pulse/control/plugin/cpolicy"
)

const (
	// Name of plugin
	Name = "dummy2"
	// Version of plugin
	Version = 2
	// Type of plugin
	Type = plugin.CollectorPluginType
)

// Dummy collector implementation used for testing
type Dummy struct {
}

//Random number generator
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// CollectMetrics collects metrics for testing
func (f *Dummy) CollectMetrics(mts []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {
	for _, p := range mts {
		log.Println("collecting", p)
	}
	rand.Seed(time.Now().UTC().UnixNano())
	for i, _ := range mts {
		data := randInt(65, 90)
		mts[i].Data_ = data
	}
	return mts, nil
}

//GetMetricTypes returns metric types for testing
func (f *Dummy) GetMetricTypes() ([]plugin.PluginMetricType, error) {
	m1 := &plugin.PluginMetricType{Namespace_: []string{"intel", "dummy", "foo"}}
	m2 := &plugin.PluginMetricType{Namespace_: []string{"intel", "dummy", "bar"}}
	return []plugin.PluginMetricType{*m1, *m2}, nil
}

//GetConfigPolicyTree returns a ConfigPolicyTree for testing
func (f *Dummy) GetConfigPolicyTree() (cpolicy.ConfigPolicyTree, error) {
	c := cpolicy.NewTree()
	rule, _ := cpolicy.NewStringRule("name", false, "bob")
	rule2, _ := cpolicy.NewStringRule("password", true)
	p := cpolicy.NewPolicyNode()
	p.Add(rule)
	p.Add(rule2)
	c.Add([]string{"intel", "dummy", "foo"}, p)
	return *c, nil
}

//Meta returns meta data for testing
func Meta() *plugin.PluginMeta {
	meta := plugin.NewPluginMeta(Name, Version, Type, []string{plugin.PulseGOBContentType}, []string{plugin.PulseGOBContentType})
	// Below we create a config policy for this plugin itself (different than the metrics above)
	// This policy ensures that enough configuration is provided on loading of this plugin to be able to run it.
	//
	// require a path to a specific binary
	rule1, _ := cpolicy.NewStringRule("some-binary-path", true)
	// optional log path that has a default
	rule2, _ := cpolicy.NewStringRule("optional-log-path", false, fmt.Sprintf("%s/pulse/dummy3", os.TempDir()))
	// optional port which defaults to 9999
	rule3, _ := cpolicy.NewIntegerRule("binary-port", false, 9999)
	// Add the rules to the meta policy
	meta.PluginConfigPolicy.Add(rule1, rule2, rule3)

	return meta
}
