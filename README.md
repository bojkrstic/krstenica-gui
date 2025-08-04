# krstenica
projekat za pravoslavnu crkvu za dodavanje krstenica

Pokrenese docker-compose
1. docker compose up -d


get http://localhost:8001/api/v1/adminv2/tample/1

1. go build -o krstenica-api  // pokrene se ovako  ./krstenica-api
2. go run main.go     // pokrene se i postmain da se prati debug

docker krstenica
docker exec -it krstenica_db psql -U admin krstenica
Kad se udje u kontejner onda se udje u posgress

1. docker exec -it krstenica_db sh
2. / # psql -U admin krstenica
3. \dt 
4. select * from public.eparhije;

brisanje servisa
1. lsof -i :8001
2. kill -9 PID    ---- PID id broj servisa
3. kill -9 5778

execution version
make build-linux

Ako imamo problema prilikom kreiranja baze, moze biti razlog tabela schema_migrations, ona ima dva polja version i dirty, dirty mora da bude na f=false
update schema_migrations set dirty=false;

1. varijanta ---  Da se pokrene ./krstenica i da se onda iz postmen-a gadja adresa (primer: GET http://localhost:8001/api/v1/adminv2/tamples/1)
2. varijatna ---  Da se pokrene debug ali naravno prethodno mora da imamo ciste portove(lsof -i :8001, i kill -9 PID). i onda da se gadja url (primer: GET http://localhost:8001/api/v1/adminv2/tamples/1). I u tom slucaju se krece kroz debug tacke.
