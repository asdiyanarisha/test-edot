package post

import (
	"context"
	"github.com/stretchr/testify/assert"
	"test-asset-fendr/constants"
	"test-asset-fendr/src/dto"
	"test-asset-fendr/src/models"
	repoMocks "test-asset-fendr/src/repository/mocks"
	"testing"
)

type TestGetPostUser struct {
	Name     string
	PostId   int
	Error    error
	Expected models.Post
}

func TestGetPostByIdService(t *testing.T) {
	ctx := context.Background()

	mockPostRepository := new(repoMocks.PostRepositoryInterface)
	MockPostData(mockPostRepository)

	s := service{PostRepository: mockPostRepository}

	testCases := []TestGetPostUser{
		{
			Name:     "test success",
			PostId:   1,
			Error:    nil,
			Expected: models.Post{Title: "test response", Content: "test content"},
		}, {
			Name:   "test failed",
			PostId: 2,
			Error:  constants.ErrorPostNotFound,
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			res, err := s.GetPostByIdService(ctx, test.PostId)
			if test.Error == nil {
				result := res.(dto.AddPost)
				assert.Equal(t, test.Expected.Title, result.Title, "test title")
				assert.Equal(t, test.Expected.Content, result.Content, "test content")
			} else {
				assert.Equal(t, test.Error, err, "test error")
			}
		})
	}

}

func MockPostData(repo *repoMocks.PostRepositoryInterface) {
	repo.On("FindOneWithTag", context.Background(), "id = ?", 1).Return(models.PostWithTag{
		ID:      1,
		Title:   "test response",
		Content: "test content",
		Slug:    "tes-slug",
		Tags: []models.Tag{
			{ID: 1, Label: "lorem", Slug: "lorem"},
		},
	}, nil)

	repo.On("FindOneWithTag", context.Background(), "id = ?", 2).Return(models.PostWithTag{}, constants.ErrorPostNotFound)
}
