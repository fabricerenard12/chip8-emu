package main

import (
	"chip8-emu/pkg/chip8"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 640
	windowHeight = 320
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalf("could not initialize sdl: %v", err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("CHIP-8 Emulator", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("could not create window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("could not create renderer: %v", err)
	}
	defer renderer.Destroy()

	chip8 := chip8.GetInstance()
	if err := chip8.LoadROM("IBM.ch8"); err != nil {
		log.Fatalf("could not load ROM: %v", err)
	}

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		if err := chip8.Cycle(); err != nil {
			log.Printf("could not emulate cycle: %v", err)
		}

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		video := chip8.GetVideo()
		for y := 0; y < 32; y++ {
			for x := 0; x < 64; x++ {
				color := uint32(video[y*64+x])
				if color != 0 {
					renderer.SetDrawColor(255, 255, 255, 255)
				} else {
					renderer.SetDrawColor(0, 0, 0, 255)
				}
				renderer.FillRect(&sdl.Rect{
					X: int32(x) * (windowWidth / 64),
					Y: int32(y) * (windowHeight / 32),
					W: windowWidth / 64,
					H: windowHeight / 32,
				})
			}
		}

		renderer.Present()
		sdl.Delay(16)
	}
}
