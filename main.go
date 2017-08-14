package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dyeduguru/wikiracer/server"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	Path      []string `json:"path"`
	TimeTaken float64  `json:"timeTaken"`
	Error     string   `json:"error"`
}

func main() {
	var startPage, endPage string
	var urlFlag bool
	var cmdServer = &cobra.Command{
		Use:   "server",
		Short: "Run the server",
		Run: func(cmd *cobra.Command, args []string) {
			log.Fatal(http.ListenAndServe(":8080", server.NewRouter()))
		},
	}
	var cmdRace = &cobra.Command{
		Use:   "race",
		Short: "Run the race",
		Run: func(cmd *cobra.Command, args []string) {
			start := base64.StdEncoding.EncodeToString([]byte(startPage))
			end := base64.StdEncoding.EncodeToString([]byte(endPage))
			runRace(start, end, urlFlag, cmd.OutOrStdout())
		},
	}
	cmdRace.Flags().StringVarP(&startPage, "start", "s", "Mike Tyson", "page to start the race")
	cmdRace.Flags().StringVarP(&endPage, "end", "e", "Hangover", "page to end the race")
	cmdRace.Flags().BoolVarP(&urlFlag, "url", "u", false, "race with URL")
	var rootCmd = &cobra.Command{Use: "wikiracer"}
	rootCmd.AddCommand(cmdServer, cmdRace)
	rootCmd.Execute()
}

func runRace(startPage, endPage string, urlFlag bool, w io.Writer) {
	var base string
	if urlFlag {
		base = "race/url"
	} else {
		base = "race"
	}
	pageURL := fmt.Sprintf("http://localhost:8080/%s/%s/%s", base, startPage, endPage)
	resp, err := http.Get(pageURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var parsedResponse Response
	if err := json.Unmarshal(body, &parsedResponse); err != nil {
		fmt.Println(string(body))
		panic(err)
	}
	if parsedResponse.Error != "" {
		fmt.Printf("Error: %v", parsedResponse.Error)
		return
	}
	fmt.Fprintf(w, "Path: ")
	for i, item := range parsedResponse.Path {
		if i == len(parsedResponse.Path)-1 {
			fmt.Fprintf(w, "%s", item)
		} else {
			fmt.Fprintf(w, "%s -> ", item)
		}
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Time Taken: %f seconds\n", parsedResponse.TimeTaken)
}
