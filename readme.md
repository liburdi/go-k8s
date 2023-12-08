#### JD
了解k8s、docker

#### 前言

今年在很多场景应用了k8s去进行服务治理。operator封装k8s的job pod等，在基础服务场景中，常常收到某通知，然后进行某项工作。
比如大数据部门的数据清洗脚本结束后，会通知大数据平台进行数据处理（缓存最新数据、清理历史数据等）的脚本。
而此时起k8s job去工作，远比传统的常驻服务效率高，因为收到回调通知后才创建job，工作结束job自动释放，节省大量不必要的资源，并且job和项目之间不会互相影响（假设你的项目和job是同一个代码仓库）。

#### 运行

kubectl apply -f serviceaccount.yaml
kubectl apply -f operator.yaml //会派生出job，取决于你的operator.yaml

#### 其他

container/main.go 是收到通知要去做的事情（job）
cmd/main.go 是operator的main文件,你可以完善此处，传入不同启动命令参数，可以触发不同的job（也就是container目录下待完善的代码）