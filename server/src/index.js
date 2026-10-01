// Anonymous message relay: CLI -> this Worker -> Telegram bot DM.
// Sender identity is never stored or logged. Everyone uses the same TEAM_KEY,
// so the server itself cannot tell who sent a message.

const TG_API = "https://api.telegram.org/bot";
const MAX_LEN = 4000; // Telegram limit is 4096

export default {
  async fetch(request, env) {
    const url = new URL(request.url);

    try {
      if (request.method === "GET" && (url.pathname === "/install.sh" || url.pathname === "/install.ps1")) {
        return await serveInstaller(request, env, url);
      }
      if (request.method === "POST" && url.pathname === "/telegram") {
        return await handleTelegram(request, env);
      }
      if (request.method === "POST" && url.pathname === "/send") {
        if (!authorized(request, env)) return json({ error: "invalid team key" }, 401);
        return await handleSend(request, env);
      }
      if (request.method === "GET" && url.pathname === "/members") {
        if (!authorized(request, env)) return json({ error: "invalid team key" }, 401);
        const list = await env.MEMBERS.list();
        return json({ members: list.keys.map((k) => k.name).sort() });
      }
      return json({ error: "not found" }, 404);
    } catch {
      return json({ error: "internal error" }, 500);
    }
  },
};

// Installer scripts are templates with the server URL filled in, so a teammate
// only needs to run one command.
async function serveInstaller(request, env, url) {
  const tpl = await env.ASSETS.fetch(new Request(new URL(url.pathname + ".tpl", url)));
  const script = (await tpl.text()).replaceAll("__SERVER__", url.origin);
  return new Response(script, { headers: { "Content-Type": "text/plain; charset=utf-8" } });
}

function authorized(request, env) {
  return request.headers.get("Authorization") === `Bearer ${env.TEAM_KEY}`;
}

// Telegram webhook: a member sends /start to the bot to register their username.
async function handleTelegram(request, env) {
  if (request.headers.get("X-Telegram-Bot-Api-Secret-Token") !== env.WEBHOOK_SECRET) {
    return new Response("forbidden", { status: 403 });
  }

  const update = await request.json();
  const msg = update.message;
  if (!msg || msg.chat.type !== "private" || !msg.text?.startsWith("/start")) {
    return new Response("ok");
  }

  const username = msg.from.username?.toLowerCase();
  if (!username) {
    await tg(env, msg.chat.id, "Apnar Telegram username set kora nai. Settings theke username set kore abar /start din.");
    return new Response("ok");
  }

  await env.MEMBERS.put(username, String(msg.chat.id));
  await tg(env, msg.chat.id, `✅ Registered as @${username}. Ekhon theke team er anonymous message ekhane ashbe.`);
  return new Response("ok");
}

async function handleSend(request, env) {
  const body = await request.json().catch(() => ({}));
  const to = String(body.to || "").replace(/^@/, "").toLowerCase();
  const message = String(body.message || "").trim();

  if (!to || !message) return json({ error: "'to' and 'message' are required" }, 400);
  if (message.length > MAX_LEN) return json({ error: `message too long (max ${MAX_LEN})` }, 400);

  const chatId = await env.MEMBERS.get(to);
  if (!chatId) return json({ error: `@${to} is not registered (they need to /start the bot)` }, 404);

  const ok = await tg(env, chatId, `📩 Anonymous message:\n\n${message}`);
  if (!ok) return json({ error: "telegram delivery failed" }, 502);
  return json({ ok: true });
}

async function tg(env, chatId, text) {
  const res = await fetch(`${TG_API}${env.BOT_TOKEN}/sendMessage`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ chat_id: chatId, text }),
  });
  return res.ok;
}

function json(data, status = 200) {
  return new Response(JSON.stringify(data), {
    status,
    headers: { "Content-Type": "application/json" },
  });
}
