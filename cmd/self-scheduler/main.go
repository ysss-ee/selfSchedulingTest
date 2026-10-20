package main

import (
	"os"

	"k8s.io/component-base/cli"
	_ "k8s.io/component-base/logs/json/register"
	_ "k8s.io/component-base/metrics/prometheus/clientgo"
	_ "k8s.io/component-base/metrics/prometheus/version"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"k8s.io/kubernetes/cmd/self-scheduler/plugin/nodelabel"
)

func main() {
	//注册 NodeLabel 插件到调度器框架中
	command := app.NewSchedulerCommand(
		app.WithPlugin(nodelabel.Name, nodelabel.New),
	)
	os.Exit(cli.Run(command))
}
