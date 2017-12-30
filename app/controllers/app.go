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

func (c App) CustomBallots(prompt string, answers string) revel.Result {
    pokemons := getPokemon(100)
    options := strings.Split(answers, "\n")
    return c.Render(prompt, options, pokemons)
}

func (c App) Ballots(prnumber int, numballots int) revel.Result {
    resp, err := http.Get("https://patch-diff.githubusercontent.com/raw/ComputerScienceHouse/Constitution/pull/" +
        strconv.Itoa(prnumber) + ".diff")
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

    resp, err = http.Get("https://api.github.com/repos/ComputerScienceHouse/Constitution/pulls/" + strconv.Itoa(prnumber) + ".diff")
    if err != nil {
        fmt.Printf("Error fetching PR title")
        return c.Render()
    }

    defer resp.Body.Close()
    titleBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("Error reading response body")
        return c.Render()
    }

    reqString := []byte(string(titleBody))
    var pr PullRequest
    err = json.Unmarshal(reqString, &pr)
    if err != nil {
        fmt.Printf("Error parsing json")
        return c.Render()
    }

    pokemons := getPokemon(100)
    return c.Render(diffString, pokemons, pr)
}

func getPokemon(numballots int) []string {
    pokefile, err := os.Open(os.Getenv("PCSV_PATH"))
    if err != nil {
        fmt.Printf("Error opening pokemon.csv")
        return nil
    }

    r := csv.NewReader(bufio.NewReader(pokefile))
    numballots = numballots + 1
    pokemons := make([]string, numballots)
    for i := 1; i < numballots; i++{
        pokemon, err := r.Read()
        if err == io.EOF {
            break
        }
        pokemons[i] = pokemon[1]
    }
    pokefile.Close()
    return pokemons
}
