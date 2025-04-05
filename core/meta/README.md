# Meta Package

The meta package provides metadata and documentation utilities for the core packages.

## Structure

- Effect metadata
- Parameter validation
- Documentation utilities

## Usage

```go
// Example: Creating effect metadata
meta := meta.NewEffect(
    "effect_name",
    "Description of the effect",
    []meta.Param{
        meta.NewParam("param1", "Description of param1", 0.0, 1.0),
        meta.NewParam("param2", "Description of param2", -1.0, 1.0),
    },
)
```

## Features

- Effect metadata:
  - Name
  - Description
  - Parameters
  - Validation rules
- Parameter validation:
  - Range checking
  - Type validation
  - Default values
- Documentation:
  - Markdown generation
  - API documentation
  - Usage examples

## Implementation Details

The package:
1. Provides a consistent interface for metadata
2. Validates parameters against defined rules
3. Generates documentation from metadata
4. Supports extensible metadata types

## Testing

Run tests with:
```bash
go test ./...
``` 