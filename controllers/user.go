package controllers

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

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	env, _ := helpers.GetEnv()

	dataUser, _ := c.Get("user")
	accountNumber := c.GetString("accountNo")

	type ResponseData struct {
		TraceID string `json:"traceId"`
		Data    struct {
			UID         int64   `json:"uid"`
			Balance     float64 `json:"balance"`
			AccountName string  `json:"accountName"`
			CreateTime  int64   `json:"createTime"`
			AccountNo   string  `json:"accountNo"`
			UpdateTime  int64   `json:"updateTime"`
			ID          int64   `json:"id"`
			Status      string  `json:"status"`
		} `json:"data"`
		Success bool `json:"success"`
	}

	type Payload struct {
		AccountNumber string `json:"accountNo"`
	}

	fmt.Println("Account No: " + accountNumber)
	payload := Payload{
		AccountNumber: accountNumber,
	}

	bodyRequest, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error saat mengubah data menjadi JSON:", err)
		return
	}

	endpointAPI := fmt.Sprintf("%s/bankAccount/info", env.API)
	requestAPI, err := http.NewRequest("POST", endpointAPI, bytes.NewBuffer(bodyRequest))
	if err != nil {
		fmt.Println("Error saat membuat request:", err)
		return
	}
	tokenAPI := c.GetString("tokenAPI")
	// Set header untuk mengatur tipe konten menjadi "application/json"
	requestAPI.Header.Set("Content-Type", "application/json")
	requestAPI.Header.Set("Authorization", "Bearer "+tokenAPI)

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
	fmt.Println("========= HIT INFO BANK ACCOUNT =========")
	fmt.Println("Response:\n", string(responseData))
	fmt.Println("===========================")

	// Mendeklarasikan variabel untuk menampung data respons JSON
	var res ResponseData

	// Menguraikan data JSON ke dalam struktur data ResponseData
	err = json.Unmarshal(responseData, &res)
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success get data user",
		"saldo":   res.Data.Balance,
		"data":    dataUser,
	})
}

func UpdateUser(c *gin.Context) {
	type requestBody struct {
		Name  string
		Email string
	}

	var json requestBody
	err := c.ShouldBindJSON(&json)
	fmt.Println("Debugging: ", json.Name == "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "data must be json",
		})
		return
	}

	if json.Email == "" || json.Name == "" {
		var message []string
		if json.Email == "" {
			message = append(message, "data email is required")
		}

		if json.Name == "" {
			message = append(message, "data name is required")
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"errors":  message,
			"message": "failed request body",
		})
		return
	}

	userId, _ := c.Get("userId")
	var modelUser models.User
	err = initializers.DB.Not("id = ?", userId).Where("email = ?", json.Email).First(&modelUser).Error
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  false,
			"message": "email was created by other user",
		})
		return
	}
	initializers.DB.Where("id = ?", userId).First(&modelUser)
	modelUser.Name = json.Name
	modelUser.Email = json.Email
	initializers.DB.Save(&modelUser)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success to updated data user",
		"data":    modelUser,
	})
}

func SetupPin(c *gin.Context) {
	type request struct {
		PinCode string `json:"pin_code"`
	}

	var bodyRequest request
	err := c.ShouldBindJSON(&bodyRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "data must be json",
		})
		return
	}

	if bodyRequest.PinCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "pin code required",
		})
		return
	}

	userId, _ := c.Get("userId")
	var modelUser models.User
	initializers.DB.Where("id = ?", userId).First(&modelUser)
	modelUser.PinNumber = bodyRequest.PinCode
	initializers.DB.Save(&modelUser)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success to set pin",
	})
	return
}

func VerifyPinNumber(c *gin.Context) {
	type request struct {
		PinCode string `json:"pin_code"`
	}

	var bodyRequest request
	err := c.ShouldBindJSON(&bodyRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "data must be json",
		})
		return
	}

	if bodyRequest.PinCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "pin code required",
		})
		return
	}

	userId, _ := c.Get("userId")
	var modelUser models.User
	initializers.DB.Where("id = ?", userId).First(&modelUser)
	if modelUser.PinNumber != bodyRequest.PinCode {
		c.JSON(http.StatusOK, gin.H{
			"status":  false,
			"message": "pin code is wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "pin code is is true",
	})
}

func ActivationBankAccount(c *gin.Context) {
	userId, _ := c.Get("userId")
	var modelUser models.User
	initializers.DB.Where("id = ?", userId).First(&modelUser)
	if modelUser.AccountName != "" || modelUser.AccountNumber != "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  true,
			"message": "bank already actived",
		})
		return
	}

	type ResponseData struct {
		TraceID string `json:"traceId"`
		Data    struct {
			UID         int64  `json:"uid"`
			Balance     int64  `json:"balance"`
			AccountName string `json:"accountName"`
			CreateTime  int64  `json:"createTime"`
			AccountNo   string `json:"accountNo"`
			UpdateTime  int64  `json:"updateTime"`
			ID          int64  `json:"id"`
			Status      string `json:"status"`
		} `json:"data"`
		Success bool `json:"success"`
	}

	env, _ := helpers.GetEnv()
	type Payload struct {
		Balance int `json:"balance"`
	}
	payload := Payload{
		Balance: 50000,
	}

	bodyRequest, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error saat mengubah data menjadi JSON:", err)
		return
	}

	endpointAPI := fmt.Sprintf("%s/bankAccount/create", env.API)
	requestAPI, err := http.NewRequest("POST", endpointAPI, bytes.NewBuffer(bodyRequest))
	if err != nil {
		fmt.Println("Error saat membuat request:", err)
		return
	}
	tokenAPI := c.GetString("tokenAPI")
	// Set header untuk mengatur tipe konten menjadi "application/json"
	requestAPI.Header.Set("Content-Type", "application/json")
	requestAPI.Header.Set("Authorization", "Bearer "+tokenAPI)

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
	fmt.Println("========= HIT API CREATE BANK =========")
	fmt.Println("Response:\n", string(responseData))
	fmt.Println("===========================")

	// Mendeklarasikan variabel untuk menampung data respons JSON
	var res ResponseData

	// Menguraikan data JSON ke dalam struktur data ResponseData
	err = json.Unmarshal(responseData, &res)

	modelUser.AccountName = strings.ToUpper(modelUser.Name)
	modelUser.AccountNumber = res.Data.AccountNo
	initializers.DB.Save(&modelUser)
	c.JSON(http.StatusUnauthorized, gin.H{
		"status":  false,
		"message": "Success activation bank account",
	})
}
