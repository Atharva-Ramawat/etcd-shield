package proxy

import (
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TransparentHandler creates a handler that proxies all unknown gRPC requests to the backend.
func TransparentHandler(backendAddr string) grpc.StreamHandler {
	return func(srv interface{}, serverStream grpc.ServerStream) error {
		fullMethodName, ok := grpc.MethodFromServerStream(serverStream)
		if !ok {
			return grpc.Errorf(grpc.Code(grpc.ErrClientConnClosing), "missing method name")
		}

		// Connect to backend
		backendConn, err := grpc.DialContext(serverStream.Context(), backendAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.CallContentSubtype((&proxyCodec{}).Name())),
		)
		if err != nil {
			return err
		}
		defer backendConn.Close()

		clientStream, err := backendConn.NewStream(serverStream.Context(), &grpc.StreamDesc{
			ServerStreams: true,
			ClientStreams: true,
		}, fullMethodName)
		if err != nil {
			return err
		}

		// Proxy data bidirectionally
		errc := make(chan error, 1)
		go func() {
			errc <- forwardStream(serverStream, clientStream)
		}()
		go func() {
			errc <- forwardStream(clientStream, serverStream)
		}()

		err = <-errc
		if err != nil && err != io.EOF {
			log.Printf("Proxy error for %s: %v", fullMethodName, err)
			return err
		}
		return nil
	}
}

func forwardStream(src grpc.Stream, dst grpc.Stream) error {
	var f []byte
	for {
		if err := src.RecvMsg(&f); err != nil {
			return err
		}
		if err := dst.SendMsg(f); err != nil {
			return err
		}
	}
}

// proxyCodec allows raw byte manipulation without protobuf decoding
type proxyCodec struct{}

func (proxyCodec) Marshal(v interface{}) ([]byte, error) { return v.([]byte), nil }
func (proxyCodec) Unmarshal(data []byte, v interface{}) error {
	*v.(*[]byte) = append(*v.(*[]byte), data...)
	return nil
}
func (proxyCodec) Name() string { return "proto" }