package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

type CommonResult struct {
	ErrorCode    int         `json:"ErrorCode"`
	ErrorMessage string      `json:"ErrorMessage"`
	Data         interface{} `json:"Data"`
}

type SearchResult struct {
	Word        string `json:"Word"`
	Translation string `json:"Translation"`
}

func (app *application) searchHome(w http.ResponseWriter, r *http.Request) {
	commonResult := CommonResult{
		ErrorCode:    1,
		ErrorMessage: "请输入要查询的单词",
		Data:         nil,
	}
	resultJsonData, err := json.Marshal(commonResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultJsonData)
}

func (app *application) searchWord(w http.ResponseWriter, r *http.Request) {
	lastIndexSlash := strings.LastIndex(r.URL.Path, "/")
	wordEscape := r.URL.Path[lastIndexSlash+1:]
	wordUnescape, err := url.QueryUnescape(wordEscape)

	var commonResult CommonResult

	if err != nil {
		commonResult = CommonResult{
			ErrorCode:    1,
			ErrorMessage: err.Error(),
			Data:         struct{}{},
		}
	} else {
		translationEscape, err := app.selectWordFromDatabase(wordUnescape)
		var data interface{}
		errorCode := 0
		errorMessage := ""
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				errorCode = 2
				errorMessage = "没有查到这个单词"
			}
			data = struct{}{}
		} else {
			data = SearchResult{
				Word:        wordEscape,
				Translation: translationEscape,
			}
		}
		commonResult = CommonResult{
			ErrorCode:    errorCode,
			ErrorMessage: errorMessage,
			Data:         data,
		}
	}
	resultJsonData, err := json.Marshal(commonResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultJsonData)
}

func (app *application) selectWordFromDatabase(wordUnescape string) (string, error) {
	var translation string
	var err error
	if len(wordUnescape) == 0 {
		return "", nil
	}
	err = app.getStatement(wordUnescape[0]).QueryRow(wordUnescape).Scan(&translation)
	if err != nil {
		return "", err
	}
	return translation, err
}

func (app *application) getStatement(firstLetter uint8) *sql.Stmt {
	if (64 < firstLetter && firstLetter < 75) || (96 < firstLetter && firstLetter < 107) {
		return app.searchStmtAK
	} else if (74 < firstLetter && firstLetter < 91) || (106 < firstLetter && firstLetter < 123) {
		return app.searchStmtLZ
	}
	return nil
}
