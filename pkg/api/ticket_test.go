package api

import (
	"context"
	"fmt"
	"testing"

	model "github.com/amankumarcs/trainticket/pkg/model/ticketing"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPurchaseTicket(t *testing.T) {
	server := NewTicketServiceServer()

	req := &model.PurchaseRequest{
		From: "City A",
		To:   "City B",
		User: &model.User{
			FirstName: "Alice",
			LastName:  "Doe",
			Email:     "alice@example.com",
		},
	}

	res, err := server.PurchaseTicket(context.Background(), req)
	if s, ok := status.FromError(err); ok {
		assert.NoError(t, err)
		assert.Equal(t, codes.OK, s.Code())
		assert.NotNil(t, res)
		assert.Equal(t, "1A", res.SeatNumber)
		assert.Equal(t, "A", res.Section)
		assert.Equal(t, "Ticket purchased successfully!", res.Message)
		assert.Equal(t, int32(1), res.TicketNumber)
		fmt.Printf("Error: %s, Code: %v\n", s.Message(), s.Code())
	}

}

func TestPurchaseTicketNotAvailable(t *testing.T) {
	server := NewTicketServiceServer()

	for i := 1; i <= 21; i++ {
		req := &model.PurchaseRequest{
			From: "City A",
			To:   "City B",
			User: &model.User{
				FirstName: "Alice",
				LastName:  "Doe",
				Email:     "alice@example.com",
			},
		}
		res, err := server.PurchaseTicket(context.Background(), req)
		if i <= 20 {
			assert.Equal(t, "Ticket purchased successfully!", res.Message)
			assert.Equal(t, int32(i), res.TicketNumber)
		} else {
			assert.Error(t, err)
			assert.Equal(t, "no available seats in any of sections", err.Error())
		}
	}
}

func TestGetReceipt(t *testing.T) {
	server := NewTicketServiceServer()

	ticketReq := &model.PurchaseRequest{
		From: "City A",
		To:   "City B",
		User: &model.User{
			FirstName: "Alice",
			LastName:  "Doe",
			Email:     "alice@example.com",
		},
	}
	_, _ = server.PurchaseTicket(context.Background(), ticketReq)

	req := &model.GetReceiptRequest{TicketNumber: 1}
	res, err := server.GetReceipt(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int32(1), res.Ticket.TicketNumber)

	//send invalid ticketid

	req = &model.GetReceiptRequest{TicketNumber: 122}
	res, err = server.GetReceipt(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("ticket not found for user: %d", 122), err.Error())
}

func TestViewUsersBySection(t *testing.T) {
	server := NewTicketServiceServer()

	for i := 1; i <= 3; i++ {
		_, _ = server.PurchaseTicket(context.Background(), &model.PurchaseRequest{
			From: "City A",
			To:   "City B",
			User: &model.User{
				FirstName: fmt.Sprintf("User%d", i),
				LastName:  "Doe",
				Email:     fmt.Sprintf("user%d@example.com", i),
			},
		})
	}

	req := &model.ViewUsersBySectionRequest{Section: "A"}
	res, err := server.ViewUsersBySection(context.Background(), req)

	assert.NoError(t, err)
	assert.Len(t, res.Tickets, 3)
}

func TestRemoveUser(t *testing.T) {
	server := NewTicketServiceServer()

	ticketReq := &model.PurchaseRequest{
		From: "City A",
		To:   "City B",
		User: &model.User{
			FirstName: "Alice",
			LastName:  "Doe",
			Email:     "alice@example.com",
		},
	}
	res, _ := server.PurchaseTicket(context.Background(), ticketReq)

	removeReq := &model.RemoveUserRequest{TicketNumber: res.TicketNumber}
	removeRes, err := server.RemoveUser(context.Background(), removeReq)

	assert.NoError(t, err)
	assert.Equal(t, "User removed successfully.", removeRes.Message)
}

func TestModifyUserSeat(t *testing.T) {
	server := NewTicketServiceServer()

	ticketReq := &model.PurchaseRequest{
		From: "City A",
		To:   "City B",
		User: &model.User{
			FirstName: "Alice",
			LastName:  "Doe",
			Email:     "alice@example.com",
		},
	}
	res, _ := server.PurchaseTicket(context.Background(), ticketReq)

	modifyReq := &model.ModifySeatRequest{
		TicketNumber:  res.TicketNumber,
		NewSeatNumber: "2A",
		NewSection:    "B",
	}
	modifyRes, err := server.ModifyUserSeat(context.Background(), modifyReq)

	assert.NoError(t, err)
	assert.Equal(t, "User seat modified successfully.", modifyRes.Message)

	receiptReq := &model.GetReceiptRequest{TicketNumber: res.TicketNumber}
	receiptRes, _ := server.GetReceipt(context.Background(), receiptReq)

	assert.Equal(t, "2A", receiptRes.Ticket.SeatNumber)
	assert.Equal(t, "B", receiptRes.Ticket.Section)
}
