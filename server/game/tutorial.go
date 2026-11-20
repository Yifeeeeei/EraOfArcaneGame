package game

import (
	"encoding/json"
	"eraofarcane/model"
	"fmt"
	"reflect"
)

const TutorialLessonID = "first-duel-v1"

type tutorialStep struct {
	id, chapter, title, text, instruction, focus, card string
	player                                             int
	action                                             string
	data                                               func(*Engine) map[string]any
}

type tutorialSession struct {
	steps   []tutorialStep
	index   int
	roles   map[string]string
	failure string
}

// TutorialDecks returns fresh, legal decks. Its declared order is private to the
// authored scenario; it never changes normal-game shuffling or deck visibility.
func TutorialDecks() (*model.Deck, *model.Deck) {
	makeDeck := func(hero string, types, opening, skills []string) *model.Deck {
		remaining := map[string]int{}
		for _, number := range types {
			remaining[number] = 2
		}
		main := append([]string{}, opening...)
		for _, number := range opening {
			remaining[number]--
		}
		for _, number := range types {
			for range remaining[number] {
				main = append(main, number)
			}
		}
		return &model.Deck{HeroID: hero, MainDeck: main, SkillPool: skills}
	}
	human := makeDeck("4111003", []string{
		"1021001", "1021002", "1021003", "1021004", "1021005", "1021013", "1121005", "1121006", "1121007", "1121015", "1121105", "2021001", "2021008", "2121006", "2121103",
	}, []string{"1021001", "2121006", "1021001", "2121006", "2021001", "2021001", "1121006", "1121006", "1121007"}, []string{
		"3121009", "3121012", "3121001", "3221002", "3321002", "3321005", "3321106", "3521006", "3621001", "3621005",
	})
	opponent := makeDeck("4211001", []string{
		"1021001", "1021002", "1021003", "1021004", "1021005", "1021013", "1221002", "1221003", "1221007", "1221009", "1221104", "1221113", "2021001", "2021008", "2221007",
	}, nil, []string{"3221002", "3121001", "3121009", "3321002", "3321005", "3321106", "3521006", "3621001", "3621005", "3121012"})
	return human, opponent
}

func arrangeScenarioDeck(ps *PlayerState) {
	queues := map[string][]*CardInstance{}
	for _, card := range ps.Deck {
		queues[card.Card.Number] = append(queues[card.Card.Number], card)
	}
	ordered := make([]*CardInstance, 0, len(ps.Deck))
	for _, number := range ps.DeckDef.MainDeck {
		ordered = append(ordered, queues[number][0])
		queues[number] = queues[number][1:]
	}
	ps.Deck = ordered
}

// SetupTutorialGame uses the real setup/turn/action rules and only fixes decks,
// first player and initial draw order. No gameplay state is patched mid-lesson.
func (e *Engine) SetupTutorialGame(name string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.tutorial != nil || e.State.Players[0] != nil {
		return fmt.Errorf("game already initialized")
	}
	p1, p2 := TutorialDecks()
	for _, deck := range []*model.Deck{p1, p2} {
		if err := deck.Validate(getCardDB()); err != nil {
			return err
		}
	}
	if err := e.setupGameWithOrder(name, p1, "教学对手 · 巴特尔", p2, 0, true); err != nil {
		return err
	}
	t := &tutorialSession{roles: map[string]string{}}
	e.tutorial = t
	t.roles["hero"] = e.State.Players[0].Hero.InstanceID
	t.roles["enemy"] = e.State.Players[1].Hero.InstanceID
	t.roles["apprentice"] = e.State.Players[0].Hand[0].InstanceID
	t.roles["mask"] = e.State.Players[0].Hand[1].InstanceID
	t.roles["spare_mask"] = e.State.Players[0].Hand[3].InstanceID
	for _, entry := range []struct {
		player       int
		number, role string
	}{{0, "3121009", "blast"}, {0, "3121012", "shield"}, {0, "3121001", "fireball"}, {1, "3221002", "ice"}} {
		for _, skill := range e.State.Players[entry.player].SkillPool {
			if skill.Card.Number == entry.number {
				t.roles[entry.role] = skill.InstanceID
			}
		}
	}
	t.steps = e.firstLessonSteps()
	return nil
}

func copyTutorialData(data map[string]any) map[string]any {
	raw, _ := json.Marshal(data)
	result := map[string]any{}
	_ = json.Unmarshal(raw, &result)
	return result
}

