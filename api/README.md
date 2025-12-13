# Shared Types Generation

This project uses **OpenAPI as the single source of truth** for types, automatically generating code for both Go and TypeScript.

## Workflow

```
┌─────────────────────┐
│ api/openapi.yaml    │  ← Source of truth
└──────────┬──────────┘
           │
           ├──────────────────────────────────────┐
           │                                      │
           ▼                                      ▼
┌─────────────────────────────────┐  ┌───────────────────────────────┐
│ oapi-codegen (Go)               │  │ openapi-typescript (TS)       │
│ internal/api/generated/         │  │ frontend/src/types/api.gen.ts │
└─────────────────────────────────┘  └───────────────────────────────┘
```

## Regenerate Types

After editing `api/openapi.yaml`, run:

```powershell
# Windows
.\scripts\generate-types.ps1

# Unix/Mac
./scripts/generate-types.sh
```

Or individually:

```bash
# Go types
oapi-codegen -package generated -generate types -o internal/api/generated/types.gen.go api/openapi.yaml

# TypeScript types
cd frontend && npm run generate
```

## Important Notes

- **Never edit generated files directly** - they will be overwritten
- All REST API types are auto-generated
- WebSocket message types are defined in the OpenAPI spec and generated for both languages
- Run `npm run typecheck` and `staticcheck ./...` after regenerating
