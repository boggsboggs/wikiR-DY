package main

import (
	"github.com/dyeduguru/wikiracer/server"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

func main() {
	var cmdServer = &cobra.Command{
		Use:   "server",
		Short: "Run the server",
		Run: func(cmd *cobra.Command, args []string) {
			log.Fatal(http.ListenAndServe(":8080", server.NewRouter()))
		},
	}
	var rootCmd = &cobra.Command{Use: "wikiracer"}
	rootCmd.AddCommand(cmdServer)
	rootCmd.Execute()
}
