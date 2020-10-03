package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	c := make(chan int, 15)
	input := bufio.NewScanner(os.Stdin)
	fmt.Println("enter a nunber:\n")
	go func() {
		// 逐行扫描
		for input.Scan() {
			line := input.Text()
			// 输入bye时 结束
			if line == "bye" {
				break
			}
			int, _ := strconv.Atoi(line)
			c <- int
		}
	}()
	for i := range c {
		i2 := i * i
		fmt.Println("--------------->\n", i, " * ", i, "=", i2, "\n enter next:\n")
	}
}
