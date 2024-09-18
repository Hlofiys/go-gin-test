package controllers

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
	"strconv"
	"time"

	db "go-gin-test/db/sqlc"
	"go-gin-test/schemas"

	"github.com/gin-gonic/gin"
)

type ContactController struct {
	db  *db.Queries
	ctx context.Context
}

func NewContactController(db *db.Queries, ctx context.Context) *ContactController {
	return &ContactController{db, ctx}
}

// CreateContact godoc
//
//	@Summary		Create new contact
//	@Tags			contacts
//	@Accept			json
//	@Produce		json
//	@Param			contact	body		schemas.CreateContact	true	"Create contact"
//	@Success		200		{object}	db.Contact
//	@Router			/contacts [post]
func (cc *ContactController) CreateContact(ctx *gin.Context) {
	var payload *schemas.CreateContact

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	now := time.Now()
	args := &db.CreateContactParams{
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		PhoneNumber: payload.PhoneNumber,
		Street:      payload.Street,
		CreatedAt:   pgtype.Timestamp{Time: now, Valid: true},
		UpdatedAt:   pgtype.Timestamp{Time: now, Valid: true},
	}

	contact, err := cc.db.CreateContact(ctx, *args)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving contact", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "successfully created contact", "contact": contact})
}

// UpdateContact godoc
//
//	@Summary		Update contact
//	@Tags			contacts
//	@Accept			json
//	@Produce		json
//	@Param          contactId   path      	int  					true  	"Contact Id"
//	@Param			contact		body		schemas.UpdateContact	true	"Update contact"
//	@Success		200		{object}	db.Contact
//	@Router			/contacts/{contactId} [patch]
func (cc *ContactController) UpdateContact(ctx *gin.Context) {
	var payload *schemas.UpdateContact
	contactId := ctx.Param("contactId")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	id, err := strconv.Atoi(contactId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed to cast id string to int", "error": err.Error()})
		return
	}

	now := time.Now()
	args := &db.UpdateContactParams{
		ContactID:   int32(id),
		FirstName:   pgtype.Text{String: payload.FirstName, Valid: payload.FirstName != ""},
		LastName:    pgtype.Text{String: payload.LastName, Valid: payload.LastName != ""},
		PhoneNumber: pgtype.Text{String: payload.PhoneNumber, Valid: payload.PhoneNumber != ""},
		Street:      pgtype.Text{String: payload.PhoneNumber, Valid: payload.Street != ""},
		UpdatedAt:   pgtype.Timestamp{Time: now, Valid: true},
	}

	contact, err := cc.db.UpdateContact(ctx, *args)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve contact with this ID"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving contact", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "successfully updated contact", "contact": contact})
}

// GetContactById godoc
//
//		@Summary		Get contact by id
//		@Tags			contacts
//		@Accept			json
//		@Produce		json
//	 	@Param          contactId   path      int  true  "Contact Id"
//		@Success		200		{object}	db.Contact
//		@Router			/contacts/{contactId} [get]
func (cc *ContactController) GetContactById(ctx *gin.Context) {
	contactId := ctx.Param("contactId")

	id, err := strconv.Atoi(contactId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed to cast id string to int", "error": err.Error()})
		return
	}

	contact, err := cc.db.GetContactById(ctx, int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve contact with this ID"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving contact", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrived id", "contact": contact})
}

// GetAllContacts godoc
//
//	@Summary		Get all contacts
//	@Tags			contacts
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	[]db.Contact
//	@Router			/contacts [get]
func (cc *ContactController) GetAllContacts(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	reqPageID, _ := strconv.Atoi(page)
	reqLimit, _ := strconv.Atoi(limit)
	offset := (reqPageID - 1) * reqLimit

	args := &db.ListContactsParams{
		Limit:  int32(reqLimit),
		Offset: int32(offset),
	}

	contacts, err := cc.db.ListContacts(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed to retrieve contacts", "error": err.Error()})
		return
	}

	if contacts == nil {
		contacts = []db.Contact{}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrieved all contacts", "size": len(contacts), "contacts": contacts})
}

// DeleteContactById godoc
//
//	@Summary		Delete contact by id
//	@Tags			contacts
//	@Accept			json
//	@Produce		json
//	@Param          contactId   path      	int  					true  	"Contact Id"
//	@Success		204
//	@Router			/contacts/{contactId} [delete]
func (cc *ContactController) DeleteContactById(ctx *gin.Context) {
	contactId := ctx.Param("contactId")

	id, err := strconv.Atoi(contactId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed to cast id string to int", "error": err.Error()})
		return
	}

	_, err = cc.db.GetContactById(ctx, int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve contact with this ID"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving contact", "error": err.Error()})
		return
	}

	err = cc.db.DeleteContact(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"status": "successfuly deleted"})

}
