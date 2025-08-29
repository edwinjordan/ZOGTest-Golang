package service

import (
	"context"

	"github.com/edwinjordan/ZOGTest-Golang.git/domain"
	"github.com/google/uuid"
)

type NewsRepository interface {
	CreateNews(ctx context.Context, news *domain.CreateNewsRequest) (*domain.News, error)
	GetNewsList(ctx context.Context, filter *domain.NewsFilter) ([]domain.News, error)
	GetNews(ctx context.Context, id uuid.UUID) (*domain.News, error)
	UpdateNews(ctx context.Context, id uuid.UUID, news *domain.News) (*domain.News, error)
	DeleteNews(ctx context.Context, id uuid.UUID) error
}

type NewsService struct {
	newsRepo NewsRepository
}

func NewNewsService(n NewsRepository) *NewsService {
	return &NewsService{
		newsRepo: n,
	}
}

func (ns *NewsService) CreateNews(
	ctx context.Context,
	u *domain.CreateNewsRequest,
) (*domain.News, error) {
	createdNews, err := ns.newsRepo.CreateNews(ctx, u)
	if err != nil {
		return nil, err
	}
	return createdNews, nil
}

func (us *NewsService) GetNews(
	ctx context.Context,
	id uuid.UUID,
) (*domain.News, error) {
	news, err := us.newsRepo.GetNews(ctx, id)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func (us *NewsService) UpdateNews(
	ctx context.Context,
	id uuid.UUID,
	u *domain.News,
) (*domain.News, error) {

	existing, err := us.newsRepo.GetNews(ctx, id)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, domain.ErrUserNotFound
	}

	existing.Title = u.Title
	existing.Slug = u.Slug
	existing.Status = u.Status
	existing.Content = u.Content
	existing.Topics = u.Topics

	_, err = us.newsRepo.UpdateNews(ctx, id, existing)
	if err != nil {
		return nil, err
	}
	return existing, nil
}

func (us *NewsService) DeleteNews(
	ctx context.Context,
	id uuid.UUID,
) error {

	news, err := us.newsRepo.GetNews(ctx, id)
	if err != nil {
		return err
	}
	if news == nil {
		return domain.ErrUserNotFound
	}

	err = us.newsRepo.DeleteNews(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (us *NewsService) GetNewsList(ctx context.Context, filter *domain.NewsFilter) ([]domain.News, error) {
	newsList, err := us.newsRepo.GetNewsList(ctx, filter)
	if err != nil {
		return nil, err
	}
	return newsList, nil
}
