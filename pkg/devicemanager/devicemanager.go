package devicemanager

import (
	"github.com/yqs112358/cross-clipboard/pkg/config"
	"github.com/yqs112358/cross-clipboard/pkg/device"
)

type DeviceManager struct {
	Devices        map[string]*device.Device
	DevicesUpdated chan struct{}

	config *config.Config
}

func NewDeviceManager(cfg *config.Config) *DeviceManager {
	return &DeviceManager{
		Devices:        make(map[string]*device.Device),
		DevicesUpdated: make(chan struct{}),
		config:         cfg,
	}
}

func (dm *DeviceManager) AddDevice(device *device.Device) {
	dm.Devices[device.AddressInfo.ID.Pretty()] = device
	dm.DevicesUpdated <- struct{}{}
}

func (dm *DeviceManager) RemoveDevice(device *device.Device) {
	// Flush and close ignore error
	device.Writer.Flush()
	device.Stream.Close()
	delete(dm.Devices, device.AddressInfo.ID.Pretty())
	dm.DevicesUpdated <- struct{}{}
}

func (dm *DeviceManager) GetDevice(id string) *device.Device {
	return dm.Devices[id]
}

func (dm *DeviceManager) UpdateDevice(device *device.Device) {
	dm.Devices[device.AddressInfo.ID.Pretty()] = device
	dm.DevicesUpdated <- struct{}{}
	dm.Save()
}
