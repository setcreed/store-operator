package dashboard

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"setcreed.github.io/store/api/v1alpha1"
)

type AdminUI struct {
	r      *gin.Engine
	client client.Client
}

func NewAdminUI(c client.Client) *AdminUI {
	r := gin.New()
	r.StaticFS("/adminui", http.Dir("./adminui"))
	r.Use(errorHandler())
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	return &AdminUI{r: r, client: c}
}

func (au *AdminUI) Start(ctx context.Context) error {
	au.r.GET("/configs", func(c *gin.Context) {
		list := &v1alpha1.DbConfigList{}
		err := au.client.List(ctx, list)
		if err != nil {
			klog.Error(err)
			return
		}
		c.JSON(200, list.Items)
	})

	au.r.POST("/configs", func(c *gin.Context) {
		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}
		cfg := &v1alpha1.DbConfig{}
		err = yaml.Unmarshal(b, cfg)
		if err != nil {
			panic(err)
		}
		if cfg.Namespace == "" {
			cfg.Namespace = "default"
		}
		err = au.client.Create(c, cfg)
		if err != nil {
			klog.Error(err)
			return
		}
		c.JSON(200, gin.H{"message": "success"})
	})

	au.r.DELETE("/configs/:ns/:name", func(c *gin.Context) {
		cfg := &v1alpha1.DbConfig{}
		err := au.client.Get(c, types.NamespacedName{
			Name:      c.Param("name"),
			Namespace: c.Param("ns"),
		}, cfg)
		if err != nil {
			panic(err)
		}
		err = au.client.Delete(c, cfg)
		if err != nil {
			panic(err)
		}
		c.JSON(200, gin.H{"message": "OK"})

	})

	au.r.GET("/events/:ns/:name", func(c *gin.Context) {
		var ns, name = c.Param("ns"), c.Param("name")

		cfg := &v1alpha1.DbConfig{}
		err := au.client.Get(c, types.NamespacedName{
			Name:      name,
			Namespace: ns,
		}, cfg)
		if err != nil {
			klog.Error(err)
			return
		}

		list := &corev1.EventList{}
		err = au.client.List(c, list, &client.ListOptions{Namespace: ns})
		if err != nil {
			klog.Error(err)
			return
		}
		ret := []corev1.Event{}
		for _, e := range list.Items {
			// 匹配 自定义资源 对应的 event
			if e.InvolvedObject.Name == name && e.InvolvedObject.UID == cfg.UID {
				ret = append(ret, e)
				continue
			}
			// 当前资源是否是dbconfig 创建出来的 deployment
			if isChildDeploy(cfg, e.InvolvedObject, au.client) {
				ret = append(ret, e)
				continue
			}
		}
		c.JSON(200, ret)

	})

	return au.r.Run(":9003")
}

func errorHandler() gin.HandlerFunc {
	gin.ErrorLogger()
	return func(c *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				c.AbortWithStatusJSON(400, gin.H{"error": e.(error).Error()})
			}
		}()
		c.Next()
	}
}

func isChildDeploy(cfg *v1alpha1.DbConfig, or corev1.ObjectReference, client client.Client) bool {
	//获取Deployment .命名规则是 store-xxxx
	dep := &appv1.Deployment{}
	err := client.Get(context.Background(), types.NamespacedName{
		Name:      "store-" + cfg.Name,
		Namespace: cfg.Namespace,
	}, dep)
	if err != nil {
		return false
	}
	if or.UID == dep.UID {
		return true
	}
	return false

}
