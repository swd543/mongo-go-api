# File Upload With Go HTTP Server

See Mongo API with Go Lang: [Part 1](http://mycodesmells.com/post/mongodb-api-with-go-ang). 

Once we finish our little HTTP server with REST API for getting, saving and deleting simple database objects, it's time to take a look on file management. In this post we create a similar API for uploading and deleting files that are related to existing database model.

### Changes in our API 

First, we extend our database Item model:

    type Item struct {
        ID     string   `json:"id" bson:"_id,omitempty"`
        Value  int      `json:"value"`
        Images []string `json:"images"`
    }

As you can see, we are going to store just file names, and the actual files will be stored directly on the disk. Our API is extended with three methods:

    [GET]    /api/items/{id}/images/{img}
    [POST]   /api/items/{id}/images/
    [DELETE] /api/items/{id}/images/{img}

### Method handlers

In order to save an image correctly, we need to persist it to the disk, and update the matching Item object. Saving part is pretty obvious, as `http.Request` allows us to get file and its headers using `req.FormFile(<file_param_name>)` method. We obviously need to check for errors, or that the parameter is not empty etc. Then, we can proceed it to actual saving file and updating the DB.

    // api/api.go
    func UploadImage(w http.ResponseWriter, req *http.Request) {
        vars := mux.Vars(req)
        itemID := vars["id"]

        file, headers, err := req.FormFile("image")

        db.SaveImage(itemID, headers, file)

        w.Write([]byte("OK"))
    }
    
The save itself is pretty simple, as we just rename file to add Item ID prefix, persist it to the disk using `os.Create` (creating a file) and `io.Copy` (coping content to a new file) methods. Then, we use `mgo`'s `Update` and everything happens almost automagically. 

    func SaveImage(itemID string, headers *multipart.FileHeader, file multipart.File) error {
        filename := fmt.Sprintf("%s_%s", itemID, headers.Filename)
        saved, err := os.Create("./media/upload/" + filename)
        
        _, err = io.Copy(saved, file)

        query := bson.M{"_id": itemID}
        ... // loading Item to "res" variable 
        res.Images = append(res.Images, headers.Filename)
        return collection().Update(query, res)
    }

Deleting an image is a bit tricky, as we need to do the reverse: delete a file and update an item. It might seem easy, but removing an element from the slice is not trivial:

    func removeImageFromItem(it *Item, imageID string) {
        for i, name := range it.Images {
            if name == imageID {
                it.Images = it.Images[:i+copy(it.Images[i:], it.Images[i+1:])]
                return
            }
        }
    }

To remove an element we first need to find its index. Then we construct a slice by **replacing** a part of the slice so that the element is not in it any more. Do play around with this piece of code, because it looks very difficult, but it's pretty easy once you take it step-by-step.

Finally, getting an image back to the user is very simple, as we just read bytes from the disk and send them further to the `http.ResponseWriter`.

**Note** that now we need to update the method for deleting an Item object, to delete related images as well.

The code is available [on Github](https://github.com/mycodesmells/mongo-go-api).
