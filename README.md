# anon: Anonymous messages

Send a message to a teammate from your terminal. They receive it on Telegram, but **they won't know who sent it**.

```sh
anon send @rahim your PR was really clean
```
```
📩 Anonymous message:

your PR was really clean
```

---

## 1. Install (one time)

Get the **team key** from your admin before you start.

**Mac / Linux** (in a terminal):
```sh
curl -fsSL https://anon-message.anonymus-message.workers.dev/install.sh | sh
```

**Windows** (in PowerShell):
```powershell
irm https://anon-message.anonymus-message.workers.dev/install.ps1 | iex
```

When it asks for `Team key:`, paste the key and press Enter. Then **open a new terminal**.

## 2. Register on Telegram (one time)

You need to do this to **receive** messages.

1. Make sure you have a Telegram **username**. Check it under Settings → Username.
2. Open **[@annonymusofficebot](https://t.me/annonymusofficebot)** in Telegram and tap **Start**.
3. The bot replies `✅ Registered as @your_username`. Anonymous messages will arrive in this chat from now on.

> If you change your Telegram username, send `/start` to the bot again.

## 3. Usage

```sh
anon members                              # see who's registered
anon send @rahim great idea in the meeting today
```

- **Press Tab:** type `anon send @` and press Tab to see your teammates' names.
- **Special characters** like `!`, `?`, `$` or `'`: wrap the message in single quotes, e.g. `anon send @rahim 'great job!'`
- **Long messages** can be piped in: `cat feedback.txt | anon send @rahim`
- Messages can be **up to 4000 characters**.

> **Windows:** in PowerShell, leave out the `@`, e.g. `anon send rahim message`. Tab works without the `@` too: `anon send ra` + Tab.

---

## Troubleshooting

| Problem | Fix |
|---|---|
| `anon: command not found` | Open a new terminal |
| `invalid team key` | Get the new key from your admin and run: `anon config https://anon-message.anonymus-message.workers.dev <TEAM_KEY>` |
| `@x is not registered` | That teammate hasn't tapped Start on the bot yet. Check the name with `anon members` |
| Tab doesn't show names | Open a new terminal. If a new teammate is missing, run `anon members` |
| Message not sent on Windows | Write the name without the `@` |

## Uninstall

**Mac / Linux:**
```sh
rm ~/.local/bin/anon
```
**Windows (PowerShell):**
```powershell
Remove-Item -Recurse "$env:LOCALAPPDATA\anon"
```
