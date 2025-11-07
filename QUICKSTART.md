# Quick Start Guide

Get started with the Go Release Conditions library in 5 minutes!

## Installation

```bash
go get github.com/parthban-db/test-go-release
```

## Common Use Cases

### 1. Check Go Version

```go
import "github.com/parthban-db/test-go-release"

// Simple check
if ok, _ := release.IsGoVersionAtLeast("1.20"); ok {
    // Use Go 1.20+ features
}

// Get specific version components
major, minor, _ := release.GetGoMajorMinor()
if major >= 1 && minor >= 21 {
    // Use Go 1.21+ features
}
```

### 2. Platform Detection

```go
// Check OS
if release.IsOS("linux") {
    // Linux-specific code
}

// Check architecture
if release.IsArch("arm64") {
    // ARM64-specific code
}

// Check specific platform
if release.IsPlatform("darwin", "arm64") {
    // macOS Apple Silicon specific code
}
```

### 3. Build Information

```go
info := release.GetBuildInfo()
fmt.Printf("Go %s on %s\n", info.GoVersion, info.Platform)
fmt.Printf("Built from commit: %s\n", info.VCSRevision)
```

### 4. Release Validation

```go
cs := release.NewConditionSet()

cs.Add("version", "Go 1.20+", func() (bool, error) {
    return release.IsGoVersionAtLeast("1.20")
})

cs.Add("platform", "Supported platform", func() (bool, error) {
    return release.IsOS("linux") || release.IsOS("darwin"), nil
})

results := cs.TestAll()
if results.AllPassed() {
    fmt.Println("Ready for release!")
}
```

## Running the Demo

```bash
cd examples/demo
go run main.go
```

## Running Tests

```bash
# Run all tests
go test -v

# Run with coverage
go test -cover

# Run benchmarks
go test -bench=.
```

## API Quick Reference

| Function | Description |
|----------|-------------|
| `GetBuildInfo()` | Get detailed build information |
| `IsGoVersionAtLeast(v)` | Check if Go version >= v |
| `CompareGoVersion(v)` | Compare current Go version with v |
| `GetGoMajorMinor()` | Get major.minor version numbers |
| `IsPlatform(os, arch)` | Check specific platform |
| `IsOS(os)` | Check operating system |
| `IsArch(arch)` | Check architecture |
| `HasVCSInfo()` | Check if VCS info available |
| `NewConditionSet()` | Create condition test set |

## Examples in Code

Check out:
- `example_test.go` - Executable examples
- `examples/demo/main.go` - Full demo application
- `release_test.go` - Unit tests showing usage

## Need Help?

See the full [README.md](README.md) for complete documentation.

