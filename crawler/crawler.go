package crawler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"go_crawling/db"
)

type NaverResponse struct {
	LastBuildDate 	string    `json:"lastBuildDate"`
	Total 			int       `json:"total"`
	Start 			int       `json:"start"`
	Display 		int       `json:"display"`
	Items 			[]Item    `json:"items"`
}

type Item struct {
	Title 		string    `json:"title"`
	Link 		string    `json:"link"`
	Description string    `json:"description"`
	BloggerName string    `json:"bloggername"`
	BloggerLink string    `json:"bloggerlink"`
	PostDate 	string `json:"postdate"`
}


// 검색 및 검색 결과 저장
func SearchAndStoreData(searchText string) error {
	items, err := SearchNaver(searchText)
	if err != nil {
		log.Println("네이버 검색 실패")
		return fmt.Errorf("검색 실패")
	}
	err = StoreNaver(items)
	if err != nil {
		log.Println("검색 결과 저장 실패")
		return fmt.Errorf("저장 실패")
	} else {
		log.Println("검색 및 결과 저장 성공")
		return nil
	}
}

// crawlNaver 네이버 api를 통해 검색
func SearchNaver(searchText string) ([]Item, error) {
	encodedSearchText := url.QueryEscape(searchText)
	url :=fmt.Sprintf("https://openapi.naver.com/v1/search/blog.json?query=%sdisplay=10&start=1&sort=sim", encodedSearchText)
	
	// 검색어 request 생성 
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Add("X-Naver-Client-Id", os.Getenv("NAVER_CLIENT_ID"))
	req.Header.Add("X-Naver-Client-Secret", os.Getenv("NAVER_CLIENT_SECRET"))
	
	// request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("네이버 요청 실패")
		return nil, fmt.Errorf("네이버 요청 실패: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	
	var response NaverResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
	}

	for _, item := range(response.Items){

		fmt.Printf("%+v\n", item)
		fmt.Println()
	}
	
	return response.Items, nil
}

// StoreNaver 네이버 데이터를 DB에 저장
func StoreNaver(items []Item) error {
	q := `INSERT INTO naver_items (title, link, description, blogger_name, blogger_link, post_date) 
	VALUES ($1, $2, $3, $4, $5, $6);`

	for _, item := range(items){
		_, err := db.Conn.Exec(q, item.Title, item.Link, item.Description, item.BloggerName, item.BloggerLink, item.PostDate)
		if err != nil {
			log.Printf("DB 입력 실패: %s", err)
			return err
		}
		log.Println("DB 입력 성공")
	}

	return nil
}

// DB에 있는 결과를 조회
func GetData(){}