package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// var cfgFile string
const DEFAULTHOST = "localhost"
const DEFAULTWEBPORT = "8080"

// Commodity structs to make simple to unmarshal json responses
type responseTks struct { // Adapted from server.Response
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *struct {
		Type    string
		Content *string
	} `json:"data"`
}

type responseTksArr struct { // Adapted from server.Response
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *struct {
		Type    string
		Content *[]string
	} `json:"data"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "A client for the Tickets server.",
	Long: `
This client connects with the Tickets server and lets you manage any number of tickets. 
* The tipical scenario is using the "New" command to create a new ticket.
* "List" and "Show" commands will let you see the status of any or all tickets available.
* The "Add" command let you add products to a given ticket..`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

var host string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&host, "server", "s", DEFAULTHOST+":"+DEFAULTWEBPORT, "Host:port of the tickets server, ie: localhost:8080")
}

func bindJSON(r io.Reader, target interface{}) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}
