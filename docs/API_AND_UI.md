# API And UI

## Auth

UI koristi session cookie `krstenica_session`. API koristi bearer token dobijen preko:

```text
POST /api/v1/auth/login
```

U `local`, `dev` i `development` okruzenjima API auth i role provere se zaobilaze prema trenutnoj implementaciji. UI auth i dalje ide kroz login/session tok.

## API rute

Prefix je:

```text
api/v1/adminv2
```

Admin-only resursi:

- `tamples`
- `priests`
- `eparhije`
- `persons`
- `users`

Krstenice su dostupne autentifikovanim korisnicima, dok servis sprovodi role/city pravila:

- `POST api/v1/adminv2/krstenice`
- `GET api/v1/adminv2/krstenice`
- `GET api/v1/adminv2/krstenice/:id`
- `PUT api/v1/adminv2/krstenice/:id`
- `DELETE api/v1/adminv2/krstenice/:id`
- `GET api/v1/adminv2/krstenice-print/:id`

## Stampanje

Primer:

```text
GET /api/v1/adminv2/krstenice-print/1?filter[format][eq]=pdf
```

Korisni filteri:

- `preview=true` - koristi preview Excel sablon.
- `format=pdf` - generise PDF umesto XLSX.
- `template_version=2` - koristi varijantu bez full-bleed pozadine.
- `font=<key>` - bira font gde je podrzano u PDF toku.

## GUI rute

Glavni ulaz:

```text
/ui
```

Sekcije:

- `/ui/krstenice`
- `/ui/eparhije`
- `/ui/hramovi`
- `/ui/svestenici`
- `/ui/osobe`
- `/ui/users` za admin korisnike
- `/ui/uputstvo`

## HTMX obrazac

Za vecinu sekcija vazi isti tok:

1. `index.html` prikazuje stranicu.
2. `GET /ui/<sekcija>/table` vraca tabelu kao fragment.
3. `GET /ui/<sekcija>/new` ili `GET /ui/<sekcija>/:id/edit` vraca modal.
4. `POST`, `PUT` ili `DELETE` handler menja podatke i salje `HX-Trigger` dogadjaj.
5. Tabela osluskuje dogadjaj i osvezava se bez reload-a cele stranice.

Kod novih UI izmena pratiti postojece sekcije umesto uvodjenja drugog frontend obrasca.
