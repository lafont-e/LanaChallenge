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

	"github.com/elafont/CbreChallenge/hangman"
	"github.com/elafont/CbreChallenge/server"
	"github.com/spf13/cobra"
)

// guessCmd represents the guess command
var guessCmd = &cobra.Command{
	Use:   "guess",
	Short: "Guess a letter in a hangmam game.",
	Long:  "Guess a letter in a hangmam game.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Guess a letter\n\n")
		guess()
	},
}

var ggame *int
var letter *string

func init() {
	rootCmd.AddCommand(guessCmd)

	letter = guessCmd.Flags().StringP("letter", "l", "", "letter to guess, only 1st character will be used")
	ggame = guessCmd.Flags().IntP("game", "g", 0, "Number of game to use for the guess.")

}

func guess() {
	// check letter is valid
	if letter == nil || len(*letter) < 1 {
		fmt.Println("Error: Can not use empty string as a guess!!")
		return
	}

	l := (*letter)[0]
	if l < 'a' || l > 'z' {
		fmt.Println("Error: Only letters from 'a' to 'z' are allowed")
		return
	}

	hs, err := guessLetter(host, *ggame, string(l))
	if err != nil {
		fmt.Println(err)
		return
	}

	if hs.Done {
		fmt.Print("Greetings!!, Word Discovered....\n\n")
	}

	fmt.Println(hs)

}

func guessLetter(srv string, game int, letter string) (*hangman.Hstatus, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/game/%d/guess/%s", srv, game, letter))
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
		return nil, fmt.Errorf("Error: Can not guess, code:%d, %s", answer.Code, answer.Message)
	}

	return answer.Data.Content, nil
}
