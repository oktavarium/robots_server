package main

import "github.com/oktavarium/sgs/internal/server"

func main() {
	if err := server.Run(); err != nil {
		panic(err)
	}
}
