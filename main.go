package main

import (
	"log"
	"miya/internal/memory"
	"miya/internal/vm"
	"os"
)

func main() {
	buffer, err := os.ReadFile("chip8-roms/Cave.ch8")
	if err != nil {
		log.Fatalf("os.ReadFile(): %v\n", err)
	}

	mem := memory.NewMemory(memory.CHIP8_MEMORY_SIZE)
	stack := memory.NewStack(memory.CHIP8_STACK_SIZE)
	vm := vm.NewVirtualMachine(mem, stack)

	mem.WriteArray(0x200, buffer)
	vm.EvalLoop()
}
