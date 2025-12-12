package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/gin-contrib/cors"
)

type Response struct {
    Desk []Item 
    Chair []Item 
    Keyboard []Item 
    Mouse []Item 
    MousePad []Item 
    Monitor []Item 
    MonitorArm []Item 
    DeskLamp []Item 
    PowerStrip []Item 
    PowerStripCase []Item 
    Speaker []Item 
    Microphone []Item 
    Camera []Item 
    Headphone []Item 
}

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
    err := godotenv.Load()
    if err != nil {
        fmt.Println(err)
    }

    router := gin.Default()
    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://127.0.0.1:3000", "http://localhost:3000"}

    router.Use(cors.New(config))
    router.GET("/items", getItems)
    router.Run("localhost:8080")
}

func searchItems(genreId int, keywords string) []Item {
    time.Sleep(200 * time.Millisecond)

    baseURL := "https://app.rakuten.co.jp/services/api/IchibaItem/Search/20220601"
    u, err := url.Parse(baseURL)
    if err != nil {
        fmt.Println(err)
        return []Item{}
    }
    values := u.Query()
    values.Add("keyword", keywords)
    values.Add("format", "json")
    values.Add("sort", "-reviewAverage")
    values.Add("carrier", "0")
    values.Add("availability", "1")
    values.Add("imageFlag", "1")
    values.Add("formatVersion", "2")
    values.Add("applicationId", os.Getenv("APPLICATION_ID"))
    values.Set("genreId", strconv.Itoa(genreId))
    u.RawQuery = values.Encode()

    req, err := http.NewRequest(http.MethodGet, u.String(), nil)
    if err != nil {
        fmt.Println(err)
    }

    client := new(http.Client)
    res, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err)
    }

    var searchResult SearchResult
    if err := json.Unmarshal(body, &searchResult); err != nil {
        fmt.Println(err)
    }

    return searchResult.Items
}

func getItems(c *gin.Context) {
    keywordsArray := c.QueryArray("keyword")
    keywords := strings.Join(keywordsArray, " ")

    var response Response
    response.Desk = searchItems(215698, keywords)
    response.Chair = searchItems(111363, keywords)
    response.Keyboard = searchItems(560088, keywords)
    response.Mouse = searchItems(565170, keywords)
    response.MousePad = searchItems(552391, keywords)
    response.Monitor = searchItems(110105, keywords)
    response.MonitorArm = searchItems(566221, keywords)
    response.DeskLamp = searchItems(500281, keywords)
    response.PowerStrip = searchItems(552481, keywords)
    response.PowerStripCase = searchItems(200166, keywords)
    response.Speaker = searchItems(208316, keywords)
    response.Microphone = searchItems(406336, keywords)
    response.Camera = searchItems(403512, keywords)
    response.Headphone = searchItems(568359, keywords)

    c.IndentedJSON(http.StatusOK, response)
}
