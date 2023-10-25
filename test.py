import requests
import base64

s = requests.session()
user = {"login": "python3", "password": "python3password"}
auth_header = "X-Cenarius-Token"

header_value = "{} {}".format(user['login'], user['password'])
header_value_bytes = header_value.encode("ascii")
headers =  {auth_header: base64.b64encode(header_value_bytes) }


# # reg
# r = s.post("http://localhost:8080/api/v1/user/login", json=user)
# print(r.status_code, r.text)
# print(r.cookies)

r = s.get("http://localhost:8080/api/v1/private/health", headers=headers)
print(r.status_code, r.text)

r = s.get("http://localhost:8080/api/v1/private/loginwithpasswords", headers=headers)
print(r.status_code, r.text)

r = s.get("http://localhost:8080/api/v1/private/creditcards", headers=headers)
print(r.status_code, r.text)