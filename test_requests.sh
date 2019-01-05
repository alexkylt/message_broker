curl -H "Accept: application/json" -H 'Content-Type: application/json' -X PUT --data '{"name": "lalai7777", "value": "hopheilalalei"}' localhost:9090

wget -S –header="Accept-Encoding: gzip, deflate" -O response.txt localhost:9090/lalai7777

curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://localhost:9090/?lalai7777
