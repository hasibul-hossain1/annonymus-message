# anon: Anonymous message

Terminal theke teammate ke message pathan. Teammate message ta paben Telegram e, kintu **janben na ke pathiyeche**.

```sh
anon send @rahim bhai tomar PR ta onek clean hoyeche
```
```
📩 Anonymous message:

bhai tomar PR ta onek clean hoyeche
```

---

## 1. Install (ekbar)

Shuru korar age admin er kach theke **team key** ta niye nin.

**Mac / Linux** (terminal e):
```sh
curl -fsSL https://anon-message.anonymus-message.workers.dev/install.sh | sh
```

**Windows** (PowerShell e):
```powershell
irm https://anon-message.anonymus-message.workers.dev/install.ps1 | iex
```

`Team key:` chaile key ta paste kore Enter chapun. Tarpor **notun ekta terminal khulun**.

## 2. Telegram e register (ekbar)

Message **receive** korte chaile ei step ta korte hobe.

1. Telegram e apnar **username** set kora thakte hobe. Settings → Username e giye dekhe nin.
2. Telegram e **[@annonymusofficebot](https://t.me/annonymusofficebot)** khule **Start** chapun.
3. Bot reply dibe `✅ Registered as @apnar_username`. Ekhon theke anonymous message ekhane ashbe.

> Telegram username bodlale bot ke abar `/start` pathan.

## 3. Use

```sh
anon members                              # ke ke ache dekhun
anon send @rahim meeting e tomar idea ta joss chilo
```

- **Tab chapun:** `anon send @` likhe Tab dile teammate der nam ashbe.
- **Special character thakle** (jemon `!`, `?`, `$`, `'`) message ta single quote e rakhun: `anon send @rahim 'darun kaj!'`
- **Lomba message** pipe kore pathano jay: `cat feedback.txt | anon send @rahim`
- **Shorboccho 4000 character** pathano jay.

> **Windows:** PowerShell e `@` chara likhun, jemon `anon send rahim message`. Tab o `@` chara kaj kore: `anon send ra` + Tab.

---

## Problem hole

| Shomossha | Shomadhan |
|---|---|
| `anon: command not found` | Notun terminal khulun |
| `invalid team key` | Admin er kach theke notun key niye chalan: `anon config https://anon-message.anonymus-message.workers.dev <TEAM_KEY>` |
| `@x is not registered` | Oi teammate ekhono bot e Start chapen nai. `anon members` diye nam ta check korun |
| Tab dile nam ashe na | Notun terminal khulun. Notun member na dekhale `anon members` chalan |
| Windows e message jay na | `@` chara nam likhun |

## Uninstall

**Mac / Linux:**
```sh
rm ~/.local/bin/anon
```
**Windows (PowerShell):**
```powershell
Remove-Item -Recurse "$env:LOCALAPPDATA\anon"
```
