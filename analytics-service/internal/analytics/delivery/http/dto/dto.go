package dto

import (
	"analytics-service/domain/entities"
	"sort"
	"time"
)

type ClickedByCountry struct {
	Country string `json:"country"`
	Count   int    `json:"count"`
}

type ClickedByCity struct {
	City  string `json:"city"`
	Count int    `json:"count"`
}

type ClickedByDate struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type ClickedByPlatform struct {
	Platform string `json:"platform"`
	Count    int    `json:"count"`
}

type ClickedByReferer struct {
	Referer string `json:"referer"`
	Count   int    `json:"count"`
}

type URLClickedResponse struct {
	ClickedByCountry  []ClickedByCountry  `json:"clicked_by_country"`
	ClickedByCity     []ClickedByCity     `json:"clicked_by_city"`
	ClickedByDate     []ClickedByDate     `json:"clicked_by_date"`
	ClickedByPlatform []ClickedByPlatform `json:"clicked_by_platform"`
	ClickedByReferer  []ClickedByReferer  `json:"clicked_by_referer"`
	TotalClicks       int                 `json:"total_clicks"`
}

func ToURLClickedResponse(in []entities.URLClicked) URLClickedResponse {
	if len(in) == 0 {
		return URLClickedResponse{
			ClickedByCountry:  []ClickedByCountry{},
			ClickedByCity:     []ClickedByCity{},
			ClickedByDate:     []ClickedByDate{},
			ClickedByPlatform: []ClickedByPlatform{},
			ClickedByReferer:  []ClickedByReferer{},
			TotalClicks:       0,
		}
	}
	countryCount := make(map[string]int)
	cityCount := make(map[string]int)
	dateCount := make(map[string]int)
	platformCount := make(map[string]int)
	refererCount := make(map[string]int)

	for _, click := range in {
		countryCount[click.Country]++
		cityCount[click.City]++
		dateCount[click.Timestamp.Format("2006-01-02")]++
		platformCount[click.Platform]++
		refererCount[click.Referer]++
	}

	clickedByCountry := make([]ClickedByCountry, 0, len(countryCount))
	for country, count := range countryCount {
		clickedByCountry = append(clickedByCountry, ClickedByCountry{
			Country: country,
			Count:   count,
		})
	}

	sort.Slice(clickedByCountry, func(i, j int) bool {
		if clickedByCountry[i].Count == clickedByCountry[j].Count {
			return clickedByCountry[i].Country < clickedByCountry[j].Country
		}
		return clickedByCountry[i].Count > clickedByCountry[j].Count
	})

	clickedByCity := make([]ClickedByCity, 0, len(cityCount))
	for city, count := range cityCount {
		clickedByCity = append(clickedByCity, ClickedByCity{
			City:  city,
			Count: count,
		})
	}
	sort.Slice(clickedByCity, func(i, j int) bool {
		if clickedByCity[i].Count == clickedByCity[j].Count {
			return clickedByCity[i].City < clickedByCity[j].City
		}
		return clickedByCity[i].Count > clickedByCity[j].Count
	})

	clickedByDate := make([]ClickedByDate, 0, len(dateCount))
	for date, count := range dateCount {
		clickedByDate = append(clickedByDate, ClickedByDate{
			Date:  date,
			Count: count,
		})
	}

	sort.Slice(clickedByDate, func(i, j int) bool {
		ti, _ := time.Parse("2006-01-02", clickedByDate[i].Date)
		tj, _ := time.Parse("2006-01-02", clickedByDate[j].Date)
		return ti.After(tj)
	})

	clickedByPlatform := make([]ClickedByPlatform, 0, len(platformCount))
	for platform, count := range platformCount {
		clickedByPlatform = append(clickedByPlatform, ClickedByPlatform{
			Platform: platform,
			Count:    count,
		})
	}
	sort.Slice(clickedByPlatform, func(i, j int) bool {
		if clickedByPlatform[i].Count == clickedByPlatform[j].Count {
			return clickedByPlatform[i].Platform < clickedByPlatform[j].Platform
		}
		return clickedByPlatform[i].Count > clickedByPlatform[j].Count
	})

	clickedByReferer := make([]ClickedByReferer, 0, len(refererCount))
	for referer, count := range refererCount {
		clickedByReferer = append(clickedByReferer, ClickedByReferer{
			Referer: referer,
			Count:   count,
		})
	}
	sort.Slice(clickedByReferer, func(i, j int) bool {
		if clickedByReferer[i].Count == clickedByReferer[j].Count {
			return clickedByReferer[i].Referer < clickedByReferer[j].Referer
		}
		return clickedByReferer[i].Count > clickedByReferer[j].Count
	})

	return URLClickedResponse{
		ClickedByCountry:  clickedByCountry,
		ClickedByCity:     clickedByCity,
		ClickedByDate:     clickedByDate,
		ClickedByPlatform: clickedByPlatform,
		ClickedByReferer:  clickedByReferer,
		TotalClicks:       len(in),
	}

}
