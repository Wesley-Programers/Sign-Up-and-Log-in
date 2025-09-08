    package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)


type Data struct {
	name string
	email string
	password int
}

var dataSlice []Data
