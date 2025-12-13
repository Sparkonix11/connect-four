#!/bin/bash
# Generate types from OpenAPI specification
# Run this after editing api/openapi.yaml

set -e

echo "Generating Go types..."
oapi-codegen -package generated -generate types -o internal/api/generated/types.gen.go api/openapi.yaml

echo "Generating TypeScript types..."
cd frontend && npm run generate

echo "✅ Types generated successfully!"
echo ""
echo "Files updated:"
echo "  - internal/api/generated/types.gen.go"
echo "  - frontend/src/types/api.gen.ts"
