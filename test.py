import requests

s = requests.session()
user = {"login": "python3", "password": "python3password"}

# reg
r = s.post("http://localhost:8080/api/v1/user/login", json=user)
print(r.status_code, r.text)
print(r.cookies)

r = s.get("http://localhost:8080/api/v1/private/health")
print(r.status_code, r.text)
print(r.cookies)