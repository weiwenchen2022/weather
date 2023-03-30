package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Conditions struct {
	Summary            string
	TemperatureCelsius float64
}

func (c Conditions) String() string {
	return fmt.Sprintf("%s %.1fºC", c.Summary, c.TemperatureCelsius+0.05)
}

type Client struct {
	token string
}

func NewClient(token string) *Client {
	return &Client{token: token}
}

func (c *Client) Current(location string) (Conditions, error) {
	return Current(location, c.token)
}

func FormatURL(location, token string) string {
	return fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&APPID=%s",
		location, token)
}

func ParseJSON(r io.Reader) Conditions {
	var temp struct {
		Weather []struct {
			Main string `json:"main"`
		} `json:"weather"`

		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
	}

	if err := json.NewDecoder(r).Decode(&temp); err != nil {
		log.Fatal(err)
	}

	summary := temp.Weather[0].Main
	t := temp.Main.Temp - 273.15
	t, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", t), 64)
	return Conditions{summary, t}
}

func Current(location, token string) (Conditions, error) {
	resp, err := http.Get(FormatURL(location, token))
	if err != nil {
		return Conditions{}, err
	}
	defer resp.Body.Close()

	return ParseJSON(resp.Body), nil
}

func LocationFromArgs(location []string) (string, error) {
	if len(location) == 0 {
		return "", errors.New("no location")
	}

	return strings.Join(location, ""), nil
}

func RunCLI() {
	location, err := LocationFromArgs(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	token := os.Getenv("OPENWEATHER_API_TOKEN")
	if token == "" {
		log.Fatal("No valid API key in the environment variable OPENWEATHER_API_TOKEN")
	}

	client := NewClient(token)

	for i := 0; i < 5; i++ {
		cond, err := client.Current(location)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(cond)

		time.Sleep(1 * time.Second)
	}
}
