# pg-migrate
Go module that handle simple migration for postgres database

## Install
```go
go install -v github.com/yatoogamii/pg-migrate@latest
```

## Use

Firstly setup your `DATABASE_URL` in the env. either by providing a .env file or by giving the env directly the in command line

Secondly, create a `migrations` folder with a folder for each migration. 
- your migration folder need to start with `{number}-`. for example `migrations/1-setup` or `migrations/5-new-migration`
- your file need to match the pattern : `{name}_{up or down}.sql`. for example `users_up.sql` or `my-new-migration_down.sql`

Example of structure :
```
migrations/
  1-users/
    users_up.sql
    users_down.sql
  2-companies/
    company_up.sql
    company_down.sql
```

Then you can use either `up` or `down` command. It will run all available migrations files 

```
$ pg-migrate up
```

```
$ pg-migrate down
```
