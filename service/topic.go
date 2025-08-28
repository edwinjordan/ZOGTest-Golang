package service

import (
	"context"

	"github.com/edwinjordan/ZOGTest-Golang.git/domain"
	"github.com/edwinjordan/ZOGTest-Golang.git/internal/logging"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type TopicRepository interface {
	CreateTopic(ctx context.Context, topic *domain.CreateTopicRequest) (*domain.Topic, error)
	GetTopicList(ctx context.Context, filter *domain.TopicFilter) ([]domain.Topic, error)
	GetTopic(ctx context.Context, id uuid.UUID) (*domain.Topic, error)
	UpdateTopic(ctx context.Context, id uuid.UUID, topic *domain.Topic) (*domain.Topic, error)
	DeleteTopic(ctx context.Context, id uuid.UUID) error
}

type TopicService struct {
	topicRepo TopicRepository
}

func NewTopicService(u TopicRepository) *TopicService {
	return &TopicService{
		topicRepo: u,
	}
}

// CreateTopic adds a new topic.
func (us *TopicService) CreateTopic(
	ctx context.Context,
	u *domain.CreateTopicRequest,
) (*domain.Topic, error) {
	createdTopic, err := us.topicRepo.CreateTopic(ctx, u)
	if err != nil {
		return nil, err
	}
	return createdTopic, nil
}

// GetTopic fetches a topic by ID.
func (us *TopicService) GetTopic(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Topic, error) {
	tracer := otel.Tracer("service.topic")
	ctxTrace, span := tracer.Start(ctx, "TopicService.GetTopic")
	defer span.End()

	topic, err := us.topicRepo.GetTopic(ctxTrace, id)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

// UpdateTopic updates name/email of an existing topic.
func (us *TopicService) UpdateTopic(
	ctx context.Context,
	id uuid.UUID,
	u *domain.Topic,
) (*domain.Topic, error) {

	existing, err := us.topicRepo.GetTopic(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, domain.ErrUserNotFound
	}

	existing.Name = u.Name
	existing.Slug = u.Slug

	_, err = us.topicRepo.UpdateTopic(ctx, id, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteUser removes a topic by ID.
func (us *TopicService) DeleteTopic(
	ctx context.Context,
	id uuid.UUID,
) error {

	topic, err := us.topicRepo.GetTopic(ctx, id)
	if err != nil {
		return err
	}
	if topic == nil {
		return domain.ErrUserNotFound
	}

	err = us.topicRepo.DeleteTopic(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (us *TopicService) GetTopicList(ctx context.Context, filter *domain.TopicFilter) ([]domain.Topic, error) {
	topics, err := us.topicRepo.GetTopicList(ctx, filter)
	if err != nil {
		logging.LogError(ctx, err, "get_topic_list_service")
		return nil, err
	}

	return topics, nil
}
