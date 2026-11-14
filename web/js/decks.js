(function () {
  const { createApp, ref, computed, watch, onMounted, onUnmounted } = Vue,
    UI = ArcaneUI,
    DS = ArcaneDecks;
  createApp({
    setup() {
      const cards = ref({}),
        decks = ref([]),
        storageError = ref(""),
        loading = ref(true),
        loadError = ref(""),
        query = ref(""),
        typeFilter = ref("全部"),
        versionFilter = ref(""),
        elementFilter = ref(""),
        page = ref(0),
        selected = ref(null),
        fullCard = ref(null),
        showImport = ref(new URLSearchParams(location.search).has("import")),
        importCode = ref(""),
        importError = ref("");
      const editingId = ref(""),
        deckName = ref("未命名卡组"),
        draft = ref({ hero: "", main: [], skills: [], extra: [] }),
        draftTab = ref("main"),
        error = ref(""),
        busy = ref(false),
        validCode = ref(""),
        savedSignature = ref("");
      const allCards = computed(() =>
        Object.values(cards.value).sort((a, b) =>
          a.number.localeCompare(b.number),
        ),
      );
      const packs = computed(() => [
        ...new Set(allCards.value.map((c) => c.version_name)),
      ]);
      const elements = ["火", "水", "气", "地", "光", "暗", "无"];
      const filtered = computed(() => {
        const q = query.value.trim().toLowerCase();
        return allCards.value.filter(
          (c) =>
            (typeFilter.value === "全部" || c.type === typeFilter.value) &&
            (!versionFilter.value || c.version_name === versionFilter.value) &&
            (!elementFilter.value || c.category === elementFilter.value) &&
            (!q ||
              [c.number, c.name, c.tag].some((v) =>
                String(v || "")
                  .toLowerCase()
                  .includes(q),
              )),
        );
      });
      const pageCount = computed(() =>
          Math.max(1, Math.ceil(filtered.value.length / 40)),
        ),
        effectivePage = computed(() =>
          Math.min(page.value, pageCount.value - 1),
        ),
        pageCards = computed(() =>
          filtered.value.slice(
            effectivePage.value * 40,
            (effectivePage.value + 1) * 40,
          ),
        );
      const deckCode = computed(() => DS.encode(draft.value)),
        draftHero = computed(() => cards.value[draft.value.hero]),
        draftRows = computed(() =>
          DS.groups(draft.value[draftTab.value]).map((row) => ({
            ...row,
            card: cards.value[row.number],
          })),
        ),
        dirty = computed(
          () =>
            JSON.stringify([deckName.value, deckCode.value]) !==
            savedSignature.value,
        );
      watch([query, typeFilter, versionFilter, elementFilter], () => {
        page.value = 0;
      });
      function signature() {
        savedSignature.value = JSON.stringify([deckName.value, deckCode.value]);
      }
      function reloadStore() {
        try {
          decks.value = DS.load().decks;
          storageError.value = "";
        } catch (e) {
          storageError.value = e.message;
        }
      }
      function discardChanges() {
        return !dirty.value || confirm("当前构筑尚未保存，放弃修改？");
      }
      function openDeck(deck) {
        if (busy.value) return;
        if (!discardChanges()) return;
        try {
          draft.value = DS.parse(deck.code);
          deckName.value = deck.name;
          editingId.value = deck.id;
          validCode.value = "";
          error.value = "";
          signature();
        } catch (e) {
          error.value = e.message;
        }
      }
      function newDeck() {
        if (busy.value) return;
        if (!discardChanges()) return;
        draft.value = { hero: "", main: [], skills: [], extra: [] };
        deckName.value = "未命名卡组";
        editingId.value = "";
        validCode.value = "";
        error.value = "";
        signature();
        typeFilter.value = "人物";
        selected.value = null;
      }
      function heroFor(deck) {
        try {
          return cards.value[DS.parse(deck.code).hero];
        } catch {
          return null;
        }
      }
      function deckSummary(deck) {
        try {
          const d = DS.parse(deck.code);
          return `${d.main.length}/30 · ${d.skills.length}/10`;
        } catch {
          return "格式异常";
        }
      }
      function tags(card) {
        return String(card?.tag || "").split(/[-,，、\s]+/);
      }
      function countInDeck(card) {
        if (card?.type === "人物")
          return draft.value.hero === card.number ? 1 : 0;
        return [...draft.value.main, ...draft.value.skills].filter(
          (n) => n === card?.number,
        ).length;
      }
      function addReason(card) {
        if (!card) return "未选择卡牌";
        if (tags(card).includes("衍生")) return "衍生牌不可构筑";
        if (card.type === "人物")
          return draft.value.hero === card.number ? "已选人物" : "";
        if (card.type === "技能") {
          if (draft.value.skills.includes(card.number)) return "技能已加入";
          if (draft.value.skills.length >= 10) return "技能已满";
          if (
            tags(card).includes("传奇") &&
            draft.value.skills.some((n) =>
              tags(cards.value[n]).includes("传奇"),
            )
          )
            return "限一张传奇技能";
          return "";
        }
        if (!["伙伴", "道具"].includes(card.type)) return "不可构筑";
        if (draft.value.main.length >= 30) return "主卡组已满";
        if (countInDeck(card) >= (tags(card).includes("传奇") ? 1 : 2))
          return "达到同名上限";
        return "";
      }
      function addCard(card) {
        const reason = addReason(card);
        if (reason) {
          UI.toast(reason);
          return;
        }
        if (card.type === "人物") draft.value.hero = card.number;
        else {
          const zone = card.type === "技能" ? "skills" : "main";
          draft.value[zone].push(card.number);
          draftTab.value = zone;
        }
      }
      function removeCard(number) {
        const zone = draft.value[draftTab.value];
        const i = zone.indexOf(number);
        if (i >= 0) zone.splice(i, 1);
      }
      function loadSample() {
        if (busy.value) return;
        if (!discardChanges()) return;
        draft.value = DS.parse(DS.SAMPLE);
        deckName.value = "大气 · 基础示例";
        editingId.value = "";
        validCode.value = "";
        error.value = "";
      }
      function importDeck() {
        if (busy.value) return;
        importError.value = "";
        try {
          const next = DS.parse(importCode.value);
          if (!discardChanges()) return;
          draft.value = next;
          editingId.value = "";
          deckName.value = "导入卡组";
          validCode.value = "";
          showImport.value = false;
          selected.value = cards.value[next.hero] || null;
          UI.toast("卡牌已读取，可继续编辑");
        } catch (e) {
          importError.value = e.message;
        }
      }
      async function validateCode(code) {
        await UI.request("/api/deck/validate", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ deck_code: code }),
        });
        validCode.value = code;
      }
      async function saveDeck(asDraft = false) {
        if (busy.value || storageError.value) return false;
        busy.value = true;
        error.value = "";
        const code = deckCode.value,
          name = deckName.value,
          id = editingId.value;
        try {
          if (!asDraft) await validateCode(code);
          const data = DS.put(DS.load(), { id, name, code });
          DS.write(data);
          editingId.value = data.selectedId;
          decks.value = data.decks;
          if (code !== deckCode.value || name !== deckName.value) {
            UI.toast("提交时的版本已保存；当前修改尚未保存");
            return false;
          }
          signature();
          UI.toast(asDraft ? "草稿已保存" : "卡组已校验并保存");
          return true;
        } catch (e) {
          error.value = e.message;
          return false;
        } finally {
          busy.value = false;
        }
      }
      async function useDeck() {
        if (await saveDeck(false)) location.href = "/";
      }
      function deleteDeck() {
        if (busy.value) return;
        if (!editingId.value || !confirm("删除这套本地卡组？")) return;
        try {
          const data = DS.load();
          data.decks = data.decks.filter((d) => d.id !== editingId.value);
          if (data.selectedId === editingId.value)
            data.selectedId = data.decks[0]?.id || "";
          DS.write(data);
          decks.value = data.decks;
          signature();
          newDeck();
        } catch (e) {
          error.value = e.message;
        }
      }
      function exportCode() {
        UI.download((deckName.value || "卡组") + ".txt", deckCode.value);
      }
      function backupAll() {
        try {
          UI.download(
            "奥术纪元-本地卡组.json",
            JSON.stringify(DS.load(), null, 2),
            "application/json",
          );
        } catch (e) {
          error.value = e.message;
        }
      }
      function backupRaw() {
        try {
          UI.download(
            "奥术纪元-原始卡组备份.txt",
            localStorage.getItem(DS.KEY) || "",
          );
        } catch (e) {
          error.value = e.message;
        }
      }
      async function restoreBackup(e) {
        const file = e.target.files?.[0];
        if (!file) return;
        try {
          const raw = await file.text();
          const restored = DS.load({ getItem: () => raw });
          restored.decks.forEach((d) => DS.parse(d.code));
          if (storageError.value)
            throw Error("请先保留原始备份，修复本地存储后重试");
          const current = DS.load();
          for (const d of restored.decks)
            current.decks.push({
              ...d,
              id:
                "d_" +
                Date.now().toString(36) +
                "_" +
                Math.random().toString(36).slice(2, 10),
            });
          DS.write(current);
          reloadStore();
          UI.toast("已追加恢复 " + restored.decks.length + " 套卡组");
        } catch (err) {
          error.value = "备份恢复失败：" + err.message;
        }
        e.target.value = "";
      }
      async function loadCards() {
        loading.value = true;
        loadError.value = "";
        try {
          cards.value = await UI.request("/api/cards");
        } catch (e) {
          loadError.value = e.message;
        } finally {
          loading.value = false;
        }
      }
      function beforeUnload(e) {
        if (dirty.value) {
          e.preventDefault();
          e.returnValue = "";
        }
      }
      function onKey(e) {
        if (e.key === "Escape") {
          if (fullCard.value) fullCard.value = null;
          else if (showImport.value) showImport.value = false;
          else selected.value = null;
        }
      }
      onMounted(() => {
        reloadStore();
        const data = (() => {
          try {
            return DS.load();
          } catch {
            return null;
          }
        })();
        const chosen =
          data?.decks.find((d) => d.id === data.selectedId) || data?.decks[0];
        if (chosen) {
          try {
            draft.value = DS.parse(chosen.code);
            editingId.value = chosen.id;
            deckName.value = chosen.name;
          } catch (e) {
            error.value = "已保存卡组无法读取：" + e.message;
          }
        }
        signature();
        loadCards();
        window.addEventListener("beforeunload", beforeUnload);
        window.addEventListener("keydown", onKey);
      });
      onUnmounted(() => {
        window.removeEventListener("beforeunload", beforeUnload);
        window.removeEventListener("keydown", onKey);
      });
      return {
        crest: UI.crest,
        image: UI.image,
        artAttrs: UI.artAttrs,
        imageError: UI.imageError,
        cost: UI.cost,
        description: (s) =>
          String(s || "").replace(
            /\\(火|水|气|地|光|暗|无|血|攻|威|持)/g,
            "$1",
          ),
        cards,
        decks,
        storageError,
        loading,
        loadError,
        query,
        typeFilter,
        versionFilter,
        elementFilter,
        page,
        selected,
        fullCard,
        showImport,
        importCode,
        importError,
        editingId,
        deckName,
        draft,
        draftTab,
        error,
        busy,
        validCode,
        allCards,
        packs,
        elements,
        filtered,
        pageCount,
        effectivePage,
        pageCards,
        deckCode,
        draftHero,
        draftRows,
        openDeck,
        newDeck,
        heroFor,
        deckSummary,
        countInDeck,
        addReason,
        addCard,
        removeCard,
        loadSample,
        importDeck,
        saveDeck,
        useDeck,
        deleteDeck,
        exportCode,
        backupAll,
        backupRaw,
        restoreBackup,
        loadCards,
      };
    },
  }).mount("#app");
})();
