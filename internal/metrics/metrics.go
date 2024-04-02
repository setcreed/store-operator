package metrics

import "github.com/prometheus/client_golang/prometheus"

var StoreReconcileTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "store_reconcile_total",
	Help: "我们自己写的reconcile触发统计",
}, []string{"controller"})
