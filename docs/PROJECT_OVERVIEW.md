# Project Overview

Krstenica GUI je Go aplikacija za evidenciju krstenica. Aplikacija ima web GUI baziran na Gin HTML sablonima i HTMX-u, JSON API za iste resurse, PostgreSQL bazu i eksport krstenica u Excel/PDF obrasce.

## Glavni domeni

- `krstenice` - centralna evidencija krstenja.
- `eparhije` - crkvene administrativne jedinice.
- `hramovi` / `tamples` - hramovi vezani za krstenice.
- `svestenici` / `priests` - svestenici i njihove titule.
- `osobe` / `persons` - roditelji, kumovi i druge povezane osobe.
- `users` - korisnici aplikacije, role i ogranicenje po gradu.

## Tehnologije

- Go `1.22.2`
- Gin HTTP framework
- HTMX u Go HTML sablonima
- GORM + PostgreSQL
- `golang-migrate` sa embedovanim migracijama
- Excelize za XLSX stampu
- `gofpdf` za PDF generisanje
- Viper za konfiguraciju

## Vazni direktorijumi

- `cmd/krstenica/` - entrypoint aplikacije i lokalni Makefile.
- `internal/handler/` - HTTP rute, auth middleware, GUI handleri, API handleri i stampanje.
- `internal/service/` - poslovna pravila, validacije i city/role ogranicenja.
- `internal/repository/` - GORM pristup bazi.
- `internal/model/` - DB modeli.
- `internal/dto/` - request/response strukture.
- `migrations/postgres/` - SQL migracije.
- `web/templates/` - layout, stranice, tabele, modali i HTMX fragmenti.
- `web/static/` - staticki fajlovi koje Gin servira preko `/static`.
- `doc/template_files/` - Excel sabloni za stampu.
- `pictures/` - slike i pozadine za stampu/dokumentaciju.
- `config/` - osnovna i lokalna konfiguracija.

## Runtime tok

`cmd/krstenica/main.go` ucitava konfiguraciju, pokrece migracije, inicijalizuje GORM, sklapa repository, service i HTTP handler, zatim pokrece Gin server na `http_port`.

Migracije se pokrecu pri startu aplikacije, pa lokalni testovi koji ulaze u `cmd/krstenica` mogu traziti dostupan PostgreSQL na konfigurisanom DSN-u.
