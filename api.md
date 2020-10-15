## 空闲教室API

### Notes

1. 南湖综合楼无课的教室会锁，只有 N101 教室开放自习

### 获取空闲教室

| Method | Header | URL |
| :----: | :----: | :-: |
| GET    | -      | api/classroom/v2 |

**URL Params**

```
week: 周次(1-21)
day: 星期(1-7)
building：楼栋(7,8,N)，分别为7号楼、8号楼和南湖综合楼
```

**Response Data**

```json
   {
      "code": 0,
      "message": "OK",
      "data": {
         "week": 8,
         "day": 4,
         "building": "8",
         "list": [
            {
               "time": 1, // 节次，对应每天的12节课
               "rooms": ["8101", "8102", "8103", "8106", "8109"]
            },{
               "time": 2,
               "rooms": ["8101"]
            },
            // ...
         ]
      }
   }

```
