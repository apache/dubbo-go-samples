# 静态路由配置

这个目录用于展示如何在代码中直接配置 Dubbo-Go 的静态路由。
这里的 `static_config` 指的是：在 consumer 侧通过代码静态注入路由规则，而不是从配置中心动态下发。

[English](README.md) | 中文

## 前置准备

- Go 1.25+。

## 子示例

- `condition`：基于直连地址的服务级静态条件路由示例
- `tag`：基于直连地址的应用级静态标签路由示例

## 如何使用

进入对应子目录，按各自 README 运行：

- `condition/README_CN.md`
- `tag/README_CN.md`

`condition` 示例只使用直连 URL。
`tag` 示例只使用直连 URL。
两个示例都不需要配置中心。
每个子示例都需要先完成各自 README 中的前置准备和启动步骤。
