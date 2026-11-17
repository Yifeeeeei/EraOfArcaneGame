# 前端独立质量复核（2026-10-10）

复核对象：`86b117c` 的重设计；只在云端通过 CLI、API、无头 Chromium 操作。该轮复核没有上传文件、推送、创建 PR、合并或部署。后续经用户授权的 GitHub 草稿交付见 [交付材料](review/frontend-f2b78f5/README.md)。

## 入口核对

新旧模板的 click/contextmenu 调用处理器集合比对没有遗漏；引擎 HandleAction 的 18 种游戏动作在新 JS 均有入口（client_log 不属于游戏动作）。这只是静态核对，不把它当作所有分支可达的证明。

真实页面新增检查 `tests/frontend_quality.py`：

- 已存在大法师之杖 → 手牌同类装备 → 替换选择 → 付款 → 提交 → 选择绑定雷击 → 从宿主详情打开绑定技能。
- 草稿存在时打开卡牌检视并关闭，再提交；此前关闭按钮会被 capture 锁拦截，现已修复。
- 星落之银叶在场 → 回合末弃牌 → 连续收纳选择 → 从技能详情打开卡下方的卡牌。
- 同场第三个观战客户端：没有双方私有手牌、没有 pending 候选、没有投降入口；当前回合标记恰好一个。
- 真实过期 expected_revision 请求被拒绝且 revision 不变；投降先展示草稿，再确认，双方与观众都收到终局。

原页面已有绑定技能入口保留。`under_cards` 和 `attached_behaviors` 此前仅有协议数据，原页面和第一版重设计都没有详情区；本次补上实际详情入口，并更正覆盖台账。附属效果名称/数据采用服务器结构化字段，不解析卡文。实际截图还发现详情费用图标受云端外部图源限制而破图，详情费用改为完整元素名称和数量，避免图标缺失时无法辨认。

测试工具面板在新对局页隐藏；独立 `/card-test.html` 仍为完整测试管理入口。这是诊断入口安排，不是某种游戏机制退回 legacy。没有按动作/卡号自动跳到旧页的分支。完整旧页仅为手动回退链接。

## revision 与 decision_id 审核

外部可调用的状态变更路径是 WebSocket → HandleAction、房间初始化 → SetupGame/SetupGameWithFirstPlayer，以及测试房间的 DebugAddCardAt/DebugUpdateElements/DebugMutateCard。规则辅助方法在这些动作内部运行，不逐帧向 UI 暴露版本或内部树。

- 所有进入规则执行的 HandleAction（包含 error）在引擎锁内递增 revision；expected_revision 在任何规则执行前检查。
- 初始化原先未递增，本次在共同 setup 实现补齐。
- DebugAddCard 多张入槽原先可能部分成功后返回错误；元素 map 和状态 map 也可能边校验边写入。本次先检查容量、构建元素结果、预校验状态，再写入，避免这些错误路径无版本地部分改动。新增 Go 测试覆盖失败原子性、成功递增和旧草稿拒绝。
- 所有 PendingAction 激活统一经过 activatePendingAction；排队的同类型窗口激活时取得新 ID，失败后恢复同一窗口保留原 ID。
- decision_id 不是动作提交的独立前置条件；当前提交使用 expected_revision。防御、起手、主行动等非 PendingAction 窗口没有 decision_id，依赖 revision 和 phase。

**后续已修复执行报错的版本/广播缺口。** 所有进入规则执行的请求现在无论成功或 error 都推进 revision，并在解锁后广播权限过滤快照；执行前拒绝保留原版本。失败的部分修改仍不回滚。详细边界和实际血沙部分付款测试见 [执行报错同步](frontend-failed-action-sync.md)。

另外，旧客户端允许省略 expected_revision；没有跨重启 epoch、持久 request_id 去重缓存或自动重放保证。revision/decision_id 是进程内对局信息。

## 尚未逐项验证

- 728 张卡全部组合；动态构造的全部 pending 类型和所有专属 payload 组合。
- 绑定技能的每一种施放/防御路径；本次新增真实宿主绑定与检视，原防御集合静态/Go 检查保留。
- 附属效果列表已接 UI，但本次没有新增真实附属效果生成的浏览器旅程。
- 所有区域法术、技能替换、吞噬/复制/控制转移等专属组合仍按覆盖台账标注，不能仅凭旧入口迁移宣称逐卡通过。
- 真机、读屏器与全面无障碍审计；外部图源不可达时仍只能显示文字。

可复跑命令与总体限制见 frontend-design.md、frontend-coverage.md。最终复核输出保存在忽略的 agent-data/frontend-review/quality-*.log，不上传。
