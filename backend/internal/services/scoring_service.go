package services

import (
    "github.com/easyhire/backend/internal/models"
)

type ScoringService interface {
    CalculateFinalScore(answers []models.CandidateAnswer, questions []models.Question) (*models.ScoreResult, error)
}

type scoringService struct{}

func NewScoringService() ScoringService {
    return &scoringService{}
}

func (s *scoringService) CalculateFinalScore(answers []models.CandidateAnswer, questions []models.Question) (*models.ScoreResult, error) {
    // Базовая реализация Fibonacci scoring system
    var totalScore float64
    var maxPossibleScore float64
    
    for _, answer := range answers {
        for _, question := range questions {
            if answer.QuestionID == question.ID {
                // Базовый вес уровня (Fibonacci)
                levelWeights := map[string]float64{
                    "junior": 1,
                    "middle": 2,
                    "senior": 3,
                    "expert": 5,
                }
                
                // Вес компетенции
                competencyWeights := map[string]float64{
                    "go_fundamentals":     1.0,
                    "concurrency":         1.3,
                    "system_design":       1.3,
                    "architecture":        1.3,
                    "data_structures_go":  1.1,
                    "memory_management":   1.1,
                    "http_go":             1.0,
                    "microservices":       1.2,
                    "reliability":         1.2,
                    "message_brokers":     1.2,
                    "software_design":     1.2,
                    "quality_assurance":   1.1,
                    "optimization":        1.1,
                    "web_security":        1.2,
                    "data_security":       1.3,
                }
                
                levelWeight := levelWeights[string(question.Difficulty)]
                competencyWeight := competencyWeights[question.Competency]
                
                maxPossibleScore += levelWeight * competencyWeight
                
                if answer.IsCorrect {
                    // Бонус за время
                    timeBonus := 1.0
                    if answer.TimeSpent > 0 {
                        timeRatio := float64(answer.TimeSpent) / 300.0 // 5 минут на вопрос
                        if timeRatio < 0.3 {
                            timeBonus = 1.2
                        } else if timeRatio < 0.7 {
                            timeBonus = 1.1
                        }
                    }
                    
                    totalScore += levelWeight * competencyWeight * timeBonus
                }
                break
            }
        }
    }
    
    percentage := 0.0
    if maxPossibleScore > 0 {
        percentage = (totalScore / maxPossibleScore) * 100
    }
    
    // Определение уровня
    level := determineLevel(percentage)
    
    return &models.ScoreResult{
        TotalScore: totalScore,
        Percentage: percentage,
        Level:      level,
    }, nil
}

func determineLevel(percentage float64) string {
    if percentage >= 85 {
        return "EXPERT"
    } else if percentage >= 70 {
        return "SENIOR"
    } else if percentage >= 55 {
        return "MIDDLE"
    } else if percentage >= 40 {
        return "JUNIOR"
    } else {
        return "TRAINEE"
    }
}
