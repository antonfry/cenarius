package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *server) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("ip")
		if len(values) > 0 {
			// ключ содержит слайс строк, получаем первую строку
			ipStr := values[0]
			ip := net.ParseIP(ipStr)
			if ip == nil {
				return nil, errors.New("no IP in metadata")
			}
			if !s.allowedSubnet.Contains(ip) {
				return nil, fmt.Errorf("IP from X-Real-IP header is not from %s %s", ip.String(), s.allowedSubnet.String())
			}
		}
	}
	// вызываем RPC-метод
	result, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}
	return result, nil
}
