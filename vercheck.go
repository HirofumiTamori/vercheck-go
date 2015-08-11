package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type RetVal struct {
	Title   string
	Version string
	RelDate string
	Diff    time.Duration
	Index   int
}

type UrlListType struct {
	Url   string
	Type  string
	Index int
}

var UrlList = []UrlListType{
	{Url: "https://github.com/jquery/jquery/releases", Type: "1"},
	{Url: "https://github.com/angular/angular/releases", Type: "1"},
	{Url: "https://github.com/facebook/react/releases", Type: "2"},
	{Url: "https://github.com/PuerkitoBio/goquery/releases", Type: "1"},
	{Url: "https://github.com/revel/revel/releases", Type: "2"},
	{Url: "https://github.com/lhorie/mithril.js/releases", Type: "2"},
	{Url: "https://github.com/muut/riotjs/releases", Type: "1"},
	{Url: "https://github.com/atom/atom/releases", Type: "2"},
	{Url: "https://github.com/Microsoft/TypeScript/releases", Type: "2"},
	{Url: "https://github.com/docker/docker/releases", Type: "1"},
	{Url: "https://github.com/JuliaLang/julia/releases", Type: "2"},
	{Url: "https://github.com/Araq/Nim/releases", Type: "1"},
	{Url: "https://github.com/rust-lang/rust/releases", Type: "1"},
	{Url: "https://github.com/elixir-lang/elixir/releases", Type: "2"},
}

var ShortTimeFormat = "2006/01/02"
var LocationName = "Asia/Tokyo"

func PutAFormattedLine(val RetVal) {
	l := val.Title
	if len(val.Title) < 8 {
		l += "\t"
	}
	l += "\t"
	l += val.Version
	if len(val.Version) < 8 {
		l += "\t"
	}
	l += "\t" + val.RelDate

	// Indacate the caution if the last release was within 2 weeks
	if val.Diff <= time.Hour*24*14 {
		l += "\t<<<<< updated at "
		l += strconv.FormatInt(int64(val.Diff/(time.Hour*24)), 10)
		l += " day(s) ago."
	}

	fmt.Println(l)

}

func GetTitleVerAndRelDate1(url string, index int) RetVal {
	doc, _ := goquery.NewDocument(url)
	var ret RetVal
	ret.Title = doc.Find(".container").Find("strong").First().Text()
	var sv = doc.Find(".tag-name").First()
	var vstr = sv.Text()
	var re = regexp.MustCompile("^[^\\d]+(.*)")
	var replaced = re.ReplaceAllString(vstr, "$1")
	ret.Version = replaced

	var st = doc.Find("span[class=date]").Children().Nodes[0]
	var loc, _ = time.LoadLocation("Asia/Tokyo")
	var tm, _ = time.Parse(time.RFC3339, st.Attr[0].Val)
	tm = tm.In(loc) // change location
	ret.RelDate = strings.Replace(tm.Format(ShortTimeFormat), "/", ".", -1)
	ret.Diff = time.Since(tm)
	ret.Index = index

	return ret
}

func GetTitleVerAndRelDate2(url string, index int) RetVal {
	doc, _ := goquery.NewDocument(url)
	// Get Title
	var ret RetVal
	ret.Title = doc.Find(".container").Find("strong").First().Text()
	// Get Version
	var sv = doc.Find("span[class=css-truncate-target]").First()
	var vstr = sv.Text()
	var re = regexp.MustCompile("^[^\\d]+(.*)")
	var replaced = re.ReplaceAllString(vstr, "$1")
	ret.Version = replaced

	// Get Release Date
	var st = doc.Find("time").Nodes[0]
	var loc, _ = time.LoadLocation("Asia/Tokyo")
	var tm, _ = time.Parse(time.RFC3339, st.Attr[0].Val)
	tm = tm.In(loc) // change location
	ret.RelDate = strings.Replace(tm.Format(ShortTimeFormat), "/", ".", -1)
	ret.Diff = time.Since(tm)
	ret.Index = index

	return ret

}

func main() {
	// Set the max process number
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	ch1 := make(chan RetVal, len(UrlList))

	Results := make([]RetVal, len(UrlList))

	for i, u := range UrlList {
		go func(url string, tp string, index int) {
			switch tp {
			case "1":
				ch1 <- GetTitleVerAndRelDate1(url, index)
			case "2":
				ch1 <- GetTitleVerAndRelDate2(url, index)
			}
		}(u.Url, u.Type, i)
	}

	// Get the result thru a channel and sort them
	for i := 0; i < len(UrlList); i++ {
		var r RetVal
		r = <-ch1
		Results[r.Index] = r
	}

	for _, r := range Results {
		PutAFormattedLine(r)
	}
}
