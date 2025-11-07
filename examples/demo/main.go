package main

import (
	"fmt"
	"os"

	release "github.com/parthban-db/test-go-release"
)

func main() {
	fmt.Println("=== Go Release Conditions Demo ===\n")

	// Display build information
	displayBuildInfo()
	fmt.Println()

	// Test version conditions
	testVersionConditions()
	fmt.Println()

	// Test platform conditions
	testPlatformConditions()
	fmt.Println()

	// Run comprehensive release checks
	runReleaseChecks()
}

func displayBuildInfo() {
	fmt.Println("📦 Build Information:")
	info := release.GetBuildInfo()

	fmt.Printf("  Go Version:  %s\n", info.GoVersion)
	fmt.Printf("  Compiler:    %s\n", info.Compiler)
	fmt.Printf("  Platform:    %s\n", info.Platform)
	fmt.Printf("  OS:          %s\n", info.OS)
	fmt.Printf("  Arch:        %s\n", info.Arch)
	fmt.Printf("  CPUs:        %d\n", info.NumCPU)

	if info.VCSRevision != "" {
		fmt.Printf("  VCS Commit:  %s\n", info.VCSRevision)
		fmt.Printf("  VCS Modified: %v\n", info.VCSModified)
		fmt.Printf("  VCS Time:    %s\n", info.VCSTime)
	} else {
		fmt.Println("  VCS Info:    Not available")
	}
}

func testVersionConditions() {
	fmt.Println("🔍 Version Checks:")

	major, minor, err := release.GetGoMajorMinor()
	if err != nil {
		fmt.Printf("  Error getting version: %v\n", err)
		return
	}
	fmt.Printf("  Current Go: %d.%d\n", major, minor)

	versions := []string{"1.16", "1.20", "1.21", "1.22"}
	for _, v := range versions {
		ok, err := release.IsGoVersionAtLeast(v)
		if err != nil {
			fmt.Printf("  Error checking version %s: %v\n", v, err)
			continue
		}
		status := "✗"
		if ok {
			status = "✓"
		}
		fmt.Printf("  %s Go >= %s\n", status, v)
	}
}

func testPlatformConditions() {
	fmt.Println("🖥️  Platform Checks:")

	platforms := []struct {
		os   string
		arch string
	}{
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"windows", "amd64"},
	}

	for _, p := range platforms {
		status := "✗"
		if release.IsPlatform(p.os, p.arch) {
			status = "✓"
		}
		fmt.Printf("  %s %s/%s\n", status, p.os, p.arch)
	}
}

func runReleaseChecks() {
	fmt.Println("✅ Release Readiness Checks:")

	cs := release.NewConditionSet()

	// Add production release conditions
	cs.Add("go-version", "Go version >= 1.20", func() (bool, error) {
		return release.IsGoVersionAtLeast("1.20")
	})

	cs.Add("supported-os", "Running on Linux, macOS, or Windows", func() (bool, error) {
		return release.IsOS("linux") || release.IsOS("darwin") || release.IsOS("windows"), nil
	})

	cs.Add("64-bit-arch", "Running on 64-bit architecture", func() (bool, error) {
		return release.IsArch("amd64") || release.IsArch("arm64"), nil
	})

	cs.Add("multi-cpu", "At least 2 CPUs available", func() (bool, error) {
		info := release.GetBuildInfo()
		return info.NumCPU >= 2, nil
	})

	// Test all conditions
	results := cs.TestAll()

	for _, result := range results {
		status := "✓"
		if !result.Passed {
			status = "✗"
		}
		fmt.Printf("  %s %s: %s\n", status, result.Name, result.Description)
		if result.Error != nil {
			fmt.Printf("     Error: %v\n", result.Error)
		}
	}

	fmt.Println()
	if results.AllPassed() {
		fmt.Println("🎉 All conditions passed - Ready for release!")
		os.Exit(0)
	} else {
		fmt.Println("❌ Some conditions failed - Not ready for release")
		os.Exit(1)
	}
}
