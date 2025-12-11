package main

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
    "net/http"
    "github.com/gin-gonic/gin"
)

type SearchResult struct {
    GenreInformation []int `json:"GenreInformation"`
    Items []Items `json:"Items"`
    TagInformation []int `json:"TagInformation"`
    Carrier int `json:"carrier"`
    Count int `json:"count"`
    First int `json:"first"`
    Hits int `json:"hits"`
    Last int `json:"last"`
    Page int `json:"page"`
    PageCount int `json:"pageCount"`
}

type Items struct {
    Item Item `json:"Item"`
}

type Item struct {
    Availability int `json:"availability"`
    Catchcopy string `json:"catchcopy"`
    ItemCaption string `json:"itemCaption"`
    ItemCode string `json:"itemCode"`
    ItemName string `json:"itemName"`
    ItemPrice int `json:"itemPrice"`
    ItemUrl string `json:"itemUrl"`
    MediumImageUrls []ImageUrl `json:"mediumImageUrls"`
    ReviewAverage int `json:"reviewAverage"`
    ReviewCount int `json:"reviewCount"`
    ShopName string `json:"shopName"`
    ShopCode string `json:"shopCode"`
}

type ImageUrl struct {
    ImageUrl string `json:"imageUrl"`
}

func main() {
    router := gin.Default()
    router.GET("/items/all", getItems)
    router.GET("/items", getItemsByColor)
    router.Run("localhost:8080")
}

func getSearchResult() ([]Items, error) {
    data, err := os.ReadFile("test_data.json")
    if err != nil {
        return nil, fmt.Errorf("fail to read file: %w", err)
    }

    var searchResult SearchResult
    if err := json.Unmarshal(data, &searchResult); err != nil {
        return nil, fmt.Errorf("fail to unmarshal json: %w", err)
    }

    return searchResult.Items, nil
}

func getItems(c *gin.Context) {
    result, err := getSearchResult()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.IndentedJSON(http.StatusOK, result)
}

func itemMatchesColor(item Item, color string) bool {
    return strings.Contains(item.Catchcopy, color) ||
           strings.Contains(item.ItemCaption, color) ||
           strings.Contains(item.ItemName, color)
}

func getItemsByColor(c *gin.Context) {
    result, err := getSearchResult()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    colors := c.QueryArray("color")
    newItems := make([]Items, 0, len(result))

    for _, r := range result {
        for _, c := range colors {
            if itemMatchesColor(r.Item, c) {
                newItems = append(newItems, r)
                break
            }
        }
    }
    if len(newItems) == 0 {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "item not found"})
        return
    }

    c.IndentedJSON(http.StatusOK, newItems)
}
