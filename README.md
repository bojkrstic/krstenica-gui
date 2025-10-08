# Krstenica GUI

Ovo je GUI i API sloj za evidenciju krstenica u okviru pravoslavne crkve.
U nastavku su objedinjena uputstva za podizanje okruzenja, debug, pristup bazi i osnovno odrzavanje.

## Pokretanje preko Docker Compose
- Podigni celokupan stack: `docker compose up -d`
- Primer API poziva za proveru: `GET http://localhost:8011/api/v1/adminv2/tamples/1`

## Lokalno pokretanje Go servisa
- Napravi binarni fajl: `go build -o krstenica-api`
- Pokreni binarno izdanje: `./krstenica-api`
- Alternativa za debug: `go run main.go`

## Web GUI (HTMX)
- Poseti `http://localhost:8011/ui` za dashboard i listu krstenica.
- Stranica koristi HTMX pa se podaci dinamicki ucitavaju iz `api/v1/adminv2/krstenice` endpoint-a.
- Pretragu po imenu pokrecemo direktno sa stranice; paginacija radi kroz HTMX bez reload-a.
- U koloni "Akcije" dostupno je dugme `Stampaj` koje generise Excel krstenicu sa pozadinskim obrascem ( `krstenica_obrada.jpg` ).
- Fajl `krstenica_obrada.jpg` treba da stoji u korenu repozitorijuma kako bi pozadina bila podvučena ispod popunjenih polja prilikom štampe.

## Rad sa PostgreSQL bazom u kontejneru
```
docker exec -it krstenica_db sh
psql -U admin krstenica
\dt
select * from public.eparhije;
```

## Gasenje procesa na portu 8011
```
lsof -i :8011
kill -9 <PID>
```

## Build za produkciono okruzenje
- Kreiraj Linux binarno izdanje: `make build-linux`

## Debug workflow
1. Uveri se da port 8011 nije zauzet (`lsof -i :8011`, pa `kill -9 <PID>` ako je potrebno).
2. Pokreni servis (`./krstenica-api` ili `go run main.go`).
3. Pogadjaj API iz Postmana / curl-a, npr. `GET http://localhost:8011/api/v1/adminv2/tamples/1`.
4. Kroz IDE kontrolisi breakpointe i tok izvrsavanja.

## Resavanje problema sa migracijama
Ako kreiranje baze zataji zbog tabele `schema_migrations`, proveri polje `dirty` i postavi ga na `false`:
```
update schema_migrations set dirty=false;
```

## Git napomena
Ovaj direktorijum je povezan sa Git remote-om `git@github.com:bojkrstic/krstenica-gui.git`.
Za promenu origin-a koristi sledece:
```
git remote remove origin
git remote add origin <novi_git_url>
git branch -M main
git push -u origin main
```
