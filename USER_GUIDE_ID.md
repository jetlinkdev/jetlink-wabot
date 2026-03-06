# Panduan Pengguna Bot WhatsApp AI

## Apa itu Bot WhatsApp AI?

Bot WhatsApp AI adalah asisten pintar yang terintegrasi dengan WhatsApp untuk menjawab pertanyaan Anda menggunakan kecerdasan buatan (AI). Bot ini menggunakan teknologi LLM (Large Language Model) dari Groq dengan model Qwen/Qwen3-32B.

## Fitur Utama

- 💬 **Chat AI**: Mengobrol dengan AI seperti ChatGPT
- 🧠 **Konteks Percakapan**: Mengingat riwayat chat Anda
- ⚡ **Respon Cepat**: Jawaban dalam hitungan detik
- 🔧 **Perintah Khusus**: Berbagai perintah untuk kontrol
- 📱 **Mudah Digunakan**: Cukup kirim pesan seperti biasa

## Cara Menggunakan

### 1. Mulai Mengobrol

Cukup kirim pesan ke nomor bot WhatsApp, dan bot akan menjawab:

```
Halo, apa kabar?
```

Bot akan menjawab dengan ramah!

### 2. Perintah yang Tersedia

#### /help - Tampilkan Bantuan

Menampilkan daftar semua perintah yang tersedia.

**Contoh:**
```
/help
```

**Respon:**
```
*Bot Commands:*
/help - Show this help message
/clear - Clear chat history
/status - Show chat status
/context <number> - Set context window size (e.g., /context 5)

Just send a message to chat with the AI assistant!
```

---

#### /clear - Hapus Riwayat Chat

Menghapus semua riwayat percakapan Anda. Berguna jika Anda ingin memulai percakapan baru dari awal.

**Contoh:**
```
/clear
```

**Respon:**
```
Chat history cleared! 🧹
```

**Kapan digunakan:**
- Ingin memulai topik baru
- Riwayat sudah terlalu panjang
- Ada informasi sensitif yang ingin dihapus

---

#### /status - Lihat Status Chat

Menampilkan informasi tentang sesi chat Anda:
- Ukuran context window
- Jumlah pesan dalam riwayat
- Status sesi

**Contoh:**
```
/status
```

**Respon:**
```
*Chat Status:*
- Context window: 10 messages
- Messages in history: 25
- Status: Active ✅
```

---

#### /context <angka> - Atur Ukuran Context

Mengatur berapa banyak pesan terakhir yang akan diingat oleh bot.

**Contoh:**
```
/context 5
```

**Respon:**
```
Context window set to 5 messages
```

**Ukuran yang direkomendasikan:**
- **3-5 pesan**: Untuk percakapan singkat, hemat token
- **10 pesan**: Default, cocok untuk kebanyakan kasus
- **15-20 pesan**: Untuk percakapan kompleks yang butuh konteks panjang
- **Maksimal**: 50 pesan

**Catatan:** Semakin besar context, semakin banyak token yang digunakan.

---

## Contoh Penggunaan

### 1. Bertanya Informasi Umum

```
Anda: Apa ibukota Indonesia?

Bot: Ibukota Indonesia adalah Jakarta. Namun, pemerintah Indonesia sedang 
dalam proses memindahkan ibukota negara ke Nusantara (IKN) yang terletak 
di Kalimantan Timur.
```

### 2. Minta Bantuan Coding

```
Anda: Buatkan fungsi Python untuk menghitung faktorial

Bot: Tentu! Berikut fungsi Python untuk menghitung faktorial:

def faktorial(n):
    if n == 0 or n == 1:
        return 1
    return n * faktorial(n - 1)

# Contoh penggunaan
print(faktorial(5))  # Output: 120
```

### 3. Diskusi dengan Konteks

```
Anda: Saya suka makan nasi goreng

Bot: Wah, nasi goreng memang enak! Nasi goreng adalah salah satu makanan 
khas Indonesia yang populer.

Anda: Apa resepnya?

Bot: Berdasarkan kesukaan Anda pada nasi goreng, berikut resep sederhana:

Bahan:
- 1 piring nasi putih
- 2 siung bawang putih
- 3 butir bawang merah
- 2 butir telur
- Kecap manis, garam, merica
...
```

