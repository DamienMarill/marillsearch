# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GrafiSearch is a personal search engine built with a hybrid Go backend and Preact frontend architecture. The application serves as a local web server that aggregates search results from multiple providers (Google, DuckDuckGo, Brave) while providing instant answers for common queries.

## Architecture

### Backend (Go)
- **Entry Point**: `main.go` - HTTP server running on port 8042
- **Search Providers**: `search/` directory contains implementations for different search engines
- **API Layer**: `server/api/` handles REST endpoints and search result parsing
- **Templates**: Uses `templ` for server-side rendering with Go templates in `templates/`
- **Utilities**: `utils/` contains helper functions for file operations, HTTP requests, and system installation

### Frontend (Preact + TypeScript)
- **Entry Point**: `assets/main.tsx` - Preact application with routing
- **Components**: Modular components in `assets/components/` following atomic design
- **Pages**: Route components in `assets/pages/`
- **Hooks**: Custom hooks in `assets/hooks/` for search functionality and navigation
- **Functions**: Business logic and utilities in `assets/functions/`

### Build System
- **Frontend**: Vite with Preact preset, builds to `public/assets/`
- **Backend**: Standard Go build with embedded assets
- **Templates**: `templ` generates Go code from `.templ` files

## Development Commands

### Start Development Server
```bash
make dev
```
This runs three parallel processes:
- Go server with hot reload (gow)
- Vite dev server for frontend assets
- Templ file watcher for template generation

### Build for Production
```bash
make grafisearch
```
Builds frontend assets, generates templates, and compiles Go binary

### Install System Service
```bash
make install
```
Builds the application and installs it as a system service

### Individual Commands
```bash
# Frontend development
bun run dev          # Start Vite dev server
bun run build        # Build frontend assets
bun run test         # Run Vitest tests
bun run check        # TypeScript type checking

# Go development
go run .             # Run Go server directly
go test ./...        # Run Go tests
templ generate       # Generate Go code from templates
```

## Key Architectural Patterns

### Hybrid Rendering
- Server-side templates (`templ`) for initial page structure
- Client-side Preact SPA for interactive search interface
- Asset injection via Go context for Vite-built resources

### Search Provider Abstraction
- Common `SearchResult` interface in `search/search_result.go`
- Individual provider implementations (google.go, ddg.go, brave.go)
- Unified API endpoints in `server/api/search_api.go`

### Middleware System
- Bang commands (!rt, !wrenfr) processed in `assets/functions/middleware/bangs.ts`
- Timer functionality for instant answers
- Result filtering and blocklist support

### SSE Communication
- Server-Sent Events for real-time search suggestions
- Weather API integration via dedicated endpoint

## Development Notes

### Template System
- `.templ` files are compiled to Go code - always run `templ generate` after changes
- Templates have access to Go context for dynamic asset injection
- Layout template handles Vite asset integration

### Asset Management
- Vite builds assets with content hashing for cache busting
- Go embeds `public/` directory at build time
- Development uses proxy for hot module replacement

### Testing
- Frontend: Vitest for unit tests
- Backend: Standard Go testing with `*_test.go` files
- String utilities have comprehensive test coverage

### Dependencies
- **Go**: Requires Go 1.23+, uses goquery for HTML parsing, templ for templates
- **Frontend**: Preact with preact-iso for routing, Vite for build tooling
- **Runtime**: Bun preferred for package management and frontend tooling