package converter

import (
	"errors"

	"github.com/twin-te/user-service/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGRPCError(err error, log func(err error)) error {
	switch {
	case errors.Is(err, entity.ErrUserNotFound):
		return status.Error(codes.NotFound, "指定されたユーザーが見つかりませんでした")
	case errors.Is(err, entity.ErrAuthenticationAlreadyExists):
		return status.Error(codes.AlreadyExists, "認証情報が既に登録されています")
	default:
		log(err)
		return status.Error(codes.Internal, "サーバー内で問題が発生しました")
	}
}
