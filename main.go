package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type BookAPIResponse struct {
	Items []struct {
		VolumeInfo struct {
			Title         string   `json:"title"`
			Subtitle      string   `json:"subtitle"`
			Authors       []string `json:"authors"`
			PublishedDate string   `json:"publishedDate"`
			Description   string   `json:"description"`
			ImageLinks    struct {
				Thumbnail string `json:"thumbnail"`
			} `json:"imageLinks"`
		} `json:"volumeInfo"`
	} `json:"items"`
}

type Book struct {
	Title         string
	Subtitle      string
	Authors       []string
	PublishedDate string
	Description   string
	Thumbnail     string
}

func fetchBook(isbn string) (Book, error) {
	fmt.Printf("Searching... ")
	url := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=isbn:%s&country=JP", isbn)
	resp, err := http.Get(url)
	if err != nil {
		return Book{}, err
	}
	fmt.Printf("Done\n")

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Book{}, err
	}

	fmt.Printf("Formatting...")
	var body BookAPIResponse
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		return Book{}, err
	}

	if len(body.Items) == 0 {
		return Book{}, err
	}

	info := body.Items[0].VolumeInfo
	book := Book{
		Title:         info.Title,
		Subtitle:      info.Subtitle,
		Authors:       info.Authors,
		PublishedDate: info.PublishedDate,
		Description:   info.Description,
		Thumbnail:     info.ImageLinks.Thumbnail,
	}
	fmt.Printf("Done\n")
	return book, nil
}

func main() {
	for {
		fmt.Printf("Input ISBN: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		isbn := scanner.Text()

		book, err := fetchBook(isbn)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", book.Title)
	}
}
