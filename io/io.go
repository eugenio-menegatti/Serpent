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
 
package io

import (
	"bufio"
	"fmt"
	"os"
)

/*PressEnter waits user to press enter*/
func PressEnter() {
	fmt.Println("Press enter...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
