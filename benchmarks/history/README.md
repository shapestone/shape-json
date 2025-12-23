# Benchmark History

This directory contains historical benchmark results for tracking performance over time.

## Overview

Each benchmark run creates a timestamped directory containing:
- Raw benchmark output
- Generated performance report
- Metadata (git commit, platform info, etc.)

## Directory Structure

```
benchmarks/history/
├── README.md                    # This file
├── .gitignore                   # Excludes large benchmark files from git
└── YYYY-MM-DD_HH-MM-SS/        # Timestamped benchmark runs
    ├── benchmark_output.txt     # Raw Go benchmark results
    ├── PERFORMANCE_REPORT.md    # Generated performance report
    └── metadata.json            # Run metadata
```

## Metadata Format

Each benchmark run includes a `metadata.json` file with the following information:

```json
{
  "timestamp": "2025-12-21_16-30-45",
  "commit": "ee2758df220d2752d4adc29522ab6dbb2c5a43d9",
  "platform": "Apple M1 Pro",
  "os": "darwin",
  "arch": "arm64",
  "go_version": "1.25.0",
  "bench_time": "3s",
  "description": "Baseline before optimization"
}
```

### Metadata Fields

- **timestamp**: Directory name and run time (YYYY-MM-DD_HH-MM-SS format)
- **commit**: Full git commit hash at time of benchmark
- **platform**: CPU/platform name (e.g., "Apple M1 Pro", "Intel Core i7")
- **os**: Operating system (e.g., "darwin", "linux", "windows")
- **arch**: CPU architecture (e.g., "arm64", "amd64")
- **go_version**: Go version used for benchmarks
- **bench_time**: Benchmark duration setting (usually "3s")
- **description**: Optional user-provided description of the run

## Usage

### Creating Benchmark History

Benchmark history is automatically created when you run:

```bash
make performance-report
```

This will:
1. Run all benchmarks
2. Generate performance report
3. Save both to `benchmarks/history/YYYY-MM-DD_HH-MM-SS/`
4. Create `.gitignore` to exclude files from version control

### Viewing History

List all available benchmark runs:

```bash
make bench-history
```

Example output:
```
Available benchmark history:

  2025-12-21_16-30-45
    "commit": "ee2758d",
    "platform": "Apple M1 Pro",

  2025-12-20_14-22-10
    "commit": "a1b2c3d",
    "platform": "Apple M1 Pro",

  2025-12-19_09-15-30
    "commit": "f4e5d6c",
    "platform": "Apple M1 Pro",
```

### Comparing Runs

Compare the latest run vs the previous run:

```bash
make bench-compare-history
```

Or compare specific runs using the comparison tool:

```bash
# Compare latest vs previous
go run scripts/compare_benchmarks/main.go latest previous

# Compare specific timestamps
go run scripts/compare_benchmarks/main.go 2025-12-20_14-22-10 2025-12-21_16-30-45

# Compare using full paths
go run scripts/compare_benchmarks/main.go \
  benchmarks/history/2025-12-20_14-22-10/benchmark_output.txt \
  benchmarks/history/2025-12-21_16-30-45/benchmark_output.txt
```

### Statistical Comparison with benchstat

The comparison tool uses benchstat (if available) to show:

```
Benchmark Comparison
====================

Old: benchmarks/history/2025-12-20_14-22-10
  Commit: a1b2c3d
  Platform: Apple M1 Pro (darwin/arm64)
  Go: 1.25.0

New: benchmarks/history/2025-12-21_16-30-45
  Commit: ee2758d
  Platform: Apple M1 Pro (darwin/arm64)
  Go: 1.25.0
  Description: After parser optimization

Statistical Comparison (benchstat)
==================================

name                              old time/op    new time/op    delta
ShapeJSON_Parse_Small-10            1.23µs ± 2%    1.15µs ± 1%   -6.50%  (p=0.000 n=10+10)
ShapeJSON_Parse_Medium-10           78.5µs ± 1%    72.3µs ± 2%   -7.90%  (p=0.000 n=10+10)
ShapeJSON_Parse_Large-10            6.89ms ± 3%    6.12ms ± 2%  -11.18%  (p=0.000 n=10+10)

name                              old alloc/op   new alloc/op   delta
ShapeJSON_Parse_Small-10            2.45kB ± 0%    2.35kB ± 0%   -4.08%  (p=0.000 n=10+10)
ShapeJSON_Parse_Medium-10            156kB ± 0%     145kB ± 0%   -7.05%  (p=0.000 n=10+10)
ShapeJSON_Parse_Large-10            14.2MB ± 0%    12.8MB ± 0%   -9.86%  (p=0.000 n=10+10)
```

Interpretation:
- **~** means no significant change (within statistical variance)
- **+** means new is slower (regression) or uses more memory
- **-** means new is faster (improvement) or uses less memory
- **p-value** < 0.05 indicates statistically significant change
- **±** shows variance across multiple runs

### Installing benchstat

If you don't have benchstat installed:

```bash
go install golang.org/x/perf/cmd/benchstat@latest
```

Without benchstat, the comparison tool shows a simple side-by-side view.

## Adding Descriptions

