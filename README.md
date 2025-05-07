# Fin AI

**Fin AI** adalah aplikasi web untuk membantu mencatat dan mengelola keuangan pribadi dengan teknologi AI. Aplikasi ini mampu memahami input transaksi dalam bentuk bahasa manusia, memberikan analisis keuangan otomatis, serta membantu merencanakan anggaran dan pengeluaran secara pintar dan efisien.

## Fitur AI

- Input transaksi keuangan dengan prompt
- Prediksi pengeluaran bulanan berdasarkan histori
- Chat AI untuk konsultasi keuangan pribadi
- OCR untuk membaca dan mencatat dari struk belanja
- Perencanaan keuangan jangka pendek dan panjang
- Rekomendasi pengelolaan keuangan bulanan
- Smart Summary (ringkasan keuangan otomatis)

## Fitur Non-AI

- Riwayat transaksi keuangan lengkap
- Laporan keuangan bulanan dan tahunan

## 🛠️ Tech Stack

- **Golang** – Backend utama
- **Fiber** – Web framework Golang
- **PostgreSQL** – Penyimpanan data transaksi
- **OpenAI GPT-4.1** – Model AI untuk input, insight, dan percakapan
- **Docker** – Containerization dan deployment

## Instalasi

- Clone repository ini

```bash
git clone <repository-url>
cd fin-ai
```

- Buat file `.env` berdasarkan file `.env.example` dan sesuaikan dengan konfigurasi Anda

```bash
cp .env.example .env
```

- Jalankan aplikasi dengan Docker Compose

```bash
docker-compose up --build
```

- Akses aplikasi di `http://localhost:8080`
