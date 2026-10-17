    /* ACTION_COST_PLAN_START */
    window.ArcaneActionCostPlan = (function () {
        const NEXT_SKILL_ZERO = 'next_skill_cost_zero';
        const NEXT_SKILL_MINUS = 'next_skill_use_cost_minus';
        const NEXT_ITEM_OR_SKILL_MINUS = 'next_item_or_skill_cost_minus';
        const NEXT_FIRE_CARD_MINUS = 'next_fire_card_play_cost_minus';

        function copyCost(cost) {
            const result = {};
            for (const [element, amount] of Object.entries(cost || {})) {
                const value = Number(amount || 0);
                if (value > 0) result[element] = value;
            }
            return result;
        }

        function copyModifiers(modifiers) {
            return (modifiers || []).map(modifier => ({ ...modifier }));
        }

        function modifierCanApply(modifier, card) {
            if (!modifier || Number(modifier.remaining_uses || 0) === 0) return false;
            return !modifier.target_instance_id || modifier.target_instance_id === card?.instance_id;
        }

        function applyReduction(cost, modifier) {
            const amount = Number(modifier.amount || 0) > 0 ? Number(modifier.amount) : 1;
            const elements = modifier.element ? [modifier.element] : Object.keys(cost);
            for (const element of elements) {
                const remaining = Math.max(0, Number(cost[element] || 0) - amount);
                if (remaining > 0) cost[element] = remaining;
                else delete cost[element];
            }
        }

        function consumeMatching(modifiers, card, allowedTypes) {
            for (const modifier of modifiers) {
                if (!allowedTypes.has(modifier.type) || !modifierCanApply(modifier, card)) continue;
                modifier.remaining_uses = Number(modifier.remaining_uses || 0) - 1;
            }
        }

        function skillBaseCost(card, purpose) {
            const fields = {
                attack: 'action_base_attack_cost',
                attack_boost: 'action_base_attack_boost_cost',
                defend: 'action_base_defense_cost',
                defense_boost: 'action_base_defense_boost_cost',
                reaction: 'action_base_reaction_cost',
            };
            return copyCost(card?.[fields[purpose]] || card?.elements_expense || {});
        }

        function planSkillCost(card, purpose, modifiers) {
            const cost = skillBaseCost(card, purpose);
            const isBoost = purpose === 'attack_boost' || purpose === 'defense_boost' || purpose === 'boost';
            for (const modifier of modifiers) {
                if (!modifierCanApply(modifier, card)) continue;
                if (modifier.type === NEXT_SKILL_MINUS || modifier.type === NEXT_ITEM_OR_SKILL_MINUS) {
                    applyReduction(cost, modifier);
                }
            }
            if (!isBoost && modifiers.some(modifier => modifier.type === NEXT_SKILL_ZERO && modifierCanApply(modifier, card))) {
                for (const element of Object.keys(cost)) delete cost[element];
            }
            consumeMatching(
                modifiers,
                card,
                isBoost
                    ? new Set([NEXT_SKILL_MINUS, NEXT_ITEM_OR_SKILL_MINUS])
                    : new Set([NEXT_SKILL_ZERO, NEXT_SKILL_MINUS, NEXT_ITEM_OR_SKILL_MINUS])
            );
            return cost;
        }

        function planCardPlayCost(card, modifiers) {
            const cost = copyCost(card?.action_base_play_cost || card?.elements_cost || {});
            const isConsumableOrSkill = !!card?.is_consumable || card?.type === '技能';
            for (const modifier of modifiers) {
                if (!modifierCanApply(modifier, card)) continue;
                if (modifier.type === NEXT_ITEM_OR_SKILL_MINUS && isConsumableOrSkill) applyReduction(cost, modifier);
                if (modifier.type === NEXT_FIRE_CARD_MINUS && card?.category === '火') applyReduction(cost, modifier);
            }
            const allowed = new Set();
            if (isConsumableOrSkill) allowed.add(NEXT_ITEM_OR_SKILL_MINUS);
            if (card?.category === '火') allowed.add(NEXT_FIRE_CARD_MINUS);
            consumeMatching(modifiers, card, allowed);
            return cost;
        }

        function mergeCost(total, cost) {
            for (const [element, amount] of Object.entries(cost || {})) {
                total[element] = Number(total[element] || 0) + Number(amount || 0);
            }
        }

        function planActionCost(entries, sourceModifiers) {
            const modifiers = copyModifiers(sourceModifiers);
            const total = {};
            for (const entry of entries || []) {
                if (!entry?.card) continue;
                const cost = entry.kind === 'card'
                    ? planCardPlayCost(entry.card, modifiers)
                    : planSkillCost(entry.card, entry.purpose, modifiers);
                mergeCost(total, cost);
            }
            return total;
        }

        return { planActionCost };
    }());
    /* ACTION_COST_PLAN_END */

    const { createApp, ref, computed, onMounted, onUnmounted, nextTick } = Vue;

    createApp({
        setup() {
            const connected = ref(false);
            const connectionStatus = ref('连接中...');
            const error = ref('');
            let ws = null;

            const gameState = ref(null);
            const mySlot = ref(-1);
            const mulliganSubmitted = ref(false);
            const detailCard = ref(null);
            const hoverCard = ref(null);
            const recentDiscardCard = ref(null);
            const graveyardFlash = ref({ me: null, opponent: null });
            let recentDiscardTimer = null;
            let graveyardFlashTimer = null;
            const selectedHandCard = ref(null);
            const selectedAttacker = ref(null);
            const selectedSkill = ref(null);
            const selectedItemSpellScroll = ref(null);
            const selectedSpellScrollSacrificeId = ref('');
            const selectedSpellScrollSupportId = ref('');
            const selectedSpellPowerSacrificeId = ref('');
            const defenseScrollSelected = ref([]);
            const defensePowerSacrificeId = ref('');
            const hoveredOpponentCell = ref(null);
            const skillReplaceRequest = ref(null);
            const equipmentReplaceRequest = ref(null);
            const pileView = ref(null);
            const actionLog = ref([]);
            const publicEventLog = ref([]);
            const isTestRoom = ref(false);
            const testCardNumber = ref('');
            const testElement = ref('无');
            const testElementAmount = ref(10);

            const session = new ArcaneDecisions.DecisionSession();
            const actionDraft = ref(null);
            const actionBusy = ref(false);
            const decisionError = ref('');
            const railTab = ref('hand');
            const fanPage = ref(0);
            const fanSpread = ref(28);
            const fanLayout = ref('fan');
            const fanPreviewID = ref('');
            const selectedBattleID = ref('');
            function reloadMatch() { window.location.reload(); }
            const loadDistribution = ref({});
            const allocationAmounts = ref({});
            const resourceModes = ref({});
            const legacyURL = '/game-legacy.html' + location.search;
            const actionLabels = { mulligan:'起手换牌', summon:'召唤伙伴', consume:'消耗获取元素', cast_spell:'声明法术', react_spell:'响应法术', defend:'防御', no_defend:'不防御', attack:'单位攻击', equip:'装备 / 盖反制', learn_skill:'学习技能', use_item:'使用道具', place_terrain:'放置地形', use_ability:'发动能力', resolve_action:'确认选择', end_turn:'结束回合', surrender:'投降', offer_draw:'请求和局', respond_draw_offer:'回应和局' };
            const draftTitle = computed(() => actionLabels[actionDraft.value?.action] || '确认行动');
            const decisionOwnerLabel = computed(() => {
                const s = gameState.value;
                if (!s) return '等待连接';
                if (s.phase === 'game_over') return '对局结束';
                if (s.phase === 'mulligan') return '双方起手';
                if (isSpectator) return `等待玩家${Number(s.decision_player)+1}`;
                return s.decision_player === mySlot.value ? '由你决定' : `等待玩家${Number(s.decision_player) + 1}`;
            });
            const workspaceMode = computed(() => {
                if (actionDraft.value) return 'review';
                if (incomingDrawOffer.value) return 'draw';
                if (paymentRequest.value) return 'payment';
                if (gameState.value?.phase === 'mulligan' && !mulliganSubmitted.value && !isSpectator) return 'mulligan';
                if (showPendingAction.value) return 'pending';
                if (showDefenseWindow.value) return 'defense';
                if (devourRequest.value) return 'devour';
                if (consumeGainRequest.value) return 'consume';
                if (skillReplaceRequest.value) return 'skill';
                if (equipmentReplaceRequest.value) return 'equipment';
                return 'idle';
            });
            function findVisibleCard(id) {
                const s = gameState.value;
                const cards = [...myHand.value, ...allMySkills.value, ...allMyFieldCards(), ...(pendingAction.value?.candidates || []), ...(s?.you?.skill_pool || [])];
                for (const player of [s?.you, s?.opponent]) {
                    if (!player) continue;
                    cards.push(player.hero, ...(player.skills || []), ...(player.equipment || []), ...(player.graveyard || []), ...(player.exile || []));
                    for (const col of player.units || []) cards.push(...(col || []));
                }
                const seen = new Set();
                while (cards.length) {
                    const c = cards.shift();
                    if (!c || seen.has(c)) continue;
                    if (c.instance_id === id) return c;
                    seen.add(c);
                    cards.push(...(c.bound_skills || []), ...(c.under_cards || []));
                }
                return null;
            }
            const draftSummary = computed(() => {
                const d = actionDraft.value?.data || {};
                const names = { instance_id:'卡牌', skill_id:'法术', attacker_id:'攻击者', selected:'选择', skill_ids:'防御', scroll_ids:'卷轴', boost_ids:'强化', overexert_ids:'透支', sacrifice_id:'祭献', support_id:'支援', replace_id:'替换', devour_ids:'吞噬', top_order:'牌顶顺序', bottom_order:'牌底顺序' };
                const lines = [];
                for (const [key, label] of Object.entries(names)) {
                    if (!d[key]) continue;
                    const ids = Array.isArray(d[key]) ? d[key] : [d[key]];
                    if (ids.length) lines.push(`${label}：${ids.map(id => { const c = findVisibleCard(id); return c ? `${c.name}${c.position ? ` (${c.position.col+1}-${c.position.row+1})` : ''}` : id; }).join(' → ')}`);
                }
                for (const target of [d.target, ...(d.extra_targets || []), d.position && {position:d.position}].filter(Boolean)) {
                    lines.push(`目标：${target.type === 'hero' ? '人物' : target.position ? `${Number(target.position.col)+1}-${Number(target.position.row)+1}` : target.instance_id ? findVisibleCard(target.instance_id)?.name || target.instance_id : '区域'}`);
                }
                if (d.target_col !== undefined) lines.push(`目标位置：${d.target_col+1}-${d.target_row+1}${d.target_owner !== undefined ? ` · 玩家${d.target_owner+1}` : ''}`);
                if (d.extra_target_col !== undefined) lines.push(`追加目标：${d.extra_target_col+1}-${d.extra_target_row+1}`);
                if (d.col !== undefined) lines.push(`放置位置：${d.col+1}-${d.row+1}`);
                if (d.payments) for (const p of d.payments) lines.push(`${findVisibleCard(p.instance_id)?.name || p.instance_id}：${p.amount} 点 · ${p.mode === 'life' ? '仅生命' : '先负载，后生命'}`);
                if (d.allocations) for (const [id,n] of Object.entries(d.allocations)) lines.push(`${findVisibleCard(id)?.name || id}：分配 ${n}`);
                if (d.keep !== undefined) lines.push(d.keep ? '保留当前全部手牌' : '重抽全部手牌');
                if (d.accept !== undefined) lines.push(d.accept ? '接受和局，结束对局' : '拒绝和局');
                if (d.payment) lines.push(`付款：${costLabel(d.payment)}（含所选透支；多余负载不会保留）`);
                if (d.gain) lines.push(`获得：${costLabel(d.gain)}`);
                if (d.load_distribution) for (const [id, load] of Object.entries(d.load_distribution)) lines.push(`${findVisibleCard(id)?.name || id} → ${costLabel(load)}`);
                if (actionDraft.value?.action === 'no_defend') lines.push('放弃本次防御，继续结算法术。');
                if (actionDraft.value?.action === 'resolve_action' && !d.selected?.length) lines.push('本次不选择 / 放弃发动。');
                if (actionDraft.value?.action === 'surrender') lines.push('本局立即结束，对手获胜。');
                return lines;
            });
            function setLoadDistribution(id, elem, raw) {
                loadDistribution.value = {...loadDistribution.value, [id]: {...loadDistribution.value[id], [elem]: Math.max(0, Math.floor(Number(raw) || 0))}};
            }
            function guardInteraction(event) {
                const inWorkspace = !!event.target.closest('#decision');
                if (!actionBusy.value && (!actionDraft.value && !paymentRequest.value || inWorkspace)) return;
                if (event.target.closest('.card-inspect, .interaction-detail-link, .detail-close, .close-btn, [data-inspect], [data-reconnect], a')) return;
                event.stopImmediatePropagation(); event.preventDefault();
            }
            function commitDraft() {
                const draft = actionDraft.value;
                if (!draft || !ws || ws.readyState !== WebSocket.OPEN) { decisionError.value = '连接未就绪，操作未发送'; return; }
                const id = `${Date.now()}-${nextRequestID++}`;
                if (!session.begin(id, draft.revision)) { decisionError.value = '状态已改变或动作结果未知，请同步战场后重新选择'; return; }
                actionBusy.value = true;
                decisionError.value = '';
                try { ws.send(JSON.stringify({action:draft.action, data:draft.data, request_id:id, expected_revision:draft.revision})); }
                catch (error) { session.disconnect(); actionBusy.value = false; decisionError.value = '发送中断，结果未知，请重新连接核对'; }
            }
            function cancelDraft() {
                if (actionBusy.value) return;
                actionDraft.value = null;
                decisionError.value = '';
            }
            function acceptSnapshot(state) {
                const changed = session.snapshot(state.revision);
                if (changed) {
                    actionDraft.value = null;
                    decisionError.value = session.failureMessage();
                    contextMenu.value.show = false;
                    paymentRequest.value = null;
                    consumeGainRequest.value = null;
                    devourRequest.value = null;
                    skillReplaceRequest.value = null;
                    equipmentReplaceRequest.value = null;
                    clearSelection();
                    pendingSummonDevour.value = {};
                    devourSelection.value = [];
                    resetPendingSelectionState();
                    defenseSelected.value = [];
                    defenseBoostSelected.value = [];
                    defenseScrollSelected.value = [];
                    defenseOverexertSelected.value = [];
                    defensePowerSacrificeId.value = '';
                    loadDistribution.value = {};
                    allocationAmounts.value = {};
                    resourceModes.value = {};
                }
                actionBusy.value = session.busy;
                if (state.mulligan_done) mulliganSubmitted.value = !!state.mulligan_done[mySlot.value];
            }

            const CLIENT_BUILD = 'unified-workspace-20261010';
            const IMAGE_BASE = 'https://yifeeeeei.github.io/ArcaneImages/';
            const ELEMENT_ICON_BASE = 'https://yifeeeeei.github.io/ArcaneComposer/resources/placeholders/';

            function imageFallback(event) {
                const img = event.target;
                if (img.dataset.fallback) { img.style.visibility = 'hidden'; return; }
                img.dataset.fallback = 'true';
                if (img.src.includes('yifeeeeei.github.io/ArcaneImages/')) img.src = img.src.replace('https://yifeeeeei.github.io/ArcaneImages/', 'https://raw.githubusercontent.com/Yifeeeeei/ArcaneImages/main/');
                else if (img.src.includes('yifeeeeei.github.io/ArcaneComposer/')) img.src = img.src.replace('https://yifeeeeei.github.io/ArcaneComposer/', 'https://raw.githubusercontent.com/Yifeeeeei/ArcaneComposer/main/');
                else img.style.visibility = 'hidden';
            }
            function getCardImage(card) {
                if (card?.is_hidden || card?.hidden) return 'images/card-back.svg';
                if (!card || !card.output_path) return '';
                return IMAGE_BASE + card.output_path.replace(/\\/g, '/');
            }

            const ELEM_FULL = { '火':'火焰','水':'水纹','地':'地脉','气':'大气','光':'光辉','暗':'暗影','无':'奥术','血':'生命' };
            const ELEM_ICON = { '火':'fire','水':'water','地':'earth','气':'air','光':'light','暗':'dark','无':'none' };
            const ELEMENT_ORDER = ['火','水','地','气','光','暗','无'];
            const DESCRIPTION_PLACEHOLDER_ICONS = {
                '\\光': 'ele_light.png',
                '\\暗': 'ele_dark.png',
                '\\水': 'ele_water.png',
                '\\火': 'ele_fire.png',
                '\\地': 'ele_earth.png',
                '\\气': 'ele_air.png',
                '\\无': 'ele_none.png',
                '\\攻': 'attack.png',
                '\\持': 'duration.png',
                '\\血': 'life.png',
                '\\威': 'power.png',
            };
            function elemFull(k) { return ELEM_FULL[k] || k; }
            function getElementIcon(k) { return `${ELEMENT_ICON_BASE}ele_${ELEM_ICON[k] || 'none'}.png`; }
            function escapeHTML(value) {
                return String(value ?? '').replace(/[&<>"']/g, ch => ({
                    '&': '&amp;',
                    '<': '&lt;',
                    '>': '&gt;',
                    '"': '&quot;',
                    "'": '&#39;'
                }[ch]));
            }
            function descriptionHTML(description) {
                let html = escapeHTML(description);
                for (const [placeholder, filename] of Object.entries(DESCRIPTION_PLACEHOLDER_ICONS)) {
                    const escapedPlaceholder = escapeHTML(placeholder);
                    const src = `${ELEMENT_ICON_BASE}${filename}`;
                    html = html.replaceAll(
                        escapedPlaceholder,
                        `<img src="${src}" alt="${escapeHTML(placeholder.slice(1))}" class="description-placeholder-icon inline-image">`
                    );
                }
                return html;
            }
            function costHTML(cost) {
                const entries = Object.entries(cost || {}).filter(([, v]) => Number(v) > 0);
                if (!entries.length) return '<span class="element-free">免费</span>';
                return entries.map(([k, v]) => {
                    const icon = getElementIcon(k);
                    const label = elemFull(k);
                    return `<span class="element-token" title="${label}${v}"><img src="${icon}" alt="${label}"><b>${v}</b></span>`;
                }).join('');
            }
            function hasCost(c) { return c && c.elements_cost && Object.keys(c.elements_cost).length > 0; }
            function hasUseCost(c) { return c && c.elements_expense && Object.keys(c.elements_expense).length > 0; }
            function hasGain(c) { return c && c.elements_gain && Object.keys(c.elements_gain).length > 0; }
            function playCost(card) { return card?.effective_elements_cost || card?.elements_cost || {}; }
            function learnCost(card) { return card?.effective_learn_cost || card?.elements_cost || {}; }
            function entryCostLabel(card) {
                if (!card) return '费用';
                if (card.type === '技能') return '入场费用';
                if (card.type === '伙伴') return '入场费用';
                return '费用';
            }
            function costLabel(cost) {
                const entries = Object.entries(cost || {}).filter(([, v]) => Number(v) > 0);
                if (!entries.length) return '免费';
                return entries.map(([k, v]) => `${elemFull(k)}${v}`).join(' ');
            }

            const params = new URLSearchParams(window.location.search);
            const roomID = params.get('room');
            const playerID = params.get('player_id');
            const playerName = params.get('player_name');
            const deckCode = params.get('deck_code');
            const isSpectator = params.get('role') === 'spectator';

            if (!roomID || !playerID || (!isSpectator && !deckCode)) {
                error.value = '缺少参数，请从大厅进入';
                connectionStatus.value = '连接失败';
            }

            const myState = computed(() => gameState.value ? (gameState.value.you || {}) : {});
            const opponent = computed(() => gameState.value ? (gameState.value.opponent || {}) : {});
            const myHand = computed(() => myState.value.hand || []);
            const fanPageCount = computed(() => Math.max(1, Math.ceil(myHand.value.length / 7)));
            const currentFanPage = computed(() => Math.min(fanPage.value, fanPageCount.value - 1));
            const fanCards = computed(() => myHand.value.slice(currentFanPage.value * 7, currentFanPage.value * 7 + 7));
            const fanPreviewCard = computed(() => myHand.value.find(c => c.instance_id === fanPreviewID.value) || null);
            const selectedBattleCard = computed(() => allMyFieldCards().find(c => c.instance_id === selectedBattleID.value) || null);
            const showFanPreview = computed(() => workspaceMode.value === 'idle' && !!fanPreviewCard.value && !hasActiveTargeting.value && !selectedHandCard.value && !selectedSkill.value && !selectedAttacker.value);
            function previewHand(card) { fanPreviewID.value = card.instance_id; previewCard(card); }
            function changeFanPage(delta) { fanPage.value = Math.max(0, Math.min(fanPageCount.value - 1, currentFanPage.value + delta)); }
            function activateFanCard(card) { fanPreviewID.value = ''; selectedBattleID.value = ''; selectHandCard(card); }
            async function navigateFan(event, card) {
                if (event.key === 'Escape') { fanPreviewID.value = ''; return; }
                if (!['ArrowLeft','ArrowRight','Home','End'].includes(event.key)) return;
                event.preventDefault();
                const index = myHand.value.findIndex(c => c.instance_id === card.instance_id);
                const target = event.key === 'Home' ? 0 : event.key === 'End' ? myHand.value.length - 1 : Math.max(0, Math.min(myHand.value.length - 1, index + (event.key === 'ArrowRight' ? 1 : -1)));
                fanPage.value = Math.floor(target / 7);
                const id = myHand.value[target]?.instance_id;
                await nextTick();
                if (id) document.querySelector('.hand-card[data-instance-id="' + CSS.escape(id) + '"]')?.focus();
            }

            const opponentHandCards = computed(() => {
                const revealed = opponent.value.revealed_hand || [];
                const hiddenCount = Math.max(0, Number(opponent.value.hand_count || 0) - revealed.length);
                const hidden = Array.from({ length: hiddenCount }, (_, i) => ({
                    key: `hidden-${i}`,
                    name: '未知手牌',
                    hidden: true,
                }));
                return [
                    ...hidden,
                    ...revealed.map(card => ({ ...card, key: card.instance_id, hidden: false })),
                ];
            });
            function slotCount(state, field) {
                return Math.max(Number(state?.[field] || 5), 5);
            }
            function paddedSlots(cards, count) {
                const result = Array.isArray(cards) ? [...cards] : [];
                while (result.length < count) result.push(null);
                return result;
            }
            const mySkillSlotCount = computed(() => slotCount(myState.value, 'skill_slot_capacity'));
            const myEquipmentSlotCount = computed(() => slotCount(myState.value, 'equipment_slot_capacity'));
            const opponentSkillSlotCount = computed(() => slotCount(opponent.value, 'skill_slot_capacity'));
            const opponentEquipmentSlotCount = computed(() => slotCount(opponent.value, 'equipment_slot_capacity'));
            const mySkills = computed(() => paddedSlots(myState.value.skills, mySkillSlotCount.value));
            const myEquipment = computed(() => paddedSlots(myState.value.equipment, myEquipmentSlotCount.value));
            const opponentSkills = computed(() => paddedSlots(opponent.value.skills, opponentSkillSlotCount.value));
            const allMyBoundSkills = computed(() => {
                const skills = [];
                for (const unit of allMyFieldCards()) {
                    for (const bound of unit.bound_skills || []) {
                        skills.push(bound);
                    }
                }
                return skills;
            });
            const allMySkills = computed(() => [...mySkills.value.filter(Boolean), ...allMyBoundSkills.value]);
            const replaceableSkills = computed(() => mySkills.value.filter(skill => skill && !skill.is_horizontal));
            const replaceableEquipment = computed(() => {
                const card = equipmentReplaceRequest.value?.card || selectedHandCard.value;
                return replaceableEquipmentForCard(card);
            });
            const isMyTurn = computed(() => !isSpectator && gameState.value && Number(gameState.value.current_turn) === currentPlayerSlot());
            const drawOfferBy = computed(() => {
                const value = Number(gameState.value?.draw_offer_by ?? -1);
                return Number.isFinite(value) ? value : -1;
            });
            const incomingDrawOffer = computed(() => !isSpectator && gameState.value?.phase !== 'game_over' && drawOfferBy.value === 1 - currentPlayerSlot());
            const outgoingDrawOffer = computed(() => !isSpectator && gameState.value?.phase !== 'game_over' && drawOfferBy.value === currentPlayerSlot());
            const canUseMatchActions = computed(() => !isSpectator && gameState.value && gameState.value.phase !== 'game_over');

            function getMyTerrain(col, row) { return myState.value.terrain?.[col]?.[row] || null; }
            function getOpponentTerrain(col, row) { return opponent.value.terrain?.[col]?.[row] || null; }

            const pendingSpellView = computed(() => normalizePendingSpell(gameState.value?.pending_spell));
            const showDefenseWindow = computed(() => {
                const spell = pendingSpellView.value;
                return !isSpectator &&
                    gameState.value?.phase === 'defense_window' &&
                    !!spell &&
                    spell.defender_id === currentPlayerSlot();
            });

            function currentPlayerSlot() {
                const slot = Number(mySlot.value);
                return Number.isFinite(slot) ? slot : -1;
            }

            function normalizePendingSpell(spell) {
                if (!spell) return null;
                const attackerRaw = spell.attacker_id ?? spell.attacker ?? spell.cast_player;
                const attackerID = Number(attackerRaw);
                if (!Number.isFinite(attackerID)) return null;
                const powerRaw = spell.power ?? spell.total_power ?? spell.TotalPower ?? 0;
                const power = Number(powerRaw);
                return {
                    raw: spell,
                    attacker_id: attackerID,
                    defender_id: 1 - attackerID,
                    skill: spell.skill || null,
                    target: spell.target || null,
                    extra_targets: Array.isArray(spell.extra_targets) ? spell.extra_targets : [],
                    affected_positions: Array.isArray(spell.affected_positions) ? spell.affected_positions : [],
                    boost_skills: Array.isArray(spell.boost_skills) ? spell.boost_skills : [],
                    power_sources: Array.isArray(spell.power_sources) ? spell.power_sources : [],
                    power: Number.isFinite(power) ? power : 0,
                };
            }

            function pendingSpellAttackerID() {
                const spell = pendingSpellView.value;
                if (!spell) return -1;
                return spell.attacker_id;
            }

            const activePreviewCard = computed(() =>
                recentDiscardCard.value ||
                hoverCard.value ||
                selectedHandCard.value ||
                selectedSkill.value ||
                selectedAttacker.value ||
                myState.value.hero ||
                null
            );

            const myFrontRow = computed(() => {
                if (!myState.value.units) return 0;
                for (let r = 0; r < 3; r++)
                    for (let c = 0; c < 3; c++)
                        if (myState.value.units[c] && myState.value.units[c][r]) return r;
                return -1;
            });
            const opponentFrontRow = computed(() => {
                if (!opponent.value.units) return 0;
                for (let r = 0; r < 3; r++)
                    for (let c = 0; c < 3; c++)
                        if (opponent.value.units[c] && opponent.value.units[c][r]) return r;
                return -1;
            });

            function getMyUnit(col, row) { return myState.value.units?.[col]?.[row] || null; }
            function getOpponentUnit(col, row) { return opponent.value.units?.[col]?.[row] || null; }
            function myBoardOwnerID() {
                return isSpectator ? 0 : currentPlayerSlot();
            }
            function opponentBoardOwnerID() {
                const ownerID = myBoardOwnerID();
                return ownerID >= 0 ? 1 - ownerID : -1;
            }
            function spellTargetPosition(target) {
                if (!target || target.type !== 'unit') return null;
                const pos = target.position || target.Position || {};
                const col = Number(pos.col ?? pos.Col);
                const row = Number(pos.row ?? pos.Row);
                return Number.isInteger(col) && Number.isInteger(row) ? { col, row } : null;
            }
            function spellTargetOwnerID(target, spell = pendingSpellView.value) {
                if (!spell) return -1;
                const ownerID = Number(target?.owner_id ?? target?.OwnerID);
                return Number.isInteger(ownerID) ? ownerID : spell.defender_id;
            }
            function isPendingSpellCenter(ownerID, col, row) {
                const spell = pendingSpellView.value;
                if (!spell || ownerID < 0) return false;
                return [spell.target, ...spell.extra_targets].some(target => {
                    const pos = spellTargetPosition(target);
                    return pos && spellTargetOwnerID(target, spell) === ownerID && pos.col === col && pos.row === row;
                });
            }
            function isPendingSpellAffected(ownerID, col, row) {
                const spell = pendingSpellView.value;
                if (!spell || ownerID < 0) return false;
                return spell.affected_positions.some(pos =>
                    Number(pos?.owner_id) === ownerID && Number(pos?.col) === col && Number(pos?.row) === row
                );
            }
            function pendingSpellTargetOwnerState() {
                const spell = pendingSpellView.value;
                if (!spell) return myState.value;
                return spell.attacker_id === currentPlayerSlot() ? opponent.value : myState.value;
            }
            function cardAtSpellTarget(target) {
                if (!target || target.type !== 'unit') return null;
                const pos = target.position || target.Position || {};
                const col = Number(pos.col ?? pos.Col);
                const row = Number(pos.row ?? pos.Row);
                if (!Number.isFinite(col) || !Number.isFinite(row)) return null;
                return pendingSpellTargetOwnerState().units?.[col]?.[row] || null;
            }
            function spellTargetLabel(target) {
                if (!target || target.type === 'none') return '无目标';
                if (target.type !== 'unit') return target.type || '未知目标';
                const pos = target.position || target.Position || {};
                const col = Number(pos.col ?? pos.Col);
                const row = Number(pos.row ?? pos.Row);
                const card = cardAtSpellTarget(target);
                const coord = Number.isFinite(col) && Number.isFinite(row) ? `${col + 1}-${row + 1}` : '未知位置';
                return card ? `${card.name}（${coord}）` : `单位（${coord}）`;
            }
            function pendingSpellTargetLabel() {
                return spellTargetLabel(pendingSpellView.value?.target);
            }
            function pendingSpellExtraTargetLabels() {
                return (pendingSpellView.value?.extra_targets || []).map(spellTargetLabel);
            }
            function pendingSpellSources() {
                const spell = pendingSpellView.value;
                if (!spell) return [];
                const fullCards = [spell.skill, ...spell.boost_skills].filter(Boolean);
                if (spell.power_sources.length) {
                    return spell.power_sources.map((source, idx) => {
                        const card = fullCards.find(c => c.instance_id === source.instance_id) || null;
                        return {
                            key: `${source.instance_id || 'source'}-${idx}`,
                            name: source.card_name || card?.name || '未知来源',
                            power: Number(source.power || 0),
                            is_main: !!source.is_main,
                            card,
                        };
                    });
                }
                return fullCards.map((card, idx) => ({
                    key: `${card.instance_id || 'card'}-${idx}`,
                    name: card.name,
                    power: Math.max(Number(card.power || 0), 0),
                    is_main: idx === 0,
                    card,
                }));
            }
            function allMyFieldCards() {
                const cards = [];
                if (myState.value.hero) cards.push(myState.value.hero);
                for (let col = 0; col < 3; col++) {
                    for (let row = 0; row < 3; row++) {
                        const unit = getMyUnit(col, row);
                        if (unit && unit !== myState.value.hero) cards.push(unit);
                    }
                }
                for (const equipment of myEquipment.value) {
                    if (equipment) cards.push(equipment);
                }
                return cards;
            }
            function hasLightWildcardPayment() {
                return !!myState.value.light_wildcard;
            }
            function showCardDetail(card) { if (!card?.is_hidden && !card?.hidden) { previewCard(card); detailCard.value = card; } }
            function previewCard(card) {
                if (card && !card.is_hidden && !card.hidden) {
                    recentDiscardCard.value = null;
                    hoverCard.value = card;
                }
            }

            function cloneElements(elements) {
                const result = {};
                for (const elem of ELEMENT_ORDER) result[elem] = Number(elements?.[elem] || 0);
                return result;
            }

            function calculatePayment(cost) {
                return calculatePaymentWithAvailable(myState.value.elements || {}, cost);
            }

            function calculatePaymentWithAvailable(availableSource, cost, lightWildcard = hasLightWildcardPayment()) {
                const remaining = {};
                for (const [elem, raw] of Object.entries(cost || {})) {
                    const value = Number(raw) || 0;
                    if (value > 0) remaining[elem] = (remaining[elem] || 0) + value;
                }
                if (!Object.keys(remaining).length) return {};

                const available = cloneElements(availableSource || {});
                const payment = {};
                const pay = (elem, amount) => {
                    if (amount <= 0) return;
                    payment[elem] = (payment[elem] || 0) + amount;
                    available[elem] = (available[elem] || 0) - amount;
                };

                for (const elem of ELEMENT_ORDER) {
                    if (elem === '无') continue;
                    const need = remaining[elem] || 0;
                    const amount = Math.min(available[elem] || 0, need);
                    pay(elem, amount);
                    remaining[elem] = need - amount;
                }

                for (const elem of ELEMENT_ORDER) {
                    if (elem === '无') continue;
                    const need = remaining[elem] || 0;
                    const amount = Math.min(available['无'] || 0, need);
                    pay('无', amount);
                    remaining[elem] = need - amount;
                    if (remaining[elem] > 0 && lightWildcard && elem !== '光') {
                        const lightAmount = Math.min(available['光'] || 0, remaining[elem]);
                        pay('光', lightAmount);
                        remaining[elem] -= lightAmount;
                    }
                    if (remaining[elem] > 0) return null;
                }

                let arcaneNeed = remaining['无'] || 0;
                const arcanePay = Math.min(available['无'] || 0, arcaneNeed);
                pay('无', arcanePay);
                arcaneNeed -= arcanePay;
                for (const elem of ELEMENT_ORDER) {
                    if (elem === '无' || arcaneNeed <= 0) continue;
                    const amount = Math.min(available[elem] || 0, arcaneNeed);
                    pay(elem, amount);
                    arcaneNeed -= amount;
                }
                return arcaneNeed > 0 ? null : payment;
            }

            function totalCost(cost) {
                return Object.values(cost || {}).reduce((sum, value) => sum + Number(value || 0), 0);
            }

            function isFireCompanion(card) {
                return !!card && card.type === '伙伴' && card.category === '火';
            }

            function cardLoadValue(card) {
                return totalCost(card?.elements_gain || {});
            }

            function myUnitCards() {
                const units = [];
                if (!myState.value?.units) return units;
                for (let col = 0; col < 3; col++) {
                    for (let row = 0; row < 3; row++) {
                        const unit = getMyUnit(col, row);
                        if (unit) units.push(unit);
                    }
                }
                return units;
            }

            function isExactPaymentForCost(payment, cost) {
                const spent = calculatePaymentWithAvailable(payment, cost);
                if (!spent) return false;
                for (const elem of ELEMENT_ORDER) {
                    if (Number(spent[elem] || 0) !== Number(payment[elem] || 0)) return false;
                }
                return true;
            }

            function canPayCost(cost) {
                return calculatePayment(cost) !== null;
            }

            function missingCostParts(cost) {
                if (!cost || !Object.keys(cost).length || canPayCost(cost)) return [];
                const remaining = {};
                for (const [elem, raw] of Object.entries(cost || {})) {
                    const value = Number(raw) || 0;
                    if (value > 0) remaining[elem] = (remaining[elem] || 0) + value;
                }
                const available = cloneElements(myState.value.elements || {});
                const spend = (elem, amount) => {
                    if (amount <= 0) return;
                    available[elem] = Math.max(0, Number(available[elem] || 0) - amount);
                };

                for (const elem of ELEMENT_ORDER) {
                    if (elem === '无') continue;
                    const amount = Math.min(Number(available[elem] || 0), Number(remaining[elem] || 0));
                    spend(elem, amount);
                    remaining[elem] = Number(remaining[elem] || 0) - amount;
                }
                for (const elem of ELEMENT_ORDER) {
                    if (elem === '无') continue;
                    const amount = Math.min(Number(available['无'] || 0), Number(remaining[elem] || 0));
                    spend('无', amount);
                    remaining[elem] = Number(remaining[elem] || 0) - amount;
                }

                let arcaneNeed = Number(remaining['无'] || 0);
                const arcanePay = Math.min(Number(available['无'] || 0), arcaneNeed);
                spend('无', arcanePay);
                arcaneNeed -= arcanePay;
                for (const elem of ELEMENT_ORDER) {
                    if (elem === '无' || arcaneNeed <= 0) continue;
                    const amount = Math.min(Number(available[elem] || 0), arcaneNeed);
                    spend(elem, amount);
                    arcaneNeed -= amount;
                }
                remaining['无'] = arcaneNeed;

                return ELEMENT_ORDER
                    .map(elem => [elem, Math.max(0, Number(remaining[elem] || 0))])
                    .filter(([, amount]) => amount > 0)
                    .map(([elem, amount]) => `${elemFull(elem)}${amount}`);
            }

            function missingCostText(cost) {
                if (!cost || !Object.keys(cost).length) return '免费';
                if (canPayCost(cost)) return '可支付';
                const missing = missingCostParts(cost);
                return missing.length ? `缺${missing.join(' ')}` : '元素不足';
            }

            function cardAvailabilityText(card) {
                if (!card || !hasCost(card)) return '';
                const cost = playCost(card);
                return canPayCost(cost) ? '可用' : missingCostText(cost);
            }

            function hasNextSkillCostZero() {
                return (myState.value.temp_modifiers || []).some(m => m.type === 'next_skill_cost_zero' && (m.remaining_uses || 0) > 0);
            }

            function skillUseCost(skill, applyNextSkillFree = true) {
                if (!skill) return {};
                if (skill.type !== '技能') return {};
                if (applyNextSkillFree && hasNextSkillCostZero()) {
                    return {};
                }
                return { ...(skill.action_base_attack_cost || skill.elements_expense || {}) };
            }

            function clearSelection() {
                selectedBattleID.value = '';
                fanPreviewID.value = '';
                selectedHandCard.value = null;
                selectedAttacker.value = null;
                selectedSkill.value = null;
                selectedItemSpellScroll.value = null;
                selectedSpellScrollSacrificeId.value = '';
                selectedSpellScrollSupportId.value = '';
                selectedSpellPowerSacrificeId.value = '';
                selectedBoostSkills.value = [];
                pendingExtraTargetCast.value = null;
                hoveredOpponentCell.value = null;
            }

            function statusEntries(card, limit = 3) {
                if (!card) return [];
                const entries = [];
                if (card.is_horizontal) entries.push({ label: '横置' });
                const statusLabels = {
                    winter_bow_water_mark: '水纹',
                };
                for (const [name, count] of Object.entries(card.statuses || {})) {
                    if (Number(count) > 0) {
                        const label = statusLabels[name] || name;
                        entries.push({ label: `${label}${Number(count) > 1 ? count : ''}` });
                    }
                }
                return limit > 0 ? entries.slice(0, limit) : entries;
            }

            function effectLabels(card) {
                if (!card) return [];
                return [...new Set([
                    ...(card.effect_categories || []),
                    ...(card.effect_optionality || []),
                ].filter(Boolean))];
            }

            function groupedCards(cards) {
                const groups = new Map();
                for (const raw of cards || []) {
                    if (!raw) continue;
                    const key = raw.number || raw.instance_id || raw.name;
                    const current = groups.get(key);
                    if (current) {
                        current.count += Number(raw.count || 1);
                    } else {
                        groups.set(key, { ...raw, count: Number(raw.count || 1) });
                    }
                }
                return [...groups.values()].sort((a, b) => String(a.number || a.name).localeCompare(String(b.number || b.name)));
            }

            function openPileView(owner, zone) {
                if (isSpectator && zone === 'deck') return;
                if (owner !== 'me' && zone === 'deck') return;
                const player = owner === 'me' ? myState.value : opponent.value;
                const ownerName = owner === 'me'
                    ? (myState.value.player_name || (isSpectator ? '玩家1' : '你'))
                    : (opponent.value.player_name || (isSpectator ? '玩家2' : '对手'));
                const cards = zone === 'deck'
                    ? groupedCards(player.deck_summary || [])
                    : groupedCards(player.graveyard || []);
                pileView.value = {
                    zone,
                    ownerName,
                    title: zone === 'deck' ? `剩余卡组 (${player.deck_count || 0})` : `弃牌堆 (${player.graveyard?.length || 0})`,
                    cards,
                };
            }

            function closePileView() {
                pileView.value = null;
            }

            const paymentRequest = ref(null);
            const paymentSelection = ref({});
            const consumeGainRequest = ref(null);
            const consumeGainSelection = ref({});
            const devourRequest = ref(null);
            const devourSelection = ref([]);
            const pendingSummonDevour = ref({});

            function hasArcanePaymentChoice(cost) {
                if (!cost || Number(cost['无'] || 0) <= 0) return false;
                return ELEMENT_ORDER.filter(elem => Number(myState.value.elements?.[elem] || 0) > 0).length > 1;
            }

            function hasLightWildcardPaymentChoice(cost) {
                if (!cost || !hasLightWildcardPayment() || Number(myState.value.elements?.['光'] || 0) <= 0) return false;
                if (calculatePaymentWithAvailable(myState.value.elements || {}, cost, false) === null) return false;
                for (const elem of ELEMENT_ORDER) {
                    if (elem === '光' || elem === '无' || Number(cost[elem] || 0) <= 0) continue;
                    const remainingCost = { ...(cost || {}) };
                    remainingCost[elem] = Number(remainingCost[elem] || 0) - 1;
                    const remainingAvailable = cloneElements(myState.value.elements || {});
                    remainingAvailable['光'] = Number(remainingAvailable['光'] || 0) - 1;
                    if (calculatePaymentWithAvailable(remainingAvailable, remainingCost, true) !== null) return true;
                }
                return false;
            }

            function needsPaymentChoice(cost) {
                return hasArcanePaymentChoice(cost) || hasLightWildcardPaymentChoice(cost);
            }

            function sendPaidAction(action, data, cost, title, afterSend, available = null) {
                if (totalCost(cost) === 0) return sendAction(action, data, { onSuccess: afterSend });
                paymentSelection.value = {};
                paymentRequest.value = { action, data, cost: {...cost}, title, afterSend, available };
            }
            function paymentAvailable() { return paymentRequest.value?.available || myState.value.elements || {}; }
            function removePayment(elem) { paymentSelection.value = { ...paymentSelection.value, [elem]: Math.max(0, (paymentSelection.value[elem] || 0) - 1) }; }

            function paymentElements() {
                return ELEMENT_ORDER.filter(elem => Number(paymentAvailable()[elem] || 0) > 0);
            }

            function selectedPaymentTotal() {
                return Object.values(paymentSelection.value).reduce((sum, value) => sum + Number(value || 0), 0);
            }

            function canAddPayment(elem) {
                if (!paymentRequest.value) return false;
                if (selectedPaymentTotal() >= totalCost(paymentRequest.value.cost)) return false;
                return Number(paymentSelection.value[elem] || 0) < Number(paymentAvailable()[elem] || 0);
            }

            function addPayment(elem) {
                if (!canAddPayment(elem)) return;
                paymentSelection.value = {
                    ...paymentSelection.value,
                    [elem]: Number(paymentSelection.value[elem] || 0) + 1,
                };
            }

            function isPaymentSelectionValid() {
                if (!paymentRequest.value) return false;
                if (selectedPaymentTotal() !== totalCost(paymentRequest.value.cost)) return false;
                return isExactPaymentForCost(paymentSelection.value, paymentRequest.value.cost);
            }

            function paymentRemainingText() {
                if (!paymentRequest.value) return '';
                if (isPaymentSelectionValid()) return '支付组合有效';
                const need = totalCost(paymentRequest.value.cost) - selectedPaymentTotal();
                if (need > 0) return `还需选择 ${need} 个元素`;
                return '当前组合不能支付这项费用';
            }

            function confirmPayment() {
                if (!isPaymentSelectionValid()) return;
                const req = paymentRequest.value;
                sendAction(req.action, { ...req.data, payment: { ...paymentSelection.value } }, {
                    onSuccess: () => {
                        if (req.afterSend) req.afterSend();
                        paymentRequest.value = null;
                        paymentSelection.value = {};
                    },
                }, true);
            }

            function cancelPayment() {
                paymentRequest.value = null;
                paymentSelection.value = {};
            }

            function positiveLoadEntries(card) {
                return Object.entries(card?.elements_gain || {})
                    .map(([elem, amount]) => [elem, Number(amount) || 0])
                    .filter(([, amount]) => amount > 0);
            }

            function totalPositiveLoad(card) {
                return positiveLoadEntries(card).reduce((sum, [, amount]) => sum + amount, 0);
            }

            function isFirstTurnHeroConsumeChoice(unit) {
                if (!unit || !gameState.value || !myState.value?.hero) return false;
                if (!gameState.value.is_first_turn) return false;
                if (gameState.value.current_turn !== gameState.value.first_player) return false;
                if (mySlot.value !== gameState.value.first_player) return false;
                if (unit.instance_id !== myState.value.hero.instance_id) return false;
                const entries = positiveLoadEntries(unit);
                return entries.length > 1 && totalPositiveLoad(unit) > 1;
            }

            function openConsumeGainChoice(unit) {
                const available = {};
                positiveLoadEntries(unit).forEach(([elem, amount]) => {
                    available[elem] = amount;
                });
                consumeGainSelection.value = {};
                consumeGainRequest.value = {
                    card: unit,
                    available,
                    limit: Math.ceil(totalPositiveLoad(unit) / 2),
                };
            }

            function consumeGainElements() {
                return ELEMENT_ORDER.filter(elem => Number(consumeGainRequest.value?.available?.[elem] || 0) > 0);
            }

            function selectedConsumeGainTotal() {
                return Object.values(consumeGainSelection.value || {}).reduce((sum, amount) => sum + (Number(amount) || 0), 0);
            }

            function canAddConsumeGain(elem) {
                if (!consumeGainRequest.value) return false;
                if (selectedConsumeGainTotal() >= consumeGainRequest.value.limit) return false;
                return Number(consumeGainSelection.value[elem] || 0) < Number(consumeGainRequest.value.available?.[elem] || 0);
            }

            function addConsumeGain(elem) {
                if (!canAddConsumeGain(elem)) return;
                consumeGainSelection.value = {
                    ...consumeGainSelection.value,
                    [elem]: Number(consumeGainSelection.value[elem] || 0) + 1,
                };
            }

            function isConsumeGainValid() {
                const req = consumeGainRequest.value;
                if (!req) return false;
                if (selectedConsumeGainTotal() !== req.limit) return false;
                return Object.entries(consumeGainSelection.value || {}).every(([elem, amount]) => {
                    const selected = Number(amount) || 0;
                    return selected >= 0 && selected <= Number(req.available?.[elem] || 0);
                });
            }

            function consumeGainRemainingText() {
                if (!consumeGainRequest.value) return '';
                if (isConsumeGainValid()) return '获得元素组合有效';
                const need = consumeGainRequest.value.limit - selectedConsumeGainTotal();
                if (need > 0) return `还需选择 ${need} 个元素`;
                return '当前组合超过可获得负载';
            }

            function confirmConsumeGain() {
                if (!isConsumeGainValid()) return;
                const req = consumeGainRequest.value;
                sendAction('consume', {
                    instance_id: req.card.instance_id,
                    gain: { ...consumeGainSelection.value },
                });
                addLog(`消耗 ${req.card.name}`);
                consumeGainRequest.value = null;
                consumeGainSelection.value = {};
            }

            function cancelConsumeGain() {
                consumeGainRequest.value = null;
                consumeGainSelection.value = {};
            }

            // Context menu for cards
            const contextMenu = ref({ show: false, x: 0, y: 0, card: null, source: '' });

            const menuItems = computed(() => {
                const c = contextMenu.value.card, source = contextMenu.value.source;
                if (!c || c.is_hidden || c.hidden) return [];
                const items = [{id:'detail', label:'放大查看'}];
                if (isSpectator || actionBusy.value || actionDraft.value || paymentRequest.value || !isMyTurn.value || gameState.value?.phase !== 'main') return items;
                if (source === 'hand') { if (canPlayHandCard(c)) items.push({id:'use',label:c.type === '伙伴' ? '选择召唤位置' : '使用 / 装备'}); }
                if (source === 'skill_pool' && canLearnSkill(c)) items.push({id:'learn',label:'学习技能'});
                if (source === 'field' || source === 'skill') {
                    if (source === 'field' && canConsumeUnit(c)) items.push({id:'consume',label:'消耗'});
                    if (source === 'field' && canAttackWith(c)) items.push({id:'attack',label:'攻击'});
                    if (canCastSkill(c)) items.push({id:'cast',label:'施放'});
                    if (canUsePerTurn(c)) items.push({id:'per_turn',label:abilityLabel(c,'per_turn')});
                    if (canUseUltimate(c)) items.push({id:'ultimate',label:'绝技'});
                }
                return items;
            });
            let menuAnchor = null;
            async function showCardContextMenu(event, card, source) {
                if (!card || card.is_hidden || card.hidden || hasActiveTargeting.value && source !== 'interaction') return;
                previewCard(card);
                menuAnchor = event?.currentTarget || document.activeElement;
                const rect = menuAnchor?.getBoundingClientRect?.();
                const x = event?.clientX || rect?.right || 20, y = event?.clientY || rect?.top || 20;
                contextMenu.value = {show:true, x:Math.min(x+8,innerWidth-230), y:Math.max(8,y-180), card, source};
                await nextTick();
                const menu = document.querySelector('.context-menu');
                if (!menu) return;
                const bounds=menu.getBoundingClientRect();
                contextMenu.value.x=Math.max(8,Math.min(x+8,innerWidth-bounds.width-8));
                contextMenu.value.y=Math.max(8,Math.min(y-bounds.height-8,innerHeight-bounds.height-8));
                menu.querySelector('button')?.focus({preventScroll:true});
            }
            function menuKeydown(event) {
                if (event.key === 'Escape') { event.preventDefault(); hideContextMenu(); menuAnchor?.focus?.({preventScroll:true}); return; }
                if (!['ArrowDown','ArrowUp','Home','End','Tab'].includes(event.key)) return;
                event.preventDefault();
                const buttons=[...event.currentTarget.querySelectorAll('button')], index=buttons.indexOf(document.activeElement);
                const n=event.key==='Home'?0:event.key==='End'?buttons.length-1:(index+(event.key==='ArrowUp'||event.shiftKey?-1:1)+buttons.length)%buttons.length;
                buttons[n]?.focus();
            }
            function previewMetadata(card) {
                if (!card || card.is_hidden || card.hidden) return [];
                const labels=statusEntries(card,0).map(s=>s.label);
                if (card.position && card.position.row === (card.owner === mySlot.value ? myFrontRow.value : opponentFrontRow.value)) labels.push('前排');
                if (card.is_defense_only) labels.push('防御技能');
                if (card.is_sorcery) labels.push('咒术');
                if (card.is_set_counter) labels.push('己方盖放牌');
                return labels;
            }

            function showInteractionCardContextMenu(event) {
                const cardEl = event.target.closest('[data-instance-id]');
                if (!cardEl || !event.currentTarget.contains(cardEl)) return;
                const card = defenseCardById(cardEl.dataset.instanceId || '');
                if (!card) return;
                showCardContextMenu(event, card, 'interaction');
            }

            function hideContextMenu() {
                contextMenu.value.show = false;
            }

            function contextMenuAction(action) {
                const card = contextMenu.value.card;
                const source = contextMenu.value.source;
                if (!card) return;

                switch (action) {
                    case 'consume': consumeCard(card); break;
                    case 'attack': selectAttacker(card); break;
                    case 'cast': selectSkillToCast(card); break;
                    case 'per_turn': useAbility(card, 'per_turn'); break;
                    case 'ultimate': useAbility(card, 'ultimate'); break;
                    case 'detail':
                        showCardDetail(card);
                        break;
                    case 'learn':
                        if (source === 'skill_pool') {
                            learnSkill(card);
                        }
                        break;
                    case 'use':
                        if (source === 'hand') {
                            if (isSpellScrollItem(card) && card.needs_target) {
                                selectedHandCard.value = null;
                                selectedItemSpellScroll.value = card;
                                selectedSpellScrollSacrificeId.value = '';
                                selectedSpellScrollSupportId.value = '';
                                const prefix = card.preparation === 'fire_sacrifice'
                                    ? `选择 ${card.name}，先点击友方火焰伙伴祭献`
                                    : card.preparation === 'companion_support'
                                        ? `选择 ${card.name}，先点击友方伙伴支援`
                                        : `选择 ${card.name}，点击法力范围内目标`;
                                addLog(prefix);
                            } else if (isConsumable(card)) {
                                selectedHandCard.value = null;
                                if (card.is_counter_trap) {
                                    setCounterCard(card);
                                } else {
                                    sendPaidAction('use_item', { instance_id: card.instance_id }, playCost(card), `使用 ${card.name}`, () => addLog(`使用 ${card.name}`));
                                }
                            } else if (card.type === '道具' && !card.is_terrain) {
                                selectedHandCard.value = null;
                                equipCard(card);
                            } else {
                                selectHandCard(card);
                            }
                        }
                        break;
                }
                hideContextMenu();
            }

            let nextRequestID = 1;

            // Explicit simple controls finalize immediately; target-based and destructive actions retain one review.
            // All writes still pass through the same revision/in-flight gate in commitDraft.
            function sendAction(action, data = {}, callbacks = null, finalize = ['consume','mulligan','resolve_action','no_defend','end_turn','use_ability','learn_skill'].includes(action) && !(action === 'use_ability' && data.ability_type === 'ultimate')) {
                if (isSpectator || actionBusy.value || actionDraft.value) return false;
                actionDraft.value = { action, data: JSON.parse(JSON.stringify(data)), callbacks, revision: gameState.value?.revision };
                decisionError.value = '';
                if (finalize) commitDraft();
                return true;
            }

            function sendClientLog(level, message, data = {}) {
                if (!ws || ws.readyState !== WebSocket.OPEN) return;
                ws.send(JSON.stringify({
                    action: 'client_log',
                    data: {
                        level,
                        message,
                        client_build: CLIENT_BUILD,
                        url: location.href,
                        user_agent: navigator.userAgent,
                        ...data
                    }
                }));
            }

            function mulligan(keep) {
                sendAction('mulligan', { keep });

            }

            function surrender() {
                if (!canUseMatchActions.value) return;

                sendAction('surrender', {});
                addLog('你选择投降');
            }

            function offerDraw() {
                if (!canUseMatchActions.value || outgoingDrawOffer.value) return;
                sendAction('offer_draw', {});
                addLog('已向对手请求和局');
            }

            function respondDrawOffer(accept) {
                if (!incomingDrawOffer.value) return;
                sendAction('respond_draw_offer', { accept: !!accept });
                addLog(accept ? '接受和局' : '拒绝和局');
            }

            function isConsumable(card) {
                if (!card || card.type !== '道具') return false;
                return !!card.is_consumable;
            }
            function isSpellScrollItem(card) {
                return !!card && card.type === '道具' && !!card.is_consumable &&
                    !!card.is_spell_scroll;
            }
            function isEquippableItem(card) {
                return !!card && card.type === '道具' && !card.is_terrain && !isConsumable(card);
            }
            function hasEmptyEquipmentSlot() {
                return myEquipment.value.some(equipment => !equipment);
            }
            function equipmentSubtype(card) { return card?.equipment_subtype || ''; }

            function replaceableEquipmentForCard(card) {
                if (!card) return myEquipment.value.filter(equipment => equipment && !equipment.is_horizontal);
                const subtype = equipmentSubtype(card);
                return myEquipment.value.filter(equipment => {
                    if (!equipment || equipment.is_horizontal) return false;
                    return !subtype || equipmentSubtype(equipment) === subtype;
                });
            }
            function sameSubtypeEquipment(card) {
                const subtype = equipmentSubtype(card);
                if (!subtype) return [];
                return myEquipment.value.filter(equipment => equipment && equipmentSubtype(equipment) === subtype);
            }
            function canSetCounterCard(card) {
                return !!card && card.is_counter_trap && isMyTurn.value && gameState.value.phase === 'main' && hasEmptyEquipmentSlot();
            }
            function setCounterCard(card) {
                if (!canSetCounterCard(card)) {
                    addLog(hasEmptyEquipmentSlot() ? `${card.name} 无法盖放` : `${card.name} 无法盖放，装备区已满`);
                    return;
                }
                selectedHandCard.value = null;
                sendAction('use_item', { instance_id: card.instance_id }, null, true);
                addLog(`盖放 ${card.name}`);
            }
            function canEquipCard(card) {
                if (!isEquippableItem(card) || !isMyTurn.value || gameState.value.phase !== 'main') return false;
                if (!canPayCost(playCost(card))) return false;
                const sameSubtype = sameSubtypeEquipment(card);
                if (sameSubtype.length > 0) {
                    return sameSubtype.some(equipment => !equipment.is_horizontal);
                }
                return hasEmptyEquipmentSlot() || replaceableEquipmentForCard(card).length > 0;
            }
            function equipCard(card) {
                if (!canEquipCard(card)) {
                    addLog(`${card.name} 无法装备`);
                    return;
                }
                selectedHandCard.value = null;
                if (sameSubtypeEquipment(card).length > 0 || !hasEmptyEquipmentSlot()) {
                    equipmentReplaceRequest.value = { card, replaceId: '' };
                    return;
                }
                sendEquipCard(card, '');
            }
            function sendEquipCard(card, replaceId) {
                sendPaidAction('equip', { instance_id: card.instance_id, replace_id: replaceId || '' }, playCost(card), `装备 ${card.name}`, () => {
                    addLog(replaceId ? `装备 ${card.name}，替换旧装备` : `装备 ${card.name}`);
                    equipmentReplaceRequest.value = null;
                });
            }
            function confirmEquipmentReplacement() {
                if (!equipmentReplaceRequest.value?.card || !equipmentReplaceRequest.value.replaceId) return;
                sendEquipCard(equipmentReplaceRequest.value.card, equipmentReplaceRequest.value.replaceId);
            }

            function selectHandCard(card) {
                if (!isMyTurn.value || gameState.value.phase !== 'main') return;
                selectedAttacker.value = null;
                selectedSkill.value = null;
                selectedItemSpellScroll.value = null;
                selectedSpellScrollSacrificeId.value = '';
                selectedSpellScrollSupportId.value = '';
                selectedSpellPowerSacrificeId.value = '';
                if (!canPlayHandCard(card)) {
                    addLog(handCardUnavailableText(card));
                    selectedHandCard.value = null;
                    return;
                }
                if (selectedHandCard.value === card) {
                    if (card.type === '道具' && !card.is_terrain) {
                        if (isConsumable(card)) {
                            if (card.is_counter_trap) {
                                setCounterCard(card);
                            } else if (isSpellScrollItem(card) && card.needs_target) {
                                selectedItemSpellScroll.value = card;
                                selectedSpellScrollSacrificeId.value = '';
                                selectedSpellScrollSupportId.value = '';
                                const prefix = card.preparation === 'fire_sacrifice'
                                    ? `选择 ${card.name}，先点击友方火焰伙伴祭献`
                                    : card.preparation === 'companion_support'
                                        ? `选择 ${card.name}，先点击友方伙伴支援`
                                        : `选择 ${card.name}，点击法力范围内目标`;
                                addLog(prefix);
                            } else {
                                sendPaidAction('use_item', { instance_id: card.instance_id }, playCost(card), `使用 ${card.name}`, () => addLog(`使用 ${card.name}`));
                            }
                        } else {
                            equipCard(card);
                        }
                        selectedHandCard.value = null;
                        return;
                    }
                    selectedHandCard.value = null;
                } else {
                    selectedHandCard.value = card;
                    if (summonDevourRequirement(card)) {
                        beginSummonDevour(card);
                        return;
                    }
                    if (card.type === '伙伴') addLog(`选择 ${card.name}，点击空格放置`);
                    else if (card.type === '道具' && card.is_terrain) addLog(`选择地形 ${card.name}，点击空格放置`);
                    else if (card.type === '道具' && isSpellScrollItem(card) && card.needs_target) {
                        selectedHandCard.value = null;
                        selectedItemSpellScroll.value = card;
                        selectedSpellScrollSacrificeId.value = '';
                        selectedSpellScrollSupportId.value = '';
                        const prefix = card.preparation === 'fire_sacrifice'
                            ? `选择 ${card.name}，先点击友方火焰伙伴祭献`
                            : card.preparation === 'companion_support'
                                ? `选择 ${card.name}，先点击友方伙伴支援`
                                : `选择 ${card.name}，点击法力范围内目标`;
                        addLog(prefix);
                    }
                    else if (card.type === '道具' && isConsumable(card)) addLog(`选择 ${card.name}，再次点击使用`);
                    else if (card.type === '道具') addLog(`选择 ${card.name}，再次点击装备`);
                }
            }

            function isHandCompanionCard(card) {
                return !!card && !!card.is_companion;
            }

            function canSummonAt(col, row) {
                if (!selectedHandCard.value || !isHandCompanionCard(selectedHandCard.value)) return false;
                const unit = getMyUnit(col, row);
                if (!unit) return true;
                return pendingSummonDevour.value[selectedHandCard.value.instance_id]?.includes(unit.instance_id);
            }

            function canUseMyCell(col, row) {
                if (!selectedHandCard.value) return false;
                if (isHandCompanionCard(selectedHandCard.value)) return canSummonAt(col, row);
                if (selectedHandCard.value.is_terrain) return !getMyTerrain(col, row);
                if (selectedHandCard.value.type === '伙伴') return !getMyUnit(col, row);
                return false;
            }

            function handleCellClick(ownerSlot, col, row) {
                if (ownerSlot === mySlot.value && selectedHandCard.value && isHandCompanionCard(selectedHandCard.value) && canSummonAt(col, row)) {
                    summonCard(selectedHandCard.value, col, row);
                    return;
                }
                if (ownerSlot === mySlot.value && selectedHandCard.value && selectedHandCard.value.is_terrain) {
                    if (!getMyTerrain(col, row)) { placeTerrainCard(selectedHandCard.value, col, row); return; }
                }
                if (ownerSlot === mySlot.value && selectedHandCard.value && selectedHandCard.value.type === '伙伴') {
                    if (!getMyUnit(col, row)) { summonCard(selectedHandCard.value, col, row); return; }
                }
                if (ownerSlot === mySlot.value && selectedItemSpellScroll.value?.preparation === 'fire_sacrifice' && canSacrificeForFlameArray(getMyUnit(col, row))) {
                    selectedSpellScrollSacrificeId.value = getMyUnit(col, row).instance_id;
                    addLog(`魔焰阵卷轴：祭献 ${getMyUnit(col, row).name}，再点击敌方目标。`);
                    return;
                }
                if (ownerSlot === mySlot.value && selectedItemSpellScroll.value?.preparation === 'companion_support' && canSupportOracleGlory(getMyUnit(col, row))) {
                    selectedSpellScrollSupportId.value = getMyUnit(col, row).instance_id;
                    addLog(`神谕卷轴 荣耀：支援 ${getMyUnit(col, row).name}，再点击敌方目标。`);
                    return;
                }
                if (ownerSlot === mySlot.value && selectedSkill.value && selectedSkillPowerSacrificeSource() && canSacrificeForHellRoar(getMyUnit(col, row))) {
                    selectedSpellPowerSacrificeId.value = getMyUnit(col, row).instance_id;
                    addLog(`炎狱怒吼：祭献 ${getMyUnit(col, row).name}，再点击敌方目标。`);
                    return;
                }
                if (ownerSlot === mySlot.value && selectedItemSpellScroll.value && isFriendlySpellTarget(col, row)) { useSpellScrollAt(col, row, ownerSlot); return; }
                if (ownerSlot === mySlot.value && selectedSkill.value && isFriendlySpellTarget(col, row)) { castSpellAt(col, row, ownerSlot); return; }
                if (ownerSlot !== mySlot.value && selectedAttacker.value && isValidOpponentTarget(col, row)) { attackTarget(col, row); return; }
                if (ownerSlot !== mySlot.value && selectedItemSpellScroll.value && isValidOpponentTarget(col, row)) { useSpellScrollAt(col, row, ownerSlot); return; }
                if (ownerSlot !== mySlot.value && selectedSkill.value && isValidOpponentTarget(col, row)) { castSpellAt(col, row, ownerSlot); return; }
            }

            function handleUnitImageClick(ownerSlot, col, row, unit, event) {
                if ((selectedSkill.value || selectedItemSpellScroll.value) && ((ownerSlot === mySlot.value && isFriendlySpellTarget(col, row)) || ownerSlot !== mySlot.value)) {
                    handleCellClick(ownerSlot, col, row);
                    return;
                }
                if (ownerSlot !== mySlot.value && selectedAttacker.value) {
                    handleCellClick(ownerSlot, col, row);
                    return;
                }
                if (ownerSlot !== mySlot.value && selectedItemSpellScroll.value) {
                    handleCellClick(ownerSlot, col, row);
                    return;
                }
                if (ownerSlot === mySlot.value && isHandCompanionCard(selectedHandCard.value) && canSummonAt(col, row)) {
                    summonCard(selectedHandCard.value, col, row);
                    return;
                }
                if (ownerSlot === mySlot.value && !isSpectator) {
                    selectedBattleID.value = unit.instance_id; fanPreviewID.value = '';
                    previewCard(unit); showCardContextMenu(event, unit, 'field');
                } else showCardDetail(unit);
            }

            function placeTerrainCard(card, col, row) {
                sendPaidAction('place_terrain', { instance_id: card.instance_id, col, row }, playCost(card), `放置 ${card.name}`, () => {
                    selectedHandCard.value = null;
                    addLog(`放置地形 ${card.name} → (${col},${row})`);
                });
            }

            function summonDevourRequirement(card) {
                if (card?.devour_requirement || card?.devour_card_requirement) return card.devour_requirement || {};
                return null;
            }

            function loadSatisfiesRequirement(unit, requirement) {
                if (!unit || !requirement) return false;
                const gain = loadContribution(unit, requirement);
                return Object.entries(requirement).every(([elem, amount]) => Number(gain[elem] || 0) >= Number(amount || 0));
            }

            function loadContribution(card, requirement = null) {
                const total = {};
                for (const [elem, amount] of Object.entries(card?.elements_gain || {})) {
                    total[elem] = Number(total[elem] || 0) + Math.max(0, Number(amount || 0));
                }
                if (requirement && Object.prototype.hasOwnProperty.call(requirement, '血')) {
                    total['血'] = Math.max(0, Number(card?.current_life ?? card?.life ?? 0));
                }
                return total;
            }

            function addLoadMaps(cards, requirement = null) {
                const total = {};
                for (const card of cards) {
                    for (const [elem, amount] of Object.entries(loadContribution(card, requirement))) {
                        total[elem] = Number(total[elem] || 0) + Math.max(0, Number(amount || 0));
                    }
                }
                return total;
            }

            const devourSelectedLoad = computed(() => {
                if (!devourRequest.value) return {};
                const selected = new Set(devourSelection.value);
                return addLoadMaps(devourCandidates(devourRequest.value.card).filter(card => selected.has(card.instance_id)), devourRequest.value.requirement);
            });

            const devourSelectionMeetsRequirement = computed(() => {
                const requirement = devourRequest.value?.requirement || {};
                const cardRequirement = devourRequest.value?.card?.devour_card_requirement || null;
                return Object.entries(requirement).every(([elem, amount]) => Number(devourSelectedLoad.value[elem] || 0) >= Number(amount || 0))
                    && selectedDevourCardCount(cardRequirement) >= Number(cardRequirement?.count || 0);
            });

            function devourCandidates(card) {
                const requirement = summonDevourRequirement(card);
                if (!requirement) return [];
                const cardRequirement = card?.devour_card_requirement || null;
                const candidates = [];
                for (let col = 0; col < 3; col++) {
                    for (let row = 0; row < 3; row++) {
                        const unit = getMyUnit(col, row);
                        if (!unit || unit.instance_id === myState.value?.hero?.instance_id) continue;
                        if (hasDevourLoadContribution(unit, requirement) || satisfiesDevourCardRequirement(unit, cardRequirement)) candidates.push(unit);
                    }
                }
                for (const equipment of (myState.value?.equipment || [])) {
                    if (equipment && hasDevourLoadContribution(equipment, requirement)) candidates.push(equipment);
                }
                return candidates;
            }

            function hasDevourLoadContribution(card, requirement) {
                const contribution = loadContribution(card, requirement);
                return Object.keys(requirement || {}).some(elem => Number(contribution[elem] || 0) > 0);
            }

            function satisfiesDevourCardRequirement(card, requirement) {
                if (!requirement || Number(requirement.count || 0) <= 0) return false;
                if (!card) return false;
                if (requirement.companion_only && card.type !== '伙伴') return false;
                if (requirement.category && card.category !== requirement.category) return false;
                return true;
            }

            function selectedDevourCardCount(requirement) {
                if (!requirement || Number(requirement.count || 0) <= 0 || !devourRequest.value) return 0;
                const selected = new Set(devourSelection.value);
                return devourCandidates(devourRequest.value.card)
                    .filter(card => selected.has(card.instance_id) && satisfiesDevourCardRequirement(card, requirement))
                    .length;
            }

            function beginSummonDevour(card, col = null, row = null) {
                const candidates = devourCandidates(card);
                const requirement = summonDevourRequirement(card) || {};
                const cardRequirement = card?.devour_card_requirement || null;
                const loadSatisfied = Object.entries(requirement).every(([elem, amount]) => Number(addLoadMaps(candidates, requirement)[elem] || 0) >= Number(amount || 0));
                const cardSatisfied = !cardRequirement || Number(cardRequirement.count || 0) <= 0
                    || candidates.filter(candidate => satisfiesDevourCardRequirement(candidate, cardRequirement)).length >= Number(cardRequirement.count || 0);
                if (!candidates.length || !loadSatisfied || !cardSatisfied) {
                    addLog(`${card.name} 没有满足吞噬要求的友方卡牌`);
                    return;
                }
                devourSelection.value = [];
                devourRequest.value = { card, col, row, requirement };
            }

            function toggleDevourSelection(candidate) {
                const id = candidate?.instance_id;
                if (!id) return;
                if (devourSelection.value.includes(id)) {
                    devourSelection.value = devourSelection.value.filter(existing => existing !== id);
                } else {
                    devourSelection.value = [...devourSelection.value, id];
                }
            }

            function confirmDevourSelection() {
                if (!devourRequest.value || !devourSelectionMeetsRequirement.value) return;
                const { card, col, row } = devourRequest.value;
                pendingSummonDevour.value = {
                    ...pendingSummonDevour.value,
                    [card.instance_id]: [...devourSelection.value],
                };
                devourRequest.value = null;
                if (col !== null && row !== null) summonCard(card, col, row);
                else addLog("Devour selected; choose a position to summon " + card.name);
            }

            function cancelDevourSelection() {
                devourRequest.value = null;
                devourSelection.value = [];
            }

            function summonCard(card, col, row) {
                const requirement = summonDevourRequirement(card);
                const devourIDs = pendingSummonDevour.value[card.instance_id] || [];
                if (requirement && !devourIDs.length) {
                    beginSummonDevour(card);
                    return;
                }
                const data = { instance_id: card.instance_id, col, row };
                if (devourIDs.length) data.devour_ids = devourIDs;
                sendPaidAction('summon', data, playCost(card), `召唤 ${card.name}`, () => {
                    const next = { ...pendingSummonDevour.value };
                    delete next[card.instance_id];
                    pendingSummonDevour.value = next;
                    selectedHandCard.value = null;
                    devourSelection.value = [];
                    addLog(`召唤 ${card.name} → (${col},${row})`);
                });
            }

            function canConsumeUnit(unit) {
                if (!unit || unit.is_horizontal) return false;
                if (unit.statuses?.['眩晕'] > 0 || unit.statuses?.['冷却'] > 0) return false;
                return isMyTurn.value && gameState.value.phase === 'main';
            }
            function unitActionChips(unit) {
                if (!isMyTurn.value || gameState.value.phase !== 'main' || !unit) return [];
                const chips = [];
                if (unit.bound_skills?.length) chips.push('绑');
                if (canConsumeUnit(unit)) chips.push('消');
                if (canAttackWith(unit)) chips.push('攻');
                if (canUsePerTurn(unit)) chips.push('技');
                if (canUseUltimate(unit)) chips.push('绝');
                return chips;
            }
            function consumeCard(unit) {
                if (isFirstTurnHeroConsumeChoice(unit)) {
                    openConsumeGainChoice(unit);
                    return;
                }
                sendAction('consume', { instance_id: unit.instance_id });
                addLog(`消耗 ${unit.name}`);
            }

            function canAttackWith(unit) {
                if (!isMyTurn.value || gameState.value.phase !== 'main') return false;
                if (!unit || unit.is_horizontal) return false;
                if ((unit.current_attack ?? unit.attack) <= 0) return false;
                if (unit.position && unit.position.row !== myFrontRow.value && !unit.can_attack_from_non_front) return false;
                return true;
            }
            function selectAttacker(unit) {
                selectedHandCard.value = null; selectedSkill.value = null;
                selectedAttacker.value = selectedAttacker.value === unit ? null : unit;
                if (selectedAttacker.value) addLog(`选择 ${unit.name} 攻击，点击敌方目标`);
            }
            function isAttackTarget(col, row) {
                return selectedAttacker.value && row === opponentFrontRow.value && getOpponentUnit(col, row);
            }
            const hasActiveTargeting = computed(() => !!selectedAttacker.value || !!selectedSkill.value || !!selectedItemSpellScroll.value);
            function isValidOpponentTarget(col, row) {
                const unit = getOpponentUnit(col, row);
                if (selectedAttacker.value) return !!unit && row === opponentFrontRow.value;
                const spell = activeTargetSpell();
                if (!spell) return false;
                if (pendingExtraTargetCast.value && spell.extra_target_ignore_range) return !!unit;
                const area = spell.spell_area || 'single';
                if (area === 'front_row') {
                    return !!unit && row === opponentFrontRow.value;
                }
                const inRange = activeTargetSpellHasPierce() || hasGlobalSpellRange() || row === opponentFrontRow.value || isTauntUnit(unit);
                if (unit) return inRange;
                if (area === 'column' && inRange) {
                    return opponentColumnHasUnit(col);
                }
                return false;
            }
            function isSpellTarget(col, row) {
                return !!activeTargetSpell() && isValidOpponentTarget(col, row);
            }
            function hoverOpponentCell(col, row) {
                hoveredOpponentCell.value = { col, row };
            }
            function clearHoverOpponentCell(col, row) {
                if (hoveredOpponentCell.value?.col === col && hoveredOpponentCell.value?.row === row) {
                    hoveredOpponentCell.value = null;
                }
            }
            function isSpellAreaPreview(col, row) {
                if (!selectedSkill.value || !hoveredOpponentCell.value) return false;
                const center = hoveredOpponentCell.value;
                if (!isValidOpponentTarget(center.col, center.row)) return false;
                const area = selectedSkill.value.spell_area || 'single';
                if (area === 'square') {
                    const startCol = Math.min(Math.max(center.col, 0), 1);
                    const startRow = Math.min(Math.max(center.row, 0), 1);
                    return col >= startCol && col < startCol + 2 && row >= startRow && row < startRow + 2;
                }
                if (area === 'all') return true;
                if (area === 'column') return col === center.col;
                if (area === 'front_row') return row === opponentFrontRow.value;
                if (area === 'splash_cross') {
                    return (col === center.col && row === center.row) ||
                        (Math.abs(col - center.col) + Math.abs(row - center.row) === 1);
                }
                return col === center.col && row === center.row;
            }
            function isFriendlySpellTarget(col, row) {
                const spell = activeTargetSpell();
                if (!spell || !getMyUnit(col, row)) return false;
                if (selectedItemSpellScroll.value?.preparation === 'fire_sacrifice' && !selectedSpellScrollSacrificeId.value) {
                    return canSacrificeForFlameArray(getMyUnit(col, row));
                }
                if (selectedItemSpellScroll.value?.preparation === 'companion_support' && !selectedSpellScrollSupportId.value) {
                    return canSupportOracleGlory(getMyUnit(col, row));
                }
                if (selectedSkill.value && selectedSkillPowerSacrificeSource() && canSacrificeForHellRoar(getMyUnit(col, row))) {
                    return true;
                }
                return !!spell.allows_friendly_target || canAttackSpellTargetFriendlyUnit(spell);
            }
            function isFriendlyHeroSpellTarget() {
                return !!selectedSkill.value && !!selectedSkill.value.allows_friendly_target && !!myState.value?.hero;
            }
            function canAttackSpellTargetFriendlyUnit(spell) {
                return !!spell && !!spell.needs_target && spell.can_attack !== false && !spell.is_sorcery;
            }
            function activeTargetSpell() {
                return selectedSkill.value || selectedItemSpellScroll.value || null;
            }
            function activeTargetSpellHasPierce() {
                const spell = activeTargetSpell();
                if (!spell) return false;
                if (spell.has_pierce) return true;
                if (!selectedSkill.value) return false;
                return selectedBoostSkills.value.some(id => {
                    const boost = defensePaymentCardById(id);
                    return boost?.boost_grants_pierce;
                });
            }
            function isTauntUnit(unit) {
                return !!unit && (unit.has_taunt || (unit.statuses?.['引魔'] > 0));
            }
            function hasGlobalSpellRange() {
                return allMyFieldCards().some(card => card && card.has_global_spell_range && !(card.statuses?.['石化'] > 0));
            }
            function opponentColumnHasUnit(col) {
                for (let row = 0; row < 3; row++) {
                    if (getOpponentUnit(col, row)) return true;
                }
                return false;
            }
            function attackTarget(col, row) {
                if (!selectedAttacker.value) return;
                sendAction('attack', { attacker_id: selectedAttacker.value.instance_id, target_col: col, target_row: row });
                addLog(`${selectedAttacker.value.name} 攻击 → (${col},${row})`);
                selectedAttacker.value = null;
            }

            function canCastSkill(skill) {
                if (!isMyTurn.value || gameState.value.phase !== 'main' || !skill || skill.is_horizontal) return false;
                if (skill.statuses?.['冷却'] > 0 || skill.statuses?.['石化'] > 0) return false;
                if (!skill.can_attack) return false;
                if (!canPayCost(attackActionCost(skill, []))) return false;
                return true;
            }
            function skillActionChips(skill) {
                if (!isMyTurn.value || gameState.value.phase !== 'main' || !skill) return [];
                const chips = [];
                if (canCastSkill(skill)) chips.push('施');
                if (selectedSkill.value && canBoostWith(skill)) chips.push('强');
                if (canUsePerTurn(skill)) chips.push('技');
                return chips;
            }
            function canShowCostHint(skill) {
                return !!skill && isMyTurn.value && gameState.value.phase === 'main' &&
                    !skill.is_horizontal && !canCastSkill(skill) && !canPayCost(attackActionCost(skill, []));
            }
            function selectSkillToCast(skill) {
                selectedHandCard.value = null; selectedAttacker.value = null;
                selectedSkill.value = selectedSkill.value === skill ? null : skill;
                selectedBoostSkills.value = [];
                selectedSpellPowerSacrificeId.value = '';
                if (!selectedSkill.value) return;
                if (!skillNeedsTarget(skill)) return;
                addLog(`选择法术 ${skill.name}，点击法力范围内目标`);
            }

            // Boost skill selection for attack casting
            const selectedBoostSkills = ref([]);
            const pendingExtraTargetCast = ref(null);

            function canBoostWith(skill) {
                if (!selectedSkill.value || !skill) return false;
                if (selectedSkill.value.can_be_boosted === false) return false;
                if (skill.instance_id === selectedSkill.value.instance_id) return false;
                if (skill.type === '道具') {
                    if (!isSpellScrollItem(skill) || skill.statuses?.['石化'] > 0 || !skill.can_attack_boost) return false;
                } else if (!canUseSkillAsAttackBoost(skill)) {
                    return false;
                }
                const boostIds = selectedBoostSkills.value.includes(skill.instance_id)
                    ? [...selectedBoostSkills.value]
                    : [...selectedBoostSkills.value, skill.instance_id];
                if (!canPayCost(attackActionCost(selectedSkill.value, boostIds))) return false;
                return true;
            }
            function selectedSkillHasPierce() {
                return activeTargetSpellHasPierce();
            }
            function matchingExtraTargetModifier(skill) {
                if (!skill) return null;
                return (myState.value.temp_modifiers || []).find(m => {
                    if (!m || (m.remaining_uses || 0) === 0) return false;
                    if (m.type === 'next_spell_extra_target') {
                        return !m.target_instance_id || m.target_instance_id === skill.instance_id;
                    }
                    if (m.type === 'next_drive_spell_extra_target') {
                        return !!skill.is_drive;
                    }
                    return false;
                }) || null;
            }
            function selectedSkillNeedsExtraTarget(skill) {
                if (!skill) return false;
                if (skill.requires_extra_target) return true;
                return !!matchingExtraTargetModifier(skill);
            }
            function extraTargetCanRepeatMain(skill) {
                if (!skill) return false;
                if (skill.extra_target_can_repeat) return true;
                return !!matchingExtraTargetModifier(skill)?.allow_same_target;
            }
            function toggleBoostSkill(skill) {
                const idx = selectedBoostSkills.value.indexOf(skill.instance_id);
                if (idx >= 0) selectedBoostSkills.value.splice(idx, 1);
                else selectedBoostSkills.value.push(skill.instance_id);
                if (!selectedSkillPowerSacrificeSource()) {
                    selectedSpellPowerSacrificeId.value = '';
                }
            }

            function useSpellScrollAt(col, row, targetOwner = 1 - mySlot.value) {
                if (!selectedItemSpellScroll.value) return;
                const card = selectedItemSpellScroll.value;
                if (card.preparation === 'fire_sacrifice' && !selectedSpellScrollSacrificeId.value) {
                    addLog('魔焰阵卷轴：先选择一个入场费用大于4的友方火焰伙伴作为祭献。');
                    return;
                }
                if (card.preparation === 'companion_support' && !selectedSpellScrollSupportId.value) {
                    addLog('神谕卷轴 荣耀：先选择一个生命值+负载大于5的友方伙伴作为支援。');
                    return;
                }
                const data = {
                    instance_id: card.instance_id,
                    target_type: 'unit',
                    target_col: col,
                    target_row: row,
                    target_owner: targetOwner,
                };
                if (selectedSpellScrollSacrificeId.value) {
                    data.sacrifice_id = selectedSpellScrollSacrificeId.value;
                }
                if (selectedSpellScrollSupportId.value) {
                    data.support_id = selectedSpellScrollSupportId.value;
                }
                sendPaidAction('use_item', data, playCost(card), `使用 ${card.name}`, () => {
                    addLog(`使用 ${card.name} → (${col},${row})`);
                    selectedItemSpellScroll.value = null;
                    selectedHandCard.value = null;
                    selectedSpellScrollSacrificeId.value = '';
                    selectedSpellScrollSupportId.value = '';
                });
            }
            function castSpellAt(col, row, targetOwner = 1 - mySlot.value) {
                if (!selectedSkill.value) return;
                if (pendingExtraTargetCast.value && selectedSkillNeedsExtraTarget(selectedSkill.value)) {
                    const data = { ...pendingExtraTargetCast.value };
                    if (extraTargetCanRepeatMain(selectedSkill.value) || data.target_col !== col || data.target_row !== row) {
                        data.extra_target_col = col;
                        data.extra_target_row = row;
                    }
                    const skillName = selectedSkill.value.name;
                    const boostCount = selectedBoostSkills.value.length;
                    const cost = attackActionCost(selectedSkill.value, selectedBoostSkills.value);
                    sendPaidAction('cast_spell', data, cost, `施放 ${skillName}`, () => {
                        addLog(`施放 ${skillName} → (${data.target_col},${data.target_row})` +
                            (data.extra_target_col !== undefined ? ` + 额外(${data.extra_target_col},${data.extra_target_row})` : '') +
                            (boostCount > 0 ? ` (+${boostCount}强化)` : ''));
                        selectedSkill.value = null;
                        selectedBoostSkills.value = [];
                        pendingExtraTargetCast.value = null;
                        selectedSpellPowerSacrificeId.value = '';
                    });
                    return;
                }
                const data = {
                    instance_id: selectedSkill.value.instance_id,
                    target_type: 'unit',
                    target_col: col,
                    target_row: row,
                    target_owner: targetOwner,
                };
                if (selectedBoostSkills.value.length > 0) {
                    data.boost_ids = [...selectedBoostSkills.value];
                }
                if (selectedSpellPowerSacrificeId.value) {
                    data.sacrifice_id = selectedSpellPowerSacrificeId.value;
                }
                if (selectedSkillNeedsExtraTarget(selectedSkill.value)) {
                    pendingExtraTargetCast.value = data;
                    const repeatHint = extraTargetCanRepeatMain(selectedSkill.value) ? '可以再次点击同一目标。' : '再次点击同一目标则不选择额外目标。';
                    addLog(`${selectedSkill.value.name}：再点击一个敌方单位作为额外目标；${repeatHint}`);
                    return;
                }
                const skillName = selectedSkill.value.name;
                const boostCount = selectedBoostSkills.value.length;
                const cost = attackActionCost(selectedSkill.value, selectedBoostSkills.value);
                sendPaidAction('cast_spell', data, cost, `施放 ${skillName}`, () => {
                    addLog(`施放 ${skillName} → (${col},${row})` +
                        (boostCount > 0 ? ` (+${boostCount}强化)` : ''));
                    selectedSkill.value = null;
                    selectedBoostSkills.value = [];
                    pendingExtraTargetCast.value = null;
                    selectedSpellPowerSacrificeId.value = '';
                });
            }

            function castSpellAtHero() {
                if (!selectedSkill.value || !isFriendlyHeroSpellTarget()) return;
                const data = {
                    instance_id: selectedSkill.value.instance_id,
                    target_type: 'hero'
                };
                if (selectedBoostSkills.value.length > 0) {
                    data.boost_ids = [...selectedBoostSkills.value];
                }
                if (selectedSpellPowerSacrificeId.value) {
                    data.sacrifice_id = selectedSpellPowerSacrificeId.value;
                }
                const skillName = selectedSkill.value.name;
                const boostCount = selectedBoostSkills.value.length;
                const cost = attackActionCost(selectedSkill.value, selectedBoostSkills.value);
                sendPaidAction('cast_spell', data, cost, `施放 ${skillName}`, () => {
                    addLog(`施放 ${skillName} → 人物` +
                        (boostCount > 0 ? ` (+${boostCount}强化)` : ''));
                    selectedSkill.value = null;
                    selectedBoostSkills.value = [];
                    pendingExtraTargetCast.value = null;
                    selectedSpellPowerSacrificeId.value = '';
                });
            }

            function handleMyHeroClick() {
                if (isFriendlyHeroSpellTarget()) {
                    castSpellAtHero();
                    return;
                }
                showCardDetail(myState.value.hero);
            }

            function castSpellWithoutTarget() {
                if (!selectedSkill.value) return;
                const data = {
                    instance_id: selectedSkill.value.instance_id,
                    target_type: 'none'
                };
                if (selectedBoostSkills.value.length > 0) {
                    data.boost_ids = [...selectedBoostSkills.value];
                }
                if (selectedSpellPowerSacrificeId.value) {
                    data.sacrifice_id = selectedSpellPowerSacrificeId.value;
                }
                const skillName = selectedSkill.value.name;
                const cost = attackActionCost(selectedSkill.value, selectedBoostSkills.value);
                sendPaidAction('cast_spell', data, cost, `施放 ${skillName}`, () => {
                    addLog(`施放 ${skillName}`);
                    selectedSkill.value = null;
                    selectedBoostSkills.value = [];
                    pendingExtraTargetCast.value = null;
                    selectedSpellPowerSacrificeId.value = '';
                });
            }

            function skillNeedsTarget(skill) {
                if (!skill) return false;
                return !!skill.needs_target;
            }

            function canSacrificeForFlameArray(unit) {
                return isFireCompanion(unit) && totalCost(playCost(unit)) > 4;
            }

            function canSacrificeForHellRoar(unit) {
                return isFireCompanion(unit);
            }

            function canSupportOracleGlory(unit) {
                return !!unit && unit.type === '伙伴' && (Number(unit.current_life || 0) + cardLoadValue(unit)) > 5;
            }

            function selectedSkillPowerSacrificeSource() {
                if (selectedSkill.value && selectedSkill.value.preparation === 'optional_fire_sacrifice') return selectedSkill.value;
                return selectedBoostSkills.value
                    .map(id => allMySkills.value.find(skill => skill && skill.instance_id === id))
                    .find(skill => skill && skill.preparation === 'optional_fire_sacrifice') || null;
            }

            // Defense skill selection state
            const defenseSelected = ref([]);
            const defenseBoostSelected = ref([]);
            const defenseOverexertSelected = ref([]);

            const availableDefenseSkills = computed(() => {
                if (!showDefenseWindow.value) return [];
                return allMySkills.value.filter(canUseSkillAsDefense);
            });

            const defensePowerSacrificeCandidates = computed(() => {
                if (!showDefenseWindow.value || !hasDefenseSource()) return [];
                const selectedSources = [...defenseSelected.value, ...defenseBoostSelected.value]
                    .map(id => allMySkills.value.find(s => s && s.instance_id === id));
                if (!selectedSources.some(s => s && s.preparation === 'optional_fire_sacrifice')) return [];
                return myUnitCards().filter(isFireCompanion);
            });

            const availableDefenseScrolls = computed(() => {
                if (!showDefenseWindow.value) return [];
                return myHand.value.filter(card => isSpellScrollItem(card) && card.can_defend);
            });

            const availableReactionSkills = computed(() => {
                if (!showDefenseWindow.value) return [];
                return [
                    ...allMySkills.value,
                    ...(myState.value.equipment || []),
                ].filter(Boolean).filter(canUseSkillAsReaction);
            });

            const defenseOverexertUnits = computed(() => {
                if (!showDefenseWindow.value || !myState.value.units) return [];
                const units = [];
                for (let col = 0; col < 3; col++) {
                    for (let row = 0; row < 3; row++) {
                        const unit = getMyUnit(col, row);
                        if (unit && canOverexertForDefense(unit)) units.push(unit);
                    }
                }
                for (const equipment of (myState.value.equipment || [])) {
                    if (equipment && canOverexertForDefense(equipment)) units.push(equipment);
                }
                return units;
            });

            const availableDefenseBoostCandidates = computed(() => {
                if (!showDefenseWindow.value) return [];
                const learnedBoosts = allMySkills.value.filter(s =>
                    canUseSkillAsDefenseBoost(s)
                );
                const scrollBoosts = myHand.value.filter(s =>
                    canUseScrollAsDefenseBoost(s)
                );
                return [...learnedBoosts, ...scrollBoosts];
            });

            const availableBoostSkills = computed(() => {
                return availableDefenseBoostCandidates.value.filter(skill => !isDefenseBoostDisabled(skill));
            });

            function hasDefenseSource() {
                return defenseSelected.value.length + defenseScrollSelected.value.length > 0;
            }

            function currentDefenseSources() {
                return [
                    ...allMySkills.value.filter(s => s && defenseSelected.value.includes(s.instance_id)),
                    ...myHand.value.filter(s => s && defenseScrollSelected.value.includes(s.instance_id)),
                ];
            }

            function isDefenseBoostDisabled(skill) {
                if (!skill || !showDefenseWindow.value || !hasDefenseSource()) return true;
                if (defenseSelected.value.includes(skill.instance_id)) return true;
                if (defenseScrollSelected.value.includes(skill.instance_id)) return true;
                return currentDefenseSources().some(source => source?.can_be_boosted === false);
            }

            function canOverexertForDefense(unit) {
                return !!unit && !unit.is_horizontal &&
                    !(unit.statuses?.['眩晕'] > 0) &&
                    !(unit.statuses?.['冷却'] > 0);
            }

            function defensePaymentCardById(id) {
                return allMySkills.value.find(s => s && s.instance_id === id) || myHand.value.find(c => c && c.instance_id === id) || null;
            }

            function defensePaymentCost(card, applyNextSkillFree = true) {
                if (!card) return {};
                return card.type === '道具' ? playCost(card) : skillUseCost(card, applyNextSkillFree);
            }

            function splitActionCostCards(ids) {
                const skills = [];
                const cards = [];
                for (const id of ids || []) {
                    const skill = allMySkills.value.find(candidate => candidate?.instance_id === id);
                    if (skill) {
                        skills.push(skill);
                        continue;
                    }
                    const card = myHand.value.find(candidate => candidate?.instance_id === id);
                    if (card) cards.push(card);
                }
                return { skills, cards };
            }

            function plannedActionCost(entries) {
                return ArcaneActionCostPlan.planActionCost(entries, myState.value.temp_modifiers || []);
            }

            function attackActionCost(mainSkill, boostIDs) {
                if (!mainSkill) return {};
                const boosts = splitActionCostCards(boostIDs);
                return plannedActionCost([
                    { kind: 'skill', purpose: 'attack', card: mainSkill },
                    ...boosts.skills.map(card => ({ kind: 'skill', purpose: 'attack_boost', card })),
                    ...boosts.cards.map(card => ({ kind: 'card', card })),
                ]);
            }

            function defenseActionCost() {
                const defenseSkills = splitActionCostCards(defenseSelected.value).skills;
                const defenseScrolls = splitActionCostCards(defenseScrollSelected.value).cards;
                const boosts = splitActionCostCards(defenseBoostSelected.value);
                return plannedActionCost([
                    ...defenseSkills.map(card => ({ kind: 'skill', purpose: 'defend', card })),
                    ...defenseScrolls.map(card => ({ kind: 'card', card })),
                    ...boosts.skills.map(card => ({ kind: 'skill', purpose: 'defense_boost', card })),
                    ...boosts.cards.map(card => ({ kind: 'card', card })),
                ]);
            }

            function reactionActionCost(skill) {
                return plannedActionCost([{ kind: 'skill', purpose: 'reaction', card: skill }]);
            }

            const totalDefenseCost = computed(defenseActionCost);

            const defenseAvailableElements = computed(() => {
                const available = cloneElements(myState.value.elements || {});
                const units = defenseOverexertUnits.value;
                for (const id of defenseOverexertSelected.value) {
                    const unit = units.find(u => u.instance_id === id);
                    if (!unit) continue;
                    for (const [elem, amount] of Object.entries(unit.elements_gain || {})) {
                        available[elem] = Number(available[elem] || 0) + Number(amount || 0);
                    }
                }
                return available;
            });

            function isSkillReady(skill) {
                return !!skill && !skill.is_horizontal &&
                    !(skill.statuses?.['冷却'] > 0) &&
                    !(skill.statuses?.['石化'] > 0);
            }

            function canUseSkillAsDefense(skill) {
                if (!isSkillReady(skill)) return false;

                return !!skill.can_defend;
            }

            function reactionAvailableElements(includeAllOverexert = false) {
                const available = cloneElements(myState.value.elements || {});
                const selected = includeAllOverexert
                    ? defenseOverexertUnits.value.map(unit => unit.instance_id)
                    : defenseOverexertSelected.value;
                for (const id of selected) {
                    const unit = defenseOverexertUnits.value.find(u => u.instance_id === id);
                    if (!unit) continue;
                    for (const [elem, amount] of Object.entries(unit.elements_gain || {})) {
                        available[elem] = Number(available[elem] || 0) + Number(amount || 0);
                    }
                }
                return available;
            }

            function canUseSkillAsReaction(skill) {
                if (!skill) return false;
                if (!skill.can_react) return false;
                if (skill.type === '技能' && !isSkillReady(skill)) return false;
                if (skill.type !== '技能' && skill.statuses?.['石化'] > 0) return false;
                return calculatePaymentWithAvailable(reactionAvailableElements(true), reactionActionCost(skill)) !== null;
            }

            function canPaySelectedReaction(skill) {
                return calculatePaymentWithAvailable(reactionAvailableElements(false), reactionActionCost(skill)) !== null;
            }

            function canUseSkillAsAttackBoost(skill) {
                if (!isSkillReady(skill)) return false;
                return !!skill.can_attack_boost;
            }

            function canUseSkillAsDefenseBoost(skill) {
                if (!isSkillReady(skill)) return false;

                return !!skill.can_defense_boost;
            }

            function canUseScrollAsDefenseBoost(scroll) {
                if (!isSpellScrollItem(scroll)) return false;

                if (scroll.statuses?.['石化'] > 0) return false;
                return !!scroll.can_defense_boost;
            }

            const totalDefensePower = computed(() => {
                let power = 0;
                const skills = allMySkills.value;
                for (const id of defenseSelected.value) {
                    const sk = skills.find(s => s && s.instance_id === id);
                    if (sk) power += defensePower(sk);
                }
                for (const id of defenseScrollSelected.value) {
                    const scroll = myHand.value.find(s => s && s.instance_id === id);
                    if (scroll) power += defensePower(scroll);
                }
                for (const id of defenseBoostSelected.value) {
                    const sk = skills.find(s => s && s.instance_id === id) || myHand.value.find(s => s && s.instance_id === id);
                    if (sk) power += defenseBoostPower(sk);
                }
                if (defensePowerSacrificeId.value) {
                    const sacrifice = defensePowerSacrificeCandidates.value.find(unit => unit.instance_id === defensePowerSacrificeId.value);
                    if (sacrifice) power += totalCost(playCost(sacrifice));
                }
                return power;
            });

            function defensePower(card) {
                return Math.max(Number(card?.effective_defense_power ?? card?.power ?? 0), 0);
            }

            function defenseBoostPower(card) {
                return Math.max(Number(card?.effective_defense_boost_power ?? card?.power ?? 0), 0);
            }

            const canSubmitDefense = computed(() => {
                if (pendingAttackPower() <= 0) return true;
                if (!hasDefenseSource()) return false;
                return calculatePaymentWithAvailable(defenseAvailableElements.value, totalDefenseCost.value) !== null;
            });

            function pendingAttackPower() {
                return pendingSpellView.value?.power || 0;
            }

            function defenseSummaryText() {
                if (!hasDefenseSource()) return '请选择防御法术或卷轴';
                if (calculatePaymentWithAvailable(defenseAvailableElements.value, totalDefenseCost.value) === null) return '费用不足';
                return totalDefensePower.value >= pendingAttackPower() ? '防御成功' : '威力不足';
            }

            function defenseSummaryClass() {
                return hasDefenseSource() && calculatePaymentWithAvailable(defenseAvailableElements.value, totalDefenseCost.value) !== null && totalDefensePower.value >= pendingAttackPower()
                    ? 'power-win'
                    : 'power-lose';
            }

            function defenseCardById(id) {
                if (!id) return null;
                return [
                    ...availableReactionSkills.value,
                    ...defenseOverexertUnits.value,
                    ...availableDefenseScrolls.value,
                    ...availableDefenseSkills.value,
                    ...defensePowerSacrificeCandidates.value,
                    ...availableDefenseBoostCandidates.value,
                    ...availableBoostSkills.value,
                ].find(card => card && card.instance_id === id) || null;
            }

            function handleDefensePanelClick(event) {
                const actionEl = event.target.closest('[data-defense-action]');
                if (!actionEl || !event.currentTarget.contains(actionEl)) return;
                const action = actionEl.dataset.defenseAction;
                const card = defenseCardById(actionEl.dataset.instanceId || '');
                event.preventDefault();
                event.stopPropagation();
                if (action === 'detail') {
                    if (card) showCardDetail(card);
                    return;
                }
                if (action === 'reaction') {
                    if (card) reactSpell(card);
                    return;
                }
                if (action === 'overexert') {
                    if (card) toggleDefenseOverexert(card);
                    return;
                }
                if (action === 'scroll') {
                    if (card) toggleDefenseScroll(card);
                    return;
                }
                if (action === 'skill') {
                    if (card) toggleDefenseSkill(card);
                    return;
                }
                if (action === 'power_sacrifice') {
                    if (card) toggleDefensePowerSacrifice(card);
                    return;
                }
                if (action === 'boost') {
                    if (card) toggleDefenseBoost(card);
                    return;
                }
                if (action === 'submit') {
                    submitDefend();
                    return;
                }
                if (action === 'no_defend') {
                    noDefend();
                }
            }

            function syncDefensePanelVisualState() {
                const apply = () => {
                    if (!showDefenseWindow.value) return;
                    const panel = document.querySelector('.defense-panel');
                    if (!panel) return;
                    panel.querySelectorAll('.defense-skill-card[data-defense-action][data-instance-id]').forEach(el => {
                        const id = el.dataset.instanceId;
                        const action = el.dataset.defenseAction;
                        const selected =
                            (action === 'skill' && defenseSelected.value.includes(id)) ||
                            (action === 'scroll' && defenseScrollSelected.value.includes(id)) ||
                            (action === 'boost' && defenseBoostSelected.value.includes(id)) ||
                            (action === 'power_sacrifice' && defensePowerSacrificeId.value === id) ||
                            (action === 'overexert' && defenseOverexertSelected.value.includes(id));
                        el.classList.toggle('selected', selected);
                        if (action === 'boost') {
                            const card = defenseCardById(id);
                            el.classList.toggle('disabled', isDefenseBoostDisabled(card));
                        }
                    });
                    syncDefenseSummary(panel);
                    const submit = panel.querySelector('[data-defense-action="submit"]');
                    if (submit) submit.classList.toggle('disabled', !canSubmitDefense.value || !defensePaymentValid.value || actionBusy.value);
                };
                apply();
                nextTick(apply);
            }

            function syncDefenseSummary(panel) {
                const power = panel.querySelector('[data-defense-summary="power"]');
                if (power) {
                    power.innerHTML = `防御威力：<strong>${totalDefensePower.value}</strong> / 进攻威力：<strong>${pendingAttackPower()}</strong>`;
                }
                const cost = panel.querySelector('[data-defense-summary="cost"]');
                if (cost) {
                    cost.innerHTML = `费用：<strong>${costLabel(totalDefenseCost.value)}</strong> / 可用：<strong>${costLabel(defenseAvailableElements.value)}</strong>`;
                }
                const result = panel.querySelector('[data-defense-summary="result"]');
                if (result) {
                    result.textContent = defenseSummaryText();
                    result.classList.toggle('power-win', defenseSummaryClass() === 'power-win');
                    result.classList.toggle('power-lose', defenseSummaryClass() === 'power-lose');
                }
            }

            function toggleDefenseOverexert(unit) {
                const idx = defenseOverexertSelected.value.indexOf(unit.instance_id);
                if (idx >= 0) defenseOverexertSelected.value.splice(idx, 1);
                else defenseOverexertSelected.value.push(unit.instance_id);
                syncDefensePanelVisualState();
                logDefenseSelection('defense_overexert_toggle', unit);
            }

            function toggleDefenseSkill(skill) {
                if (!defenseSelected.value.includes(skill.instance_id)) {
                    defenseSelected.value = [skill.instance_id];
                    defenseScrollSelected.value = [];
                    defenseBoostSelected.value = [];
                    defensePowerSacrificeId.value = '';
                }
                syncDefensePanelVisualState();
                logDefenseSelection('defense_skill_toggle', skill);
            }
            function toggleDefenseScroll(scroll) {
                if (!defenseScrollSelected.value.includes(scroll.instance_id)) {
                    defenseScrollSelected.value = [scroll.instance_id];
                    defenseSelected.value = [];
                    defenseBoostSelected.value = [];
                    defensePowerSacrificeId.value = '';
                }
                syncDefensePanelVisualState();
                logDefenseSelection('defense_scroll_toggle', scroll);
            }

            function toggleDefensePowerSacrifice(unit) {
                defensePowerSacrificeId.value = defensePowerSacrificeId.value === unit.instance_id ? '' : unit.instance_id;
                syncDefensePanelVisualState();
                logDefenseSelection('defense_power_sacrifice_toggle', unit);
            }

            function toggleDefenseBoost(skill) {
                if (isDefenseBoostDisabled(skill)) {
                    logDefenseSelection('defense_boost_blocked', skill);
                    return;
                }
                const idx = defenseBoostSelected.value.indexOf(skill.instance_id);
                if (idx >= 0) {
                    defenseBoostSelected.value.splice(idx, 1);
                } else {
                    defenseBoostSelected.value.push(skill.instance_id);
                }
                syncDefensePanelVisualState();
                logDefenseSelection('defense_boost_toggle', skill);
            }

            function logDefenseSelection(message, card = null) {
                sendClientLog('info', message, {
                    card: card ? { instance_id: card.instance_id, number: card.number, name: card.name } : null,
                    defense_selected: [...defenseSelected.value],
                    defense_scroll_selected: [...defenseScrollSelected.value],
                    defense_boost_selected: [...defenseBoostSelected.value],
                    defense_overexert_selected: [...defenseOverexertSelected.value],
                    defense_power_sacrifice_id: defensePowerSacrificeId.value || '',
                    can_submit_defense: canSubmitDefense.value,
                    total_defense_cost: totalDefenseCost.value,
                    defense_available_elements: defenseAvailableElements.value,
                    total_defense_power: totalDefensePower.value,
                    pending_spell_power: pendingAttackPower(),
                });
            }

            function reactSpell(skill) {
                if (!canUseSkillAsReaction(skill)) {
                    addLog('无法响应该法术');
                    return;
                }
                if (!canPaySelectedReaction(skill)) {
                    addLog('响应费用不足，请先选择透支单位');
                    return;
                }
                sendPaidAction('react_spell', { instance_id: skill.instance_id, overexert_ids: [...defenseOverexertSelected.value] }, reactionActionCost(skill), `响应 ${skill.name}`, null, reactionAvailableElements());
                addLog(`响应 ${skill.name}`);
            }

            const defensePaymentOverride = ref(null);
            const defensePaymentKey = computed(() => JSON.stringify([gameState.value?.revision, totalDefenseCost.value, defenseAvailableElements.value, defenseSelected.value, defenseScrollSelected.value, defenseBoostSelected.value, defenseOverexertSelected.value]));
            const defensePaymentPlan = computed(() => defensePaymentOverride.value?.key === defensePaymentKey.value ? defensePaymentOverride.value.payment : calculatePaymentWithAvailable(defenseAvailableElements.value,totalDefenseCost.value) || {});
            const defensePaymentValid = computed(() => isExactPaymentForCost(defensePaymentPlan.value,totalDefenseCost.value) && ELEMENT_ORDER.every(e=>(defensePaymentPlan.value[e]||0)<=(defenseAvailableElements.value[e]||0)));
            function adjustDefensePayment(elem,delta) {
                const p={...defensePaymentPlan.value};
                p[elem]=Math.max(0,Math.min(defenseAvailableElements.value[elem]||0,(p[elem]||0)+delta));
                defensePaymentOverride.value={key:defensePaymentKey.value,payment:p};
            }
            function resetDefenseDraft() {
                defenseSelected.value=[]; defenseScrollSelected.value=[]; defenseBoostSelected.value=[]; defenseOverexertSelected.value=[]; defensePowerSacrificeId.value=''; defensePaymentOverride.value=null;
            }
            function submitDefend() {
                if (!canSubmitDefense.value || !defensePaymentValid.value) {
                    const hasDefenseSource = defenseSelected.value.length + defenseScrollSelected.value.length > 0;
                    addLog(hasDefenseSource ? '防御元素不足' : '请选择防御法术或防御卷轴');
                    logDefenseSelection('defense_submit_blocked');
                    return;
                }
                const data = {
                    skill_ids: [...defenseSelected.value],
                    scroll_ids: [...defenseScrollSelected.value],
                    boost_ids: [...defenseBoostSelected.value],
                    overexert_ids: [...defenseOverexertSelected.value],
                };
                if (defensePowerSacrificeId.value) {
                    data.sacrifice_id = defensePowerSacrificeId.value;
                }
                const skillCount = defenseSelected.value.length + defenseScrollSelected.value.length;
                const boostCount = defenseBoostSelected.value.length;
                const overexertCount = defenseOverexertSelected.value.length;
                sendAction('defend', {...data,payment:{...defensePaymentPlan.value}},null,true);
            }

            function noDefend() { sendAction('no_defend', {}); }

            // Pending action selection
            const pendingSelected = ref([]);
            const pendingDeckOrder = ref([]);
            const pendingDeckBottom = ref([]);
            const pendingDeckKey = ref('');
            const pendingOverexertSelected = ref([]);

            const pendingAction = computed(() => gameState.value?.pending_action || null);
            function pendingPlayerID(action) {
                if (!action || action.player_id === undefined || action.player_id === null) return null;
                const id = Number(action.player_id);
                return Number.isFinite(id) ? id : null;
            }
            const isOwnPendingAction = computed(() =>
                !isSpectator &&
                pendingAction.value != null &&
                (pendingPlayerID(pendingAction.value) === null || pendingPlayerID(pendingAction.value) === mySlot.value)
            );
            const showPendingAction = computed(() => isOwnPendingAction.value);
            const isCounterPending = () => pendingAction.value?.type === 'counter_trigger';
            const pendingCanOverexert = () => !!pendingAction.value?.can_overexert;

            const pendingOverexertUnits = computed(() => {
                if (!pendingCanOverexert() || !myState.value.units) return [];
                const units = [];
                const selectedCounter = new Set(pendingSelected.value);
                for (let col = 0; col < 3; col++) {
                    for (let row = 0; row < 3; row++) {
                        const unit = getMyUnit(col, row);
                        if (unit && canOverexertForDefense(unit)) units.push(unit);
                    }
                }
                for (const equipment of (myState.value.equipment || [])) {
                    if (equipment && canOverexertForDefense(equipment) && !selectedCounter.has(equipment.instance_id)) units.push(equipment);
                }
                return units;
            });

            function togglePendingOverexert(unit) {
                const idx = pendingOverexertSelected.value.indexOf(unit.instance_id);
                if (idx >= 0) pendingOverexertSelected.value.splice(idx, 1);
                else pendingOverexertSelected.value.push(unit.instance_id);
            }

            function isDeckReorderPending() {
                return pendingAction.value?.ui?.kind === 'reorder';
            }

            function ensurePendingDeckOrder() {
                if (!isDeckReorderPending()) return;
                const ids = (pendingAction.value?.candidates || []).map(candidate => candidate.instance_id);
                const key = `${pendingAction.value?.type || ''}:${ids.join('|')}`;
                if (pendingDeckKey.value === key) return;
                pendingDeckKey.value = key;
                pendingDeckOrder.value = [...ids];
                pendingDeckBottom.value = [];
                pendingSelected.value = pendingSelected.value.filter(id => ids.includes(id));
            }

            function resetPendingSelectionState() {
                pendingSelected.value = [];
                pendingDeckOrder.value = [];
                pendingDeckBottom.value = [];
                pendingDeckKey.value = '';
                pendingOverexertSelected.value = [];
            }

            function pendingCandidateByID(id) {
                return (pendingAction.value?.candidates || []).find(candidate => candidate.instance_id === id);
            }

            function pendingCandidateDetail(candidate) {
                if (!candidate) return '';
                const parts = [];
                if (typeof candidate.power === 'number' && candidate.power >= 0) {
                    const role = typeof candidate.is_main === 'boolean' ? (candidate.is_main ? '主法术 · ' : '强化法术 · ') : '';
                    parts.push(`${role}威力 ${candidate.power}`);
                }
                if (candidate.elements_gain && totalCost(candidate.elements_gain) > 0) {
                    parts.push(`负载 ${costLabel(candidate.elements_gain)}`);
                }
                if (candidate.position && typeof candidate.position.col === 'number' && typeof candidate.position.row === 'number') {
                    parts.push(`位置 ${candidate.position.col + 1}-${candidate.position.row + 1}`);
                }
                if (candidate.side) parts.unshift(candidate.side === 'own' ? '己方' : candidate.side === 'enemy' ? '敌方' : candidate.side);
                else if (candidate.owner !== undefined) parts.unshift(pendingContextPlayerLabel(candidate.owner));
                if (candidate.zone) parts.push(zoneLabel(candidate.zone));
                if (candidate.current_life >= 0) parts.push(`血 ${candidate.current_life}`);
                if (typeof candidate.is_horizontal === 'boolean') parts.push(candidate.is_horizontal ? '横置' : '竖置');
                for (const [name, amount] of Object.entries(candidate.statuses || {})) {
                    if (amount > 0) parts.push(`${name} ${amount}`);
                }
                if (candidate.number && !candidate.position) {
                    const index = (pendingAction.value?.candidates || []).findIndex(c => c.instance_id === candidate.instance_id);
                    if (index >= 0) parts.push(`候选 ${index + 1}`);
                }
                return parts.join(' · ');
            }

            function pendingConfirmLabel() {
                if (!pendingSelected.value.length && pendingAction.value?.min_select === 0 && !isDeckReorderPending()) return '不选择 / 放弃发动';
                return '确认选择';
            }

            function pendingContextPlayerLabel(playerID) {
                if (playerID === undefined || playerID === null) return '';
                const id = Number(playerID);
                if (!Number.isFinite(id)) return '';
                if (isSpectator) return `玩家${id + 1}`;
                if (id === mySlot.value) return '你';
                if (id === 1 - mySlot.value) return '对手';
                return `玩家${id + 1}`;
            }

            function pendingContextCardLabel(card, playerID) {
                if (!card) return '';
                const owner = pendingContextPlayerLabel(playerID);
                const name = card.name || card.number || card.instance_id || '未知卡牌';
                return owner ? `${owner}的${name}` : name;
            }

            function pendingContextText() {
                const ctx = pendingAction.value?.context;
                if (!ctx) return '';
                const parts = [];
                if (ctx.trigger_label) parts.push(`触发:${ctx.trigger_label}`);
                if (ctx.source) parts.push(`来源:${pendingContextCardLabel(ctx.source, ctx.source_player)}`);
                if (ctx.target) parts.push(`目标:${pendingContextCardLabel(ctx.target, ctx.target_player)}`);
                if (ctx.draw_count) parts.push(`抽牌:${ctx.draw_count}`);
                if (ctx.power) parts.push(`威力:${ctx.power}`);
                if (ctx.spell_target && !ctx.target) {
                    const targetType = ctx.spell_target.type || '';
                    if (targetType === 'hero') parts.push('目标:人物');
                    else if (ctx.spell_target.position) {
                        const pos = ctx.spell_target.position;
                        parts.push(`目标:${Number(pos.col) + 1}-${Number(pos.row) + 1}`);
                    }
                }
                return parts.join(' | ');
            }

            function pendingOrderedCandidates() {
                if (!isDeckReorderPending()) return pendingAction.value?.candidates || [];
                ensurePendingDeckOrder();
                const ordered = pendingDeckOrder.value
                    .map(id => pendingCandidateByID(id))
                    .filter(Boolean);
                const seen = new Set(ordered.map(candidate => candidate.instance_id));
                for (const candidate of pendingAction.value?.candidates || []) {
                    if (!seen.has(candidate.instance_id)) ordered.push(candidate);
                }
                return ordered;
            }

            function pendingCandidatePosition(candidate) {
                if (!candidate) return null;
                const raw = candidate.position;
                if (raw && typeof raw === 'object') {
                    const col = Number(raw.col ?? raw.Col);
                    const row = Number(raw.row ?? raw.Row);
                    if (Number.isInteger(col) && Number.isInteger(row)) return { col, row };
                }
                const match = String(candidate.instance_id || '').match(/^pos:(\d+):(\d+)$/);
                if (!match) return null;
                return { col: Number(match[1]), row: Number(match[2]) };
            }

            function isPendingPositionCandidate(candidate) {
                return candidate?.zone === 'unit_position' || !!pendingCandidatePosition(candidate);
            }

            function pendingPositionLabel(candidate) {
                const pos = pendingCandidatePosition(candidate);
                if (!pos) return candidate?.name || '位置';
                const colLabels = ['左列', '中列', '右列'];
                const rowLabels = ['前排', '中排', '后排'];
                const col = colLabels[pos.col] || `${pos.col + 1}列`;
                const row = rowLabels[pos.row] || `${pos.row + 1}排`;
                return `${row}${col}`;
            }

            function pendingPositionCells(candidate) {
                const pos = pendingCandidatePosition(candidate);
                const cells = [];
                for (const row of [2, 1, 0]) {
                    for (const col of [0, 1, 2]) {
                        cells.push({
                            key: `${col}-${row}`,
                            col,
                            row,
                            active: !!pos && pos.col === col && pos.row === row
                        });
                    }
                }
                return cells;
            }

            function isPendingBottom(id) {
                ensurePendingDeckOrder();
                return pendingDeckBottom.value.includes(id);
            }

            function movePendingCandidate(id, direction) {
                ensurePendingDeckOrder();
                const idx = pendingDeckOrder.value.indexOf(id);
                const next = idx + direction;
                if (idx < 0 || next < 0 || next >= pendingDeckOrder.value.length) return;
                const order = [...pendingDeckOrder.value];
                [order[idx], order[next]] = [order[next], order[idx]];
                pendingDeckOrder.value = order;
            }

            function togglePendingBottom(id) {
                ensurePendingDeckOrder();
                if (pendingSelected.value.includes(id)) return;
                const idx = pendingDeckBottom.value.indexOf(id);
                if (idx >= 0) {
                    pendingDeckBottom.value.splice(idx, 1);
                } else {
                    pendingDeckBottom.value.push(id);
                }
            }

            function pendingTopOrder() {
                ensurePendingDeckOrder();
                return pendingDeckOrder.value.filter(id => !pendingSelected.value.includes(id) && !pendingDeckBottom.value.includes(id));
            }

            function pendingBottomOrder() {
                ensurePendingDeckOrder();
                return pendingDeckOrder.value.filter(id => !pendingSelected.value.includes(id) && pendingDeckBottom.value.includes(id));
            }

            function togglePendingSelect(candidate) {
                if (candidate.can_select === false) return;
                const id = candidate.instance_id;
                const idx = pendingSelected.value.indexOf(id);
                if (idx >= 0) {
                    pendingSelected.value.splice(idx, 1);
                } else {
                    if (pendingAction.value?.ui?.unique_by === 'side' && candidate.side) {
                        const existingSameSide = pendingSelected.value.some(selectedID => {
                            const selected = pendingAction.value.candidates.find(c => c.instance_id === selectedID);
                            return selected?.side === candidate.side;
                        });
                        if (existingSameSide) return;
                    }
                    const max = pendingAction.value?.max_select || 1;
                    if (max === 1) {
                        pendingSelected.value = [id];
                    } else if (pendingSelected.value.length < max) {
                        pendingSelected.value.push(id);
                    }
                }
            }

            function submitPendingAction() {
                if (!pendingAction.value) {
                    sendClientLog('warn', 'submit_pending_without_action', {
                        phase: gameState.value?.phase || '',
                        my_slot: mySlot.value,
                    });
                    return;
                }
                sendClientLog('info', 'pending_selection_submitted', { type: pendingAction.value.type, count: pendingSelected.value.length });
                const payload = { selected: [...pendingSelected.value] };
                if (pendingAction.value.ui?.kind === 'load_distribution') payload.load_distribution = Object.fromEntries(pendingSelected.value.map(id => [id, loadDistribution.value[id] || {}]));
                if (pendingAction.value.ui?.kind === 'allocation') payload.allocations = Object.fromEntries(pendingSelected.value.map(id => [id, Number(allocationAmounts.value[id] || 0)]));
                if (pendingAction.value.ui?.kind === 'resource_distribution') {
                    payload.payments = pendingSelected.value.map(id => ({instance_id:id, amount:Number(allocationAmounts.value[id] || 1), mode:resourceModes.value[id] || 'auto'}));
                    const total = payload.payments.reduce((sum,p)=>sum+p.amount,0);
                    if (total > pendingAction.value.ui.max_total || payload.payments.some(p=> !Number.isInteger(p.amount) || p.amount<1)) { decisionError.value = '请检查分配数量，不得超过总上限'; return; }
                    for (const p of payload.payments) {
                        const c = pendingAction.value.candidates.find(c=>c.instance_id===p.instance_id);
                        const capacity = Math.max(0,Number(c.current_life)-1) + (p.mode==='life' ? 0 : totalCost(c.elements_gain));
                        if (p.amount > capacity) { decisionError.value = '所选单位没有足够的可支付负载或生命'; return; }
                    }
                }
                if (isDeckReorderPending()) {
                    payload.top_order = pendingTopOrder();
                    payload.bottom_order = pendingBottomOrder();
                }
                if (pendingCanOverexert()) {
                    payload.overexert_ids = [...pendingOverexertSelected.value];
                }
                const cost = pendingSelected.value.length ? (pendingAction.value.cost || {}) : {};
                const available = cloneElements(myState.value.elements || {});
                for (const id of pendingOverexertSelected.value) {
                    const unit = pendingOverexertUnits.value.find(u => u.instance_id === id);
                    for (const [elem, value] of Object.entries(unit?.elements_gain || {})) available[elem] = (available[elem] || 0) + value;
                }
                sendPaidAction('resolve_action', payload, cost, pendingAction.value.prompt, null, available);
                // Keep the draft until an authoritative state transition.
            }

            function pendingRequiresPaymentChoice() {
                return !!pendingAction.value?.cost && totalCost(pendingAction.value.cost) > 0 && !pendingCanOverexert();
            }

            function zoneLabel(zone) {
                return {
                    hand: '手牌',
                    equipment: '装备',
                    skill: '技能',
                    unit: '战场',
                    deck: '卡组',
                }[zone] || zone;
            }

            // Ability functions
            function canUsePerTurn(unit) {
                if (!unit || !unit.has_per_turn) return false;
                if (unit.statuses?.['石化'] > 0 || unit.statuses?.['眩晕'] > 0) return false;
                const limit = unit.per_turn_limit || 1;
                return (unit.used_this_turn || 0) < limit;
            }

            function canUseUltimate(unit) {
                if (!unit || !unit.has_ultimate) return false;
                if (unit.statuses?.['石化'] > 0 || unit.statuses?.['眩晕'] > 0) return false;
                return !unit.ultimate_used;
            }

            function abilityLabel(unit, abilityType) {
                if (abilityType === 'ultimate') return '绝技';
                return unit?.per_turn_label || '回合技';
            }

            function compactAbilityLabel(unit, abilityType) {
                const label = abilityLabel(unit, abilityType);

                if (label.length > 4) return '技能';
                return label;
            }

            function useAbility(unit, abilityType) {
                sendAction('use_ability', {
                    instance_id: unit.instance_id,
                    ability_type: abilityType
                });
                addLog(`${unit.name} 使用${abilityLabel(unit, abilityType)}`);
            }

            function castDetailSkill(skill) {
                if (!canCastSkill(skill)) return;
                selectSkillToCast(skill);
                detailCard.value = null;
            }

            function useDetailAbility(unit, abilityType) {
                if (abilityType === 'per_turn' && !canUsePerTurn(unit)) return;
                if (abilityType === 'ultimate' && !canUseUltimate(unit)) return;
                useAbility(unit, abilityType);
                detailCard.value = null;
            }

            function learnSkill(card) {
                if (!isMyTurn.value || gameState.value.phase !== 'main') return;
                if (!canLearnSkill(card)) {
                    addLog(`${card.name} ${missingCostText(learnCost(card))}`);
                    return;
                }
                if (!mySkills.value.some(skill => !skill)) {
                    skillReplaceRequest.value = { card, replaceId: '' };
                    return;
                }
                sendLearnSkill(card, '');
            }
            function canLearnSkill(card) {
                const hasEmptySkillSlot = mySkills.value.some(skill => !skill);
                const canReplace = replaceableSkills.value.length > 0;
                return !!card && (hasEmptySkillSlot || canReplace) && isMyTurn.value && gameState.value.phase === 'main' && canPayCost(learnCost(card));
            }
            function sendLearnSkill(card, replaceId) {
                sendPaidAction('learn_skill', { instance_id: card.instance_id, replace_id: replaceId || '' }, learnCost(card), `学习 ${card.name}`, () => {
                    addLog(replaceId ? `学习技能 ${card.name}，替代旧技能` : `学习技能 ${card.name}`);
                    skillReplaceRequest.value = null;
                });
            }
            function confirmSkillReplacement() {
                if (!skillReplaceRequest.value?.card || !skillReplaceRequest.value.replaceId) return;
                sendLearnSkill(skillReplaceRequest.value.card, skillReplaceRequest.value.replaceId);
            }
            function canPlayHandCard(card) {
                if (!card || !isMyTurn.value || gameState.value.phase !== 'main') return false;
                if (card.type === '技能') return false;
                if (card.unavailable_reason) return false;
                if (card.is_counter_trap) return canSetCounterCard(card);
                if (isEquippableItem(card)) return canEquipCard(card);
                return canPayCost(playCost(card));
            }
            function handCardUnavailableText(card) {
                if (!card) return '无法使用';
                if (card.unavailable_reason) return card.unavailable_reason;
                if (card.is_counter_trap && !hasEmptyEquipmentSlot()) return `${card.name} 无法盖放，装备区已满`;
                if (isEquippableItem(card) && !hasEmptyEquipmentSlot() && replaceableEquipment.value.length === 0) {
                    return `${card.name} 无法装备，装备区已满且没有竖置装备可替换`;
                }
                return `${card.name} ${missingCostText(playCost(card))}`;
            }
            function endTurn() {
                if (!isMyTurn.value || gameState.value?.phase !== 'main') {
                    sendClientLog('warn', 'blocked_end_turn_not_main', {
                        phase: gameState.value?.phase || '',
                        my_slot: mySlot.value,
                        has_pending: !!pendingAction.value,
                        pending_type: pendingAction.value?.type || '',
                        pending_player_id: pendingPlayerID(pendingAction.value),
                    });
                    addLog(waitingNotice.value || '当前不能结束回合');
                    return;
                }
                sendAction('end_turn', {});
                selectedHandCard.value = null; selectedAttacker.value = null; selectedSkill.value = null; selectedItemSpellScroll.value = null; selectedSpellScrollSacrificeId.value = ''; selectedSpellScrollSupportId.value = ''; selectedSpellPowerSacrificeId.value = ''; hoveredOpponentCell.value = null;
                addLog('结束回合');
            }

            const guideKicker = computed(() => {
                if (!gameState.value) return '连接';
                if (isSpectator) return '观战';
                if (gameState.value.phase === 'mulligan') return '起手';
                if (showDefenseWindow.value) return '防御窗口';
                if (showPendingAction.value) return '选择';
                if (gameState.value.phase === 'game_over') return '终局';
                return isMyTurn.value ? '你的回合' : '对手回合';
            });

            const waitingNotice = computed(() => {
                if (!gameState.value) return '';
                if (isSpectator) {
                    if (gameState.value.phase === 'waiting_action' && pendingAction.value) {
                        return `等待玩家${pendingPlayerID(pendingAction.value) + 1}完成选择...`;
                    }
                    if (gameState.value.phase === 'defense_window' && gameState.value.pending_spell) {
                        return `玩家${pendingSpellAttackerID() === 0 ? 2 : 1}正在防御...`;
                    }
                    if (gameState.value.phase === 'mulligan') return '双方正在确认起手...';
                    if (gameState.value.phase === 'main') return `玩家${gameState.value.current_turn + 1}回合进行中...`;
                    return '观战中。';
                }
                if (gameState.value.phase === 'mulligan' && mulliganSubmitted.value) return '等待对手确认起手...';
                if (showPendingAction.value) return pendingAction.value?.prompt || '等待你完成选择...';
                if (gameState.value.phase === 'waiting_action') return '等待对手完成选择...';
                if (showDefenseWindow.value) return '对手施放法术，等待你选择防御方式。';
                if (gameState.value.phase === 'defense_window' && pendingSpellAttackerID() === mySlot.value) return '等待对手选择防御或不防御...';
                if (gameState.value.phase === 'main' && !isMyTurn.value) return '等待对手行动...';
                return '';
            });

            const currentHint = computed(() => {
                if (!gameState.value) return '等待游戏状态同步。';
                if (isSpectator) return '观战模式会隐藏双方手牌、技能池和卡组顺序，只显示公开信息。';
                if (gameState.value.phase === 'mulligan') return '查看起手牌，选择保留或全部重抽。';
                if (showDefenseWindow.value) return '选择防御技能，或直接不防御。强化技能会叠加防御威力。';
                if (showPendingAction.value) return pendingAction.value?.prompt || '请选择要处理的卡牌。';
                if (gameState.value.phase === 'waiting_action') return '等待对手完成选择。';
                if (gameState.value.phase === 'game_over') return '对局结束，可以返回大厅重新开局。';
                if (!isMyTurn.value) return '等待对手行动。';
                if (selectedItemSpellScroll.value) {
                    const name = selectedItemSpellScroll.value.name;
                    if (selectedItemSpellScroll.value.preparation === 'fire_sacrifice' && !selectedSpellScrollSacrificeId.value) return `先选择友方火焰伙伴祭献，再选择目标使用 ${name}。`;
                    if (selectedItemSpellScroll.value.preparation === 'companion_support' && !selectedSpellScrollSupportId.value) return `先选择生命值+负载大于5的友方伙伴支援，再选择目标使用 ${name}。`;
                    return `选择敌方目标使用 ${name}。`;
                }
                if (selectedHandCard.value) {
                    const name = selectedHandCard.value.name;
                    if (!canPayCost(selectedHandCard.value.elements_cost || {})) return `${name} 元素不足，先消耗竖置卡牌获得元素。`;
                    if (selectedHandCard.value.type === '伙伴') return `选择己方空格召唤 ${name}。`;
                    if (selectedHandCard.value.is_terrain) return `选择己方空格放置地形 ${name}。`;
                    return `再次点击或右键使用 ${name}。`;
                }
                if (selectedAttacker.value) return `选择敌方前排目标，让 ${selectedAttacker.value.name} 攻击。`;
                if (pendingExtraTargetCast.value && selectedSkillNeedsExtraTarget(selectedSkill.value)) return `选择 ${selectedSkill.value.name} 的额外目标。`;
                if (selectedSkill.value && selectedSkillPowerSacrificeSource() && !selectedSpellPowerSacrificeId.value) return `可先选择友方火焰伙伴祭献，或直接选择目标施放 ${selectedSkill.value.name}。`;
                if (selectedSkill.value) return `选择敌方目标施放 ${selectedSkill.value.name}。`;
                return '消耗竖置卡牌获得元素，然后召唤、学习技能、施法或攻击。';
            });

            function addLog(msg) {
                actionLog.value.unshift(`[${new Date().toLocaleTimeString()}] ${msg}`);
                if (actionLog.value.length > 50) actionLog.value.pop();
            }

            async function postTestRoomTool(path, payload) {
                const resp = await fetch(path, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(payload)
                });
                const data = await resp.json().catch(() => ({}));
                if (!resp.ok) throw new Error(data.error || `HTTP ${resp.status}`);
                return data;
            }

            async function testAddCard(zone) {
                const cardNumber = testCardNumber.value.trim();
                if (!cardNumber) {
                    addLog('测试工具：请输入卡号');
                    return;
                }
                try {
                    await postTestRoomTool('/api/test-room/add-card', {
                        room_id: roomID,
                        player: mySlot.value,
                        card_number: cardNumber,
                        zone,
                        count: 1
                    });
                    addLog(`测试工具：已添加 ${cardNumber} -> ${zone}`);
                } catch (err) {
                    addLog(`测试工具错误：${err.message}`);
                }
            }

            async function testAddElements() {
                const amount = Number(testElementAmount.value || 0);
                if (!amount) {
                    addLog('测试工具：元素调整值不能为 0');
                    return;
                }
                try {
                    await postTestRoomTool('/api/test-room/elements', {
                        room_id: roomID,
                        player: mySlot.value,
                        elements: { [testElement.value]: amount },
                        mode: 'add'
                    });
                    addLog(`测试工具：${elemFull(testElement.value)} ${amount > 0 ? '+' : ''}${amount}`);
                } catch (err) {
                    addLog(`测试工具错误：${err.message}`);
                }
            }

            // WebSocket
            function connect() {
                const wsProto = location.protocol === 'https:' ? 'wss:' : 'ws:';
                const wsParams = new URLSearchParams({
                    room: roomID,
                    player_id: playerID,
                    player_name: playerName || '观战者',
                });
                if (isSpectator) {
                    wsParams.set('role', 'spectator');
                } else {
                    wsParams.set('deck_code', deckCode);
                }
                const wsUrl = `${wsProto}//${location.host}/ws?${wsParams.toString()}`;
                ws = new WebSocket(wsUrl);
                ws.onopen = () => {
                    connected.value = true;
                    connectionStatus.value = '已连接';
                    addLog(isSpectator ? '已进入观战' : '已连接');
                    sendClientLog('info', 'websocket_open');
                };
                ws.onmessage = (e) => {
                    try {
                        handleMessage(JSON.parse(e.data));
                    } catch (err) {
                        sendClientLog('error', 'message_parse_or_handle_failed', { error: String(err), raw: e.data });
                        throw err;
                    }
                };
                ws.onclose = () => {
                    connected.value = false;
                    connectionStatus.value = '连接断开';
                    session.disconnect();
                    actionBusy.value = false;
                    decisionError.value = '连接已断开。结果未知的动作不会自动重发；重新连接后请核对战场。';
                };
                ws.onerror = () => {
                    error.value = '连接错误';
                    sendClientLog('error', 'websocket_error');
                };
            }

            function handleMessage(msg) {
                if (msg.type === 'joined') {
                    mySlot.value = msg.data.slot;
                    isTestRoom.value = !!msg.data.room?.test_mode;
                    if (isSpectator || msg.data.role === 'spectator') {
                        mySlot.value = -1;
                        addLog(`观战房间 #${msg.data.room_id}`);
                    } else if (msg.data.reconnected) {
                        mulliganSubmitted.value = true;
                        addLog('重新连接成功');
                    } else {
                        addLog(`加入房间 #${msg.data.room_id}`);
                    }
                }
                else if (msg.type === 'game_event') handleGameEvent(msg.event);
                else if (msg.type === 'action_result') {
                    if (!msg.ok && msg.message) addLog(`操作失败: ${msg.message}`);
                    const outcome = session.result(msg);
                    actionBusy.value = session.busy;
                    if (outcome.startsWith('error')) decisionError.value = msg.message || '操作失败，请编辑后重试';
                    if (outcome === 'success') decisionError.value = '';
                }
                else if (msg.type === 'error') { decisionError.value = msg.message; addLog(`错误: ${msg.message}`); }
            }

            function handleGameEvent(evt) {
                const d = evt.data || {};
                if (evt.player === -1 && evt.type !== 'state_sync') {
                    const labels = {game_setup:'对局准备', game_start:'对局开始', turn_start:'回合开始', turn_end:'回合结束', summon:'召唤', equip:'装备', learn_skill:'学习技能', counter_set:'盖放反制', counter_revealed:'揭示反制', spell_cast:'声明法术', spell_reaction:'法术响应', defense_success:'防御成功', defense_attempt:'尝试防御', spell_hit:'法术命中', unit_attack:'单位攻击', unit_destroyed:'单位死亡', unit_returned:'返回手牌', unit_moved:'单位移动', ability_used:'发动能力', place_terrain:'放置地形', discard:'弃置', use_item:'使用道具', consume:'消耗', game_over:'对局结束', draw_offer:'请求和局', draw_offer_declined:'拒绝和局', effect_trigger:'规则效果'};
                    const name = d.card?.name || d.skill?.name || d.source?.name || '';
                    const detail = typeof d.damage === 'number' ? ` · 伤害 ${d.damage}` : typeof d.power === 'number' ? ` · 威力 ${d.power}` : '';
                    publicEventLog.value.unshift(`${labels[evt.type] || '规则结算'}${name ? ' · '+name : ''}${detail}`);
                    publicEventLog.value = publicEventLog.value.slice(0,100);
                }
                switch (evt.type) {
                    case 'state_sync':
                        if (session.revision !== null && d.revision < session.revision) break;
                        acceptSnapshot(d);
                        gameState.value = d;
                        if (detailCard.value?.instance_id) detailCard.value = findVisibleCard(detailCard.value.instance_id) || null;
                        if (hoverCard.value?.instance_id) hoverCard.value = findVisibleCard(hoverCard.value.instance_id) || null;
                        logDefenseWindowStateSync(gameState.value);
                        break;
                    case 'game_setup': addLog('游戏准备完成'); break;
                    case 'game_start': addLog('游戏开始！'); break;
                    case 'turn_start': addLog(`回合 ${d.turn_number} — 玩家${d.current_player+1}`); break;
                    case 'spell_cast': addLog(`玩家${d.attacker+1} 施放 ${d.skill.name} 威力:${d.power}`); break;
                    case 'spell_reaction': addLog(`${d.card?.name || '响应咒术'} 使法术威力变为 ${d.power}`); break;
                    case 'defense_success': addLog('防御成功！'); break;
                    case 'spell_hit': addLog(`法术命中 ${d.damage} 伤害`); break;
                    case 'unit_attack': addLog(`${d.attacker.name} 攻击 ${d.damage} 伤害`); break;
                    case 'unit_moved': addLog(`${d.card.name} 移动到 (${d.to.col},${d.to.row})`); break;
                    case 'unit_destroyed': addLog(`${d.card.name} 被摧毁`); break;
                    case 'unit_returned': addLog(`${d.card.name} 返回手牌`); break;
                    case 'game_over':
                        if (d.reason === 'surrender') addLog(`玩家${d.actor+1} 投降，玩家${d.winner+1} 获胜`);
                        else if (d.reason === 'draw_agreement') addLog('双方同意和局');
                        else addLog(d.winner === -2 ? '游戏结束 平局' : `游戏结束 玩家${d.winner+1} 获胜`);
                        break;
                    case 'draw_offer': addLog(`玩家${d.player+1} 请求和局`); break;
                    case 'draw_offer_declined': addLog(`玩家${d.player+1} 拒绝和局`); break;
                    case 'draw_card': addLog(`抽到 ${d.card.name}`); break;
                    case 'effect_trigger':
                        if (d.effect === 'charge') addLog(`充能+${d.amount}`);
                        else if (d.effect === 'damage') addLog(`${d.source?.name} 造成${d.amount}伤害`);
                        else if (d.effect === 'apply_status') addLog(`${d.target?.name} ${d.status}+${d.amount}`);
                        else if (d.effect === 'gain_shield') addLog(`${d.source?.name} 获得护盾${d.amount}`);
                        else if (d.effect === 'gain_stealth') addLog(`${d.source?.name} 获得隐蔽`);
                        else if (d.effect === 'heal') addLog(`${d.target?.name} 回复${d.amount}生命`);
                        else if (d.effect === 'skip_draw') addLog('跳过本回合抽牌');
                        else if (d.effect === 'temp_modifier') addLog('获得临时效果');
                        break;
                    case 'ability_used': addLog(`${d.card.name} 使用${abilityLabel(d.card, d.ability)}`); break;
                    case 'defense_attempt': addLog(`防御 威力${d.defense_power} vs ${d.attack_power}`); break;
                    case 'pending_action':
                        // Wait for the authoritative filtered snapshot.
                        break;
                    case 'pending_action_cleared':
                        // The next snapshot identifies the new decision.
                        break;
                    case 'discard':
                        addLog(`弃置 ${d.card.name}`);
                        flashGraveyardCard(d.player, d.card);
                        break;
                    case 'use_item':
                        addLog((d.card && d.card.name ? d.card.name : 'item') + ' to graveyard');
                        flashGraveyardCard(d.player, d.card);
                        break;
                    case 'search_card': addLog(`检索 ${d.card.name}`); break;
                    case 'place_terrain': addLog(`放置地形 ${d.card.name}`); break;
                    case 'consume': addLog(`消耗 ${d.instance_id}`); break;
                }
            }

            function flashGraveyardCard(player, card) {
                if (!card) return;
                const owner = card.owner ?? player;
                const side = Number(owner) === mySlot.value ? 'me' : 'opponent';
                graveyardFlash.value = { ...graveyardFlash.value, [side]: card };
                recentDiscardCard.value = card;
                clearTimeout(graveyardFlashTimer);
                clearTimeout(recentDiscardTimer);
                graveyardFlashTimer = setTimeout(() => {
                    graveyardFlash.value = { me: null, opponent: null };
                }, 1300);
                recentDiscardTimer = setTimeout(() => {
                    if (!hoverCard.value) recentDiscardCard.value = null;
                }, 2200);
            }
            function logDefenseWindowStateSync(state) {
                if (!state || state.phase !== 'defense_window') return;
                const spell = normalizePendingSpell(state.pending_spell);
                sendClientLog('info', 'defense_window_state_sync', {
                    my_slot: currentPlayerSlot(),
                    attacker_id: spell?.attacker_id ?? null,
                    defender_id: spell?.defender_id ?? null,
                    power: spell?.power ?? null,
                    has_spell: !!spell,
                    show_defense_window: showDefenseWindow.value,
                    is_spectator: isSpectator,
                });
                nextTick(() => {
                    sendClientLog('info', 'defense_window_dom_check', {
                        my_slot: currentPlayerSlot(),
                        show_defense_window: showDefenseWindow.value,
                        overlay_exists: !!document.querySelector('.defense-overlay'),
                        panel_exists: !!document.querySelector('.defense-panel'),
                    });
                });
            }

            function onWindowError(evt) {
                sendClientLog('error', 'window_error', {
                    error: evt?.message || '',
                    source: evt?.filename || '',
                    line: evt?.lineno || 0,
                    column: evt?.colno || 0,
                    stack: evt?.error?.stack || ''
                });
            }

            function onUnhandledRejection(evt) {
                const reason = evt?.reason;
                sendClientLog('error', 'unhandled_rejection', {
                    error: String(reason),
                    stack: reason?.stack || ''
                });
            }

            onMounted(() => {
                if (typeof window !== 'undefined') {
                    window.addEventListener('error', onWindowError);
                    window.addEventListener('unhandledrejection', onUnhandledRejection);
                    window.__arcaneDebug = {
                        actionBusy, actionDraft,
                        gameState,
                        mySlot,
                        pendingAction,
                        showPendingAction,
                        incomingDrawOffer,
                        outgoingDrawOffer,
                        pendingSpellView,
                        showDefenseWindow,
                        isMyTurn,
                        myHand,
                        mySkills,
                        hoverCard,
                        pileView,
                        paymentRequest,
                        devourRequest,
                        skillReplaceRequest,
                        equipmentReplaceRequest,
                        skillNeedsTarget,
                        summonDevourRequirement,
                        devourCandidates,
                        selectedSkill,
                        selectedSpellScrollSupportId,
                        selectedSpellPowerSacrificeId,
                        pendingExtraTargetCast,
                        beginSummonDevour,
                        showCardDetail,
                        sendAction: (action,data,callbacks) => sendAction(action,data,callbacks,false),
                        defenseSelected,
                        defenseScrollSelected,
                        defenseBoostSelected,
                        defenseOverexertSelected,
                        defensePowerSacrificeId,
                        defensePowerSacrificeCandidates,
                        canSubmitDefense,
                        availableDefenseBoostCandidates,
                        syncDefensePanelVisualState,
                    };
                }
                if (roomID && playerID && (isSpectator || deckCode)) connect();
            });
            onUnmounted(() => {
                if (typeof window !== 'undefined') {
                    window.removeEventListener('error', onWindowError);
                    window.removeEventListener('unhandledrejection', onUnhandledRejection);
                }
                if (ws) ws.close();
            });

            return {
                fanPage, fanSpread, fanLayout, fanPreviewID, fanPageCount, currentFanPage, fanCards, fanPreviewCard, showFanPreview, previewHand, changeFanPage, activateFanCard, navigateFan, selectedBattleCard,
                publicEventLog, resourceModes, skillNeedsTarget, castSpellWithoutTarget, railTab, reloadMatch, allocationAmounts, imageFallback, actionDraft, actionBusy, decisionError, workspaceMode, decisionOwnerLabel, draftTitle, draftSummary, commitDraft, cancelDraft, guardInteraction, legacyURL, loadDistribution, setLoadDistribution, paymentAvailable, removePayment,
                connected, connectionStatus, error,
                gameState, mySlot, roomID, isSpectator, mulliganSubmitted, detailCard, hoverCard, recentDiscardCard, graveyardFlash, activePreviewCard,
                selectedHandCard, selectedAttacker, selectedSkill, selectedItemSpellScroll, selectedSpellScrollSacrificeId, selectedSpellScrollSupportId, selectedSpellPowerSacrificeId, selectedBoostSkills, pileView, actionLog,
                isTestRoom, testCardNumber, testElement, testElementAmount, ELEMENT_ORDER,
                myState, opponent, myHand, opponentHandCards, mySkills, myEquipment, mySkillSlotCount, myEquipmentSlotCount, opponentSkillSlotCount, opponentEquipmentSlotCount, allMySkills, allMyBoundSkills,
                isMyTurn, showDefenseWindow, pendingSpellView, myFrontRow, opponentFrontRow,
                incomingDrawOffer, outgoingDrawOffer, canUseMatchActions,
                getCardImage, getElementIcon, getMyUnit, getOpponentUnit, getMyTerrain, getOpponentTerrain, elemFull, costHTML, descriptionHTML, hasCost, hasUseCost, hasGain, playCost, learnCost, entryCostLabel, effectLabels,
                pendingSpellTargetLabel, pendingSpellExtraTargetLabels, pendingSpellSources,
                showCardDetail, previewCard, mulligan, selectHandCard, isConsumable, canSummonAt,
                handleCellClick, canConsumeUnit, consumeCard,
                handleUnitImageClick,
                canAttackWith, selectAttacker, isAttackTarget,
                abilityLabel, compactAbilityLabel, castDetailSkill, useDetailAbility,
                canCastSkill, selectSkillToCast, canBoostWith, toggleBoostSkill, noDefend,
                guideKicker, currentHint, waitingNotice, clearSelection, statusEntries,
                hasActiveTargeting, activeTargetSpell, isValidOpponentTarget, isSpellTarget, isFriendlySpellTarget, isFriendlyHeroSpellTarget, handleMyHeroClick, canUseMyCell,
                myBoardOwnerID, opponentBoardOwnerID, isPendingSpellCenter, isPendingSpellAffected,
                canShowCostHint, missingCostText, skillUseCost, attackActionCost, reactionActionCost, canLearnSkill, costLabel,
                skillReplaceRequest, replaceableSkills, confirmSkillReplacement,
                equipmentReplaceRequest, replaceableEquipment, confirmEquipmentReplacement,
                paymentRequest, paymentSelection, paymentElements, selectedPaymentTotal,
                canAddPayment, addPayment, isPaymentSelectionValid, paymentRemainingText, confirmPayment, cancelPayment, totalCost,
                consumeGainRequest, consumeGainSelection, consumeGainElements, selectedConsumeGainTotal,
                canAddConsumeGain, addConsumeGain, isConsumeGainValid, consumeGainRemainingText, confirmConsumeGain, cancelConsumeGain,
                devourRequest, devourSelection, devourSelectedLoad, devourSelectionMeetsRequirement,
                devourCandidates, toggleDevourSelection, confirmDevourSelection, cancelDevourSelection,
                canSubmitDefense,
                testAddCard, testAddElements,
                canPlayHandCard, canPayCost,
                // Defense
                defenseSelected, defenseScrollSelected, defenseBoostSelected, defenseOverexertSelected, defensePowerSacrificeId, availableDefenseSkills, availableDefenseScrolls, availableReactionSkills, availableDefenseBoostCandidates, availableBoostSkills, defenseOverexertUnits, defensePowerSacrificeCandidates,
                totalDefensePower, totalDefenseCost, defenseAvailableElements, pendingAttackPower, defenseSummaryText, defenseSummaryClass, defensePower, defenseBoostPower, defensePaymentCost, isDefenseBoostDisabled, handleDefensePanelClick, showInteractionCardContextMenu, toggleDefenseSkill, toggleDefenseScroll, toggleDefenseBoost, toggleDefenseOverexert, toggleDefensePowerSacrifice, reactSpell, submitDefend,
                hoverOpponentCell, clearHoverOpponentCell, isSpellAreaPreview,
                // Pending action
                pendingAction, showPendingAction, pendingSelected, togglePendingSelect, submitPendingAction, pendingConfirmLabel, pendingCandidateDetail, pendingContextText,
                isPendingPositionCandidate, pendingPositionLabel, pendingPositionCells,
                pendingCanOverexert, isCounterPending, pendingOverexertUnits, pendingOverexertSelected, togglePendingOverexert,
                isDeckReorderPending, pendingOrderedCandidates, movePendingCandidate, togglePendingBottom, isPendingBottom,
                zoneLabel,
                openPileView, closePileView,
                // Abilities
                canUsePerTurn, canUseUltimate, useAbility, unitActionChips, skillActionChips, cardAvailabilityText,
                learnSkill, endTurn, surrender, offerDraw, respondDrawOffer,
                // Context menu
                menuItems, menuKeydown, previewMetadata, defensePaymentPlan, defensePaymentValid, adjustDefensePayment, resetDefenseDraft, ELEMENT_ORDER,
                contextMenu, showCardContextMenu, hideContextMenu, contextMenuAction
            };
        }
    }).mount('#app');
