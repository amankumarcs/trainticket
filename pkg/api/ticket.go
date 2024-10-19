package api

import (
	"context"
	"fmt"
	"sync"

	model "github.com/amankumarcs/trainticket/pkg/model/ticketing"
)

type TicketServiceServer struct {
	model.UnimplementedTicketServiceServer
	mu             sync.Mutex
	tickets        map[int32]*model.Ticket // Store tickets in memory
	availableSeats map[string][]string     // Available seats for each section
	sections       [2]string
}

// Constructor for TicketServiceServer
func NewTicketServiceServer() *TicketServiceServer {
	return &TicketServiceServer{
		tickets: make(map[int32]*model.Ticket),
		availableSeats: map[string][]string{
			"A": {"1A", "1B", "1C", "1D", "1E", "1F", "1G", "1H", "1I", "1J"},
			"B": {"2A", "2B", "2C", "2D", "2E", "2F", "2G", "2H", "2I", "2J"},
		},
		sections: [2]string{"A", "B"},
	}
}

// PurchaseTicket implementation
func (s *TicketServiceServer) PurchaseTicket(ctx context.Context, req *model.PurchaseRequest) (*model.PurchaseResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var section string

	for _, v := range s.sections {
		if len(s.availableSeats[v]) == 0 {
			continue
		} else {
			section = v
			break
		}
	}
	if section == "" {
		return nil, fmt.Errorf("no available seats in any of sections")
	}
	seat := s.availableSeats[section][0]
	s.availableSeats[section] = s.availableSeats[section][1:]

	ticket_number := int32(len(s.tickets) + 1)
	ticket := &model.Ticket{
		From:         req.From,
		To:           req.To,
		User:         req.User,
		PricePaid:    20.0,
		SeatNumber:   seat,
		Section:      section,
		TicketNumber: ticket_number,
	}

	// Store ticket in memory
	s.tickets[ticket_number] = ticket

	return &model.PurchaseResponse{
		TicketNumber: ticket_number,
		SeatNumber:   seat,
		Section:      section,
		Message:      "Ticket purchased successfully!",
	}, nil
}

// GetReceipt implementation
func (s *TicketServiceServer) GetReceipt(ctx context.Context, req *model.GetReceiptRequest) (*model.GetReceiptResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ticket, exists := s.tickets[req.TicketNumber]
	if !exists {
		return nil, fmt.Errorf("ticket not found for user: %d", req.TicketNumber)
	}

	return &model.GetReceiptResponse{
		Ticket: ticket,
	}, nil
}

// ViewUsersBySection implementation
func (s *TicketServiceServer) ViewUsersBySection(ctx context.Context, req *model.ViewUsersBySectionRequest) (*model.ViewUsersBySectionResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var tickets []*model.Ticket
	for _, ticket := range s.tickets {
		if ticket.Section == req.Section {
			tickets = append(tickets, ticket)
		}
	}

	return &model.ViewUsersBySectionResponse{
		Tickets: tickets,
	}, nil
}

// RemoveUser implementation
func (s *TicketServiceServer) RemoveUser(ctx context.Context, req *model.RemoveUserRequest) (*model.RemoveUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ticket, exists := s.tickets[req.TicketNumber]
	if exists {
		s.availableSeats[ticket.Section] = append(s.availableSeats[ticket.Section], ticket.SeatNumber)
		delete(s.tickets, req.TicketNumber)
		return &model.RemoveUserResponse{Message: "User removed successfully."}, nil
	}
	return &model.RemoveUserResponse{Message: "User not found."}, nil
}

// ModifyUserSeat implementation
func (s *TicketServiceServer) ModifyUserSeat(ctx context.Context, req *model.ModifySeatRequest) (*model.ModifySeatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ticket, exists := s.tickets[req.TicketNumber]
	if !exists {
		return nil, fmt.Errorf("ticket not found for ticket: %d", req.TicketNumber)
	}

	// Check if the new section has available seats
	if len(s.availableSeats[req.NewSection]) == 0 {
		return nil, fmt.Errorf("no available seats in section %s", req.NewSection)
	}

	// Add the old seat back to available seats
	s.availableSeats[ticket.Section] = append(s.availableSeats[ticket.Section], ticket.SeatNumber)

	// Allocate the new seat
	newSeat := s.availableSeats[req.NewSection][0]
	s.availableSeats[req.NewSection] = s.availableSeats[req.NewSection][1:]

	// Update the ticket
	ticket.SeatNumber = newSeat
	ticket.Section = req.NewSection

	return &model.ModifySeatResponse{Message: "User seat modified successfully."}, nil
}
