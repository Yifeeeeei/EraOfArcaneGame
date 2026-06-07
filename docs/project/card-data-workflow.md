# 卡牌数据更新流程

本文档说明当前仓库的卡牌信息从哪里读取，以及重建卡牌文本时推荐怎样更新。

## 1. 当前运行时卡牌来源

当前游戏运行时不从 JSON 读取卡牌。

后端启动时通过 Go 代码加载卡牌定义：

- `server/cards/definitions_gen.go`：当前基础包卡牌的编译后定义，是运行时卡牌数据来源。
- `server/cards/loader.go`：调用编译定义，构建 `CardDB`、`BaseCardDB` 和 `PlayableCardDB`。
- `cards.PlayableCardDB`：API、牌组校验和游戏对局实际使用的卡池。

前端 `/api/cards` 读取的是后端 `cards.PlayableCardDB` 暴露出来的数据，包括 `description`、费用、负载、威力、攻击等字段。

`data/supported_card_infos.json` 是当前基础包卡牌数据快照，也是生成编译定义的输入。后端启动时仍不直接读取 JSON；更新 JSON 后必须重新生成 Go 定义。

新版卡牌数据可以额外填写以下结构化字段：

- `effect_categories`：效果类别数组，例如 `["主动"]`、`["条件"]`、`["光环"]`、`["入场", "条件"]`、`["反制"]`、`["响应"]`。
- `effect_optionality`：效果是否可选数组，例如 `["强制"]`、`["可选"]`，一张牌同时有强制和可选效果时可写 `["强制", "可选"]`。

这些字段用于读牌、审计和前端展示，不用于从文本自动推断运行时效果。真正的卡牌行为仍必须写在 Go 行为代码里。

## 2. 当前单卡行为来源

卡牌效果的运行时行为来自 Go 代码，而不是中文描述文本。

主要位置：

- `server/game/card_<卡号>_<名称>.go`
- `server/game/card_behavior.go`
- `server/game/card_effects_catalog.go`
- `server/game/effect_system.go`

后续重建时仍应坚持：卡牌描述给玩家读，卡牌行为用结构化 Go 代码实现。不要通过解析描述文本推断效果。

## 3. 建议的卡牌文本重建路线

为了方便你批量改文本，又避免一上来改运行时代码，建议分三层推进。

### 第一步：建立人工审计表

新增一个卡牌文本工作稿，例如：

```text
docs/cards/base-set-text-audit.md
```

或如果你更想批量编辑，也可以使用：

```text
docs/cards/base-set-text-audit.csv
```

建议字段：

```text
卡号
名称
卡牌类型
当前描述
新规则文本
effect_categories
effect_optionality
触发时机
目标/选择
费用/额外费用
持续时间
实现备注
待确认问题
```

这一层方便你先把基础包卡牌描述改严谨，不影响当前游戏运行。

### 第二步：确认文本后重新生成编译定义

等一批卡牌文本确认后，先更新：

```text
data/all_card_infos.json
```

然后从 `server` 目录依次运行：

```powershell
cd server
go run ./cmd/extract-supported-cards
go run ./cmd/generate-card-definitions
go test ./...
```

其中 `extract-supported-cards` 会读取 `../data/all_card_infos.json`，筛选 `version_name == "基础包"`，并写入 `../data/supported_card_infos.json`。

生成器会读取 `../data/supported_card_infos.json`，并输出：

```text
server/cards/definitions_gen.go
server/cards/category_markers_gen.go
```

如果需要单独检查卡牌是否已经补齐结构化效果分类，可以从 `server` 目录运行：

```powershell
go run ./cmd/check-card-metadata
```

该命令只做数据审计，不参与游戏运行。当前旧数据没有填写 `effect_categories` / `effect_optionality` 时，它会返回非零并列出缺失项；等表格补齐后应逐步变成通过。

注意：

- 不要手改 `server/cards/definitions_gen.go` 或 `server/cards/category_markers_gen.go`。
- 生成器兼容带 UTF-8 BOM 的 JSON。
- `definitions_gen.go` 只同步结构化卡牌字段；运行时效果仍必须写 Go 行为，不能解析中文描述。

### 第三步：按机制和卡牌批次更新 Go 行为

卡牌文本确认后，再逐批更新 Go 行为和测试。建议顺序：

1. 无效果或纯关键词卡。
2. 入场、遗言、回合技、绝技。
3. 法术攻击、防御法术、强化法术。
4. 反制、反应、触发队列相关卡。
5. 光环、临时效果、绑定、衍生卡。

每批更新后运行：

```powershell
cd server
go test ./...
```

涉及前端窗口、支付、目标、反制、防御时，还需要做浏览器级检查。

## 4. 推荐你现在怎么改最方便

短期最方便的做法是：

1. 先不要直接改 `definitions_gen.go`。
2. 直接从表格导出并更新 `data/all_card_infos.json`。
3. 运行 `go run ./cmd/extract-supported-cards`，提取 `version_name == "基础包"` 的卡牌到 `data/supported_card_infos.json`。
4. 在 `docs/cards/base-set-text-review.md` 记录术语、规则冲突和实现备注。
5. 可选运行 `go run ./cmd/check-card-metadata`，检查效果分类字段是否漏填；当前数据尚未补齐分类时该命令会返回非零。
6. 文本确认后，从 `server` 目录运行 `go run ./cmd/generate-card-definitions`。
7. 再运行 `go test ./...`。

这样可以避免文本还没稳定就频繁改生成文件和 Go 行为。
