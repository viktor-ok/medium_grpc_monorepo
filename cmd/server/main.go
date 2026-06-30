package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"acme/gen/billingv1"
	"acme/gen/commonv1"
	"acme/gen/userv1"

	"google.golang.org/grpc"
)

// creamos un type userServer
// que embebed la interfaz UserServiceServer del servidor grpc
// a traves del type UnimplementedUserServiceServer que la implementa
// y tenemos acceso a todos sus metodos
type userServer struct {
	userv1.UnimplementedUserServiceServer
}

// GetUser sobreescribe el metodo GetUser del servidor grpc
// recibe como parametro un GetUserRequest mensaje protobuf
// y lo personalizamos para que nos devuelva un GetUserResponse mensaje protobuf
func (s *userServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	if req.GetId() == "" {
		return &userv1.GetUserResponse{
			Error: &commonv1.ErrorStatus{
				Code:    "INVALID_ARGUMENT",
				Message: "id is required",
			},
		}, nil
	}

	return &userv1.GetUserResponse{
		User: &userv1.User{
			Id:    req.GetId(),
			Email: fmt.Sprintf("%s@exmaple.com", req.GetId()),
		},
	}, nil
}

// -------------------------------------------------------------------------------------

// creamos un type billingServer
// que embebed la interfaz BillingServiceServer del servidor grpc
// a traves del type UnimplementedBillingServiceServer que la implementa
// y tenemos acceso a todos sus metodos
type billingServer struct {
	billingv1.UnimplementedBillingServiceServer
}

// GetInvoice sobreescribe el metodo GetInvoice del servidor grpc
// recibe como parametro un GetInvoiceRequest mensaje protobuf
// y lo personalizamos para que nos devuelva un GetInvoiceResponse mensaje protobuf
func (s *billingServer) GetInvoice(ctx context.Context, req *billingv1.GetInvoiceRequest) (*billingv1.GetInvoiceResponse, error) {
	if req.GetId() == "" {
		return &billingv1.GetInvoiceResponse{
			Error: &commonv1.ErrorStatus{
				Code:    "INVALID_ARGUMENT",
				Message: "id is required",
			},
		}, nil
	}

	return &billingv1.GetInvoiceResponse{
		Invoice: &billingv1.Invoice{
			Id:          req.GetId(),
			UserId:      "user-123",
			AmountCents: 4999,
			Currency:    "USD",
		},
	}, nil
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		// escuchamos por medio de la red local en el puerto 50051
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen on 50051: %v", err)
		}

		// creamos un servidor grpc
		grpcServer := grpc.NewServer()

		// registramos el servidor grpc y el tipo UserServer
		// que contiene el metodo GetUser como respuesta al cliente
		userv1.RegisterUserServiceServer(grpcServer, &userServer{})
		log.Println("UserService listening on :50051")

		// aceptamos las peticiones de los clientes
		// en al listener de la red local (lis)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("UserService failet: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		// escuchamos por medio de la red local en el puerto 50052
		lis, err := net.Listen("tcp", ":50052")
		if err != nil {
			log.Fatalf("failed to listen on 50052: %v", err)
		}

		// creamos un servidor grpc
		grpcServer := grpc.NewServer()

		// registramos el servidor grpc y el tipo billingServer
		// que contiene el metodo GetInvoice como respuesta al cliente
		billingv1.RegisterBillingServiceServer(grpcServer, &billingServer{})
		log.Println("BillingService listening on :50052")

		// aceptamos las peticiones de los clientes
		// en al listener de la red local (lis)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("BillingService failet: %v", err)
		}
	}()

	wg.Wait()
}
