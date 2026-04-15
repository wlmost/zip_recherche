# zip_recherche

Kommandozeilenwerkzeug zur Recherche in ZIP-Archiven, die vom **Eddyson ySE Konverter** erzeugt werden.

## Funktionen

- Inhalt eines ZIP-Archivs auflisten
- Nach Dateinamen innerhalb eines Archivs suchen
- Textzeilen in Archiv-Inhalten per regulärem Ausdruck durchsuchen (mehrere Muster als AND-Verknüpfung)
- Dateien aus einem Archiv extrahieren
- Verzeichnisbäume mit mehreren Archiven rekursiv durchsuchen

## Voraussetzungen

- Go 1.15 oder neuer

## Build

```bash
go build -o zip_recherche .
```

## Verwendung

```
zip_recherche [Optionen]
```

### Optionen

| Option | Standard     | Beschreibung |
|--------|-------------|--------------|
| `-d`   | `.`         | Verzeichnis mit ZIP-Archiven, das durchsucht werden soll |
| `-f`   | _(leer)_    | Einzelnes ZIP-Archiv, das durchsucht werden soll |
| `-l`   | `false`     | Inhalt des Archivs auflisten |
| `-p`   | _(leer)_    | Suchmuster (regulärer Ausdruck); mehrere Muster kommagetrennt (AND-Verknüpfung) |
| `-s`   | _(leer)_    | Nach einem bestimmten Dateinamen im Archiv suchen |
| `-v`   | `false`     | Zeigt die Zeile, in der das mit `-p` angegebene Muster gefunden wurde |
| `-x`   | `false`     | Extrahiert die mit `-s` gefundenen Dateien |
| `-t`   | `/tmp`      | Zielverzeichnis für die Extraktion |

## Beispiele

**Inhalt eines Archivs auflisten:**
```bash
zip_recherche -f 1.zip -l
```

**Nach einer Datei im Archiv suchen:**
```bash
zip_recherche -f 1.zip -s bestellung.xml
```

**Alle Archive in einem Verzeichnis nach einem Muster durchsuchen:**
```bash
zip_recherche -d ./2021/09/01 -p "Artikelnummer"
```

**Mehrere Muster kombiniert suchen (AND):**
```bash
zip_recherche -d ./2021/09/01 -p "Artikelnummer,10042"
```

**Datei aus einem Archiv extrahieren:**
```bash
zip_recherche -f 1.zip -s bestellung.xml -x -t C:\Temp\export
```

**Gesamtes Verzeichnis rekursiv mit Mustersuche durchsuchen:**
```bash
zip_recherche -d ./zip_recherche/2021 -p "Lieferant"
```

## Verzeichnisstruktur der Archive

Die vom Eddyson ySE Konverter erzeugten ZIP-Archive liegen typischerweise in einer Datumsstruktur vor:

```
zip_recherche/
  2021/
    09/
      01/
      02/
      ...
```

Das Tool durchläuft diese Struktur automatisch, wenn mit `-d` ein übergeordnetes Verzeichnis angegeben wird.

## Abhängigkeiten

- [`github.com/gabriel-vasile/mimetype`](https://github.com/gabriel-vasile/mimetype) – MIME-Typ-Erkennung zur Unterscheidung von Text- und Binärdateien
