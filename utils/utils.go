package utils

import (
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("secret")

func HashPassword(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPassword(hash, pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) == nil
}

func GenerateJWT(id uuid.UUID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": id,
		"role":    role,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtSecret)
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		token, _ := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		claims := token.Claims.(jwt.MapClaims)

		c.Set("userID", claims["user_id"].(string))
		c.Set("role", claims["role"].(string))

		c.Next()
	}
}

func RoleMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")

		if userRole != role {
			c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}

type visitor struct {
	lastSeen time.Time
	count    int
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		userID, exists := c.Get("userID")
		key := ip

		if exists {
			key = userID.(string)
		}

		mu.Lock()

		v, found := visitors[key]

		if !found {
			visitors[key] = &visitor{time.Now(), 1}
			mu.Unlock()
			c.Next()
			return
		}

		if time.Since(v.lastSeen) > time.Minute {
			v.count = 0
			v.lastSeen = time.Now()
		}

		v.count++

		if v.count > 5 {
			mu.Unlock()
			c.AbortWithStatusJSON(429, gin.H{"error": "too many requests"})
			return
		}

		mu.Unlock()
		c.Next()
	}
}
