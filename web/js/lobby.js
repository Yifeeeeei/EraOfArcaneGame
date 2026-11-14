(function () {
  const { createApp, ref, computed, onMounted, onUnmounted } = Vue;
  const UI = ArcaneUI,
    DS = ArcaneDecks;
  createApp({
    setup() {
      const playerName = ref(""),
        decks = ref([]),
        selectedId = ref(""),
        cards = ref({}),
        error = ref(""),
        busy = ref(false),
        rooms = ref([]),
        joinRoomID = ref(
          new URLSearchParams(location.search).get("room") || "",
        ),
        loadingRooms = ref(false),
        roomError = ref(""),
        deckExpanded = ref(false),
        inspect = ref(null);
      let timer;
      let store;
      try {
        playerName.value = localStorage.getItem("eraofarcane.playerName") || "";
        store = DS.load();
        decks.value = store.decks;
        selectedId.value = store.selectedId || decks.value[0]?.id || "";
      } catch (e) {
        error.value = e.message;
      }
      const selectedDeck = computed(() =>
        decks.value.find((d) => d.id === selectedId.value),
      );
      const selectedData = computed(() => {
        try {
          return DS.parse(selectedDeck.value?.code || " // // ");
        } catch {
          return { main: [], skills: [], extra: [] };
        }
      });
      const selectedHero = computed(() => cards.value[selectedData.value.hero]);
      const selectedCounts = computed(
        () =>
          `人物 ${selectedData.value.hero ? 1 : 0} · 主卡组 ${selectedData.value.main.length}/30 · 技能 ${selectedData.value.skills.length}/10`,
      );
      const selectedGroups = computed(() =>
        [
          { title: "主卡组", ids: selectedData.value.main },
          { title: "技能池", ids: selectedData.value.skills },
          { title: "额外卡组", ids: selectedData.value.extra },
        ]
          .filter((g) => g.ids.length)
          .map((g) => ({
            ...g,
            cards: DS.groups(g.ids).map((row) => ({
              ...row,
              card: cards.value[row.number],
            })),
          })),
      );
      function saveName() {
        try {
          localStorage.setItem("eraofarcane.playerName", playerName.value);
        } catch {
          UI.toast("名称未能保存");
        }
      }
      function chooseDeck() {
        try {
          store = DS.load();
          store.selectedId = selectedId.value;
          DS.write(store);
        } catch (e) {
          error.value = e.message;
        }
      }
      async function sample() {
        busy.value = true;
        error.value = "";
        try {
          await UI.request("/api/deck/validate", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ deck_code: DS.SAMPLE }),
          });
          store = DS.put(DS.load(), {
            name: "大气 · 基础示例",
            code: DS.SAMPLE,
          });
          DS.write(store);
          decks.value = store.decks;
          selectedId.value = store.selectedId;
        } catch (e) {
          error.value = e.message;
        } finally {
          busy.value = false;
        }
      }
      async function validate() {
        if (!playerName.value.trim()) {
          error.value = "请输入玩家名称";
          document.querySelector("#player-name").focus();
          return false;
        }
        if (!selectedDeck.value) {
          error.value = "请先选择卡组";
          return false;
        }
        await UI.request("/api/deck/validate", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ deck_code: selectedDeck.value.code }),
        });
        saveName();
        return true;
      }
      function enter(options = {}) {
        const p = new URLSearchParams({
          room: joinRoomID.value,
          player_id:
            "p_" +
            Date.now().toString(36) +
            "_" +
            Math.random().toString(36).slice(2, 9),
          player_name: playerName.value || "观战者",
        });
        if (options.spectator) p.set("role", "spectator");
        else p.set("deck_code", selectedDeck.value.code);
        if (options.testMode) p.set("test_mode", "1");
        location.href = "/game.html?" + p;
      }
      async function joinRoom(options = {}) {
        if (busy.value) return;
        error.value = "";
        if (!joinRoomID.value) {
          error.value = "请输入房间号";
          return;
        }
        busy.value = true;
        try {
          if (options.spectator || (await validate())) enter(options);
        } catch (e) {
          error.value = e.message;
        } finally {
          busy.value = false;
        }
      }
      async function createRoom() {
        if (busy.value) return;
        error.value = "";
        busy.value = true;
        try {
          if (!(await validate())) return;
          const room = await UI.request("/api/room/create", { method: "POST" });
          joinRoomID.value = room.room_id;
          enter();
        } catch (e) {
          error.value = e.message;
        } finally {
          busy.value = false;
        }
      }
      async function refreshRooms() {
        if (loadingRooms.value) return;
        loadingRooms.value = true;
        try {
          const data = await UI.request("/api/room/list");
          rooms.value = Array.isArray(data) ? data : [];
          roomError.value = "";
        } catch (e) {
          roomError.value = "无法读取房间：" + e.message;
        } finally {
          loadingRooms.value = false;
        }
      }
      function onKey(e) {
        if (e.key === "Escape") inspect.value = null;
      }
      function onStorage(e) {
        if (e.key === DS.KEY) {
          try {
            store = DS.load();
            decks.value = store.decks;
            selectedId.value = store.selectedId || "";
          } catch (err) {
            error.value = err.message;
          }
        }
      }
      onMounted(async () => {
        refreshRooms();
        timer = setInterval(() => {
          if (!document.hidden) refreshRooms();
        }, 10000);
        document.addEventListener("keydown", onKey);
        window.addEventListener("storage", onStorage);
        try {
          cards.value = await UI.request("/api/cards");
        } catch (e) {
          error.value = "卡牌数据加载失败：" + e.message;
        }
      });
      onUnmounted(() => {
        clearInterval(timer);
        document.removeEventListener("keydown", onKey);
        window.removeEventListener("storage", onStorage);
      });
      return {
        crest: UI.crest,
        image: UI.image,
        artAttrs: UI.artAttrs,
        imageError: UI.imageError,
        playerName,
        decks,
        selectedId,
        selectedDeck,
        selectedHero,
        selectedCounts,
        selectedGroups,
        error,
        busy,
        rooms,
        joinRoomID,
        loadingRooms,
        roomError,
        deckExpanded,
        inspect,
        saveName,
        chooseDeck,
        sample,
        createRoom,
        joinRoom,
        refreshRooms,
      };
    },
  }).mount("#app");
})();
