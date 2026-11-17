# Sprint 1 — 18–24 November 2026

**Tujuan sprint:** backend Go live dan stabil, sisa Express bersih, angka simulasi jujur, dan dasar operasional (backup, keamanan, CI) siap.
Rujukan jangka panjang: `PLAN.md`.

| Hari | Fokus | Selesai jika |
|---|---|---|
| **Rab 18** | Cutover ke Go | `ascard.hallojanu.me` dilayani image `ascard-app:go`, login admin + artikel + semua halaman 200 |
| **Kam 19** | Pantau + bersih-bersih Express | 24 jam tanpa error, `rewrite/go` di-merge ke `main`, Express dihapus |
| **Jum 20** | Backup & keamanan | Backup harian ascard jalan, `npm audit` tanpa high |
| **Sab 21** | Simulasi jujur (1) | Mesin simulasi jadi modul murni + RTP teruji |
| **Min 22** | Simulasi jujur (2) | UI pakai mesin baru, bug totalLoss/reset beres, deploy |
| **Sen 23** | SEO + jalur bantuan | Preview link muncul di WhatsApp, tombol "Butuh Bantuan" ada |
| **Sel 24** | Review sprint | Semua test hijau, CHANGELOG, rilis `v1.1.0` |

---

## Rabu 18 — Cutover ke Go
- [x] Backup ulang DB live tepat sebelum ganti (`ascard_pre-go-cutover-2.sql.gz`)
- [x] `docker compose --env-file .env.production -f docker-compose.go.yml up -d`
- [x] Cek: health, `/`, `/simulation`, `/admin/login`, login admin, buat & hapus artikel uji, rate limit login
- [x] Cek dari publik (HTTPS lewat Cloudflare)
- [x] `deploy.sh` diarahkan ke `docker-compose.go.yml`
- **Rollback:** `docker compose --env-file .env.production -f docker-compose.yml up -d`

## Kamis 19 — Pantau + bersih-bersih Express
- [ ] Cek log `ascard-app` (error 5xx, panic) setelah ±24 jam
- [ ] Merge `rewrite/go` → `main`, push; tag `legacy/express` tetap untuk arsip
- [ ] Hapus `server/`, `prisma/`, `prisma.config.ts`, `Dockerfile` (Node), dependensi backend di `package.json` (express, prisma, bcryptjs, jsonwebtoken, pg, cors, helmet, express-rate-limit, dotenv, nodemon, concurrently)
- [ ] `Dockerfile.go` → `Dockerfile`, `docker-compose.go.yml` → `docker-compose.yml`
- [ ] Script dev: `npm run dev` (Vite) + `make run` di `api/`; update README
- [ ] Hapus 3 worktree Orca + branch `go/*`
- [ ] Image lama `ascard-app:latest` (908 MB) dihapus setelah yakin tidak perlu rollback

## Jumat 20 — Backup & keamanan
- [ ] Backup harian `ascard-db` masuk ke pola `/opt/stacks/backup/databases` seperti app lain (+ retensi)
- [ ] Uji restore dari backup ke DB sementara
- [ ] `npm audit fix` (sisa dependensi frontend), `govulncheck ./...` di `api/`
- [ ] Ganti password Postgres lokal yang pernah bocor di `.env` lama
- [ ] CI GitHub Actions (`.github/workflows/go.yml`) jalan hijau di `main`

## Sabtu 21 — Simulasi jujur (1): mesin
- [ ] `src/lib/sim/engine.js`, `paytable.js`, `stats.js` — logika murni tanpa React
- [ ] Paytable berbasis kelipatan taruhan, RTP teoritis tepat sesuai yang ditampilkan (96,4%)
- [ ] Vitest: RTP teoritis = angka di UI; Monte Carlo 1 juta spin dalam ±0,5%
- [ ] RNG pakai `crypto.getRandomValues()`

## Minggu 22 — Simulasi jujur (2): UI
- [ ] `SimulationLab.jsx` memakai mesin baru
- [ ] Perbaiki totalLoss saat menang, reset grafik & sesi, `SimpleLineChart` keluar dari komponen
- [ ] README & teks UI: hit rate dan RTP konsisten dengan kode
- [ ] Lint frontend 0 error
- [ ] Deploy + cek manual di HP

## Senin 23 — SEO + jalur bantuan
- [ ] `robots.txt`, `sitemap.xml` asli (bukan fallback SPA)
- [ ] Meta description + Open Graph + gambar preview
- [ ] Tombol "Butuh Bantuan" di navigasi → halaman bantuan; nomor layanan **diverifikasi ke sumber resmi** sebelum ditampilkan
- [ ] Peringatan konten singkat di Simulation Lab (tema visualnya mirip slot)

## Selasa 24 — Review sprint
- [ ] `make test`, `make contract`, Vitest, lint — semua hijau
- [ ] Update `CHANGELOG.md`, `VERSION` → `1.1.0`, tag `v1.1.0`
- [ ] Update `PLAN.md` (centang yang selesai) dan rencana Sprint 2 (mode eksperimen, konten artikel, monitoring)

## Di luar sprint ini (Sprint 2+)
Mode eksperimen (Martingale, 1.000 pemain, kalkulator kerugian), editor Markdown + 20 artikel, monitoring uptime/Sentry, provably fair RNG, aksesibilitas & performa.
