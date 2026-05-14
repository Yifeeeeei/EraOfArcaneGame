# 奥术纪元开发进度

## 当前范围

项目当前只支持基础包。非基础包卡牌已经从运行时代码和仓库数据中清掉。

## 技术架构

- 后端：Go `net/http` + `gorilla/websocket`
- 前端：静态 HTML/CSS/JS + Vue 3 CDN
- 卡牌定义：`server/cards/definitions_gen.go` 中的 Go card definitions
- 基础包快照：`data/supported_card_infos.json`
- 卡牌图片：`https://yifeeeeei.github.io/ArcaneImages/{output_path}`

## 当前状态

- 基础包 378 张卡是唯一 playable card pool。
- 后端不再运行时读取全量 JSON。
- 每张基础包卡都有 Go definition，暴露 `ID/Name/Kind/Element/Card`。
- 特殊效果按卡号拆到 `server/game/card_effects_*.go`。
- 前端已有双客户端操作回归脚本：`tools/frontend-card-operation-test.js`。

## 验证

```bash
cd server
go test ./...
```

前端卡牌操作回归报告保存在 `tmp/`，该目录不入库。
