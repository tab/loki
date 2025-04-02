package rpcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	proto "loki/internal/app/rpcs/proto/sso/v1"
)

type permissionService struct {
	proto.UnimplementedPermissionServiceServer
}

type roleService struct {
	proto.UnimplementedRoleServiceServer
}

type scopeService struct {
	proto.UnimplementedScopeServiceServer
}

func Test_Registry_RegisterAll(t *testing.T) {
	registry := NewRegistry(
		&permissionService{},
		&roleService{},
		&scopeService{},
	)
	assert.NotNil(t, registry)

	server := grpc.NewServer()
	registry.RegisterAll(server)

	serviceInfo := server.GetServiceInfo()
	assert.Contains(t, serviceInfo, "sso.v1.PermissionService")
}
