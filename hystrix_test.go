package hystrix

import (
	"context"
	"testing"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry/memory"
	"github.com/micro/go-micro/v2/router"
	rrouter "github.com/micro/go-micro/v2/router/registry"
)

func TestBreaker(t *testing.T) {
	// setup
	registry := memory.NewRegistry()

	c := client.NewClient(
		// set the selector
		client.Router(rrouter.NewRouter(router.Registry(registry))),
		// add the breaker wrapper
		client.Wrap(NewClientWrapper()),
	)

	req := c.NewRequest("test.service", "Test.Method", map[string]string{
		"foo": "bar",
	}, client.WithContentType("application/json"))

	var rsp map[string]interface{}

	// Force to point of trip
	for i := 0; i < (hystrix.DefaultVolumeThreshold * 3); i++ {
		c.Call(context.TODO(), req, rsp)
	}

	err := c.Call(context.TODO(), req, rsp)
	if err == nil {
		t.Error("Expecting tripped breaker, got nil error")
	}

	if err.Error() != "hystrix: circuit open" {
		t.Errorf("Expecting tripped breaker, got %v", err)
	}
}
