# 桌面对阵与扇形手牌评审

本地工作基线：`865c211f3211d6da2b367a8909fc5b6b80bdae49`，分支 `codex/frontend-unified-workspace`。截图对应布局提交 `cd176f8160f304c736df4481af0c807d40d67a48`。本轮后续授权将更新原草稿 PR #165；不合并、不部署。

这批图片由云端真实 Go 服务、两份 Vue 客户端和 headless Chromium 生成，并逐张查看。测试使用 `/api/test-room/*` 建立密集场景，游戏动作经过实际网页和 WebSocket；无用户窗口操作。没有上传用户草图或玩家记录。Library 草图下载失败，布局依据用户已确认的文字描述。

| 图片 | 视口 | 检查内容 |
| --- | --- | --- |
| [满场](dense-1440.png) | 1440×1000 | 双方 18 单位、四条技能/装备栏、动态前排、固定资源与手牌 |
| [较矮桌面](dense-1366.png) | 1366×768 | 战场需局部滚动；右下手牌和右上工作区保持固定 |
| [窄桌面滚到底部](dense-1024-arena-bottom.png) | 1024×768 | 己方后排可达；较长技能/装备栏需独立横向滚动 |
| [32 张手牌与预览](fan-interaction-1440.png) | 1440×1000 | 分页、稳定预览、查看/选择入口；预览不占战场 |
| [窄桌面预览](fan-interaction-1024.png) | 1024×768 | 扇形边缘保留在手牌区；预览正文可滚动 |
| [实际反制窗口](counter-1366.png) | 1366×768 | 回合为对方、决策为自己；候选与规则上的放弃按钮 |
| [防御滚到底部](defense-skills-1024-bottom.png) | 1024×768 | 候选独立滚动，汇总/确认/不防御固定在底栏 |

## 验证

- `go test -count=1 ./...`：通过；必须禁用缓存，因为模板测试读取 Go 包外的 HTML/CSS。原模板检查的手牌强化选择器已更新到集中入口，另有真实卷轴强化增减检查。
- `go vet ./...`：通过。
- Node `decision-session.test.js`、`game-costs.test.js` 与 `node --check web/js/game.js`：通过。
- `frontend_match.py`、`frontend_reactions.py`、`frontend_quality.py`、`frontend_primitives.py`、`frontend_failed_action.py`：通过。包括反制→嵌套选择→响应→防御、连续红月触发刷新、重复点击、错误保留/重试、断线不重发、旁观者隐私和付款增减/重置。
- `frontend_card_visibility.py`、`frontend_fan.py`：通过。仅运行三个桌面视口；32 张牌的实例级键盘导航、稳定鼠标命中、悬停零提交、无效牌、选择/取消、平铺/间距、卷轴强化增减。空手牌边界使用独立前端展示 fixture，未冒充真实引擎弃牌流程。
- `go test -race -count=1 ./...` 最终复跑通过；此前全量与单包运行均发现 race，因此不视为稳定通过。Race：出现可重复但非每次必现的 API 全局行为注册竞争。`RegisterAllCardEffects` (`card_effects_catalog.go:652`) 在 `pr163StartMatch` (`pr163_api_matrix_test.go:128`) 写全局 registry；上一个测试房间 WebSocket 的 `BroadcastState` 经 `behaviorForNumber` (`card_traits.go:287`) 同时读。单独 `go test -race ./api` 也复现。相关后端生产代码本轮没有改动，不把一次复跑通过视为已解决。当前日志位于忽略/临时目录，未将房间原始日志提交。

## 边界

扇形仍是供评审的交互探索：七张一页重叠较多，窄桌面会进一步限制间距，可切换平铺。满场卡面的小字应使用详情阅读。768px 高度不能同时看完整双方所有排，使用战场局部滚动；长栏同理。状态/地形/附属进一步增加时仍可能需要滚动。

仅验证桌面，未新增移动端适配、触控设备测试或手机截图。没有穷举全部 728 张卡组合。浏览器测试为了云环境证书链使用 `ignore_https_errors=True`，生产页面未关闭 TLS 校验；卡图仍依赖远端图片服务及现有 raw GitHub 回退。截图和几何断言分别进行，API 通过不代表视觉通过。

## 发布前复核与基线 race 诊断

已实际查看本目录最终七张 PNG，另重新打开满场和窄桌面预览原图复核。结论：单位与技能/装备卡图保持完整比例，附加状态/控件在图外；扇形本身按设计遮叠邻牌，可抬升或在固定预览区读全牌。对手前排朝下、己方前排朝上，向中线相对，动态前排文字与所在行一致。四条长栏归属清楚。扇形外缘不越出手牌区，预览/决策不覆盖战场或确认底栏。窄桌面中预览按钮、战场后排和长栏末尾需要各自区域滚动，并非承诺所有内容同时显示。

独立基线复现：将远端仍指向的 `865c211f3211d6da2b367a8909fc5b6b80bdae49` 用 `git archive` 解包到 `/tmp/era-race-baseline`，未复制本轮改动。Go 1.26.1，执行位置均为该目录的 `server/`：

```sh
# 实际运行的缩小复现（单个既有测试，不新增测试或修改引擎）：exit 1
GOCACHE=/tmp/era-go-cache GOPATH=/tmp/era-go-path \
  /tmp/era-toolchain/go/bin/go test -race -count=10 \
  -run '^TestPR163APIRedMoonExtraTargetAndDefenseHandshake$' ./api

# 同一原始基线的 API 包复现：exit 1
GOCACHE=/tmp/era-go-cache GOPATH=/tmp/era-go-path \
  /tmp/era-toolchain/go/bin/go test -race -count=3 ./api
```

单次 `-count=1` 跑同一测试得到 exit 0，说明调度相关、并非每次必现。当前布局提交先前运行 `go test -race ./...` 和 `go test -race ./api` 均复现；最终 `go test -race -count=1 ./...` 得到 exit 0。不能把最后一次成功表述为竞争已解决。

基线缩小复现中的冲突堆栈：

- 写：`game.RegisterAllCardEffects`，`server/game/card_effects_catalog.go:652`，赋值 `globalRegistry = NewEffectRegistry()`。
- 调用写入的测试入口：`server/api/pr163_api_matrix_test.go:128`；同一测试两个 subtest 从 `:345` 创建房间，`:346–348` 关闭服务器/连接。
- 读：`server/game/rules.go:207` 的 `effectiveCardPlayCost`，以及 `server/game/card_traits.go:287` 的 `behaviorForNumber`，访问同一个全局 registry。
- 仍在运行的旧房间路径：`server/api/ws.go:250` → `server/match/room.go:401` → `server/game/state_view.go:102` → 卡牌序列化/费用计算。

据堆栈判断，测试关闭连接后旧 WebSocket handler 的广播尚未结束，下一个测试/子测试就重新注册了全局行为表；这不是前端手牌交互导致的竞争。生产入口通常只启动注册一次，本次证据不能扩展为“正常线上每场都会发生”。全局注册是否应禁止重复/并发，或测试是否应等待 handler 退出，应在单独后端修复中决定。

`git diff 865c211f HEAD -- server ':!server/game/frontend_template_test.go'` 为空；本轮 server 目录唯一变化是适配新强化入口的既有模板断言，没有混入后端生产修复。原始运行日志保留在临时目录，仅提交本诊断，不发布房间日志或用户原图。
