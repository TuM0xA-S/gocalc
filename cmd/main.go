package main

import (
	"bufio"
	"fmt"
	"gorpn"
	"os"
)

func main() {
	scn := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("eval> ")
		if !scn.Scan() {
			break
		}
		expr := scn.Text()
		res, err := gorpn.CalculateExpression(expr)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Printf("= %.3f\n", res)
		}
	}
}
