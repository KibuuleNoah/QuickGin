









# curl -X POST -H "Content-Type: application/json" -d '{"name":"Tristar","Identifier":"tristar02@gmail.com","password":"woshwosh"}' http://127.0.0.1:9000/v1/user/

# curl -X POST -H "Content-Type: application/json" -d '{"Identifier":"tristar3@gmail.com","password":"woshwosh"}' http://127.0.0.1:9000/v1/auth/with-password

# curl -X GET http://127.0.0.1:9000/v1/articles

# curl -X POST -H "Content-Type: application/json" -d '{"Identifier":"tristar3@gmail.com"}' http://127.0.0.1:9000/v1/auth/request-otp
# curl -X POST -H "Content-Type: application/json" -d '{"userId":"Ffibc9lZ", "otp": "832974"}' http://127.0.0.1:9000/v1/auth/with-otp

# {"message":"Successfully logged in","token":{"
# access_token":"
# "},"user":{"id":2,"identifier":"tristar3@gmail.com","verified":false,"name":"Tristar"}}%
#
#!/bin/bash

# Set your JWT token here
TOKEN=$(jq -r '.accessToken' token.json)
RTOKEN=$(jq -r '.refreshToken' token.json)


BASE_URL="http://127.0.0.1:9000/v1" 

create_user() {
  curl -X POST \
    $BASE_URL/user/ \
    -H "Content-Type: application/json" \
    -d '{"name":"'$1'","Identifier":"'$2'","password":"'$3'"}'
}


auth_with_password() {
  curl -X POST \
    $BASE_URL/auth/with-password \
    -H "Content-Type: application/json" \
    -d '{"Identifier":"'$1'","password":"'$2'"}'
}


get_articles() {
  curl -X GET $BASE_URL/articles \
       -H "Authorization: Bearer: $TOKEN"
}

request_otp() {
  curl -X POST \
    $BASE_URL/auth/request-otp \
    -H "Content-Type: application/json" \
    -d '{"Identifier":"'$1'"}'
}


auth_with_otp() {
  curl -X POST \
    $BASE_URL/auth/with-otp \
    -H "Content-Type: application/json" \
    -d '{"userId":"'$1'", "otp": "'$2'"}'
}

auth_token_refresh(){
  curl -X POST 
      $BASE_URL/auth/token/refresh \
      -H "Content-Type: application/json" \
      -d '{"token":"'$RTOKEN'"}'

}





# Example JSON file (data.json): { "name": "John", "age": 30 }






# create_user "Moxie" "moxie3@gmail.com"

# auth_with_password "tristar3@gmail.com" "woshwosh"
# get_articles
request_otp "moxie3@gmail.com"
# auth_with_otp "QMSpesgN" "015012"
# auth_token_refresh $RTOKEN
