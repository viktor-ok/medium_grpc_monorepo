package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"acme/gen/billingv1"
	"acme/gen/userv1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// creamos un cliente grpc
	// que hace llamadas a un servidor es este puerto
	userConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("dial user: %v", err)
	}
	defer userConn.Close()

	// creamos UserClient del servidor grpc y pasamos la conexion
	userClient := userv1.NewUserServiceClient(userConn)

	// el cliente hace la peticion al endpoint GetUser del sevicio grpc
	// y pasa el GetUserRequest type struct del mensaje protobuf .pb.go (msg)
	userRes, err := userClient.GetUser(ctx, &userv1.GetUserRequest{Id: "user-123"})
	if err != nil {
		log.Fatalf("GetUser: %v", err)
	}
	// log.Printf("UserService.GetUser => user=%+v error=%+v", userRes.GetUser(), userRes.GetError())

	_, err = fmt.Fprintf(
		os.Stdout,
		`{"id": %q, "email": %q, "error": %q}`,
		userRes.GetUser().GetId(),
		userRes.GetUser().GetEmail(),
		userRes.GetError(),
	)
	if err != nil {
		log.Fatalf("error Fprintf: %v", err)
	}

	// --------------------------------------------------------------------------------------------------

	// creamos un cliente grpc
	// que hace llamadas a un servidor es este puerto
	billingConn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("dial billing: %v", err)
	}
	defer billingConn.Close()

	// creamos billignClient del servidor grpc y pasamos la conexion
	billignClient := billingv1.NewBillingServiceClient(billingConn)

	// el cliente hace la peticion al endpoint GetInvoice del sevicio grpc
	// y pasa el GetInvoiceRequest type struct del mensaje protobuf .pb.go (msg)
	invRes, err := billignClient.GetInvoice(ctx, &billingv1.GetInvoiceRequest{Id: "inv-999"})
	if err != nil {
		log.Fatalf("GetInvoice: %v", err)
	}
	// log.Printf("BillingService.GetInvoice => invoice=%+v error=%+v", invRes.GetInvoice(), invRes.GetError())

	_, err = fmt.Fprintf(
		os.Stdout,
		`{"id": %q, "user_id": %q, "amount_cents": %v, "currency": %q, "error": %q}`,
		invRes.GetInvoice().GetId(),
		invRes.GetInvoice().GetUserId(),
		invRes.GetInvoice().GetAmountCents(),
		invRes.GetInvoice().GetCurrency(),
		invRes.GetError(),
	)
	if err != nil {
		log.Fatalf("error Fprintf: %v", err)
	}
}
