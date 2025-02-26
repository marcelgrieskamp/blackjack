# Blackjack Spiel

Ein browser-basiertes Blackjack-Spiel mit Casino-Atmosphäre, implementiert in Go und JavaScript.

## Features

- Authentisches Blackjack-Spielerlebnis mit klassischen Regeln
- Visuell ansprechendes Casino-Design mit animierten Elementen
- Realistisches Wett-System mit verschiedenen Chip-Werten
- Responsive Design für verschiedene Bildschirmgrößen
- Intuitive Benutzeroberfläche mit visuellen Effekten

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

### Grundregeln
1. Klicken Sie auf "Neues Spiel" um zu beginnen
2. Platzieren Sie Ihren Einsatz mit den Chips oder der "ALL IN"-Option
3. Klicken Sie auf "Einsatz platzieren" um die Runde zu starten
4. Nutzen Sie "Hit" um eine weitere Karte zu ziehen
5. Nutzen Sie "Stand" um Ihre Hand zu halten
6. Ziel ist es, näher an 21 Punkte heranzukommen als der Dealer, ohne dabei über 21 zu gehen

### Wettsystem
- Sie beginnen mit 1000 Chips
- Wählen Sie Chips mit unterschiedlichen Werten (10, 25, 50, 100, 200, 500, 1000)
- Mit "ALL IN" setzen Sie alle verfügbaren Chips
- Bei Gewinn verdoppeln Sie Ihren Einsatz
- Bei Verlust verlieren Sie Ihren Einsatz
- Haben Sie keine Chips mehr, können Sie das Spiel mit dem "Neues Guthaben ($1000)"-Button zurücksetzen

## Technologien

- **Backend**: Go (Golang) für den Server und die Spiellogik
- **Frontend**: HTML5, CSS3 und JavaScript
- **Design**: Responsives Casino-Theming mit CSS-Animationen
- **APIs**: RESTful API für die Kommunikation zwischen Frontend und Backend

## Design & Theming

Die Blackjack-Anwendung verfügt über ein ansprechendes Casino-Theming mit folgenden visuellen Elementen:

- Vegas-Stil-Hintergrund mit subtilen Texturen
- Animierte Geldscheine und Lichteffekte
- Neon-Titel und pulsierende Animationen
- Realistische Chipdesigns mit unterschiedlichen Farben
- Kartendesigns mit goldenen Akzenten
- Responsives Layout für verschiedene Geräte

## Konfiguration

Der Server läuft standardmäßig auf Port 8080. Sie können den Port durch Setzen der Umgebungsvariable `PORT` ändern:

```bash
PORT=3000 ./blackjack-server
```

## Fehlerbehebung

### Server startet nicht
- Stellen Sie sicher, dass der Port nicht bereits belegt ist
- Überprüfen Sie die Firewall-Einstellungen
- Stellen Sie sicher, dass Sie die erforderlichen Berechtigungen haben

### Anzeige-Probleme
- Stellen Sie sicher, dass Sie einen modernen Browser verwenden
- Leeren Sie den Cache Ihres Browsers
- Vergewissern Sie sich, dass JavaScript aktiviert ist

### Spielprobleme
- Wenn das Spiel nicht reagiert, laden Sie die Seite neu
- Bei fehlenden Chips verwenden Sie den "Neues Guthaben" Button
- Bei anderen Problemen starten Sie den Server neu

## Lizenz

Dieses Projekt ist unter der MIT-Lizenz lizenziert. Siehe die [LICENSE](LICENSE) Datei für Details. 