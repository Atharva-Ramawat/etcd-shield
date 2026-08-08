package main

import (
	"context"
	"log"
	"net"
	"strings"

	extauthz "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type authServer struct{}

func (s *authServer) Check(ctx context.Context, req *extauthz.CheckRequest) (*extauthz.CheckResponse, error) {
	path := "unknown"
	// Safely extract the path to avoid nil pointer panics
	if req != nil && req.Attributes != nil && req.Attributes.Request != nil && req.Attributes.Request.Http != nil {
		path = req.Attributes.Request.Http.Path
	}

	log.Printf("[Firewall] Intercepted traffic for path: %s", path)

	// The Virtual Patch
	if strings.HasPrefix(path, "/etcdserverpb.Maintenance/") {
		log.Println("[Firewall] 🚨 BLOCKED: Malicious Maintenance API call detected!")
		return &extauthz.CheckResponse{
			Status: &status.Status{Code: int32(codes.PermissionDenied)},
		}, nil
	}

	log.Println("[Firewall] ✅ ALLOWED: Safe traffic.")
	return &extauthz.CheckResponse{
		Status: &status.Status{Code: int32(codes.OK)},
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	extauthz.RegisterAuthorizationServer(s, &authServer{})

	log.Println("Shield Activated: Ext-Authz Firewall listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}