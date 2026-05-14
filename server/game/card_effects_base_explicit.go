package game

import "eraofarcane/model"

type Card1021006Grocer struct{}

func (Card1021006Grocer) ID() string   { return "1021006" }
func (Card1021006Grocer) Name() string { return "杂货商贩" }
func (Card1021006Grocer) OnEnter(ctx *EffectContext) error {
	return DrawCards(2)(ctx)
}

type Card1121002LivelyHearth struct{}

func (Card1121002LivelyHearth) ID() string   { return "1121002" }
func (Card1121002LivelyHearth) Name() string { return "活泼的炉火" }
func (Card1121002LivelyHearth) OnEnter(ctx *EffectContext) error {
	return DrawCards(1)(ctx)
}

type Card1121014Firethorn struct{}

func (Card1121014Firethorn) ID() string   { return "1121014" }
func (Card1121014Firethorn) Name() string { return "火荆" }
func (Card1121014Firethorn) OnDeath(ctx *EffectContext) error {
	return ApplyStatusAuto(StatusBurn, 1)(ctx)
}

type Card1221004FrostPuppet struct{}

func (Card1221004FrostPuppet) ID() string   { return "1221004" }
func (Card1221004FrostPuppet) Name() string { return "寒霜傀儡" }
func (Card1221004FrostPuppet) OnEnter(ctx *EffectContext) error {
	if ctx.Target != nil {
		ctx.Target.Statuses[StatusFreeze]++
		return nil
	}
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	frontRow := opponent.GetFrontRow()
	if frontRow < 0 {
		return nil
	}
	for col := 0; col < 3; col++ {
		if opponent.Units[col][frontRow] != nil {
			opponent.Units[col][frontRow].Statuses[StatusFreeze]++
			return nil
		}
	}
	return nil
}

type Card1221008IcefieldDemon struct{}

func (Card1221008IcefieldDemon) ID() string   { return "1221008" }
func (Card1221008IcefieldDemon) Name() string { return "冰域恶魔" }
func (Card1221008IcefieldDemon) OnEnter(ctx *EffectContext) error {
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	for _, target := range ctx.Engine.getAllFieldCards(opponent) {
		target.Statuses[StatusFreeze]++
	}
	return nil
}

type Card1321002WindTraveler struct{}

func (Card1321002WindTraveler) ID() string   { return "1321002" }
func (Card1321002WindTraveler) Name() string { return "随风旅行者" }
func (Card1321002WindTraveler) OnEnter(ctx *EffectContext) error {
	ctx.Engine.State.Players[ctx.PlayerID].GainElements(map[string]int{model.ElementAir: 2})
	return nil
}
func (Card1321002WindTraveler) OnDeath(ctx *EffectContext) error {
	return DrawCards(1)(ctx)
}

type Card1321003MagicDandelion struct{}

func (Card1321003MagicDandelion) ID() string   { return "1321003" }
func (Card1321003MagicDandelion) Name() string { return "魔法蒲公英" }
func (Card1321003MagicDandelion) OnEnter(ctx *EffectContext) error {
	// Runtime still lacks per-instance draw-turn tracking. Keep the explicit
	// playable behavior that used to come from text parsing: entering draws 1.
	return DrawCards(1)(ctx)
}

type Card1321004LightningElemental struct{}

func (Card1321004LightningElemental) ID() string   { return "1321004" }
func (Card1321004LightningElemental) Name() string { return "雷电元素" }
func (Card1321004LightningElemental) OnEnter(ctx *EffectContext) error {
	return ApplyStatusAuto(StatusStun, 1)(ctx)
}

type Card1421001SandMage struct{}

func (Card1421001SandMage) ID() string   { return "1421001" }
func (Card1421001SandMage) Name() string { return "流沙法师" }
func (Card1421001SandMage) OnEnter(ctx *EffectContext) error {
	return ApplyStatusAuto(StatusPetrify, 1)(ctx)
}

