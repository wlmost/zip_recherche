package main

import (
	"archive/zip"
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func doAction(path string, patterns []*regexp.Regexp, file string, list bool, line bool, extract bool, targetdir string) {
	buf, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	mtype := mimetype.Detect(buf)
	if mtype.String() != "application/zip" {
		msg := "Datei %s ist kein zip-Archiv"
		log.Fatalf(msg, path)
	}
	zh, err := zip.OpenReader(path)
	if err != nil {
		msg := "Konnte zip-Archiv %s nicht öffnen: %s\n"
		log.Fatalf(msg, path, err)
	}
	defer zh.Close()

	// Dateilisting
	if list {
		for _, f := range zh.File {
			fmt.Printf("-> %s\n", f.Name)
		}
	}

	// Extraktion der gegebenen Datei
	if file != "" {
		for _, fn := range strings.Split(file, ",") {
			for _, fh := range zh.Reader.File {
				if filepath.Base(fh.Name) == fn {
					fmt.Printf("Datei '%s' gefunden\n", fh.Name)
					if extract {
						fmt.Printf("Extrahiere %s nach %s\n", fh.Name, targetdir)
						zfile, err := fh.Open()
						if err != nil {
							log.Fatalf(err.Error())
						}
						defer zfile.Close()
						if _, err := os.Stat(targetdir); os.IsNotExist(err) {
							os.MkdirAll(targetdir, fh.Mode())
						}
						extractedFilePath := filepath.Join(targetdir, fn)
						if fh.FileInfo().IsDir() {
							os.MkdirAll(extractedFilePath, fh.Mode())
						} else {
							output, err := os.OpenFile(
								extractedFilePath,
								os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
								fh.Mode(),
							)
							if err != nil {
								log.Fatal(err.Error())
							}
							defer output.Close()
							_, err = io.Copy(output, zfile)
							if err != nil {
								log.Fatalf(err.Error())
							}
						}
					}
				}
			}
		}
	}

	// Suche nach Muster
	if len(patterns) > 0 {
		for _, fh := range zh.Reader.File {
			zfile, err := fh.Open()
			if err != nil {
				log.Fatalf(err.Error())
			}
			defer zfile.Close()
			scanner := bufio.NewScanner(zfile)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				allMatch := true
				var lastMatches []string
				for _, re := range patterns {
					m := re.FindStringSubmatch(scanner.Text())
					if len(m) == 0 {
						allMatch = false
						break
					}
					lastMatches = m
				}
				if allMatch {
					mtype := mimetype.Detect([]byte(scanner.Text()))
					if strings.Index(mtype.String(), "text/plain") > -1 {
						fmt.Printf("%s: %s %s\n", fh.Name, scanner.Text(), lastMatches)
					} else {
						fmt.Printf("binäre Datei %s enthält Suchmuster\n", fh.Name)
					}
				}
			}
		}
	}
}

func main() {
	// flag declaration
	optDir := flag.String("d", ".", "Angabe des Verzeichnes mit Zip-Archiven, die durchsucht werden sollen. Vorgabe ist das aktuelle Verzeichnis.")
	optZip := flag.String("f", "", "Angabe des Zip-Archives, das durchsucht werden soll. Vorgabe ist 1.zip.")
	optLst := flag.Bool("l", false, "Listet den Inhalt des angegebenen Zip-Archives.")
	optPat := flag.String("p", "", "Angabe eines Suchmusters. Liefert den Dateinamen in dem das Muster gefunden wurde.")
	optFnam := flag.String("s", "", "Sucht nach dem angegebenen Dateinamen im Zip-Archiv.")
	optLine := flag.Bool("v", false, "Zeigt die Zeile, in der das mit '-p' angebebene Muster gefunden wurde.")
	optExt := flag.Bool("x", false, "Extrahiert die Datei(en) die mit dem Paramter '-s' gesucht wurden.")
	optTD := flag.String("t", "/tmp", "Extrahiert die Datei(en) in das angegebene Verzeichnis.")

	// Hilfeausgabe
	//flag.Usage = func() {

	//}

	// parse flags
	flag.Parse()

	if len(os.Args) == 1 {
		fmt.Println("Bitte ein Argument angeben")
		flag.Usage()
		os.Exit(1)
	}

	if *optDir == "" {
		*optDir = "."
	}

	if *optTD == "" {
		*optTD = "/tmp"
	}

	// Suchmuster kompilieren (kommagetrennt, werden der Reihe nach als AND-Filter angewendet)
	var patterns []*regexp.Regexp
	if *optPat != "" {
		for _, p := range strings.Split(*optPat, ",") {
			re, err := regexp.Compile(strings.TrimSpace(p))
			if err != nil {
				log.Fatalf("Ungültiges Suchmuster '%s': %s\n", p, err)
			}
			patterns = append(patterns, re)
		}
	}

	if *optZip != "" {
		path := *optDir + "/" + *optZip
		abspath, _ := filepath.Abs(path)
		doAction(abspath, patterns, *optFnam, *optLst, *optLine, *optExt, *optTD)
	} else {
		filepath.Walk(*optDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatalf(err.Error())
			}
			if info.IsDir() {
				fmt.Printf("Dir: %s\n", info.Name())
			} else {
				fmt.Printf("File: %s\n", info.Name())
				abspath, _ := filepath.Abs(path)
				doAction(abspath, patterns, *optFnam, *optLst, *optLine, *optExt, *optTD)
			}
			return nil
		})
	}
}
