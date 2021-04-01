package grpcconnect

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	pb "github.com/minmax1996/aoimdb/api/proto/command"
	db "github.com/minmax1996/aoimdb/internal/aoimdb"
	"github.com/minmax1996/aoimdb/logger"
)

type Server struct {
	pb.UnimplementedDatabaseControllerServer
}

func HandleContextAuthentication(ctx context.Context) error {
	logger.Warn("CONMETESAEDW")
	return nil
	userPass, ok := ctx.Value("userPass").(string)
	if !ok {
		return errors.New("context has no userCredentials")
	}
	userPassParts := strings.SplitN(userPass, ":", 2)
	if len(userPassParts) < 2 {
		return errors.New("not enought args to process auth")
	}
	return db.AuthentificateByUserPass(userPassParts[0], userPassParts[1])
}

func (s *Server) SelectDatabase(ctx context.Context, req *pb.SelectDatabaseRequest) (*pb.SelectDatabaseResponse, error) {
	if err := HandleContextAuthentication(ctx); err != nil {
		return nil, err
	}

	db.SelectDatabase(req.Name)
	return &pb.SelectDatabaseResponse{
		Name: req.Name,
	}, nil
}

func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	if err := HandleContextAuthentication(ctx); err != nil {
		return nil, err
	}

	var err error
	db.SelectDatabase(req.DatabaseName)
	val, err := db.Get(req.DatabaseName, req.Key)
	if err != nil {
		return nil, err
	}
	var returnValue []byte

	switch data := val.(type) {
	case string:
		returnValue = []byte(data)
	default:
		returnValue, err = json.Marshal(data)
	}

	if err != nil {
		return nil, err
	}

	return &pb.GetResponse{
		Key:   req.Key,
		Value: returnValue,
	}, nil
}

func (s *Server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	if err := HandleContextAuthentication(ctx); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Server) GetKeys(ctx context.Context, req *pb.GetKeysRequest) (*pb.GetKeysResponse, error) {
	if err := HandleContextAuthentication(ctx); err != nil {
		return nil, err
	}

	result, err := db.GetKeys(req.DatabaseName, req.KeyPattern)
	if err != nil {
		return nil, err
	}

	return &pb.GetKeysResponse{
		Keys: result,
	}, nil
}
