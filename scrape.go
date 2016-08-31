package main

import (
  "strings"
  "strconv"
  "regexp"
  "log"
  "fmt"
  "encoding/json"
  "github.com/PuerkitoBio/goquery"
)

func main() {
  doc, err := goquery.NewDocument("http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html")
  if err != nil {
    log.Fatal(err)
  }

  total := 0.0

  type T struct {
    Title       string  `json:"title"`
    UnitPrice   float64 `json:"unit_price"`
    Description string  `json:"description"`
    Size        string  `json:"size"`
  }

  // Find the product items
  var items []T
  doc.Find(".productInner").Each(func(i int, s *goquery.Selection) {
    // For each item found, get the title, price and description
    anchor := s.Find(".productInfo h3 a")
    title := strings.TrimSpace(anchor.Text())
    re := regexp.MustCompile("[0-9]+\\.[0-9]+")
    str := re.FindString(strings.TrimSpace(s.Find(".pricePerUnit").Text()))
    price, _ := strconv.ParseFloat(str, 64)
    description := ""
    size := 0
    // Read in each item's own page too..
    url, use := anchor.Attr("href")
    if use {
      doc2, err := goquery.NewDocument(url)
      if err != nil {
        log.Fatal(err)
      }
      meta := doc2.Find("meta")
      description = meta.AttrOr("content", "")
      html, err := doc2.Html()
      if err != nil {
        log.Fatal(err)
      }
      size = len(html)
    }
    total += price
    item := T {
      Title: title,
      UnitPrice: price,
      Description: description,
      Size: fmt.Sprintf("%.1fkb", float64(size)/1000.0),
    }
    items = append(items, item)
  })

  type Response struct {
    Results []T     `json:"results"`
    Total   float64 `json:"total"`
  }
  data := &Response{Results: items, Total: total}

  b, err := json.MarshalIndent(data, "", "  ")
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println(string(b))
}
