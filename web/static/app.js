(function () {
  const shell = document.querySelector("[data-chat-endpoint]");
  const form = document.getElementById("chat-form");
  const messageInput = document.getElementById("message");
  const studentInput = document.getElementById("student-id");
  const conversationInput = document.getElementById("conversation-id");
  const conversation = document.getElementById("conversation");
  const intent = document.getElementById("intent");
  const traceID = document.getElementById("trace-id");
  const sources = document.getElementById("sources");
  const actions = document.getElementById("actions");
  const escalation = document.getElementById("escalation");

  if (!shell || !form) {
    return;
  }

  form.addEventListener("submit", async function (event) {
    event.preventDefault();
    const message = messageInput.value.trim();
    if (!message) {
      return;
    }

    appendMessage("user", message);
    setBusy(true);

    const payload = {
      conversation_id: conversationInput.value,
      channel: "web",
      message: message,
      student_id: studentInput.value.trim()
    };

    try {
      const response = await fetch(shell.dataset.chatEndpoint, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload)
      });
      const body = await response.json();
      if (!response.ok) {
        appendMessage("assistant error", "The request could not be handled safely. Check the input and try the local demo again.");
        return;
      }
      conversationInput.value = body.conversation_id || "";
      appendMessage("assistant", body.answer || "The local demo returned no answer.");
      renderDetails(body);
    } catch (error) {
      appendMessage("assistant error", "The local demo API is not available.");
    } finally {
      setBusy(false);
    }
  });

  function appendMessage(role, text) {
    const article = document.createElement("article");
    article.className = "message " + role;
    article.dataset.role = role.split(" ")[0];
    const paragraph = document.createElement("p");
    paragraph.textContent = text;
    article.appendChild(paragraph);
    conversation.appendChild(article);
    conversation.scrollTop = conversation.scrollHeight;
  }

  function renderDetails(body) {
    intent.textContent = body.intent ? body.intent.name + " (" + asPercent(body.intent.confidence || 0) + ")" : "unknown";
    traceID.textContent = body.trace_id || "Not assigned";
    renderSourceEvidence(sources, body.sources || [], body);
    renderActionEvidence(actions, body.actions || [], body);
    escalation.textContent = formatEscalation(body.escalation);
  }

  function renderSourceEvidence(node, items, body) {
    node.replaceChildren();
    if (!items.length) {
      node.appendChild(emptyItem(body && body.intent && body.intent.confidence < 0.5 ? "Low-confidence fallback" : "No approved source"));
      return;
    }
    items.forEach(function (source) {
      const item = document.createElement("li");
      const title = document.createElement("span");
      title.className = "evidence-title";
      title.textContent = source.title || source.chunk_id || "Approved source";
      item.appendChild(title);

      const meta = document.createElement("span");
      meta.className = "evidence-meta";
      meta.appendChild(chip("chunk " + (source.chunk_id || source.id || "source"), "evidence-code"));
      meta.appendChild(chip("confidence " + asPercent(source.confidence || 0), source.confidence < 0.5 ? "warning" : "success"));
      if (source.risk_level) {
        meta.appendChild(chip(source.risk_level === "high" ? "High risk" : "risk " + source.risk_level, source.risk_level === "high" ? "error" : ""));
      }
      if (source.freshness_status) {
        meta.appendChild(chip("freshness " + source.freshness_status, source.freshness_status === "stale" ? "warning" : ""));
      }
      if (source.caution) {
        meta.appendChild(chip(source.caution, "warning"));
      }
      item.appendChild(meta);
      node.appendChild(item);
    });
  }

  function renderActionEvidence(node, items, body) {
    node.replaceChildren();
    const allItems = body && body.trace_id ? [{ type: "trace_id", status: "completed", reference_id: body.trace_id, message: "Trace ID links chat, audit, and eval evidence." }].concat(items) : items;
    if (!allItems.length) {
      node.appendChild(emptyItem("No synthetic integration action yet"));
      return;
    }
    allItems.forEach(function (action) {
      const item = document.createElement("li");
      const title = document.createElement("span");
      title.className = "evidence-title";
      title.textContent = labelAction(action.type) + " - " + (action.status || "pending");
      item.appendChild(title);

      const meta = document.createElement("span");
      meta.className = "evidence-meta";
      meta.appendChild(chip("Synthetic integration", ""));
      meta.appendChild(chip(action.type || "action", "evidence-code"));
      if (action.trace_id) {
        meta.appendChild(chip("trace_id " + action.trace_id, "evidence-code"));
      }
      if (action.reference_id) {
        meta.appendChild(chip(referenceLabel(action.type) + " " + action.reference_id, "evidence-code"));
      }
      if (action.idempotency_key) {
        meta.appendChild(chip("idempotency_key " + action.idempotency_key, "evidence-code"));
      }
      if (body && body.escalation && body.escalation.case_id && action.type === "crm_case_created") {
        meta.appendChild(chip("crm_case_id " + body.escalation.case_id, "evidence-code"));
      }
      if (body && body.escalation && body.escalation.priority) {
        meta.appendChild(chip("priority " + body.escalation.priority, body.escalation.priority === "high" ? "error" : ""));
      }
      item.appendChild(meta);
      if (action.message) {
        const message = document.createElement("span");
        message.textContent = action.message;
        item.appendChild(message);
      }
      node.appendChild(item);
    });
  }

  function formatEscalation(value) {
    if (!value) {
      return "None";
    }
    const parts = [value.status || "pending"];
    if (value.queue) {
      parts.push("queue " + value.queue);
    }
    if (value.priority) {
      parts.push("priority " + value.priority);
    }
    if (value.case_id) {
      parts.push("crm_case_id " + value.case_id);
    }
    return parts.join(" - ");
  }

  function labelAction(type) {
    if (type === "trace_id") {
      return "Trace ID";
    }
    return String(type || "action").replace(/_/g, " ");
  }

  function referenceLabel(type) {
    if (type === "payment_reminder_triggered") {
      return "workflow_id";
    }
    if (type === "crm_case_created") {
      return "crm_case_id";
    }
    return "reference_id";
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
    title.className = "evidence-title";
    title.textContent = text;
    item.appendChild(title);
    return item;
  }

  function asPercent(value) {
    return Math.round(value * 100) + "%";
  }

  function setBusy(busy) {
    form.querySelector("button").disabled = busy;
  }
}());
