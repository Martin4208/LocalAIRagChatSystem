users = ["admin", "guest", "guest", "user"]

for u in users:
    if u == "guest":
        users.remove(u)
        
print(users)