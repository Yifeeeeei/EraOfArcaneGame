/* Static printed-card help. This module never derives abilities or reads battle state. */
(function (root) {
  const source = root.ArcaneKeywordIndex || { definitions: [], cards: {} };
  const definitions = new Map(source.definitions.map((entry) => [entry.id, Object.freeze(entry)]));
  function entries(card) {
    if (!card || card.is_hidden || card.hidden) return [];
    return (source.cards[String(card.number)] || []).map((id) => definitions.get(id)).filter(Boolean);
  }
  const component = {
    props: { card: { type: Object, default: null } },
    computed: { notes() { return entries(this.card); } },
    template: `<section v-if="notes.length" class="keyword-notes" aria-label="关键词释义">
      <header class="keyword-notes-heading"><span aria-hidden="true">✧</span><h4>关键词释义</h4><small>规则参考</small></header>
      <dl><template v-for="note in notes" :key="note.id"><dt>{{ note.name }}</dt><dd>{{ note.text }}</dd></template></dl>
    </section>`,
  };
  root.ArcaneCardKeywords = { entries, component };
  if (typeof module !== "undefined") module.exports = root.ArcaneCardKeywords;
})(typeof window === "undefined" ? globalThis : window);
