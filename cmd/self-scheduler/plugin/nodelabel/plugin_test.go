package nodelabel

import (
	"context"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	framework "k8s.io/kube-scheduler/framework"
)

func TestFilter(t *testing.T) {
	plugin := &Plugin{}
	tests := []struct {
		name   string
		labels map[string]string
		want   framework.Code
	}{
		{name: "enabled node passes", labels: map[string]string{RequiredLabel: "true"}, want: framework.Success},
		{name: "unlabeled node is rejected", labels: nil, want: framework.Unschedulable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := plugin.Filter(context.Background(), nil, &v1.Pod{}, nodeInfo(tt.labels))
			got := framework.Success
			if status != nil {
				got = status.Code()
			}
			if got != tt.want {
				t.Fatalf("Filter() code = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScore(t *testing.T) {
	plugin := &Plugin{}
	tests := []struct {
		name   string
		labels map[string]string
		want   int64
	}{
		{name: "preferred node receives maximum score", labels: map[string]string{PreferredLabel: "true"}, want: framework.MaxNodeScore},
		{name: "ordinary node receives zero score", labels: nil, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, status := plugin.Score(context.Background(), nil, &v1.Pod{}, nodeInfo(tt.labels))
			if status != nil {
				t.Fatalf("Score() status = %v", status)
			}
			if got != tt.want {
				t.Fatalf("Score() = %d, want %d", got, tt.want)
			}
		})
	}
}

func nodeInfo(labels map[string]string) framework.NodeInfo {
	return fakeNodeInfo{node: &v1.Node{ObjectMeta: metav1.ObjectMeta{Name: "test-node", Labels: labels}}}
}

// fakeNodeInfo embeds the framework interface so unneeded methods remain
// available to the compiler. This plugin reads only Node().
type fakeNodeInfo struct {
	framework.NodeInfo
	node *v1.Node
}

func (f fakeNodeInfo) Node() *v1.Node {
	return f.node
}
