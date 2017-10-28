package controllers

import (
    "github.com/revel/revel"
    "net/http"
    "io/ioutil"
    "io"
    "encoding/json"
    "fmt"
    "strconv"
    "strings"
    "encoding/csv"
    "bufio"
    "os"
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
    Body string
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

func (c App) Ballots(number int) revel.Result {
    resp, err := http.Get("https://patch-diff.githubusercontent.com/raw/ComputerScienceHouse/Constitution/pull/" +
        strconv.Itoa(number) + ".diff")
    if err != nil {
        fmt.Printf("Error fetching PR diff")
        return c.Render()
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("Error reading response body")
        return c.Render()
    }

    diffString := string(body)
    strings.Replace(diffString, `\n`, "\n", -1)

    pokefile, err := os.Open("pokemon.csv")
    if err != nil {
        fmt.Printf("Error opening pokemon.csv")
        return c.Render()
    }

    r := csv.NewReader(bufio.NewReader(pokefile))
    for i := 0; i < 100; i++{
        pokemon, err := r.Read()
        if err == io.EOF {
            break
        }
    }

    return c.Render(diffString, pokemon)
}
