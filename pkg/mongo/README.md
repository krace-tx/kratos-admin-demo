# MongoDB 基本用法指南

## 1. 安装与启动

### 安装
- Ubuntu: `sudo apt-get install mongodb`
- Mac: `brew install mongodb`
- Windows: 从官网下载安装包

### 启动服务
```bash
mongod --dbpath /path/to/data/directory
```

### 连接客户端
```bash
mongo
```

## 2. 数据库操作

```javascript
// 显示所有数据库
show dbs

// 使用/创建数据库
use mydb

// 删除当前数据库
db.dropDatabase()
```

## 3. 集合(表)操作

```javascript
// 创建集合
db.createCollection("users")

// 显示所有集合
show collections

// 删除集合
db.users.drop()
```

## 4. CRUD 操作

### 插入文档
```javascript
// 插入单个文档
db.users.insertOne({
  name: "John Doe",
  age: 30,
  email: "john@example.com"
})

// 插入多个文档
db.users.insertMany([
  {name: "Alice", age: 25},
  {name: "Bob", age: 35}
])
```

### 查询文档
```javascript
// 查询所有文档
db.users.find()

// 带条件查询
db.users.find({age: {$gt: 25}})

// 查询第一个匹配的文档
db.users.findOne({name: "John Doe"})

// 限制返回字段
db.users.find({}, {name: 1, email: 1})

// 排序
db.users.find().sort({age: 1}) // 1升序，-1降序

// 分页
db.users.find().skip(10).limit(5)
```

### 更新文档
```javascript
// 更新单个文档
db.users.updateOne(
  {name: "John Doe"},
  {$set: {age: 31}}
)

// 更新多个文档
db.users.updateMany(
  {age: {$lt: 30}},
  {$set: {status: "young"}}
)

// 替换文档
db.users.replaceOne(
  {name: "John Doe"},
  {name: "John Doe", age: 32, email: "john.doe@example.com"}
)
```

### 删除文档
```javascript
// 删除单个文档
db.users.deleteOne({name: "John Doe"})

// 删除多个文档
db.users.deleteMany({age: {$gt: 30}})
```

## 5. 索引操作

```javascript
// 创建索引
db.users.createIndex({email: 1}, {unique: true})

// 查看索引
db.users.getIndexes()

// 删除索引
db.users.dropIndex("email_1")
```

## 6. 聚合操作

```javascript
// 简单聚合
db.users.aggregate([
  {$match: {age: {$gt: 25}}},
  {$group: {_id: "$status", total: {$sum: 1}}}
])

// 常用聚合阶段:
// $match: 过滤文档
// $group: 分组
// $sort: 排序
// $project: 投影/重塑文档
// $limit: 限制数量
// $skip: 跳过数量
```

## 7. 其他实用操作

```javascript
// 统计文档数量
db.users.countDocuments({age: {$gt: 25}})

// 去重查询
db.users.distinct("status")

// 批量操作
var bulk = db.users.initializeUnorderedBulkOp();
bulk.insert({name: "User1"});
bulk.insert({name: "User2"});
bulk.execute();
```

## 8. 数据备份与恢复

```bash
# 备份
mongodump --db mydb --out /backup/

# 恢复
mongorestore --db mydb /backup/mydb/
```

# MongoDB 进阶用法指南

## 1. 高级查询技巧

### 复杂条件查询
```javascript
// 多条件组合
db.users.find({
  $and: [
    {age: {$gt: 25}},
    {status: "active"},
    {$or: [
      {email: /@company.com$/},
      {department: "IT"}
    ]}
  ]
})

// 数组查询
db.products.find({
  tags: {$all: ["electronics", "new"]}, // 包含所有指定元素
  ratings: {$elemMatch: {score: {$gt: 4}}} // 数组元素匹配
})

// 正则表达式查询
db.users.find({
  name: {$regex: /^J/, $options: "i"} // 以J开头，不区分大小写
})
```

## 2. 聚合框架高级用法

