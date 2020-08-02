package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	///"text/template"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rogpeppe/go-charset/charset"
)

// create array of maps
// M is an alias for map[string]interface{}
type M map[string]interface{}

type RssFeed struct {
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

type RssFeed2 struct {
	XMLName xml.Name   `xml:"urlset"`
	Items   []RssItem2 `xml:"url"`
}

//Keywords    string   `xml:"news>keywords"`
type RssItem2 struct {
	XMLName     xml.Name `xml:"url"`
	Title       string   `xml:"news>title"`
	Link        string   `xml:"loc"`
	PubDate     string   `xml:"news>publication_date"`
	Description string   `xml:"news>title"`
}

type NewsAggPage struct {
	Title string
	News  []NewsMap
}

type NewsMap struct {
	Srce        string
	Url         string
	Description string
	PubDate     string
}

// postgres credentials
const (
	host     = "127.0.0.1"
	port     = 54320
	user     = "postgres"
	password = "postgres"
	dbname   = "ddoor_db"
)

// relative include path, converted to absolute path by packr
var templatesBox = packr.New("Templates", "./templates")
var assetBox = packr.New("Assets", "./assets")
var dataBox = packr.New("Data", "./data")
var url = "https://www.washingtonpost.com/news-sitemaps/business.xml"
var dbDriver *sql.DB

func logRequest(req *http.Request) {
	now := time.Now()
	log.Printf("%s - %s [%s] \"%s %s %s\" ",
		req.RemoteAddr,
		"",
		now.Format("02/Jan/2007:15:04:05 -0700"),
		req.Method,
		req.URL.RequestURI(),
		req.Proto)
}

func PgConn() *sql.DB {
	var err error
	fmt.Println("#*******Connecting to DB******#")
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

	//fmt.Println("#****Successfully connected*****#")
	return dbDriver
}

func testDBConn(conn *sql.DB) {

	fmt.Println("Check News Table=>")
	rows, err := conn.Query("SELECT description, pub_date FROM news order by pub_date desc LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println("DB Data=>")
	for rows.Next() {
		var desc, pub_date string
		if err := rows.Scan(&desc, &pub_date); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("%s %s\n", desc, pub_date)
		log.Println(desc, pub_date)
		///w.Write([]byte(fmt.Sprintf("Top News Item, %s, %s\n", desc, pub_date)))
	}
}

func fetchURL(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("unable to GET '%s': %s", url, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("unable to read body '%s': %s", url, err)
	}
	return body
}

func parseXML(xmlDoc []byte, target interface{}) {
	reader := bytes.NewReader(xmlDoc)
	decoder := xml.NewDecoder(reader)
	// Fixes "xml: encoding \"windows-1252\" declared but Decoder.CharsetReader is nil"
	decoder.CharsetReader = charset.NewReader
	if err := decoder.Decode(target); err != nil {
		log.Fatalf("unable to parse XML '%s':\n%s", err, xmlDoc)
		fmt.Println("unable to parse XML '%s':\n%s", err)
	}
}

func getNewsData(srce string, url string, tname string, rss *RssFeed) {
	//func getNewsData(srce string, url string, tname string, rss interface{}) {

	//if srce == "ycombinator" {
	//var rssFeed = &RssFeed{}
	xmlDoc := fetchURL(url)
	//parseXML(xmlDoc, &rssFeed)
	parseXML(xmlDoc, rss)
	var sStmt string = "insert into finance_news (srce,url,description,pub_date) values ($1, $2, $3, $4)"
	dbConn := PgConn()
	defer dbConn.Close()
	stmt, err := dbConn.Prepare(sStmt)
	if err != nil {
		log.Panic(err)
	}
	//}

	for _, item := range rss.Items {
		//log.Printf("%s: %s: %s", item.Title, item.PubDate, item.Link)
		//fmt.Fprintf(w, item.Title, item.PubDate, item.Link, "\n")
		res, err := stmt.Exec(srce, item.Link, item.Title, item.PubDate)
		if err != nil || res == nil {
			// do nothing...
			//log.Println(err)
		}
	}
	// close prepared stmt
	stmt.Close()

}

func getNewsSrce(conn *sql.DB) {

	// Get News Sources...
	//rows, err := conn.Query("SELECT srce, cat, sub_cat, url, tab_name FROM news_srce where srce = 'ycombinator'")
	rows, err := conn.Query("SELECT srce, cat, sub_cat, url, tab_name FROM news_srce where srce in ('ycombinator','bbc','yahoo')")

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	defer rows.Close()
	//fmt.Println("DB Data=>")
	for rows.Next() {
		var srce, cat, sub_cat, url, tab_name string
		if err := rows.Scan(&srce, &cat, &sub_cat, &url, &tab_name); err != nil {
			log.Fatal(err)
		}
		//fmt.Println("%s %s\n", srce, cat, sub_cat, url, tab_name)
		//log.Println("DB SRCE ROW=>", srce, cat, sub_cat, url, tab_name)
		if srce == "ycombinator" || srce == "bbc" || srce == "yahoo" {
			getNewsData(srce, url, tab_name, &RssFeed{})
		} //else if srce == "cnbc" {
		//	getNewsData(srce, url, tab_name, &RssFeed2{})
		//}
	}
}

func getNews() {
	for {

		// Get News from sources
		log.Println("Activating news service at=>", time.Now())
		//fmt.Println("#******Activating news service..*******#")
		//log.Println(time.Now().UTC())
		//log.Println(time.Now())
		//rssHandler()
		getNewsSrce(PgConn())
		//time.Sleep(100000 * time.Millisecond)
		time.Sleep(3600 * time.Second)
	}
}

func newsHandler(w http.ResponseWriter, r *http.Request) {

	var ps []NewsMap
	dbConn := PgConn()
	defer dbConn.Close()
	// Get News Sources...
	rows, err := dbConn.Query("SELECT srce, url, description, pub_date FROM finance_news order by pub_date desc")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	//fmt.Println("DB Data=>")
	var srce, url, description, pubDate string
	for rows.Next() {
		if err := rows.Scan(&srce, &url, &description, &pubDate); err != nil {
			log.Fatal(err)
		}
		//fmt.Println("%s %s\n", srce, url, description, pubDate)
		//log.Println("DB NEWS ROW=>", srce, url, description, pubDate)
		// push data template
		ps = append(ps, NewsMap{Srce: srce, Url: url, Description: description, PubDate: pubDate})
	}
	p := NewsAggPage{Title: "News Aggregator", News: ps}
	log.Println("Inside News Handler2")
	templateHome, err := templatesBox.FindString("news.html")
	if err != nil {
		fmt.Println("Template Parsing Error")
		log.Fatal(err)
	}
	t := template.New("")
	t.Parse(templateHome)
	t.Execute(w, p)

}

func helloHandler(res http.ResponseWriter, req *http.Request) {
	//fmt.Fprintf(res, "<img src='assets/gopher.jpeg' alt='gopher' style='width:235px;height:320px;'>")
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	io.WriteString(
		res,
		`<doctype html>
			<html>
				<head>
					<title>Hello Gopher</title>
				</head>
				<body>
					Hello Gopher!!!!!! </br>
					It is really awesome that both Docker and Kubernetes are written in Go!
				</body>
			</html>`,
	)
}
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Go web app powered by Docker")
}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "Ninja!!!"
	}
	log.Printf("Received request for %s\n", name)
	w.Write([]byte(fmt.Sprintf("Hello, %s\n", name)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Health Status=> %d\n", http.StatusOK)))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Root Handler!!!!")
	fmt.Fprintf(w, "<h1>Whoa, Go is neat!</h1>")
	w.Write([]byte(string(9786)))
	fmt.Fprintln(w, "Smiley!!!")
	w.Write([]byte("Docker Gorilla!\n"))
	dbConn := PgConn()
	fmt.Println("Testing DB Connection")
	testDBConn(dbConn)
	defer dbConn.Close()
	//fmt.Println(" DB Connection Test Done")
	fmt.Fprintf(w, "<h4>DB Connection Test Successful!!!</h4>")

}

func main() {

	dbConn := PgConn()
	fmt.Println("Testing DB Connection")
	testDBConn(dbConn)

	// Create Server and Route Handlers
	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/health", healthHandler)
	///r.HandleFunc("/readiness", readinessHandler)
	r.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	r.HandleFunc("/hello", helloHandler)
	r.HandleFunc("/news", newsHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":9000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
