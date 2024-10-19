# Ticketing Service

## Overview

The Ticketing Service is a gRPC-based API designed for managing ticket purchases, receipts, user views by section, and seat modifications. This service enables users to purchase tickets for events, retrieve their receipts, view users in specific sections, remove users, and modify seat assignments.

## Features

- **Purchase Ticket**: Allows users to purchase tickets for an event.
- **Get Receipt**: Enables users to retrieve the receipt for a purchased ticket.
- **View Users by Section**: Lists all users and their tickets in a specific section.
- **Remove User**: Removes a user and frees up their assigned seat.
- **Modify User Seat**: Changes a user's seat assignment.

## Requirements

- Go 1.20+
- gRPC
- Protobuf

## Installation

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/amankumarcs/trainticket.git
   cd trainticket
   

2. **Install Dependencies**:

   ```bash
   go mod tidy

## Start
1. **Server Start**:
   
   
   go inside cmd 
   ```bash
   go run main.go

2. **Run Unit Test Start**:
   
   
   go inside cmd 
   ```bash
   go test -v





