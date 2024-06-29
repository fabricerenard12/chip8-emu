package chip8

import (
	"errors"
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
const FONTSET_START_ADDRESS int = 0x50
const VIDEO_WIDTH uint8 = 64
const VIDEO_HEIGHT uint8 = 32

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

func (c *Chip8) GetVideo() [64 * 32]uint32 {
	return c.video
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
		chip8.memory[FONTSET_START_ADDRESS+i] = fonts[i]
	}
}

func (c *Chip8) RandomByte() uint8 {
	return uint8(rand.Intn(0x100))
}

func (c *Chip8) op00E0() {
	for i := range c.video {
		c.video[i] = 0
	}
}

func (c *Chip8) op00EE() {
	c.pc = c.stack[c.sp]
	c.sp--
}

func (c *Chip8) op1nnn() {
	c.pc = c.opcode & 0x0FFF
}

func (c *Chip8) op2nnn() {
	c.sp++
	c.stack[c.sp] = c.pc
	c.pc = c.opcode & 0x0FFF
}

func (c *Chip8) op3xkk() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	var kk uint8 = uint8(c.opcode & 0x00FF)
	if c.registers[x] == kk {
		c.pc += 2
	}
}

func (c *Chip8) op4xkk() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	var kk uint8 = uint8(c.opcode & 0x00FF)
	if c.registers[x] != kk {
		c.pc += 2
	}
}

func (c *Chip8) op5xy0() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	var y uint8 = uint8((c.opcode & 0x00F0) >> 4)
	if c.registers[x] == c.registers[y] {
		c.pc += 2
	}
}

func (c *Chip8) op6xkk() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	var kk uint8 = uint8(c.opcode & 0x00FF)
	c.registers[x] = kk
}

func (c *Chip8) op7xkk() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	var kk uint8 = uint8(c.opcode & 0x00FF)
	c.registers[x] += kk
}

func (c *Chip8) op9xy0() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	var y uint8 = uint8((c.opcode & 0x00F0) >> 4)
	if c.registers[x] != c.registers[y] {
		c.pc += 2
	}
}

func (c *Chip8) opAnnn() {
	c.index = c.opcode & 0x0FFF
}

func (c *Chip8) opBnnn() {
	c.pc = uint16(c.registers[0]) + (c.opcode & 0x0FFF)
}

func (c *Chip8) opCxkk() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	var kk uint8 = uint8(c.opcode & 0x00FF)
	c.registers[x] = c.RandomByte() & kk
}

func (c *Chip8) opDxyn() {
	var x uint8 = c.registers[(c.opcode&0x0F00)>>8]
	var y uint8 = c.registers[(c.opcode&0x00F0)>>4]
	var height uint16 = c.opcode & 0x000F
	c.registers[0x0F] = 0

	for yline := uint16(0); yline < height; yline++ {
		var pixel uint8 = c.memory[c.index+yline]
		for xline := uint16(0); xline < 8; xline++ {
			if (pixel & (0x80 >> xline)) != 0 {
				var xCoord uint8 = (x + uint8(xline)) % VIDEO_WIDTH
				var yCoord uint8 = (y + uint8(yline)) % VIDEO_HEIGHT

				if c.video[xCoord+(yCoord*VIDEO_WIDTH)] == 1 {
					c.registers[0xF] = 1
				}

				c.video[xCoord+(yCoord*VIDEO_WIDTH)] ^= 1
			}
		}
	}
}

func (c *Chip8) op8xy0() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	var y uint8 = uint8((c.opcode & 0x00F0) >> 4)
	c.registers[x] = c.registers[y]
}

func (c *Chip8) op8xy1() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	var y uint8 = uint8((c.opcode & 0x00F0) >> 4)
	c.registers[x] |= c.registers[y]
}

func (c *Chip8) op8xy2() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	var y uint8 = uint8((c.opcode & 0x00F0) >> 4)
	c.registers[x] &= c.registers[y]
}

func (c *Chip8) op8xy3() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	var y uint8 = uint8((c.opcode & 0x00F0) >> 4)
	c.registers[x] ^= c.registers[y]
}

func (c *Chip8) op8xy4() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	var y uint8 = uint8((c.opcode & 0x00F0) >> 4)
	var sum uint16 = uint16(c.registers[x]) + uint16(c.registers[y])
	c.registers[0x0F] = 0
	if sum > 0xFF {
		c.registers[0x0F] = 1
	}
	c.registers[x] = uint8(sum)
}

func (c *Chip8) op8xy5() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	var y uint8 = uint8((c.opcode & 0x00F0) >> 4)
	c.registers[0x0F] = 0
	if c.registers[x] > c.registers[y] {
		c.registers[0x0F] = 1
	}
	c.registers[x] -= c.registers[y]
}

func (c *Chip8) op8xy6() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	c.registers[0x0F] = c.registers[x] & 0x01
	c.registers[x] >>= 1
}

