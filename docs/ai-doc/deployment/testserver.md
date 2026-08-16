# Deployment: Testserver

## Purpose

Run optional HTTP lab harness. No Docker/Compose/Nginx in repo.

## Startup sequence

1. `loadDotEnv(".env")`
2. Parse port, JWT lifetime, app name
3. `auth.LoadUsersFile` → `NewStore`
4. `EventRegisterDefaultImpl` + background `EventRunDefaultImpl` loop
5. Register mux routes + Swagger
6. If TLS cert+key set → `ListenAndServeTLS`; else plaintext listen

## Run

```bash
cd cmd/testserver
cp .env_dummy .env   # set JWT key
cp users.example.yaml users.yaml
go run .
```

Swagger: `http://localhost:8000/swagger/index.html`

## OpenAPI regen

From repo root: `make swagger` (swag v1.16.4; `-d cmd/testserver,cmd/testserver/handlers`).

## Reverse proxy / Cloudflare / Apache

**N/A** in-repo. Operator-owned if exposed.

## Related Docs

- [configs/environment.md](../configs/environment.md)
- [features/testserver.md](../features/testserver.md)
- [deployment/ci.md](ci.md)
- [INDEX.md](../INDEX.md)