// Remove only empty optional protocol fields. Zero coordinates remain material.
func normalizedTutorialData(data map[string]any) map[string]any {
	result := copyTutorialData(data)
	for k, v := range result {
		if k == "payment" {
			if elements, ok := v.(map[string]any); ok {
				for element, amount := range elements {
					if amount == float64(0) {
						delete(elements, element)
					}
				}
			}
		}
		empty := v == nil
		switch value := v.(type) {
		case string:
			empty = value == ""
		case bool:
			empty = !value
		case []any:
			empty = len(value) == 0
		case map[string]any:
			empty = len(value) == 0
		}
		if empty {
			delete(result, k)
		}
	}
	return result
}

func (e *Engine) tutorialStepAction(step tutorialStep) ActionMessage {
	data := map[string]any{}
	if step.data != nil {
		data = copyTutorialData(step.data(e))
	}
	if step.action == "resolve_action" && e.State.PendingAction != nil {
		data["decision_id"] = e.State.PendingAction.DecisionID
	}
	action := step.action
	if action == "" {
		action = "tutorial_continue"
	}
	revision := e.revision
	return ActionMessage{Action: action, Data: data, ExpectedRevision: &revision}
}

func (e *Engine) validateTutorialAction(player int, action *ActionMessage) error {
	t := e.tutorial
	if t == nil {
		return nil
	}
	if action.ExpectedRevision == nil {
		return fmt.Errorf("教学战场已更新，请同步后重试")
	}
	if player == 0 && action.Action == "surrender" {
		return nil
	}
	if t.failure != "" {
		return fmt.Errorf("教学暂停，请重新开始：%s", t.failure)
	}
	if t.index >= len(t.steps) || e.State.Phase == PhaseGameOver {
		return fmt.Errorf("教学已结束")
	}
	step := t.steps[t.index]
	expected := e.tutorialStepAction(step)
	if player != step.player || action.Action != expected.Action {
		return fmt.Errorf("请按当前教学提示操作：%s", step.title)
	}
	// A unique-cost client may omit payment; pin that omission to the authored
	// payment rather than allowing the engine to pick a different resource.
	data := copyTutorialData(action.Data)
	if payment, ok := expected.Data["payment"]; ok && data["payment"] == nil {
		data["payment"] = payment
	}
	if decisionID, ok := expected.Data["decision_id"]; ok && data["decision_id"] == nil {
		data["decision_id"] = decisionID // Existing UI binds pending choices by revision.
	}
	if !reflect.DeepEqual(normalizedTutorialData(data), normalizedTutorialData(expected.Data)) {
		return fmt.Errorf("本步已限定卡牌、目标和支付方案，请按提示选择")
	}
	action.Data = data
	return nil
}

func (e *Engine) completeTutorialStep(player int, action ActionMessage) {
	t := e.tutorial
	if t == nil || t.index >= len(t.steps) || action.Action == "surrender" {
		return
	}
	step := t.steps[t.index]
	if player == step.player {
		t.index++
	}
	if e.State.Phase == PhaseGameOver && (e.State.Winner != 0 || t.index != len(t.steps)) {
		t.failure = "对局结果与教学流程不一致"
	}
}

// NextTutorialOpponentAction returns only one detached scripted action. Human
// pauses are never auto-acknowledged. ExpectedRevision closes the query/act gap.
func (e *Engine) NextTutorialOpponentAction() *ActionMessage {
	e.mu.Lock()
	defer e.mu.Unlock()
	t := e.tutorial
	if t == nil || t.failure != "" || t.index >= len(t.steps) || e.State.Phase == PhaseGameOver || t.steps[t.index].player != 1 {
		return nil
	}
	action := e.tutorialStepAction(t.steps[t.index])
	return &action
}

func (e *Engine) PauseTutorial(reason string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.tutorial != nil {
		e.tutorial.failure = reason
	}
}

func (e *Engine) tutorialInfo(player int) any {
	t := e.tutorial
	if t == nil || player != 0 {
		return nil
	}
	info := map[string]any{"lesson_id": TutorialLessonID, "index": t.index, "total": len(t.steps), "failure": t.failure,
		"completed": e.State.Phase == PhaseGameOver && e.State.Winner == 0 && t.index == len(t.steps)}
	if t.index >= len(t.steps) {
		return info
	}
	step := t.steps[t.index]
	info["id"] = step.id
	info["chapter"] = step.chapter
	info["title"] = step.title
	info["text"] = step.text
	info["instruction"] = step.instruction
	info["focus"] = step.focus
	info["waiting"] = step.player != 0
	info["reading"] = step.action == "" && step.player == 0
	info["lesson_skill_ids"] = []string{t.roles["blast"], t.roles["shield"], t.roles["fireball"]}
	if step.player == 0 {
		info["card_id"] = t.roles[step.card]
		info["allowed"] = e.tutorialStepAction(step)
	}
	return info
}
