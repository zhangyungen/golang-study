package cache

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/maypok86/otter/v2"
	"github.com/maypok86/otter/v2/stats"
	"time"
	"zyj.com/golang-study/util/data"
	"zyj.com/golang-study/xorm/biz"
	"zyj.com/golang-study/xorm/model"
)

func ExampleOtter() {
	// 创建基本缓存实例，容量10000个条目
	cache := otter.Must(&otter.Options[string, string]{
		MaximumSize:     10_000,
		InitialCapacity: 100,
		Weigher: func(key string, value string) uint32 {
			return uint32(len(key) + len(value))
		},
		ExpiryCalculator:  otter.ExpiryAccessing[string, string](time.Hour), // 写入后1小时过期
		RefreshCalculator: otter.RefreshWriting[string, string](5 * time.Minute),
		StatsRecorder:     stats.NewCounter(),
	})

	// 添加缓存项
	cache.Set("user:1001", `{"name":"Alice","age":30}`)

	// 获取缓存项（检查是否存在）
	if value, ok := cache.GetIfPresent("user:1001"); ok {
		fmt.Printf("用户信息: %s\n", value) // 输出: 用户信息: {"name":"Alice","age":30}
	}

	// 删除缓存项
	if value, deleted := cache.Invalidate("user:1001"); deleted {
		fmt.Printf("已删除: %s\n", value)
	}

	// 带过期时间的缓存配置
	cacheWithExpiry := otter.Must(&otter.Options[int, []byte]{
		MaximumSize:      1_000,
		ExpiryCalculator: otter.ExpiryAccessing[int, []byte](5 * time.Minute), // 访问后5分钟过期
	})

	cacheWithExpiry.Set(1, []byte("重要数据"))

	// 应用自定义过期策略
	cache2 := otter.Must(&otter.Options[string, Session]{
		MaximumSize:      10_000,
		ExpiryCalculator: &SessionExpiry{},
	})
	cache2.Set("session:1", Session{IsPremium: true})

	ctx := context.Background()

	// 定义数据库加载器
	loader := otter.LoaderFunc[int64, *model.User](func(ctx context.Context, key int64) (*model.User, error) {
		// 模拟数据库查询
		user, err := biz.UserQueryBizIns.GetUserById(key)
		if err == sql.ErrNoRows {
			// 明确指示未找到，可缓存空值避免穿透
			return &model.User{}, otter.ErrNotFound
		}
		return user, err
	})

	cache3 := otter.Must(&otter.Options[int64, *model.User]{
		MaximumSize:      10_000,
		InitialCapacity:  100,
		ExpiryCalculator: otter.ExpiryAccessing[int64, *model.User](time.Hour), // 写入后1小时过期
		//BulkLoader:      otter.BulkLoaderFunc[int64, *model.User](loader.BulkLoad),
	})

	// 获取或加载数据（自动处理并发请求）
	user, err := cache3.Get(ctx, int64(10001), loader)
	if err != nil {
		if errors.Is(err, otter.ErrNotFound) {
			fmt.Println("用户不存在")
		} else {
			fmt.Printf("查询失败: %v\n", err)
		}
	} else {
		fmt.Printf("用户信息: %+v\n", user)
	}

	// 批量加载器定义
	bulkLoader := otter.BulkLoaderFunc[int64, *model.User](func(ctx context.Context, keys []int64) (map[int64]*model.User, error) {
		// 批量查询数据库，减少IO次数
		users, err := biz.UserQueryBizIns.ListUserByIds(keys)
		if err != nil {
			return nil, err
		}
		maps := data.MapKeyObjectPtr[model.User, int64](users,
			func(user model.User) int64 { return user.Id })
		return maps, nil
	})

	// 批量获取多个键
	userIDs := []int64{1001, 1002, 1003, 1004}
	results, err := cache3.BulkGet(ctx, userIDs, bulkLoader)
	if err != nil {
		fmt.Printf("批量查询失败: %v\n", err)
	} else {
		for id, user := range results {
			fmt.Printf("用户%d: %+v\n", id, user)
		}
	}

}

type SessionExpiry struct{}

type Session struct {
	IsPremium bool
	expiry    time.Time
}

func (e *SessionExpiry) ExpireAfterCreate(entry otter.Entry[string, Session]) time.Duration {
	if entry.Value.IsPremium {
		return 24 * time.Hour // 高级会话缓存24小时
	}
	return 2 * time.Hour // 普通会话缓存2小时
}

func (e *SessionExpiry) ExpireAfterUpdate(entry otter.Entry[string, Session], oldValue Session) time.Duration {
	return e.ExpireAfterCreate(entry) // 更新后沿用相同策略
}

func (e *SessionExpiry) ExpireAfterRead(entry otter.Entry[string, Session]) time.Duration {
	return e.ExpireAfterCreate(entry) // 更新后沿用相同策略
}

func PreloadHotData(ctx context.Context, cache *otter.Cache[int, model.User]) error {
	// 获取热点数据ID列表（从分析服务或历史数据）
	//hotKeys, err := analyticsService.GetTopNKeys(100)
	//if err != nil {
	//	return err
	//}

	//hotKeys:= []int64{1,2,3,4,5,6,7,8,9,10}
	//// 批量加载器
	//bulkLoader := otter.BulkLoaderFunc[int, model.User](func(cacheCtx context.Context, keys []int64) (map[int]model.User, error) {
	//	return biz.UserQueryBizIns.ListUserByIds(keys)
	//})
	//
	//// 批量预热
	//_, err = cache.BulkGet(cacheCtx, hotKeys, bulkLoader)
	//return err
	return nil
}

func GetProduct(ctx context.Context, cache *otter.Cache[int, model.User], id int64) (model.User, error) {
	// 1. 检查布隆过滤器（如果使用）
	//if !bloomFilter.Test(id) {
	//	return model.User{}, ErrNotFound
	//}
	//
	//// 2. 从缓存获取
	//product, err := cache.get(cacheCtx, id, func(cacheCtx context.Context, id int) (model.User, error) {
	//	// 3. 数据库查询
	//	p, err := db.GetProduct(id)
	//	if err == sql.ErrNoRows {
	//		// 4. 缓存空值（设置较短过期时间）
	//		return Product{}, otter.ErrNotFound
	//	}
	//	return p, err
	//})

	//return product, err
	return model.User{}, nil
}
