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
	"github.com/google/uuid"
)

func CreateTransaction(c *gin.Context) {
	type bodyRequest struct {
		Amount int `json:"amount"`
	}

	var payload bodyRequest

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "data must be json",
		})
		return
	}

	uuidReceiver, _ := uuid.Parse(c.Param("uuid"))
	uuidSender, _ := uuid.Parse(c.GetString("userId"))

	var userSender models.User
	err = initializers.DB.Where("id = ?", uuidSender.String()).First(&userSender).Error
	if userSender.LimitCredit < (uint64(payload.Amount) + userSender.TotalCredit) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "failed to create transaction",
			"errors":  "limit credit tidak mencukupi",
		})
		return
	}

	createTransaction := models.Transaction{
		ReceiverID:      uuidReceiver,
		SenderID:        uuidSender,
		Amount:          payload.Amount,
		TransactionCode: strings.ToUpper("TRX-" + helpers.RandomString(5) + "-" + helpers.RandomString(5) + "-" + helpers.RandomString(5)),
	}

	result := initializers.DB.Create(&createTransaction)

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
	type PayloadAPI struct {
		SenderAccountNumber   string `json:"senderAccountNo"`
		ReceiverAccountNumber string `json:"receiverAccountNo"`
		Amount                int    `json:"amount"`
	}

	var userReceiver models.User

	fmt.Println("UUID Receiver: " + uuidReceiver.String())
	err = initializers.DB.Where("id = ?", uuidReceiver.String()).First(&userReceiver).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "failed to create transaction",
			"errors":  "code uuid not valid",
		})
		return
	}

	requestPayloadAPI := PayloadAPI{
		SenderAccountNumber:   userSender.AccountNumber,
		ReceiverAccountNumber: userReceiver.AccountNumber,
		Amount:                payload.Amount,
	}

	fmt.Println("request: " + requestPayloadAPI.ReceiverAccountNumber)

	bodyReq, err := json.Marshal(requestPayloadAPI)
	if err != nil {
		fmt.Println("Error saat mengubah data menjadi JSON:", err)
		return
	}

	endpointAPI := fmt.Sprintf("%s/bankAccount/transaction/create", env.API)
	requestAPI, err := http.NewRequest("POST", endpointAPI, bytes.NewBuffer(bodyReq))
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
	fmt.Println("========= HIT API CREATE TRANSACTION =========")
	fmt.Println("Response:\n", string(responseData))
	fmt.Println("===========================")

	// Mendeklarasikan variabel untuk menampung data respons JSON
	var res ResponseData

	// Menguraikan data JSON ke dalam struktur data ResponseData
	err = json.Unmarshal(responseData, &res)

	userSender.TotalCredit += uint64(payload.Amount)
	initializers.DB.Save(&userSender)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "failed to create transaction",
			"errors":  result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success create transaction",
		"data":    createTransaction,
	})
}

func GetListTransaction(c *gin.Context) {
	type DataTransaction struct {
		IN  []models.Transaction
		OUT []models.Transaction
	}

	userId := c.GetString("userId")

	var transIn []models.Transaction
	var transOut []models.Transaction
	_ := initializers.DB.Where("receiver_id = ?", userId).Find(&transIn)

	_ = initializers.DB.Where("sender_id = ?", userId).Find(&transOut)

	dataTrans := DataTransaction{
		IN:  transIn,
		OUT: transOut,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Berhasil get data transaction",
		"data":    dataTrans,
	})
	return
}
