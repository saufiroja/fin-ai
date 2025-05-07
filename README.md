# Fin AI

**Fin AI** adalah aplikasi web untuk membantu mencatat dan mengelola keuangan pribadi dengan teknologi AI. Aplikasi ini mampu memahami input transaksi dalam bentuk bahasa manusia, memberikan analisis keuangan otomatis, serta membantu merencanakan anggaran dan pengeluaran secara pintar dan efisien.

## ✨ Fitur

### 🔍 Fitur AI

- Input transaksi keuangan dengan prompt
- Prediksi pengeluaran bulanan berdasarkan histori
- Chat AI untuk konsultasi keuangan pribadi
- OCR untuk membaca dan mencatat dari struk belanja
- Perencanaan keuangan jangka pendek dan panjang
- Rekomendasi pengelolaan keuangan bulanan
- Smart Summary (ringkasan keuangan otomatis)

### 📋 Fitur Non-AI

- Riwayat transaksi keuangan lengkap
- Laporan keuangan bulanan dan tahunan

## 🛠️ Tech Stack

- **Golang** – Backend utama
- **Fiber** – Web framework Golang
- **PostgreSQL** – Penyimpanan data transaksi
- **OpenAI GPT-4.1** – Model AI untuk input, insight, dan percakapan
- **Docker** – Containerization dan deployment

---

## 🚀 Instalasi

```bash
git clone https://github.com/saufiroja/fin-ai.git
cd fin-ai

cp .env.example .env

docker-compose up --build
```

Akses aplikasi di: `http://localhost:8080`
Database migration berjalan otomatis saat container aktif.

---

## 📚 API Documentation

### 1. 🔐 Authentication

| Method | Endpoint                | Deskripsi     |
| ------ | ----------------------- | ------------- |
| POST   | `/api/v1/auth/register` | Register user |
| POST   | `/api/v1/auth/login`    | Login user    |

### 2. 👤 User

| Method | Endpoint                | Deskripsi     |
| ------ | ----------------------- | ------------- |
| GET    | `/api/v1/user`          | Get user info |
| PUT    | `/api/v1/user/:user_id` | Update user   |
| DELETE | `/api/v1/user/:user_id` | Delete user   |

### 3. 💸 Transactions

| Method | Endpoint                               | Deskripsi                |
| ------ | -------------------------------------- | ------------------------ |
| GET    | `/api/v1/transactions`                 | Get all transactions     |
| GET    | `/api/v1/transactions/:transaction_id` | Get transaction by ID    |
| POST   | `/api/v1/transactions`                 | Create new transaction   |
| PUT    | `/api/v1/transactions/:transaction_id` | Update transaction by ID |
| DELETE | `/api/v1/transactions/:transaction_id` | Delete transaction by ID |

### 4. 📊 Reports

| Method | Endpoint                     | Deskripsi                  |
| ------ | ---------------------------- | -------------------------- |
| GET    | `/api/v1/reports`            | Get monthly/yearly reports |
| GET    | `/api/v1/reports/:report_id` | Get report by ID           |
| POST   | `/api/v1/reports`            | Create new report          |
| PUT    | `/api/v1/reports/:report_id` | Update report by ID        |
| DELETE | `/api/v1/reports/:report_id` | Delete report by ID        |

### 5. 🧠 AI Features

| Method | Endpoint          | Deskripsi                                                  |
| ------ | ----------------- | ---------------------------------------------------------- |
| POST   | `/api/v1/ai/chat` | Endpoint terpadu untuk AI chat, OCR, prediksi, dan lainnya |

#### ✉️ Contoh Payload untuk `/api/v1/ai/chat`

```json
{
  "mode": "ocr", // Atau: "consultation", "summary", "prediction", "planner"
  "message": "Ini struk belanja bulan ini...",
  "file_base64": "..." // Opsional, digunakan jika mode = "ocr"
}
```

- `mode`: Menentukan jenis respons AI yang diminta.
- `message`: Prompt utama dari user.
- `file_base64`: Opsional. Base64 encoded string dari file struk belanja.

Jika `mode` tidak diisi, maka default dianggap sebagai "consultation".

---

## 🔧 Contributing

Pull request sangat diterima! Untuk perubahan besar, mohon buka issue terlebih dahulu untuk didiskusikan.

---
