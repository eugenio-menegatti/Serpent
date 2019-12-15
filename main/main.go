/*
SERPENT - a simple program to play a famous game in text mode
Copyright 2019 Eugenio Menegatti
myindievg@gmail.com

	 This file is part of SERPENT.
	 The file COPYING describes the terms under which SERPENT is distributed.

   SERPENT is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   SERPENT is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with SERPENT.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

// GOPATH = %USERPROFILE%\go
import (
	"bufio"
	"fmt"
	"os"
	"serpent/io"
	"serpent/piton"
	"strings"
)

func main() {
	fmt.Println("Severus start")

	io.PressEnter()
	quit := false
	reader := bufio.NewReader(os.Stdin)

	piton.Init()
	defer piton.Close()

	for !quit {

		piton.ClearConsole()

		fmt.Println("Enter one of the following:")
		fmt.Println(" p to play")
		fmt.Println(" q to quit")
		fmt.Println("")
		fmt.Print("Choice? ")

		text, _ := reader.ReadString('\n')
		fmt.Println("Your choice: ", text)
		fmt.Println()

		if strings.Contains(text, "q") || strings.Contains(text, "Q") {
			quit = true
		}
		if strings.Contains(text, "p") || strings.Contains(text, "P") {
			piton.NewGame(nil)
			score := piton.HumanPlay()
			fmt.Print("Game over. Your score is ", score, ".  Press Enter")
			reader.ReadString('\n')
		}
	}
	fmt.Println("Severus end")
}
