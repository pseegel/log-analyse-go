# Log-Analyse (Go)

Nebenläufiger Parser für Webserver-Logs in Go.

## Motivation

Dieses Projekt ist die Go-Variante des Python-Projekts [log-analyse-script](https://github.com/pseegel/log-analyse-script). Ziel ist die Vertiefung idiomatischer Go-Concurrency-Muster (Goroutines, Channels, Context) an einem konkreten Problem.

## Anforderungen

- Go >= 1.26
- Keine externen Abhängigkeiten zur Laufzeit

## Nutzung

Direkt mit `go run`:

```bash
go run . -input access.log -output report.json
```

Oder als Binary bauen und aufrufen:

```bash
go build
./log-analyse-go -input access.log -output report.json      # Linux / macOS
.\log-analyse-go.exe -input access.log -output report.json  # Windows
```

Hilfe anzeigen:

```bash
./log-analyse-go -h
```

Standardmäßig wird `access.log` im aktuellen Verzeichnis erwartet und nach `report.json` geschrieben. Input- und Output-Dateien können per Flag überschrieben werden.

## Beispiel-Input

[Beispiel](access.log)

## Architektur

_Wird nach der ersten Implementierung ergänzt._

## Entwicklung

_Wird nach der ersten Implementierung ergänzt._

## Lizenz

MIT — siehe [LICENSE](LICENSE.md).
