// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/lafont-e/LanaChallenge/server"
	"github.com/spf13/cobra"
)

// addCmd represents the addProduct command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a product to a ticket in the server.",
	Long:  "add a product to a ticket in the server.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("add a product\n\n")
		add()
	},
}

var tkID *int
var prodCode *string
var prodQty *int

func init() {
	rootCmd.AddCommand(addCmd)

	tkID = addCmd.Flags().IntP("ticket", "t", 0, "Number of ticket to use.")
	prodCode = addCmd.Flags().StringP("product", "p", "", "code of product to add.")
	prodQty = addCmd.Flags().IntP("quantity", "q", 1, "Number of product to add.")

}

func add() {
	// check Code is valid
	if prodCode == nil {
		fmt.Println("Error: Can not use empty string as a product code.")
		return
	}

	tkStatus, err := addProduct(host, *tkID, *prodQty, prodCode)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(*tkStatus)

}

func addProduct(srv string, tkID int, qty int, pcode *string) (*string, error) {
	form := url.Values{}
	form.Set("code", *pcode)
	form.Set("quantity", strconv.Itoa(qty))

	resp, err := http.PostForm(fmt.Sprintf("http://%s/ticket/%d", srv, tkID), form)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var answer responseTks

	if err := bindJSON(bytes.NewReader(body), &answer); err != nil {
		return nil, fmt.Errorf("error reading response %v", err)
	}

	if answer.Status == server.StatusFail {
		return nil, fmt.Errorf("Error: Can not guess, code:%d, %s", answer.Code, answer.Message)
	}

	return answer.Data.Content, nil
}
