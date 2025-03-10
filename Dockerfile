# Etap budowania: używamy obrazu golang do kompilacji aplikacji
FROM golang:1.19 AS builder
WORKDIR /app

# Skopiuj pliki modułów i pobierz zależności
COPY go.mod go.sum ./
RUN go mod download

# Skopiuj cały kod źródłowy
COPY . .

# Budujemy aplikację w trybie statycznym (bez zależności dynamicznych)
RUN CGO_ENABLED=0 GOOS=linux go build -a -o go-ping .

# Etap końcowy: minimalny obraz do uruchomienia aplikacji
FROM alpine:latest
# Kopiujemy skompilowany plik binarny z etapu builder
COPY --from=builder /app/go-ping /go-ping

# Ustawiamy punkt wejścia
ENTRYPOINT ["/go-ping"]
