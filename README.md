# backendprojectbinar - dreamcatchertimche

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

Pada project ini, cloud yang digunakan adalah Heroku. Untuk unggah aplikasi ke heroku, daftar terlebih dahulu ke heroku lewat https://signup.heroku.com/. Kemudian install Heroku CLI melalui https://devcenter.heroku.com/articles/heroku-cli. Setelah Heroku CLI terinstall, ketikkan perintah berikut ini di command line untuk membuat aplikasi Heroku:
```
heroku create
```

Kemudian lakukan login:
```
heroku login
```

Setelah itu buat database postgresql di Heroku:
```
heroku addons:create heroku-postgresql
```

Database yang digunakan akan terdeteksi secara otomatis karena digunakan kode ```os.Getenv($DATABASE_URL)``` pada variabel ```db_url``` untuk mendeteksi URL database default yang digunakan pada Heroku.
Kemudian masukkan perintah ini untuk deploy repository ke Heroku:
```
git add . -A
git commit -m "Deploy Heroku"
git push heroku master
```

Setelah berhasil di-deploy ke Heroku, jalankan perintah ```heroku open``` untuk membuka URL tempat deploy aplikasi atau lakukan tes melalui Postman dengan menggunakan URL yang digunakan. 

## Pengujian dengan Postman

Untuk pengujian melalui Postman dilakukan dengan menggunakan <nama_URL_heroku>/v1/<nama_resource>. Adapun list dari resource yang dapat diakses adalah :
```
Authentication
POST /v1/register
POST /v1/login
POST /v1/logout

Posts
POST /v1/posts
GET /v1/posts
GET /v1/posts/{id_post}
PUT /v1/posts/{id_post}
DELETE /v1/posts/{id_post}

Categories
GET /v1/categories
GET /v1/categories/showposts

Profile
GET /v1/profile
PUT /v1/profile
PATCH /v1/profile
GET /v1/profile/posts
POST /v1/profile/categories

Bookmarks
POST /v1/bookmarks
GET /v1/bookmarks
DELETE /v1/bookmarks/{id_bookmark}

Comments
POST /v1/posts/{id_post}/comments
GET /v1/posts/{id_post}/comments
GET /v1/posts/{id_post}/comments/{id_comment}
PUT /v1/posts/{id_post}/comments/{id_comment}
DELETE /v1/posts/{id_post}/comments/{id_comment}
```

# Penjelasan Aplikasi

## Autentikasi

Autentikasi sign up dilakukan dengan memasukkan input request yang berupa nama, email, dan password pada JSON body untuk disimpan ke dalam database. Password disimpan dalam database dengan enkripsi bcrypt. Setelah melakukan sign up, user bisa melakukan login dengan memasukkan email dan password pada input request dengan keluaran berupa token yang di-generate secara random.

## Operasi CRUD pada Posts, Users, Bookmarks, Categories, dan Comments

Operasi CRUD (create, read, update, dan delete) pada tabel posts, users, bookmarks, categories, dan comments dilakukan dengan menggunakan token yang didapatkan saat login. Token tersebut dimasukkan ke header bagian Authorization.

** File blueprint API : apiary.apib **
