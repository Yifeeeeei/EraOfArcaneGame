package game

// precondition.go — 通用前提条件检查系统
//
// 设计原则:
// 1. 可组合: 多个条件可以组合使用 (AllOf, AnyOf)
// 2. 可扩展: 新条件类型可以方便地添加
// 3. 可复用: 常用条件作为预定义常量提供
// 4. 错误信息: 检查失败时提供清晰的错误信息

import (
	"fmt"
	"strings"
)

// ═══════════════════════════════════════════════════════════
// 核心类型定义
// ═══════════════════════════════════════════════════════════

// Precondition 是一个可执行的条件检查
type Precondition interface {
	// Check 执行条件检查，返回检查结果和可选的错误信息
	Check(ctx *ConditionContext) PreconditionResult

	// String 返回条件的描述
	String() string
}

// PreconditionResult 检查结果
type PreconditionResult struct {
	Passed  bool   // 是否通过
	Reason  string // 失败原因（如果失败）
	Details string // 详细信息
}

// ConditionContext 条件检查的上下文
type ConditionContext struct {
	Engine     *Engine        // 游戏引擎
	PlayerID   int            // 当前玩家ID
	SourceCard *CardInstance  // 来源卡牌（如果是卡牌效果）
	ExtraData  map[string]any // 额外数据
}

// ═══════════════════════════════════════════════════════════
// 组合条件
// ═══════════════════════════════════════════════════════════

// AllOf 要求所有条件都满足
type AllOf struct {
	Conditions []Precondition
	Name       string // 可选的名称
}

func (a *AllOf) Check(ctx *ConditionContext) PreconditionResult {
	var failedReasons []string

	for _, cond := range a.Conditions {
		result := cond.Check(ctx)
		if !result.Passed {
			failedReasons = append(failedReasons, result.Reason)
		}
	}

	if len(failedReasons) > 0 {
		return PreconditionResult{
			Passed: false,
			Reason: fmt.Sprintf("需要满足所有条件，但 %d 个条件未满足: %s",
				len(failedReasons), strings.Join(failedReasons, "; ")),
		}
	}

	return PreconditionResult{Passed: true}
}

func (a *AllOf) String() string {
	if a.Name != "" {
		return a.Name
	}
	return fmt.Sprintf("满足所有条件: %v", a.Conditions)
}

// AnyOf 要求至少一个条件满足
type AnyOf struct {
	Conditions []Precondition
	Name       string
}

func (o *AnyOf) Check(ctx *ConditionContext) PreconditionResult {
	var details []string

	for _, cond := range o.Conditions {
		result := cond.Check(ctx)
		if result.Passed {
			return PreconditionResult{
				Passed:  true,
				Details: fmt.Sprintf("条件满足: %s", cond.String()),
			}
		}
		details = append(details, result.Reason)
	}

	return PreconditionResult{
		Passed: false,
		Reason: fmt.Sprintf("需要满足至少一个条件，但所有条件都未满足: %s",
			strings.Join(details, "; ")),
	}
}

func (o *AnyOf) String() string {
	if o.Name != "" {
		return o.Name
	}
	return fmt.Sprintf("满足任一条件: %v", o.Conditions)
}

// Not 否定一个条件
type Not struct {
	Condition Precondition
}

func (n *Not) Check(ctx *ConditionContext) PreconditionResult {
	result := n.Condition.Check(ctx)
	if result.Passed {
		return PreconditionResult{
			Passed: false,
			Reason: fmt.Sprintf("要求条件不满足，但条件已满足: %s", n.Condition.String()),
		}
	}
	return PreconditionResult{Passed: true}
}

func (n *Not) String() string {
	return fmt.Sprintf("不满足: %s", n.Condition.String())
}

// ═══════════════════════════════════════════════════════════
// 具体条件实现
// ═══════════════════════════════════════════════════════════

// PlayerHasFieldUnits 检查玩家是否有场上单位
type PlayerHasFieldUnits struct {
	PlayerID int // -1 表示当前玩家, -2 表示对手
}

func (p *PlayerHasFieldUnits) Check(ctx *ConditionContext) PreconditionResult {
	pid := p.PlayerID
	if pid == -1 {
		pid = ctx.PlayerID
	} else if pid == -2 {
		pid = 1 - ctx.PlayerID
	}

	ps := ctx.Engine.State.Players[pid]
	hasUnits := false
	for col := 0; col < 3 && !hasUnits; col++ {
		for row := 0; row < 3 && !hasUnits; row++ {
			if ps.Units[col][row] != nil {
				hasUnits = true
			}
		}
	}

	if !hasUnits {
		return PreconditionResult{
			Passed: false,
			Reason: "场上没有单位",
		}
	}
	return PreconditionResult{Passed: true}
}

