# krstenica
projekat za pravoslavnu crkvu za dodavanje krstenica


get http://localhost:8001/api/v1/adminv2/tample/1

1. go build -o krstenica-api  // pokrene se ovako  ./krstenica-api
2. go run main.go     // pokrene se i postmain da se prati debug

docker krstenica
docker exec -it krstenica_db psql -U admin krstenica

brisanje servisa
1. lsof -i :8001
2. kill -9 PID    ---- PID id broj servisa
3. kill -9 5778