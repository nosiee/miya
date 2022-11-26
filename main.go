package main

import (
	"log"
	"miya/internal/memory"
	"miya/internal/screen"
	"miya/internal/vm"
	"os"
)

func main() {
	buffer, err := os.ReadFile("chip8-roms/Cave.ch8")
	if err != nil {
		log.Fatalf("os.ReadFile(): %v\n", err)
	}

	screen := screen.NewScreen(640, 320, "CHIP8")
	mem := memory.NewMemory(memory.CHIP8_MEMORY_SIZE)
	stack := memory.NewStack(memory.CHIP8_STACK_SIZE)
	vm := vm.NewVirtualMachine(mem, stack, screen)

	mem.WriteArray(0x200, buffer)
	go vm.EvalLoop()

	screen.Show()
}
