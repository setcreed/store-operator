package builders

import (
	"bytes"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"setcreed.github.io/store/internal/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"text/template"

	"setcreed.github.io/store/api/v1alpha1"
)

type ConfigMapBuilder struct {
	cm     *corev1.ConfigMap
	config *v1alpha1.DbConfig
	client.Client
	DataKey string // 将configmap中的config.yaml的数据 进行md5加密后存入DataKey
}

const configMapKey = "config.yaml"

func NewConfigMapBuilder(config *v1alpha1.DbConfig, client client.Client) (*ConfigMapBuilder, error) {
	cm := &corev1.ConfigMap{}
	err := client.Get(context.Background(), types.NamespacedName{Namespace: config.Namespace, Name: deployName(config.Name)}, cm)
	if err != nil {
		// 没取到，赋值
		cm.Namespace, cm.Name = config.Namespace, deployName(config.Name)
		cm.Data = make(map[string]string)
	}
	return &ConfigMapBuilder{cm: cm, config: config, Client: client}, nil
}

func (cb *ConfigMapBuilder) setOwner() *ConfigMapBuilder {
	cb.cm.OwnerReferences = append(cb.cm.OwnerReferences,
		metav1.OwnerReference{
			APIVersion: cb.config.APIVersion,
			Kind:       cb.config.Kind,
			Name:       cb.config.Name,
			UID:        cb.config.UID,
		})
	return cb
}

func (cb *ConfigMapBuilder) apply() *ConfigMapBuilder {
	tpl, err := template.New("config").Delims("[[", "]]").Parse(cmTemplate)
	if err != nil {
		fmt.Printf("parse template error:%v\n", err)
		return cb
	}
	var tplRet bytes.Buffer
	err = tpl.Execute(&tplRet, cb.config.Spec)
	if err != nil {
		fmt.Printf("execute template error:%v\n", err)
		return cb
	}
	cb.cm.Data[configMapKey] = tplRet.String()
	return cb
}

// configmap构建
func (cb *ConfigMapBuilder) Build(ctx context.Context) error {
	if cb.cm.CreationTimestamp.IsZero() {
		cb.apply().setOwner().parseKey()
		err := cb.Create(ctx, cb.cm)
		if err != nil {
			return err
		}
	} else {
		patch := client.MergeFrom(cb.cm.DeepCopy())
		cb.apply().parseKey()
		err := cb.Patch(ctx, cb.cm, patch)
		if err != nil {
			return err
		}
	}
	return nil
}

// 把configmap里面的 key=config.yaml的内容 取出变成md5
func (cb *ConfigMapBuilder) parseKey() *ConfigMapBuilder {
	if appData, ok := cb.cm.Data[configMapKey]; ok {
		cb.DataKey = utils.Md5(appData)
		return cb
	}
	cb.DataKey = ""
	return cb
}
