# domloc

Map local domains to ports. Zero config. Reliable HTTPS.

```bash
domloc add app.test 3000
# ✓ app.test -> localhost:3000 (https)
```

## Install

**macOS**
```bash
brew install wemit/domloc/domloc
```

**Linux**
```bash
curl -fsSL https://github.com/wemit/domloc/releases/latest/download/domloc-linux-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/').tar.gz | tar -xz
sudo mv domloc /usr/local/bin/
```

Then:
```bash
domloc init
```

Requires `dnsmasq` and `caddy`. If [Homebrew](https://brew.sh) is available they are installed automatically. Otherwise install them manually before running `domloc init`:

```bash
# Ubuntu / Debian
sudo apt install -y dnsmasq
# caddy: https://caddyserver.com/docs/install#debian-ubuntu-raspbian
```

---

## Commands

### `domloc init`

Install dependencies, start services, and trust the local HTTPS CA.

```bash
domloc init             # HTTPS — trusts local CA (sudo once)
domloc init --no-https  # HTTP only — skips CA trust, sets HTTP as default for add/wildcard
```

### `domloc add <domain> <port>`

Route a domain to a local port.

```bash
domloc add app.test 3000            # uses HTTPS default from init
domloc add api.test 4000 --no-https # HTTP regardless of init setting
```

DNS for the TLD is configured automatically on first use — writes a resolver config (requires sudo once per TLD).

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
| dnsmasq agent | No | `domloc init` — user-space, port 5300 |
| Caddy service | Yes | `domloc init` — once, ports 80/443 are privileged |
| CA trust | Yes | `domloc init` — once, installs CA in system trust store |
| DNS resolver config | Yes | First `add` per TLD — once per TLD, never again |
| `caddy reload` | No | Every `add` / `remove` — uses admin API |

`--no-https` on `init` sets HTTP as the default for all subsequent `add` and `wildcard` calls. Per-route override always available:

```bash
domloc init --no-https
domloc add app.test 3000                 # http://app.test (default)
domloc add admin.test 5000 --no-https=false  # https://admin.test (explicit)
```

---

## Existing Caddy or dnsmasq

domloc does not touch your existing Caddy or dnsmasq configuration.

- **Caddy already running**: domloc detects the admin API at `localhost:2019` and injects its routes into a dedicated `domloc` server block via the JSON API. Your other servers and config are untouched.
- **dnsmasq already running**: domloc runs its own dnsmasq instance on port 5300 as a separate user-space service. It writes its own config at `~/.config/domloc/dnsmasq.conf` and never touches your existing dnsmasq.

---

## How it works

```
Browser
  ↓
Caddy (ports 443/80)       ← system service, runs as root
  ↓
localhost:PORT

dnsmasq (port 5300)        ← user-space service
  resolves *.test → 127.0.0.1
```

**macOS**: Caddy runs as a LaunchDaemon. dnsmasq runs as a LaunchAgent. DNS routing via `/etc/resolver/<tld>`.

**Linux**: Caddy runs as a systemd system service. dnsmasq runs as a systemd user service. DNS routing via systemd-resolved drop-in at `/etc/systemd/resolved.conf.d/domloc-<tld>.conf`.

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
| Linux | ✓ Supported |
| Windows | Planned |
