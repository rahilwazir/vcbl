package main

import (
	"fmt"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli"
)

const (
	url = "https://www.vocabulary.com/dictionary/definition.ajax?search=%s"
)

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

func getDefinition(c *cli.Context, lookup string) {
	doc, err := goquery.NewDocument(fmt.Sprintf(url, lookup))
	if err != nil {
		log.Fatal(err)
	}

	if doc.Find("div.noresults").Length() == 1 {
		fmt.Println("No results.")
		return
	}

	shortDesc := doc.Find("p.short").Text()
	longDesc := doc.Find("p.long").Text()

	text := fmt.Sprintln()

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

	doc.Find("table.definitionNavigator tr").Each(func(i int, tr *goquery.Selection) {
		text += fmt.Sprintf("%s %s\n", tr.Find(".groupNumber").Text(), tr.Find(".def").Text())
	})

	fmt.Print(text)
}
