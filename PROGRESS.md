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

- 基础包 393 张卡是唯一 playable card pool。
- 后端不再运行时读取全量 JSON。
- 每张基础包卡都有 Go definition，暴露 `ID/Name/Kind/Element/Card`。
- 基础包卡定义实现了人物、伙伴、技能、道具等类别 interface。
- 自定义效果不再由文本 parser 推断；基础包效果由具体卡牌 struct 实现 `OnEnter`、`OnDeath`、`OnUltimate` 等接口。
- 行为 catalog 只注册 lazy factory，不会启动时实例化全部 behavior。
- 非基础包效果代码已从运行时 catalog 清理掉。
- 前端已有双客户端操作回归脚本：`tools/frontend-card-operation-test.js`。
- `docs/project/` 是本地 agent 协作文档区，不再提交到 GitHub；需要公开给 review 的状态应同步到根目录文档或 `docs/rules/`。

## 最新收尾记录（2026-06-25）

- PR #58 / Issue #59：机器测试清单中的 40 张基础包卡已完成后端检查，重点覆盖 `PendingAction` 目标选择、模式选择、可选诱发和连续选择窗口。
- `3021006 洞察之眼` 已修复：施放时打开敌方盖放卡选择窗口，不再默认摧毁第一张敌方装备/盖牌。
- 新增/补强了 `预见`、`万灵药`、`黑市商贩`、`新生卷轴`、`血魔爆`、`元素附魔`、`专精法师`、`伦德萨尔`、`食腐者` 等牌的后端语义测试。
- 人工手测清单仍以本地 `docs/project/base-set-test-marking-sheet.md` 为准；不要从 GitHub 公开文档反向覆盖本地标注表。

## 验证

```bash
cd server
go test ./...
```

前端卡牌操作回归报告保存在 `tmp/`，该目录不入库。
