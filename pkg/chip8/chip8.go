package chip8

import (
	"io"
	"math/rand"
	"os"
	"sync"
)

var (
	emu  *Chip8
	once sync.Once
)

const START_ADDRESS int = 0x200

type Chip8 struct {
	registers  [16]uint8
	memory     [4096]uint8
	index      uint16
	pc         uint16
	stack      [16]uint16
	sp         uint8
	delayTimer uint8
	soundTimer uint8
	keypad     [16]uint8
	video      [64 * 32]uint32
	opcode     uint16
}

func GetInstance() *Chip8 {
	once.Do(func() {
		emu = &Chip8{}
		emu.initialize()
	})
	return emu
}

func (c *Chip8) initialize() {
	c.pc = uint16(START_ADDRESS)
	c.loadFonts()
}

func (c *Chip8) LoadROM(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}
	size := info.Size()

	buffer := make([]byte, size)
	_, err = io.ReadFull(file, buffer)
	if err != nil {
		return err
	}

	for i := 0; i < int(size); i++ {
		c.memory[START_ADDRESS+i] = buffer[i]
	}

	return nil
}

func (chip8 *Chip8) loadFonts() {
	fonts := [80]byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	for i := 0; i < len(fonts); i++ {
		chip8.memory[0x50+i] = fonts[i]
	}
}

func (c *Chip8) RandomByte() uint8 {
	return uint8(rand.Intn(0xFF))
}

func (c *Chip8) op00E0() {

}

func (c *Chip8) op00EE() {

}

func (c *Chip8) op1nnn() {

}

func (c *Chip8) op2nnn() {

}

func (c *Chip8) op3xkk() {

}

func (c *Chip8) Run() {
	c.opcode = (uint16(c.memory[c.pc]) << 8) | uint16(c.memory[c.pc+1])
	c.pc += 2

	switch c.opcode {
	case 0x00E0:
	case 0x00EE:
	}

	if c.delayTimer > 0 {
		c.delayTimer--
	}

	if c.soundTimer > 0 {
		c.soundTimer--
	}

}
