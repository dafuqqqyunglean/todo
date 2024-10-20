# To Do Project
## Для сборки проекта:
```bash
docker-compose build
```
Для запуска проекта:
```bash
docker-compose up
```
Миграции:
```bash
docker compose exec app migrate -path ./schema -database 'postgres://todo_user:todo_password@postgres:5432/postgres?sslmode=disable' up
```
## Запросы
### **Аутентификация:**

![image](https://github.com/user-attachments/assets/0fbd2750-e715-48dc-ada3-98aa0b905e97)

![image](https://github.com/user-attachments/assets/1aa9e14a-9b5f-49d4-83e4-65582155cbfa)

### **Post:**

![image](https://github.com/user-attachments/assets/6d7c3027-9fb5-4cfb-91e6-d52a0f974446)

### **Get:**

![image](https://github.com/user-attachments/assets/9539da9a-7d19-49bc-aae6-1f621a4d55d5)

### **Put:**

![image](https://github.com/user-attachments/assets/21d8dece-7714-4c59-8a52-fc80dfb41973)

### **Delete:**

![image](https://github.com/user-attachments/assets/94d91422-a62b-44b5-8c0a-8cb813facf23)
