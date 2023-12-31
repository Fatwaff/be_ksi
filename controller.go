package be_ksi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/badoux/checkmail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/argon2"
)

// mongo
func MongoConnect(MongoString, dbname string) *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv(MongoString)))
	if err != nil {
		fmt.Printf("MongoConnect: %v\n", err)
	}
	return client.Database(dbname)
}

// crud
func GetAllDocs(db *mongo.Database, col string, docs interface{}) interface{} {
	collection := db.Collection(col)
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error GetAllDocs %s: %s", col, err)
	}
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		return err
	}
	return docs
}

func InsertOneDoc(db *mongo.Database, col string, doc interface{}) (insertedID primitive.ObjectID, err error) {
	result, err := db.Collection(col).InsertOne(context.Background(), doc)
	if err != nil {
		return insertedID, fmt.Errorf("kesalahan server : insert")
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

func UpdateOneDoc(id primitive.ObjectID, db *mongo.Database, col string, doc interface{}) (err error) {
	filter := bson.M{"_id": id}
	result, err := db.Collection(col).UpdateOne(context.Background(), filter, bson.M{"$set": doc})
	if err != nil {
		return fmt.Errorf("error update: %v", err)
	}
	if result.ModifiedCount == 0 {
		err = fmt.Errorf("tidak ada data yang diubah")
		return
	}
	return nil
}

func DeleteOneDoc(_id primitive.ObjectID, db *mongo.Database, col string) error {
	collection := db.Collection(col)
	filter := bson.M{"_id": _id}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", _id, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", _id)
	}

	return nil
}

// get user
func GetUserFromID(_id primitive.ObjectID, db *mongo.Database) (doc User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
}

func GetUserFromEmail(email string, db *mongo.Database) (doc User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"email": email}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("email tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

// get user login
func GetUserLogin(PASETOPUBLICKEYENV string, r *http.Request) (Payload, error) {
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		return payload, err
	}
	return payload, nil
}

// get id
func GetID(r *http.Request) string {
    return r.URL.Query().Get("id")
}

// return struct
func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

// register
func SignUp(db *mongo.Database, insertedDoc User) error {
	if insertedDoc.NamaLengkap == "" || insertedDoc.Email == "" || insertedDoc.Password == "" || insertedDoc.NoHp == "" || insertedDoc.KTP == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	if err := checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return fmt.Errorf("email tidak valid")
	}
	userExists, _ := GetUserFromEmail(insertedDoc.Email, db)
	if insertedDoc.Email == userExists.Email {
		return fmt.Errorf("email sudah terdaftar")
	}
	if insertedDoc.Confirmpassword != insertedDoc.Password {
		return fmt.Errorf("konfirmasi password salah")
	}
	if strings.Contains(insertedDoc.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Password) < 8 {
		return fmt.Errorf("password terlalu pendek")
	}
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Password), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"namalengkap": insertedDoc.NamaLengkap,
		"email":    insertedDoc.Email,
		"password": hex.EncodeToString(hashedPassword),
		"nohp":     insertedDoc.NoHp,
		"ktp":      insertedDoc.KTP,
		"salt":     hex.EncodeToString(salt),
	}
	_, err = InsertOneDoc(db, "user", user)
	if err != nil {
		return err
	}
	return nil
}

// login
func LogIn(db *mongo.Database, insertedDoc User) (user User, err error) {
	if insertedDoc.Email == "" || insertedDoc.Password == "" {
		return user, fmt.Errorf("mohon untuk melengkapi data")
	}
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return user, fmt.Errorf("email tidak valid")
	}
	existsDoc, err := GetUserFromEmail(insertedDoc.Email, db)
	if err != nil {
		return
	}
	salt, err := hex.DecodeString(existsDoc.Salt)
	if err != nil {
		return user, fmt.Errorf("kesalahan server : salt")
	}
	hash := argon2.IDKey([]byte(insertedDoc.Password), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hash) != existsDoc.Password {
		return user, fmt.Errorf("password salah")
	}
	return existsDoc, nil
}

// billboard
func CheckLatitudeLongitude(db *mongo.Database, insertedDoc Billboard) bool {
	collection := db.Collection("billboard")
	filter := bson.M{"latitude": insertedDoc.Latitude, "longitude": insertedDoc.Longitude}
	err := collection.FindOne(context.Background(), filter).Decode(&Billboard{})
	return err == nil
}

func TambahBillboardOlehAdmin(db *mongo.Database, insertedDoc Billboard) error {
	if insertedDoc.Kode == "" ||  insertedDoc.Panjang == 0 || insertedDoc.Lebar == 0 || insertedDoc.Latitude == 0 || insertedDoc.Longitude == 0 || insertedDoc.Lokasi == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	if CheckLatitudeLongitude(db, insertedDoc) {
		return fmt.Errorf("billboard sudah terdaftar")
	}
	_, err := InsertOneDoc(db, "billboard", insertedDoc)
	if err != nil {
		return err
	}
	return nil
}

func GetAllBillboard(db *mongo.Database) (docs []Billboard, err error) {
	collection := db.Collection("billboard")
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return docs, fmt.Errorf("error GetAllBillboard: %s", err)
	}
	err = cursor.All(context.Background(), &docs)
	if err != nil {
		return docs, err
	}
	return docs, nil
}

func GetBillboardFromID(_id primitive.ObjectID, db *mongo.Database) (doc Billboard, err error) {
	collection := db.Collection("billboard")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.Background(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
}

func EditBillboardOlehAdmin(_id primitive.ObjectID, db *mongo.Database, insertedDoc Billboard) error {
	if insertedDoc.Kode == "" || insertedDoc.Panjang == 0 || insertedDoc.Lebar == 0 || insertedDoc.Latitude == 0 || insertedDoc.Longitude == 0 || insertedDoc.Lokasi == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	err := UpdateOneDoc(_id, db, "billboard", insertedDoc)
	if err != nil {
		return err
	}
	return nil
}

func HapusBillboardOlehAdmin(_id primitive.ObjectID, db *mongo.Database) error {
	err := DeleteOneDoc(_id, db, "billboard")
	if err != nil {
		return err
	}
	return nil
}

// sewa
func SewaBillboard(db *mongo.Database, insertedDoc Sewa) error {
	if insertedDoc.Billboard.ID == primitive.NilObjectID || insertedDoc.User.ID == primitive.NilObjectID || insertedDoc.Content == "" || insertedDoc.TanggalMulai.IsZero() || insertedDoc.TanggalSelesai.IsZero() {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	if CheckSewa(db, insertedDoc) {
		return fmt.Errorf("billboard sudah disewa")
	}
	_, err := InsertOneDoc(db, "sewa", insertedDoc)
	if err != nil {
		return err
	}
	return nil
}

func CheckSewa(db *mongo.Database, insertedDoc Sewa) bool {
	collection := db.Collection("sewa")
	filter := bson.M{"sewa": insertedDoc.ID}
	err := collection.FindOne(context.Background(), filter).Decode(&Sewa{})
	return err == nil
}