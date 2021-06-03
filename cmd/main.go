package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/TuM0xA-S/gocalc"
)

func main() {
	scn := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("eval> ")
		if !scn.Scan() {
			break
		}
		expr := scn.Text()
		res, err := gocalc.CalculateExpression(expr)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Printf("= %.3f\n", res)
		}
	}
}
