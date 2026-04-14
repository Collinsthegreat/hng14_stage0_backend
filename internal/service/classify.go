package service

import (
	"context"
	"errors"
	"time"

	"github.com/Collinsthegreat/hng14_stage0_backend/internal/client"
	"github.com/Collinsthegreat/hng14_stage0_backend/internal/model"
)

var ErrNoPrediction = errors.New("No prediction available for the provided name")

type ClassifyService interface {
	Classify(ctx context.Context, name string) (*model.ClassifyResponseData, error)
}

type classifyService struct {
	client client.GenderizeClient
}

func NewClassifyService(client client.GenderizeClient) ClassifyService {
	return &classifyService{client: client}
}

func (s *classifyService) Classify(ctx context.Context, name string) (*model.ClassifyResponseData, error) {
	resp, err := s.client.Predict(ctx, name)
	if err != nil {
		return nil, err
	}

	if resp.Gender == nil || resp.Count == 0 {
		return nil, ErrNoPrediction
	}

	isConfident := resp.Probability >= 0.7 && resp.Count >= 100

	return &model.ClassifyResponseData{
		Name:        resp.Name,
		Gender:      *resp.Gender,
		Probability: resp.Probability,
		SampleSize:  resp.Count,
		IsConfident: isConfident,
		ProcessedAt: time.Now().UTC(),
	}, nil
}
