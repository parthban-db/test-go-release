package release_test

import (
	"fmt"
	"runtime"

	release "github.com/parthban-db/test-go-release"
)

// Example demonstrates basic usage of the release library
func Example() {
	// Get build information
	info := release.GetBuildInfo()
	fmt.Printf("Go Version: %s\n", info.GoVersion)
	fmt.Printf("Platform: %s\n", info.Platform)

	// Check Go version
	isAtLeast, _ := release.IsGoVersionAtLeast("1.20")
	fmt.Printf("Go version >= 1.20: %v\n", isAtLeast)

	// Check platform
	isLinux := release.IsOS("linux")
	fmt.Printf("Running on Linux: %v\n", isLinux)
}

// ExampleConditionSet demonstrates how to create and test release conditions
func ExampleConditionSet() {
	cs := release.NewConditionSet()

	// Add production readiness checks
	cs.Add("Go Version", "Check minimum Go version", func() (bool, error) {
		return release.IsGoVersionAtLeast("1.20")
	})

	cs.Add("Platform", "Check if running on supported platform", func() (bool, error) {
		return release.IsOS("linux") || release.IsOS("darwin"), nil
	})

	cs.Add("Architecture", "Check if running on 64-bit architecture", func() (bool, error) {
		return release.IsArch("amd64") || release.IsArch("arm64"), nil
	})

	// Test all conditions
	results := cs.TestAll()

	fmt.Println("Release Condition Tests:")
	for _, result := range results {
		status := "✓"
		if !result.Passed {
			status = "✗"
		}
		fmt.Printf("%s %s: %v\n", status, result.Name, result.Passed)
	}

	if results.AllPassed() {
		fmt.Println("All conditions passed - ready for release!")
	}
}

// ExampleCompareGoVersion demonstrates version comparison
func ExampleCompareGoVersion() {
	cmp, err := release.CompareGoVersion("1.20")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	switch cmp {
	case -1:
		fmt.Println("Current Go version is older than 1.20")
	case 0:
		fmt.Println("Current Go version is exactly 1.20")
	case 1:
		fmt.Println("Current Go version is newer than 1.20")
	}
}

// ExampleGetBuildInfo demonstrates retrieving build information
func ExampleGetBuildInfo() {
	info := release.GetBuildInfo()

	fmt.Println("Build Information:")
	fmt.Printf("  Go Version: %s\n", info.GoVersion)
	fmt.Printf("  Compiler: %s\n", info.Compiler)
	fmt.Printf("  OS: %s\n", info.OS)
	fmt.Printf("  Architecture: %s\n", info.Arch)
	fmt.Printf("  Number of CPUs: %d\n", info.NumCPU)

	if info.VCSRevision != "" {
		fmt.Printf("  VCS Revision: %s\n", info.VCSRevision)
		fmt.Printf("  VCS Modified: %v\n", info.VCSModified)
		fmt.Printf("  VCS Time: %s\n", info.VCSTime)
	}
}

// ExampleIsPlatform demonstrates platform checking
func ExampleIsPlatform() {
	// Check specific platforms
	if release.IsPlatform("linux", "amd64") {
		fmt.Println("Running on Linux AMD64")
	} else if release.IsPlatform("darwin", "arm64") {
		fmt.Println("Running on macOS ARM64 (Apple Silicon)")
	} else if release.IsPlatform("windows", "amd64") {
		fmt.Println("Running on Windows AMD64")
	} else {
		fmt.Printf("Running on %s/%s\n", runtime.GOOS, runtime.GOARCH)
	}
}

// ExampleGetGoMajorMinor demonstrates getting Go version components
func ExampleGetGoMajorMinor() {
	major, minor, err := release.GetGoMajorMinor()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Go version: %d.%d\n", major, minor)

	// Use for conditional logic
	if major > 1 || (major == 1 && minor >= 21) {
		fmt.Println("Using Go 1.21+ features")
	}
}
