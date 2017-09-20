# dreamcatchertimche

Repository ini digunakan untuk membuat kode backend dari aplikasi Dreamcatcher yang digunakan pada project Binar Batch 3 di Tim C. Pada kode ini terdapat API untuk melakukan autentikasi (sign up, login, dan logout) dan operasi CRUD yang hanya pada tabel posts, users, categories, bookmarks dan comments yang digunakan oleh tim Front End (Android & iOS) untuk menampilkan data di sisi client. 

## Instalasi

Install Golang terlebih dahulu kemudian setting $GOPATH. Adapun petunjuk untuk instalasi Golang dan setting $GOPATH bisa dilihat di https://golang.org/doc/install.

Jalankan perintah ini di command line Linux/Unix:
```
cd $GOPATH/src
```

Atau pada Windows:
```
cd %GOPATH%/src
```

Kemudian ketikkan perintah ini pada Windows/Linux/Unix untuk membuat folder baru tempat menaruh project:
```
mkdir dreamcatchertimche
cd dreamcatchertimche
```

Setelah itu, download seluruh dependency untuk kode yang digunakan dengan menggunakan godeps. Masukkan perintah ````go get -u github.com/tools/godep```` untuk instalasi/update godep. Supaya godep dapat menentukan, menyimpan, dan menulis ulang dependency dari aplikasi yang digunakan, masuk ke direktori aplikasi yang berada di ````$GOPATH/src/dreamcatchertimche```` dan jalankan perintah ini di command line:
```
godep save ./...
```

Dari perintah tersebut didapat file ```Godeps/Godeps.json``` yang berisi representasi JSON dari dependency aplikasi beserta kopian dependency yang disimpan di subdirectory ```vendor/```.

Setelah itu, aplikasi bisa dijalankan dengan mengetikkan perintah ```go run main.go``` di command line. 

## Deploy ke Cloud

Pada 

## Autentikasi

Autentikasi sign up dilakukan dengan memasukkan input request yang berupa nama, email, dan password pada JSON body untuk disimpan ke dalam database. Password disimpan dalam database dengan enkripsi bcrypt. Setelah melakukan sign up, user bisa melakukan login dengan memasukkan email dan password pada token di-generate secara random dan berbeda dengan masukan seed berupa waktu saat ini.

## Rancangan Database MySQL

Database yang digunakan adalah binar_backend_test.sql di mana terdapat empat tabel, yaitu tabel data pengguna, soal ujian, jenis ujian, dan jawaban ujian. Adapun keterangan detail terkait kolom dan tipe data yang digunakan dapat dilihat dengan melakukan import file database binar_backend_test.sql.

## Rancangan RESTful API Sederhana

API yang akan dibuat dirancang terlebih dahulu untuk mempermudah proses pembuatan API dengan menggunakan Apiary (apiary.io). Pada API ini, keamanan komunikasi data dilakukan dengan autentikasi terlebih dahulu dengan OAuth 2.0 atau JWT. Setelah pengguna mendapatkan token autentikasi, pengguna dapat mengakses API berdasarkan privilege yang dimiliki. Operasi yang dapat dilakukan melalui API ini antara lain adalah :
- Operasi create, read, delete, dan update (CRUD) pada data pengguna dan soal ujian dengan privilege admin
- Operasi autentikasi user dan pemilihan serta penyimpanan jawaban ujian

** File blueprint API : apiary.apib **
