#JWPlayer

```shell script
docker build -t jwplayer .
docker run --name jwplayer --restart=always -p 80:80 -e DOMAIN={DOMAIN} jwplayer
```

```html
<script src="http://DOMAIN/player/v/8.15.0/jwplayer.js"></script>
```

```javascript
jwplayer.key = "enterprise/canPlayAds/1800000000000";
```