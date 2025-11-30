## 项目数据库选型与设计
用了PostgreSQL ,因为业务需要如评分精确到小数和创建电影时要upsert,包括组合查询什么的。PostgreSQL方便点

**数据库表设计**


**电影表 (movies)**


标题加 UNIQUE 限制
电影标题在业务里必须唯一，所以用了UNIQUE。这样当有重复创建同名电影时，数据库会返回错误（PG 会抛 23505），应用拿到后转成 422 状态码，让客户端能明确知道是“标题重复”。
票房信息是通过外部 API 拉取的，但外部服务不一定每次都能成功。所以我把这些字段都设计成可空。
这样即使上游挂了，也不影响电影的创建流程，系统能继续正常工作。
上映日期用的是 DATE，因为后面可能要按年份筛选；预算和票房用了 BIGINT


```sql
CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,              -- 电影标题（唯一索引）
    release_date DATE NOT NULL,              -- 上映日期
    genre TEXT NOT NULL,                     -- 类型
    distributor TEXT,                        -- 发行商（可选）
    budget BIGINT,                           -- 预算（可选，单位：美元）
    mpa_rating TEXT,                         -- 分级（可选）
    
    -- 票房信息字段（来自外部 API）
    boxoffice_revenue_worldwide BIGINT,      -- 全球票房
    boxoffice_revenue_opening_weekend_usa BIGINT,  -- 美国首周末票房
    boxoffice_currency TEXT,                 -- 货币单位
    boxoffice_source TEXT,                   -- 数据来源
    boxoffice_last_updated TIMESTAMP,        -- 数据更新时间
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```


**Rating表**

复合主键 (movie_title, rater_id)，每个用户对每部电影只能有一条评分记录，当同一用户再次提交评分时，会更新现有记录而非插入新记录。
外键级联删除：当电影被删除时，相关的所有评分也会自动删除

评分值约束：通过 CHECK 约束确保评分值合法性：

范围：0.5 ≤ rating ≤ 5.0
步进：必须是 0.5 的倍数（即只能是 0.5, 1.0, 1.5, ..., 5.0）
约束表达式 MOD((rating * 2)::INTEGER, 1) = 0 的含义是：评分乘以 2 后必须是整数，这样就只允许半星步进。 这是Claude 建议的，我自己没想出来



```sql
CREATE TABLE ratings (
    movie_title TEXT NOT NULL,
    rater_id TEXT NOT NULL,
    rating NUMERIC(2,1) NOT NULL 
        CHECK (rating >= 0.5 AND rating <= 5.0 
               AND MOD((rating * 2)::INTEGER, 1) = 0), 
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (movie_title, rater_id),
    FOREIGN KEY (movie_title) REFERENCES movies(title) ON DELETE CASCADE
);
```

用基于文件的迁移方案：
- **001_create_movies.sql**：创建电影表
- **002_create_rating.sql**：创建评分表
迁移文件按顺序编号，在应用启动时自动执行。




## 后端服务选型与设计
- 用了Go , Gin。 数据库驱动：sqlx + pq 采用分层架构



##  优化方向
- **缓存机制**：引入 Redis 缓存热门电影信息和评分聚合结果
- **批量操作**：支持批量导入电影数据，使用事务和批量插入提升性能


















