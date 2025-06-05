package main

import (
	"fmt"
)

func main() {
	var age int
	var name string

	fmt.Print("กรอกชื่อและอายุ (เว้นด้วยช่องว่าง): ")
	fmt.Scanln(&name, &age)

	fmt.Printf("ชื่อ: %s, อายุ: %d ปี\n", name, age)
}
