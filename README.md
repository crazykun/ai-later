# AI Tools Navigator

这是一个使用 Golang 和 Gin 框架开发的 AI 工具导航网站。

## 功能特点

- 响应式设计，支持多种设备访问
- 按分类浏览 AI 工具
- 搜索功能支持按名称和描述搜索
- 支持分类筛选

## 技术栈

- Backend: Golang + Gin Framework
- Frontend: Tailwind CSS
- Template Engine: Go HTML Templates

## 本地运行

```bash
mv config.demo.yaml config.yaml
go run main.go
```

## 项目结构

```html
.
├── config       # 配置文件
│   └── config.go
├── config.demo.yaml # 配置文件, 用于示例 可以重命名成 config.yaml 使用
├── data
│   └── ai.json # AI 工具数据
├── global      # 全局变量
│   └── global.go
├── go.mod
├── go.sum
├── handlers    # 处理器
│   └── handlers.go
├── main.go
├── middleware # 中间件
│   └── globalmiddleware.go
├── models    # 模型
│   └── site.go
├── README.md
├── static   # 静态文件
│   ├── css
│   │   └── style.css
│   ├── img
│   │   └── logo.png
│   └── js
│       └── main.js
└── templates # 模板文件
    ├── index.html
    └── layout.html
```
