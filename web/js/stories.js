/* Story content is replaceable Markdown data, independent of table/UI rules. */
(function () {
  const { createApp, ref, onMounted } = Vue;
  const stories = [
    { id: "EmptyThrone", title: "王座无人", card: "4111101" },
    { id: "Knight", title: "光辉骑士", card: "4511003" },
    { id: "Moon", title: "第九轮红月", card: "1621109" },
  ].map((s) => ({
    ...s,
    source: "https://yifeeeeei.github.io/EraOfArcane/htmls/" + s.id + ".html",
  }));
  createApp({
    setup() {
      const active = ref(stories[0]),
        blocks = ref([]),
        cards = ref({}),
        error = ref(""),
        loading = ref(false),
        reader = ref(null);
      let sequence = 0;
      async function read(story) {
        active.value = story;
        loading.value = true;
        error.value = "";
        const ticket = ++sequence;
        try {
          const response = await fetch("/assets/stories/" + story.id + ".md");
          if (!response.ok) throw Error("无法读取故事");
          const text = await response.text();
          if (ticket !== sequence) return;
          blocks.value = text
            .replace(/\r/g, "")
            .split(/\n\s*\n/)
            .map((t) => t.trim())
            .filter((t) => t && !/^#\s/.test(t))
            .map((t) =>
              /^#{2,6}\s/.test(t)
                ? { kind: "heading", text: t.replace(/^#{2,6}\s+/, "") }
                : /^[-*_]{3,}$/.test(t)
                  ? { kind: "rule", text: "" }
                  : { kind: "text", text: t },
            );
          reader.value?.scrollTo({ top: 0 });
        } catch (e) {
          if (ticket === sequence) error.value = e.message;
        } finally {
          if (ticket === sequence) loading.value = false;
        }
      }
      onMounted(() => {
        read(active.value);
        ArcaneUI.request("/api/cards")
          .then((data) => (cards.value = data))
          .catch(() => {});
      });
      return {
        stories,
        active,
        blocks,
        cards,
        error,
        loading,
        reader,
        read,
        crest: ArcaneUI.crest,
        image: ArcaneUI.image,
        artAttrs: ArcaneUI.artAttrs,
        imageError: ArcaneUI.imageError,
      };
    },
  }).mount("#app");
})();
