package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/shlee2227/go_crawling/internal/service/search"
)


type Controller interface {
    RegisterRoutes(r *gin.Engine)
	searchAndStore(ctx *gin.Context)
	getItems(ctx *gin.Context)
}


type controller struct {
    service service.Service
}

func NewNaverController(service service.Service) Controller {
    return &controller{
        service: service,
    }
}

// Gin 라우터 등록
func (c *controller) RegisterRoutes(r *gin.Engine) {
    r.GET("/search", c.searchAndStore)
    r.GET("/items", c.getItems)
}

func (c *controller) searchAndStore(ctx *gin.Context) {
    searchWord := ctx.Query("search_word")
    if searchWord == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "검색어를 입력해주세요"})
        return
    }

    err := c.service.SearchAndStoreItems(searchWord)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, gin.H{"message": "검색 및 저장 성공"})
}

func (c *controller) getItems(ctx *gin.Context) {
    items, err := c.service.GetAllItems()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, items)
}
