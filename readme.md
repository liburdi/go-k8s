# 基于go的k8s operator

#### 运行
```
kubectl apply -f serviceaccount.yaml
kubectl apply -f operator.yaml //会派生出job，取决于你的operator.yaml
```

#### 其他

`container/main.go` 是收到通知要去做的事情（job）

`cmd/main.go` 是`operator`的`main`文件,你可以完善此处，传入不同启动命令参数，可以触发不同的`job`
