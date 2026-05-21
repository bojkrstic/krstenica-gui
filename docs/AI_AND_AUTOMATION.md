# AI And Automation Guide

Ovaj fajl je namenjen AI agentima i buducim pomocnim skriptama. Cilj je da se promene rade predvidljivo, bez gubljenja konteksta i bez gazenja korisnickih izmena.

## Obavezna pravila za AI agente

- Raditi iz root-a projekta: `/home/bojan/develop/horisen/Krstenica/Krstenica-gui/krstenica`.
- Pre izmena pokrenuti `git status --short`.
- Ne revertovati i ne prepisivati tudje izmene osim ako korisnik to eksplicitno trazi.
- Drzati izmene fokusirane na trazeni zadatak.
- Pre kraja rada obavezno azurirati `SESSION.md`.
- U `SESSION.md` upisati sta je zavrseno, koji fajlovi su menjani, sta je provereno, sta nije provereno i zasto.
- Ne upisivati tajne u dokumentaciju, komentare, testove ili skripte.

## Kada se dodaje nova funkcionalnost

1. Procitati postojece handler/service/repository obrasce za isti tip resursa.
2. Dodati DTO polja eksplicitno sa JSON/form tagovima.
3. Poslovna pravila staviti u `internal/service`.
4. DB pristup staviti u `internal/repository`.
5. GUI izmenu uklopiti u postojece `index/table/new/edit` HTMX sablone.
6. Ako se menja schema, dodati novu numerisanu migraciju u `migrations/postgres`.
7. Pokrenuti `gofmt` za Go fajlove.
8. Pokrenuti najuzze relevantne testove, pa sire testove ako uslovi postoje.

## Kada se dodaje nova skripta

Skripte treba drzati male i eksplicitne. Preporuceni direktorijum za buduce skripte je `scripts/`.

Svaka skripta treba da ima:

- jasan naziv, npr. `scripts/export-db.sh` ili `scripts/check-config.sh`;
- `set -euo pipefail` ako je Bash;
- proveru da se pokrece iz root-a projekta ili automatski prelazak u root;
- `--help` ili kratak usage blok ako prima argumente;
- dry-run opciju za operacije koje menjaju stanje;
- jasne poruke greske;
- bez hardkodovanih tajni.

Primer skeletona:

```bash
#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

usage() {
  echo "Usage: $0 [--dry-run]"
}

DRY_RUN=false
for arg in "$@"; do
  case "$arg" in
    --dry-run) DRY_RUN=true ;;
    --help|-h) usage; exit 0 ;;
    *) echo "Unknown argument: $arg" >&2; usage; exit 1 ;;
  esac
done
```

## Korisne provere za skripte

- `go test ./internal/...` za brzu proveru internih paketa.
- `go test ./...` samo kada je dostupna lokalna baza.
- `docker compose config` za proveru compose sintakse.
- `docker build -t krstenica-svc:local .` za proveru image build-a.

## Automatizacija deploymenta

Deployment skripta ne treba automatski da brise ili zamenjuje produkcioni container bez eksplicitne potvrde, osim ako korisnik bas trazi non-interactive rollout.

Za rollout skriptu koristiti promenljive:

- `IMAGE_TAG`
- `CONTAINER_NAME`
- `NETWORK_NAME`
- `CONFIG_DIR`
- `DOC_DIR`
- `HTTP_PORT`

Preporuceni redosled:

1. proveri da je `IMAGE_TAG` prosledjen;
2. proveri da je korisnik ulogovan na Docker registry;
3. build;
4. push;
5. na serveru pull;
6. proveri config;
7. zaustavi stari container;
8. startuj novi container;
9. proveri logs i `/ui`.

## Poznati rizici

- `go test ./...` moze pasti bez lokalnog PostgreSQL-a.
- Export krstenice zavisi od `doc/template_files` i `pictures`.
- Server i lokalni DSN se razlikuju; ne menjati server config na lokalni DSN u commitu.
- UI je HTMX/Go template aplikacija; uvodjenje SPA frameworka nije kompatibilno sa trenutnim stilom projekta bez posebne odluke.
