package vm

import (
	"fmt"
	"math/rand"
	"miya/internal/memory"
	"miya/internal/screen"
	"time"
)

type VirtualMachine struct {
	registers    Registers
	memory       *memory.Memory
	stack        *memory.Stack
	screen       *screen.Screen
	instructions map[uint16]func(uint16)
}

type Registers struct {
	I  uint16
	PC uint16
	V  []byte
}

func NewVirtualMachine(memory *memory.Memory, stack *memory.Stack, screen *screen.Screen) *VirtualMachine {
	return &VirtualMachine{
		registers: Registers{
			PC: 0x200,
			V:  make([]byte, 0x10),
		},
		memory:       memory,
		stack:        stack,
		screen:       screen,
		instructions: make(map[uint16]func(uint16)),
	}
}

func (vm *VirtualMachine) EvalLoop() {
	vm.instructions[CLC] = vm.clc
	vm.instructions[JP] = vm.jp
	vm.instructions[CALL] = vm.call
	vm.instructions[SE_VX] = vm.sevx
	vm.instructions[SNE] = vm.sne
	vm.instructions[SE_VX_VY] = vm.sevxvy
	vm.instructions[LD_VX] = vm.ldvx
	vm.instructions[ADD] = vm.add
	vm.instructions[VX_VY] = vm.vxvy
	vm.instructions[SNE_VX_VY] = vm.snevxvy
	vm.instructions[LD_I] = vm.ldi
	vm.instructions[JP_V0] = vm.jpv0
	vm.instructions[RND] = vm.rnd
	vm.instructions[DRW] = vm.drw
	vm.instructions[SKP] = vm.skp
	vm.instructions[LDF] = vm.ldf

	for {
		opcode := vm.memory.ReadOpcode(vm.registers.PC)
		vm.instructions[(opcode & 0xF000)](opcode)

		time.Sleep(time.Second / 4)
	}
}

func (vm *VirtualMachine) clc(opcode uint16) {
	otype := (opcode & 0x0FFF)

	if otype == 0x0E0 {
		vm.screen.Clear()
		vm.registers.PC += 2
	}

	if otype == 0x0EE {
		vm.registers.PC = vm.stack.Pop()
	}
}

func (vm *VirtualMachine) jp(opcode uint16) {
	vm.registers.PC = (opcode & 0x0FFF)
}

func (vm *VirtualMachine) call(opcode uint16) {
	vm.stack.Push(vm.registers.PC)
	vm.registers.PC = (opcode & 0x0FFF)
}

func (vm *VirtualMachine) sevx(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	nn := (opcode & 0x00FF)

	if vm.registers.V[x] == byte(nn) {
		vm.registers.PC += 4
		return
	}

	vm.registers.PC += 2
}

func (vm *VirtualMachine) sne(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	nn := (opcode & 0x00FF)

	if vm.registers.V[x] != byte(nn) {
		vm.registers.PC += 4
		return
	}

	vm.registers.PC += 2
}

func (vm *VirtualMachine) sevxvy(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4

	if vm.registers.V[x] == vm.registers.V[y] {
		vm.registers.PC += 4
		return
	}

	vm.registers.PC += 2
}

func (vm *VirtualMachine) ldvx(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	nn := (opcode & 0x00FF)

	vm.registers.V[x] = byte(nn)
	vm.registers.PC += 2
}

func (vm *VirtualMachine) add(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	nn := (opcode & 0x00FF)

	vm.registers.V[x] += byte(nn)
	vm.registers.PC += 2
}

func (vm *VirtualMachine) vxvy(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4

	switch opcode & 0x000F {
	case 0:
		vm.registers.V[x] = vm.registers.V[y]
	case 1:
		vm.registers.V[x] |= vm.registers.V[y]
	case 2:
		vm.registers.V[x] &= vm.registers.V[y]
	case 3:
		vm.registers.V[x] ^= vm.registers.V[y]
	case 4:
		vm.registers.V[x] += vm.registers.V[y]
	case 5:
		vm.registers.V[x] -= vm.registers.V[y]
	case 6:
		vm.registers.V[x] >>= 1
	case 7:
		vm.registers.V[x] = vm.registers.V[y] - vm.registers.V[x]
	case 0xe:
		vm.registers.V[x] <<= 1
	}

	vm.registers.PC += 2
}

func (vm *VirtualMachine) snevxvy(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4

	if vm.registers.V[x] != vm.registers.V[y] {
		vm.registers.PC += 4
		return
	}

	vm.registers.PC += 2
}

func (vm *VirtualMachine) ldi(opcode uint16) {
	vm.registers.I = (opcode & 0x0FFF)
	vm.registers.PC += 2
}

func (vm *VirtualMachine) jpv0(opcode uint16) {
	vm.registers.PC = (uint16(vm.registers.V[0]) + (opcode & 0x0FFF))
}

func (vm *VirtualMachine) rnd(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	nn := (opcode & 0x00FF)

	rand.Seed(time.Now().UnixNano())
	vm.registers.V[x] = byte((rand.Intn(0xff-0x00) + 0x00) & int(nn))
	vm.registers.PC += 2
}

func (vm *VirtualMachine) drw(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	n := (opcode & 0x000F)

	fmt.Println(x, y, n)
	vm.registers.PC += 2
}

func (vm *VirtualMachine) skp(opcode uint16) {
	println("0xe000 TODO")
	vm.registers.PC += 2
}

func (vm *VirtualMachine) ldf(opcode uint16) {
	println("0xf000 TODO")
	vm.registers.PC += 2
}
