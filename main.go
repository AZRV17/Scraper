package main

import (
	"encoding/csv"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"log"
	"os"
)

type Influencer struct {
	Rank        string   `csv:"Rank"`
	Name        string   `csv:"Name"`
	Title       string   `csv:"Title"`
	Topics      []string `csv:"Topics"`
	Subscribers string   `csv:"Subscribers"`
	Audience    string   `csv:"Audience"`
	Authentic   string   `csv:"Authentic"`
	Engagement  string   `csv:"Engagement"`
}

func main() {
	g := geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{"https://hypeauditor.com/top-instagram-all-russia/"},
		ParseFunc: parseInfluencers,
	})

	g.Start()
}

func parseInfluencers(g *geziyor.Geziyor, r *client.Response) {
	err := os.Remove("influencers.csv")
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatal(err)
		}
	}

	file, err := os.Create("influencers.csv")
	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(file)

	defer writer.Flush()

	headers := []string{
		"Rank",
		"Nick",
		"Name",
		"Category",
		"Followers",
		"Country",
		"Eng (Auth.)",
		"Eng (Avg.)",
	}

	_ = writer.Write(headers)

	r.HTMLDoc.Find("div.row__top").Each(func(i int, s *goquery.Selection) {
		influencer := &Influencer{
			Rank:  s.Find("div.row-cell.rank span").First().Text(),
			Name:  s.Find("div.contributor__name-content").Text(),
			Title: s.Find("div.contributor__title").Text(),
		}

		s.Find("div.tag__content.ellipsis").Each(func(i int, s *goquery.Selection) {
			influencer.Topics = append(influencer.Topics, s.Text())
		})

		influencer.Subscribers = s.Find("div.row-cell.subscribers").Text()
		influencer.Audience = s.Find("div.row-cell.audience").Text()
		influencer.Authentic = s.Find("div.row-cell.authentic").Text()
		influencer.Engagement = s.Find("div.row-cell.engagement").Text()

		topicsString := ""
		for i, topic := range influencer.Topics {
			topicsString += topic
			if i < len(influencer.Topics)-1 {
				topicsString += ", "
			}
		}
		influencer.Topics = []string{topicsString}

		err = writer.Write([]string{
			influencer.Rank,
			influencer.Name,
			influencer.Title,
			topicsString,
			influencer.Subscribers,
			influencer.Audience,
			influencer.Authentic,
			influencer.Engagement,
		})
	})
}
