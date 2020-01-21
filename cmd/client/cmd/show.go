package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/elafont/LanaChallenge/server"

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details of a given ticket.",
	Long:  "Show details of a given ticket.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Showing a ticket\n\n")
		show()
	},
}

var tkid *int

func init() {
	rootCmd.AddCommand(showCmd)

	tkid = showCmd.Flags().IntP("ticket", "t", 0, "Number of ticket to show.")

}

func show() {
	hs, err := getTicket(host, *tkid)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(hs)
}

func getTicket(srv string, tkid int) (*string, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/ticket/%d", srv, tkid))
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
		return nil, fmt.Errorf("Error: Can not show given ticket, code:%d, %s", answer.Code, answer.Message)
	}

	return answer.Data.Content, nil
}
