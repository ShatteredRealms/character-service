package skills

type AbilityName string

const (
	Strength   AbilityName = "Strength"
	Stamina    AbilityName = "Stamina"
	Agility    AbilityName = "Agility"
	Perception AbilityName = "Perception"
	Intellect  AbilityName = "Intellect"
	Willpower  AbilityName = "Willpower"
)

type SkillName string

const (
	// Physical Skills
	Endurance           SkillName = "Endurance"
	EnergyReserve       SkillName = "Energy Reserve"
	CloseCombatEvasion  SkillName = "Close Combat Evasion"
	RangedCombatEvasion SkillName = "Ranged Combat Evasion"
	Deflection          SkillName = "Deflection"

	// Melee Skills
	MartialArts SkillName = "Martial Arts"
	OneHanded   SkillName = "One-Handed"
	TwoHanded   SkillName = "Two-Handed"
	Piercing    SkillName = "Piercing"
	WeaponSpeed SkillName = "Weapon Speed"
	MeleeSpeed  SkillName = "Melee Speed"

	// Melee Specials
	SneakAttack SkillName = "Sneak Attack"
	QuickStrike SkillName = "Quick Strike"
	PowerStrike SkillName = "Power Strike"
	RapidBlows  SkillName = "Rapid Blows"

	// Ranged Skills
	BowSkill      SkillName = "Bow"
	PistolSkill   SkillName = "Pistol"
	RifleSkill    SkillName = "Rifle"
	ShotgunSkill  SkillName = "Shotgun"
	SniperSkill   SkillName = "Sniper"
	ThrowingSkill SkillName = "Throwing"
	RangedSpeed   SkillName = "Ranged Speed"

	// Ranged Specials
	SteadyShot SkillName = "Steady Shot"
	QuickShot  SkillName = "Quick Shot"
	PowerShot  SkillName = "Power Shot"
	RapidFire  SkillName = "Rapid Fire"
	SilentShot SkillName = "Silent Shot"

	// Energy Skills
	Creation      SkillName = "Creation"
	Metamorphosis SkillName = "Metamorphosis"
	Remedy        SkillName = "Remedy"
	Rejuvenation  SkillName = "Rejuvenation"
	SpaceTime     SkillName = "Space Time"
	Destruction   SkillName = "Destruction"
	EnergySpeed   SkillName = "Energy Speed"

	// Tradeskills
	Alchemy      SkillName = "Alchemy"
	Programming  SkillName = "Programming"
	Construction SkillName = "Construction"
	Engineering  SkillName = "Engineering"

	// Healing
	FirstAid SkillName = "First Aid"
	Surgery  SkillName = "Surgery"

	// Adventuring
	Speed   SkillName = "Speed"
	Flying  SkillName = "Flying"
	Stealth SkillName = "Stealth"
)
