/* Payment vectors, not individual identical element tokens, are choices. */
(function (root, factory) {
  const planner = factory();
  if (typeof module === "object" && module.exports) module.exports = planner;
  else root.ArcanePayments = planner;
})(typeof window === "undefined" ? globalThis : window, function () {
  function analyze(available, cost, solve) {
    const payment = solve(available, cost);
    if (!payment) return { kind: "unavailable", payment: null };
    // Every different exact payment must spend less of at least one element.
    // Cap each used element below this plan to find a witness, in <= 8 solves.
    for (const [element, amount] of Object.entries(payment)) {
      if (amount <= 0) continue;
      const capped = { ...available, [element]: amount - 1 };
      if (solve(capped, cost)) return { kind: "choice", payment: null };
    }
    return { kind: "unique", payment };
  }
  function context(normalAvailable, solve, strictArcane = 0) {
    if (strictArcane <= 0) return { available: { ...normalAvailable }, solve };
    const normalArcane = Number(normalAvailable["无"] || 0);
    const available = { ...normalAvailable, 无: normalArcane + strictArcane };
    function constrainedSolve(limits, cost) {
      // Dedicated arcane pays generic costs only. An arcane limit caps the
      // combined spend, so uniqueness compares the actual wire payment vector.
      const maximum = Math.min(
        strictArcane,
        cost?.["无"] || 0,
        limits["无"] || 0,
      );
      for (let used = maximum; used >= 0; used--) {
        const remaining = { ...cost, 无: (cost?.["无"] || 0) - used };
        const pool = {
          ...limits,
          无: Math.min(normalArcane, (limits["无"] || 0) - used),
        };
        const payment = solve(pool, remaining);
        if (payment)
          return used > 0
            ? { ...payment, 无: (payment["无"] || 0) + used }
            : payment;
      }
      return null;
    }
    return { available, solve: constrainedSolve };
  }
  return { analyze, context };
});
