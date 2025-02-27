package parser

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

const baseURL = "https://lenta.com"

var categories = map[string]string{
	"Молочные продукты": "moloko-syr-yaytsa",
	"Хлеб":              "hleb-i-vypechka",
	"Овощи и фрукты":    "ovoshchi-i-frukty",
	"Мясо и птица":      "myaso-ptitsa-kolbasa",
}

func ParserCategories(city string) ([][]string, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("lenta.com"),
	)

	var products [][]string
	products = append(products, []string{"Категория", "Название товара", "Цена", "Ссылка"})

	for category, slug := range categories {
		url := fmt.Sprintf("%s/%s/catalog/%s/", baseURL, city, slug)
		c.OnHTML("div.catalog-product", func(e *colly.HTMLElement) {
			name := e.ChildText("div.product-card__title")
			price := e.ChildText("div.product-card__price-current")
			link := e.ChildAttr("a.product-card__link", "href")

			if name == "" || price == "" || link == "" {
				return
			}

			link = baseURL + link
			products = append(products, []string{category, name, strings.TrimSpace(price), link})

		})
		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})
		err := c.Visit(url)
		if err != nil {
			return nil, err
		}
	}
	return products, nil
}
