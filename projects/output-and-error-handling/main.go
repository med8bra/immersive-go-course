package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/med8bra/immersive-go-course/projects/output-and-error-handling/client"
)

var (
	weatherURL = flag.String("weather-url", "http://localhost:8080", "weather URL")
)

func main() {
	flag.Parse()

	weatherClient := client.NewWeatherClient(http.DefaultClient, *weatherURL)
	var command string

	for {

		fmt.Print("\n> ")
		if _, err := fmt.Scan(&command); err != nil {
			if err == io.EOF {
				os.Exit(0)
			}
			fmt.Printf("Failed to read the command: %s -- Please try again.", err)
			continue
		}

		switch strings.TrimSpace(command) {
		case "exit":
			fmt.Println("Bye!")
			os.Exit(0)
		default:
			s, err := weatherClient.GetWeather()
			if err != nil {
				fmt.Printf("Failed to get weather: %s -- Please try again.", err)
				continue
			}
			fmt.Printf("\t Weather: %s\n", s)
		}
	}
}
