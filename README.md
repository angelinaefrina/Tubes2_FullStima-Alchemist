# FullStima Alchemist
## Tugas Besar 2 - IF2211 Strategi Algoritma
### Pemanfaatan Algoritma BFS dan DFS dalam Pencarian Recipe pada Permainan Little Alchemy 2
![creators](https://github.com/angelinaefrina/Tubes2_FullStima-Alchemist/blob/main/doc/creators.jpg?raw=true)
## Kontributor
| Nama  | NIM | 
| ------------- | ------------- |
| Shannon Aurellius Anastasya Lie  | 13523019 |
| Angelina Efrina Prahastaputri | 13523060 |
| Sebastian Hung Yansen  | 13523070 | 

# Terkait Program
## **Element Recipe Finding in Little Alchemy 2**

Tugas besar ini berupa website desktop yang memungkinkan pengguna untuk mencari satu atau lebih resep untuk suatu elemen dalam permainan Little Alchemy 2 menggunakan algoritma penelusuran **BFS** dan **DFS**. Backend menggunakan **Golang** dan frontend menggunakan **Next.js**. Program hanya dapat dijalankan melalui browser dengan spesifikasi minimal **Mozilla 5.0**.

---

## Struktur Repository
```
|   README.md
|   
+---backend
|   |   algorithm.go
|   |   go.mod
|   |   go.sum
|   |   main.go
|   |   recipe.json
|   |   scrapper.go
|   |   server.js
|   |   
|   \---public
|       \---svgs //berisi .svg seluruh element dan default image
|                   
+---doc
|       laporan.txt
|       
\---frontend
    |   .eslintrc.json
    |   .gitignore
    |   jsconfig.json
    |   next.config.mjs
    |   package-lock.json
    |   package.json
    |   postcss.config.js
    |   postcss.config.mjs
    |   tailwind.config.js
    |   
    +---.next
    |   |   
    |   +---build
    |   |   \---chunks
    |   |           
    |   +---cache
    |   |       
    |   +---server
    |   |   |   
    |   |   +---app
    |   |   |   |   page.js
    |   |   |   |   page.js.map
    |   |   |   |   page_client-reference-manifest.js
    |   |   |   |   
    |   |   |   +---aboutpage
    |   |   |   |   |   page.js
    |   |   |   |   |   page.js.map
    |   |   |   |   |   page_client-reference-manifest.js
    |   |   |   |   |   
    |   |   |   |   \---page
```

## **Instalasi**

### **Langkah 1: Menyiapkan Backend**

1. **Clone atau Unduh Repositori:**
   - Jika menggunakan Git:
     ```bash
     git clone <URL_REPOSITORI>
     cd Tubes2_FullStima-Alchemist
     ```

2. **Jalankan Backend:**
   - Masuk ke direktori **backend** dan jalankan backend:
     ```bash
     cd backend
     go run .
     ```

### **Langkah 2: Menyiapkan Frontend**
1.  **Buka Terminal Baru**
   - Jangan tutup terminal backend, buka terminal baru untuk menjalankan frontend

2. **Instalansi Dependensi Node.js**
   - Masuk ke direktori **frontend** dan install dependensi:
     ```bash
     cd frontend
     npm install
     ```
2. **Jalankan Frontend:**
   - Setelah semua dependensi terinstall, jalankan server pengembangan Next.js:
     ```bash
     npm run dev
     ```
   - Aplikasi akan berjalan di `http://localhost:3000`.   

---

## **Penggunaan Aplikasi**

### **Frontend:**

1. **Menampilkan Pilihan Algoritma Pencarian Recipe**
   - Aplikasi frontend akan menampilkan pilihan pencarian resep, dapat menggunakan BFS maupun DFS serta dapat mencari satu resep maupun banyak resep

2. **Menampilkan Hasil Pencarian berupa Visualisasi Recipe**
   - Aplikasi frontend akan menampilkan visualisasi resep yang didapatkan berupa tree, dengan masing-masing simpul berisi nama dan gambar dari elemen penyusun

3. **Proses Pencarian:**
   - Setelah pengguna memasukkan nama elemen (dan jumlah resep yang diinginkan jika memilih banyak resep), backend akan melakukan penelusuran satu atau lebih resep sesuai permintaan pengguna

### **Backend:**

- Mengembalikan (kumpulan) resep sesuai algoritma yang dipakai

---

## **Penting:**

### **Menangani Kerentanan:**

Jika Anda melihat peringatan tentang kerentanan atau dependensi yang rentan saat menginstall dependensi, Anda bisa mengatasinya dengan menjalankan:
```bash
npm audit fix
npm audit fix --force
```

### **Spesifikasi Browser**
Program hanya dapat dijalankan melalui browser dengan spesifikasi minimal **Mozilla 5.0**.