func (p *PlayerHasFieldUnits) String() string {
	switch p.PlayerID {
	case -1:
		return "当前玩家场上有单位"
	case -2:
		return "对手场上有单位"
	default:
		return fmt.Sprintf("玩家%d场上有单位", p.PlayerID)
	}
}

// PlayerHasHandCards 检查玩家是否有手牌
type PlayerHasHandCards struct {
	PlayerID int
	MinCount int // 最少手牌数，0表示只要有手牌
}

func (p *PlayerHasHandCards) Check(ctx *ConditionContext) PreconditionResult {
	pid := p.PlayerID
	if pid < 0 {
		pid = ctx.PlayerID
	}

	ps := ctx.Engine.State.Players[pid]
	handCount := len(ps.Hand)

	if handCount < p.MinCount {
		if p.MinCount == 1 {
			return PreconditionResult{
				Passed: false,
				Reason: "没有手牌",
			}
		}
		return PreconditionResult{
			Passed: false,
			Reason: fmt.Sprintf("手牌不足，需要%d张，当前有%d张", p.MinCount, handCount),
		}
	}
	return PreconditionResult{Passed: true}
}

func (p *PlayerHasHandCards) String() string {
	if p.MinCount <= 1 {
		return "有手牌"
	}
	return fmt.Sprintf("手牌数≥%d", p.MinCount)
}

// PlayerHasCharge 检查玩家是否有足够的充能
type PlayerHasCharge struct {
	MinCharge int
}

func (p *PlayerHasCharge) Check(ctx *ConditionContext) PreconditionResult {
	ps := ctx.Engine.State.Players[ctx.PlayerID]
	if ps.Charge < p.MinCharge {
		return PreconditionResult{
			Passed: false,
			Reason: fmt.Sprintf("充能不足，需要%d点，当前有%d点", p.MinCharge, ps.Charge),
		}
	}
	return PreconditionResult{Passed: true}
}

func (p *PlayerHasCharge) String() string {
	return fmt.Sprintf("充能≥%d", p.MinCharge)
}

// CardHasStatus 检查卡牌是否有指定状态
type CardHasStatus struct {
	Status    string
	MinAmount int    // 最小层数，0表示只要有
	Target    string // "source" 或 "target"
}

func (c *CardHasStatus) Check(ctx *ConditionContext) PreconditionResult {
	var card *CardInstance
	if c.Target == "target" {
		card = ctx.ExtraData["target_card"].(*CardInstance)
	} else {
		card = ctx.SourceCard
	}

	if card == nil {
		return PreconditionResult{
			Passed: false,
			Reason: "没有目标卡牌",
		}
	}

	amount := card.Statuses[c.Status]
	if amount < c.MinAmount {
		if c.MinAmount == 0 {
			return PreconditionResult{
				Passed: false,
				Reason: fmt.Sprintf("卡牌没有%s状态", c.Status),
			}
		}
		return PreconditionResult{
			Passed: false,
			Reason: fmt.Sprintf("%s层数不足，需要%d层，当前有%d层", c.Status, c.MinAmount, amount),
		}
	}
	return PreconditionResult{Passed: true}
}

func (c *CardHasStatus) String() string {
	if c.MinAmount <= 0 {
		return fmt.Sprintf("有%s状态", c.Status)
	}
	return fmt.Sprintf("%s层数≥%d", c.Status, c.MinAmount)
}

// ═══════════════════════════════════════════════════════════
// 条件上下文辅助函数
// ═══════════════════════════════════════════════════════════

// HasFieldUnits 检查玩家场上是否有单位（快捷函数）
func HasFieldUnits(playerID int) Precondition {
	return &PlayerHasFieldUnits{PlayerID: playerID}
}

// HasHandCards 检查玩家是否有手牌（快捷函数）
func HasHandCards(playerID int, minCount int) Precondition {
	return &PlayerHasHandCards{PlayerID: playerID, MinCount: minCount}
}

// HasCharge 检查玩家是否有足够充能（快捷函数）
func HasCharge(minCharge int) Precondition {
	return &PlayerHasCharge{MinCharge: minCharge}
}

// SourceHasStatus 检查来源卡牌是否有状态（快捷函数）
func SourceHasStatus(status string, minAmount int) Precondition {
	return &CardHasStatus{Status: status, MinAmount: minAmount, Target: "source"}
}

// TargetHasStatus 检查目标卡牌是否有状态（快捷函数）
func TargetHasStatus(status string, minAmount int) Precondition {
	return &CardHasStatus{Status: status, MinAmount: minAmount, Target: "target"}
}
