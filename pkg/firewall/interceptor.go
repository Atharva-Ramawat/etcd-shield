package firewall

import (
	"context"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SecurityFirewallInterceptor blocks malicious requests based on your config
func SecurityFirewallInterceptor(blockMaintenance bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Printf("[Firewall] Intercepted request to: %s", info.FullMethod)

		if blockMaintenance && strings.HasPrefix(info.FullMethod, "/etcdserverpb.Maintenance/") {
			log.Printf("[Firewall] BLOCKED: Unauthorized access to %s", info.FullMethod)
			return nil, status.Errorf(codes.PermissionDenied, "Firewall blocked Maintenance API call")
		}

		return handler(ctx, req)
	}
}