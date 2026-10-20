# 基于 Scheduler Framework 的自定义调度器

插件名为 `NodeLabel`，它接入两个 Scheduler Framework 扩展点：

1. **Filter**：只有带 `scheduler.example.com/enabled=true` 标签的节点才能通过。
2. **Score**：在已通过 Filter 的节点中，带 `scheduler.example.com/preferred=true` 的节点得到 `100` 分；其他节点得到 `0` 分。

## 结构

```text
Pod（schedulerName=self-scheduler）
              │
              ▼
   kube-scheduler 调度框架
              │
      ┌───────┴────────┐
      ▼                ▼
 NodeLabel Filter   NodeLabel Score
  enabled=true      preferred=true → 100
      │                │
      └───────┬────────┘
              ▼
  默认插件继续处理并完成 Bind
```

## 构建镜像

```powershell
cd D:\k8s\selfSchedulingTest
docker build -t example/self-scheduler:v1.36.1 .
```

确保目标集群能使用该镜像：Docker Desktop Kubernetes 通常可直接使用本地镜像；kind 需要执行 `kind load docker-image example/self-scheduler:v1.36.1`；其他集群则需要推送到镜像仓库，并修改 `manifests/scheduler.yaml` 中的 `image`。

## 运行与验证

首先为测试节点添加标签。至少需要一个 `enabled=true` 节点；如果有多个节点，再给其中一个添加 `preferred=true` 以观察打分偏好。

```powershell
kubectl get nodes
kubectl label node <候选节点> scheduler.example.com/enabled=true
kubectl label node <首选节点> scheduler.example.com/enabled=true scheduler.example.com/preferred=true
```

部署调度器和测试 Pod：

```powershell
kubectl apply -f manifests/scheduler.yaml
kubectl -n self-scheduler-system logs deploy/self-scheduler -f

kubectl apply -f manifests/demo-pod.yaml
kubectl get pod self-scheduler-demo -w
kubectl get pod self-scheduler-demo -o wide
```

若没有任何节点带 `scheduler.example.com/enabled=true`，Pod 会保持 `Pending`。

## 代码阅读顺序

1. 先看 `plugin.go` 中的 `Filter`：它决定一个节点是否有资格继续参与调度。
2. 再看 `Score`：对通过所有 Filter 的节点执行，并打分。
3. 查看 `main.go`：如何注册到 kube-scheduler。
4. 查看 ConfigMap：`schedulerName` 和 `plugins.filter/score.enabled` 如何启用插件。

## 下一步

建议在这个最小示例跑通后，按顺序扩展：

1. 用 `PreFilter` 预先读取 GPU / 资源请求。
2. 在 `Filter` 中淘汰 GPU 型号、显存或节点标签不满足要求的节点。
3. 在 `Score` 中按剩余 GPU、显存碎片或业务优先级计算得分。
4. 学习 `Reserve`、`Permit`、`PreBind` 和 `PostFilter` 等扩展点。
