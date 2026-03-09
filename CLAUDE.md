# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is `kratos-scaffold`, a CLI tool for generating Kratos-layout style Go microservice projects and components. The tool supports generating complete projects (mono-repo, single service), as well as individual components (proto files, business logic, data access, and service layers).

## Architecture

The project is structured as a CLI application built with Cobra:

- **Main Entry Point**: `main.go:6` calls `cmd.Execute()`
- **Command Structure**: `cmd/` contains all CLI commands (root, new, biz, data, service, proto, generate)
- **Code Generation**:
  - `generator/` contains template-based generators for individual components
  - `project_generator/` handles full project scaffolding
  - Templates are in `generator/tmpl/` and `project_generator/resources/`
- **Core Utilities**:
  - `pkg/field/` defines field types, predicates, and validation for code generation
  - `pkg/util/` provides file operations and Go tooling helpers
  - `pkg/cli/` handles environment settings and CLI configuration
  - `pkg/ast/` provides AST manipulation utilities (e.g., wire injection)
  - `internal/merge/` handles YAML configuration merging

## Development Commands

### Building and Testing
```bash
# Build the project
go build -o kratos-scaffold .

# Install locally for development
go install .

# Run tests
go test ./...

# Run specific test file
go test ./cmd -run TestBiz

# Update golden test files
go test ./cmd -update
```

### Code Generation Testing
The project uses golden file testing for verifying code generation output. Test files are in `cmd/testdata/`:
- `biz-user.txt` - Expected biz layer output
- `data-ent-user.txt` - Expected data layer with Ent output
- `data-ent-user_transfer.txt` - Expected data transfer layer
- `service-user.txt` - Expected service layer output
- `service-user_transfer.txt` - Expected service transfer layer
- `proto.txt` - Expected proto file output

## Key Components

### Field System (`pkg/field/`)
Defines a comprehensive type system for code generation with:
- **Field Types**: Basic types (int32, string, time), complex types (json, arrays)
- **Predicates**: Query conditions (eq, cont, gt, gte, lt, lte, in)
- **Nullability**: Optional fields marked with 'nil' or 'null'
- **Database Mapping**: Automatic mapping to appropriate database types
- **Field Style**: Configurable proto field naming (snake, camel, low-camel)

### Project Types (`project_generator/project.go`)
Supports three project types:
1. **Mono**: Main monorepo with `app/` directory structure
2. **SubMono**: Sub-service within existing monorepo (detected by presence of `go.mod`)
3. **Single**: Standalone single service

Detection logic in `project_generator/project.go:22-29` and `project_generator/project.go:41-60`.

### Template System
Uses Go templates for code generation:
- Business logic templates in `generator/tmpl/biz.tmpl`
- Data layer templates support both Ent ORM (`data_ent_*.tmpl`) and proto-based data access (`data_proto_*.tmpl`)
- Service templates generate gRPC service implementations
- Proto templates create protobuf definitions with proper field mappings
- Migration SQL templates in `generator/tmpl/migration.sql.tmpl`

### Wire Integration (`pkg/ast/wire.go`)
The tool automatically injects constructors into `biz.go` using AST manipulation when generating biz layer code.

## Field Type Mapping

The tool maintains consistent type mappings across layers:
- **Go Entity → Go Params → Proto Entity → Proto Params → Database**
- Example: `int64 → *int64 → int64 → optional int64 → bigint`
- Special handling for time fields, JSON types, and arrays
- Full mapping table available in README.md:22-34

## Command Structure

Commands follow a consistent pattern:
1. **new**: Creates new project (mono, sub-mono, or single)
2. **proto**: Generates proto files and compiles them with protoc
3. **biz**: Generates business logic layer
4. **data**: Generates data access layer (supports `--orm-type=ent` or `--orm-type=proto`)
5. **service**: Generates service layer
6. **generate (g)**: One-shot generation of proto, biz, data, and service

## Configuration

Configuration is loaded from (in order of precedence):
1. Environment variables (prefix: `KRATOS_`)
2. `.kratos-scaffold.yaml` in current directory
3. `~/.kratos/config`

Key settings in `pkg/cli/environment.go:24-32`:
- `AppDirName`: Directory name for mono repo apps (default: "app")
- `ApiDirName`: Directory name for API definitions (default: "api")
- `Namespace`: Target sub-service name
- `FieldStyle`: Proto field naming style (snake/camel/low-camel)
- `PrimaryKey`: Primary key field name (default: "id")

## Usage Patterns

When working with this codebase:
1. New commands should follow the pattern in `cmd/` with corresponding test files
2. Template modifications require updating golden test files (`go test ./cmd -update`)
3. Field type additions need updates in `pkg/field/field_type.go`
4. Project generation logic is centralized in `project_generator/`
5. Generators embed templates using `//go:embed` directive

## Testing Strategy

- Unit tests for individual components
- Golden file tests for code generation verification
- Test workspace creation utilities in `cmd/root_test.go:40-58`
- File operation tests in `pkg/util/file_test.go:9-22`
