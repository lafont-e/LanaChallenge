package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/elafont/LanaChallenge/server"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available tickets",
	Long:  `List all available tickets.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("List of tickets\n\n")
		list()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func list() {
	tkstatus, err := listTicket(host)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, tk := range *tkstatus {
		fmt.Println(tk)
	}

}

func listTicket(srv string) (*[]string, error) {
	resp, err := http.Get("http://" + srv + "/tickets")
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var answer responseTksArr

	if err := bindJSON(bytes.NewReader(body), &answer); err != nil {
		return nil, fmt.Errorf("error reading response %v", err)
	}

	if answer.Status == server.StatusFail {
		return nil, fmt.Errorf("Error: Can not list available tickets, code:%d, %s", answer.Code, answer.Message)
	}

	return answer.Data.Content, nil
}
