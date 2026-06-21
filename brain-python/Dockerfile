# (Keep all previous compilation lines exactly the same...)
FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .

# CHANGED RUNTIME INSTRUCTION: Run our consumer script directly
CMD ["python", "-m", "app.consumer"]