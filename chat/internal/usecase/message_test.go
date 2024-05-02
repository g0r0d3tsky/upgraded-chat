package usecase

import (
	"context"
	"errors"
	"testing"

	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/domain"
	mock_cache "2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/usecase/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestPush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBroker := mock_cache.NewMockBroker(ctrl)
	mockCache := mock_cache.NewMockCache(ctrl)
	messageService := &MessageService{
		broker: mockBroker,
		cache:  mockCache,
	}
	message := &domain.Message{
		ID:           uuid.New(),
		UserNickname: "User1",
		Content:      "Test message",
	}
	topic := "test-topic"
	testCases := []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Push broker successfully",
			mockFunc: func() {
				mockBroker.EXPECT().Push(topic, message).Return(nil)
				mockCache.EXPECT().SetMessages(gomock.Any(), []*domain.Message{message}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Push broker fail",
			mockFunc: func() {
				mockBroker.EXPECT().Push(topic, message).Return(errors.New("push fail"))
			},
			wantErr: true,
		},
		{
			name: "Set messages error",
			mockFunc: func() {
				mockBroker.EXPECT().Push(topic, message).Return(nil)
				mockCache.EXPECT().SetMessages(gomock.Any(), []*domain.Message{message}).Return(errors.New("setting messages"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()
			err := messageService.Push(context.Background(), topic, message)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetMessages() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}
