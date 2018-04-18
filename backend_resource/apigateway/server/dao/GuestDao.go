package dao

import (
	"apigateway/dao/database"	
)

func Login(loginCall,loginPassword string) (bool,int){
	db := database.GetConnection()
	var userID int
	err := db.QueryRow("SELECT userID FROM eo_admin WHERE loginCall = ? AND loginPassword = ?;",loginCall,loginPassword).Scan(&userID)
	if err != nil{
		return false,0
	}
	return true,userID
}

func Register(loginCall,loginPassword string) (bool){
	db := database.GetConnection()
	var userID int
	err := db.QueryRow("SELECT userID FROM eo_admin WHERE loginCall = ?",loginCall).Scan(&userID)
	if err == nil{
		return false
	}else{
		stmt,err:= db.Prepare("INSERT INTO eo_admin (loginCall,loginPassword) VALUES (?,?);")
		defer stmt.Close()
		_, err = stmt.Exec(loginCall,loginPassword)
		if err != nil {
			return false
		} else{
			return true
		}
	}
	
}

func CheckUserNameExist(loginCall string) bool{
	db := database.GetConnection()
	var userID int
	err := db.QueryRow("SELECT userID FROM eo_admin WHERE loginCall = ?",loginCall).Scan(&userID)
	if err != nil{
		return false
	}else{
		return true
	}
}