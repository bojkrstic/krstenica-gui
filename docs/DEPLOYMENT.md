# Deployment

## Docker image

`Dockerfile` pravi Go binary iz `./cmd/krstenica` i u runtime image kopira:

- `server`
- `web/`
- `pictures/`

Za kompletan PDF/XLSX export runtime mora imati i `doc/template_files/`. Ako image ne kopira taj direktorijum, montirati ga kao volume na `/app/doc`.

## Lokalni Docker Compose

`docker-compose.yaml` trenutno podize PostgreSQL:

```bash
docker compose up -d krstenica_database
```

Baza je izlozena na host portu `5560`, a container se zove `krstenica-db`.

## Server konfiguracija

Na serveru aktivan DB host treba da bude Docker network naziv baze:

```yaml
db:
  url: postgresql://admin:secret@krstenica-db:5432/krstenica?sslmode=disable
```

Lokalne varijante treba ostaviti zakomentarisane ili premestiti u `config/config.local.yaml`.

`auth.session_secret`, JWT tajne i lozinke ne treba drzati u javnim dokumentima. Za produkciju ih podesavati kroz privatni config ili environment promenljive.

## Build i push

Script:

```bash
./build-and-push.sh
```

Pre pokretanja proveriti tag u `build-and-push.sh`, npr:

```bash
IMAGE="bojankrlekrstic/krstenica-svc:version1.1.4"
```

Script radi `docker build` i `docker push`.

## Rucni rollout na serveru

Tipican tok:

```bash
docker ps
docker rm -f krle-krstenica-svc-1
docker run -d --name krle-krstenica-svc-1 \
  --network krstenica-gui_global \
  -p 8011:8011 \
  -v ~/app/krstenica-gui/config:/app/config \
  -v ~/app/krstenica-gui/doc:/app/doc \
  bojankrlekrstic/krstenica-svc:versionX.Y.Z
docker logs krle-krstenica-svc-1 --tail=50
```

## Pre deployment provera

- Proveriti branch i poslednji commit.
- Proveriti da image tag nije vec koriscen za drugu verziju.
- Proveriti `config/config.yaml` za server DB host.
- Proveriti da su `doc/template_files` dostupni u runtime-u.
- Nakon starta otvoriti `/ui` i probati jedan export krstenice ako je menjana stampa.
