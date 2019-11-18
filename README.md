# backendprojectbinar - dreamcatchertimche

Repository ini digunakan untuk membuat kode backend dari aplikasi Dreamcatcher yang digunakan pada project Binar Batch 3 di Tim C. Pada kode ini terdapat API untuk melakukan autentikasi (sign up, login, dan logout) dan operasi CRUD yang hanya pada tabel posts, users, categories, bookmarks dan comments yang digunakan oleh tim Front End (Android & iOS) untuk menampilkan data di sisi client. 

[![license](https://img.shields.io/github/license/mashape/apistatus.svg)]()

## Instalasi

Install Golang terlebih dahulu kemudian setting $GOPATH. Adapun petunjuk untuk instalasi Golang dan setting $GOPATH bisa dilihat di https://golang.org/doc/install.

<p align="center">![Peek recording itself](https://github.com/cyberboyz/backendprojectbinar/blob/master/asset/runningdreamcatcher.gif)</p>

Jalankan perintah go get :
```
go get github.com/cyberboyz/backendprojectbinar
```

Kemudian masuk directory dengan ketik di command line Linux/Unix:
```
cd $GOPATH/src/github.com/cyberboyz/backendprojectbinar
```

Atau pada Windows:
```
cd %GOPATH%\src\github.com\cyberboyz\backendprojectbinar
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

Untuk pengujian melalui Postman dilakukan dengan menggunakan <your_url>/v1/<nama_resource>. Adapun list dari resource yang dapat diakses adalah :

| Name                  | URL                                | HTTP Method  |
| ----------------------|:----------------------------------:|:------------:|
| Register User         | `<your_url>/v1/register`           |   **POST**   |
| Login User            | `<your_url>/v1/login`              |   **POST**   |
| Logout User           | `<your_url>/v1/logout`             |   **GET**    |
| Get All Post          | `<your_url>/v1/posts`              |   **GET**    |
| Get Post Detail       | `<your_url>/v1/posts/<id_post>`    |   **GET**    |
| Update Post           | `<your_url>/v1/posts/<id_post>`    |   **PUT**    |
| Delete Post           | `<your_url>/v1/posts/<id_post>`    |   **DELETE** |
| Show All Users        | `<your_url>/v1/profile`            |   **GET**    |
| Show Profile Detail   | `<your_url>/v1/profile/<id_user>`  |   **GET**    |
| Update Profile Detail | `<your_url>/v1/profile/<id_user>`  |   **PUT**    |
| Delete Profile        | `<your_url>/v1/profile/<id_user>`  |   **DELETE** |
| Add Category          | `<your_url>/v1/categories`         |   **POST**   |
| Show All Categories   | `<your_url>/v1/categories`         |   **GET**    |
| Show All Posts Based on Categories | `<your_url>/v1/categories`         |   **GET**    |
| Show All Posts Based on Several Categories| `<your_url>/v1/3categoriesposts` |   **POST** |
| Add Bookmark          | `<your_url>/v1/bookmarks`          |   **POST**   |
| Delete Bookmark       | `<your_url>/v1/bookmarks/<id_post>`|   **DELETE** |
| Show Own Bookmarks    | `<your_url>/v1/bookmarks`          |   **GET**    |
| Show Own Profile      | `<your_url>/v1/ownprofile`         |   **GET**    |
| Show Own Posts        | `<your_url>/v1/ownposts`           |   **GET**    |
| Add Categories by User| `<your_url>/v1/owncategory`        |   **POST**   |
| Update Categories by User| `<your_url>/v1/owncategory`     |   **PUT**    |
| Delete Categories by User| `<your_url>/v1/owncategory`     |   **DELETE** |
| Get All Posts Based on User| `<your_url>/v1/profile/<id_user>/posts` |   **GET** |

# Penjelasan Aplikasi

## Autentikasi

Autentikasi sign up dilakukan dengan memasukkan input request yang berupa nama, email, dan password pada JSON body untuk disimpan ke dalam database. Password disimpan dalam database dengan enkripsi bcrypt. Setelah melakukan sign up, user bisa melakukan login dengan memasukkan email dan password pada input request dengan keluaran berupa token yang di-generate secara random.

## Operasi CRUD pada Posts, Users, Bookmarks, Categories, dan Comments

Operasi CRUD (create, read, update, dan delete) pada tabel posts, users, bookmarks, categories, dan comments dilakukan dengan menggunakan token yang didapatkan saat login. Token tersebut dimasukkan ke header bagian Authorization.

** File blueprint API : apiary.apib **

# Programmer
- Fattah fattahazzuhry@gmail.com
- Riska rizkawidarsono29@gmail.com

# Library yang Digunakan
- Gin
- BCrypt

# Credits
- Binar Academy <3
- Mentor Backend Binar Academy Batch #3 : mas Prima, mas Gean, mas Andi, mas Estu
