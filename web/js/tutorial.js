/* Presentation and input hints only. The engine independently enforces steps. */
(function (root) {
  function permits(course, action, id) {
    if (!course) return true;
    if (course.failure || course.waiting || (course.reading && action !== 'tutorial_continue') || course.allowed?.action !== action) return false;
    if (!id) return true;
    const data = course.allowed.data || {};
    return data.instance_id === id || (data.skill_ids || []).includes(id) || (data.overexert_ids || []).includes(id) || (data.selected || []).includes(id);
  }
  function defenseHint(course, step, selected = [], overexert = []) {
    if (course?.allowed?.action !== 'defend' || course.reading || course.waiting) return null;
    const data = course.allowed.data;
    const shield = data.skill_ids[0], mask = data.overexert_ids[0];
    const canReview = selected.includes(shield) && overexert.includes(mask);
    const hints = [
      { title: '① 防御：选择主防御法术', text: '点击下方「选择烈焰护体」；也可点击己方技能区高亮的卡面。它的威力3足以挡住冰锥术，选中后自动进入②强化。卡牌上的「查看」只打开详情，不会选牌。此时还没有施法或扣费。', instruction: '点击右侧「选择烈焰护体」。选中后会自动进入②强化。', cardId: shield, choiceLabel: '选择烈焰护体', selector: '.tutorial-defense-choice' },
      { title: '② 强化：本次跳过', text: '强化可以叠加防御威力，但本次主法术已经足够，不需要再选牌。点击本面板下方「跳过强化 → 透支」，进入③透支。', instruction: '点击右侧「跳过强化 → 透支」。本步不用选择卡牌。', selector: '.defense-step-actions .step-advance' },
      overexert.includes(mask)
        ? { title: '③ 透支：面甲已选中', text: '已选择火焰面甲，付款可用1火。现在只是预选，面甲还没有横置；点击「核对付款 →」进入④付款。', instruction: '面甲已选好。点击右侧「核对付款 →」。', selector: '.defense-step-actions .step-advance' }
        : { title: '③ 透支：选择火焰面甲', text: '当前元素池为空。点击下方「透支火焰面甲」，为这次防御提供1火；也可点击己方装备区高亮的面甲卡面。这里是透支，不要打开普通「消耗」菜单。选中后再点「核对付款 →」。', instruction: '点击右侧「透支火焰面甲」，选中后再核对付款。', cardId: mask, choiceLabel: '透支火焰面甲', selector: '.tutorial-defense-choice' },
      { title: '④ 付款：核对后一起提交', text: '核对下方组合：防御为烈焰护体、强化为无、透支为火焰面甲，付款为火1。防御威力3达到进攻威力3，点击「确认防御」。此时才会支付费用、横置卡牌并结算。', instruction: '核对烈焰护体＋火焰面甲、火1，再点击右侧「确认防御」。', selector: '[data-defense-action="submit"]' },
    ];
    return { ...hints[step], canReview, step };
  }
  function highlight(course, defense) {
    document.querySelectorAll('.tutorial-focus').forEach((node) => node.classList.remove('tutorial-focus'));
    if (!course || course.waiting || course.completed) return;
    const selectors = [];
    if (defense?.selector) selectors.push(defense.selector);
    const cardId = defense ? defense.cardId : course.card_id;
    if (cardId) selectors.push(`[data-instance-id="${CSS.escape(cardId)}"]`);
    const focus = {
      board: '.boards-wrapper', elements: '.player-panel.is-active .elem-grid',
      hand: '.hand-dock', skills: '.skill-pool-panel', hero: '.hero-row',
      opponent: '.opponent-board', end_turn: '.btn-end-turn',
      mulligan: '[data-tutorial-keep]', incoming: '.defense-incoming', pending: '.pending-action-panel', no_defend: '[data-defense-action="no_defend"]',
    };
    if (!course.card_id && focus[course.focus]) selectors.push(focus[course.focus]);
    const data = course.allowed?.data || {};
    if (course.allowed?.action === 'summon') selectors.push(`.my-board .unit-cell[data-col="${data.col}"][data-row="${data.row}"]`);
    if (course.allowed?.action === 'cast_spell') selectors.push(`.opponent-board .unit-cell[data-col="${data.target_col}"][data-row="${data.target_row}"]`);
    for (const selector of selectors) document.querySelectorAll(selector).forEach((node) => node.classList.add('tutorial-focus'));
  }
  const component = {
    props: { course: Object, busy: Boolean, defense: Object }, emits: ['continue', 'locate', 'restart', 'exit'],
    template: `<section id="tutorial-guide" class="tutorial-guide" aria-label="新手教程指引" aria-live="polite">
      <div class="tutorial-eyebrow">新手教程 · 认识你的第一场对局</div>
      <div class="tutorial-progress"><span :style="{width:(100*(course.index+1)/course.total)+'%'}"></span></div>
      <small>{{ course.chapter }} · {{ course.index+1 }} / {{ course.total }}</small>
      <h2>{{ course.title }}</h2><p>{{ course.text }}</p>
      <p class="tutorial-instruction">{{ defense ? defense.instruction : course.instruction }}</p>
      <p v-if="course.failure" class="notice-error" role="alert">{{ course.failure }}</p>
      <button v-if="course.reading && !course.failure" data-tutorial-continue class="btn btn-primary" :disabled="busy" @click="$emit('continue')">继续</button>
      <button v-else-if="!course.waiting && !course.failure" data-tutorial-locate data-inspect class="btn" @click="$emit('locate')">定位当前操作</button>
      <p v-else-if="course.waiting" class="tutorial-wait">对手正在按规则行动…</p>
      <div class="tutorial-tools"><button class="mini-btn" data-inspect @click="$emit('restart')">从头开始</button><button class="mini-btn" data-inspect @click="$emit('exit')">退出教程</button></div>
    </section>`,
  };
  const api = { permits, defenseHint, highlight, component };
  root.ArcaneTutorial = api;
  if (typeof module !== 'undefined') module.exports = api;
})(typeof window === 'undefined' ? globalThis : window);
