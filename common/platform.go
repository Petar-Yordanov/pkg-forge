package common

import "runtime"

type Platform string

const (
	PlatformLinux   Platform = "linux"
	PlatformMacOS   Platform = "macos"
	PlatformWindows Platform = "windows"
)

func CurrentPlatform() Platform {
	switch runtime.GOOS {
	case "windows":
		return PlatformWindows
	case "darwin":
		return PlatformMacOS
	default:
		return PlatformLinux
	}
}
