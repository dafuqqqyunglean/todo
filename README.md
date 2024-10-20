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

![image](https://github.com/user-attachments/assets/36f0d69b-1723-4eb3-9570-039d4ebf456e)

### **Post:**

![image](https://github.com/user-attachments/assets/b4fe9a0d-3848-4cb4-bf4c-6005950b901a)

### **Get:**

![image](https://github.com/user-attachments/assets/57cb29ef-30a5-4ac6-93c9-305b8b6d34f9)

### **Put:**

![image](https://github.com/user-attachments/assets/1891045e-1f0d-4da8-8828-e87fc5e8624d)

### **Delete:**

![image](https://github.com/user-attachments/assets/17159d67-fd5a-41e4-bf78-64ec23abcdd2)
