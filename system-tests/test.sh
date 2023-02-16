curl -X POST -H "Content-Type: application/json" \
"http://localhost:8080/api/register" \
-d '{"teacher": "test@gmail.com", "students":["test1@gmail.com","test2@gmail.com"]}'
echo "\n"
curl -X POST -H "Content-Type: application/json" \
"http://localhost:8080/api/register" \
-d '{}'
echo "\n"
curl -X POST -H "Content-Type: application/json" \
"http://localhost:8080/api/register" \
-d '{"teacher": "test@gmail.com", "students":["test1@gmail.com","test2@gmail.com"], "extra": "extra"}'
echo "\n"
curl -X GET "http://localhost:8080/api/commonstudents?teacher=test%40gmail.com"
echo "\n"
curl -X GET "http://localhost:8080/api/commonstudents?teacher=nani%40gmail.com"
echo "\n"
curl -X POST -H "Content-Type: application/json" \
"http://localhost:8080/api/suspend" \
-d '{"student":"test1@gmail.com"}'
echo "\n"
curl -X POST -H "Content-Type: application/json" \
"http://localhost:8080/api/suspend" \
-d '{"student":"test3@gmail.com"}'
echo "\n"
curl -X POST -H "Content-Type: application/json" \
"http://localhost:8080/api/retrievefornotifications" \
-d '{"teacher":"test@gmail.com", "notification":"hello @test3@gmail.com and @test4@gmail.com and @test5@gmail.com"}'
echo "\n"
curl -X POST -H "Content-Type: application/json" \
"http://localhost:8080/api/retrievefornotifications" \
-d '{"teacher":"test@gmail.com", "notification":"hello @test2@gmail.com @test3@gmail.com and @test4@gmail.com and @test5@gmail.com"}'
echo "\n"
curl -X POST -H "Content-Type: application/json" \
"http://localhost:8080/api/register" \
-d '{"teacher": "tes8gmail.com", "students":["test7gmail.com","test6gmail.com"]}'