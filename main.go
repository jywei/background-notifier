package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"
)

func main() {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	for {
		fmt.Println(">>Please enter some messages for notification Enter it or \"q\" to exit")
		bio := bufio.NewReader(os.Stdin)
		line, _, _ := bio.ReadLine()

		if string(line) == "q" {
			break
		}

		c.Do("SET", "notification", line)
		world, err := redis.String(c.Do("GET", "notification"))
		if err != nil {
			fmt.Println("key not found")
		}
		fmt.Println(world)
	}
}
