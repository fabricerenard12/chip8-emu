package main

import (
	"chip8-emu/pkg/chip8"
	"fmt"
)

func main() {
	var emulator *chip8.Chip8 = chip8.GetInstance()
	fmt.Println(emulator.RandomByte())
	fmt.Println(emulator.RandomByte())
	emulator.Run()
}
