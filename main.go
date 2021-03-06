package main

import (
	"fmt"
	"math"
)

// lc3vm - Tyler Mace <tyler@madhive.com>
//
// Notes:
//  - Simulates the LC-3 educational computer platform
//  - The ISA is small but includes most of the common features
//	  found in modern architectures

// Register is a makeshift enum type
type Register int

const (
	r0  Register = iota // General purpose
	r1                  // General purpose
	r2                  // General purpose
	r3                  // General purpose
	r4                  // General purpose
	r5                  // General purpose
	r6                  // General purpose
	r7                  // General purpose
	rPC                 // Program Counter
	rCD                 // Conditional
	rCN                 // Count
)

// Flag is a makeshift enum type
type Flag int

const (
	flPOS Flag = 1 << 0 // Positive
	flZRO Flag = 1 << 1 // Zero
	flNEG Flag = 1 << 2 // Negative
)

// OpCode is a makeshift enum type
type OpCode int

const ( // Order is important: BR = 0000, ADD = 0001, LD = 0010, ...
	opBR  OpCode = iota // Branch
	opADD               // Add
	opLD                // Load
	opST                // Store
	opJSR               // Jump register
	opAND               // Bitwise and
	opLDR               // Load register
	opSTR               // Store register
	opRTI               // Return from inturrupt (Unimplemented)
	opNOT               // Bitwise not
	opLDI               // Load indirect
	opSTI               // Store indirect
	opJMP               // Jump
	opRES               // Reserve (Unimplemented)
	opLEA               // Load effective address
	opTRA               // Trap (Clock resets and halt operations, for example)
)

// RegisterMM is a makeshift enum
type RegisterMM int

const (
	rKBSR RegisterMM = 0xFE00
	rKBDR RegisterMM = 0xFE02
)

const (
	pcStart = 0x3000 // OS space < 3000 -- everything else can be used as program memory
)

// Memory

var memory []uint16
var registers []uint16

func main() {
	memory = make([]uint16, math.MaxUint16) // The available memory in the LC-3 is limited to 128kb (65k addressable locations)
	registers = make([]uint16, rCN)

	// Start our execution (PC) at 0x3000
	registers[rPC] = pcStart

	// Hardcode an initial instruction at rPC + 1
	memory[pcStart+1] = 0x1001 // This instruction is an add instruction that adds values from r0 and r1 and stores back in r0
	memory[pcStart+2] = 0x1001
	memory[pcStart+3] = 0x1001
	memory[pcStart+4] = 0x1024 // This instruction uses the imm scalar value 4 instead of r1 as its second operand resulting in r0 += 4

	// ... and set some initial values in registers
	registers[r1] = 0x1

	// Execution procedure:
	// 1. Load instruction from memory at the address held by PC
	// 2. Increment the PC register (resulting in new instruction address; this may be subsequently changed by our resulting instruction execution)
	// 3. Read the OpCode
	// 4. Perform the instruction/operation
	// 5. Start over

	// Event loop
	running := true
	for {
		registers[rPC] = registers[rPC] + 1

		if registers[rPC] >= (uint16)(len(memory)) {
			break
		}

		// Fetch
		instruction := read(registers[rPC])

		// Decode
		op, dst, srcA, flag, srcB := decode(instruction)

		if op != 0x0 { // Dont log no-ops
			fmt.Printf("Instruction: %b\r\nOperation: %d\r\nOperands: %d, %d, %d, %d\r\n", instruction, op, dst, srcA, flag, srcB)
		}

		// Execute
		switch op {
		case 0x0: // Branch
		case 0x1: // Add
			// f(a,b) == a+b
			if flag == 0x1 { // We know that the srcB is a 5-bit unsigned int used as immediate scalar
				registers[dst] = registers[srcA] + srcB
			} else {
				registers[dst] = registers[srcA] + registers[srcB]
			}

			updateFlag(dst)
		case 0x2: // Load
		case 0x3: // Store
		case 0x4: // Jump register
		case 0x5: // Bitwise and
		case 0x6: // Load register
		case 0x7: // Store register
		case 0x8: // Return from interrupt
		case 0x9: // Bitwise not
		case 0xA: // Load indirect
		case 0xB: // Store indirect
		case 0xC: // Jump
		case 0xD: // Unused -- can use for testing
		case 0xE: // Load effective address
		case 0xF: // Trap
			running = false
		default: // OpCode == HCF ;)
			fmt.Println("Halting & catching fire...")
			running = false
		}

		// Exit our event loop
		if !running {
			break
		}
	}

	fmt.Printf("Registers: %X\r\n", registers)
}

func decode(r uint16) (uint16, uint16, uint16, uint16, uint16) {
	return r >> 12, (r >> 9) & 0x7, (r >> 6) & 0x7, (r >> 5) & 0x1, r & 0x1F
}

func read(address uint16) uint16 {
	return memory[address]
}

func write(address uint16, value uint16) {
	memory[address] = value
}

func signExtend(x uint16, count uint) uint16 {
	if ((x >> (count - 1)) & 0x1) == 0x1 {
		x = (0xFFFF << count)
	}

	return x
}

func swap(x uint16) uint16 {
	return x<<8 | x>>8
}

func updateFlag(f uint16) {
	if registers[f] == 0x0 {
		registers[rCD] = (uint16)(flZRO)
	} else if registers[f]>>15 == 0x1 {
		registers[rCD] = (uint16)(flNEG)
	} else {
		registers[rCD] = (uint16)(flPOS)
	}
}
