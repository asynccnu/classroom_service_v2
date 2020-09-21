package service

import (
	"fmt"
	"github.com/asynccnu/classroom_service_v2/model"
	"log"
	"github.com/tealeg/xlsx/v3"
)



func GetAvailableClassrooms(){
	// 7,8号楼和南湖所有教室
	allClassroomsIn7:=[]int {7101,7102,7103,7104,7105,7106,7107,7108,7109,
		7201,7202,7203,7204,7205,7206,7207,7208,7209,7211,
		7301,7302,7303,7304,7305,7306,7307,7308,7309,7311,
		7401,7402,7403,7404,7405,7406,7407,7408,7409,7410,7411,
		7501,7503,7505,
		}

	allClassroomsIn8:=[]int{
		8101,8102,8103,8104,8105,8106,8107,8108,8109,
		8110,8111,8112,821,8202,8203,8204,8205,8206,
		8207,8208,8209,8210,8211,8212,8213,8214,8215,
		8216,8301,8302,8303,8304,8305,8306,8307,8308,
		8309,8310,8311,8312,8313,8314,8315,8316,8401,
		8402,8403,8404,8405,8406,8407,8408,8409,8410,
		8411,8412,8413,8414,8415,8416,8501,8502,8503,
		8504,8505,8506,8507,8508,8509,8510,8511,8512,
		8513,8514,8515,8516,8716,8717,
	}
	// N都表示南湖
	allClassroomsInN:=string[]{

	}

	var building7 model.Rooms
	var building8 model.Rooms
	var buildingN model.Rooms

	building7={
		One
	}

	model.Insert()

	//file,err:=xlsx.OpenFile("/home/wency/data/2020—2021学年第一学期选课手册.xlsx")
	//if err!=nil{
	//	log.Printf("open file faild: %s",err)
	//}
	//for _,sheet:=range file.Sheets{
	//	fmt.Println(sheet.Name)
	//	for i:=0;i<=sheet.MaxRow;i++{
	//		cell,_:=sheet.Cell(i,10)
	//		fmt.Println(cell)
	//	}
	//	break
	//	sheet.Close()
	//}

}