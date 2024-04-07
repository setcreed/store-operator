package builders

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"setcreed.github.io/store/api/v1alpha1"
)

type DeployBuilder struct {
	deploy    *appsv1.Deployment
	cmBuilder *ConfigMapBuilder
	config    *v1alpha1.DbConfig

	client.Client
}

const CMAnnotation = "store.config/md5"

// 定义了命名规则
func deployName(name string) string {
	return "store-" + name
}

// 构建deployment 构造器
func NewDeployBuilder(config *v1alpha1.DbConfig, client client.Client) (*DeployBuilder, error) {
	deployment := &appsv1.Deployment{}
	err := client.Get(context.Background(), types.NamespacedName{Namespace: config.Namespace, Name: deployName(config.Name)}, deployment)
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
	// configmap 构建
	cmBuilder, err := NewConfigMapBuilder(config, client)
	if err != nil {
		fmt.Println("cm error:", err)
		return nil, err
	}

	return &DeployBuilder{
		deploy:    deployment,
		cmBuilder: cmBuilder,
		config:    config,
		Client:    client,
	}, nil
}

func (db *DeployBuilder) setCMAnnotation(configStr string) {
	db.deploy.Spec.Template.Annotations[CMAnnotation] = configStr
}

// 同步属性
func (db *DeployBuilder) apply() *DeployBuilder {
	// 同步副本
	*db.deploy.Spec.Replicas = int32(db.config.Spec.DbConfig.Replicas)
	return db
}

func (db *DeployBuilder) setOwner() *DeployBuilder {
	db.deploy.OwnerReferences = append(db.deploy.OwnerReferences, metav1.OwnerReference{
		APIVersion: db.config.APIVersion,
		Kind:       db.config.Kind,
		Name:       db.config.Name,
		UID:        db.config.UID,
	})
	return db
}

// 构建出deployment
// 包含创建和更新
func (db *DeployBuilder) Build(ctx context.Context) error {
	if db.deploy.CreationTimestamp.IsZero() {
		db.apply().setOwner()

		// 先创建configmap，再创建deployment
		err := db.cmBuilder.Build(ctx)
		if err != nil {
			return err
		}
		// 设置 md5
		db.setCMAnnotation(db.cmBuilder.DataKey)
		err = db.Create(ctx, db.deploy)
		if err != nil {
			return err
		}
	} else {
		err := db.cmBuilder.Build(ctx)
		if err != nil {
			return err
		}
		// 更新
		patch := client.MergeFrom(db.deploy.DeepCopy())
		db.apply()
		db.setCMAnnotation(db.cmBuilder.DataKey)
		err = db.Patch(ctx, db.deploy, patch)
		if err != nil {
			return err
		}

		// 获取当前deployment ready的副本数
		replicas := db.deploy.Status.ReadyReplicas
		db.config.Status.Ready = fmt.Sprintf("%d/%d", replicas, db.config.Spec.DbConfig.Replicas)
		db.config.Status.Replicas = replicas
		// 设置状态
		err = db.Client.Status().Update(ctx, db.config)
		if err != nil {
			return err
		}
	}
	return nil
}
