package controllers


import (
    "github.com/revel/revel"
    "net/http"
    "io/ioutil"
    "encoding/json"
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
    Url string `json:"url"`
}

func (c App) Index() revel.Result {
    resp, err := http.Get("https://api.github.com/repos/ComputerScienceHouse/Constitution/pulls")
    if err != nil {
        return c.Render()
    }
    defer resp.Body.Close()
    body, err1 := ioutil.ReadAll(resp.Body)
    if err1 != nil {
        return c.Render()
    }
    responseString := []byte(string(body))
    var prs []PullRequest

    err2 := json.Unmarshal(responseString, &prs)
    if err2 != nil {
        return c.Render()
    }

    return c.Render(prs)
}
