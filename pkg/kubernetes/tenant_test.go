package kubernetes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTenant(t *testing.T) {
	t.Run("tenantCreate", func(t *testing.T) {
		tenantTitle := "tenantTest1"
		data, err := TenantCreate(tenantTitle)
		assert.Nil(t, err, "unexpected error")
		assert.NotNil(t, data, "Expected create tenant empty")
		t.Log(data)
	})
	t.Run("tenantList", func(t *testing.T) {
		data, err := TenantList()
		assert.Nil(t, err, "unexpected error")
		assert.NotNil(t, data, "Expected list tenant empty")
		t.Log(data)
	})
}
