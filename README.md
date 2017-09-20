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

## Autentikasi

Autentikasi sign up dilakukan dengan memasukkan input request yang berupa nama, email, dan password pada JSON body untuk disimpan ke dalam database. Password disimpan dalam database dengan enkripsi bcrypt. Setelah melakukan sign up, user bisa melakukan login dengan memasukkan email dan password pada token di-generate secara random dan berbeda dengan masukan seed berupa waktu saat ini.

** File blueprint API : apiary.apib **
