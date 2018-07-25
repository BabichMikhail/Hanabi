# Hanabi - cooperative game with partitial information

# Requirements

- MySQL
- NodeJS
- Golang 1.6+

## Installation

### Repository

- ```go get github.com/BabichMikhail/Hanabi```
- copy app/app.conf.example to app/app.conf
- ```npm install```

### Database

- create database using MySQL
- set orm credentials in app/app.conf

## Run

### If you have bee tools
```  
bee run
```
### Else
```
go build main.go
```
Windows:
```
main
```
Linux:
```
./main
```

# Best Results

AI with recommendation strategy:
```
+------------+-------+-------+-------+-------+
| Players    |   2   |   3   |   4   |   5   |
+------------+-------+-------+-------+-------+
| Average    | 18.11 | 19.59 | 19.03 | 19.00 |
+------------+-------+-------+-------+-------+
| Perfect, % |  1.07 |  0.07 |  0.05 |     0 |
+------------+-------+-------+-------+-------+
```

AI with information strategy:

```
+------------+-------+-------+-------+-------+
| Players    |   2   |   3   |   4   |   5   |
+------------+-------+-------+-------+-------+
| Average    | 19.73 | 24.50 | 24.83 | 24.83 |
+------------+-------+-------+-------+-------+
| Perfect, % |  2.19 | 70.68 | 87.30 | 87.10 |
+------------+-------+-------+-------+-------+
```