type Card1421014WindbreathMerchant struct{}

func (Card1421014WindbreathMerchant) ID() string   { return "1421014" }
func (Card1421014WindbreathMerchant) Name() string { return "风息谷旅商" }
func (Card1421014WindbreathMerchant) OnEnter(ctx *EffectContext) error {
	count := 0
	for _, unit := range ctx.Engine.getAllFieldCards(ctx.Engine.State.Players[ctx.PlayerID]) {
		if unit == ctx.Source || !unit.Card.IsCompanion() {
			continue
		}
		if hasTag(unit.Card.Tag, "野兽", "精灵", "植物") {
			count++
			if count >= 3 {
				break
			}
		}
	}
	if count == 0 {
		return nil
	}
	return DrawCards(count)(ctx)
}

type Card1521002LightforgedTitan struct{}

func (Card1521002LightforgedTitan) ID() string   { return "1521002" }
func (Card1521002LightforgedTitan) Name() string { return "光铸泰坦" }
func (Card1521002LightforgedTitan) OnEnter(ctx *EffectContext) error {
	return DrawCards(2)(ctx)
}

type Card1521014TorchWitch struct{}

func (Card1521014TorchWitch) ID() string   { return "1521014" }
func (Card1521014TorchWitch) Name() string { return "炬之女巫" }
func (Card1521014TorchWitch) OnEnter(ctx *EffectContext) error {
	return ApplyStatusToSelf(StatusBurn, 2)(ctx)
}

type Card1521015EmberWitch struct{}

func (Card1521015EmberWitch) ID() string   { return "1521015" }
func (Card1521015EmberWitch) Name() string { return "烬之女巫" }
func (Card1521015EmberWitch) OnEnter(ctx *EffectContext) error {
	return ApplyStatusToSelf(StatusBurn, 3)(ctx)
}

type Card1611001ObserverOkoru struct{}

func (Card1611001ObserverOkoru) ID() string   { return "1611001" }
func (Card1611001ObserverOkoru) Name() string { return "\"观察者\" 欧柯茹" }
func (Card1611001ObserverOkoru) OnEnter(ctx *EffectContext) error {
	if err := DrawCards(1)(ctx); err != nil {
		return err
	}
	return DealDamageToSelfHero(1)(ctx)
}

type Card1621001UnderworldPigeon struct{}

func (Card1621001UnderworldPigeon) ID() string   { return "1621001" }
func (Card1621001UnderworldPigeon) Name() string { return "冥界信鸽" }
func (Card1621001UnderworldPigeon) OnDeath(ctx *EffectContext) error {
	return DrawCards(1)(ctx)
}

type Card1621005CursedGolem struct{}

func (Card1621005CursedGolem) ID() string   { return "1621005" }
func (Card1621005CursedGolem) Name() string { return "诅咒魔像" }
func (Card1621005CursedGolem) OnEnter(ctx *EffectContext) error {
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	for _, skill := range opponent.Skills {
		if skill != nil {
			skill.Statuses[StatusWeaken] += 2
			return nil
		}
	}
	return nil
}

type Card2321007WindwhisperRing struct{}

func (Card2321007WindwhisperRing) ID() string   { return "2321007" }
func (Card2321007WindwhisperRing) Name() string { return "风语之戒" }
func (Card2321007WindwhisperRing) OnEnter(ctx *EffectContext) error {
	return DrawCards(1)(ctx)
}

type Card2601002Spellbook struct{}

func (Card2601002Spellbook) ID() string   { return "2601002" }
func (Card2601002Spellbook) Name() string { return "咒言书" }
func (Card2601002Spellbook) OnEnter(ctx *EffectContext) error {
	opponent := ctx.Engine.State.Players[ctx.OpponentID]
	for _, skill := range opponent.Skills {
		if skill != nil {
			skill.Statuses[StatusWeaken]++
		}
	}
	return nil
}
