# k8s-camp

## 2023-01-15 更新内容 
添加 Pod类型对象 yaml 和 configMap类型对象yaml
* 1、探活 使用livenessProbe 探测 web /healthz接口 
* 2、资源需求 设置 100Mi ～ 500Mi 内存要求
* 3、使用 env 固定设置环境参数
* 4、日常运维日志级别，通过configMap 注入到 Pod环境变量，应用启动读取对应参数 本地参数通过 /healthz 直接对外输出

## 2023-01-21 更新内容
对http服务 创建两个 入口/business 和 /sale，功能为读取环境参数。相同代码分别发布 business 和 sale 两组deployment，区别在与环境参数的变量值（各有2个duplicateSet 保持冗余），分别创建service。
最后使用ingress的 对不同入口，分别指向两个不同的后段服务。 通过入口调用可以得到不同的环境参数值，代表配置准确性。
