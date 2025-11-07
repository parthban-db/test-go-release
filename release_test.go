package release

import (
	"runtime"
	"testing"
)

func TestGetBuildInfo(t *testing.T) {
	info := GetBuildInfo()

	if info == nil {
		t.Fatal("GetBuildInfo returned nil")
	}

	if info.GoVersion == "" {
		t.Error("GoVersion should not be empty")
	}

	if info.OS == "" {
		t.Error("OS should not be empty")
	}

	if info.Arch == "" {
		t.Error("Arch should not be empty")
	}

	if info.NumCPU <= 0 {
		t.Error("NumCPU should be positive")
	}

	t.Logf("Build Info: %+v", info)
}

func TestCompareGoVersion(t *testing.T) {
	tests := []struct {
		name    string
		target  string
		wantErr bool
	}{
		{"Valid version 1.20", "1.20", false},
		{"Valid version go1.20", "go1.20", false},
		{"Valid version 1.21.0", "1.21.0", false},
		{"Invalid version", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CompareGoVersion(tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareGoVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				t.Logf("Current version compared to %s: %d", tt.target, result)
			}
		})
	}
}

func TestIsGoVersionAtLeast(t *testing.T) {
	// Test with a very old version
	result, err := IsGoVersionAtLeast("1.10")
	if err != nil {
		t.Errorf("IsGoVersionAtLeast() error = %v", err)
	}
	if !result {
		t.Error("Current Go version should be at least 1.10")
	}

	// Test with a future version (should fail)
	result, err = IsGoVersionAtLeast("99.99")
	if err != nil {
		t.Errorf("IsGoVersionAtLeast() error = %v", err)
	}
	if result {
		t.Error("Current Go version should not be at least 99.99")
	}
}

func TestGetGoMajorMinor(t *testing.T) {
	major, minor, err := GetGoMajorMinor()
	if err != nil {
		t.Fatalf("GetGoMajorMinor() error = %v", err)
	}

	if major < 1 {
		t.Errorf("Major version should be at least 1, got %d", major)
	}

	if minor < 0 {
		t.Errorf("Minor version should be non-negative, got %d", minor)
	}

	t.Logf("Go version: %d.%d", major, minor)
}

func TestIsPlatform(t *testing.T) {
	// Test current platform
	if !IsPlatform(runtime.GOOS, runtime.GOARCH) {
		t.Error("IsPlatform should return true for current platform")
	}

	// Test non-existent platform
	if IsPlatform("fakeos", "fakearch") {
		t.Error("IsPlatform should return false for non-existent platform")
	}
}

func TestIsOS(t *testing.T) {
	if !IsOS(runtime.GOOS) {
		t.Error("IsOS should return true for current OS")
	}

	if IsOS("fakeos") {
		t.Error("IsOS should return false for non-existent OS")
	}
}

func TestIsArch(t *testing.T) {
	if !IsArch(runtime.GOARCH) {
		t.Error("IsArch should return true for current architecture")
	}

	if IsArch("fakearch") {
		t.Error("IsArch should return false for non-existent architecture")
	}
}

func TestConditionSet(t *testing.T) {
	cs := NewConditionSet()

	// Add some test conditions
	cs.Add("Go Version Check", "Check if Go version is at least 1.10", func() (bool, error) {
		return IsGoVersionAtLeast("1.10")
	})

	cs.Add("Platform Check", "Check if running on a supported platform", func() (bool, error) {
		return IsPlatform(runtime.GOOS, runtime.GOARCH), nil
	})

	cs.Add("CPU Check", "Check if we have at least 1 CPU", func() (bool, error) {
		return runtime.NumCPU() >= 1, nil
	})

	results := cs.TestAll()

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	for _, result := range results {
		t.Logf("Condition: %s - Passed: %v - Error: %v", result.Name, result.Passed, result.Error)
		if !result.Passed {
			t.Errorf("Condition %s failed", result.Name)
		}
		if result.Error != nil {
			t.Errorf("Condition %s returned error: %v", result.Name, result.Error)
		}
	}

	if !results.AllPassed() {
		t.Error("Not all conditions passed")
	}
}

func TestNormalizeGoVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"go1.21.0", "v1.21.0"},
		{"1.21.0", "v1.21.0"},
		{"v1.21.0", "v1.21.0"},
		{"go1.20", "v1.20"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeGoVersion(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeGoVersion(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHasVCSInfo(t *testing.T) {
	// This test just ensures the function doesn't panic
	result := HasVCSInfo()
	t.Logf("Has VCS Info: %v", result)
}

func BenchmarkGetBuildInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetBuildInfo()
	}
}

func BenchmarkCompareGoVersion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CompareGoVersion("1.20")
	}
}
