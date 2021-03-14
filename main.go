package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var audioFilename, audioFormat, performer, title, genre, date string
	flag.StringVar(&audioFilename, "filename", "", "name of the file containing the audio data")
	flag.StringVar(&audioFormat, "format", "", "format of the file containing the audio data")
	flag.StringVar(&title, "title", "", "album title")
	flag.StringVar(&performer, "performer", "", "album performer")
	flag.StringVar(&genre, "genre", "", "album genre")
	flag.StringVar(&date, "date", "", "album date")
	flag.Parse()

	if genre != "" {
		fmt.Printf("REM GENRE %s\n", genre)
	}

	if date != "" {
		fmt.Printf("REM DATE %s\n", date)
	}

	if performer != "" {
		fmt.Printf("PERFORMER \"%s\"\n", performer)
	}

	if title != "" {
		fmt.Printf("TITLE \"%s\"\n", title)
	}

	if audioFilename != "" {
		if audioFormat == "" {
			audioFormat = filepath.Ext(audioFilename)
		}

		fmt.Printf("FILE \"%s\" %s\n", audioFilename, audioFormat)
	}

	convert := func(input io.Reader) {
		scanner := bufio.NewScanner(input)
		for i := 1; ; {
			if !scanner.Scan() {
				err := scanner.Err()
				if err != nil {
					log.Println(err)
				}
				break
			}

			line := scanner.Text()
			if line == "" {
				continue
			}

			indexArr := strings.SplitN(line, " ", 2)
			index := indexArr[0]
			if index == "" {
				log.Println("no index for track")
				continue
			}
			if len(indexArr) < 2 {
				log.Printf("no artist or title info for track at position %s\n", index)
				continue
			}

			info := indexArr[1]
			infoArr := strings.Split(info, "-")

			var performer, title string
			switch len(infoArr) {
			case 0:
				log.Printf("no artist or title info for track at position %s\n", index)
			case 1:
				title = strings.TrimSpace(infoArr[0])
			default:
				if len(infoArr) > 2 {
					log.Printf("too many dashes in track info: \"%s\"\n", info)
				}
				performer = strings.TrimSpace(infoArr[0])
				title = strings.TrimSpace(infoArr[1])
			}

			fmt.Printf("  TRACK %d AUDIO\n", i)
			if title != "" {
				fmt.Printf("    TITLE \"%s\"\n", title)
			}
			if performer != "" {
				fmt.Printf("    PERFORMER \"%s\"\n", performer)
			}
			fmt.Printf("    INDEX 01 %s\n", index+":00")

			i++
		}
	}

	args := flag.Args()
	if len(args) == 0 {
		convert(os.Stdin)
	} else {
		for _, filename := range args {
			file, err := os.Open(filename)
			if err != nil {
				log.Println(err)
				continue
			}

			convert(file)
		}
	}
}
