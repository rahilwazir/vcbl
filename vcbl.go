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
			scrapeDefinition(c, lookup)
		} else {
			fmt.Println("No lookup!")
		}
		return nil
	}
	app.Run(os.Args)
}

func scrapeDefinition(c *cli.Context, lookup string) {
	doc, err := goquery.NewDocument(fmt.Sprintf(url, lookup))
	if err != nil {
		log.Fatal(err)
	}

	shortDesc := doc.Find("p.short").Text()
	longDesc := doc.Find("p.long").Text()

	fmt.Println()

	switch c.String("desc") {
	case "long":
		fmt.Println(longDesc)
	case "both":
		fmt.Println(shortDesc)
		fmt.Println("------------------")
		fmt.Println(longDesc)
	default:
		fmt.Println(shortDesc)
	}

	doc.Find("table.definitionNavigator tr").Each(func(i int, tr *goquery.Selection) {
		fmt.Printf("%s %s\n", tr.Find(".groupNumber").Text(), tr.Find(".def").Text())
	})
}
