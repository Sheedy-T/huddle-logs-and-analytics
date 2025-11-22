package grpc

import (
	"context"
	"log"

	// IMPORTS: Replace 'github.com/YOUR_USER/huddle-backend' with your actual module name from go.mod
	pb "github.com/Sheedy-T/huddle-backend/proto" 
)

// AnalyticsServer must embed the "Unimplemented" struct for forward compatibility
type AnalyticsServer struct {
	pb.UnimplementedAnalyticsServiceServer
}

// Notice the 'pb.' prefix before LogEventRequest and LogEventResponse
func (s *AnalyticsServer) LogEvent(ctx context.Context, req *pb.LogEventRequest) (*pb.LogEventResponse, error) {
	
	log.Printf("Received Event from User: %s", req.UserId)

	// Successful response
	return &pb.LogEventResponse{
		Success: true,
		Message: "Event logged successfully",
	}, nil
};