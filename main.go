package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type SearchResult struct {
	GenreInformation []int  `json:"GenreInformation"`
	Items            []Item `json:"Items"`
	TagInformation   []int  `json:"TagInformation"`
	Carrier          int    `json:"carrier"`
	Count            int    `json:"count"`
	First            int    `json:"first"`
	Hits             int    `json:"hits"`
	Last             int    `json:"last"`
	Page             int    `json:"page"`
	PageCount        int    `json:"pageCount"`
}

type Item struct {
	Availability    int      `json:"availability"`
	Catchcopy       string   `json:"catchcopy"`
	ItemCaption     string   `json:"itemCaption"`
	ItemCode        string   `json:"itemCode"`
	ItemName        string   `json:"itemName"`
	ItemPrice       int      `json:"itemPrice"`
	ItemUrl         string   `json:"itemUrl"`
	MediumImageUrls []string `json:"mediumImageUrls"`
	ReviewAverage   float64  `json:"reviewAverage"`
	ReviewCount     int      `json:"reviewCount"`
	ShopName        string   `json:"shopName"`
	ShopCode        string   `json:"shopCode"`
}

func main() {
	router := gin.Default()
	router.GET("/items", getItems)
	router.Run("localhost:8080")
}

func searchDesk(keywords string) ([]Item, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("fail to load .env file: %w", err)
	}

	baseURL := "https://app.rakuten.co.jp/services/api/IchibaItem/Search/20220601"
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("fail to parse url: %w", err)
	}
	values := u.Query()
	values.Add("format", "json")
	values.Add("sort", "-reviewAverage")
	values.Add("carrier", "0")
	values.Add("availability", "1")
	values.Add("imageFlag", "1")
	values.Add("formatVersion", "2")
	values.Add("applicationId", os.Getenv("APPLICATION_ID"))
	values.Add("keyword", keywords)

	deskGenreId := [3]int{215698, 215702, 215706}

	items := []Item{}

	for i := 0; i < len(deskGenreId); i++ {
		values.Set("genreId", strconv.Itoa(deskGenreId[i]))
		u.RawQuery = values.Encode()

		req, _ := http.NewRequest(http.MethodGet, u.String(), nil)

		client := new(http.Client)
		res, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fail to access Rakuten api: %w", err)
		}
		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)

		var searchResult SearchResult
		if err := json.Unmarshal(body, &searchResult); err != nil {
			return nil, fmt.Errorf("fail to unmarshal json: %w", err)
		}
		items = append(items, searchResult.Items...)
	}

	return items, nil
}

func getItems(c *gin.Context) {
	keywords := c.Query("keyword")

	result, err := searchDesk(keywords)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(result) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "item not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}
