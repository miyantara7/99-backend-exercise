# 99 Backend Exercise

Project ini terdiri dari 3 service:

- `listing_service.py`
  Listing service berbasis Python Tornado
- `user-service`
  User service berbasis Go + Gin
- `public-api`
  Public API berbasis Go + Gin

Flow utamanya:

`public-api -> user-service -> listing-service`

## Struktur

```text
.
├── cmd/dev
├── internal/runner
├── listing_service.py
├── public-api
└── user-service
```

- `cmd/dev`
  Entrypoint untuk menjalankan semua service sekaligus saat development
- `internal/runner`
  Orchestrator lokal yang menyalakan `listing-service`, `user-service`, dan `public-api`
- `public-api/cmd/api`
  Entrypoint service `public-api`
- `user-service/cmd/api`
  Entrypoint service `user-service`

## Requirement

- Go
- Python 3
- Virtual environment Python di `.venv`
- Dependency Python dari `python-libs.txt`

Kalau belum install dependency Python:

```bash
python3 -m venv .venv
source .venv/bin/activate
python -m pip install -U pip
python -m pip install -r python-libs.txt
```

## Menjalankan Semua Service

Dari root project:

```bash
make run
```

Command ini akan menjalankan:

- `listing-service` di `http://127.0.0.1:6300`
- `user-service` di `http://127.0.0.1:6301`
- `public-api` di `http://127.0.0.1:7300`

Untuk stop semua service:

- tekan `Ctrl+C`

## Port Default

- `listing-service`: `6300`
- `user-service`: `6301`
- `public-api`: `7300`

Port alternatif:

```bash
make run-alt
```

Port yang dipakai:

- `listing-service`: `6400`
- `user-service`: `6401`
- `public-api`: `7400`

Jika port default sedang dipakai proses lain, gunakan:

```bash
make run-alt
```

Atau jalankan manual dengan port sendiri:

```bash
go run ./cmd/dev --listing-port=6500 --user-port=6501 --public-port=7500
```

## Build Dev Runner

```bash
make build
```

Output binary:

```text
.run/dev-runner
```

## API Endpoint

### User Service

- `GET /users`
- `GET /users/{id}`
- `POST /users`

Base URL:

```text
http://127.0.0.1:6301
```

### Listing Service

- `GET /listings/ping`
- `GET /listings`
- `POST /listings`

Base URL:

```text
http://127.0.0.1:6300
```

### Public API

- `GET /public-api/listings`
- `POST /public-api/users`
- `POST /public-api/listings`

Base URL:

```text
http://127.0.0.1:7300
```

## Contoh Curl

### User Service

Get all users:

```bash
curl -i http://127.0.0.1:6301/users
```

Get users dengan pagination:

```bash
curl -i "http://127.0.0.1:6301/users?page_num=1&page_size=10"
```

Get user by id:

```bash
curl -i http://127.0.0.1:6301/users/1
```

Create user:

```bash
curl -i -X POST http://127.0.0.1:6301/users \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "name=Suresh Subramaniam"
```

### Listing Service

Ping:

```bash
curl -i http://127.0.0.1:6300/listings/ping
```

Get listings:

```bash
curl -i http://127.0.0.1:6300/listings
```

Get listings dengan pagination:

```bash
curl -i "http://127.0.0.1:6300/listings?page_num=1&page_size=10"
```

Get listings by user:

```bash
curl -i "http://127.0.0.1:6300/listings?page_num=1&page_size=10&user_id=1"
```

Create listing:

```bash
curl -i -X POST http://127.0.0.1:6300/listings \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "user_id=1&listing_type=rent&price=6000"
```

### Public API

Get public listings:

```bash
curl -i http://127.0.0.1:7300/public-api/listings
```

Get public listings dengan pagination:

```bash
curl -i "http://127.0.0.1:7300/public-api/listings?page_num=1&page_size=10"
```

Get public listings by user:

```bash
curl -i "http://127.0.0.1:7300/public-api/listings?page_num=1&page_size=10&user_id=1"
```

Create user via public API:

```bash
curl -i -X POST http://127.0.0.1:7300/public-api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Lorel Ipsum"}'
```

Create listing via public API:

```bash
curl -i -X POST http://127.0.0.1:7300/public-api/listings \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"listing_type":"rent","price":6000}'
```