### 复杂聚合管道
```javascript
db.orders.aggregate([
  // 阶段1: 匹配条件
  {$match: {
    date: {$gte: ISODate("2023-01-01")},
    status: "completed"
  }},
  
  // 阶段2: 按客户分组
  {$group: {
    _id: "$customerId",
    totalSpent: {$sum: "$amount"},
    avgOrder: {$avg: "$amount"},
    orderCount: {$sum: 1},
    firstOrder: {$min: "$date"},
    lastOrder: {$max: "$date"}
  }},
  
  // 阶段3: 关联客户信息
  {$lookup: {
    from: "customers",
    localField: "_id",
    foreignField: "_id",
    as: "customer"
  }},
  
  // 阶段4: 解构数组
  {$unwind: "$customer"},
  
  // 阶段5: 投影重塑
  {$project: {
    customerName: "$customer.name",
    email: "$customer.email",
    totalSpent: 1,
    avgOrder: 1,
    orderCount: 1,
    loyaltyLevel: {
      $switch: {
        branches: [
          {case: {$gte: ["$totalSpent", 1000]}, then: "Gold"},
          {case: {$gte: ["$totalSpent", 500]}, then: "Silver"}
        ],
        default: "Bronze"
      }
    }
  }},
  
  // 阶段6: 排序
  {$sort: {totalSpent: -1}},
  
  // 阶段7: 分页
  {$skip: 10},
  {$limit: 5}
])
```

### 聚合运算符进阶
```javascript
// 日期处理
db.events.aggregate([
  {$project: {
    year: {$year: "$date"},
    month: {$month: "$date"},
    dayOfWeek: {$dayOfWeek: "$date"},
    durationMinutes: {$divide: [{$subtract: ["$endTime", "$startTime"]}, 60000]}
  }}
])

// 条件逻辑
db.products.aggregate([
  {$project: {
    priceCategory: {
      $cond: {
        if: {$gte: ["$price", 100]},
        then: "premium",
        else: "standard"
      }
    },
    discountPrice: {
      $multiply: [
        "$price",
        {$subtract: [1, {$ifNull: ["$discount", 0]}]}
      ]
    }
  }}
])
```

## 3. 索引优化

### 高级索引策略
```javascript
// 复合索引
db.users.createIndex({lastName: 1, firstName: 1})

// 多键索引(数组字段)
db.products.createIndex({tags: 1})

// 文本索引(全文搜索)
db.articles.createIndex({content: "text"})
db.articles.find({$text: {$search: "mongodb tutorial"}})

// 地理空间索引
db.places.createIndex({location: "2dsphere"})
db.places.find({
  location: {
    $near: {
      $geometry: {
        type: "Point",
        coordinates: [longitude, latitude]
      },
      $maxDistance: 1000 // 1公里内
    }
  }
})

// TTL索引(自动过期)
db.sessions.createIndex({lastAccess: 1}, {expireAfterSeconds: 3600})

// 部分索引(条件索引)
db.users.createIndex(
  {email: 1},
  {partialFilterExpression: {email: {$exists: true}}}
)

// 隐藏索引(测试索引效果)
db.users.createIndex({age: 1}, {hidden: true})
db.unhideIndex("users", "age_1")
```

## 4. 事务处理

```javascript
// 多文档事务
const session = db.getMongo().startSession();
try {
  session.startTransaction({
    readConcern: {level: "snapshot"},
    writeConcern: {w: "majority"}
  });
  
  const users = session.getDatabase("mydb").users;
  const accounts = session.getDatabase("mydb").accounts;
  
  users.deleteOne({_id: "user1"}, {session});
  accounts.deleteMany({userId: "user1"}, {session});
  
  session.commitTransaction();
} catch (error) {
  session.abortTransaction();
  throw error;
} finally {
  session.endSession();
}
```

## 5. 性能优化技巧

### 查询分析
```javascript
// 启用查询分析
db.setProfilingLevel(2) // 0=关闭, 1=慢查询, 2=全部

// 查看分析结果
db.system.profile.find().sort({ts: -1}).limit(10)

// 解释查询计划
db.users.find({age: {$gt: 25}}).explain("executionStats")
```

### 性能优化实践
1. **使用投影减少返回数据量**：
   ```javascript
   db.users.find({status: "active"}, {name: 1, email: 1})
   ```

