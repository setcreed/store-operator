package builders

import (
	"bytes"
	"context"
	"text/template"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"setcreed.github.io/store/api/v1alpha1"
)

type DeployBuilder struct {
	deploy *v1.Deployment
	config *v1alpha1.DbConfig

	client.Client
}

// 定义了命名规则
func deployName(name string) string {
	return "store-" + name
}

// 构建deployment 构造器
func NewDeployBuilder(config *v1alpha1.DbConfig, c client.Client) (*DeployBuilder, error) {
	deployment := &v1.Deployment{}
	err := c.Get(context.Background(), client.ObjectKey{Namespace: config.Namespace, Name: deployName(config.Name)}, deployment)
	if err != nil {
		// 表示deployment不存在
		deployment.Namespace, deployment.Name = config.Namespace, config.Name
		tpl, err := template.New("deploy").Parse(deployTemplate)
		if err != nil {
			return nil, err
		}
		var tplRet bytes.Buffer
		err = tpl.Execute(&tplRet, deployment)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(tplRet.Bytes(), deployment)
		if err != nil {
			return nil, err
		}
	}

	return &DeployBuilder{deploy: deployment, Client: c, config: config}, nil
}

// 同步属性
func (db *DeployBuilder) apply() *DeployBuilder {
	// 同步副本
	db.deploy.Spec.Replicas = db.config.Spec.Replicas
	return db
}

// 构建出deployment
// 包含创建和更新
func (db *DeployBuilder) Build() error {
	if db.deploy.CreationTimestamp.IsZero() {
		db.apply()
		err := db.Create(context.Background(), db.deploy)
		if err != nil {
			return err
		}
	} else {
		// 更新
		patch := client.MergeFrom(db.deploy.DeepCopy())
		db.apply()
		err := db.Patch(context.Background(), db.deploy, patch)
		if err != nil {
			return err
		}
	}
	return nil
}
