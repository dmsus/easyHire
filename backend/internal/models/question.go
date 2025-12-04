package models

// Question представляет вопрос
type Question struct {
    BaseModel
    Title       string           `gorm:"type:varchar(500);not null" json:"title"`
    Description string           `gorm:"type:text" json:"description"`
    Type        QuestionType     `gorm:"type:varchar(50);not null" json:"type"`
    Difficulty  DifficultyLevel  `gorm:"type:varchar(20);not null" json:"difficulty"`
    Competency  string           `gorm:"type:varchar(100);not null" json:"competency"`
    Tags        []QuestionTag    `gorm:"foreignKey:QuestionID" json:"tags"`
    Options     []QuestionOption `gorm:"foreignKey:QuestionID" json:"options"`
    TestCases   []TestCase       `gorm:"foreignKey:QuestionID" json:"test_cases"`
    Explanation string           `gorm:"type:text" json:"explanation"`
    TimeLimit   int              `gorm:"default:300" json:"time_limit"` // в секундах
    Points      int              `gorm:"default:1" json:"points"`
    IsActive    bool             `gorm:"default:true" json:"is_active"`
    CreatedBy   string           `gorm:"type:uuid;not null" json:"created_by"`
}

// QuestionTag тег вопроса
type QuestionTag struct {
    BaseModel
    QuestionID string `gorm:"type:uuid;not null;index" json:"question_id"`
    Tag        string `gorm:"type:varchar(100);not null" json:"tag"`
}

// QuestionOption вариант ответа
type QuestionOption struct {
    BaseModel
    QuestionID string `gorm:"type:uuid;not null;index" json:"question_id"`
    Text       string `gorm:"type:text;not null" json:"text"`
    IsCorrect  bool   `gorm:"default:false" json:"is_correct"`
    Order      int    `gorm:"not null" json:"order"`
}

// TestCase тестовый случай
type TestCase struct {
    BaseModel
    QuestionID  string `gorm:"type:uuid;not null;index" json:"question_id"`
    Input       string `gorm:"type:text" json:"input"`
    Expected    string `gorm:"type:text" json:"expected"`
    IsHidden    bool   `gorm:"default:false" json:"is_hidden"`
    Order       int    `gorm:"not null" json:"order"`
}
