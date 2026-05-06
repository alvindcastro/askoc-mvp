(function () {
  const shell = document.querySelector("[data-chat-endpoint]");
  const form = document.getElementById("chat-form");
  const messageInput = document.getElementById("message");
  const studentInput = document.getElementById("student-id");
  const conversationInput = document.getElementById("conversation-id");
  const conversation = document.getElementById("conversation");
  const intent = document.getElementById("intent");
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
        appendMessage("assistant", body.error ? body.error.message : "The request could not be handled.");
        return;
      }
      conversationInput.value = body.conversation_id || "";
      appendMessage("assistant", body.answer || "");
      renderDetails(body);
    } catch (error) {
      appendMessage("assistant", "The local demo API is not available.");
    } finally {
      setBusy(false);
    }
  });

  function appendMessage(role, text) {
    const article = document.createElement("article");
    article.className = "message " + role;
    const paragraph = document.createElement("p");
    paragraph.textContent = text;
    article.appendChild(paragraph);
    conversation.appendChild(article);
    conversation.scrollTop = conversation.scrollHeight;
  }

  function renderDetails(body) {
    intent.textContent = body.intent ? body.intent.name + " (" + Math.round((body.intent.confidence || 0) * 100) + "%)" : "unknown";
    renderList(sources, body.sources || [], function (source) {
      const details = [source.title, source.chunk_id];
      if (source.confidence) {
        details.push(Math.round(source.confidence * 100) + "%");
      }
      if (source.risk_level) {
        details.push(source.risk_level);
      }
      if (source.freshness_status) {
        details.push(source.freshness_status);
      }
      if (source.caution) {
        details.push(source.caution);
      }
      return details.join(" - ");
    });
    renderList(actions, body.actions || [], function (action) {
      return action.type + " - " + action.status;
    });
    escalation.textContent = body.escalation ? body.escalation.status + " - " + body.escalation.queue : "None";
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

  function setBusy(busy) {
    form.querySelector("button").disabled = busy;
  }
}());
