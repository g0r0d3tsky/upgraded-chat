package cache

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	mock_cache "2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/cache/mocks"
	mock_repo "2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/cache/mocks"
	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/domain"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/mock/gomock"
)

func TestGetMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockRepository(ctrl)
	mockCache := mock_cache.NewMockRedis(ctrl)
	cacheService := NewCache(mockCache, mockRepo)

	testCases := []struct {
		name         string
		mockFunc     func()
		amount       int
		mockRepoResp []*domain.Message
		mockCacheErr error
		want         []*domain.Message
		wantErr      bool
	}{
		{
			name: "Cache hit, return messages from cache",
			mockFunc: func() {
				mockCache.EXPECT().LLen(gomock.Any(), gomock.Eq("messages")).Return(redis.NewIntResult(5, nil))
				mockCache.EXPECT().LRange(gomock.Any(), gomock.Eq("messages"), int64(2), int64(4)).Return(redis.NewStringSliceResult([]string{
					`{"id": "ab32e4f0-a283-4ee3-99b0-75f9ab82c478", "nickname": "User1", "content": "Message 1", "time": "2022-01-01T00:00:00Z"}`,
					`{"id": "99efa34c-dc05-4b24-aa1c-c70db9cae892", "nickname": "User2", "content": "Message 2", "time": "2022-01-02T00:00:00Z"}`,
					`{"id": "943f9aae-797e-422f-8a8b-9012f9ad7565", "nickname": "User3", "content": "Message 3", "time": "2022-01-03T00:00:00Z"}`,
				}, nil))
			},
			amount: 3,
			want: []*domain.Message{
				{ID: uuid.MustParse("ab32e4f0-a283-4ee3-99b0-75f9ab82c478"), UserNickname: "User1", Content: "Message 1", Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				{ID: uuid.MustParse("99efa34c-dc05-4b24-aa1c-c70db9cae892"), UserNickname: "User2", Content: "Message 2", Time: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
				{ID: uuid.MustParse("943f9aae-797e-422f-8a8b-9012f9ad7565"), UserNickname: "User3", Content: "Message 3", Time: time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC)},
			},
			wantErr: false,
		},
		{
			name: "Cache miss, return messages from repository and update cache",
			mockFunc: func() {
				mockCache.EXPECT().LLen(gomock.Any(), gomock.Eq("messages")).Return(redis.NewIntResult(0, nil))
				mockRepo.EXPECT().GetAmountMessage(gomock.Any(), 3).Return([]*domain.Message{
					{ID: uuid.MustParse("b1da7a09-9bc3-4326-8d29-cb237631d298"), UserNickname: "User6", Content: "Message 6", Time: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)},
					{ID: uuid.MustParse("87b56064-77ad-44e4-8abb-4a82224092cd"), UserNickname: "User7", Content: "Message 7", Time: time.Date(2022, 1, 7, 0, 0, 0, 0, time.UTC)},
					{ID: uuid.MustParse("0f0d7aa7-de40-45e5-8880-062900397181"), UserNickname: "User8", Content: "Message 8", Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
				}, nil)
				mockCache.EXPECT().RPush(gomock.Any(), gomock.Eq("messages"), gomock.Any()).Return(redis.NewIntCmd(context.Background(), nil)).AnyTimes()
			},
			amount: 3,
			want: []*domain.Message{
				{ID: uuid.MustParse("b1da7a09-9bc3-4326-8d29-cb237631d298"), UserNickname: "User6", Content: "Message 6", Time: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)},
				{ID: uuid.MustParse("87b56064-77ad-44e4-8abb-4a82224092cd"), UserNickname: "User7", Content: "Message 7", Time: time.Date(2022, 1, 7, 0, 0, 0, 0, time.UTC)},
				{ID: uuid.MustParse("0f0d7aa7-de40-45e5-8880-062900397181"), UserNickname: "User8", Content: "Message 8", Time: time.Date(2022, 1, 8, 0, 0, 0, 0, time.UTC)},
			},
			wantErr: false,
		},
		{
			name: "Error getting cache length",
			mockFunc: func() {
				mockCache.EXPECT().LLen(gomock.Any(), gomock.Eq("messages")).Return(redis.NewIntResult(0, errors.New("cache error")))
			},
			amount:  3,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error getting messages from cache",
			mockFunc: func() {
				mockCache.EXPECT().LLen(gomock.Any(), gomock.Eq("messages")).Return(redis.NewIntResult(5, nil))
				mockCache.EXPECT().LRange(gomock.Any(), "messages", int64(2), int64(4)).Return(redis.NewStringSliceResult(nil, errors.New("cache error")))
			},
			amount:  3,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error getting messages from repository",
			mockFunc: func() {
				mockCache.EXPECT().LLen(gomock.Any(), gomock.Eq("messages")).Return(redis.NewIntResult(0, nil))
				mockRepo.EXPECT().GetAmountMessage(gomock.Any(), 3).Return(nil, errors.New("repository error"))
			},
			amount:  3,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()
			got, err := cacheService.GetMessages(context.Background(), tc.amount)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetMessages() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if len(got) != len(tc.want) {
				t.Errorf("GetMessages() got %d messages, want %d", len(got), len(tc.want))
				return
			}
			for i, msg := range got {
				if !reflect.DeepEqual(msg, tc.want[i]) {
					t.Errorf("GetMessages() got[%d] = %v, want[%d] = %v", i, msg, i, tc.want[i])
				}
			}
		})
	}
}
