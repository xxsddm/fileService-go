# 文件服务微服务 (Go + Gin 实现)

这是一个使用 Go 和 Gin 框架实现的文件管理微服务，功能与原有的 Java 版本相同。

## 功能特性

1. 文件上传（支持批量上传）
2. 文件下载
3. 文件列表分页查询
4. 文件名过滤查询
5. 文件删除
6. 自动清理过期文件（7天）
7. 前端页面管理文件

## 项目结构
file-service/
├── config.yaml          # 配置文件
├── go.mod               # Go 模块定义
├── main.go              # 主程序入口
├── README.md            # 说明文档
├── controllers/         # 控制器层
├── models/              # 数据模型
├── dto/                 # 数据传输对象
├── service/             # 业务逻辑层
├── database/            # 数据库连接
├── config/              # 配置管理
├── utils/               # 工具类
└── static/              # 静态文件（前端页面）


## 配置说明

在 `config.yaml` 文件中可以配置：

- 服务器端口
- 数据库连接信息
- 文件上传路径、大小限制和允许的文件类型

## 数据库表结构

```sql
CREATE TABLE `fileInfo` (
  `id` bigint(20) NOT NULL,
  `fileName` varchar(255) DEFAULT NULL,
  `filePath` varchar(255) DEFAULT NULL,
  `fileSize` int(11) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `uploadDate` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## API 接口

1. `POST /upload/` - 上传文件
2. `GET /download/{fileName}` - 下载文件
3. `GET /files/` - 获取文件列表（支持分页和过滤）

## 前端页面

访问 `http://localhost:8080` 可以使用前端页面进行文件管理。

## 部署说明

1. 确保已安装 Go 1.21+
2. 创建 MySQL 数据库 `file_service`
3. 根据需要修改 `config.yaml` 配置文件
4. 运行以下命令启动服务：

```bash
go mod tidy
go run main.go
```

服务默认运行在 8080 端口。