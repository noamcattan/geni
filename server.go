package main

import (
	"github.com/gin-gonic/gin"
	"github.com/noamcattan/geni/ent"
	"log"
	"net/http"
)

type server struct {
	client *ent.Client
	*gin.Engine
}

func newServer(client *ent.Client) *server {
	r := gin.Default()
	s := &server{client: client, Engine: r}
	r.POST("/v1/account", s.createAccount)
	r.GET("/v1/account", s.getAccounts)
	r.POST("/v1/user", s.createUser)
	r.GET("/v1/user", s.getUsers)
	return s
}

// createAccount creates a new account.
func (s *server) createAccount(c *gin.Context) {
	var payload ent.Account
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	account, err := s.client.Account.Create().
		SetName(payload.Name).
		SetSheetsCredentials(payload.SheetsCredentials).
		SetSpreadsheetID(payload.SpreadsheetID).
		Save(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("account created successfully")
	c.JSON(http.StatusOK, gin.H{"id": account.ID})
}

func (s *server) getAccounts(c *gin.Context) {
	accounts, err := s.client.Account.Query().WithMember().All(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

// createUser creates a new user.
func (s *server) createUser(c *gin.Context) {
	var payload ent.User
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx, err := s.client.Tx(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := tx.User.Create().
		SetName(payload.Name).
		SetTelegramID(payload.TelegramID).
		Save(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = tx.Account.UpdateOneID(payload.Edges.Account.ID).
		AddMember(user).
		Save(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": user.ID})
}

func (s *server) getUsers(c *gin.Context) {
	users, err := s.client.User.Query().All(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}
