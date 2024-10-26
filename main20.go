package main

import "fmt"

func main() {

}
func PrettyArrayOutput(input [9]string) {
	for i := 0; i < 7; i++ {
		fmt.Println(input[i+1], "я уже сделал:", input[i])
	}
	fmt.Println(8, "не успел сделать:", input[7])
	fmt.Println(9, "не успел сделать:", input[8])
}
