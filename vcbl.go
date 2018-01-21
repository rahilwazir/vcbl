package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli"
)

const (
	dictionaryURL = "https://www.vocabulary.com/dictionary/definition.ajax?search=%s"
	audioURL      = "https://audio.vocab.com/1.0/us/%s.mp3"
)

var (
	retry = 0
	doc   *goquery.Document
	c     *cli.Context
)

func cliAction(cli *cli.Context) error {
	c = cli

	lookup := c.Args().Get(0)

	if lookup == "" {
		fmt.Println("No lookup!")
		return nil
	}

	if !queryDocument(lookup) {
		return nil
	}

	definition := getDefinition(lookup)
	if definition == "" {
		return nil
	}

	fmt.Print(definition)

	pronounceWord()

	return nil
}

func queryDocument(lookup string) bool {
	document, err := goquery.NewDocument(fmt.Sprintf(dictionaryURL, lookup))
	if err != nil {
		log.Fatal(err)
	}

	doc = document

	ret, _ := doc.Html()
	verbose("HTML: " + ret)

	if doc.Find("div.noresults").Length() == 1 {
		if retry == 0 {
			fmt.Println("Not found.")
			lookup = strings.Title(lookup)
			fmt.Printf("Trying %q...\n", lookup)
			retry++
			return queryDocument(lookup)
		}

		fmt.Println("Not found.")
		return false
	}

	return true
}

func getDefinition(lookup string) string {
	shortDesc := doc.Find("p.short").Text()
	longDesc := doc.Find("p.long").Text()

	text := fmt.Sprintln()

	if doc.Find(".blurb").Length() != 0 {
		switch c.String("desc") {
		case "long":
			text += fmt.Sprintln(longDesc)
		case "both":
			text += fmt.Sprintln(shortDesc)
			text += fmt.Sprintln("------------------")
			text += fmt.Sprintln(longDesc)
		default:
			text += fmt.Sprintln(shortDesc)
		}
	}

	if doc.Find("table.definitionNavigator").Length() != 0 {
		doc.Find("table.definitionNavigator tr").Each(func(i int, tr *goquery.Selection) {
			groupNumber := tr.Find(".groupNumber").Text()
			def := tr.Find(".def").Text()

			text += fmt.Sprintf("%s %s\n", groupNumber, def)
		})
	} else {
		doc.Find("div.ordinal").Each(func(i int, ordinal *goquery.Selection) {
			if i == 0 {
				groupNumber := ordinal.Closest("div.group").Find(".groupNumber").Text()
				text += fmt.Sprintf("%s.\n", groupNumber)
			}

			wordType := ordinal.Find("h3.definition").Find(".anchor").Text()
			text += fmt.Sprintf("(%s) ", wordType)

			h3Definition := ordinal.Find("h3.definition").First().Contents().Eq(2).Text()
			text += fmt.Sprintf("%s\n", strings.TrimSpace(h3Definition))
		})
	}

	return text
}

func pronounceWord() {
	if c.Bool("play") == false {
		return
	}

	audioEl := doc.Find(".audio")

	if audioEl.Length() == 0 {
		verbose("Audio: Documnt element not found.")
		return
	}

	verbose("Audio: Playing...")

	uri, _ := audioEl.Attr("data-audio")
	qualifiedAudioURL := fmt.Sprintf(audioURL, uri)
	err := exec.Command("play", qualifiedAudioURL).Run()
	if err != nil {
		log.Fatal(err)
	}
}

func verbose(output string) {
	if c.Bool("verbose") {
		log.Println(output)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "vcbl-cli"
	app.Usage = "Vocabulary.com CLI dictionary"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "desc",
			Value: "short",
			Usage: "Description type of the lookup word. Possible values are: short, long, both",
		},
		cli.BoolFlag{
			Name:  "play",
			Usage: "Play the word pronounciation with SoX cli. SoX must be installed",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Debug output",
		},
	}

	app.Action = cliAction

	app.Run(os.Args)
}
