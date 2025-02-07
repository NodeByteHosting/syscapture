package test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/nodebytehosting/syscapture/internal/metric"
	"github.com/nodebytehosting/syscapture/internal/sysfs"

	"github.com/stretchr/testify/assert"
)

const (
	cmdNodeName       = "uname -n"
	cmdKernelVersion  = "uname -r"
	errorGetHostInfo  = "Error getting host information: "
	errorGetNodeName  = "Error getting node name: "
	errorGetKernelVer = "Error getting kernel version: "
)

// TestHostLinux tests the GetHostInformation function
// It interacts with the host system to get the NodeName and Kernel Version
// It then compares the values with the ones returned by the GetHostInformation function
func TestHostLinux(t *testing.T) {
	osPlatform, osPlatformErr := sysfs.ShellExec(cmdNodeName)  // Nodename
	osKernel, osKernelErr := sysfs.ShellExec(cmdKernelVersion) // Kernel version
	info, infoErr := metric.GetHostInformation()

	if infoErr != nil {
		t.Error(errorGetHostInfo + infoErr.Error())
		t.FailNow()
	}

	if osKernelErr != nil {
		t.Error(errorGetKernelVer + osKernelErr.Error())
		t.FailNow()
	}

	if osPlatformErr != nil {
		t.Error(errorGetNodeName + osPlatformErr.Error())
		t.FailNow()
	}

	assert.Equal(t, runtime.GOOS, info.Os)
	assert.Equal(t, strings.TrimSuffix(osPlatform, "\n"), info.Platform)
	assert.Equal(t, strings.TrimSuffix(osKernel, "\n"), info.KernelVersion)
}
