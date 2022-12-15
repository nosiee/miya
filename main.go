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

	flag.StringVar(&fname, "fname", "", "Rom filename")
	flag.Uint64Var(&delay, "delay", 10, "Delay in ms for virtualmachine and screen")
	flag.Parse()

	buffer, err := os.ReadFile(fname)
	if err != nil {
		log.Fatalf("os.ReadFile(): %v\n", err)
	}

	screen, err := screen.NewScreen(640, 320, fmt.Sprintf("CHIP8 - %s", fname), delay)
	if err != nil {
		log.Fatalf("screen.NewScreen(): %v\n", err)
	}

	mem := memory.NewMemory(memory.CHIP8_MEMORY_SIZE)
	stack := memory.NewStack(memory.CHIP8_STACK_SIZE)
	vm := vm.NewVirtualMachine(mem, stack, screen, delay)

	mem.WriteArray(0x200, buffer)

	go screen.Show()
	vm.EvalLoop()
}
