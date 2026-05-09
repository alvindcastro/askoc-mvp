(function () {
  const root = document.querySelector("[data-metrics-endpoint]");
  const form = document.getElementById("token-form");
  const tokenInput = document.getElementById("admin-token");
  const status = document.getElementById("admin-status");
  const exportButton = document.getElementById("export-audit");
  const purgeButton = document.getElementById("purge-audit");
  const resetButton = document.getElementById("reset-audit");
  const reviewFilter = document.getElementById("review-filter");
  let currentReviewItems = [];

  if (!root || !form) {
    return;
  }

  const storedToken = window.sessionStorage.getItem("askoc_admin_token") || "";
  if (storedToken) {
    tokenInput.value = storedToken;
    refresh();
  }

  form.addEventListener("submit", function (event) {
    event.preventDefault();
    window.sessionStorage.setItem("askoc_admin_token", tokenInput.value.trim());
    refresh();
  });

  if (reviewFilter) {
    reviewFilter.addEventListener("change", function () {
      renderReviewItems(document.getElementById("review-items"), currentReviewItems);
    });
  }

  exportButton.addEventListener("click", async function () {
    const body = await request(root.dataset.exportEndpoint, { method: "GET" });
    if (body) {
      status.textContent = "Exported " + ((body.events || []).length) + " redacted events";
    }
  });

  resetButton.addEventListener("click", async function () {
    const body = await request(root.dataset.resetEndpoint, { method: "POST" });
    if (body) {
      status.textContent = "Demo audit data reset";
      refresh();
    }
  });

  purgeButton.addEventListener("click", async function () {
    const body = await request(root.dataset.purgeEndpoint, { method: "POST" });
    if (body) {
      status.textContent = "Purged " + (body.pruned || 0) + " expired events";
      refresh();
    }
  });

  async function refresh() {
    const body = await request(root.dataset.metricsEndpoint, { method: "GET" });
    if (!body) {
      return;
    }
    document.getElementById("total-conversations").textContent = String(body.total_conversations || 0);
    document.getElementById("containment-rate").textContent = asPercent(body.containment_rate || 0);
    document.getElementById("escalations").textContent = String(body.escalations || 0);
    document.getElementById("workflows").textContent = String((body.automation && body.automation.workflow_events) || 0);
    document.getElementById("payment-reminders").textContent = String((body.automation && body.automation.payment_reminders_sent) || 0);
    document.getElementById("workflow-failures").textContent = String((body.automation && body.automation.workflow_failures) || 0);
    document.getElementById("low-confidence-count").textContent = String((body.review_queue && body.review_queue.low_confidence_answers) || 0);
    document.getElementById("stale-source-count").textContent = String((body.review_queue && body.review_queue.stale_source_questions) || 0);
    renderTopIntents(document.getElementById("top-intents"), body.top_intents || []);
    let reviewItems = (body.review_queue && body.review_queue.items) || [];
    if (root.dataset.reviewEndpoint) {
      const reviewBody = await request(root.dataset.reviewEndpoint, { method: "GET" });
      if (reviewBody) {
        reviewItems = reviewItems.length ? reviewItems.concat(reviewBody.items || []) : (reviewBody.items || reviewItems);
      }
    }
    currentReviewItems = reviewItems.length ? reviewItems : [];
    renderReviewItems(document.getElementById("review-items"), currentReviewItems);
    status.textContent = "Dashboard refreshed - Evaluation gate evidence remains in reports/eval-summary.md";
  }

  async function request(url, options) {
    const token = tokenInput.value.trim();
    if (!token) {
      status.textContent = "Admin token required";
      return null;
    }
    setBusy(true);
    try {
      const response = await fetch(url, {
        method: options.method,
        headers: { "Authorization": "Bearer " + token }
      });
      const body = await response.json();
      if (!response.ok) {
        status.textContent = body.error ? body.error.message : "Request failed";
        return null;
      }
      return body;
    } catch (error) {
      status.textContent = "Admin API unavailable";
      return null;
    } finally {
      setBusy(false);
    }
  }

  function renderTopIntents(node, items) {
    node.replaceChildren();
    if (!items.length) {
      node.appendChild(emptyItem("No intent metrics yet"));
      return;
    }
    items.forEach(function (item) {
      const row = document.createElement("li");
      const title = document.createElement("span");
      title.className = "row-title";
      title.textContent = item.intent || "unknown";
      row.appendChild(title);
      const meta = document.createElement("span");
      meta.className = "row-meta";
      meta.appendChild(chip("count " + (item.count || 0), "code"));
      row.appendChild(meta);
      node.appendChild(row);
    });
  }

  function renderReviewItems(node, items) {
    node.replaceChildren();
    const filtered = filterReviewItems(items);
    if (!filtered.length) {
      node.appendChild(emptyItem("No review items"));
      return;
    }
    filtered.forEach(function (item) {
      const row = document.createElement("li");
      const title = document.createElement("span");
      title.className = "row-title";
      title.textContent = item.reason || item.status || "review item";
      row.appendChild(title);

      const meta = document.createElement("span");
      meta.className = "row-meta";
      meta.appendChild(chip("trace_id " + (item.trace_id || "not recorded"), "code"));
      meta.appendChild(chip("queue " + (item.queue || "review"), ""));
      meta.appendChild(chip("priority " + (item.priority || "normal"), item.priority === "high" ? "error" : ""));
      meta.appendChild(chip("status " + (item.status || "pending"), item.status === "completed" ? "success" : "warning"));
      meta.appendChild(chip("redacted", "success"));
      row.appendChild(meta);

      const question = document.createElement("span");
      question.textContent = item.question || item.prompt || "Redacted review proof only";
      row.appendChild(question);
      node.appendChild(row);
    });
  }

  function filterReviewItems(items) {
    const mode = reviewFilter ? reviewFilter.value : "all";
    if (mode === "high") {
      return items.filter(function (item) { return item.priority === "high"; });
    }
    if (mode === "redacted") {
      return items.filter(function (item) { return item.question || item.prompt || item.trace_id; });
    }
    return items;
  }

  function chip(text, tone) {
    const node = document.createElement("span");
    node.className = ["chip", tone].filter(Boolean).join(" ");
    node.textContent = text;
    return node;
  }

  function emptyItem(text) {
    const item = document.createElement("li");
    const title = document.createElement("span");
    title.className = "row-title";
    title.textContent = text;
    item.appendChild(title);
    return item;
  }

  function asPercent(value) {
    return Math.round(value * 100) + "%";
  }

  function setBusy(busy) {
    form.querySelector("button").disabled = busy;
    exportButton.disabled = busy;
    purgeButton.disabled = busy;
    resetButton.disabled = busy;
  }
}());
