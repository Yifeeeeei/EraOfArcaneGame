# 奥术纪元开发进度

## 当前范围

项目当前只支持基础包。非基础包卡牌已经从运行时代码和仓库数据中清掉。

## 技术架构

- 后端：Go `net/http` + `gorilla/websocket`
- 前端：静态 HTML/CSS/JS + Vue 3 CDN
- 卡牌定义：`server/cards/definitions_gen.go` 中的 Go card definitions
- 卡牌类别：`server/cards/interfaces.go` + `server/cards/category_markers_gen.go` 的 Go interfaces/markers
- 卡牌行为：`server/game/card_behavior.go` 的接口 + `server/game/card_<number>_<name>.go` 的一牌一文件 struct
- 基础包快照：`data/supported_card_infos.json`
- 卡牌图片：`https://yifeeeeei.github.io/ArcaneImages/{output_path}`

## 当前状态

- 基础包 378 张卡是唯一 playable card pool。
- 后端不再运行时读取全量 JSON。
- 每张基础包卡都有 Go definition，暴露 `ID/Name/Kind/Element/Card`。
- 基础包卡定义实现了人物、伙伴、技能、道具等类别 interface。
- 自定义效果不再由文本 parser 推断；基础包效果由具体卡牌 struct 实现 `OnEnter`、`OnDeath`、`OnUltimate` 等接口。
- 行为 catalog 只注册 lazy factory，不会启动时实例化全部 behavior。
- 非基础包效果代码已从运行时 catalog 清理掉。
- 前端已有双客户端操作回归脚本：`tools/frontend-card-operation-test.js`。

## 验证

```bash
cd server
go test ./...
```

前端卡牌操作回归报告保存在 `tmp/`，该目录不入库。
