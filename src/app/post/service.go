package post

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
	"test-asset-fendr/constants"
	"test-asset-fendr/src/dto"
	"test-asset-fendr/src/factory"
	"test-asset-fendr/src/models"
	"test-asset-fendr/src/repository"
	"time"
)

type Service interface {
	AddPostService(ctx context.Context, payload dto.AddPost) error
	GetPostService(ctx context.Context, payload dto.ParamGetPost) (any, error)
}

type service struct {
	Log              *zap.Logger
	PostRepository   repository.PostRepositoryInterface
	TagRepository    repository.TagRepositoryInterface
	PosTagRepository repository.PostTagRepositoryInterface
}

func NewService(f *factory.Factory) Service {
	return &service{
		Log:              f.Log,
		PostRepository:   f.PostRepository,
		TagRepository:    f.TagRepository,
		PosTagRepository: f.PostTagRepository,
	}
}

func (s *service) GetPostService(ctx context.Context, payload dto.ParamGetPost) (any, error) {
	limit := 20
	if payload.Limit != 0 {
		limit = payload.Limit
	}

	posts, err := s.PostRepository.GetPosts(ctx, payload.Offset, limit)
	if err != nil {
		return nil, err
	}

	return s.FormatterPosts(posts), nil
}

func (s *service) FormatterPosts(posts []models.PostWithTag) []dto.AddPost {
	var results []dto.AddPost
	for _, post := range posts {
		var tags []string
		for _, tag := range post.Tags {
			tags = append(tags, tag.Label)
		}

		result := dto.AddPost{
			Title:   post.Title,
			Content: post.Content,
			Tags:    tags,
		}

		results = append(results, result)
	}
	return results
}

func (s *service) AddPostService(ctx context.Context, payload dto.AddPost) error {
	slug := s.CreateSlug(payload.Title)

	// check title by slug today is already inserted or not
	data, err := s.PostRepository.FindOne(ctx, "slug = ?", slug)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if data != (models.Post{}) {
		return constants.ErrorPostAlreadyInserted
	}

	tx := s.PostRepository.Begin(ctx)
	defer func() {
		if r := recover(); r != nil {
			s.PostRepository.Rollback()
		}
	}()

	if err := s.PostRepository.Error(); err != nil {
		return err
	}

	postId, err := s.InsertNewPost(payload, slug)
	if err != nil {
		return err
	}

	if err := s.ProcessTag(ctx, tx, postId, payload.Tags); err != nil {
		s.PostRepository.Rollback()
		return err
	}

	if err := s.PostRepository.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *service) ProcessTag(ctx context.Context, tx *gorm.DB, postId int, tags []string) error {
	mapInsertedTag := map[string]bool{}

	for _, tag := range tags {
		tagLower := strings.ReplaceAll(strings.ToLower(tag), " ", "_")
		if !mapInsertedTag[tagLower] {
			dTag, err := s.TagRepository.FindOne(ctx, "slug = ?", tagLower)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if dTag == (models.Tag{}) {
				data := models.Tag{
					Label: tag,
					Slug:  tagLower,
				}

				if err := s.TagRepository.Create(tx, &data); err != nil {
					return err
				}

				dTag.ID = data.ID
			}

			if err := s.InsertPostTag(tx, postId, dTag.ID); err != nil {
				return err
			}

			mapInsertedTag[tagLower] = true
		}
	}

	return nil
}

func (s *service) InsertPostTag(tx *gorm.DB, postId, tagId int) error {
	data := models.PostTag{
		PostId: postId,
		TagId:  tagId,
	}
	if err := s.PosTagRepository.Create(tx, &data); err != nil {
		return err
	}

	return nil
}

func (s *service) InsertNewPost(payload dto.AddPost, slug string) (int, error) {
	data := models.Post{
		Title:     payload.Title,
		Content:   payload.Content,
		Slug:      slug,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.PostRepository.Create(&data); err != nil {
		return 0, err
	}

	s.Log.Info("finish insert post", zap.Int("post Id", data.ID))
	return data.ID, nil
}

func (s *service) CreateSlug(title string) string {
	dateNow := time.Now().Format("2006-01-02")
	titleLower := strings.ReplaceAll(strings.ToLower(title), " ", "_")
	return fmt.Sprintf("%s_%s", dateNow, titleLower)
}
