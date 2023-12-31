package be_ksi

import (
	"encoding/json"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	credential Credential
	response Response
	user User
	billboard Billboard
	sewa Sewa
)

func SignUpHandler(MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	response.Status = false
	//
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(response)
	}
	err = SignUp(conn, user)
	if err != nil {
		response.Message = err.Error()
		return GCFReturnStruct(response)
	}
	//
	response.Status = true
	response.Message = user.NamaLengkap
	return GCFReturnStruct(response)
}

func LogInHandler(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	response.Status = false
	//
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(response)
	}
	user, err := LogIn(conn, user)
	if err != nil {
		response.Message = err.Error()
		return GCFReturnStruct(response)
	}
	tokenstring, err := Encode(user.ID, user.Email, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		response.Message = "Gagal Encode Token : " + err.Error()
		return GCFReturnStruct(response)
	}
	//
	credential.Message = "Selamat Datang " + user.Email
	credential.Token = tokenstring
	credential.Status = true
	return GCFReturnStruct(credential)
}

func GetProfileHandler(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	response.Status = false
	//
	payload, err := GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		response.Message = err.Error()
		return GCFReturnStruct(response)
	}
	user, err := GetUserFromID(payload.Id, conn)
	if err != nil {
		response.Message = err.Error()
		return GCFReturnStruct(response)
	}
	//
	response.Status = true
	response.Message = user.NamaLengkap
	return GCFReturnStruct(response)
}

func TambahBillboardHandler(MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	response.Status = false
	//
	user, err := GetUserLogin(os.Getenv("PASETOPUBLICKEYENV"), r)
	if err != nil {
		response.Message = err.Error()
		return GCFReturnStruct(response)
	}
	if user.Email != "admin@gmail.com" {
		response.Message = "Anda tidak memiliki akses"
		return GCFReturnStruct(response)
	}
	err = json.NewDecoder(r.Body).Decode(&billboard)
	if err != nil {
		response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(response)
	}
	err = TambahBillboardOlehAdmin(conn, billboard)
	if err != nil {
		response.Message = err.Error()
		return GCFReturnStruct(response)
	}
	//
	response.Status = true
	response.Message = billboard.Kode
	return GCFReturnStruct(response)
}

func GetBillboarHandler(MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	response.Status = false
	//
	id := GetID(r)
	if id == "" {
		billboard, err := GetAllBillboard(conn)
		if err != nil {
			response.Message = err.Error()
			return GCFReturnStruct(response)
		}
		//
		return GCFReturnStruct(billboard)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		response.Message = err.Error()
		return GCFReturnStruct(response)
	}
	billboard, err := GetBillboardFromID(idparam, conn)
	if err != nil {
		response.Message = err.Error()
		return GCFReturnStruct(response)
	}
	//
	response.Status = true
	return GCFReturnStruct(billboard)
}

func EditBillboardHandler(MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	response.Status = false
	//
	user, err := GetUserLogin(os.Getenv("PASETOPUBLICKEYENV"), r)
	if err != nil {
		response.Message = err.Error()
		return GCFReturnStruct(response)
	}
	if user.Email != "admin@gmail.com" {
		response.Message = "Anda tidak memiliki akses"
		return GCFReturnStruct(response)
	}
	id := GetID(r)
	if id == "" {
		response.Message = "Wrong parameter"
		return GCFReturnStruct(response)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		response.Message = "Invalid id parameter"
		return GCFReturnStruct(response)
	}
	err = json.NewDecoder(r.Body).Decode(&billboard)
	if err != nil {
		response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(response)
	}
	err = EditBillboardOlehAdmin(idparam, conn, billboard)
	if err != nil {
		response.Message = err.Error()
		return GCFReturnStruct(response)
	}
	//
	response.Status = true
	response.Message = billboard.Kode
	return GCFReturnStruct(response)
}

func HapusBillboardHandler(MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	response.Status = false
	//
	user, err := GetUserLogin(os.Getenv("PASETOPUBLICKEYENV"), r)
	if err != nil {
		response.Message = err.Error()
		return GCFReturnStruct(response)
	}
	if user.Email != "admin@gmail.com" {
		response.Message = "Anda tidak memiliki akses"
		return GCFReturnStruct(response)
	}
	id := GetID(r)
	if id == "" {
		response.Message = "Wrong parameter"
		return GCFReturnStruct(response)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		response.Message = "Invalid id parameter"
		return GCFReturnStruct(response)
	}
	err = HapusBillboardOlehAdmin(idparam, conn)
	if err != nil {
		response.Message = err.Error()
		return GCFReturnStruct(response)
	}
	//
	response.Status = true
	response.Message = "Berhasil menghapus billboard"
	return GCFReturnStruct(response)
}

//sewa
func SewaHandler(MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	response.Status = false
	//
	_, err := GetUserLogin(os.Getenv("PASETOPUBLICKEYENV"), r)
	if err != nil {
		response.Message = err.Error()
		return GCFReturnStruct(response)
	}
	err = json.NewDecoder(r.Body).Decode(&billboard)
	if err != nil {
		response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(response)
	}
	err = SewaBillboard(conn, sewa)
	if err != nil {
		response.Message = err.Error()
		return GCFReturnStruct(response)
	}
	//
	response.Status = true
	response.Message = billboard.Kode
	return GCFReturnStruct(response)
}