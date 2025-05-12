package cache

import (
	"context"
	"strconv"
)

func (Cache) keyUsers() string      { return "users" }
func (Cache) keyUsersInit() string  { return "users_init" }
func (Cache) keyUsersAdmin() string { return "users_admin" }

func (Cache) keyUserPoints(userID int64) string {
	return "user_points:" + strconv.FormatInt(userID, 10)
}

func (Cache) keyUserRefPoints(userID int64) string {
	return "ref_points:" + strconv.FormatInt(userID, 10)
}
func (Cache) keyUserRefCount(userID int64) string {
	return "ref_count:" + strconv.FormatInt(userID, 10)
}

func (c Cache) CheckUser(ctx context.Context, userID int64) (exist bool, err error) {
	return c.c.SIsMember(ctx, c.keyUsers(), userID).Result()
}

func (c Cache) AddUser(ctx context.Context, userID int64) error {
	return c.c.SAdd(ctx, c.keyUsers(), userID).Err()
}

func (c Cache) CheckUserInit(ctx context.Context, userID int64) (exist bool, err error) {
	return c.c.SIsMember(ctx, c.keyUsersInit(), userID).Result()
}

func (c Cache) AddUserInit(ctx context.Context, userID int64) error {
	return c.c.SAdd(ctx, c.keyUsersInit(), userID).Err()
}

func (c Cache) CheckUserAdmin(ctx context.Context, userID int64) (exist bool, err error) {
	return c.c.SIsMember(ctx, c.keyUsersAdmin(), userID).Result()
}

func (c Cache) AddUserAdmin(ctx context.Context, userID int64) error {
	return c.c.SAdd(ctx, c.keyUsersAdmin(), userID).Err()
}

func (c Cache) GetUserPoints(ctx context.Context, userID int64) (points int64, err error) {
	val, err := c.c.Get(ctx, c.keyUserPoints(userID)).Result()
	if err != nil {
		if isNil(err) {
			return 0, nil
		}
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}

func (c Cache) IncUserPoints(ctx context.Context, userID, incPoints int64) error {
	return c.c.IncrBy(ctx, c.keyUserPoints(userID), incPoints).Err()
}

func (c Cache) DecUserPoints(ctx context.Context, userID, decPoints int64) (newBalance int64, err error) {
	return c.c.DecrBy(ctx, c.keyUserPoints(userID), decPoints).Result()
}

func (c Cache) SetUserPoints(ctx context.Context, userID, points int64) error {
	return c.c.Set(ctx, c.keyUserPoints(userID), points, 0).Err()
}

func (c Cache) GetUserRefPoints(ctx context.Context, userID int64) (int64, error) {
	val, err := c.c.Get(ctx, c.keyUserRefPoints(userID)).Result()
	if err != nil {
		if isNil(err) {
			return 0, nil
		}
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}

func (c Cache) IncUserRefPoints(ctx context.Context, userID, incPoints int64) error {
	return c.c.IncrBy(ctx, c.keyUserRefPoints(userID), incPoints).Err()
}

func (c Cache) SetUserRefPoints(ctx context.Context, userID, points int64) error {
	return c.c.Set(ctx, c.keyUserRefPoints(userID), points, 0).Err()
}

func (c Cache) GetUserRefCount(ctx context.Context, userID int64) (int64, error) {
	val, err := c.c.Get(ctx, c.keyUserRefCount(userID)).Result()
	if err != nil {
		if isNil(err) {
			return 0, nil
		}
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}

func (c Cache) IncUserRefCount(ctx context.Context, userID int64) error {
	return c.c.Incr(ctx, c.keyUserRefCount(userID)).Err()
}

func (c Cache) SetUserRefCount(ctx context.Context, userID, count int64) error {
	return c.c.Set(ctx, c.keyUserRefCount(userID), count, 0).Err()
}
