package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)


func SetUpRouter() *gin.Engine{
    router := gin.Default()
    return router
}

func TestHomepageHandler(t *testing.T) {
    mockResponse := `{"message":"Welcome to the Tech Company listing API with Golang"}`
    r := SetUpRouter()
    r.GET("/", HomepageHandler)
    req, _ := http.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    responseData, _ := io.ReadAll(w.Body)
    assert.Equal(t, mockResponse, string(responseData))
    assert.Equal(t, http.StatusOK, w.Code)
}

func TestNewCompanyHandler(t *testing.T) {
    r := SetUpRouter()
    r.POST("/company", NewCompanyHandler)
    companyId := xid.New().String()
    company := Company{
        ID: companyId,
        Name: "Demo Company",
        CEO: "Demo CEO",
        Revenue: "35 million",
    }
    jsonValue, _ := json.Marshal(company)
    req, _ := http.NewRequest("POST", "/company", bytes.NewBuffer(jsonValue))

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetCompaniesHandler(t *testing.T) {
    r := SetUpRouter()
    r.GET("/companies", GetCompaniesHandler)
    req, _ := http.NewRequest("GET", "/companies", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    var companies []Company
    json.Unmarshal(w.Body.Bytes(), &companies)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.NotEmpty(t, companies)
}

func TestUpdateCompanyHandler(t *testing.T) {
    r := SetUpRouter()
    r.PUT("/company/:id", UpdateCompanyHandler)
    company := Company{
        ID: `2`,
        Name: "Demo Company",
        CEO: "Demo CEO",
        Revenue: "35 million",
    }
    jsonValue, _ := json.Marshal(company)
    reqFound, _ := http.NewRequest("PUT", "/company/"+company.ID, bytes.NewBuffer(jsonValue))
    w := httptest.NewRecorder()
    r.ServeHTTP(w, reqFound)
    assert.Equal(t, http.StatusOK, w.Code)

    reqNotFound, _ := http.NewRequest("PUT", "/company/12", bytes.NewBuffer(jsonValue))
    w = httptest.NewRecorder()
    r.ServeHTTP(w, reqNotFound)
    assert.Equal(t, http.StatusNotFound, w.Code)
}