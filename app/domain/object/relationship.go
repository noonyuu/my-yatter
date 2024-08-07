package object

import (
)

// 本当はPasswordHashがハッシュされたパスワードであることを型で保証したい。
// ハッシュ化されたパスワード用の型を用意してstringと区別して管理すると良い。
// 今回は簡単のためstringで管理している。

type Relationship struct {
	FollowerID int64 `db:"follower_id"`
	FolloweeID int64 `db:"followee_id"`
}

func Follow(username string) (*Account, error) {
	return &Account{
		Username: username,
	}, nil
}