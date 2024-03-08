# # dead-mans-switch
本项目是一个简单的用于接收alertmanager的watchdog类告警，并进行反转逻辑（即不告警的时候就是告警），然后把下线告警转发给钉钉的转发器。
      
## Develop
本地环境编译：
```sh
make build
```

本地环境运行：
```sh
dms -config ./manifest/config.example.yaml
```

发送一个webhook告警（用于测试）：
```sh
curl -H "Content-Type: application/json" --data @payload.json http://localhost:8080/webhook
```

## Deploy

`manifest/deploy` 目录有 k8s 部署 yaml 文件，您可以复制它并更新 configmap和deployment。
manifest/monitoring 目录下有 ServiceMonitor 和 PrometheusRule crd 文件，如果你使用 prometheus-operator 监控你的 k8s 集群，可以尝试一下。

### AlertManager config
让watchdog告警发送到本接收器:
```yaml
route:
  routes:
    - receiver: dead-mans-switch
      group_wait: 10s
      group_interval: 30s
      repeat_interval: 15s
      match:
        alertname: 'Watchdog'
```


添加 Dead Mans Switch 服务作为新的 Webhook 接收器:
```yaml
receivers:
- name: dead-mans-switch
  webhook_configs:
  - url: http://dead-mans-switch:8080/webhook/alertmanager
```