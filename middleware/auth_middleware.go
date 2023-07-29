package middleware

import (
	"backend-technoscape/helpers"
	"backend-technoscape/initializers"
	"backend-technoscape/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func Auth(c *gin.Context) {
	fmt.Println("========= IN MIDDLEWARE AUTH =============")

	tokenstring := c.GetHeader("Authorization")
	tokenstring = strings.Replace(tokenstring, "Bearer ", "", -1)

	fmt.Println("Tokenstring: " + tokenstring)

	token, err := jwt.Parse(string(tokenstring), func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected sigining method: %v", t.Header["alg"])
		}

		return []byte(viper.Get("JWT_SECRET").(string)), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		var user models.User
		err = initializers.DB.Where("id = ?", claims["sub"]).First(&user).Error
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		env, _ := helpers.GetEnv()

		type Payload struct {
			Username string `json:"username"`
			Password string `json:"loginPassword"`
		}

		type ResponseData struct {
			TraceID string `json:"traceId"`
			Data    struct {
				AccessToken   string `json:"accessToken"`
				Balance       int    `json:"balance"`
				AccountName   string `json:"accountName"`
				AccountNumber string `json:"accountNo"`
			} `json:"data"`
			Success bool   `json:"success"`
			Error   string `json:"errMsg"`
		}

		payload := Payload{
			Username: env.Username,
			Password: env.Password,
		}

		bodyRequest, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error saat mengubah data menjadi JSON:", err)
			return
		}

		endpointAPI := fmt.Sprintf("%s/user/auth/token", env.API)
		requestAPI, err := http.NewRequest("POST", endpointAPI, bytes.NewBuffer(bodyRequest))
		if err != nil {
			fmt.Println("Error saat membuat request:", err)
			return
		}

		// Set header untuk mengatur tipe konten menjadi "application/json"
		requestAPI.Header.Set("Content-Type", "application/json")

		// Membuat client HTTP untuk mengirim request
		client := http.DefaultClient

		// Mengirimkan request POST
		resp, err := client.Do(requestAPI)
		if err != nil {
			fmt.Println("Error saat mengirimkan request:", err)
			return
		}

		defer resp.Body.Close()

		// Membaca respons dari API pihak ketiga (opsional)
		// Jika Anda tidak memerlukan respons, Anda dapat menghapus bagian ini
		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error saat membaca respons:", err)
			return
		}
		fmt.Println("========= HIT API LOGIN =========")
		fmt.Println("Response:\n", string(responseData))
		fmt.Println("===========================")

		// Mendeklarasikan variabel untuk menampung data respons JSON
		var res ResponseData
		err = json.Unmarshal(responseData, &res)

		c.Set("user", user)
		c.Set("userId", claims["sub"])
		c.Set("accountNo", user.AccountNumber)
		c.Set("tokenAPI", res.Data.AccessToken)
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
