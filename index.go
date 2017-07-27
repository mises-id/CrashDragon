package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"git.1750studios.com/GSoC/CrashDragon/database"

	"github.com/gin-gonic/gin"
	colorful "github.com/lucasb-eyer/go-colorful"
)

type chartDataset struct {
	Label           string `json:"label"`
	BackgroundColor string `json:"backgroundColor"`
	Data            []int  `json:"data"`
}

type chartData struct {
	Labels   []string       `json:"labels"`
	Datasets []chartDataset `json:"datasets"`
}

// GetIndex returns index page with stats
func GetIndex(c *gin.Context) {
	type statResult struct {
		Date  string
		Field string
		Count int
	}
	type genericResult struct {
		Result string
	}

	var DateResult []genericResult
	database.Db.Raw("SELECT (generate_series(CURRENT_DATE - '30 days'::interval, CURRENT_DATE, '1 day'::interval))::date::text AS result;").Scan(&DateResult)

	// Version Stats diagram
	var VersionStatResult []statResult
	var VersionResult []genericResult
	var VersionData chartData
	database.Db.Raw("SELECT version AS result FROM reports WHERE (created_at > now() - '30 days'::interval) GROUP BY version;").Scan(&VersionResult)
	database.Db.Raw("SELECT created_at::date::text AS date, version AS field, count(*) FROM (SELECT * FROM reports WHERE created_at > now() - '30 days'::interval) AS dates GROUP BY created_at::date, version ORDER BY created_at::date ASC, version ASC;").Scan(&VersionStatResult)
	for _, row := range DateResult {
		VersionData.Labels = append(VersionData.Labels, row.Result)
	}
	for _, row := range VersionResult {
		var ChartDataset chartDataset
		ChartDataset.Label = row.Result
		sum := sha256.Sum256([]byte(ChartDataset.Label))
		c, _ := colorful.Hex(fmt.Sprintf("#%x", sum[7:10]))
		r, g, b := c.RGB255()
		ChartDataset.BackgroundColor = fmt.Sprintf("rgba(%d,%d,%d,0.2)", r, g, b)
		VersionData.Datasets = append(VersionData.Datasets, ChartDataset)
	}
	for _, date := range VersionData.Labels {
		for i, version := range VersionData.Datasets {
			versionAndDateInRow := false
			for _, row := range VersionStatResult {
				if row.Date == date && row.Field == version.Label {
					versionAndDateInRow = true
					VersionData.Datasets[i].Data = append(version.Data, row.Count)
				}
			}
			if !versionAndDateInRow {
				VersionData.Datasets[i].Data = append(version.Data, 0)
			}
		}
	}

	// Product stats diagram
	var ProductStatResult []statResult
	var ProductResult []genericResult
	var ProductData chartData
	database.Db.Raw("SELECT product AS result FROM reports WHERE (created_at > now() - '30 days'::interval) GROUP BY product;").Scan(&ProductResult)
	database.Db.Raw("SELECT created_at::date::text AS date, product AS field, count(*) FROM (SELECT * FROM reports WHERE created_at > now() - '30 days'::interval) AS dates GROUP BY created_at::date, product ORDER BY created_at::date ASC, product ASC;").Scan(&ProductStatResult)
	for _, row := range DateResult {
		ProductData.Labels = append(ProductData.Labels, row.Result)
	}
	for _, row := range ProductResult {
		var ChartDataset chartDataset
		ChartDataset.Label = row.Result
		sum := sha256.Sum256([]byte(ChartDataset.Label))
		c, _ := colorful.Hex(fmt.Sprintf("#%x", sum[7:10]))
		r, g, b := c.RGB255()
		ChartDataset.BackgroundColor = fmt.Sprintf("rgba(%d,%d,%d,0.2)", r, g, b)
		ProductData.Datasets = append(ProductData.Datasets, ChartDataset)
	}
	for _, date := range ProductData.Labels {
		for i, product := range ProductData.Datasets {
			productAndDateInRow := false
			for _, row := range ProductStatResult {
				if row.Date == date && row.Field == product.Label {
					productAndDateInRow = true
					ProductData.Datasets[i].Data = append(product.Data, row.Count)
				}
			}
			if !productAndDateInRow {
				ProductData.Datasets[i].Data = append(product.Data, 0)
			}
		}
	}

	// Platform stats diagram
	var PlatformStatResult []statResult
	var PlatformResult []genericResult
	var PlatformData chartData
	database.Db.Raw("SELECT os AS result FROM reports WHERE (created_at > now() - '30 days'::interval) GROUP BY os;").Scan(&PlatformResult)
	database.Db.Raw("SELECT created_at::date::text AS date, os AS field, count(*) FROM (SELECT * FROM reports WHERE created_at > now() - '30 days'::interval) AS dates GROUP BY created_at::date, os ORDER BY created_at::date ASC, os ASC;").Scan(&PlatformStatResult)
	for _, row := range DateResult {
		PlatformData.Labels = append(PlatformData.Labels, row.Result)
	}
	for _, row := range PlatformResult {
		var ChartDataset chartDataset
		ChartDataset.Label = row.Result
		sum := sha256.Sum256([]byte(ChartDataset.Label))
		c, _ := colorful.Hex(fmt.Sprintf("#%x", sum[7:10]))
		r, g, b := c.RGB255()
		ChartDataset.BackgroundColor = fmt.Sprintf("rgba(%d,%d,%d,0.2)", r, g, b)
		PlatformData.Datasets = append(PlatformData.Datasets, ChartDataset)
	}
	for _, date := range PlatformData.Labels {
		for i, platform := range PlatformData.Datasets {
			platformAndDateInRow := false
			for _, row := range PlatformStatResult {
				if row.Date == date && row.Field == platform.Label {
					platformAndDateInRow = true
					PlatformData.Datasets[i].Data = append(platform.Data, row.Count)
				}
			}
			if !platformAndDateInRow {
				PlatformData.Datasets[i].Data = append(platform.Data, 0)
			}
		}
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":     "Stats",
		"versions":  VersionData,
		"products":  ProductData,
		"platforms": PlatformData,
	})
}