func (c *Chip8) op8xy7() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	var y uint8 = uint8((c.opcode & 0x00F0) >> 4)
	c.registers[0x0F] = 0
	if c.registers[y] > c.registers[x] {
		c.registers[0x0F] = 1
	}
	c.registers[x] = c.registers[y] - c.registers[x]
}

func (c *Chip8) op8xyE() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	c.registers[0x0F] = (c.registers[x] & 0x80) >> 7
	c.registers[x] <<= 1
}

func (c *Chip8) opExA1() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	if c.keypad[c.registers[x]] == 0 {
		c.pc += 2
	}
}

func (c *Chip8) opEx9E() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	if c.keypad[c.registers[x]] != 0 {
		c.pc += 2
	}
}

func (c *Chip8) opFx07() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	c.registers[x] = c.delayTimer
}

func (c *Chip8) opFx0A() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	for i, key := range c.keypad {
		if key != 0 {
			c.registers[x] = uint8(i)
			return
		}
	}
	c.pc -= 2
}

func (c *Chip8) opFx15() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	c.delayTimer = c.registers[x]
}

func (c *Chip8) opFx18() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	c.soundTimer = c.registers[x]
}

func (c *Chip8) opFx1E() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	c.index += uint16(c.registers[x])
}

func (c *Chip8) opFx29() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	c.index = uint16(c.registers[x]) * 5
}

func (c *Chip8) opFx33() {
	var x uint8 = uint8((c.opcode & 0x0F00) >> 8)
	c.memory[c.index] = c.registers[x] / 100
	c.memory[c.index+1] = (c.registers[x] / 10) % 10
	c.memory[c.index+2] = (c.registers[x] % 100) % 10
}

func (c *Chip8) opFx55() {
	var x uint16 = (c.opcode & 0x0F00) >> 8
	for i := uint16(0); i <= x; i++ {
		c.memory[c.index+i] = c.registers[i]
	}
}

func (c *Chip8) opFx65() {
	var x uint16 = (c.opcode & 0x0F00) >> 8
	for i := uint16(0); i <= x; i++ {
		c.registers[i] = c.memory[c.index+i]
	}
}

func (c *Chip8) Cycle() error {
	c.opcode = (uint16(c.memory[c.pc]) << 8) | uint16(c.memory[c.pc+1])
	c.pc += 2

	var firstDigit uint8 = uint8(c.opcode >> 12)
	var lastDigit uint8 = uint8(c.opcode & 0x000F)
	var lastTwoDigits uint8 = uint8(c.opcode & 0x00FF)
	var unknownInstruction error = errors.New("unknown instruction")

	switch firstDigit {
	case 0x00:
		switch c.opcode {
		case 0x00E0:
			c.op00E0()
		case 0x00EE:
			c.op00EE()
		default:
			return unknownInstruction
		}
	case 0x01:
		c.op1nnn()
	case 0x02:
		c.op2nnn()
	case 0x03:
		c.op3xkk()
	case 0x04:
		c.op4xkk()
	case 0x05:
		switch lastDigit {
		case 0x00:
			c.op5xy0()
		default:
			return unknownInstruction
		}
	case 0x06:
		c.op6xkk()
	case 0x07:
		c.op7xkk()
	case 0x08:
		switch lastDigit {
		case 0x0:
			c.op8xy0()
		case 0x1:
			c.op8xy1()
		case 0x2:
			c.op8xy2()
		case 0x3:
			c.op8xy3()
		case 0x4:
			c.op8xy4()
		case 0x5:
			c.op8xy5()
		case 0x6:
			c.op8xy6()
		case 0x7:
			c.op8xy7()
		case 0xE:
			c.op8xyE()
		default:
			return unknownInstruction
		}
	case 0x09:
		switch lastDigit {
		case 0x00:
			c.op9xy0()
		default:
			return unknownInstruction
		}
	case 0x0A:
		c.opAnnn()
	case 0x0B:
		c.opBnnn()
	case 0x0C:
		c.opCxkk()
	case 0x0D:
		c.opDxyn()
	case 0x0E:
		switch lastTwoDigits {
		case 0x9E:
			c.opEx9E()
		case 0xA1:
			c.opExA1()
		default:
			return unknownInstruction
		}
	case 0x0F:
		switch lastTwoDigits {
		case 0x07:
			c.opFx07()
		case 0x0A:
			c.opFx0A()
		case 0x15:
			c.opFx15()
		case 0x18:
			c.opFx18()
		case 0x1E:
			c.opFx1E()
		case 0x29:
			c.opFx29()
		case 0x33:
			c.opFx33()
		case 0x55:
			c.opFx55()
		case 0x65:
			c.opFx65()
		default:
			return unknownInstruction
		}
	default:
		return unknownInstruction
	}

	if c.delayTimer > 0 {
		c.delayTimer--
	}

	if c.soundTimer > 0 {
		c.soundTimer--
	}

	return nil
}
