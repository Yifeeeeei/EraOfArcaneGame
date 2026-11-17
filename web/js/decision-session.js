/* Transport state is independent of Vue and of card rules. No automatic replay. */
(function (root) {
    class DecisionSession {
        constructor() { this.revision = null; this.flight = null; this.receipt = null; this.failure = null; this.unknown = false; }
        get busy() { return !!this.flight; }
        failureMessage() { return this.failure?.revision === this.revision ? this.failure.message : ''; }
        begin(requestID, revision) {
            if (this.busy || this.unknown || revision !== this.revision) return false;
            this.receipt = null;
            this.failure = null;
            this.flight = { requestID, revision, accepted: false };
            return true;
        }
        snapshot(revision) {
            if (this.revision !== null && revision < this.revision) return false;
            const changed = revision !== this.revision;
            this.revision = revision;
            if (this.flight && revision > this.flight.revision) {
                this.receipt = this.flight;
                this.flight = null;
            }
            this.unknown = false;
            return changed;
        }
        result(message) {
            if (!this.flight || this.flight.requestID !== message.request_id) {
                // A snapshot can precede its error receipt. Preserve the original
                // error at that exact revision without touching the new decision.
                if (!this.flight && this.receipt?.requestID === message.request_id &&
                    typeof message.revision === 'number' && !message.ok && message.revision === this.revision) {
                    this.failure = {revision:message.revision, message:message.message || '操作失败，请重新选择'};
                    this.receipt = null;
                    return 'error_synced';
                }
                return 'stale';
            }
            if (!message.ok) {
                if (typeof message.revision === 'number') this.failure = {revision:message.revision, message:message.message || '操作失败，请重新选择'};
                if (message.state_advanced) {
                    // The rejected handler may already have changed the board.
                    // Keep the lock until the authoritative snapshot invalidates the draft.
                    return 'error_waiting';
                }
                this.flight = null;
                return 'error';
            }
            this.flight.accepted = true;
            return 'success';
        }
        disconnect() { this.unknown = !!this.flight; this.flight = null; this.receipt = null; }
    }
    const api = { DecisionSession };
    if (typeof module !== 'undefined') module.exports = api;
    root.ArcaneDecisions = api;
})(typeof window === 'undefined' ? globalThis : window);
