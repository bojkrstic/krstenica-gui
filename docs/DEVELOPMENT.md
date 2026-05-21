# Development

## Preduslovi

- Go `1.22.x`
- Docker i Docker Compose za lokalni PostgreSQL
- PostgreSQL dostupan na `127.0.0.1:5560` kada se pokrecu integracioni/full testovi

## Lokalna baza

Pokretanje baze:

```bash
docker compose up -d krstenica_database
```

Konekcija:

```bash
docker exec -it krstenica-db sh
psql -U admin krstenica
```

## Konfiguracija

Osnovni fajl je `config/config.yaml`. Lokalna podesavanja idu u `config/config.local.yaml`, koji je ignorisan u gitu.

Preporuceni lokalni override:

```yaml
env: "local"

db:
  local_url: postgresql://admin:secret@localhost:5560/krstenica?sslmode=disable
```

`config.Load()` merge-uje lokalni fajl preko osnovnog i zatim cita environment promenljive sa prefiksom `KRSTENICA_`, npr. `KRSTENICA_HTTP_PORT`.

## Pokretanje aplikacije

Iz root-a projekta:

```bash
go run ./cmd/krstenica
```

GUI je na:

```text
http://localhost:8011/ui
```

API login je:

```text
POST /api/v1/auth/login
```

## Build

```bash
go build -o krstenica-api ./cmd/krstenica
```

Linux build preko Makefile-a:

```bash
cd cmd/krstenica
make build-linux
```

## Testiranje

Targeted testovi su bolji prvi korak:

```bash
go test ./internal/handler ./internal/service ./internal/config
```

Pun test:

```bash
go test ./...
```

Poznato ogranicenje: `go test ./...` moze zahtevati lokalni PostgreSQL na `127.0.0.1:5560`, jer `cmd/krstenica` startuje migracije.

Ako sandbox blokira Go build cache, koristiti dozvoljeni writable cache ili traziti elevaciju za test komandu.

## Stil izmena

- Pokrenuti `gofmt` posle izmena Go fajlova.
- Cuvati postojece Gin template/HTMX obrasce.
- Ne prebacivati poslovna pravila u handler ako pripadaju domenu.
- DTO promene moraju biti eksplicitne i kompatibilne sa postojecim JSON/form tagovima.
- Ne menjati deployment config usput ako zadatak nije deployment.
