# k8s-camp

## 2023-01-15 更新内容 
添加 Pod类型对象 yaml 和 configMap类型对象yaml
* 1、探活 使用livenessProbe 探测 web /healthz接口 
* 2、资源需求 设置 100Mi ～ 500Mi 内存要求
* 3、使用 env 固定设置环境参数
* 4、日常运维日志级别，通过configMap 注入到 Pod环境变量，应用启动读取对应参数 本地参数通过 /healthz 直接对外输出
