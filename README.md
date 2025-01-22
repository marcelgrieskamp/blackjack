# Blackjack Spiel

Ein browser-basiertes Blackjack-Spiel, implementiert in Go und JavaScript.

## Voraussetzungen

- Go 1.16 oder höher
- Ein moderner Webbrowser

## Installation

1. Klonen Sie das Repository:
```bash
git clone https://github.com/marcelgrieskamp/blackjack.git
cd blackjack
```

2. Abhängigkeiten installieren:
```bash
go mod tidy
```

## Starten des Servers

1. Starten Sie den Server mit:
```bash
go run cmd/main.go
```

2. Öffnen Sie einen Webbrowser und navigieren Sie zu:
```
http://localhost:8080
```

## Deployment auf einem Server

### Option 1: Direktes Deployment

1. Kompilieren Sie das Programm für Ihr Zielsystem:
```bash
GOOS=linux GOARCH=amd64 go build -o blackjack-server cmd/main.go
```

2. Übertragen Sie die folgenden Dateien auf Ihren Server:
   - Der kompilierte `blackjack-server`
   - Der komplette `static/` Ordner

3. Starten Sie den Server:
```bash
./blackjack-server
```

### Option 2: Docker Deployment

1. Erstellen Sie ein Docker Image:
```bash
docker build -t blackjack-game .
```

2. Starten Sie den Container:
```bash
docker run -p 8080:8080 blackjack-game
```

## Spielanleitung

1. Klicken Sie auf "Neues Spiel" um zu beginnen
2. Nutzen Sie "Hit" um eine weitere Karte zu ziehen
3. Nutzen Sie "Stand" um Ihre Hand zu halten
4. Ziel ist es, näher an 21 Punkte heranzukommen als der Dealer, ohne dabei über 21 zu gehen

## Konfiguration

Der Server läuft standardmäßig auf Port 8080. Sie können den Port durch Setzen der Umgebungsvariable `PORT` ändern:

```bash
PORT=3000 ./blackjack-server
``` 