(function () {
  const root = document.querySelector("[data-metrics-endpoint]");
  const form = document.getElementById("token-form");
  const tokenInput = document.getElementById("admin-token");
  const status = document.getElementById("admin-status");
  const exportButton = document.getElementById("export-audit");
  const purgeButton = document.getElementById("purge-audit");
  const resetButton = document.getElementById("reset-audit");

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
    renderList(document.getElementById("top-intents"), body.top_intents || [], function (item) {
      return item.intent + " - " + item.count;
    });
    let reviewItems = (body.review_queue && body.review_queue.items) || [];
    if (root.dataset.reviewEndpoint) {
      const reviewBody = await request(root.dataset.reviewEndpoint, { method: "GET" });
      if (reviewBody) {
        reviewItems = reviewBody.items || reviewItems;
      }
    }
    renderList(document.getElementById("review-items"), reviewItems, function (item) {
      return item.reason + " - " + (item.question || item.trace_id || "review item");
    });
    status.textContent = "Dashboard refreshed";
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

  function renderList(node, items, label) {
    node.replaceChildren();
    if (!items.length) {
      const item = document.createElement("li");
      item.textContent = "None";
      node.appendChild(item);
      return;
    }
    items.forEach(function (value) {
      const item = document.createElement("li");
      item.textContent = label(value);
      node.appendChild(item);
    });
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
