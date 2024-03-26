package builders

import (
	"context"
	"fmt"
	"setcreed.github.io/store/api/v1alpha1"
	"testing"
)

func TestNewDeployBuilder(t *testing.T) {
	config := &v1alpha1.DbConfig{}
	config.Namespace = "default"
	config.Name = "test"
	builder, err := NewDeployBuilder(config, nil)
	if err != nil {
		t.Error(err)
	}
	build := builder.Build(context.Background())
	t.Log(fmt.Sprintf("%+v", build))
}
