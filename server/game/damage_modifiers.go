package game

// DamageModifierConsumption separates paying for a reduction from calculating
// it. Modifier queries may be repeated; consumption occurs once on commit.
type DamageModifierConsumption interface {
	ConsumeDamageModifier(*EffectContext, int, int)
}

type damageModifierUse struct {
	behavior      DamageModifierConsumption
	context       *EffectContext
	before, after int
}

type damageModifierPlan struct {
	Amount int
	uses   []damageModifierUse
}

// commit is idempotent. A plan is prepared and committed within one engine
// operation, before any replacement choice can suspend execution.
func (plan *damageModifierPlan) commit() {
	uses := plan.uses
	plan.uses = nil
	for _, use := range uses {
		use.behavior.ConsumeDamageModifier(use.context, use.before, use.after)
	}
}

func (plan *damageModifierPlan) record(behavior CardBehavior, ctx *EffectContext, before int) {
	if plan.Amount == before {
		return
	}
	if consumer, ok := behavior.(DamageModifierConsumption); ok {
		plan.uses = append(plan.uses, damageModifierUse{consumer, ctx, before, plan.Amount})
	}
}

// planDamageModifiers is a pure query. The target layer runs before friendly
// field reductions; both use stable field order and the same incoming event.
func (e *Engine) planDamageModifiers(target *CardInstance, amount, ownerID int, data map[string]any) damageModifierPlan {
	plan := damageModifierPlan{Amount: max(0, amount)}
	if target == nil || target.Card == nil || plan.Amount == 0 || ownerID < 0 || ownerID >= len(e.State.Players) {
		return plan
	}
	context := func(source *CardInstance) *EffectContext {
		return &EffectContext{Engine: e, Source: source, Target: target, PlayerID: ownerID, OpponentID: 1 - ownerID, ExtraData: data}
	}
	behavior := cardBehavior(target)
	if modifier, ok := behavior.(DamageAmountModifier); ok && modifier.HasActiveDamageAmountModifier(target) {
		ctx, before := context(target), plan.Amount
		plan.Amount = max(0, modifier.ModifyDamageAmount(ctx, before))
		plan.record(behavior, ctx, before)
	}
	for _, source := range e.getAllFieldCards(e.State.Players[ownerID]) {
		if plan.Amount == 0 {
			break
		}
		if source == nil || source.Card == nil || source == target || e.hasEffectiveStatus(source, StatusPetrify) {
			continue
		}
		behavior := cardBehavior(source)
		modifier, ok := behavior.(FieldDamageAmountModifier)
		if !ok || !modifier.HasActiveFieldDamageAmountModifier(source) {
			continue
		}
		ctx, before := context(source), plan.Amount
		plan.Amount = max(0, modifier.ModifyFieldDamageAmount(ctx, before))
		plan.record(behavior, ctx, before)
	}
	return plan
}
