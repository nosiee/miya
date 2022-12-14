package main

import (
	"flag"
	"fmt"
	"log"
	"miya/internal/memory"
	"miya/internal/screen"
	"miya/internal/vm"
	"os"
)

func main() {
	var fname string
	var delay uint64
	var backgroundColor uint64
	var pixelColor uint64
	var debugMode bool

	flag.StringVar(&fname, "fname", "", "Rom filename")
	flag.Uint64Var(&delay, "delay", 1, "Delay in ms for virtualmachine and screen")
	flag.Uint64Var(&backgroundColor, "background-color", 0x00000000, "Background color in uint32 for the screen")
	flag.Uint64Var(&pixelColor, "pixel-color", 0xFFFFFF00, "Pixel color in uint32 for the screen")
	flag.BoolVar(&debugMode, "debug-mode", false, "Run in debug mode")
	flag.Parse()

	buffer, err := os.ReadFile(fname)
	if err != nil {
		log.Fatalf("os.ReadFile(): %v\n", err)
	}

	mw, err := screen.NewMainWindow(fmt.Sprintf("CHIP8 - %s | %dms", fname, delay), 640, 320, backgroundColor, pixelColor)
	if err != nil {
		log.Fatalf("screen.NewMainWindow(): %v\n", err)
	}

	mem := memory.NewMemory(memory.CHIP8_MEMORY_SIZE)
	stack := memory.NewStack(memory.CHIP8_STACK_SIZE)
	vm := vm.NewVirtualMachine(mem, stack, mw, delay, debugMode)

	mem.WriteArray(0x200, buffer)
	go vm.EvalLoop()

	if debugMode {
		dw, err := screen.NewDebugWindow("Debug", 380, 120)
		if err != nil {
			log.Fatalf("screen.NewDebugWindow(): %v\n", err)
		}

		go vm.Debug()
		screen.ShowWindows(delay, mw, dw)
	}

	screen.ShowWindows(delay, mw)
}
