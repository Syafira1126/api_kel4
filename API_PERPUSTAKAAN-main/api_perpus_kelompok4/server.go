package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Buku struct {
	Id_Buku      string `json:"id"`
	Judul_Buku   string `json:"jb"`
	Id_Pengarang string `json:"idp"`
	Id_Penerbit  string `json:"id_p"`
	Tahun_terbit string `json:"t_te"`
	gambar_buku  string `json "g_buku"`
	Id_Genre     string `json:"id_genre"`
	Jumlah_buku  string `json:"jmh_bu"`
	Sinopsis     string `json:"ss"`
}
type Genre struct {
	Id_Genre        string `json:"id_genre"`
	Genre           string `json:"genre"`
	Deskripsi_Genre string `json:"deskripsi_genre"`
}
type Peminjaman struct {
	Id_Peminjaman        string `json:"id_peminjaman"`
	Id_Buku              string `json:"id_buku"`
	Tanggal_peminjaman   string `json:"tanggal_peminjaman"`
	Tanggal_pengembalian string `json:"tanggal_pengembalian"`
}
type Penerbit struct {
	Id_Penerbit string `json:"id_penerbit"`
	Penerbit    string `json:"penerbit"`
}
type Pengarang struct {
	Id_Pengarang   string `json:"id_pengarang"`
	Nama_Pengarang string `json:"nama_pengarang"`
}
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {

	// database connection
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_perpustakaan")
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}
	// database connection

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Service API!")
	})

	// Buku
	e.GET("/buku", func(c echo.Context) error {
		res, err := db.Query("CALL show_elmu")

		defer res.Close()

		if err != nil {
			log.Fatal(err)
		}
		var buku []Buku
		for res.Next() {
			var m Buku
			_ = res.Scan(&m.Id_Buku, &m.Judul_Buku, &m.Id_Pengarang, &m.Tahun_terbit, &m.Sinopsis, &m.gambar_buku)
			buku = append(buku, m)
		}

		return c.JSON(http.StatusOK, buku)
	})

	e.POST("/buku/add", func(c echo.Context) error {
		var buku Buku
		c.Bind(&buku)

		sqlStatement := "CALL add_buku"
		res, err := db.Query(sqlStatement, c.Param("id"), c.Param("jb"), c.Param("idp"), c.Param("id_p"), c.Param("t_ter"), c.Param("id_genre"), c.Param("jmh_bu"), c.Param("ss"))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
			return c.JSON(http.StatusCreated, buku)
		}
		return c.String(http.StatusOK, "ok")
	})

	e.PUT("/buku/:id_buku", func(c echo.Context) error {
		var buku Buku
		c.Bind(&buku)

		sqlStatement := "CALL edit_buku"
		res, err := db.Query(sqlStatement, c.Param("jb"), c.Param("idp"), c.Param("id_p"), c.Param("t_ter"), c.Param("id_genre"), c.Param("jmh_bu"), c.Param("ss"), c.Param("id_buku"))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
			return c.JSON(http.StatusCreated, buku)
		}
		return c.String(http.StatusOK, "ok")
	})

	// Buku

	// Peminjaman
	e.GET("/peminjaman", func(c echo.Context) error {
		res, err := db.Query("SELECT * FROM peminjaman")

		defer res.Close()

		if err != nil {
			log.Fatal(err)
		}
		var peminjaman []Peminjaman
		for res.Next() {
			var m Peminjaman
			_ = res.Scan(&m.Id_Peminjaman, &m.Id_Buku, &m.Tanggal_peminjaman, &m.Tanggal_pengembalian)
			peminjaman = append(peminjaman, m)
		}

		return c.JSON(http.StatusOK, peminjaman)
	})

	e.POST("/peminjaman/add", func(c echo.Context) error {
		var peminjaman Peminjaman
		c.Bind(&peminjaman)

		sqlStatement := "CALL mau_pinjam (?, ?, ?)"
		res, err := db.Query(sqlStatement, peminjaman.Id_Buku, peminjaman.Tanggal_peminjaman, peminjaman.Tanggal_pengembalian)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
			return c.JSON(http.StatusCreated, peminjaman)
		}
		return c.String(http.StatusOK, "ok")
	})
	// Peminjaman

	// User
	e.GET("/login/user/:id/:pw", func(c echo.Context) error {
		req := new(User)
		if err := c.Bind(req); err != nil {
			return err
		}

		sqlStatement := "CALL valid_user(?, ?)"
		res, err := db.Query(sqlStatement, c.Param("id"), c.Param("pw"))
		defer res.Close()

		if err != nil {
			log.Fatal(err)
		}

		var users []User
		for res.Next() {
			var u User
			err := res.Scan(&u.Username, &u.Password)
			if err != nil {
				log.Fatal(err)
			}
			users = append(users, u)
		}

		return c.JSON(http.StatusOK, users)
	})

	e.POST("/user/add", func(c echo.Context) error {
		var user User
		c.Bind(&user)
		sqlStatement := "CALL add_user(?, ?)"
		res, err := db.Query(sqlStatement, user.Username, user.Password)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
			return c.JSON(http.StatusCreated, user)
		}
		return c.String(http.StatusOK, "ok")
	})
	// User

	e.Logger.Fatal(e.Start(":1231"))
}
