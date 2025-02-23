package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

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
	var wg sync.WaitGroup

	errorchan := make(chan error)
	resultchan := make(chan []entities.Item)

	for i := 1; i <= 100; i += 10 {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			encodedSearchText := url.QueryEscape(searchText)
			url := fmt.Sprintf("https://openapi.naver.com/v1/search/blog.json?query=%s&display=10&start=%d&sort=sim", encodedSearchText, start)

			// 검색어 request 생성
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Println(err)
				errorchan <- fmt.Errorf("naver API 생성 실패: %v", err)
				return
			}

			req.Header.Add("X-Naver-Client-Id", os.Getenv("NAVER_CLIENT_ID"))
			req.Header.Add("X-Naver-Client-Secret", os.Getenv("NAVER_CLIENT_SECRET"))

			// request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				errorchan <- fmt.Errorf("naver API 요청 실패: %v", err)
				return
			}
			defer resp.Body.Close() // 응답 본문 닫아 줘야 함

			if resp.StatusCode != http.StatusOK {
				log.Printf("StatusCode: %d", resp.StatusCode)
				errorchan <- fmt.Errorf("naver API 요청 실패: %d", resp.StatusCode)
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
				errorchan <- fmt.Errorf("naver API 응답 읽기 실패: %v", err)
				return
			}

			var response entities.NaverResponse
			err = json.Unmarshal(body, &response)
			if err != nil {
				log.Println(err)
				errorchan <- fmt.Errorf("naver API 응답 parsing 실패: %v", err)
				return
			}

			resultchan <- response.Items
		}(i)
	}
	// for _, item := range response.Items {
	// 	fmt.Printf("%+v\n", item) //구조체 필드와 값 출력
	// 	fmt.Println()
	// }

	go func() {
		wg.Wait()
		close(resultchan)
		close(errorchan)
	}()

	select {
	case err := <-errorchan:
		fmt.Println("에러가 발생하였습니다", err)
		return nil, err
	default: //에러 없으면 그냥 진행
	}

	var result []entities.Item
	for rs := range resultchan {
		result = append(result, rs...)
	}
	return result, nil
}

// DB에 있는 결과를 조회
func (s *service) GetAllItems() ([]entities.Item, error) {
	return s.repo.GetAll()
}
