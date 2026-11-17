/* Both player entry points use the same pinned, validated official catalog. */
(function () {
  const { ref, computed, onMounted, onUnmounted, nextTick } = Vue;
  window.ArcanePreconstructedPicker = {
    props: {
      cards: { type: Object, required: true },
      busy: Boolean,
      error: { type: String, default: "" },
      actionLabel: { type: String, default: "选用" },
    },
    emits: ["pick", "close"],
    setup(props, { emit }) {
      const dialog = ref(null),
        catalog = ref(null),
        group = ref("新手入门"),
        loading = ref(true),
        loadError = ref("");
      const controller = new AbortController();
      const opener = document.activeElement;
      const groups = computed(() => [
        "全部",
        ...new Set((catalog.value?.decks || []).map((d) => d.id.split("/")[0])),
      ]);
      const recipes = computed(() =>
        (catalog.value?.decks || []).filter(
          (d) => group.value === "全部" || d.id.split("/")[0] === group.value,
        ),
      );
      function hero(recipe) {
        return props.cards[ArcaneDecks.parse(recipe.deckCode).hero];
      }
      function close() {
        if (!props.busy) emit("close");
      }
      function keydown(event) {
        if (event.key === "Escape") {
          event.stopPropagation();
          close();
        }
        if (event.key !== "Tab") return;
        const nodes = [
          ...dialog.value.querySelectorAll("button:not(:disabled), a[href]"),
        ];
        const first = nodes[0],
          last = nodes.at(-1);
        if (event.shiftKey && document.activeElement === first) {
          event.preventDefault();
          last?.focus();
        } else if (!event.shiftKey && document.activeElement === last) {
          event.preventDefault();
          first?.focus();
        }
      }
      async function load() {
        loading.value = true;
        loadError.value = "";
        try {
          catalog.value = await ArcaneUI.request(
            "/assets/preconstructed-decks.json",
            {
              signal: controller.signal,
            },
          );
          if (!groups.value.includes(group.value)) group.value = "全部";
        } catch (e) {
          if (e.name !== "AbortError")
            loadError.value = "预组卡组加载失败，请重试";
        } finally {
          loading.value = false;
        }
      }
      onMounted(async () => {
        await nextTick();
        dialog.value.querySelector("button")?.focus();
        load();
      });
      onUnmounted(() => {
        controller.abort();
        if (opener?.isConnected) opener.focus();
      });
      return {
        dialog,
        catalog,
        group,
        loading,
        loadError,
        groups,
        recipes,
        hero,
        close,
        keydown,
        load,
        artAttrs: ArcaneUI.artAttrs,
        imageError: ArcaneUI.imageError,
      };
    },
    template: `
      <div class="arcane-dialog-backdrop" @click.self="close">
        <section ref="dialog" class="arcane-dialog preconstructed-dialog" role="dialog"
          aria-modal="true" aria-label="预组卡组" @keydown="keydown">
          <button class="btn dialog-close" aria-label="关闭预组卡组" @click="close" :disabled="busy">×</button>
          <h2>预组卡组</h2>
          <p class="preconstructed-intro">从新手入门开始，或选择不同版本的预组。</p>
          <p v-if="loading" role="status">正在读取预组…</p>
          <div v-else-if="loadError" class="notice-error" role="alert">
            {{ loadError }} <button class="mini-btn" @click="load">重试</button>
          </div>
          <template v-else>
            <div class="preconstructed-tabs" aria-label="预组分类">
              <button v-for="item in groups" :key="item" class="mini-btn"
                :aria-pressed="group===item" @click="group=item" :disabled="busy">{{ item }}</button>
            </div>
            <div class="preconstructed-list">
              <button v-for="recipe in recipes" :key="recipe.id" class="preconstructed-deck"
                :aria-label="actionLabel+' '+recipe.name+' · '+recipe.version.join(' / ')"
                @click="$emit('pick',recipe)" :disabled="busy">
                <span class="preconstructed-hero card-art">
                  <img v-bind="artAttrs(hero(recipe))" @error="imageError" alt="" />
                </span>
                <span class="preconstructed-info"><strong>{{ recipe.name }}</strong>
                  <span>{{ hero(recipe)?.name || '人物资料读取中…' }}</span>
                  <small>{{ recipe.version.join(' / ') }} · 30 张主卡 · 10 张技能</small>
                </span>
                <span class="preconstructed-action">{{ busy?'处理中…':actionLabel }} ›</span>
              </button>
            </div>
            <p v-if="error" class="notice-error" role="alert">{{ error }}</p>
            <footer class="preconstructed-source">
              <span>{{ catalog.decks.length }} 套可用预组</span>
              <a href="https://github.com/Yifeeeeei/ArcaneDecks" target="_blank" rel="noopener">预组来源 ↗</a>
            </footer>
          </template>
        </section>
      </div>`,
  };
})();