2. **批量操作代替循环**：
   ```javascript
   var bulk = db.users.initializeOrderedBulkOp();
   for (let i = 0; i < 1000; i++) {
     bulk.insert({user: "user"+i});
   }
   bulk.execute();
   ```

3. **合理使用游标**：
   ```javascript
   const cursor = db.users.find().batchSize(100);
   while (cursor.hasNext()) {
     processDocument(cursor.next());
   }
   ```

4. **读写分离**：
   ```javascript
   // 从secondary节点读取
   db.getMongo().setReadPref("secondaryPreferred")
   ```

## 6. 数据建模高级技巧

### 关系建模
```javascript
// 引用式(适合一对多)
{
  _id: "order123",
  customerId: "cust789",
  items: ["item456", "item789"]
}

// 嵌入式(适合一对一或少量一对多)
{
  _id: "user123",
  name: "Alice",
  addresses: [
    {type: "home", street: "123 Main St"},
    {type: "work", street: "456 Office Ave"}
  ]
}
```

### 分桶模式(时间序列数据)
```javascript
// 每小时一个文档，存储该小时的所有事件
{
  _id: "2023-01-01T13:00:00Z",
  events: [
    {time: "2023-01-01T13:05:00Z", type: "login"},
    {time: "2023-01-01T13:15:00Z", type: "purchase"}
  ],
  count: 2,
  firstEvent: ISODate("2023-01-01T13:05:00Z"),
  lastEvent: ISODate("2023-01-01T13:15:00Z")
}
```

### 多态模式
```javascript
// 不同类型的文档存储在同一个集合
{
  _id: "doc1",
  type: "book",
  title: "MongoDB Guide",
  author: "John Doe",
  pages: 300
},
{
  _id: "doc2",
  type: "movie",
  title: "MongoDB: The Movie",
  director: "Jane Smith",
  duration: 120
}
```

## 7. 复制集与分片集群

### 复制集管理
```bash
# 初始化复制集
rs.initiate({
  _id: "replSet",
  members: [
    {_id: 0, host: "mongo1:27017", priority: 2},
    {_id: 1, host: "mongo2:27017", priority: 1},
    {_id: 2, host: "mongo3:27017", priority: 1, arbiterOnly: true}
  ]
})

# 查看复制集状态
rs.status()
```

### 分片集群配置
```bash
# 添加分片
sh.addShard("shard1/mongo1:27017,mongo2:27017")

# 启用分片数据库
sh.enableSharding("mydb")

# 分片集合
sh.shardCollection("mydb.users", {userId: "hashed"})

# 查看分片状态
sh.status()
```

## 8. Change Streams (变更流)

```javascript
// 监听集合变更
const changeStream = db.orders.watch([
  {$match: {"operationType": "insert"}},
  {$project: {"fullDocument": 1}}
]);

changeStream.on("change", function(change) {
  console.log("New order:", change.fullDocument);
});

// 监听数据库变更
const dbStream = db.watch();
// 监听整个集群变更
const clusterStream = db.getMongo().watch();
```

## 9. MongoDB Atlas 云服务特性

1. **Atlas Search** (基于Lucene的全文搜索):
   ```javascript
   db.articles.aggregate([
     {$search: {
       index: "default",
       text: {
         query: "mongodb tutorial",
         path: {wildcard: "*"}
       }
     }}
   ])
   ```

2. **Atlas Data Lake** (跨数据源查询):
   ```javascript
   db.getSiblingDB("external").sales.find({
     date: {$gt: ISODate("2023-01-01")}
   })
   ```

3. **Atlas Triggers** (事件驱动):
   ```javascript
   // 配置数据库触发器响应数据变更
   ```

## 10. 安全最佳实践

```javascript
// 创建角色
db.createRole({
  role: "analyst",
  privileges: [
    {
      resource: {db: "reports", collection: ""},
      actions: ["find", "aggregate"]
    }
  ],
  roles: []
})

// 创建用户
db.createUser({
  user: "reportUser",
  pwd: "securePassword",
  roles: ["analyst"]
})

// 启用加密
// 在配置文件中设置:
// security:
//   authorization: enabled
//   keyFile: /path/to/keyfile
//   enableEncryption: true
//   encryptionKeyFile: /path/to/encryption/key
```