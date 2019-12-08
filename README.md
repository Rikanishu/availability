Assume you have something that can be booked (hotels, meeting rooms, etc.) so you store all date ranges when object is taken in a format `(object_id, date_from, date_to)` to recognize if the object is reserved for that time. 

Every time you want to check if the object is available in a particular time range, as a worst case you need to select all objects and exclude the object that are taken for this date ranges.

Assume you want to speed up this search, what do you need to do? The decision you get first is reverse the logic and keep a copy of the data but this time store available time. If you do it at db level you’ll have a duplication of data but let’s assume we can sacrifice the additional space. And of course, we’ll add index for this date fields to get extra speed. I thought so, what’s next? Can we improve the speed further? What if we store available time range in service memory, in the same format as db index tree? We’ll definitely get an improvement (of corse if we can sacrifice extra RAM).
This repository shows an example of using BTrees as in-memory cache for date intervals. So the lookup of objects takes O(log n) time complexity plus stores in RAM.
The repository has `calendar.csv` as some payload to initialize tree (it’s an AirBnB reservations dataset from Kaggle). The CSV has data in a format `(object_id, checkin_timestamp, checkout_timestamp)` to describe time ranges when the objects are occupied.

To compile and run just execute the following commands:

```
go build -o ./app && ./app
```
