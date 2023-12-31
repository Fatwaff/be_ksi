package be_ksi

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	NamaLengkap     string             `bson:"namalengkap,omitempty" json:"namalengkap,omitempty"`
	Email           string             `bson:"email,omitempty" json:"email,omitempty"`
	Password        string             `bson:"password,omitempty" json:"password,omitempty"`
	Confirmpassword string             `bson:"confirmpass,omitempty" json:"confirmpass,omitempty"`
	NoHp			string             `bson:"nohp,omitempty" json:"nohp,omitempty"`
	KTP				string             `bson:"ktp,omitempty" json:"ktp,omitempty"`
	Salt            string             `bson:"salt,omitempty" json:"salt,omitempty"`
}

type Billboard struct {
	ID 				primitive.ObjectID 	`bson:"_id,omitempty" json:"_id,omitempty"`
	Kode 			string				`bson:"kode,omitempty" json:"kode,omitempty"`
	Nama 			string 				`bson:"nama,omitempty" json:"nama,omitempty"`
	Gambar 			string 				`bson:"gambar,omitempty" json:"gambar,omitempty"`
	Panjang 		float64 			`bson:"panjang,omitempty" json:"panjang,omitempty"`
	Lebar 			float64 			`bson:"lebar,omitempty" json:"lebar,omitempty"`
	Latitude  		float64 			`bson:"latitude,omitempty" json:"latitude,omitempty"`
	Longitude 		float64 			`bson:"longitude,omitempty" json:"longitude,omitempty"`
	Lokasi 			string 				`bson:"lokasi,omitempty" json:"lokasi,omitempty"`
}

type Sewa struct {
	ID 				primitive.ObjectID 	`bson:"_id,omitempty" json:"_id,omitempty"`
	Billboard		Billboard			`bson:"billboard,omitempty" json:"billboard,omitempty"`
	User			User				`bson:"user,omitempty" json:"user,omitempty"`
	Content			string				`bson:"content,omitempty" json:"content,omitempty"`
	TanggalMulai	time.Time			`bson:"tanggal_mulai,omitempty" json:"tanggal_mulai,omitempty"`
	TanggalSelesai	time.Time			`bson:"tanggal_selesai,omitempty" json:"tanggal_selesai,omitempty"`
	Status			bool				`bson:"status,omitempty" json:"status,omitempty"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Response struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Payload struct {
	Id   		primitive.ObjectID `json:"id"`
	Email 		string             `json:"email"`
	Exp  		time.Time          `json:"exp"`
	Iat  		time.Time          `json:"iat"`
	Nbf  		time.Time          `json:"nbf"`
}