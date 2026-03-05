# Contributing

Thanks for your interest in contributing to namecheap-dns-exporter!

## Getting started

1. Fork and clone the repository
2. Install Go 1.26+
3. Run `go build ./...` to verify the build

## Development

### Project structure

```
cmd/                        # CLI entrypoint
internal/namecheap/
  types.go                  # JSON deserialization types
  exporter.go               # Zone file and Route 53 export logic
```

### Building

```bash
go build -o namecheap-dns-exporter ./cmd
```

### Running tests

```bash
go test ./...
```

## Making changes

1. Create a feature branch from `main`
2. Make your changes
3. Ensure `go build ./...` and `go test ./...` pass
4. Submit a pull request

## Guidelines

- Keep changes focused and minimal
- Follow existing code style and conventions
- Add tests for new functionality
- Update the README if adding new flags or features

## Adding a new output format

1. Add a new `Format` constant in `exporter.go`
2. Implement an `exportXxx(w io.Writer, records []Record, domain string) error` function
3. Add a case to the `switch` in `Export()`
4. Update the `--format` flag help text in `cmd/main.go`
5. Document the format in the README

## Adding a new record type

1. Add the constant to `RecordType` in `types.go`
2. Add the integer mapping in `recordTypeFromInt()`
3. Handle the type in `formatZoneRecord()` and `formatRoute53Value()`
4. Update the supported record types table in the README

## Reporting issues

Open an issue on GitHub with:
- The command you ran
- Expected vs actual output
- A sanitized sample of your input JSON (if applicable)