Bot mengingat percakapan sebelumnya berkat context window!

### 4. Terjemahan

```
Anda: Terjemahkan "Good morning, how are you?" ke Bahasa Indonesia

Bot: "Selamat pagi, apa kabar?"
```

### 5. Ringkasan Teks

```
Anda: [paste panjang teks]
Anda: Ringkas teks di atas

Bot: [memberikan ringkasan]
```

## Tips Penggunaan

### ✅ Do (Lakukan)

1. **Gunakan kalimat jelas**: Semakin jelas pertanyaan, semakin baik jawaban
2. **Gunakan /clear**: Jika ingin ganti topik sepenuhnya
3. **Atur context**: Sesuaikan dengan kebutuhan percakapan
4. **Simpan info penting**: Bot tidak menyimpan data permanen

### ❌ Don't (Jangan)

1. **Jangan share info sensitif**: Password, data pribadi, dll.
2. **Jangan spam**: Beri waktu bot untuk merespon
3. **Jangan harap 100% akurat**: AI bisa membuat kesalahan

## Troubleshooting

### Bot Tidak Menjawab

**Kemungkinan penyebab:**
- Koneksi internet bermasalah
- Server bot sedang down
- Token API habis

**Solusi:**
1. Cek koneksi internet
2. Tunggu beberapa menit
3. Hubungi administrator

### Respon Lama

**Penyebab:**
- Server Groq sedang sibuk
- Context window terlalu besar
- Pertanyaan kompleks

**Solusi:**
1. Kurangi ukuran context: `/context 5`
2. Tanyakan dengan lebih spesifik
3. Tunggu sebentar

### Jawaban Tidak Relevan

**Penyebab:**
- Context window terlalu kecil
- Pertanyaan ambigu

**Solusi:**
1. Besarkan context: `/context 15`
2. Ulangi pertanyaan dengan lebih jelas
3. Gunakan `/clear` dan mulai baru

## Pertanyaan Umum (FAQ)

### Q: Apakah percakapan saya disimpan?
**A:** Ya, percakapan disimpan di database SQLite untuk menjaga context. Anda bisa menghapusnya kapan saja dengan `/clear`.

### Q: Apakah bot ini gratis?
**A:** Tergantung konfigurasi administrator. Bot menggunakan Groq API yang memiliki quota gratis dan berbayar.

### Q: Berapa lama percakapan disimpan?
**A:** Percakapan disimpan sampai Anda menghapusnya dengan `/clear` atau sampai context window penuh.

### Q: Bisa bot mengirim gambar/file?
**A:** Saat ini bot hanya mendukung teks. Fitur gambar/file mungkin ditambahkan di masa depan.

### Q: Apakah bot bisa dipanggil 24/7?
**A:** Ya, bot berjalan terus menerus. Namun, ada rate limiting dari Groq API.

### Q: Berapa maksimal panjang pesan?
**A:** Bot mendukung pesan hingga 60,000 karakter, tapi untuk performa optimal, gunakan pesan yang lebih singkat.

## Keamanan & Privasi

### Data yang Disimpan

- **Nomor WhatsApp**: Untuk identifikasi sesi
- **Riwayat Chat**: Untuk context window
- **Timestamp**: Waktu pesan dikirim

### Data yang TIDAK Disimpan

- Informasi login
- Data pembayaran
- File atau media

### Tips Keamanan

1. Jangan share password atau informasi sensitif
2. Gunakan `/clear` secara berkala
3. Jangan percaya 100% pada informasi dari AI

## Dukungan

Jika mengalami masalah atau punya pertanyaan:
1. Cek dokumentasi ini
2. Hubungi administrator bot
3. Laporkan bug jika ditemukan

## Update & Perubahan

Bot ini terus dikembangkan. Fitur baru akan ditambahkan secara berkala. Untuk melihat perubahan, cek file `CHANGELOG.md`.

---

**Selamat mengobrol dengan Bot WhatsApp AI! 🎉**
