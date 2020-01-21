package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/elafont/LanaChallenge/server"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new ticket.",
	Long:  "Creates a new ticket.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("New ticket\n\n")
		new()
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func new() {
	tkstatus, err := newTicket(host)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tkstatus)
}

func newTicket(srv string) (*string, error) {
	resp, err := http.Get("http://" + srv + "/newticket")
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var answer responseHs

	if err := bindJSON(bytes.NewReader(body), &answer); err != nil {
		return nil, fmt.Errorf("error reading response %v", err)
	}

	if answer.Status == server.StatusFail {
		return nil, fmt.Errorf("Error: Can not generate a new ticket, code:%d, %s", answer.Code, answer.Message)
	}

	return answer.Data.Content, nil
}
