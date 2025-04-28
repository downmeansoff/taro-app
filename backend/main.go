package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type Reading struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	Cards     []string  `json:"cards"`
	Timestamp time.Time `json:"timestamp"`
}

type TarotCard struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Keywords    string `json:"keywords"`
}

var (
	readingsDB = make(map[string][]Reading)
	cardsDB    = generateTarotCards()
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	r := gin.Default()
	r.Static("/api/images", "./static/images")

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		c.Next()
	})

	r.POST("/api/auth", telegramAuthHandler)
	r.GET("/api/cards", getCardsHandler)
	r.POST("/api/daily", dailyReadingHandler)
	r.POST("/api/ask", askQuestionHandler)
	r.GET("/api/history/:userID", getHistoryHandler)
	r.GET("/api/random-question", randomQuestionHandler)

	r.Run(":8080")
}

func telegramAuthHandler(c *gin.Context) {
	data := c.Request.URL.Query()
	if !validateTelegramData(data, os.Getenv("TELEGRAM_BOT_TOKEN")) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication"})
		return
	}

	claims := jwt.MapClaims{
		"tg_id": data.Get("id"),
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	c.JSON(http.StatusOK, gin.H{
		"token":   tokenStr,
		"user_id": data.Get("id"),
	})
}

func dailyReadingHandler(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id"`
		Card   string `json:"card"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	response, err := generateDeepSeekResponse(req.Card, "daily")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Гадание не удалось. Попробуйте позже",
		})
		return
	}

	newReading := Reading{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Question:  "Ежедневный расклад",
		Answer:    formatResponse(response),
		Cards:     []string{req.Card},
		Timestamp: time.Now(),
	}

	readingsDB[req.UserID] = append(readingsDB[req.UserID], newReading)

	c.JSON(http.StatusOK, gin.H{
		"reading":   newReading,
		"remaining": decrementFreeReadings(req.UserID),
	})
}

func askQuestionHandler(c *gin.Context) {
	var req struct {
		UserID   string   `json:"user_id"`
		Question string   `json:"question"`
		Cards    []string `json:"cards"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if len(req.Question) > 200 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Question too long"})
		return
	}

	response, err := generateDeepSeekResponse(strings.Join(req.Cards, ","), "question", req.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Анализ карт не удался. Попробуйте другой расклад",
		})
		return
	}

	newReading := Reading{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Question:  req.Question,
		Answer:    formatResponse(response),
		Cards:     req.Cards,
		Timestamp: time.Now(),
	}

	readingsDB[req.UserID] = append(readingsDB[req.UserID], newReading)

	c.JSON(http.StatusOK, gin.H{
		"reading":   newReading,
		"remaining": decrementFreeReadings(req.UserID),
	})
}

func getHistoryHandler(c *gin.Context) {
	userID := c.Param("userID")
	readings, exists := readingsDB[userID]
	if !exists {
		c.JSON(http.StatusOK, []Reading{})
		return
	}
	c.JSON(http.StatusOK, readings)
}

func randomQuestionHandler(c *gin.Context) {
	response, err := generateDeepSeekResponse("", "random_question")
	questions := []string{
		"Что меня ждет в любовных отношениях?",
		"Какие препятствия меня ожидают?",
		"Как улучшить финансовое положение?",
	}

	if err == nil {
		if generated, ok := response["questions"].([]string); ok {
			questions = append(questions, generated...)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"question": questions[time.Now().Unix()%int64(len(questions))],
	})
}

func generateDeepSeekResponse(input string, promptType string, extras ...string) (map[string]interface{}, error) {
	client := openai.NewClient(os.Getenv("DEEPSEEK_API_KEY"))

	var (
		prompt    string
		maxTokens = 500
	)

	switch promptType {
	case "daily":
		prompt = fmt.Sprintf(`CARTA TAROT: %s
Формат ответа ТОЛЬКО JSON:
{
	"forecast": "Прогноз на день...",
	"tags": ["тег1", "тег2", "тег3"],
	"do": ["совет1", "совет2", "совет3"],
	"dont": ["совет1", "совет2", "совет3"]
}`, input)

	case "question":
		prompt = fmt.Sprintf(`ВОПРОС: %s 
КАРТЫ: %s
Формат ответа ТОЛЬКО JSON:
{
	"answer": "Полный ответ...",
	"interpretation": "Толкование комбинации...",
	"tags": ["тег1", "тег2", "тег3"],
	"do": ["совет1", "совет2", "совет3"],
	"dont": ["совет1", "совет2", "совет3"]
}`, extras[0], input)

	case "random_question":
		prompt = `Сгенерируй 5 уникальных вопросов для гадания на Таро через запятую. 
Формат: {"questions": ["вопрос1", "вопрос2", ...]}`
		maxTokens = 300

	default:
		return nil, fmt.Errorf("unknown prompt type")
	}

	resp, err := client.CreateChatCompletion(c.Request.Context(), openai.ChatCompletionRequest{
		Model:     "deepseek-chat",
		Messages:  []openai.ChatCompletionMessage{{Role: "user", Content: prompt}},
		MaxTokens: maxTokens,
	})

	if err != nil {
		return nil, fmt.Errorf("API error: %v", err)
	}

	return parseDeepSeekResponse(resp.Choices[0].Message.Content, promptType)
}

func parseDeepSeekResponse(response string, promptType string) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("invalid response format")
	}

	switch promptType {
	case "daily":
		required := []string{"forecast", "tags", "do", "dont"}
		for _, key := range required {
			if _, exists := result[key]; !exists {
				return nil, fmt.Errorf("missing required field: %s", key)
			}
		}

	case "question":
		required := []string{"answer", "interpretation", "tags", "do", "dont"}
		for _, key := range required {
			if _, exists := result[key]; !exists {
				return nil, fmt.Errorf("missing required field: %s", key)
			}
		}

	case "random_question":
		if _, ok := result["questions"]; !ok {
			return nil, fmt.Errorf("invalid questions format")
		}

	default:
		return nil, fmt.Errorf("unknown prompt type")
	}

	return result, nil
}

func validateTelegramData(query map[string][]string, botToken string) bool {
	hash := query["hash"][0]
	delete(query, "hash")

	var keys []string
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var dataArr []string
	for _, k := range keys {
		dataArr = append(dataArr, fmt.Sprintf("%s=%s", k, query[k][0]))
	}
	dataStr := strings.Join(dataArr, "\n")

	secretKey := sha256.Sum256([]byte(botToken))
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataStr))
	return hex.EncodeToString(h.Sum(nil)) == hash
}

func generateTarotCards() []TarotCard {
	return []TarotCard{
		{
			ID:          uuid.New().String(),
			Name:        "Шут",
			Description: "Символ новых начинаний, невинности и спонтанности...",
			Keywords:    "новые начинания, свобода, риск",
			ImageURL:    "fool.jpg",
		},
		{
			ID:          uuid.New().String(),
			Name:        "Маг",
			Description: "Карта творческой силы и реализации потенциала...",
			Keywords:    "воля, мастерство, концентрация",
			ImageURL:    "magician.jpg",
		},
		//потом
	}
}

func decrementFreeReadings(userID string) int {
	// later
	return 2
}

func getCardsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, cardsDB)
}

func formatResponse(response map[string]interface{}) string {
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "Ошибка форматирования ответа"
	}
	return string(jsonData)
}
