# library-under-the-sea
A minimal web based REST API for a fictional public library

# Architecture
Follows a clean architecture layout with presentation layer separated from the application business rules as well as the database handlers.
HTTP request |-> library-api -> library -> library-repo |-> mongodb etc. 

Network communication between services via grpc which means:
* Simpler code as it handles many of the more complex code, such as request-response serialisation, for you and can be extended with middleware 
for logging, authentication, rate limiting and tracing as well as other.  
* More maintainable with backward compatibility via versioning as well as type checking to prevent breaking changes.
* More secure when used with TLS encryption
* Easier to use as requests and responses are defined in types. GRPC type checking for users of the API as well as copying comments to code.
* More performant as GRPC is built on solid foundations with protobuf and HTTP/2 because protobuf performs very well at serialization and HTTP/2 
provides a means for long-lasting connections, which gRPC takes advantage of.

## Todo
- [ ] Dockerise services
- [ ] Add Logger
- [ ] Add Mongo DB seeder
- [ ] Write more test cases incl. error scenarios
- [ ] End-to-end testing
- [ ] GRPC TLS encryption