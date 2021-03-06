package docker

import (
	"testing"

	"github.com/Netflix/titus-executor/api/netflix/titus"
	runtimeTypes "github.com/Netflix/titus-executor/executor/runtime/types"
	"github.com/docker/docker/api/types/container"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestDefaultProfile(t *testing.T) {
	c := runtimeTypes.Container{
		TitusInfo: &titus.ContainerInfo{},
	}
	hostConfig := container.HostConfig{}

	assert.NoError(t, setupAdditionalCapabilities(&c, &hostConfig))

	assert.Len(t, hostConfig.CapAdd, 0)
	assert.Len(t, hostConfig.CapDrop, 0)
	assert.Len(t, hostConfig.SecurityOpt, 1)
}

func TestFuseProfile(t *testing.T) {
	c := runtimeTypes.Container{
		TitusInfo: &titus.ContainerInfo{
			PassthroughAttributes: map[string]string{
				runtimeTypes.FuseEnabledParam: "true",
			},
		},
	}
	hostConfig := container.HostConfig{}

	assert.NoError(t, setupAdditionalCapabilities(&c, &hostConfig))

	assert.Contains(t, hostConfig.CapAdd, "SYS_ADMIN")
	assert.Len(t, hostConfig.CapDrop, 0)
	assert.Len(t, hostConfig.SecurityOpt, 2)
	assert.Contains(t, hostConfig.SecurityOpt, "apparmor:docker-fuse")
}

func TestNestedContainerProfile(t *testing.T) {
	c := runtimeTypes.Container{
		Env: map[string]string{},
		TitusInfo: &titus.ContainerInfo{
			AllowNestedContainers: proto.Bool(true),
		},
	}
	hostConfig := container.HostConfig{}

	assert.NoError(t, setupAdditionalCapabilities(&c, &hostConfig))

	assert.Contains(t, hostConfig.CapAdd, "SYS_ADMIN")
	assert.Len(t, hostConfig.CapDrop, 0)
	assert.Len(t, hostConfig.SecurityOpt, 2)
	assert.Contains(t, hostConfig.SecurityOpt, "apparmor:docker-nested")

}

func TestFuseAndNestedContainerProfileProfile(t *testing.T) {
	c := runtimeTypes.Container{
		Env: map[string]string{},
		TitusInfo: &titus.ContainerInfo{
			AllowNestedContainers: proto.Bool(true),
			PassthroughAttributes: map[string]string{
				runtimeTypes.FuseEnabledParam: "true",
			},
		},
	}
	hostConfig := container.HostConfig{}

	assert.NoError(t, setupAdditionalCapabilities(&c, &hostConfig))

	assert.Contains(t, hostConfig.CapAdd, "SYS_ADMIN")
	assert.Len(t, hostConfig.CapDrop, 0)
	assert.Len(t, hostConfig.SecurityOpt, 2)
	assert.Contains(t, hostConfig.SecurityOpt, "apparmor:docker-nested")
}
