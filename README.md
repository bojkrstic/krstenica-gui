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
sudo lsof -i :8011
sudo kill -9 <PID>
```

## Build za produkciono okruzenje
- Kreiraj Linux binarno izdanje: `make build-linux`

## Debug workflow
1. Uveri se da port 8011 nije zauzet (`sudo lsof -i :8011`, pa `kill -9 <PID>` ako je potrebno).
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

## Glavni tok

  - Aplikacija se startuje iz cmd/krstenica/main.go:1: učitava konfiguraciju, pokreće migracije, inicijalizuje bazu i sklapa repository → service → http
  handler.
  - HTTP sloj se priprema u internal/handler/handler.go:32: kreira se Gin router, registruje servis statičkih fajlova (/static), dodaju se zajedničke
  template funkcije i učitavaju svi HTML šabloni iz web/templates.
  - Rute za GUI su definisane u internal/handler/gui.go:24; ovde svaki URL (/ui/…) mapira na funkcije koje sklapaju podatke i renderuju konkretne šablone.
  - Interakcija frontenda oslanja se na HTMX (učitan u web/templates/layouts/base.html:7): forme i dugmad šalju asinhro zahteve i dobijaju delimične HTML
  odgovore koje router vraća.

  Struktura direktorijuma

  - web/templates/layouts/base.html:1 drži glavni “layout” (header, navigacija, <main> blok); sve ostale stranice ga uključuju.
  - web/templates/<entitet>/*.html sadrži stranice za dashboard, krštenice, eparhije, hramove, svештенike и особе; svaka sekcija ima index, table, new,
  edit, eventualno picker.
  - web/templates/partials/error.html:1 je univerzalni isječak za prikaz greške.
  - web/static/ je predviđen za dodatne CSS/JS fajlove ili slike koje se služe preko /static rute (trenutno prazan).
  - internal/handler/*.go fajlovi grupisani po resursima (npr. krstenice.go, eparhije.go) obrađuju HTTP zahteve, komuniciraju sa servisnim slojem i
  prosleđuju podatke šablonima.

  Rad sa šablonima

  - Svaki index.html definše stranicu i prosleđuje koji “content” blok da se ubaci u layout (web/templates/krstenice/index.html:1).
  - table.html fajlovi renderuju tablice sa paginacijom; odgovori dolaze preko HTMX-a kako bi se osvežio samo deo stranice (web/templates/osobe/
  table.html:8).
  - new.html i edit.html su modali (HTML <dialog>) koji hvataju submit preko HTMX-a i posle uspeha šalju događaj za osvežavanje odgovarajuće tabele (web/
  templates/hramovi/new.html:1, web/templates/hramovi/edit.html:1).
  - Pickeri za izbor entiteta rade kao ugnježdeni modali: osnovni prikaz (web/templates/osobe/picker.html:1) otvara tabelu (web/templates/osobe/picker-
  table.html:8) i šalje prilagođene događaje na izbor reda.
  - Dodatne funkcije (format datuma, konverzija ID vrednosti) se dodaju preko router.SetFuncMap i dostupne su u svim šablonima (internal/handler/
  handler.go:39).

  Ako želiš da proširiš GUI:

  1. Definiši novu rutu i handler u odgovarajućem internal/handler/*.go.
  2. Dodaj nove HTML šablone u web/templates/<sekcija>/.
  3. Po potrebi ubaci stilove ili skripte u web/static/ i referenciraj ih iz layout-a.
  4. Testiraj interakcije tako što ćeš otvoriti /ui i pratiti HTMX zahteve u mrežnom panelu.
``

## Osnovni layout je Go HTML šablon definisan u web/templates/layouts/base.html:1. Evo kako se učitava i koristi:

  - При иницијализацији HTTP sloja, у internal/handler/handler.go:35 Gin позива router.LoadHTMLFiles над свим .html фајловима под web/templates. То учитава
  и региструје и layouts/base.html.
  - Layout је дефинисан као именовани шаблон: {{ define "layouts/base" }} на врху фајла. Унутар њега су <html>, <head>, <body>, навигација, и блок где се
  убацује специфичан садржај.
  - Свака страница која треба овај layout прво га укључи: нpr. web/templates/dashboard/index.html:1-2 позива {{ template "layouts/base" . }}. То каже html/
  template механизму да покрене базни шаблон и да му проследи текући gin.H контекст.
  - Унутар layout-а, позиција за садржај је у <main> преко {{ block "content" . }} … {{ end }} (видљиво око web/templates/layouts/base.html:150). Страница
  која се рендерује дефинише тај блок: нpr. web/templates/dashboard/index.html:5 има {{ define "dashboard/content" }} и тај блок се убацује када хендлер
  попуни ContentTemplate поље у контексту.
  - Када рутер одговори на захтев, нпр. renderDashboard у internal/handler/gui.go:54, он враћа ctx.HTML(http.StatusOK, "dashboard/index.html",
  gin.H{ ... }). Gin/html/template тада изврши dashboard/index.html, што аутоматски позове layout и попуни блок.

  Значи: Go template систем прво учита све шablone, затим за сваку страницу покреће layout као базу и убацује конкретан садржај преко block/template
  механизма.
  ``
## Polja 
Trenutno E17 nema poseban unos u mapi cellOffsets, pa koristi podrazumevani pomak (dx = 0.0, dy = -0.9). Ako želiš da ga gurneš udesno, dodaj prilagođeni
  offset u mapu, npr.:

  var cellOffsets = map[string]textOffset{
      "H11": {dx: 1.6, dy: -0.9},
      "K11": {dx: -8.0, dy: -0.9},
      "F14": {dx: 2.0, dy: -0.3},
      "E20": {dx: 4.0, dy: -0.3},
      "F24": {dx: 2.2, dy: -0.9},
      "E17": {dx: 2.0, dy: -0.9}, // ← nov offset
  }

  – dx је хоризонтални помак у милиметрима (већа вредност = више удесно)
  – dy је вертикални (негативна вредност помера нагоре, позитивна надоле)
  ``