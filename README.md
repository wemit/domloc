# domloc

Map local domains to ports. Zero config. Reliable HTTPS.

```bash
domloc add app.test 3000
# ✓ app.test -> localhost:3000 (https)
```

## Install

```bash
brew install wemit/domloc/domloc
domloc init
```

Requires [Homebrew](https://brew.sh). Installs `dnsmasq` and `caddy` automatically if missing.

---

## Commands

### `domloc init`

Install dependencies, start services, and trust the local HTTPS CA.

```bash
domloc init           # HTTPS — trusts local CA (sudo/Touch ID once)
domloc init --no-https  # HTTP only — skips CA trust, sets HTTP as default for add/wildcard
```

### `domloc add <domain> <port>`

Route a domain to a local port.

```bash
domloc add app.test 3000            # uses HTTPS default from init
domloc add api.test 4000 --no-https # HTTP regardless of init setting
```

DNS for the TLD is configured automatically on first use — writes `/etc/resolver/<tld>` (requires sudo once per TLD).

### `domloc wildcard <pattern> <port>`

Route all subdomains matching a pattern to a port.

```bash
domloc wildcard "*.app.test" 3000
```

### `domloc open <domain>`

Open a route in the browser.

```bash
domloc open app.test
```

### `domloc remove <domain>`

Remove a route. Alias: `domloc rm`

### `domloc list`

Show all routes.

```
  DOMAIN                         PORT     HTTPS   WILDCARD   PROXY
  app.test                       3000     yes     no         running
  *.app.test                     3000     yes     yes        running
  api.test                       4000     no      no         running
```

### `domloc doctor`

Diagnose environment issues.

### `domloc reset`

Stop all services and remove generated configs.

```bash
domloc reset           # keeps routes.json
domloc reset --hard    # removes everything including routes
```

---

## HTTPS

domloc uses [Caddy](https://caddyserver.com)'s built-in local CA — no mkcert, no manual cert management.

| Step | Sudo | When |
|---|---|---|
| dnsmasq LaunchAgent | No | `domloc init` — user-space, port 5300 |
| Caddy LaunchDaemon | Yes | `domloc init` — once, ports 80/443 are privileged |
| CA trust | Yes (Touch ID) | `domloc init` — once, installs CA in system keychain |
| DNS resolver file | Yes | First `add` per TLD — writes `/etc/resolver/<tld>` |
| `caddy reload` | No | Every `add` / `remove` — uses admin API |

`--no-https` on `init` sets HTTP as the default for all subsequent `add` and `wildcard` calls. Per-route override always available:

```bash
domloc init --no-https
domloc add app.test 3000              # http://app.test (default)
domloc add api.test 4000              # http://api.test (default)
domloc add admin.test 5000 --no-https=false  # https://admin.test (explicit override)
```

---

## How it works

```
Browser
  ↓
Caddy (ports 443/80)       ← LaunchDaemon, runs as root
  ↓
localhost:PORT

dnsmasq (port 5300)        ← LaunchAgent, runs as user
  resolves *.test → 127.0.0.1

/etc/resolver/test
  nameserver 127.0.0.1
  port 5300
```

State lives in `~/.config/domloc/`:

| File | Purpose |
|---|---|
| `routes.json` | Route registry |
| `Caddyfile` | Generated Caddy config |
| `caddy-data/` | Caddy PKI / local CA |
| `dnsmasq.conf` | Generated dnsmasq config |
| `caddy.log` | Caddy stdout/stderr |
| `dnsmasq.log` | dnsmasq stdout/stderr |

---

## Platform support

| Platform | Status |
|---|---|
| macOS | ✓ Supported |
| Linux | Planned |
| Windows | Planned |

---

## License

MIT
