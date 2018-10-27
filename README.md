# Event sourcing

## Tài liệu tham khảo
* [Sách Domain Driven Design](https://www.amazon.com/Patterns-Principles-Practices-Domain-Driven-Design/dp/1118714709)
* [Sourcecode tham khảo](https://github.com/z0mbie42/goes)
* [Event Store trong Postgres](https://dev.to/kspeakman/event-storage-in-postgres-4dk2)
## Thiết kế CSDL
* Mỗi một microservice lưu trữ một event store và danh sách snapshot riêng (nếu có).
* Cấu trúc bảng:
    **Bảng event_store**: Lưu trữ event theo thời gian 
    | Tên cột | Kiểu dữ liệu | Chú thích | |
    | ------ | ------ | ------ | ------- |
    | AggregateID | string | Global unique id (ví dụ user_id) | PK|
    | AggregateType | string | Loại aggregate, là danh từ (ví dụ user, course, class) | |
    | Time | time | Thời gian tạo event | |
    | Version | integer| Version của aggregate, tăng dần khi có event mới | PK |
    | Data | jsonb | Nội dung event | |
    | EventType | string | Loại event tác động lên aggregate, là động từ (ví dụ UserCreated) | |
    | UserID | string | Id của user tạo event | |
    | Revision | integer |  Phòng trường hợp cấu trúc event payload bị thay đổi không để de-serialize từ JSONP về Golang Struct. | |
    | TenantID | string | Id cuả tenant | |
    
    **Bảng snapshot**: Lưu tạm trạng thái trung gian, tránh việc phải project lại chuỗi quá nhiều event
    | Tên cột | Kiểu dữ liệu | Chú thích | |
    | ------ | ------ | ------ | ------- |
    | AggregateID | string | Global unique id (ví dụ user_id) | PK|
    | Time | time | Thời gian tạo event | |
    | Version | integer| Version của aggregate tại thời điểm tạo snapshot | PK |
    | Data | jsonb | Trạng thái của aggregate tại thời điểm tạo snapshot | |
    | TenantID | string | Id của tenant | |
    

## Code

1. Tạo event 
    Định nghĩa:
    * Struct EventSource: đối tượng đại diện cho bảng event_store
    * Struct Snapshot: đối tượng đại diện cho bảng snapshot
    * Interface EventInterface: 
        ```go
            type EventInterface interface {
            	SaveReadDB(Event, config.Config) error
            	Apply(Aggregate, Event) error
            }
        ```
        =>  Các type event của aggregate cần triển khai (implement) interface EventInterface. Method SaveReadDB thực hiện thao tác với các bảng lưu trữ trạng thái cuối cùng của aggregate dùng để query. Method Apply được dùng khi getHistory hoặc project, set trạng thái mới cho aggregate khi có event mới. 
        => Với mỗi EventType, phải định nghĩa cách thức decode từ jsonb -> struct của Go trong hàm Decode. 
        Ví dụ: Event AddTeacherEvent 
        ```go
            type AddTeacherEvent struct {
            	CourseID string      `json:"course_id"`
            	Teacher  TeacherInfo `json:"teacher"`
            }
            
            func (event AddTeacherEvent) SaveReadDB(ev Event, db *pg.DB) error {
                // Do something
            }
            
            func (event AddTeacherEvent) Apply(agg Aggregate, ev Event) error {
                // Do something
            }
        ```
    * Struct Event
        ```go
            type Event struct {
            	AggregateId string      `json:"aggregate_id"`
            	Time        time.Time   `json:"time"`
            	Version     int32       `json:"version"`
            	Data        EventInterface `json:"data"`
            	EventType   string      `json:"event_type"`
            	UserID      string      `json:"user_id"`
            	Revision    int32       `json:"revision"`
            }
        ```
    * Interface Aggregate
        ```go
            type Aggregate interface {
                UpdateVersion()
            }
        ```
    * Struct BaseAggregate: implement interface Aggregate
        ```go
            type BaseAggregate struct {
                Version int32 `json:"version"`
            }
            func (agg *BaseAggregate) UpdateVersion() {
                agg.Version++
            }
        ```
        Định nghĩa truct aggregate chứa BaseAggregate. Ví dụ trong service-auth có aggregate User. Struct User có dạng:
        ```go
            type User struct {
                ...
                BaseAggregate
            }
        ```

    Để thêm event, gọi method SaveEvent trong file controller/event.go.
2. Tạo snapshot
    Cách n version của aggregate thì tạo một snapshot. n được quy định trong constant. 
3. Get Event stream 
4. Project

***
***Tham khảo code đầy đủ trong service Auth.***