Track the purpose of each benchmark run:

```bash
# Run with description
go run scripts/generate_benchmark_report/main.go \
  -description "Baseline before string interning optimization"

# The description will appear in metadata.json and comparisons
```

This is useful for:
- Tracking optimization attempts
- Before/after feature changes
- Platform comparisons
- Go version upgrades

## Workflow Examples

### Tracking an Optimization

```bash
# 1. Baseline benchmark
go run scripts/generate_benchmark_report/main.go \
  -description "Baseline before parser optimization"

# 2. Make your code changes
# ... edit code ...

# 3. Run new benchmark
make performance-report

# 4. Compare results
make bench-compare-history
```

### Comparing Across Branches

```bash
# On main branch
git checkout main
make performance-report

# On feature branch
git checkout feature/parser-optimization
make performance-report

# Compare the two runs
go run scripts/compare_benchmarks/main.go \
  2025-12-21_10-00-00 \  # main branch timestamp
  2025-12-21_10-30-00    # feature branch timestamp
```

### Tracking Performance Over Time

View historical performance trends:

```bash
# List all runs
make bench-history

# Compare milestones
go run scripts/compare_benchmarks/main.go \
  2025-12-01_09-00-00 \  # v1.0.0 release
  2025-12-20_14-00-00    # v1.1.0 release
```

## Git Integration

The `.gitignore` file in this directory excludes benchmark files from version control:

```gitignore
# Benchmark history files are large and change frequently
# Only commit the directory structure and README
*
!.gitignore
!README.md
```

This means:
- ✅ Directory structure is tracked
- ✅ This README is tracked
- ❌ Benchmark data is NOT tracked (local only)

### Why exclude from git?

1. **Size**: Benchmark files can be large (100+ KB per run)
2. **Frequency**: Frequent benchmark runs create many files
3. **Machine-specific**: Results are platform/hardware-specific
4. **Reproducibility**: Can be regenerated with `make performance-report`

If you need to share benchmark results:
- Share specific `PERFORMANCE_REPORT.md` files via other means
- Include benchmark summaries in commit messages or PRs
- Use CI/CD to run benchmarks on standardized hardware

## CI/CD Integration

For consistent benchmark tracking across your team:

```yaml
# .github/workflows/benchmark.yml
name: Benchmark

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  benchmark:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Run benchmarks
        run: make performance-report

      - name: Upload benchmark results
        uses: actions/upload-artifact@v4
        with:
          name: benchmark-results
          path: PERFORMANCE_REPORT.md

      - name: Compare with main
        if: github.event_name == 'pull_request'
        run: |
          # Download main branch benchmarks
          # Compare current vs main
          # Post results to PR
```

## Maintenance

### Cleaning Old History

Benchmark history can accumulate over time. To clean up:

```bash
# Remove all history
rm -rf benchmarks/history/2025-*

# Keep only last 10 runs
cd benchmarks/history
ls -t | tail -n +11 | xargs rm -rf

# Keep only runs from this month
cd benchmarks/history
ls -d 2025-11-* | xargs rm -rf  # Remove November runs
```

### Archiving Important Runs

To preserve specific benchmark runs:

```bash
# Archive a specific run
tar -czf baseline-v1.0.tar.gz benchmarks/history/2025-12-01_09-00-00/

# Archive all history
tar -czf all-benchmarks.tar.gz benchmarks/history/
```

## Best Practices

1. **Run benchmarks on consistent hardware**
   - Same machine for meaningful comparisons
   - Minimal background processes
   - Consistent power/thermal state

2. **Add descriptions for important runs**
   - Track optimization goals
   - Note what changed
   - Reference issue/PR numbers

3. **Use benchstat for accuracy**
   - Install with: `go install golang.org/x/perf/cmd/benchstat@latest`
   - Shows statistical significance
   - Filters out noise

4. **Regular cleanup**
   - Archive old runs
   - Keep 1-2 months of history
   - Document major milestones

5. **Track with git commits**
   - Benchmark after significant changes
   - Reference commits in descriptions
   - Compare across releases

## Troubleshooting

### No benchmark history found

```bash
# Create your first benchmark
make performance-report
```

### benchstat not found

```bash
# Install benchstat
go install golang.org/x/perf/cmd/benchstat@latest

# Verify installation
benchstat -h
```

### Comparison shows unexpected results

- Check metadata files to verify:
  - Different git commits
  - Same platform/OS/arch
  - Same Go version
- Ensure consistent benchmark conditions:
  - Same machine
  - Minimal background processes
  - Multiple runs for statistical significance

### Cannot find specific timestamp

```bash
# List all available runs
make bench-history

# Use exact timestamp format: YYYY-MM-DD_HH-MM-SS
go run scripts/compare_benchmarks/main.go 2025-12-21_16-30-45 latest
```

## Additional Resources

- [Go Benchmarking Documentation](https://pkg.go.dev/testing#hdr-Benchmarks)
- [benchstat Documentation](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)
- [PERFORMANCE_REPORT.md](../../PERFORMANCE_REPORT.md) - Latest benchmark results
- [README.md](../../README.md#performance-benchmarking) - Benchmark system overview
