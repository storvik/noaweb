/*

noaweb - Not Another WEB framework

This project was created as I wanted to reuse my go code for several web projects. Use at your own risk!
External libraries used:
 - github.com/gorilla/csrf
 - github.com/gorilla/mux
 - github.com/lib/pq


Typical noaweb project

This describes a simple noaweb project with a database connection and two registered routes.
For a better understanding of all functions, please refere to the function documentation.

 // Initiate instance with optional config file.
 // An optional field UserConfig can be included in dev and prod in config.json.
 // This can be used for stuff that needs different values depending on dev or prod.
 s := noaweb.NewInstanceParse("config.json")

 var err error

 // Connect to a database. Postgres is the only suppported database.
 err = s.ConnectDB("postgres://user:password@localhost:port/databasename?sslmode=disable")
 if err != nil {
     panic("Could not connect to database")
 }
 // After connected to DB refere to given db connection with a string, user@databasename
 defer s.DisconnectDB("user@databasename")

 // Migrate database. Migrations should be created using the noaweb/noawebmigration tool.
 err = s.MigrateFolder("migrations/", "user@databasename")
 if err != nil {
     panic(err.Error())
 }

 // Initiate mailgun config. The config can later be refered to by using domain name.
 noaweb.Mailgun.Init("mg.mailgunhost.com", "key-mailgunkey", "pubkey-mailgunkey")

 // Register a route. This registers a GET route which renders a page.tmpl.html inside a layout.tmpl.html.
 // Note that the noaweb.Middleware.Logger middleware is used.
 s.RegRoute("/", "GET", noaweb.Middleware.Logger, func(w http.ResponseWriter, r *http.Request) {
     content, _ := noaweb.View.Parse("resources/views/page.tmpl.html")
     str, _ := noaweb.View.ParseTmpl("resources/views/layouts/layout.tmpl.html", map[string]interface{}{
         "csrfToken": noaweb.GetCSRF(r),      // noaweb.GetCSRF is a helper to the gorilla/csrf package.
         "content":   template.HTML(content),
     })
     fmt.Fprintf(w, str)
 })

 // Register API route
 s.RegRoute("/api/getdata", "POST", apiGetdata)

 // Register static routes
 s.RegStatic("/img/", "public/img/")

 // Print all routes
 fmt.Println("Registered routes:")
 s.ListAllRoutes()

 // Serve application
 s.Serve()


Development vs production

The config file consists of two different configuration objects, dev and prod.
This is useful to distinguish between development and production version of the app.
If the app is started with --dev, development config is used.
Note that the UserConfig section can be used to create custom options / variables.
These variables must be defined as a string(map[string]string), hence convertion must be handled.

When referring to files, ex. templates or other assets, AssetsDir in the configuration struct will be used.
AssetsDir is set to the location of the executable when running NewInstance().
If AssetsDir is included in the config file it will be overwritten.
For example this can be utilized to get assets from project directory in development mode, and use /var/appname/ as AssetsDir in production.
Note that all paths not prefixed with / will use path relative to AssetsDir, hence pahts starting with / will be the same path undependant of AssetsDir.


*/
package noaweb
