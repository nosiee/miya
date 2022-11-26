package vm

const (
	CLC       = 0x0000
	JP        = 0x1000
	CALL      = 0x2000
	SE_VX     = 0x3000
	SNE       = 0x4000
	SE_VX_VY  = 0x5000
	LD_VX     = 0x6000
	ADD       = 0x7000
	VX_VY     = 0x8000
	SNE_VX_VY = 0x9000
	LD_I      = 0xA000
	JP_V0     = 0xB000
	RND       = 0xC000
	DRW       = 0xD000
	SKP       = 0xE000
	LDF       = 0xF000
)
