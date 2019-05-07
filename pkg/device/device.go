package device

import (
	"plugin"

	"github.com/microsoft/KubeDevice-API/pkg/types"
)

type Mount struct {
	ContainerPath string
	HostPath      string
	Readonly      bool
}

// Device is a device to use
type Device interface {
	// New creates the device and initializes it
	New() error
	// Start logically initializes the device
	Start() error
	// UpdateNodeInfo - updates a node info structure by writing capacity, allocatable, used, scorer
	UpdateNodeInfo(*types.NodeInfo) error
	// Allocate attempts to allocate the devices
	// Returns list of Mounts, and list of Devices to use
	// Returns an error on failure.
	Allocate(*types.PodInfo, *types.ContainerInfo) ([]Mount, []string, map[string]string, error)
	// GetName returns the name of a device
	GetName() string
}

// CreateDeviceFromPlugin returns a device from a plugin name
func CreateDeviceFromPlugin(pluginName string) (Device, error) {
	p, err := plugin.Open(pluginName)
	if err != nil {
		return nil, err
	}
	f, err := p.Lookup("CreateDevicePlugin")
	if err != nil {
		return nil, err
	}
	d, err := f.(func() (Device, error))()
	return d, err
}
