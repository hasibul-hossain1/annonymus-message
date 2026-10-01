# anon: Team er jonno anonymous message

Terminal theke teammate ke message pathan. Teammate message ta paben **Telegram bot er DM e**, kintu janben na ke pathiyeche.

```
Apnar terminal                          Cloudflare Worker                 Teammate er Telegram
anon send @rahim "PR ta review koro" ──► sender save kore na ──────────► 📩 Anonymous message:
                                                                             PR ta review koro
```

- Mac, Windows ar Linux, tin OS ei cholbe. Install korte ekta command lage.
- Sender er nam, IP ba kono info kothao save hoy na.
- Team er shobai **eki team key** use kore, tai server nijeo jane na message ta ke pathiyeche.

---

## Table of contents

- [Team member der jonno](#team-member-der-jonno)
  - [1. Install](#1-install)
  - [2. Telegram e register](#2-telegram-e-register)
  - [3. Message pathano](#3-message-pathano)
  - [Uninstall](#uninstall)
- [Problem hole (Troubleshooting)](#problem-hole-troubleshooting)
- [Admin der jonno (ekbar setup)](#admin-der-jonno-ekbar-setup)
- [Maintenance](#maintenance)
- [Project structure](#project-structure)
- [Privacy kibhabe kaj kore](#privacy-kibhabe-kaj-kore)

---

## Team member der jonno

Admin er kach theke duita jinish nin:
1. **Server URL**, jemon `https://anon-message.xyz.workers.dev`
2. **Team key**, ekta lomba secret string

### 1. Install

**Mac / Linux:** Terminal khule chalan.
```sh
curl -fsSL https://<SERVER_URL>/install.sh | sh
```

**Windows:** PowerShell khule chalan. Start menu te "PowerShell" likhe search korlei paben.
```powershell
irm https://<SERVER_URL>/install.ps1 | iex
```

Installer ja ja korbe:
- Apnar OS ar CPU (Intel, Apple Silicon ba ARM) nije detect korbe.
- `anon` download kore PATH e add korbe. Kono sudo ba admin permission lage na.
- **Team key** chaibe. Admin er deya key ta paste kore Enter chapun.

Install shesh hole **notun ekta terminal khulun**. Puraton terminal e PATH update hoy na.

### 2. Telegram e register

Message **receive** korte chaile ei step ta korte hobe.

1. Telegram e apnar ekta **username** set kora thakte hobe. Settings → Username e giye dekhe nin.
2. Admin er deya bot ta Telegram e khuje ber korun, tarpor **`/start`** pathan.
3. Bot reply dibe: `✅ Registered as @apnar_username`.

Ekhon theke apnar kache anonymous message ashle ei bot er chat e dekhaben.

> Telegram username bodlale bot ke abar `/start` pathate hobe.

### 3. Message pathano

**Ke ke register koreche dekhte:**
```sh
anon members
```
```
@karim
@rahim
@sadia
```

**Message pathate:**
```sh
anon send @rahim bhai tomar PR ta onek clean hoyeche
```
```
✓ sent anonymously to @rahim
```

`@` dewa na dewa duita e cholbe, ar username er boro-chhoto hater lekha (case) matter kore na.

**Quote use korun** jodi message e `!`, `?`, `'`, `$`, `*` er moto special character thake:
```sh
anon send @rahim 'meeting e tomar idea ta joss chilo!'
```

**Lomba ba multi-line message** stdin diye pathano jay:
```sh
echo "line 1
line 2" | anon send @rahim

cat feedback.txt | anon send @rahim
```

**Help dekhte:**
```sh
anon help
```

| Command | Ki kore |
|---|---|
| `anon send @user <message>` | Anonymous message pathay |
| `anon members` | Register kora teammate der list dekhay |
| `anon config <server-url> <team-key>` | Server URL ar team key notun kore set kore |
| `anon help` | Help dekhay |

> Ekta message e shorboccho 4000 character pathano jay.

### Uninstall

**Mac / Linux:**
```sh
rm ~/.local/bin/anon
rm -rf ~/.config/anon                          # Linux
rm -rf ~/Library/Application\ Support/anon     # Mac
```

**Windows (PowerShell):**
```powershell
Remove-Item -Recurse "$env:LOCALAPPDATA\anon", "$env:APPDATA\anon"
```

Telegram e bot ke block korle ar kono message ashbe na.

---

## Problem hole (Troubleshooting)

| Error / shomossha | Karon | Shomadhan |
|---|---|---|
| `anon: command not found` ba `'anon' is not recognized` | Terminal e PATH update hoy nai | Notun terminal khulun |
| `error: not configured, run: anon config ...` | Config file nai | `anon config <SERVER_URL> <TEAM_KEY>` chalan |
| `error: invalid team key` | Team key bhul, ba admin key bodlechen | Admin er kach theke notun key niye `anon config` abar chalan |
| `error: @x is not registered (they need to /start the bot)` | Oi teammate bot ke `/start` pathan nai, ba username bhul likhechen | `anon members` diye username check korun, ba teammate ke `/start` dite bolun |
| `error: telegram delivery failed` | Receiver bot ke block koreche, ba bot token e shomossha | Receiver ke bot unblock kore `/start` dite bolun. Na hole admin ke janan |
| `error: message too long (max 4000)` | Message onek boro | Kaykta bhage bhag kore pathan |
| Bot `/start` er reply dey na | Webhook set kora nai | Admin ke janan ([Webhook check](#webhook-check)) |
| Windows e PowerShell script block kore | Execution policy | `Set-ExecutionPolicy -Scope CurrentUser RemoteSigned` diye abar try korun |

---

## Admin der jonno (ekbar setup)

Shuru korar age egula lagbe:
- **Node.js** 18 ba er beshi version (wrangler chalanor jonno)
- **Go** 1.22 ba er beshi version (binary build korar jonno)
- Ekta free **Cloudflare** account
- **Telegram** account

### Step 1: Telegram bot banan

1. Telegram e **@BotFather** khuje ber kore `/newbot` pathan.
2. Bot er nam din, jemon `Team Anon`, ar username din, jemon `myteam_anon_bot`.
3. BotFather ekta **token** dibe, dekhte `123456:ABC-DEF...` er moto. Eta copy kore rakhun.

> Token ta kauke deben na. Eta diye je keu apnar bot chalate parbe.

### Step 2: CLI binary build korun

```sh
./cli/build.sh
```

Ete 6 ta binary `server/public/dl/` folder e toiri hobe:
`darwin-amd64`, `darwin-arm64`, `linux-amd64`, `linux-arm64`, `windows-amd64.exe` ar `windows-arm64.exe`.

### Step 3: Cloudflare e deploy korun

```sh
cd server
npm install
npx wrangler login            # browser khulbe, Cloudflare e login korun
```

**KV namespace banan.** Ekhane member list save thake.
```sh
npx wrangler kv namespace create MEMBERS
```
Output e ekta `id` paben. Seta `server/wrangler.toml` e boshan:
```toml
[[kv_namespaces]]
binding = "MEMBERS"
id = "ekhane_id_ta_paste_korun"
```

**Secret set korun.** Protiti command input chaibe.
```sh
openssl rand -hex 24          # TEAM_KEY er jonno random string, copy korun
openssl rand -hex 24          # WEBHOOK_SECRET er jonno arekta

npx wrangler secret put BOT_TOKEN        # BotFather er token
npx wrangler secret put TEAM_KEY         # prothom random string
npx wrangler secret put WEBHOOK_SECRET   # ditiyo random string
```

**Deploy korun:**
```sh
npx wrangler deploy
```
Output e apnar **server URL** paben, jemon `https://anon-message.xyz.workers.dev`.

### Step 4: Telegram webhook set korun

Ei step er por Telegram bot er `/start` ar onno message apnar server e pathabe.

```sh
curl "https://api.telegram.org/bot<BOT_TOKEN>/setWebhook?url=<SERVER_URL>/telegram&secret_token=<WEBHOOK_SECRET>"
```
Response e `"ok":true` ashle shob thik ache.

#### Webhook check
```sh
curl "https://api.telegram.org/bot<BOT_TOKEN>/getWebhookInfo"
```
Ekhane `url` field e apnar server URL dekhano uchit, ar `last_error_message` thaka uchit na.

### Step 5: Nije test korun

```sh
curl -fsSL <SERVER_URL>/install.sh | sh     # install korun, team key din
```
Tarpor Telegram e bot ke `/start` pathan, ar nijeke ekta message pathan:
```sh
anon send @apnar_username test message
```
Telegram e `📩 Anonymous message: test message` ashle setup shesh.

### Step 6: Team ke janan

Team e ei message ta share korun:

> **Anonymous message tool install koro:**
>
> Mac/Linux: `curl -fsSL <SERVER_URL>/install.sh | sh`
> Windows (PowerShell): `irm <SERVER_URL>/install.ps1 | iex`
>
> Team key: `<TEAM_KEY>` (DM e pathano)
>
> Install er por Telegram e @<bot_username> ke `/start` pathao.
> Use: `anon send @username message`

> Team key ta public channel e na diye DM e pathan.

---

## Maintenance

### CLI update
`cli/main.go` change korle abar build ar deploy korun:
```sh
./cli/build.sh
cd server && npx wrangler deploy
```
Team member ra install command ta abar chalale notun version peye jabe.

### Team key bodlano
Kono member team chere chole gele team key bodle din:
```sh
cd server
npx wrangler secret put TEAM_KEY     # notun random string
```
Tarpor team ke notun key din. Tara chalabe `anon config <SERVER_URL> <NEW_KEY>`.

### Kono member ke list theke shorano
```sh
cd server
npx wrangler kv key delete --binding MEMBERS <username>
```
Username likhte hobe chhoto hater lekhay (lowercase) ar `@` chara.

### Local e test
```sh
cd server
printf 'BOT_TOKEN=x\nTEAM_KEY=test\nWEBHOOK_SECRET=test\n' > .dev.vars
npx wrangler dev
```
Tarpor `anon config http://localhost:8787 test` diye CLI ta local server e point korun.

---

## Project structure

```
anonymus-message/
├── cli/
│   ├── main.go              # anon CLI (shudhu Go standard library)
│   ├── go.mod
│   └── build.sh             # Shob OS er jonno binary build kore
└── server/
    ├── src/index.js         # Cloudflare Worker: API + Telegram webhook
    ├── public/
    │   ├── install.sh.tpl   # Mac/Linux installer template
    │   ├── install.ps1.tpl  # Windows installer template
    │   └── dl/              # Build kora binary (build.sh toiri kore)
    ├── wrangler.toml        # Cloudflare config
    └── package.json
```

### Server API

| Endpoint | Auth | Ki kore |
|---|---|---|
| `POST /send` | `Authorization: Bearer <TEAM_KEY>` | `{"to": "rahim", "message": "..."}` ke Telegram e pathay |
| `GET /members` | `Authorization: Bearer <TEAM_KEY>` | Register kora username er list |
| `POST /telegram` | Telegram secret header | Telegram webhook, `/start` diye register kore |
| `GET /install.sh` | Nai | Mac/Linux installer, server URL boshano thake |
| `GET /install.ps1` | Nai | Windows installer |
| `GET /dl/<file>` | Nai | CLI binary |

---

## Privacy kibhabe kaj kore

- **Sender er info save hoy na.** Server shudhu `to` ar `message` ney, ar message Telegram e pathiye dey. Kono database e message ba sender save hoy na.
- **Sender ke alada kora jay na.** Shobar team key eki, tai server bujhte pare na kon member pathiyeche.
- **Log off kora.** `wrangler.toml` e `observability` off kora ache, tai Cloudflare request log (IP shoho) save kore na.
- **Ja save hoy:** shudhu Telegram username ar chat_id, message pathanor jonno. `anon members` diye eta dekha jay.

> Cloudflare ar Telegram platform hishebe network level e request dekhte pare. Kintu ei app er code kono sender info rakhe na.
