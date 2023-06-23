SERVER:
```
sfetch -p 8888 -h xxx.com
```

CLIENT:
```
curl -i http://127.0.0.1:8888/xxx.com/a.jpg
```

Docker:
```
docker run -d -p 6000:6000 wanglei999/sfetch -p 6000 -h www.example.com
```
