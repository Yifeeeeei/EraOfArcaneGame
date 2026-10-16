# 对局工作区设计与审查说明

整合基线：PR #163，`dc2fa01196fea6be06313d048d7098b7871eebeb`（2026-10-10 API 复核为 open、未合并）。main 当时为 `862e8de520f5960c5753b349c642518b38ae955c`。实现位于独立 `codex/frontend-unified-workspace`，不改原 PR 分支。实现完成时未发布；后续经用户授权创建依赖 #163 的前端草稿 PR，见 [交付材料](review/frontend-f2b78f5/README.md)。

## 信息架构

- 状态条分开显示回合玩家、当前决策玩家、连接、房间与回合数。
- 桌面双战场、人物/资源栏、手牌与唯一决策工作区；技能池/公开事件使用切换页签。页面按内容高度滚动，决策栏保持可达，不以裁掉卡面换取一屏容纳。卡牌检视使用单独详情层。
- 深褐桌面、方形控件、细金边。单位始终显示名称、生命、状态，选择有勾选文字；有点击检视与明确学习/施放入口，不依赖悬停。
- 手机按资源、双方战场、手牌、决策、技能/事件纵向排列，装备/技能行与密集手牌自身横向滚动，页首锚点导航直达当前决策（导航不浮盖卡图）。手机不承诺整场同时塞入一屏。
- `web/game-legacy.html` 是基线完整原页。新页无 emergency DOM、无第二套 pending 面板；旧页只是独立回退入口。

## 数据和提交边界

`gameState` 是过滤后的服务器事实。pending 事件只进入叙事，真正可操作的窗口来自 state_sync。本地目标/候选/排序/数量/付款形成草稿；`sendAction` 只暂存，只有 `commitDraft` 发 WebSocket。

`web/js/decision-session.js` 单独管理在途状态，可用 Node 直接测试。执行前拒绝保留草稿；执行报错保持锁定至新快照清除旧草稿，原错误仍就地呈现；成功回执后仍等待较新状态；较新状态可先于回执到达，旧回执不会清除新窗口。重复点击不能开启第二个在途动作。断线不自动重连重发，用户可手动重连核对。

付款统一支持增、减、重置和明确确认。反制/响应/防御的透支计入同一付款池；返回编辑保留原选择。规则上的“不防御/放弃发动”和“返回编辑/取消本地选择”分别呈现。无目标法术也先选强化，再准备施放。

## 最小协议增强

- `revision`：引擎锁下每次进入规则执行的请求（成功或 error）及成功测试工具写入递增；新版动作带可选 `expected_revision`，过时草稿先于规则处理被拒绝。旧客户端不带字段仍兼容。
- `decision_id`：当前 PendingAction 激活时分配，不由类型、候选、文字指纹猜测。相同类型连续窗口也不同。
- `decision_player`：公开的当前交接对象，和 current_turn 分离。观众的等待信息只显示通用文案，不再显示可能揭露盖牌身份的私有 prompt/type。
- `ui`：候选选择、牌顶/牌底排序、整数分配、负载重分配、负载/生命付款。现有具名窗口的协议投影集中在 `decision_view.go`；新机制需声明原语，不能在前端增加卡号/卡文判断。
- 数量分配通过 `allocations` 映射校验总数后再适配既有 callback；旧客户端不传映射仍保留既有默认分配。未重写结算帧或规则调度。

这些字段是过期校验，不是持久化幂等：没有服务端 request_id 结果缓存，不解决跨重启恢复、旧客户端重复提交或全部既有 callback 的事务回滚。结算帧、种子、隐藏队列从未传给前端。

## 模块与依赖

保留静态 Vue，无 npm 生产构建。模板在 `game.html`，真实机制适配在 `js/game.js`，传输状态模块独立，新增视觉层在 `css/workspace.css`。Vue 固定 3.5.22 本地随包提供，许可证在 `web/vendor/VUE-LICENSE`。旧版仍使用基线 CDN。

卡图继续使用原资源，并提供同一 ArcaneImages 仓库 raw 地址故障回退。两条路径都不可用时卡名/状态仍在；此实现不内置整套卡图。云环境原 GitHub Pages 域名被网络策略阻断；无头 Chromium 通过同仓库 raw 原图验证，使用测试环境 `ignore_https_errors=True`，生产代码没有绕过 TLS。

## 运行和审查

在 `server/` 用 Go >= 1.26.1 运行 `go run .`，浏览器访问 `/game.html`；从现有大厅/测试工作台创建的链接直接进入新页。无额外前端构建命令。

自动化：

```sh
cd server
go test ./...
go vet ./...
go test -race ./...
cd ..
node web/js/game-costs.test.js
node web/js/decision-session.test.js
# 另一个终端已经从 server/ 启动 go run .
python tests/frontend_match.py
python tests/frontend_reactions.py
python tests/frontend_primitives.py
python tests/frontend_dense.py
python tests/frontend_card_visibility.py
python tests/frontend_quality.py
python tests/frontend_failed_action.py
```

Python 检查需要 Playwright、Pillow 和 `/usr/bin/chromium`，只运行 headless。`ARCANE_URL` 可指定测试服务器。脚本只创建和修改 test-room，真实动作通过两份实际 Vue 页面/WebSocket；部分复杂 fixture 入口借助页面已存在的 debug.sendAction 进入同一个确认/提交路径。断线、错误与重复点击有专门断言。

截图、房间日志及临时结果在忽略的 `agent-data/frontend-review/` / `server/logs/rooms/`，不提交玩家记录。机制覆盖与未穷举项见 [覆盖台账](frontend-coverage.md)，289 种静态命名 pending 创建点见 [决策清单](frontend-decision-inventory.md)。

执行报错后的同步、锁顺序和具体保证见 [执行报错同步](frontend-failed-action-sync.md)。

## 卡面完整性回归

PR #165 的初始版将单位图固定为 86px 高并使用 `object-fit:cover`，同时以绝对定位叠放名称、状态和操作；手牌、技能池与手机布局还有独立的高度约束。2026-10-16 的修复统一使用原图比例及正常文档流，将所有文字/按钮放到卡图外。技能池提供触控可用的查看入口；详情弹窗的关闭按钮独占一行。

桌面决策栏从最多 390px 收窄到 300px（超宽屏 340px）；技能池独占手牌下方区域。821–1100px 与手机纵向排列双方战场，避免六列卡牌挤在窄屏；卡图完整显示仍不等于缩略图上每个印刷字都清晰，完整效果可点击详情阅读。满场页面会比初版更长，这是保留完整卡面的取舍。

`frontend_card_visibility.py` 在真实双客户端房间检查满场、密集手牌、横置/状态、技能/装备、最后一张手牌及详情弹窗，视口为 1440×1000、1024×768、390×844。图片必须成功加载、保持自然比例、不与兄弟控件相交、不被卡容器上下裁切；页面不能横向溢出。`frontend_match.py` 和 `frontend_reactions.py` 复用检查验证服务器实际发出的防御与反制窗口。视觉产物在忽略目录 `agent-data/card-visibility/`；自动化几何断言不能替代查看这些实际截图。
