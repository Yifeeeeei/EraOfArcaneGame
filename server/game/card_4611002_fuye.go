package game

type Card4611002Fuye struct{}

func (Card4611002Fuye) ID() string   { return "4611002" }
func (Card4611002Fuye) Name() string { return "芙雅夫人" }
func (Card4611002Fuye) OnUltimate(ctx *EffectContext) error {
	if ctx.Target != nil {
		ctx.Target.CurrentAttack *= 2
		gain := effectiveElementsGain(ctx.Target)
		doubled := make(map[string]int, len(gain))
		for elem, amount := range gain {
			doubled[elem] = amount * 2
		}
		setElementsGain(ctx.Target, doubled)
		ctx.Target.Statuses["临时"] = 1
	}
	return nil
}
