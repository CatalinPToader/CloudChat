FROM python:3.10
COPY requirements.txt /tmp
RUN pip install -r /tmp/requirements.txt
RUN mkdir /app
COPY fbAdminConfig_private_key.json /app/
COPY firebaseConfig.json /app/
COPY app.py /app/
ADD templates /app/templates
WORKDIR /app
EXPOSE 9000
ENV POST_USER="custom_user"
ENV POST_PASS="custom_passwd"
ENV POST_HOST="postgresql.postgresql-namespace"
ENV POST_DB="custom_db"
CMD ["python", "app.py"]