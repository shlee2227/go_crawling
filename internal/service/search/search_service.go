package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	entities "github.com/shlee2227/go_crawling/internal/entities/search"
	repository "github.com/shlee2227/go_crawling/internal/repository/search"
)

// 정의
type Service interface {
	SearchAndStoreItems(searchText string) error
	GetAllItems() ([]entities.Item, error)
}

type service struct {
	repo repository.Repository
}

// 객체 생성자 
func NewService(repo repository.Repository) Service {
	return &service{
		repo: repo,
	}
}

// 함수 구현
func (s *service) SearchAndStoreItems(searchText string) error {
	items, err := s.SearchNaver(searchText)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("naver API 검색 실패")
	}
	err = s.repo.Create(items)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("naver API 검색 결과 저장 실패")
	} else {
		log.Println("Naver API 검색 및 결과 저장 성공")
		return nil
	}
}


// crawlNaver 네이버 api를 통해 검색
func (s *service) SearchNaver(searchText string) ([]entities.Item, error) {
	encodedSearchText := url.QueryEscape(searchText)
	url := fmt.Sprintf("https://openapi.naver.com/v1/search/blog.json?query=%sdisplay=10&start=1&sort=sim", encodedSearchText) // 이 부분을 고루틴으로 여러페이지 한번에 해보기

	// 검색어 request 생성
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("naver API 생성 실패")
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
		log.Printf("StatusCode: %d", resp.StatusCode)
		return nil, fmt.Errorf("naver API 요청 실패: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("naver API 응답 읽기 실패: %d", resp.StatusCode)
	}

	var response entities.NaverResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("naver API 응답 parsing 실패: %d", resp.StatusCode)
	}

	// for _, item := range response.Items {
	// 	fmt.Printf("%+v\n", item)
	// 	fmt.Println()
	// }
	return response.Items, nil
}

// DB에 있는 결과를 조회
func (s *service) GetAllItems() ([]entities.Item, error) {
	return s.repo.GetAll()
}
