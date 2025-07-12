# go-ping

`go-ping` to narzędzie diagnostyczne napisane w Go, które działa podobnie do polecenia `ping`. Umożliwia wysyłanie pakietów ICMP Echo Request do wybranego hosta i analizowanie odpowiedzi.

## 🛠 Funkcjonalność

- Pingowanie dowolnego hosta (IP lub hostname)
- Konfigurowalna liczba pakietów
- Konfigurowalny timeout i interwał między pakietami
- Możliwość przesyłania własnego payloadu w pakietach ICMP
- Statystyki RTT i utraconych pakietów

## 📦 Wymagania

- Go 1.21+
- Linux (wymagane uprawnienia root lub `CAP_NET_RAW`)

## 🔧 Instalacja

```bash
git clone https://github.com/twoje-repozytorium/go-ping.git
cd go-ping
go build -o go-ping main.go
```

## 🚀 Użycie

```bash
sudo ./go-ping -host=8.8.8.8 -count=4 -timeout=3s -interval=1s -payload="PING"
```

### Dostępne flagi

| Flaga    | Opis                                                              |
|----------|-------------------------------------------------------------------|
| `-host`  | Adres IP lub nazwa hosta do pingowania (wymagane)                 |
| `-count` | Liczba wysyłanych pakietów (domyślnie 4)                           |
| `-timeout` | Timeout dla jednej odpowiedzi, np. 3s (domyślnie 3s)             |
| `-interval` | Odstęp między kolejnymi pingami (domyślnie 1s)                 |
| `-payload`  | Treść wysyłana jako dane ICMP (domyślnie "PING")               |

## 🧪 Testy

Uruchomienie testów jednostkowych:

```bash
go test
```

Testy obejmują poprawność budowy pakietów ICMP (`createICMPEchoMessage`).

## 🐳 Docker

Możesz uruchomić aplikację w kontenerze:

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app
COPY . .

RUN go build -o go-ping main.go

ENTRYPOINT ["./go-ping"]
```

Uruchamianie kontenera:

```bash
docker build -t go-ping .
docker run --cap-add=NET_RAW --network=host go-ping -host=8.8.8.8
```

## 📄 Licencja

Projekt udostępniany jest na licencji Apache 2.0 — szczegóły w pliku `LICENSE`.

## ✍️ Autor

Maciej Brzeźny – 2025

---

Aplikacja nie posiada jeszcze interfejsu gRPC ani UI. W razie potrzeby można je rozbudować w przyszłości (np. w oddzielnym kontenerze).
