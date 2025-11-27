# Assessment Framework & Scoring System

## Fibonacci-Based Scoring System

### Level Weights (Fibonacci Sequence)
```go
LEVEL_WEIGHTS = {
    "junior": 1,    // F1
    "middle": 2,    // F2  
    "senior": 3,    // F3
    "expert": 5     // F4
}
```
## Competency Weights
```go
COMPETENCY_WEIGHTS = {
    // Core Go
    "go_fundamentals": 1.0,
    "data_structures_go": 1.1,
    "memory_management": 1.1,
    "concurrency": 1.3,
    "http_go": 1.0,
    
    // System Design & Architecture
    "system_design": 1.3,
    "microservices": 1.2,
    "containerization": 1.1,
    "reliability": 1.2,
    "message_brokers": 1.2,
    
    // Software Engineering
    "software_design": 1.2,
    "architecture": 1.3,
    "quality_assurance": 1.1,
    "optimization": 1.1,
    
    // General Development
    "sdlc": 1.0,
    "requirements": 1.0,
    "ci_cd": 1.1,
    "git": 1.0,
    
    // Security
    "web_security": 1.2,
    "data_security": 1.3
}
```
## Scoring Formula
```go
func calculateTaskScore(competency string, level string, isCorrect bool, timeBonus float64) float64 {
    baseScore := LEVEL_WEIGHTS[level] * COMPETENCY_WEIGHTS[competency]
    if isCorrect {
        return baseScore * timeBonus
    }
    return 0
}

func calculateFinalScore(testResults []TestResult) ScoreReport {
    totalScore := 0.0
    maxPossibleScore := 0.0
    competencyScores := make(map[string]CompetencyScore)
    
    for _, result := range testResults {
        taskMaxScore := calculateTaskScore(result.Competency, result.Level, true, 1.0)
        maxPossibleScore += taskMaxScore
        
        if result.IsCorrect {
            timeBonus := calculateTimeBonus(result.TimeSpent, result.MaxTime)
            taskScore := calculateTaskScore(result.Competency, result.Level, true, timeBonus)
            totalScore += taskScore
            
            // Aggregate by competency
            compScore := competencyScores[result.Competency]
            compScore.Achieved += taskScore
            compScore.Possible += taskMaxScore
            competencyScores[result.Competency] = compScore
        }
    }
    
    percentage := (totalScore / maxPossibleScore) * 100
    return ScoreReport{
        TotalScore: totalScore,
        Percentage: percentage,
        Level: determineLevel(percentage, competencyScores),
        CompetencyBreakdown: competencyScores,
    }
}
```
## Level Thresholds (Fibonacci-Based)
```go
LEVEL_THRESHOLDS = {
    "junior": 8,      // F6
    "middle": 21,     // F8
    "senior": 55,     // F10
    "expert": 144     // F12
}

func determineLevel(percentage float64, competencies map[string]CompetencyScore) string {
    // Check if candidate has solved at least one task at each required level
    maxSolvedLevel := getMaxSolvedLevel(competencies)
    
    if percentage >= 85 && maxSolvedLevel == "expert":
        return "EXPERT"
    elif percentage >= 70 && maxSolvedLevel in ["senior", "expert"]:
        return "SENIOR" 
    elif percentage >= 55 && maxSolvedLevel in ["middle", "senior", "expert"]:
        return "MIDDLE"
    elif percentage >= 40 && maxSolvedLevel in ["junior", "middle", "senior", "expert"]:
        return "JUNIOR"
    else:
        return "TRAINEE"
}
```
## Time-Based Bonuses
```go
func calculateTimeBonus(actualTime, maxTime time.Duration) float64 {
    timeRatio := float64(actualTime) / float64(maxTime)
    
    if timeRatio < 0.3:
        return 1.2  // +20% for very fast solution
    elif timeRatio < 0.7:
        return 1.1  // +10% for fast solution
    else:
        return 1.0  // Standard score
}
```
## Anti-Farming Protection
```go
// Prevent farming easy tasks
const (
    MAX_JUNIOR_SCORE_PER_DAY = 20
    MAX_SAME_LEVEL_RATIO = 0.5  // Max 50% of total score from one level
)

func validateScoreDistribution(scores map[string]int) bool {
    totalScore := calculateTotalScore(scores)
    
    // Check if junior tasks exceed daily limit
    if scores["junior"] * LEVEL_WEIGHTS["junior"] > MAX_JUNIOR_SCORE_PER_DAY {
        return false
    }
    
    // Check if one level dominates
    for level, count := range scores {
        levelScore := count * LEVEL_WEIGHTS[level]
        if float64(levelScore) / float64(totalScore) > MAX_SAME_LEVEL_RATIO {
            return false
        }
    }
    
    return true
}
```
## Test Structure & Distribution
## 60-Minute Assessment (20 Questions)
```go
TEST_STRUCTURE = {
    "question_types": {
        "multiple_choice": 10,    // 50% - Theory and quick knowledge checks
        "coding": 6,              // 30% - Practical coding tasks
        "architecture": 3,        // 15% - System design questions
        "debugging": 1            // 5%  - Code review and bug finding
    },
    
    "level_distribution": {
        "junior": 6,    // 30% - Foundational knowledge
        "middle": 8,    // 40% - Core competency
        "senior": 4,    // 20% - Advanced topics  
        "expert": 2     // 10% - Expert-level challenges
    },
    
    "time_allocation": {
        "multiple_choice": 1.5,   // minutes per question
        "coding": 5.0,            // minutes per question  
        "architecture": 7.0,      // minutes per question
        "debugging": 4.0          // minutes per question
    }
}
```
## Competency Clusters
```go
COMPETENCY_CLUSTERS = {
    "backend_core": [
        "go_fundamentals",
        "data_structures_go", 
        "concurrency",
        "memory_management"
    ],
    
    "system_design": [
        "architecture",
        "microservices",
        "reliability",
        "message_brokers"
    ],
    
    "software_engineering": [
        "software_design", 
        "quality_assurance",
        "optimization",
        "ci_cd"
    ],
    
    "security_infra": [
        "web_security",
        "data_security",
        "containerization"
    ]
}
```
## Adaptive Assessment Generation
```go
func generateAdaptiveAssessment(roleRequirements RoleRequirements) Assessment {
    assessment := Assessment{
        TotalQuestions: 20,
        TimeLimit: 60 * time.Minute,
        Questions: []Question{},
    }
    
    // Distribute questions based on role requirements
    for cluster, weight := range roleRequirements.ClusterWeights {
        clusterQuestions := int(float64(20) * weight)
        questions := generateClusterQuestions(cluster, clusterQuestions, roleRequirements.TargetLevel)
        assessment.Questions = append(assessment.Questions, questions...)
    }
    
    // Ensure coverage across levels
    assessment.Questions = balanceLevelDistribution(assessment.Questions, roleRequirements.TargetLevel)
    
    return assessment
}
```
## Benefits of This System
## 1. Mathematical Foundation
- **Fibonacci sequence** naturally represents increasing complexity
- **Non-linear progression** matches real-world skill growth
- **Clear milestones** at 8, 21, 55, 144 points
## 2. Flexibility
```python
# Multiple paths to reach Middle level (21 points):
path1 = 10 * middle * 2 = 20 + time bonuses = 21+ ✅
path2 = 5 * middle * 2 + 10 * junior * 1 = 10 + 10 = 20 + bonuses ✅  
path3 = 3 * senior * 3 + 6 * junior * 1 = 9 + 6 = 15 + bonuses ❌ (need more)
path4 = 2 * expert * 5 + 2 * junior * 1 = 10 + 2 = 12 + bonuses ❌
```
## 3. Incentive Structure
- Quality over quantity - harder tasks worth more
- Time efficiency rewarded but not required
- Balanced skill development encouraged
## 4.  Anti-Gaming Protection
- Daily limits on easy task farming
- Level diversity requirements
- Quality multipliers for clean code

