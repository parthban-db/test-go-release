package release

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"

	"golang.org/x/mod/semver"
)

// BuildInfo contains information about the build
type BuildInfo struct {
	GoVersion   string
	Compiler    string
	Platform    string
	OS          string
	Arch        string
	NumCPU      int
	BuildTime   string
	VCSRevision string
	VCSModified bool
	VCSTime     string
}

// GetBuildInfo returns detailed build information
func GetBuildInfo() *BuildInfo {
	info := &BuildInfo{
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		NumCPU:    runtime.NumCPU(),
	}

	// Get VCS information from build info
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range buildInfo.Settings {
			switch setting.Key {
			case "vcs.revision":
				info.VCSRevision = setting.Value
			case "vcs.modified":
				info.VCSModified = setting.Value == "true"
			case "vcs.time":
				info.VCSTime = setting.Value
			}
		}
	}

	return info
}

// IsDebugMode checks if the binary is built in debug mode (no optimizations)
// This is a heuristic based on available information
func IsDebugMode() bool {
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range buildInfo.Settings {
			// Check for -N flag (disable optimizations)
			if setting.Key == "CGO_ENABLED" && setting.Value == "0" {
				continue
			}
		}
	}
	return false
}

// CompareGoVersion compares the current Go version with a target version
// Returns:
//
//	-1 if current < target
//	 0 if current == target
//	 1 if current > target
func CompareGoVersion(targetVersion string) (int, error) {
	current := runtime.Version()

	// Normalize versions for semver comparison
	currentNorm := normalizeGoVersion(current)
	targetNorm := normalizeGoVersion(targetVersion)

	if !semver.IsValid(currentNorm) {
		return 0, fmt.Errorf("invalid current version: %s", current)
	}
	if !semver.IsValid(targetNorm) {
		return 0, fmt.Errorf("invalid target version: %s", targetVersion)
	}

	return semver.Compare(currentNorm, targetNorm), nil
}

// normalizeGoVersion converts Go version format to semver format
// e.g., "go1.21.0" -> "v1.21.0"
func normalizeGoVersion(version string) string {
	version = strings.TrimPrefix(version, "go")
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}
	return version
}

// IsGoVersionAtLeast checks if the current Go version is at least the specified version
func IsGoVersionAtLeast(minVersion string) (bool, error) {
	cmp, err := CompareGoVersion(minVersion)
	if err != nil {
		return false, err
	}
	return cmp >= 0, nil
}

// GetGoMajorMinor returns the major and minor version of the current Go runtime
func GetGoMajorMinor() (major, minor int, err error) {
	version := runtime.Version()
	version = strings.TrimPrefix(version, "go")

	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return 0, 0, fmt.Errorf("invalid version format: %s", version)
	}

	major, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid major version: %s", parts[0])
	}

	minor, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid minor version: %s", parts[1])
	}

	return major, minor, nil
}

// Environment represents different deployment environments
type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvStaging     Environment = "staging"
	EnvProduction  Environment = "production"
	EnvTest        Environment = "test"
)

// Condition represents a testable release condition
type Condition struct {
	Name        string
	Description string
	Check       func() (bool, error)
}

// ConditionSet is a collection of conditions to test
type ConditionSet struct {
	conditions []Condition
}

// NewConditionSet creates a new condition set
func NewConditionSet() *ConditionSet {
	return &ConditionSet{
		conditions: make([]Condition, 0),
	}
}

// Add adds a condition to the set
func (cs *ConditionSet) Add(name, description string, check func() (bool, error)) {
	cs.conditions = append(cs.conditions, Condition{
		Name:        name,
		Description: description,
		Check:       check,
	})
}

// TestResult represents the result of testing a condition
type TestResult struct {
	Name        string
	Description string
	Passed      bool
	Error       error
}

// TestResults represents a collection of test results
type TestResults []TestResult

// TestAll tests all conditions in the set
func (cs *ConditionSet) TestAll() TestResults {
	results := make(TestResults, 0, len(cs.conditions))

	for _, cond := range cs.conditions {
		passed, err := cond.Check()
		results = append(results, TestResult{
			Name:        cond.Name,
			Description: cond.Description,
			Passed:      passed,
			Error:       err,
		})
	}

	return results
}

// AllPassed returns true if all conditions passed
func (results TestResults) AllPassed() bool {
	for _, r := range results {
		if !r.Passed || r.Error != nil {
			return false
		}
	}
	return true
}

// IsPlatform checks if the current platform matches the specified OS and architecture
func IsPlatform(os, arch string) bool {
	return runtime.GOOS == os && runtime.GOARCH == arch
}

// IsOS checks if the current OS matches the specified value
func IsOS(os string) bool {
	return runtime.GOOS == os
}

// IsArch checks if the current architecture matches the specified value
func IsArch(arch string) bool {
	return runtime.GOARCH == arch
}

// HasVCSInfo checks if VCS information is available in the build
func HasVCSInfo() bool {
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range buildInfo.Settings {
			if setting.Key == "vcs.revision" && setting.Value != "" {
				return true
			}
		}
	}
	return false
}
