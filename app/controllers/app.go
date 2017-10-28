package controllers


import (
    "github.com/revel/revel"
    "net/http"
    "io/ioutil"
)

type App struct {
    *revel.Controller
}

func (c App) Index() revel.Result {
    resp, err := http.Get("https://api.github.com/repos/ComputerScienceHouse/Constitution/pulls")
    if err != nil {
        return c.Render()
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    responseString := string(body)
    return c.Render(responseString)
}
