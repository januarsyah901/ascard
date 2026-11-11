# PLAN.md — Rancangan Pengembangan Bongkar Ilusi

> Dokumen induk rencana pengembangan jangka pendek, menengah, dan panjang.
> Melanjutkan `.planning/ROADMAP.md` (Milestone 1 selesai) dan menggantikan "Phase 5: To be determined" di `.planning/STATE.md`.
>
> Versi saat ini: **1.0.0** · Disusun: 11 November 2026

---

## Daftar Isi

1. [Visi & Prinsip](#1-visi--prinsip)
2. [Kondisi Saat Ini](#2-kondisi-saat-ini)
3. [Hasil Audit — Masalah yang Harus Diperbaiki Dulu](#3-hasil-audit--masalah-yang-harus-diperbaiki-dulu)
4. [Peta Jalan (Milestone 2–10)](#4-peta-jalan)
5. [Arsitektur Target](#5-arsitektur-target)
6. [Model Data Target](#6-model-data-target)
7. [Rancangan API Target](#7-rancangan-api-target)
8. [Strategi Testing & Kualitas](#8-strategi-testing--kualitas)
9. [Etika, Keamanan Pengguna & Privasi](#9-etika-keamanan-pengguna--privasi)
10. [Metrik Keberhasilan](#10-metrik-keberhasilan)
11. [Risiko & Mitigasi](#11-risiko--mitigasi)
12. [Timeline Indikatif](#12-timeline-indikatif)
13. [Backlog Ide (Belum Dijadwalkan)](#13-backlog-ide-belum-dijadwalkan)
14. [Keputusan yang Masih Terbuka](#14-keputusan-yang-masih-terbuka)

---

## 1. Visi & Prinsip

**Visi:** Menjadi platform edukasi rujukan di Indonesia yang membuktikan — dengan matematika yang bisa diperiksa sendiri — bahwa sistem judi berbasis probabilitas dirancang untuk membuat pemain kalah, dan membantu orang yang terjebak menemukan jalan keluar.

**Prinsip yang tidak boleh dilanggar:**

| # | Prinsip | Artinya dalam praktik |
|---|---------|------------------------|
| P1 | **Jujur secara matematis** | Setiap angka di UI (RTP, hit rate, house edge) harus sama persis dengan yang dihitung kode. Diuji otomatis. |
| P2 | **Transparan & bisa diverifikasi** | Pengguna bisa memeriksa sendiri hasil RNG (provably fair), bukan sekadar teks "HASH". |
| P3 | **Tidak pernah menjadi alat judi** | Tidak ada uang sungguhan, tidak ada hadiah, tidak ada tautan ke situs judi, tidak ada efek "menang" yang memikat. |
| P4 | **Mengutamakan keselamatan pengguna** | Setiap halaman menyediakan akses cepat ke bantuan profesional. |
| P5 | **Privasi sejak awal** | Default anonim. Tidak mengumpulkan data pribadi tanpa alasan kuat dan persetujuan. |
| P6 | **Berbasis bukti** | Klaim klinis/psikologis wajib punya sumber yang dikutip. |
| P7 | **Bilingual penuh** | Semua teks tersedia dalam Bahasa Indonesia dan Inggris, Bahasa Indonesia sebagai default. |

---

## 2. Kondisi Saat Ini

**Stack:** React 19 + Vite 8 + Tailwind 4 + React Router 7 (frontend) · Express 5 + Prisma 7 + PostgreSQL + JWT (backend).

**Yang sudah ada (Milestone 1):**
- Landing page, Knowledge Base, Simulation Lab, Intervention ("Myth Lab")
- Multi-bahasa ID/EN (`LanguageContext`) dan mata uang IDR/USD (`CurrencyContext`)
- Admin login + CRUD artikel (`/admin/dashboard`)
- Artikel dinamis dari database, halaman detail artikel

**Yang belum ada:** testing, CI/CD, deployment, migrasi Prisma yang tercatat (`prisma/migrations` belum ada), dokumentasi API, aksesibilitas, SEO, analitik, sumber bantuan nyata.

---

## 3. Hasil Audit — Masalah yang Harus Diperbaiki Dulu

Ditemukan saat menyusun rencana ini. Diurutkan dari yang paling berbahaya.

### 🔴 Kritis (keamanan)

| # | Masalah | Lokasi | Dampak |
|---|---------|--------|--------|
| A1 | Endpoint **`POST /api/auth/seed` terbuka tanpa autentikasi** | `server/routes/auth.js:35` | Siapa pun bisa membuat akun admin sendiri lalu mengubah atau menghapus semua artikel. |
| A2 | **File `.env` ikut ter-commit ke git** (berisi `DATABASE_URL` dan `JWT_SECRET`) | `.env`, commit `7b57716` | Kredensial database dan secret JWT bocor ke siapa pun yang punya akses repo. **Rotasi kedua nilai itu**, bukan cuma hapus filenya. |
| A3 | Fallback `JWT_SECRET = 'fallback-secret'` | `server/middleware/auth.js:2`, `server/routes/auth.js:7` | Kalau env tidak terbaca, token bisa dipalsukan oleh siapa saja. |
| A4 | Password admin default `password123` di seed | `prisma/seed.js:9` | Mudah ditebak kalau seed dijalankan di production. |
| A5 | Artikel draft (`published: false`) bisa dibaca publik | `server/routes/articles.js:10,23` | Filter hanya di frontend (`KnowledgeBase.jsx`); API tetap mengembalikan draft. |
| A6 | Tidak ada rate limit pada login, CORS terbuka untuk semua origin, `error.message` dikirim ke klien | `server/index.js`, `server/routes/*.js` | Login bisa di-brute force; detail internal database bocor. |

### 🔴 Kritis (integritas edukasi — melanggar P1)

| # | Masalah | Lokasi | Dampak |
|---|---------|--------|--------|
| B1 | **RTP yang diklaim 96,4% tidak sesuai kode.** Hasil simulasi 5 juta putaran: RTP sebenarnya **±12%** pada taruhan default (Rp10.000) dan **±1,2%** pada taruhan Rp100.000. Ini karena nilai menang tidak dikaitkan dengan besar taruhan (rata-rata menang selalu ±$1,60). | `src/pages/SimulationLab.jsx:116-121` | Platform yang mengklaim "data driven" justru menampilkan angka yang salah. Pengguna yang kritis akan kehilangan kepercayaan. |
| B2 | README menyebut hit rate **0,05% (1 banding 2000)**, kode memakai **5%** | `README.md:65` vs `SimulationLab.jsx:116` | Dokumen dan kode tidak konsisten. |
| B3 | Perhitungan `totalLoss` salah saat menang: taruhan tidak dikurangkan (`netChange = winAmount`, seharusnya `winAmount - bid`) | `SimulationLab.jsx:124` | Statistik kerugian di layar terlalu kecil. |
| B4 | "HASH … (UNVERIFIED)" hanyalah angka acak dekoratif | `SimulationLab.jsx:165` | Bertentangan dengan klaim "RNG Transparency". |
| B5 | `handleReset` tidak mereset grafik `history` dan `sessionSpins` | `SimulationLab.jsx:183` | Grafik saldo membawa data sesi lama. |

### 🟠 Tinggi (teknis)

- URL API `http://localhost:5001` di-hardcode di 4 file (`AdminLogin`, `AdminDashboard`, `KnowledgeBase`, `ArticleDetail`) → tidak bisa di-deploy.
- Gambar grid simulasi memakai URL sementara `lh3.googleusercontent.com/aida-public/...` → bisa mati kapan saja.
- `AuthContext` tidak memverifikasi token ke server ("Mock user for now"); dashboard dirender sebelum redirect.
- `dotenv` dipakai tapi tidak dideklarasikan di `package.json` (hanya kebetulan terpasang sebagai dependensi tidak langsung).
- Tidak ada validasi input di API (judul kosong, konten raksasa, tipe salah).
- Belum ada `prisma/migrations` → skema database tidak tercatat.
- Artikel ditampilkan sebagai teks polos; admin tidak bisa memformat (padahal rencana Phase 4 menyebut markdown).
- Pesan "AI Analyst" dan beberapa teks di Simulation Lab masih hardcode bahasa Inggris.
- Tombol "Share Analysis" di `ArticleDetail` tidak berfungsi.

### 🟡 Kebersihan repo

- `Archive.zip` dan file HTML lama (`landing.html`, `knowledge.html`, `intervention.html`) masih ter-commit.
- `package.json` masih berisi placeholder (`Your Name`, `yourusername`).
- `CHANGELOG.md` menyebut "AsCard" dan `simulation.html` yang sudah tidak relevan.
- `SimpleLineChart` didefinisikan di dalam komponen → dibuat ulang setiap render.

---

## 4. Peta Jalan

Setiap milestone punya **tujuan**, **phase**, **tugas**, dan **kriteria selesai**. Milestone 2 dan 3 wajib diselesaikan sebelum aplikasi dibuka ke publik.

---

### Milestone 2 — Hardening & Fondasi Produksi
**Tujuan:** Aplikasi aman dan siap di-deploy. **Target versi:** `1.1.0`

#### Phase 5: Perbaikan Keamanan Darurat
- [ ] Hapus endpoint `/api/auth/seed` (pembuatan admin hanya lewat `prisma/seed.js` atau CLI)
- [ ] Keluarkan `.env` dari git (`git rm --cached .env`), tambahkan ke `.gitignore`, buat `.env.example`
- [ ] **Rotasi** password database dan `JWT_SECRET` yang sudah bocor
- [ ] Pertimbangkan membersihkan riwayat git (`git filter-repo`) jika repo pernah/akan publik
- [ ] Server gagal start kalau `JWT_SECRET` atau `DATABASE_URL` kosong (hapus fallback)
- [ ] Seed admin membaca password dari env (`ADMIN_INITIAL_PASSWORD`), wajib ganti saat login pertama
- [ ] API publik hanya mengembalikan artikel `published: true`; endpoint admin terpisah untuk draft
- [ ] Tambah `helmet`, `express-rate-limit` (login: 5 percobaan / 15 menit / IP), CORS whitelist dari env
- [ ] Error handler terpusat: log detail di server, kirim pesan generik ke klien

**Selesai jika:** tidak ada endpoint tulis tanpa auth; `npm audit` bersih dari level high/critical; secret tidak ada di repo.

#### Phase 6: Konfigurasi & Struktur Kode
- [ ] `VITE_API_URL` di env + satu modul `src/lib/api.js` (fetch wrapper: base URL, token, error handling)
- [ ] Proxy Vite `/api → localhost:5001` untuk development
- [ ] Validasi input dengan **Zod** (dipakai bersama di server dan form admin)
- [ ] `prisma migrate dev --name init` → commit folder `prisma/migrations`
- [ ] Tambah `dotenv` ke dependencies
- [ ] Verifikasi token: endpoint `GET /api/auth/me`, `ProtectedRoute` di frontend, tampilkan loading selama cek
- [ ] Pecah `LanguageContext.jsx` → `src/i18n/id.json` + `src/i18n/en.json` (sesuai `CONCERNS.md`)
- [ ] Unduh gambar simulasi ke `src/assets/sim/` (atau ganti dengan ilustrasi SVG sendiri)
- [ ] Hapus `Archive.zip`, `landing.html`, `knowledge.html`, `intervention.html`
- [ ] Rapikan `package.json`, `CHANGELOG.md`, `README.md` (nama, repo, angka RTP)
- [ ] Prettier + `lint-staged` + Husky pre-commit

**Selesai jika:** `npm run build` menghasilkan build yang bisa jalan di domain mana pun hanya dengan mengubah env.

#### Phase 7: Testing & CI
- [ ] Vitest + React Testing Library (frontend), Supertest (API)
- [ ] Database test terpisah (Docker Postgres / Testcontainers)
- [ ] GitHub Actions: lint → test → build pada setiap PR
- [ ] Target awal coverage: 80% untuk `src/lib/` (logika simulasi), 70% untuk `server/`
- [ ] Playwright E2E: alur spin, ganti bahasa/mata uang, login admin → buat artikel → tampil di Knowledge Base

**Selesai jika:** PR tidak bisa di-merge kalau CI merah.

---

### Milestone 3 — Mesin Simulasi yang Jujur & Bisa Diverifikasi
**Tujuan:** Simulation Lab menjadi bukti matematis yang bisa diperiksa sendiri, bukan sekadar animasi. **Target versi:** `1.2.0`

#### Phase 8: Refactor Mesin Simulasi
- [ ] Pindahkan logika dari `SimulationLab.jsx` (473 baris) ke modul murni `src/lib/sim/`:
  - `engine.js` — `spin(state, bet, rng) → result` (tanpa React, mudah diuji)
  - `paytable.js` — tabel pembayaran berbasis **kelipatan taruhan**, bukan nilai absolut
  - `stats.js` — RTP teoritis, RTP aktual, varians, house edge
- [ ] Paytable dirancang supaya RTP teoritis **tepat 96,4%** (atau angka apa pun yang ditampilkan), dihitung dari rumus `Σ p(i) × multiplier(i)`
- [ ] **Unit test wajib:** RTP teoritis dari paytable = angka di UI; simulasi Monte Carlo 1 juta spin berada dalam ±0,5% dari teoritis
- [ ] Perbaiki bug B3 (totalLoss), B5 (reset)
- [ ] State simulasi disimpan di `localStorage` (sesuai `CONCERNS.md`) dengan tombol "mulai sesi baru"
- [ ] Pindahkan semua teks analyst ke file i18n

#### Phase 9: RNG Provably Fair
- [ ] Ganti `Math.random()` dengan `crypto.getRandomValues()`
- [ ] Skema commit–reveal: `serverSeed` (hash ditampilkan sebelum spin) + `clientSeed` (bisa diubah pengguna) + `nonce`
- [ ] Hasil spin = `HMAC-SHA256(serverSeed, clientSeed:nonce)` → dipetakan ke paytable
- [ ] Halaman **"Verifikasi Sendiri"**: tempel seed, lihat hasil yang sama → membuktikan hasil tidak dimanipulasi *dan* bahwa sistem yang adil pun tetap membuat pemain rugi
- [ ] Tampilkan penjelasan: "Adil ≠ menguntungkan. Hasil ini 100% acak dan tetap membuat Anda kalah."

#### Phase 10: Mode Eksperimen (inti nilai edukasi)
Setiap mode adalah pelajaran interaktif dengan narasi, kontrol, dan grafik.

- [ ] **Auto-spin 10.000×** — grafik saldo turun ke arah garis ekspektasi (hukum bilangan besar)
- [ ] **1.000 pemain sekaligus** — histogram hasil akhir: berapa % yang untung setelah 100 / 1.000 / 10.000 spin
- [ ] **Martingale** — gandakan taruhan setelah kalah, tunjukkan kebangkrutan akibat batas saldo
- [ ] **Gambler's Fallacy** — "sudah kalah 20× berturut-turut, peluang menang berikutnya?" lalu tunjukkan datanya
- [ ] **Near-miss** — demonstrasi bagaimana hampir menang dirancang untuk memicu dopamin
- [ ] **Losses Disguised as Wins (LDW)** — "menang" Rp2.000 dari taruhan Rp10.000 tetap diberi efek perayaan
- [ ] **Kalkulator "Berapa yang Hilang"** — input taruhan/hari × lama bermain → proyeksi kerugian per bulan/tahun, dibandingkan dengan biaya hidup nyata (UMR, biaya sekolah, cicilan)
- [ ] **Bandingkan permainan** — slot vs roulette vs togel vs judi bola: house edge berdampingan
- [ ] Ekspor hasil eksperimen (PNG/CSV) untuk dipakai guru dan peneliti

**Selesai jika:** setiap angka di Simulation Lab punya test yang membuktikan kebenarannya.

---

### Milestone 4 — Knowledge Base Kelas Produksi
**Tujuan:** Konten yang layak dikutip, mudah ditemukan, dan mudah dikelola. **Target versi:** `1.3.0`

#### Phase 11: Konten Kaya
- [ ] Editor Markdown dengan pratinjau langsung (mis. `@uiw/react-md-editor`), render dengan `react-markdown` + `rehype-sanitize`
- [ ] Slug URL (`/artikel/anatomi-jebakan-dopamin`) menggantikan UUID
- [ ] Field baru: ringkasan, gambar sampul, waktu baca, penulis, **daftar sumber/referensi**, tanggal ditinjau
- [ ] Artikel bilingual: satu artikel punya versi ID dan EN
- [ ] Kategori & tag sebagai tabel terpisah (bukan string bebas)
- [ ] Artikel terkait di bagian bawah
- [ ] Komponen interaktif yang bisa disisipkan ke artikel (mis. `<SimEmbed mode="martingale" />`)

#### Phase 12: Penemuan Konten
- [ ] Pencarian full-text (PostgreSQL `tsvector`, mendukung Bahasa Indonesia)
- [ ] Filter kategori/tag, paginasi
- [ ] Tombol Share berfungsi (Web Share API + salin tautan + WhatsApp)
- [ ] Glosarium istilah (RTP, house edge, varians, LDW, chasing losses) dengan tooltip di seluruh situs

#### Phase 13: Kurikulum Konten Awal (minimal 20 artikel)
| Kategori | Contoh judul |
|----------|-------------|
| Matematika | RTP dijelaskan · House edge · Varians & volatilitas · Hukum bilangan besar · Mengapa "pola" itu ilusi |
| Psikologi | Dopamin & antisipasi · Gambler's fallacy · Hot hand · Illusion of control · Near-miss · LDW · Sunk cost |
| Desain | Dark pattern di aplikasi judi · Bonus & "free spin" · Desain suara & visual |
| Finansial | Pinjol & judi online · Dampak pada keluarga · Menghitung kerugian nyata |
| Pemulihan | Tanda-tanda kecanduan · Cara membantu teman/keluarga · Langkah pertama berhenti |
| Hukum | Status hukum judi online di Indonesia · Cara melapor situs judi |

Setiap artikel klinis ditinjau oleh profesional (psikolog/psikiater) sebelum terbit.

---

### Milestone 5 — Intervensi & Dukungan Nyata
**Tujuan:** Mengubah kesadaran menjadi tindakan. **Target versi:** `1.4.0`

#### Phase 14: Jalur Bantuan
- [ ] Tombol **"Butuh Bantuan?"** yang selalu terlihat di semua halaman
- [ ] Halaman Bantuan: layanan kesehatan jiwa nasional, layanan konseling, puskesmas/RS jiwa per provinsi, komunitas pemulihan
  > ⚠️ Semua nomor dan alamat **wajib diverifikasi ulang** dengan sumber resmi sebelum terbit, dan dicek ulang tiap 3 bulan. Simpan di database (tabel `HelpResource`) agar admin bisa memperbarui tanpa deploy.
- [ ] Panduan blokir situs judi (DNS, aplikasi pemblokir, kontrol orang tua) dan cara melapor situs judi ke kanal resmi pemerintah
- [ ] Panduan untuk keluarga & pasangan

#### Phase 15: Alat Bantu Diri
- [ ] **Asesmen mandiri** berbasis instrumen tervalidasi (mis. PGSI — Problem Gambling Severity Index, 9 pertanyaan), diproses di browser saja, tidak dikirim ke server
  - Hasil disertai interpretasi dan rekomendasi langkah, **bukan diagnosis**
  - Periksa lisensi/izin penggunaan instrumen dan versi terjemahan Indonesia yang tervalidasi
- [ ] **Pelacak hari bebas judi** (lokal di perangkat)
- [ ] **Kalkulator "Uangmu Bisa Jadi Apa"** — kerugian diubah menjadi barang/tujuan nyata
- [ ] **Rencana darurat saat ingin bermain** — teknik urge surfing, jeda 10 menit, kontak darurat pribadi
- [ ] Reality check di Simulation Lab yang lebih kuat: setelah N spin atau kerugian X, tampilkan intervensi penuh layar yang menghubungkan ke halaman Bantuan

#### Phase 16: Myth Lab Interaktif
- [ ] Kuis "Mitos atau Fakta" dengan penjelasan + tautan ke artikel dan mode eksperimen terkait
- [ ] Skor pengetahuan sebelum/sesudah (mengukur efek edukasi)

---

### Milestone 6 — Kualitas Pengalaman Pengguna
**Tujuan:** Cepat, mudah diakses, ditemukan mesin pencari. **Target versi:** `1.5.0`

#### Phase 17: Aksesibilitas (WCAG 2.2 AA)
- [ ] Navigasi keyboard penuh, fokus terlihat, label ARIA pada kontrol simulasi
- [ ] Kontras warna lulus AA; dukungan `prefers-reduced-motion` (animasi spin bisa dimatikan)
- [ ] Teks alternatif gambar; audit dengan axe-core di CI

#### Phase 18: Performa
- [ ] Code-splitting per route (`React.lazy`), terutama halaman admin
- [ ] Gambar lokal dalam format WebP/AVIF, font di-self-host
- [ ] Target Lighthouse: Performance ≥ 90 di perangkat mobile kelas menengah, LCP < 2,5 detik di 4G
- [ ] Mode hemat data

#### Phase 19: SEO & Berbagi
- [ ] Prerender/SSG halaman publik (mis. `vite-plugin-ssr`/Vike, atau migrasi ke React Router framework mode)
- [ ] Meta tag per halaman, Open Graph & gambar preview otomatis per artikel
- [ ] `sitemap.xml`, `robots.txt`, data terstruktur (`Article`, `FAQPage`)

#### Phase 20: PWA
- [ ] Bisa di-install di HP, Simulation Lab dan artikel tersimpan bisa dibuka offline
- [ ] Penting untuk pengguna dengan kuota/sinyal terbatas

---

### Milestone 7 — CMS & Admin Lanjutan
**Tujuan:** Tim konten bisa bekerja tanpa developer. **Target versi:** `1.6.0`

- [ ] Peran: `SUPER_ADMIN`, `EDITOR`, `REVIEWER` (peninjau klinis), `CONTRIBUTOR`
- [ ] Alur kerja artikel: Draft → Ditinjau → Disetujui → Terbit; jadwal terbit
- [ ] Riwayat revisi & rollback
- [ ] Unggah media (S3/R2/Cloudinary) dengan kompresi otomatis
- [ ] Audit log semua aksi admin
- [ ] Refresh token + logout semua sesi; token di cookie `httpOnly` (bukan `localStorage`)
- [ ] 2FA (TOTP) untuk admin
- [ ] Kelola sumber bantuan (`HelpResource`), glosarium, kuis, dan paytable dari dashboard
- [ ] Dashboard ringkas: artikel terpopuler, halaman bantuan dikunjungi, hasil kuis agregat

---

### Milestone 8 — Deploy, Operasi & Observabilitas
**Tujuan:** Berjalan stabil di internet publik. **Target versi:** `2.0.0` (rilis publik)

- [ ] Docker + `docker-compose` (app, API, Postgres) untuk development dan production
- [ ] Hosting: frontend statis di CDN (Cloudflare Pages/Vercel/Netlify), API di container (Railway/Fly.io/Render/VPS), Postgres terkelola (Neon/Supabase)
- [ ] Environment terpisah: `development`, `staging`, `production`
- [ ] Pipeline deploy otomatis dari `main` ke staging, rilis tag ke production
- [ ] Backup database harian + uji restore bulanan
- [ ] Monitoring error (Sentry), uptime check, log terstruktur (pino)
- [ ] Endpoint `/api/health` mengecek koneksi database
- [ ] Domain, HTTPS, header keamanan (CSP, HSTS)
- [ ] Rencana respons insiden sederhana (`SECURITY.md`)

---

### Milestone 9 — Analitik Etis & Riset
**Tujuan:** Mengukur dampak tanpa melanggar privasi. **Target versi:** `2.1.0`

- [ ] Analitik tanpa cookie & self-host (Plausible/Umami) — tanpa Google Analytics
- [ ] Event agregat anonim: mode eksperimen dipakai, durasi sesi simulasi, klik halaman Bantuan, skor kuis sebelum/sesudah
- [ ] Banner persetujuan untuk data riset opsional
- [ ] Laporan dampak publik berkala (transparansi untuk mitra & donor)
- [ ] Kebijakan Privasi & Syarat Penggunaan (sesuai UU PDP No. 27/2022)
- [ ] Peluang kolaborasi dengan universitas untuk studi efektivitas

---

### Milestone 10 — Ekspansi & Keberlanjutan
**Tujuan:** Menjangkau lebih banyak orang. **Target versi:** `3.x`

- [ ] **Mode Kelas** untuk guru: kode kelas, siswa ikut eksperimen bersama, hasil agregat di layar proyektor, rencana pembelajaran siap pakai (selaras kurikulum Matematika — peluang & statistik)
- [ ] **Paket workshop** untuk sekolah, kampus, kantor, komunitas (slide, panduan fasilitator, lembar kerja)
- [ ] **Widget embed** Simulation Lab untuk media berita dan situs mitra
- [ ] **API publik** baca-saja untuk artikel & data house edge
- [ ] Bahasa daerah (Jawa, Sunda) dan bahasa lain di Asia Tenggara
- [ ] Konten video pendek & infografik untuk media sosial (dibuat dari data eksperimen)
- [ ] Aplikasi mobile (Capacitor/React Native) jika PWA tidak cukup
- [ ] Program kontributor artikel dengan peninjauan ahli
- [ ] Pendanaan: hibah, CSR, kemitraan lembaga kesehatan — **tanpa iklan** dan tanpa sponsor dari industri yang bertentangan dengan misi

---

## 5. Arsitektur Target

```
┌──────────────────────────────┐        ┌──────────────────────────────┐
│  Frontend (React + Vite)     │        │  API (Express 5)             │
│  CDN statis / prerender      │ HTTPS  │  /api/v1/...                 │
│                              │───────▶│  helmet · rate-limit · zod   │
│  src/                        │        │  auth (JWT httpOnly cookie)  │
│   ├─ pages/                  │        │  pino logging · Sentry       │
│   ├─ components/             │        └──────────────┬───────────────┘
│   ├─ features/               │                       │ Prisma 7
│   │   ├─ simulation/         │                       ▼
│   │   ├─ knowledge/          │        ┌──────────────────────────────┐
│   │   ├─ intervention/       │        │  PostgreSQL                  │
│   │   └─ admin/              │        │  articles · tags · admins    │
│   ├─ lib/                    │        │  revisions · audit_logs      │
│   │   ├─ api.js              │        │  help_resources · glossary   │
│   │   └─ sim/  (logika murni,│        └──────────────────────────────┘
│   │            tanpa React)  │
│   └─ i18n/  id.json en.json  │        Media → Object storage (R2/S3)
└──────────────────────────────┘
```

**Keputusan penting:**
- **Simulasi tetap berjalan di browser.** Tidak butuh server, tidak ada data taruhan yang dikirim — lebih privat dan murah. Server hanya dipakai untuk konten dan admin.
- Logika simulasi berupa **fungsi murni** di `src/lib/sim/`, sehingga bisa diuji, dipakai ulang di widget embed, dan diverifikasi.
- API diberi versi (`/api/v1`) sebelum ada konsumen eksternal.
- Pertimbangkan migrasi ke **TypeScript** bertahap mulai dari `src/lib/` dan `server/` (Milestone 2–3).

---

## 6. Model Data Target

```prisma
enum Role { SUPER_ADMIN EDITOR REVIEWER CONTRIBUTOR }
enum ArticleStatus { DRAFT IN_REVIEW APPROVED PUBLISHED ARCHIVED }
enum Locale { id en }

model User {
  id           String    @id @default(uuid())
  email        String    @unique
  name         String
  passwordHash String
  role         Role      @default(CONTRIBUTOR)
  totpSecret   String?
  mustChangePw Boolean   @default(true)
  createdAt    DateTime  @default(now())
  articles     Article[] @relation("author")
  reviews      Article[] @relation("reviewer")
  auditLogs    AuditLog[]
}

model Article {
  id           String        @id @default(uuid())
  status       ArticleStatus @default(DRAFT)
  coverImage   String?
  authorId     String
  author       User          @relation("author", fields: [authorId], references: [id])
  reviewerId   String?
  reviewer     User?         @relation("reviewer", fields: [reviewerId], references: [id])
  reviewedAt   DateTime?
  publishAt    DateTime?
  categoryId   String?
  category     Category?     @relation(fields: [categoryId], references: [id])
  tags         Tag[]
  translations ArticleTranslation[]
  revisions    ArticleRevision[]
  sources      Source[]
  viewCount    Int           @default(0)
  createdAt    DateTime      @default(now())
  updatedAt    DateTime      @updatedAt
}

model ArticleTranslation {
  id        String  @id @default(uuid())
  articleId String
  article   Article @relation(fields: [articleId], references: [id], onDelete: Cascade)
  locale    Locale
  slug      String
  title     String
  summary   String
  content   String  @db.Text   // Markdown
  @@unique([locale, slug])
  @@unique([articleId, locale])
}

model ArticleRevision { id String @id @default(uuid()) articleId String article Article @relation(fields: [articleId], references: [id], onDelete: Cascade) snapshot Json editorId String createdAt DateTime @default(now()) }
model Category  { id String @id @default(uuid()) slug String @unique nameId String nameEn String articles Article[] }
model Tag       { id String @id @default(uuid()) slug String @unique nameId String nameEn String articles Article[] }
model Source    { id String @id @default(uuid()) articleId String article Article @relation(fields: [articleId], references: [id], onDelete: Cascade) citation String url String? }

model HelpResource {
  id          String   @id @default(uuid())
  name        String
  type        String   // hotline, klinik, komunitas, online
  phone       String?
  url         String?
  province    String?
  hours       String?
  verifiedAt  DateTime // wajib diisi; tampilkan peringatan jika > 90 hari
  active      Boolean  @default(true)
}

model GlossaryTerm { id String @id @default(uuid()) slug String @unique termId String termEn String defId String @db.Text defEn String @db.Text }

model AuditLog {
  id        String   @id @default(uuid())
  userId    String
  user      User     @relation(fields: [userId], references: [id])
  action    String   // article.publish, user.create, ...
  target    String?
  meta      Json?
  createdAt DateTime @default(now())
}
```

Migrasi dari model `Admin`/`Article` saat ini dilakukan bertahap dengan skrip migrasi data (artikel lama → `ArticleTranslation` locale `id`).

---

## 7. Rancangan API Target

| Method | Endpoint | Akses | Keterangan |
|--------|----------|-------|------------|
| GET | `/api/v1/health` | publik | status + koneksi DB |
| POST | `/api/v1/auth/login` | publik (rate limited) | set cookie httpOnly |
| POST | `/api/v1/auth/logout` | login | |
| GET | `/api/v1/auth/me` | login | verifikasi sesi |
| POST | `/api/v1/auth/2fa/verify` | login | |
| GET | `/api/v1/articles?locale=&q=&category=&tag=&page=` | publik | hanya `PUBLISHED` |
| GET | `/api/v1/articles/:locale/:slug` | publik | |
| GET | `/api/v1/categories`, `/tags`, `/glossary` | publik | |
| GET | `/api/v1/help-resources?province=` | publik | |
| GET/POST/PATCH/DELETE | `/api/v1/admin/articles[/:id]` | EDITOR+ | termasuk draft |
| POST | `/api/v1/admin/articles/:id/submit` · `/approve` · `/publish` | sesuai peran | alur kerja |
| GET | `/api/v1/admin/articles/:id/revisions` | EDITOR+ | |
| POST | `/api/v1/admin/media` | EDITOR+ | unggah gambar |
| CRUD | `/api/v1/admin/users`, `/help-resources`, `/glossary` | SUPER_ADMIN | |
| GET | `/api/v1/admin/audit-logs` | SUPER_ADMIN | |

Dokumentasi otomatis dengan OpenAPI (mis. `zod-to-openapi`) dan halaman Swagger di staging.

---

## 8. Strategi Testing & Kualitas

| Lapisan | Alat | Yang diuji | Prioritas |
|---------|------|-----------|-----------|
| Unit | Vitest | `src/lib/sim/*` — RTP teoritis = UI, Monte Carlo, provably fair deterministik, kalkulator | **Tertinggi** |
| Unit | Vitest | i18n: setiap key ada di `id.json` dan `en.json` | Tinggi |
| Komponen | React Testing Library | form admin, kontrol simulasi, toggle bahasa/mata uang | Sedang |
| API | Supertest + Postgres test | auth, izin per peran, draft tidak bocor, validasi input | Tinggi |
| E2E | Playwright | alur utama pengguna & admin, di Chrome + mobile viewport | Sedang |
| Aksesibilitas | axe-core (Playwright) | semua halaman publik | Sedang |
| Keamanan | `npm audit`, Dependabot, CodeQL | dependensi & pola rawan | Tinggi |

**Test yang wajib ada sebelum rilis publik:**
1. `paytable RTP equals displayed RTP`
2. `monte carlo RTP within 0.5% of theoretical after 1M spins`
3. `same seeds produce same outcome` (provably fair)
4. `GET /articles never returns unpublished`
5. `admin endpoints reject missing/invalid/expired token`
6. `every i18n key exists in both locales`

---

## 9. Etika, Keamanan Pengguna & Privasi

- **Tidak boleh menjadi pintu masuk judi.** Tidak ada uang sungguhan, kredit yang bisa ditukar, leaderboard "pemenang", atau tautan afiliasi. Desain menang dibuat datar, tidak meriah — kecuali saat sengaja mendemonstrasikan trik desain (mode LDW/near-miss), yang selalu diberi label jelas.
- **Peringatan konten** di Simulation Lab bagi orang yang sedang dalam pemulihan, dengan pilihan langsung ke halaman Bantuan.
- **Pembatasan usia:** konten tetap aman untuk pelajar (sasaran edukasi sekolah), tanpa unsur yang mengglorifikasi.
- **Klaim klinis** hanya dari sumber yang dikutip; tagline "Clinically Validated" sebaiknya ditinjau ulang sampai benar-benar ada validasi klinis atau peninjau profesional — klaim yang belum terbukti melemahkan kredibilitas.
- **Privasi:** asesmen mandiri & pelacak diproses lokal; tidak ada akun pengguna publik kecuali benar-benar dibutuhkan; patuh UU PDP.
- **Status hukum:** pastikan materi tidak bisa ditafsirkan sebagai promosi judi (UU ITE Pasal 27 ayat 2); konsultasikan dengan penasihat hukum sebelum rilis publik dan kemitraan.

---

## 10. Metrik Keberhasilan

| Area | Metrik | Target awal (6 bulan setelah rilis publik) |
|------|--------|---------------------------------------------|
| Edukasi | Kenaikan skor kuis sebelum → sesudah | +30% |
| Edukasi | % pengguna yang mencoba ≥ 1 mode eksperimen | 40% |
| Dampak | Klik ke halaman Bantuan / sumber bantuan | dipantau, laporan bulanan |
| Jangkauan | Pengunjung unik per bulan | 10.000 |
| Jangkauan | Sekolah/komunitas yang memakai Mode Kelas | 20 |
| Konten | Artikel terbit & ditinjau ahli | 30 |
| Kualitas | Lighthouse (mobile) Performance / Accessibility | ≥ 90 / ≥ 95 |
| Keandalan | Uptime | 99,5% |
| Teknis | Coverage `src/lib/sim` | ≥ 90% |

---

## 11. Risiko & Mitigasi

| Risiko | Kemungkinan | Dampak | Mitigasi |
|--------|-------------|--------|----------|
| Kredensial yang bocor di `.env` disalahgunakan | Tinggi | Tinggi | Rotasi segera (Phase 5) |
| Angka simulasi salah merusak kredibilitas | Sudah terjadi | Tinggi | Phase 8 + test wajib |
| Simulasi justru memicu keinginan berjudi | Sedang | Tinggi | Desain menang datar, peringatan konten, uji pengguna dengan konselor |
| Situs disalahartikan sebagai promosi judi / diblokir | Sedang | Tinggi | Konsultasi hukum, bahasa edukatif jelas, tanpa elemen hadiah |
| Nomor bantuan kedaluwarsa | Sedang | Tinggi | `verifiedAt` + pengingat 90 hari di dashboard |
| Konten klinis tidak akurat | Sedang | Tinggi | Alur REVIEWER wajib |
| Proyek bergantung pada satu developer | Tinggi | Sedang | Dokumentasi, CONTRIBUTING, test & CI |
| Biaya hosting | Rendah | Sedang | Frontend statis + free tier; cari hibah |

---

## 12. Timeline Indikatif

Asumsi: 1–2 developer paruh waktu. Sesuaikan dengan kapasitas nyata.

| Periode | Milestone | Versi |
|---------|-----------|-------|
| Minggu 1 | **Phase 5 (darurat keamanan)** — kerjakan sekarang | 1.0.1 |
| Minggu 2–4 | Sisa Milestone 2 | 1.1.0 |
| Bulan 2 | Milestone 3 — simulasi jujur & provably fair | 1.2.0 |
| Bulan 3 | Milestone 4 — Knowledge Base | 1.3.0 |
| Bulan 4 | Milestone 5 — Intervensi & bantuan | 1.4.0 |
| Bulan 5 | Milestone 6 + 7 — kualitas & CMS | 1.5.0–1.6.0 |
| Bulan 6 | Milestone 8 — **rilis publik** | 2.0.0 |
| Bulan 7–9 | Milestone 9 — analitik & riset | 2.1.0 |
| Bulan 10+ | Milestone 10 — ekspansi | 3.x |

**Urutan kerja minggu ini:**
1. Hapus `/api/auth/seed`
2. `git rm --cached .env`, buat `.env.example`, rotasi password DB & `JWT_SECRET`
3. Hapus fallback JWT secret
4. Filter `published` di API
5. Perbaiki angka RTP/hit rate di README dan UI agar tidak menyesatkan sampai Phase 8 selesai

---

## 13. Backlog Ide (Belum Dijadwalkan)

- Simulasi **judi bola/parlay**: menunjukkan bagaimana odds parlay memperbesar house edge
- Simulasi **togel**: peluang 4D = 1/10.000 dengan pembayaran jauh di bawah peluang wajar
- Simulasi **"pinjol untuk judi"**: bunga berbunga + kerugian judi dalam satu grafik
- Mode **"Bandar"**: pengguna berperan sebagai kasino dan melihat keuntungan dari 10.000 pemain
- Cerita interaktif (*choose your own path*) berbasis kisah nyata yang dianonimkan
- Kartu statistik yang bisa dibagikan ke media sosial ("Dalam 1.000 spin, saya kehilangan Rp…")
- Integrasi chatbot pendamping berbasis LLM dengan batasan keamanan ketat dan rujukan ke profesional
- Ekstensi browser yang menampilkan pengingat saat membuka situs judi
- Mode gelap
- Podcast/audio artikel untuk aksesibilitas

---

## 14. Keputusan yang Masih Terbuka

| # | Pertanyaan | Opsi | Rekomendasi |
|---|-----------|------|-------------|
| K1 | Angka hit rate & RTP mana yang dipakai? | 5% vs 0,05%; RTP 96,4% | RTP 96,4% (realistis seperti slot sungguhan) dengan paytable multi-tingkat; hit rate keseluruhan ±20–30% tapi kebanyakan LDW — ini justru yang paling mendidik |
| K2 | Pindah ke TypeScript? | Ya / tidak / bertahap | Bertahap, mulai `src/lib/sim` dan `server/` |
| K3 | SSR/prerender untuk SEO? | Vike / React Router framework mode / Next.js | React Router framework mode (sudah pakai RR7, migrasi paling kecil) |
| K4 | Perlu akun pengguna publik? | Ya / tidak | Tidak dulu; semua fitur pribadi disimpan lokal |
| K5 | Tagline "Clinically Validated" | Pertahankan / ubah | Ubah menjadi "Ditinjau Profesional" setelah ada peninjau; hapus sampai saat itu |
| K6 | Lisensi konten artikel | MIT (kode) + ? | CC BY-SA 4.0 untuk konten agar bisa dipakai sekolah |
| K7 | Hosting & domain | — | Tentukan sebelum Milestone 8 |

---

*Perbarui dokumen ini setiap akhir milestone. Status phase yang sedang berjalan dicatat di `.planning/STATE.md`.*
