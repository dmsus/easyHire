package models

type LevelWeight struct {
	Level  string `gorm:"type:varchar(20);primaryKey" json:"level"`
	Weight int    `gorm:"not null" json:"weight"`
}

func (LevelWeight) TableName() string {
	return "level_weights"
}

type LevelThreshold struct {
	Level     string `gorm:"type:varchar(20);primaryKey" json:"level"`
	Threshold int    `gorm:"not null" json:"threshold"`
}

func (LevelThreshold) TableName() string {
	return "level_thresholds"
}

// ScoringConfig для хранения формулы оценки
type ScoringConfig struct {
	LevelWeights     map[string]int     `json:"level_weights"`
	CompetencyWeights map[string]float64 `json:"competency_weights"`
	LevelThresholds  map[string]int     `json:"level_thresholds"`
	TimeBonuses      TimeBonusConfig    `json:"time_bonuses"`
	AntiFarming      AntiFarmingConfig  `json:"anti_farming"`
}

type TimeBonusConfig struct {
	VeryFastRatio float64 `json:"very_fast_ratio"`  // < 0.3
	FastRatio     float64 `json:"fast_ratio"`       // < 0.7
	VeryFastBonus float64 `json:"very_fast_bonus"`  // 1.2
	FastBonus     float64 `json:"fast_bonus"`       // 1.1
}

type AntiFarmingConfig struct {
	MaxJuniorScorePerDay int     `json:"max_junior_score_per_day"`
	MaxSameLevelRatio    float64 `json:"max_same_level_ratio"`
}
