package cache

import (
	"github.com/code-ready/crc/pkg/crc/constants"
	"github.com/code-ready/crc/pkg/crc/machine/libvirt"
)

func NewMachineDriverLibvirtCache() *Cache {
	return New(libvirt.MachineDriverCommand, libvirt.MachineDriverDownloadURL, constants.CrcBinDir, libvirt.MachineDriverVersion, getCurrentLibvirtDriverVersion)
}

func getCurrentLibvirtDriverVersion(executablePath string) (string, error) {
	return getVersionGeneric(executablePath, "version")
}
