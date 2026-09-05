# Agents Guide: go-alfred

## Structure
- `api/`: Go backend (Telegram bot and API).
  - `cmd/api/`: Main entrypoint.
  - `internal/`: Business logic (chat, dinner, cron, etc.).
  - `migrations/`: SQL migrations.
- `ui/`: Nuxt frontend (Telegram Mini App).
- `supabase/`: Database and infra configuration.

## Commands

### Backend (Go)
- Run: `go run api/cmd/api/main.go` (Verify entrypoint path in `api/cmd/api/`)
- Test: `go test ./api/...`

### Frontend (Nuxt)
- Workdir: `ui/`
- Dev: `npm run dev`
- Build: `npm run build`
- Test: `npm run test`

## Key Constraints
- **Infrastructure**: Uses Supabase for DB and DigitalOcean for hosting.
- **API**: Integrates with Telegram Bot API and Telegram Mini Apps SDK.
- **Database**: Migrations are managed in `api/migrations/`.
