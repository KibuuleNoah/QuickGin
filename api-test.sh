









# curl -X POST -H "Content-Type: application/json" -d '{"name":"Tristar","Identifier":"tristar3@gmail.com","password":"woshwosh"}' http://127.0.0.1:9000/v1/user/

# curl -X POST -H "Content-Type: application/json" -d '{"Identifier":"tristar3@gmail.com","password":"woshwosh"}' http://127.0.0.1:9000/v1/auth/with-password

curl -X GET http://127.0.0.1:9000/v1/articles

# {"message":"Successfully logged in","token":{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjIxMGQxMDI1LTFlYjgtNDdhMy04NDQ1LTEwOWM2YjMyZTAzMSIsImF1dGhvcml6ZWQiOnRydWUsImV4cCI6MTc3NDk2MTgzMCwidXNlcl9pZCI6Mn0.Zn_AGxZaxNwzdc5PekW51p5mw56gPwOFrmRB7MrbZIw","refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzU1NjU3MzAsInJlZnJlc2hfdXVpZCI6ImMyMWQ5MGQxLTc0YWQtNDAwOC04NGU0LWZmY2NmODgyMWY1ZSIsInVzZXJfaWQiOjJ9.FOL1L8Bx8V2Yo0teE1AWL7iRVUxIqiHr3M9moX9jV5g"},"user":{"id":2,"identifier":"tristar3@gmail.com","verified":false,"name":"Tristar"}}%
