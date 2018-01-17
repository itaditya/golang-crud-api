# golang-crud-api using PostgreSQL db

Seed the database as follows

```sql

CREATE TABLE books (
  isbn    char(14) NOT NULL,
  title   varchar(255) NOT NULL,
  author  varchar(255) NOT NULL,
  price   decimal(5,2) NOT NULL
);

INSERT INTO books (isbn, title, author, price) VALUES
('978-1503261969', 'Emma', 'Jayne Austen', 9.44),
('978-1505255607', 'The Time Machine', 'H. G. Wells', 5.99),
('978-1503379640', 'The Prince', 'Niccolò Machiavelli', 6.99);

ALTER TABLE books ADD PRIMARY KEY (isbn);

```


Check out the following routes

<a href="/">Home</a>

<a href="/books">Books Index</a>


From terminal curl, use these


```sh

 curl -i -X POST -d "isbn=998-1470184841&title=Bad&author=Hali&price=1.90" localhost:8080/books
 
 curl -i -X GET localhost:8080/books
 
 curl -i -X GET localhost:8080/books/998-1470184841

```
