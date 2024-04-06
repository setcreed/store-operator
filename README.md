# store-operator

创建项目
```
kubebuilder init --domain setcreed.github.io --repo setcreed.github.io/store
```

创建API
```
kubebuilder create api --group apps --version v1alpha1 --kind DbConfig
```

创建webhook
```
kubebuilder create webhook --group apps --version v1alpha1 --kind DbConfig --defaulting --programmatic-validation
```
通过设置--defaulting可创建mutatingadmissionwebhook类型准入控制器，用来修改传入资源；参数--programmatic-validation可创建validatingadmissionwebhook，用来验证传入资源

