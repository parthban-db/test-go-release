# Go Release Conditions Library

A small Go library for testing various Go release conditions, runtime information, and build metadata. This library is useful for:

- Testing Go version compatibility
- Checking runtime platform and architecture
- Validating release readiness conditions
- Gathering build and VCS information
- Creating conditional behavior based on Go runtime

## Features

- 🔍 **Version Detection**: Check and compare Go versions
- 🖥️ **Platform Detection**: Identify OS and architecture
- 📦 **Build Information**: Access build metadata and VCS info
- ✅ **Condition Testing**: Create and test custom release conditions
- 🚀 **Production Ready**: Lightweight with no external runtime dependencies (only uses `golang.org/x/mod` for semver)

## Installation

```bash
go get github.com/parthban-db/test-go-release
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/parthban-db/test-go-release"
)

func main() {
    // Get build information
    info := release.GetBuildInfo()
    fmt.Printf("Running Go %s on %s\n", info.GoVersion, info.Platform)
    
    // Check Go version
    if ok, _ := release.IsGoVersionAtLeast("1.20"); ok {
        fmt.Println("Go 1.20+ detected")
    }
    
    // Check platform
    if release.IsOS("linux") {
        fmt.Println("Running on Linux")
    }
}
```

## API Reference

### Build Information

#### `GetBuildInfo() *BuildInfo`

Returns detailed information about the current build:

```go
info := release.GetBuildInfo()
fmt.Printf("Go Version: %s\n", info.GoVersion)
fmt.Printf("OS: %s\n", info.OS)
fmt.Printf("Arch: %s\n", info.Arch)
fmt.Printf("Platform: %s\n", info.Platform)
fmt.Printf("Compiler: %s\n", info.Compiler)
fmt.Printf("CPUs: %d\n", info.NumCPU)
```

`BuildInfo` struct contains:
- `GoVersion`: Go version string (e.g., "go1.21.0")
- `Compiler`: Compiler name (e.g., "gc")
- `Platform`: OS/Arch combination (e.g., "linux/amd64")
- `OS`: Operating system (e.g., "linux", "darwin", "windows")
- `Arch`: Architecture (e.g., "amd64", "arm64")
- `NumCPU`: Number of logical CPUs
- `VCSRevision`: Git commit hash (if available)
- `VCSModified`: Whether VCS tree had uncommitted changes
- `VCSTime`: Commit timestamp

### Version Checking

#### `CompareGoVersion(targetVersion string) (int, error)`

Compare current Go version with a target version. Returns:
- `-1` if current < target
- `0` if current == target
- `1` if current > target

```go
cmp, err := release.CompareGoVersion("1.20")
if err != nil {
    log.Fatal(err)
}
if cmp >= 0 {
    fmt.Println("Go 1.20 or newer")
}
```

#### `IsGoVersionAtLeast(minVersion string) (bool, error)`

Check if current Go version meets minimum requirement:

```go
if ok, _ := release.IsGoVersionAtLeast("1.21"); ok {
    // Use Go 1.21+ features
}
```

#### `GetGoMajorMinor() (major, minor int, err error)`

Extract major and minor version numbers:

```go
major, minor, _ := release.GetGoMajorMinor()
fmt.Printf("Go %d.%d\n", major, minor)
```

### Platform Detection

#### `IsPlatform(os, arch string) bool`

Check if running on specific OS and architecture:

```go
if release.IsPlatform("linux", "amd64") {
    // Linux AMD64 specific code
}
```

#### `IsOS(os string) bool`

Check operating system:

```go
if release.IsOS("linux") {
    // Linux-specific code
}
```

Common OS values: `linux`, `darwin`, `windows`, `freebsd`, `openbsd`, `netbsd`

#### `IsArch(arch string) bool`

Check architecture:

```go
if release.IsArch("arm64") {
    // ARM64-specific code
}
```

Common arch values: `amd64`, `arm64`, `386`, `arm`

### Condition Testing

Create and test custom release conditions:

```go
cs := release.NewConditionSet()

// Add conditions
cs.Add("Go Version", "Check minimum Go version", func() (bool, error) {
    return release.IsGoVersionAtLeast("1.20")
})

cs.Add("Platform", "Check supported platform", func() (bool, error) {
    return release.IsOS("linux") || release.IsOS("darwin"), nil
})

cs.Add("VCS Info", "Check VCS information is present", func() (bool, error) {
    return release.HasVCSInfo(), nil
})

// Test all conditions
results := cs.TestAll()

for _, result := range results {
    fmt.Printf("%s: %v\n", result.Name, result.Passed)
    if result.Error != nil {
        fmt.Printf("  Error: %v\n", result.Error)
    }
}

if results.AllPassed() {
    fmt.Println("Ready for release!")
}
```

### VCS Information

#### `HasVCSInfo() bool`

Check if VCS (version control) information is embedded in the binary:

```go
if release.HasVCSInfo() {
    info := release.GetBuildInfo()
    fmt.Printf("Built from commit: %s\n", info.VCSRevision)
}
```

## Use Cases

### 1. Version-Dependent Features

```go
major, minor, _ := release.GetGoMajorMinor()
if major > 1 || (major == 1 && minor >= 21) {
    // Use Go 1.21+ features
    useNewFeature()
} else {
    // Fallback for older versions
    useLegacyFeature()
}
```

### 2. Platform-Specific Code

```go
func setupPlatform() {
    if release.IsPlatform("linux", "amd64") {
        setupLinuxAMD64()
    } else if release.IsPlatform("darwin", "arm64") {
        setupMacOSArm64()
    } else {
        setupGeneric()
    }
}
```

### 3. Release Validation

```go
func validateRelease() error {
    cs := release.NewConditionSet()
    
    cs.Add("Go Version", "Go 1.20+", func() (bool, error) {
        return release.IsGoVersionAtLeast("1.20")
    })
    
    cs.Add("VCS Info", "Build has VCS metadata", func() (bool, error) {
        return release.HasVCSInfo(), nil
    })
    
    cs.Add("Production Platform", "Linux or macOS only", func() (bool, error) {
        return release.IsOS("linux") || release.IsOS("darwin"), nil
    })
    
    results := cs.TestAll()
    if !results.AllPassed() {
        return fmt.Errorf("release validation failed")
    }
    
    return nil
}
```

### 4. Build Information Display

```go
func printBuildInfo() {
    info := release.GetBuildInfo()
    fmt.Println("Build Information:")
    fmt.Printf("  Version: %s\n", info.GoVersion)
    fmt.Printf("  Platform: %s\n", info.Platform)
    fmt.Printf("  Compiler: %s\n", info.Compiler)
    
    if info.VCSRevision != "" {
        fmt.Printf("  Commit: %s\n", info.VCSRevision)
        fmt.Printf("  Modified: %v\n", info.VCSModified)
        fmt.Printf("  Time: %s\n", info.VCSTime)
    }
}
```

## Testing

Run the tests:

```bash
go test -v
```

Run benchmarks:

```bash
go test -bench=.
```

View examples:

```bash
go test -v -run=Example
```

## Building with VCS Information

To include VCS information in your builds, use Go 1.18+ and build with:

```bash
go build -buildvcs=true
```

Or ensure you're building from within a Git repository with Go 1.18+, as VCS stamping is automatic.

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.
