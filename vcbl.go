package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli"
)

const (
	url = "https://www.vocabulary.com/dictionary/definition.ajax?search=%s"
)

var (
	retry = 0
)

func getDefinition(c *cli.Context, lookup string) {
	doc, err := goquery.NewDocument(fmt.Sprintf(url, lookup))
	if err != nil {
		log.Fatal(err)
	}

	if doc.Find("div.noresults").Length() == 1 {
		fmt.Println("No results.")

		if retry == 0 {
			lookup = strings.Title(lookup)
			fmt.Printf("Trying %q...", lookup)
			getDefinition(c, lookup)
			retry++
		}

		return
	}

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

	fmt.Print(text)
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
	}

	app.Action = func(c *cli.Context) error {
		if lookup := c.Args().Get(0); lookup != "" {
			getDefinition(c, lookup)
		} else {
			fmt.Println("No lookup!")
		}
		return nil
	}
	app.Run(os.Args)
}
