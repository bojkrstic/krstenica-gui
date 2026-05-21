# Architecture

Arhitektura je jednostavan slojeviti Go monolit:

`Gin handler -> service -> repository -> PostgreSQL`

GUI i API dele isti servisni sloj. GUI vraca HTML stranice ili HTMX fragmente, a API vraca JSON ili fajl za download.

## Start aplikacije

1. `config.Load()` ucitava `config.yaml`, opciono merge-uje `config.local.yaml`, zatim cita `KRSTENICA_*` environment promenljive.
2. `config.PostgresMigrate()` pokrece embedovane SQL migracije iz `migrations/postgres`.
3. `repository.InitORM()` otvara GORM konekciju.
4. `repository.NewRepository()` kreira DB adapter.
5. `service.NewService()` vezuje poslovna pravila za repository.
6. `handler.NewHttpHandler().Init()` registruje static assets, template funkcije, HTML sablone, auth rute, API rute i GUI rute.

## Handler sloj

`internal/handler` sadrzi:

- `handler.go` - Gin router, static path resolve, template loading i zajednicki render helper.
- `routes.go` - JSON/API rute pod `api/v1/adminv2`.
- `gui.go` - UI rute pod `/ui`, table renderi, forme i HTMX tokovi.
- `auth.go` - UI session cookie, API bearer token i role middleware.
- `krstenice-print.go` i `krstenice-pdf.go` - Excel/PDF generisanje.
- fajlove po resursima, npr. `krstenice.go`, `eparhije.go`, `persons.go`.

Handler treba da ostane tanak: bind requesta, poziv servisa, format odgovora. Poslovna pravila idu u `internal/service`.

## Service sloj

`internal/service` sadrzi validacije, mapiranje DTO-a u modele i pravila pristupa. Posebno je vazno:

- ne-admin korisnici rade samo nad krstenicama svog grada;
- ako ne-admin korisnik nema grad, servis vraca gresku;
- brisanje je uglavnom soft-delete kroz `status`, gde je to vec postojeci obrazac;
- validacije i default vrednosti treba drzati u servisu kada uticu na domen.

## Repository sloj

`internal/repository` koristi GORM. U ovom sloju treba drzati SQL/GORM detalje: filtere, joinove, count, paging i update mape.

Paging helper koristi query parametre parsirane u `pkg.FilterAndSort`. Ako se dodaju novi list endpointi, pratiti postojece `List*` obrasce.

## Template sloj

`web/templates/layouts/base.html` je zajednicki layout. Stranice definisu content blok i handler im prosledjuje `ContentTemplate`.

Tipican UI resurs ima:

- `index.html` - puna stranica;
- `table.html` - HTMX fragment za tabelu i paginaciju;
- `new.html` - modal za kreiranje;
- `edit.html` - modal za izmenu;
- `picker.html` i `picker-table.html` kada se bira povezan entitet.

## Stampanje

Stampanje koristi Excel sablone iz `doc/template_files` i slike iz `pictures`. Endpoint `GET api/v1/adminv2/krstenice-print/:id` podrzava XLSX i PDF izlaz, kao i `preview`, `format`, `template_version` i `font` filtere.

Runtime image mora imati `web/`, `pictures/` i `doc/template_files/` dostupne u `/app`, inace GUI ili export mogu otkazati.
