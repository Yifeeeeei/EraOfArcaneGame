package game

func (e *Engine) dealDamageWithExtra(target *CardInstance, amount, ownerID int, data map[string]any) {
	e.resolveDamage(target, amount, ownerID, data)
}
