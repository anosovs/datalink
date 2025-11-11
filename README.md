# datalink

Allow to share some secret data between users for selected count of views.
If some of bots request url - app will return 404 page


### ENV
Required .env file in app root directory.
.env consists:
- **SERVER_HOST** - where app hosted, example "127.0.0.1:8080". Default: ":8080"
- **STORAGE** - enum(sqlite, ram). Which type of storage to use
- **DELETE_AFTER** - Message will be erased after this count of days
- **HTML_TITLE** - Title for html-pages. Default: Datalink

### Usage  
Create .env in app dir  
In app dir `docker build -t datalink .`  
After `docker run -it -p 8080:8080 datalink`  
Visit http://localhost:8080/  