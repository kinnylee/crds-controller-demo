### 创建自定义资源：
```bash
kubectl apply -f crds/scaling_crd.yaml
kubectl get crds |grep scaling
```

### 创建自定义资源对象

```bash
kubectl apply -f crds/scaling_test.yaml
kubectl get scaling
# 或者
kubectl get sca

kubectl get sca scalingtest -o yaml
```

### 创建含校验功能的资源对象

```bash

```

### 用code-generate生成代码

```bash
go get k8s.io/code-generator
./hack/update-codegen.sh
```