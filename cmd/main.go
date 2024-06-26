package main

import (
	"fmt"
	"github.com/JaKu01/LocalMail"
)

func main() {
	fmt.Println("Hello, World!")

	LocalMail.StartServer(true, "test", "test")
}
