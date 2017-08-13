package main

import (
	"encoding/json"
	"fmt"
	"github.com/dyeduguru/wikiracer/server"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	Path      []string `json:"path"`
	TimeTaken float64  `json:"timeTaken"`
}

func main() {
	var startPage, endPage string
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
			resp, err := http.Get(fmt.Sprintf("http://localhost:8080/race/%s/%s", startPage, endPage))
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
				panic(err)
			}
			fmt.Printf("Path: ")
			for i, item := range parsedResponse.Path {
				if i == len(parsedResponse.Path)-1 {
					fmt.Printf("%s", item)
				} else {
					fmt.Printf("%s -> ", item)
				}
			}
			fmt.Println()
			fmt.Printf("Time Taken: %f seconds\n", parsedResponse.TimeTaken)
		},
	}
	cmdRace.Flags().StringVarP(&startPage, "start", "s", "Mike Tyson", "page to start the race")
	cmdRace.Flags().StringVarP(&endPage, "end", "e", "Hangover", "page to end the race")
	var rootCmd = &cobra.Command{Use: "wikiracer"}
	rootCmd.AddCommand(cmdServer, cmdRace)
	rootCmd.Execute()
}
