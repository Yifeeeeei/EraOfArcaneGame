/* Transport state is independent of Vue and of card rules. No automatic replay. */
(function (root) {
    class DecisionSession {
        constructor() { this.revision = null; this.flight = null; this.unknown = false; }
        get busy() { return !!this.flight; }
        begin(requestID, revision) {
            if (this.busy || this.unknown || revision !== this.revision) return false;
            this.flight = { requestID, revision, accepted: false };
            return true;
        }
        snapshot(revision) {
            if (this.revision !== null && revision < this.revision) return false;
            const changed = revision !== this.revision;
            this.revision = revision;
            if (this.flight && revision > this.flight.revision) this.flight = null;
            this.unknown = false;
            return changed;
        }
        result(message) {
            if (!this.flight || this.flight.requestID !== message.request_id) return 'stale';
            if (!message.ok) { this.flight = null; return 'error'; }
            // Success receipt alone is not a new decision. Wait for its snapshot.
            this.flight.accepted = true;
            return 'success';
        }
        disconnect() { this.unknown = !!this.flight; this.flight = null; }
    }
    const api = { DecisionSession };
    if (typeof module !== 'undefined') module.exports = api;
    root.ArcaneDecisions = api;
})(typeof window === 'undefined' ? globalThis : window);
