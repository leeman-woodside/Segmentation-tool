package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
)

var folderStructure = make(map[string][]string)

type Image struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ImgName string             `json:"imgName,omitempty" bson:"imgName,omitempty"`
	ImgType string             `json:"imgType,omitempty" bson:"imgType,omitempty"`
	ImgSize int                `json:"imgSize,omitempty" bson:"imgSize,omitempty"`
}

var client *mongo.Client

func CreateImageEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var image Image
	json.NewDecoder(request.Body).Decode(&image)
	collection := client.Database("Segmentation_Web_App").Collection("images")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, image)
	json.NewEncoder(response).Encode(result)
}

func GetImagesEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var images []Image
	collection := client.Database("Segmentation_Web_App").Collection("images")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var image Image
		cursor.Decode(&image)
		images = append(images, image)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(images)
}

func GetImageEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var image Image
	collection := client.Database("Segmentation_Web_App").Collection("images")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Image{ID: id}).Decode(&image)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(image)
}

func UploadFile(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	uploadType := params["type"]
	err := request.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		fmt.Println("err 1", err)
	}

	file, handler, err := request.FormFile(uploadType)
	if err != nil {
		fmt.Println("err 2", err)
		return
	}
	defer file.Close()

	err = os.MkdirAll("../ui/dist/"+uploadType+"/"+filepath.Dir(handler.Filename), 0755)
	if err != nil {
		fmt.Println("err 2.5", err)
	}

	resFile, err := os.Create("../ui/dist/" + uploadType + "/" + handler.Filename)
	if err != nil {
		fmt.Println("err 3", err)
	}
	defer resFile.Close()

	io.Copy(resFile, file)
	defer resFile.Close()

	fmt.Println("File Info")
	fmt.Println("File Name:", handler.Filename)
	fmt.Println("File Size:", handler.Size)
	fmt.Println("File Type:", handler.Header)

}

func GetFolderLocation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fileType := params["type"]
	w.Header().Add("content-type", "application/json")
	var currentFolder = ""
	err := filepath.Walk("../ui/dist/"+fileType+"/", func(path string, info os.FileInfo, err error) error {
		if path == "../ui/dist/"+fileType+"/" {
			return nil
		}

		fi, err := os.Stat(path)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		fileBase := filepath.Base(path)

		mode := fi.Mode()
		if mode.IsDir() {
			folderStructure[fileBase] = []string{}
			currentFolder = fileBase
			// do directory stuff
			fmt.Println("directory")
		}
		// case mode.IsRegular():
		// 	// do file stuff
		// 	folderStructure[currentFolder] = append(folderStructure[currentFolder], fileBase)
		// 	fmt.Println("file")

		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		fmt.Println("visited file or dir: %q\n", path)
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", err)
		return
	}
	fmt.Println("FINAL", folderStructure)
	output, _ := json.Marshal(folderStructure)
	w.Write([]byte(output))

}

func GetFileLocation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	folder := params["type"]
	var files []string
	w.Header().Add("content-type", "application/json")
	//iterate thru a certain number of files for the first Get and
	//then only load the next unloaded image as masks are saved
	files = append(files, "../ui/dist/images/"+folder+"/")
	folderStructure[folder] = append(folderStructure[folder], files)

}

func main() {
	fmt.Println("Starting the application...")
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// client, _ = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	// router.HandleFunc("/image", CreateImageEndpoint).Methods("POST")
	// router.HandleFunc("/images", GetImagesEndpoint).Methods("GET")
	// router.HandleFunc("/image/{id}", GetImageEndpoint).Methods("GET")
	router := mux.NewRouter()
	router.HandleFunc("/upload/{type}", UploadFile).Methods("POST")
	router.HandleFunc("/imageLocation/{type}", GetFolderLocation).Methods("GET")
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	static := negroni.NewStatic(http.Dir("../ui/dist"))
	n.Use(static)
	n.UseHandler(router)
	http.ListenAndServe(":8082", n)
}
