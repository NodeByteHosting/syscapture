package test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/nodebytehosting/syscapture/internal/metric"
	"github.com/nodebytehosting/syscapture/internal/sysfs"
	"github.com/stretchr/testify/assert"
)

// TestHostLinux tests the GetHostInformation function
// It interacts with the host system to get the NodeName and Kernel Version
// It then compares the values with the ones returned by the GetHostInformation function
func TestHostLinux(t *testing.T) {
	// Get the NodeName and Kernel Version using shell commands
	osPlatform, osPlatformErr := sysfs.ShellExec("uname -n") // Nodename
	osKernel, osKernelErr := sysfs.ShellExec("uname -r")     // Kernel version

	// Get host information using the GetHostInformation function
	info, infoErr := metric.GetHostInformation()

	// Check for errors in getting host information
	if len(infoErr) != 0 {
		t.Error(infoErr)
		t.FailNow()
	}

	// Check for errors in getting NodeName
	if osPlatformErr != nil {
		t.Error(osPlatformErr.Error())
		t.FailNow()
	}

	// Check for errors in getting Kernel Version
	if osKernelErr != nil {
		t.Error(osKernelErr.Error())
		t.FailNow()
	}

	// Assert that the OS matches the runtime OS
	assert.Equal(t, info.Os, runtime.GOOS)

	// Assert that the Platform matches the NodeName
	assert.Equal(t, info.Platform, strings.TrimSuffix(osPlatform, "\n"))

	// Assert that the Kernel Version matches
	assert.Equal(t, info.KernelVersion, strings.TrimSuffix(osKernel, "\n"))
}
