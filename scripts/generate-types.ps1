# Generate types from OpenAPI specification
# Run this after editing api/openapi.yaml

Write-Host "Generating Go types..." -ForegroundColor Cyan
oapi-codegen -package generated -generate types -o internal/api/generated/types.gen.go api/openapi.yaml

Write-Host "Generating TypeScript types..." -ForegroundColor Cyan
Push-Location frontend
npm run generate
Pop-Location

Write-Host ""
Write-Host "Types generated successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Files updated:"
Write-Host "  - internal/api/generated/types.gen.go"
Write-Host "  - frontend/src/types/api.gen.ts"
