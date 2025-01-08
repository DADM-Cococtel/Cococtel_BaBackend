package scrappingrepository

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

type (
	IScrapping interface {
		GetLiquor(*gin.Context) error
	}
	scrappingRepository struct {
		url *string
	}
)

func NewScrappingRepository(url *string) IScrapping {
	return &scrappingRepository{url: url}
}

func (sr *scrappingRepository) GetLiquor(ctx *gin.Context) error {
	c := colly.NewCollector(colly.AllowedDomains(*sr.url))
	c.OnHTML("div.col-md-3 product-image", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
	})
	url := fmt.Sprintf("https://www.gs1.org/services/verified-by-gs1/iframe?gtin=%s#productInformation", "7702049101337")
	err := c.Visit(url)
	if err != nil {
		return err
	}
	return nil
}
