package main

import (
    "bufio"
    "errors"
    "flag"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "time"

    "debug/macho"
)

// detectArchs inspects a Mach-O binary and reports whether it contains Intel (x86_64) and/or ARM64 slices.
func detectArchs(binPath string) (hasIntel bool, hasArm bool, _ error) {
    f, err := os.Open(binPath)
    if err != nil {
        return false, false, err
    }
    defer f.Close()

    // First try to open as Fat/Universal binary
    if ff, err := macho.NewFatFile(f); err == nil {
        for _, a := range ff.Arches {
            switch a.Cpu {
            case macho.CpuAmd64:
                hasIntel = true
            case macho.CpuArm64:
                hasArm = true
            }
        }
        _ = ff.Close()
        return hasIntel, hasArm, nil
    }

    // Reset and try single-arch
    if _, err := f.Seek(0, io.SeekStart); err != nil {
        return false, false, err
    }
    if mf, err := macho.NewFile(f); err == nil {
        switch mf.Cpu {
        case macho.CpuAmd64:
            hasIntel = true
        case macho.CpuArm64:
            hasArm = true
        }
        _ = mf.Close()
        return hasIntel, hasArm, nil
    } else {
        return false, false, err
    }
}

// getBundleExecutable reads Contents/Info.plist to find CFBundleExecutable.
// It implements a minimal XML scan to avoid external dependencies.
func getBundleExecutable(appPath string) (string, error) {
    plist := filepath.Join(appPath, "Contents", "Info.plist")
    b, err := os.ReadFile(plist)
    if err != nil {
        return "", err
    }
    s := string(b)
    // Look for <key>CFBundleExecutable</key> then the next <string>value</string>
    keyTag := "<key>CFBundleExecutable</key>"
    i := strings.Index(s, keyTag)
    if i == -1 {
        return "", errors.New("CFBundleExecutable not found in Info.plist")
    }
    rest := s[i+len(keyTag):]
    // Find <string>...</string>
    open := strings.Index(rest, "<string>")
    close := strings.Index(rest, "</string>")
    if open == -1 || close == -1 || close <= open+len("<string>") {
        return "", errors.New("CFBundleExecutable string not found in Info.plist")
    }
    val := rest[open+len("<string>") : close]
    val = strings.TrimSpace(val)
    if val == "" {
        return "", errors.New("empty CFBundleExecutable in Info.plist")
    }
    return val, nil
}

func isAppBundle(name string) bool {
    return strings.HasSuffix(strings.ToLower(name), ".app")
}

func main() {
    var (
        dir     string
        execute bool
        verbose bool
        dest    string
        intelOnly bool
    )

    flag.StringVar(&dir, "dir", "/Applications", "Applications directory to scan (default: /Applications)")
    flag.BoolVar(&execute, "execute", false, "Perform moves (otherwise dry-run)")
    flag.BoolVar(&verbose, "v", false, "Verbose logging")
    flag.StringVar(&dest, "dest", "", "Destination folder for Intel apps (default: <dir>/intel-apps)")
    flag.BoolVar(&intelOnly, "intel-only", true, "Move only Intel-only apps (no ARM64 slice); if false, moves any app containing Intel slice")
    flag.Parse()

    // Normalize legacy path: some users might specify /Application
    if dir == "/Application" {
        dir = "/Applications"
    }
    if dest == "" {
        dest = filepath.Join(dir, "intel-apps")
    }

    start := time.Now()
    fmt.Printf("Scanning %s for Intel architecture apps...\n", dir)

    entries, err := os.ReadDir(dir)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", dir, err)
        os.Exit(2)
    }

    toMove := make([]string, 0)
    for _, e := range entries {
        if !e.IsDir() || !isAppBundle(e.Name()) {
            continue
        }
        appPath := filepath.Join(dir, e.Name())
        exeName, err := getBundleExecutable(appPath)
        if err != nil {
            if verbose {
                fmt.Fprintf(os.Stderr, "Skip %s: %v\n", appPath, err)
            }
            continue
        }
        exePath := filepath.Join(appPath, "Contents", "MacOS", exeName)
        if _, err := os.Stat(exePath); err != nil {
            if verbose {
                fmt.Fprintf(os.Stderr, "Skip %s: missing executable %s\n", appPath, exePath)
            }
            continue
        }
        hasIntel, hasArm, err := detectArchs(exePath)
        if err != nil {
            if verbose {
                fmt.Fprintf(os.Stderr, "Skip %s: %v\n", appPath, err)
            }
            continue
        }
        if !hasIntel {
            if verbose {
                fmt.Printf("ARM-only: %s\n", e.Name())
            }
            continue
        }
        if intelOnly && hasArm {
            if verbose {
                fmt.Printf("Universal (ARM+Intel): %s\n", e.Name())
            }
            continue
        }
        toMove = append(toMove, appPath)
        if verbose {
            if intelOnly {
                fmt.Printf("Intel-only: %s\n", e.Name())
            } else {
                fmt.Printf("Contains Intel: %s\n", e.Name())
            }
        }
    }

    if len(toMove) == 0 {
        fmt.Println("No matching apps found.")
        return
    }

    fmt.Printf("Found %d app(s):\n", len(toMove))
    w := bufio.NewWriter(os.Stdout)
    for _, p := range toMove {
        fmt.Fprintln(w, " - "+filepath.Base(p))
    }
    _ = w.Flush()

    if !execute {
        fmt.Printf("Dry-run complete in %s. Re-run with -execute to move.\n", time.Since(start).Truncate(time.Millisecond))
        return
    }

    if err := os.MkdirAll(dest, 0755); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create destination %s: %v\n", dest, err)
        os.Exit(2)
    }

    moved := 0
    for _, p := range toMove {
        target := filepath.Join(dest, filepath.Base(p))
        if err := os.Rename(p, target); err != nil {
            fmt.Fprintf(os.Stderr, "Failed to move %s → %s: %v\n", p, target, err)
            continue
        }
        moved++
        if verbose {
            fmt.Printf("Moved %s → %s\n", filepath.Base(p), target)
        }
    }

    fmt.Printf("Moved %d/%d app(s) to %s in %s.\n", moved, len(toMove), dest, time.Since(start).Truncate(time.Millisecond))
}

