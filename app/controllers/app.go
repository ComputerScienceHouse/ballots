package controllers

import (
    "github.com/revel/revel"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "fmt"
)

type App struct {
    *revel.Controller
}

type PullRequest struct {
    Title string
    Number int
    Html_url string
    Diff_url string
    User User `json:"user"`
}

type User struct {
    Login string `json:"login"`
    Html_url string `json:"html_url"`
    Avatar_url string `json:"avatar_url"`
}

func (c App) Index() revel.Result {
    resp, err := http.Get("https://api.github.com/repos/ComputerScienceHouse/Constitution/pulls")
    if err != nil {
        fmt.Printf("Error fetching github information")
        return c.Render()
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("Error reading response body")
        return c.Render()
    }
    responseString := []byte(string(body))
    var prs []PullRequest

    err = json.Unmarshal(responseString, &prs)
    if err != nil {
        fmt.Printf("Error parsing json")
        return c.Render()
    }

    return c.Render(prs)
}
