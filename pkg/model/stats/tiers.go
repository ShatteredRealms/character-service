package stats

import "github.com/ShatteredRealms/character-service/pkg/model/game"

type CharacterTierLevel int
type CharacterTier struct {
	Level     CharacterTierLevel
	Name      map[game.Profession]string
	Threshold int
}

const (
	TierLevel1 CharacterTierLevel = iota + 1
	TierLevel2
)

var (
	Tier1 = CharacterTier{
		Level: TierLevel1,
		Name: map[game.Profession]string{
			game.ProfessionNecromancer: "Novice",
		},
		Threshold: 0,
	}
	Tier2 = CharacterTier{
		Level: TierLevel2,
		Name: map[game.Profession]string{
			game.ProfessionNecromancer: "Apprentice",
		},
		Threshold: 10,
	}
)
