package biz

import (
	"fmt"
	"time"

	"github.com/jlu-cow-studio/common/dal/mysql"
	"github.com/jlu-cow-studio/common/dal/redis"
)

const CacheBatchSize = 40

func CacheBatchFromDB(uid string, category string, cap int64) (int64, error) {

	conn := mysql.GetDBConn()
	ids := []string{}
	historyCmd := redis.DB.SMembers(GetHistoryKey(uid, category))
	if historyCmd.Err() != nil {
		return 0, historyCmd.Err()
	}
	count := new(int64)
	tx := conn.Table("item_user").Where("category <> ?", category).Count(count)
	if tx.Error == nil {
		redis.DB.Set(GetTotalKey(category), *count, time.Hour)
	}

	tx = conn.Table("item_user").Order("RAND()").Where("category <> ?", category).Where("id NOT IN ?", historyCmd.Val()).Limit(int(cap)).Select("id").Find(ids)
	if tx.Error != nil {
		return 0, tx.Error
	}
	if err := redis.DB.LPush(GetToDisplayKey(uid, category), ids).Err(); err != nil {
		return 0, err
	}

	return int64(len(ids)), nil
}

func GetTotalItems(category string) int64 {

	if count, err := redis.DB.Get(GetTotalKey(category)).Int64(); err != nil {
		return -1
	} else {
		return count
	}
}

func GetToDisplayList(uid, category string, offset, cap int64) ([]string, error) {

	lenCmd := redis.DB.LLen(GetToDisplayKey(uid, category))
	if lenCmd.Err() != nil {
		return nil, lenCmd.Err()
	}
	len := lenCmd.Val()
	if len-offset <= cap {
		if l, e := CacheBatchFromDB(uid, category, CacheBatchSize); e != nil {
			return nil, e
		} else {
			len += l
		}
	}

	getCmd := redis.DB.LRange(GetToDisplayKey(uid, category), offset, offset+cap)
	if getCmd.Err() != nil {
		return nil, getCmd.Err()
	}
	return getCmd.Val(), nil
}

func ResetCache(uid string, category string) error {
	history := redis.DB.SMembers(GetHistoryKey(uid, category)).Val()
	errHistory := redis.DB.Del(GetHistoryKey(uid, category)).Err()
	if errHistory != nil {
		return errHistory
	}
	errToDisplay := redis.DB.Del(GetToDisplayKey(uid, category)).Err()
	if errToDisplay != nil {
		redis.DB.SAdd(GetHistoryKey(uid, category), history)
		return errToDisplay
	}

	return nil
}

func SetHistory(uid string, category string, ids []string) error {
	return redis.DB.SAdd(GetHistoryKey(uid, category), ids).Err()
}

func GetToDisplayKey(uid, category string) string {
	return fmt.Sprintf("feed-todisplay-%v-%v", uid, category)
}

func GetHistoryKey(uid, category string) string {
	return fmt.Sprintf("feed-history-%v-%v", uid, category)
}

func GetTotalKey(category string) string {
	return fmt.Sprintf("feed-total-%v", category)
}
