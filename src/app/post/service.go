package post

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
	"test-edot/constants"
	"test-edot/src/dto"
	"test-edot/src/factory"
	"test-edot/src/models"
	"test-edot/src/repository"
	"time"
)

type Service interface {
	AddPostService(ctx context.Context, payload dto.AddPost) error
	GetPostService(ctx context.Context, payload dto.ParamGetPost) (any, error)
	GetPostByIdService(ctx context.Context, postId int) (any, error)
	UpdatePostService(ctx context.Context, postId int, payload dto.AddPost) error
	DeletePostService(ctx context.Context, postId int) error
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

func (s *service) DeletePostService(ctx context.Context, postId int) error {
	post, err := s.PostRepository.FindOneWithTag(ctx, "id = ?", postId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrorPostNotFound
		}

		return err
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

	if err := s.PostRepository.DeleteOne(tx, "id = ?", postId); err != nil {
		return err
	}

	for _, tag := range post.Tags {
		if err := s.PosTagRepository.DeleteOne(tx, "id_post = ? and id_tag = ?", postId, tag.ID); err != nil {
			return err
		}
	}

	if err := s.PostRepository.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *service) UpdatePostService(ctx context.Context, postId int, payload dto.AddPost) error {
	post, err := s.PostRepository.FindOneWithTag(ctx, "id = ?", postId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrorPostNotFound
		}

		return err
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

	// process
	if err := s.UpdatePost(tx, postId, payload, post); err != nil {
		return err
	}

	if err := s.UpdateTag(ctx, tx, postId, payload, post); err != nil {
		return err
	}

	if err := s.PostRepository.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateTag(ctx context.Context, tx *gorm.DB, postId int, payload dto.AddPost, post models.PostWithTag) error {
	mapUpdatedTag := map[string]bool{}

	mapTags := map[string]bool{}
	mapTagsId := map[string]int{}
	for _, tag := range post.Tags {
		mapTags[tag.Slug] = true
		mapTagsId[tag.Slug] = tag.ID
	}

	for _, updateTag := range payload.Tags {
		tagLower := strings.ReplaceAll(strings.ToLower(updateTag), " ", "_")
		if mapTags[tagLower] {
			delete(mapTags, tagLower)
			delete(mapTagsId, tagLower)
			continue
		}

		dTag, err := s.TagRepository.FindOne(ctx, "slug = ?", tagLower)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if dTag == (models.Tag{}) {
			data := models.Tag{
				Label: updateTag,
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

		mapUpdatedTag[tagLower] = true
	}

	// delete in active tag
	for _, i := range mapTagsId {
		if err := s.PosTagRepository.DeleteOne(tx, "id_post = ? and id_tag = ?", postId, i); err != nil {
			return err
		}
	}

	return nil
}
func (s *service) UpdatePost(tx *gorm.DB, postId int, payload dto.AddPost, post models.PostWithTag) error {
	updatedField := "content,updated_at"
	updatedData := models.Post{
		Content:   payload.Content,
		UpdatedAt: time.Now(),
	}

	slugNew := s.CreateSlug(payload.Title)
	if slugNew != post.Slug {
		updatedData.Slug = slugNew
		updatedData.Title = payload.Title
		updatedField = updatedField + ",slug,title"
	}

	if err := s.PostRepository.UpdateOne(tx, updatedData, updatedField, "id = ?", postId); err != nil {
		return err
	}

	return nil
}

func (s *service) GetPostByIdService(ctx context.Context, postId int) (any, error) {
	var (
		tags []string
	)

	post, err := s.PostRepository.FindOneWithTag(ctx, "id = ?", postId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constants.ErrorPostNotFound
		}

		return nil, err
	}

	for _, tag := range post.Tags {
		tags = append(tags, tag.Label)
	}

	return dto.AddPost{
		Title:   post.Title,
		Content: post.Content,
		Tags:    tags,
	}, err
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
			Id:      post.ID,
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
