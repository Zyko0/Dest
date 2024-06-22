package mod

// Difficulty

type Difficulty byte

const (
	Balanced Difficulty = iota
	Hard
	HardestHard
	// TODO: TwoHard (?) same as hardest but two bosses each stage
	diffMax
)

var current = Balanced

func SetDifficulty(diff Difficulty) {
	current = diff
}

// Modifiers

type Modifier byte

const (
	BossHPMult Modifier = iota
	BossDamageMult
	BossProjectileMult
	BossProjectileSpeed
	BossAoESpawnSpeed
	BossDashSpeed
	BossPatternDelay
	modMax
)

var (
	mods = [diffMax][modMax]float64{
		Balanced: {
			BossHPMult:          1.,
			BossDamageMult:      1.,
			BossProjectileMult:  1.,
			BossProjectileSpeed: 1.,
			BossAoESpawnSpeed:   1.,
			BossDashSpeed:       1.,
			BossPatternDelay:    1.,
		},
	}
)

func Get(m Modifier) float64 {
	return mods[current][m]
}
