package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/edwinjordan/ZOGTest-Golang.git/domain"
	"github.com/edwinjordan/ZOGTest-Golang.git/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsRepository struct {
	Conn *pgxpool.Pool
}

func NewNewsRepository(conn *pgxpool.Pool) *NewsRepository {
	return &NewsRepository{Conn: conn}
}

func (u *NewsRepository) CreateNews(ctx context.Context, news *domain.CreateNewsRequest) (*domain.News, error) {
	query := `
		INSERT INTO news (title, slug, status, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id`

	var id uuid.UUID

	//var createdAt, updatedAt string
	err := u.Conn.QueryRow(ctx, query, news.Title, utils.Slugify(news.Title), news.Status, news.Content).Scan(&id)
	if err != nil {
		return nil, err
	}
	//panic(news.Topic)
	for _, detail := range news.Topic {

		topicID, err := uuid.Parse(detail.TopicId)
		if err != nil {
			return nil, errors.New("invalid topic ID: " + detail.TopicId)
		}

		topicQuery := `
			INSERT INTO news_topic (news_id, topic_id, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW())`

		_, err = u.Conn.Exec(ctx, topicQuery, id, topicID)
		if err != nil {
			return nil, err
		}
	}
	//createdNews.Details = make([]domain.NewsDetail, 0, len(news.Details))

	return &domain.News{
		ID:      id.String(),
		Title:   news.Title,
		Slug:    utils.Slugify(news.Title),
		Status:  news.Status,
		Content: news.Content,
		//CreatedAt: createdAt,
		//UpdatedAt: updatedAt,
	}, nil
}

func (u *NewsRepository) GetNewsList(ctx context.Context, filter *domain.NewsFilter) ([]domain.News, error) {
	query := `
		SELECT
			n.id,
			n.title,
			n.slug,
			n.status,
			n.content,
			n.created_at,
			n.updated_at,
			 (
                SELECT COALESCE(json_agg(
                    json_build_object(
                        'id', t.id,
                        'name', t.name,
                        'slug', t.slug
                    )
                ), null)
                FROM news_topic nt
                LEFT JOIN topik t ON t.id = nt.topic_id
				WHERE nt.news_id = n.id AND t.deleted_at IS NULL
			) as topics_list
		FROM news n
		WHERE n.deleted_at is NULL`

	var args []interface{}
	var conditions []string
	if filter != nil && filter.Search != "" {
		conditions = append(conditions, "(n.title ILIKE $1 OR n.content ILIKE $1)")
		args = append(args, "%"+filter.Search+"%")
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	rows, err := u.Conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsList []domain.News
	for rows.Next() {
		var news domain.News
		//var topicsJSON sql.NullString
		err := rows.Scan(
			&news.ID,
			&news.Title,
			&news.Slug,
			&news.Status,
			&news.Content,
			&news.CreatedAt,
			&news.UpdatedAt,
			&news.TopicList,
		)
		if err != nil {
			return nil, err
		}

		//	if topicsJSON.Valid && topicsJSON.String != "" {
		//		if err := json.Unmarshal([]byte(topicsJSON.String), &news.Topics); err != nil {
		//			return nil, err
		//		}
		//	} else {
		//		news.Topics = []domain.NewsTopic{}
		//	}
		newsList = append(newsList, news)
	}

	return newsList, nil
}
func (u *NewsRepository) GetNews(ctx context.Context, id uuid.UUID) (*domain.News, error) {
	query := `
		SELECT
			id,
			title,
			slug,
			status,
			content,
			created_at,
			updated_at,
			 (
                SELECT COALESCE(json_agg(
                    json_build_object(
                        'id', t.id,
                        'name', t.name,
                        'slug', t.slug
                    )
                ), null)
                FROM news_topic nt
                LEFT JOIN topik t ON t.id = nt.topic_id
				WHERE nt.news_id = n.id AND t.deleted_at IS NULL
			) as topics
		FROM news as n
		WHERE n.id = $1 AND n.deleted_at IS NULL`

	var news domain.News
	err := u.Conn.QueryRow(ctx, query, id).Scan(
		&news.ID,
		&news.Title,
		&news.Slug,
		&news.Status,
		&news.Content,
		&news.CreatedAt,
		&news.UpdatedAt,
		&news.Topics,
	)
	if err != nil {
		return nil, err
	}

	return &news, nil
}

func (u *NewsRepository) UpdateNews(ctx context.Context, id uuid.UUID, news *domain.News) (*domain.News, error) {
	query := `
		UPDATE news
		SET title = $1,
			slug = $2,
			status = $3,
			content = $4,
			updated_at = NOW()
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING id, title, slug, status, content, updated_at`

	var updatedNews domain.News
	err := u.Conn.QueryRow(ctx, query, news.Title, utils.Slugify(news.Title), news.Status, news.Content, id).Scan(
		&updatedNews.ID,
		&updatedNews.Title,
		&updatedNews.Slug,
		&updatedNews.Status,
		&updatedNews.Content,
		&updatedNews.UpdatedAt,
	//	&updatedNews.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound

		}
		return nil, err
	}

	query_topik := `
		DELETE FROM news_topic
		WHERE news_id = $1 `

	_, err2 := u.Conn.Exec(ctx, query_topik, id)
	if err2 != nil {
		return nil, err2
	}

	for _, detail := range news.Topics {

		topicID, err := uuid.Parse(detail.TopicId)
		if err != nil {
			return nil, errors.New("invalid News ID: " + detail.TopicId)
		}

		topicQuery := `
			INSERT INTO news_topic (news_id, topic_id, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW())`

		_, err = u.Conn.Exec(ctx, topicQuery, id, topicID)
		if err != nil {
			return nil, err
		}
	}
	//news.UpdatedAt = utils.ParseTime(updatedAt)
	return &updatedNews, nil
}

func (u *NewsRepository) DeleteNews(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE news
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
