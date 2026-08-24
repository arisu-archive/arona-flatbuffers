# Arona FlatBuffers

[![Python Package](https://github.com/arisu-archive/arona-flatbuffers/actions/workflows/python.yml/badge.svg)](https://github.com/arisu-archive/arona-flatbuffers/actions/workflows/python.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/arisu-archive/arona-flatbuffers.svg)](https://pkg.go.dev/github.com/arisu-archive/arona-flatbuffers)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Generated FlatBuffers schemas and ready-to-use Go and Python bindings for Arona game data.

## Contents

- [Features](#features)
- [Installation](#installation)
- [Python quickstart](#python-quickstart)
- [Development](#development)
- [Regenerating schemas](#regenerating-schemas)
- [Versioning](#versioning)
- [Contributing](#contributing)
- [License](#license)

## Features

- FlatBuffers schema definitions for Arona data
- Pre-generated Go bindings in the `flatdata` and `excel` packages
- Pre-generated Python bindings in the `FlatData` and `MX.Data.Excel` packages
- Python object APIs with byte serialization helpers
- Automated schema updates and releases

## Installation

### Python with uv

```bash
uv add git+https://github.com/arisu-archive/arona-flatbuffers --tag v0.11.0
```

Use `v0.11.0` or a newer release tag. The older `v0.10.0` tag predates the
Python package metadata and cannot be installed as a package. uv resolves the
tag to an immutable commit in your lockfile.

### Go

```bash
go get github.com/arisu-archive/arona-flatbuffers
```

Import the generated Go bindings from:

```go
import "github.com/arisu-archive/arona-flatbuffers/go/flatdata"
```

## Python quickstart

```python
from FlatData.BlendInfo import BlendInfoT

original = BlendInfoT(from_=11, to=29, blend=0.5)
payload = original.to_bytes()
restored = BlendInfoT.from_bytes(payload)

assert restored.from_ == 11
assert restored.to == 29
assert restored.blend == 0.5
```

Excel bindings use the generated `MX.Data.Excel` namespace:

```python
from MX.Data.Excel.WorldRaidStageRewardExcel import WorldRaidStageRewardExcelT
```

## Development

Requirements:

- [uv](https://docs.astral.sh/uv/)
- [Go](https://go.dev/dl/) 1.22 or newer for Go tooling and regeneration
- Bash for the generation scripts

Set up the Python project and run its tests:

```bash
uv sync --locked
uv run python -m unittest discover -s tests/python -v
uv build --no-sources
```

Run the Go tests separately:

```bash
go test ./...
```

## Regenerating schemas

The generated schema and language bindings live in:

```text
.schema/   FlatBuffers schemas
go/        Go bindings
python/    Python bindings
```

The automation in `.github/workflows/generate-schema.yml` downloads the latest
supported game data, regenerates the schemas, compiles both language bindings,
and opens a reviewable pull request. Some generator dependencies are private,
so local regeneration requires the corresponding repository access.

Generated source files should not be edited manually.

## Versioning

The project has two intentionally separate versions:

- `pyproject.toml` and `.release-please-manifest.json` contain the package and
  GitHub release version.
- `version.txt` contains the APK/schema source version used by generation
  automation.

Release Please updates the package version and creates GitHub release tags.
Python consumers install from those tags. Schema-update automation updates the
APK/schema version.

## Contributing

1. Fork the repository and create a focused branch.
2. Make the change without manually editing generated files.
3. Run the relevant Python and Go tests.
4. Open a pull request describing the behavior and verification performed.

## License

Arona FlatBuffers is available under the [MIT License](LICENSE).
