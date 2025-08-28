package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/edwinjordan/ZOGTest-Golang.git/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TopicRepository struct {
	Conn *pgxpool.Pool
}

func NewTopicRepository(conn *pgxpool.Pool) *TopicRepository {
	return &TopicRepository{Conn: conn}
}

func (u *TopicRepository) CreateTopic(ctx context.Context, topic *domain.CreateTopicRequest) (*domain.Topic, error) {
	query := `
		INSERT INTO topic (name, slug, password, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id`

	// hashedPassword, err := utils.HashPassword(Topic.Password)
	// if err != nil {
	// 	return nil, err
	// }

	var id uuid.UUID
	err := u.Conn.QueryRow(ctx, query, topic.Name, topic.Slug).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &domain.Topic{
		ID:   id.String(),
		Name: topic.Name,
		Slug: topic.Slug,
	}, nil
}

func (u *TopicRepository) GetTopicList(ctx context.Context, filter *domain.TopicFilter) ([]domain.Topic, error) {
	query := `
		SELECT
			u.id,
			u.name,
			u.slug,
            u.created_at,
            u.updated_at
		FROM topic u
        WHERE u.deleted_at is NULL`

	var args []interface{}
	var conditions []string
	if filter != nil && filter.Search != "" {
		conditions = append(conditions, `(u.name ILIKE $1 OR u.slug ILIKE $1)`)
		args = append(args, "%"+filter.Search+"%")
	}

	if len(conditions) > 0 {
		query += strings.Join(conditions, " AND ")
	}
	rows, err := u.Conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []domain.Topic

	for rows.Next() {
		var topic domain.Topic
		err := rows.Scan(
			&topic.ID,
			&topic.Name,
			&topic.Slug,
			&topic.CreatedAt,
			&topic.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}

	return topics, nil
}

func (u *TopicRepository) GetTopic(ctx context.Context, id uuid.UUID) (*domain.Topic, error) {
	// tracer := otel.Tracer("repo.Topic")
	// ctx, span := tracer.Start(ctx, "TopicRepository.GetTopic")
	// defer span.End()

	query := `
		SELECT
			id,
			name,
			slug,
			created_at,
			updated_at
		FROM topic
		WHERE id = $1 AND deleted_at IS NULL`

	// span.SetAttributes(attribute.String("query.statement", query))
	// span.SetAttributes(attribute.String("query.parameter", id.String()))
	row := u.Conn.QueryRow(ctx, query, id)

	var topic domain.Topic
	err := row.Scan(
		&topic.ID,
		&topic.Name,
		&topic.Slug,
		&topic.CreatedAt,
		&topic.UpdatedAt,
	)
	if err != nil {
		//span.RecordError(err)
		//		u.Metrics.TopicRepoCalls.WithLabelValues("GetTopic", "error").Inc()
		return nil, err
	}

	//	u.Metrics.TopicRepoCalls.WithLabelValues("GetTopic", "success").Inc()
	return &topic, nil
}

func (u *TopicRepository) UpdateTopic(ctx context.Context, id uuid.UUID, topic *domain.Topic) (*domain.Topic, error) {
	query := `
		UPDATE topic
		SET name = $1,
			slug = $2,
			updated_at = NOW()
		WHERE id = $3 AND deleted_at IS NULL
		RETURNING id, name, slug, created_at, updated_at`

	var updatedTopic domain.Topic
	err := u.Conn.QueryRow(ctx, query, topic.Name, topic.Slug, id).Scan(
		&updatedTopic.ID,
		&updatedTopic.Name,
		&updatedTopic.Slug,
		&updatedTopic.CreatedAt,
		&updatedTopic.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &updatedTopic, nil
}

func (u *TopicRepository) DeleteTopic(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE topic
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := u.Conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
