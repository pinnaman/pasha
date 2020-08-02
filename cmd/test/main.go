package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/mmcdole/gofeed"
	"github.com/rogpeppe/go-charset/charset"
)

type jsonFeedObject struct {
	Name string
	Url  string
}

type RssFeed struct {
	XMLName xml.Name  `xml:"urlset"`
	News    []NewsMap `xml:"url"`
}

type NewsMap struct {
	XMLName  xml.Name `xml:"url"`
	Link     string   `xml:"loc"`
	Keywords string   `xml:"news>keywords"`
	//Title    string   `xml:"news>title"`
	PubDate     string `xml:"news>publication_date"`
	Description string `xml:"news>title"`
}

type RssFeed2 struct {
	XMLName xml.Name  `xml:"rss"`
	Items   []RssItem `xml:"channel>item"`
}

type RssItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	PubDate     string   `xml:"pubDate"`
	Description string   `xml:"description"`
}

//var m = map[string]interface{}
//var url = "https://www.cnbc.com/sitemap_news.xml"
//var rssFeed = &RssFeed{}

const (
	localFilename = "./config.json"
)

// define global map; initialize as empty with the trailing {}
var (
	rssMap = make(map[string]interface{})
	//powers := make(map[Nums]int)
)

// postgres credentials
const (
	host     = "localhost"
	port     = 54320
	user     = "postgres"
	password = "postgres"
	dbname   = "ddoor_db"
)

func init() {
	rssMap["https://www.cnbc.com/sitemap_news.xml"] = "RssFeed"
	// 'https://news.ycombinator.com/rss'
	// https://edition.cnn.com/sitemaps/news.xml
}

// PgConn return conn handle
func PgConn() *sql.DB {
	var err error
	fmt.Println("connecting...")
	//dbDriver, err = sql.Open("postgres", dataSourceName)
	// lazily open db (doesn't truly open until first request)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	dbDriver, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Panic(err)
	}
	//defer db.Close()

	if err = dbDriver.Ping(); err != nil {
		log.Panic(err)
	}

	fmt.Println("#****Successfully connected*****#")
	return dbDriver
}

func parseXML(xmlDoc []byte, target interface{}) {
	reader := bytes.NewReader(xmlDoc)
	decoder := xml.NewDecoder(reader)
	// Fixes "xml: encoding \"windows-1252\" declared but Decoder.CharsetReader is nil"
	decoder.CharsetReader = charset.NewReader
	if err := decoder.Decode(target); err != nil {
		log.Fatalf("unable to parse XML '%s':\n%s", err, xmlDoc)
	}
}

func fetchURL(url string) []byte {
	/*
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("unable to GET '%s': %s", url, err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
	*/
	b, err := ioutil.ReadFile("../data/cnbc_sitemap_news.xml") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	//fmt.Println(b)   // print the content as 'bytes'
	//str := string(b) // convert content to a 'string'
	fmt.Println("****STRING*******")
	//fmt.Println(str) // print the content as a 'string'
	if err != nil {
		log.Fatalf("unable to read body '%s': %s", url, err)
	}

	return b
}

func getRSSData() {

	//var rssFeed RssFeed
	for k, v := range rssMap {
		fmt.Println("k:", k, "v:", v)
		xmlDoc := fetchURL(k)
		//print(xmlDoc)
		//switch x := v.(type) {
		//case string:
		var rssFeed = &RssFeed{}
		//}

		parseXML(xmlDoc, rssFeed)

		for _, item := range rssFeed.News {
			//log.Printf("%s: %s: %s", item.Title, item.PubDate, item.Link)
			//fmt.Fprintf(w, item.Title, item.PubDate, item.Link, "\n")
			//fmt.Println("%s: %s: %s", item.Title, item.PubDate, item.Link)
			fmt.Println("#*******"+item.Link, item.Keywords, item.Description, item.PubDate+"******#")
		}
	}

}

func getRSSDataAgg() {

	var sStmt string = "insert into finance_news (srce,url,description,pub_date) values ($1, $2, $3, $4)"
	dbConn := PgConn()
	defer dbConn.Close()
	stmt, err := dbConn.Prepare(sStmt)
	if err != nil {
		log.Panic(err)
	}

	fp := gofeed.NewParser()
	//b, err := ioutil.ReadFile(localFilename) // just pass the file name
	file, err := ioutil.ReadFile(localFilename)
	//fmt.Println(b)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("// defining array of struct jsonFeedObject")
	var feeds []jsonFeedObject
	err2 := json.Unmarshal(file, &feeds)
	if err2 != nil {
		fmt.Println("error:", err2)
		os.Exit(1)
	}

	fmt.Println("// loop over array of structs of feedObject")

	for k := range feeds {
		fmt.Printf("The feed '%s' is located at=> '%s'\n", feeds[k].Name, feeds[k].Url)
		//feed, _ := fp.ParseURL("http://feeds.twit.tv/twit.xml")
		feed, _ := fp.ParseURL(feeds[k].Url)
		//fmt.Println(feed.Title + " - " + feed.Description)
		fmt.Println(feed.Title)
		for i := 0; i <= len(feed.Items)-1; i++ {
			// display title content

			fmt.Println(feed.Items[i].Title)
			fmt.Println(feed.Items[i].Description)
			fmt.Println(feed.Items[i].Link)
			fmt.Println(feed.Items[i].Published)

			res, err := stmt.Exec(feeds[k].Name, feed.Items[i].Link, feed.Items[i].Title, feed.Items[i].Published)
			if err != nil || res == nil {
				// do nothing...
				//log.Println(err)
			}

		}

	}
	// close prepared stmt
	stmt.Close()

}

func main() {

	//fmt.Printf("globalMap:%#+v", rssMap)
	//getRSSData()
	getRSSDataAgg()

}
