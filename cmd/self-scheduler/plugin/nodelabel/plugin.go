package nodelabel

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	framework "k8s.io/kube-scheduler/framework"
)

const (
	Name = "NodeLabel"

	RequiredLabel = "scheduler.example.com/enabled"
	PreferredLabel = "scheduler.example.com/preferred"
)

var _ framework.FilterPlugin = &Plugin{}
var _ framework.ScorePlugin = &Plugin{}


type Plugin struct{}

func New(_ context.Context, _ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
	return &Plugin{}, nil
}

func (p *Plugin) Name() string {
	return Name
}

// Filter 
func (p *Plugin) Filter(_ context.Context, _ framework.CycleState, _ *v1.Pod, nodeInfo framework.NodeInfo) *framework.Status {
	node := nodeInfo.Node()
	if node == nil {
		return framework.NewStatus(framework.Error, "node object is missing")
	}
	if node.Labels[RequiredLabel] != "true" {
		return framework.NewStatus(framework.Unschedulable, "node is not enabled for self-scheduler")
	}
	return nil
}

// Score 
func (p *Plugin) Score(_ context.Context, _ framework.CycleState, _ *v1.Pod, nodeInfo framework.NodeInfo) (int64, *framework.Status) {
	node := nodeInfo.Node()
	if node == nil {
		return 0, framework.NewStatus(framework.Error, "node object is missing")
	}
	if node.Labels[PreferredLabel] == "true" {
		return framework.MaxNodeScore, nil
	}
	return 0, nil
}

func (p *Plugin) ScoreExtensions() framework.ScoreExtensions {
	return nil
}
