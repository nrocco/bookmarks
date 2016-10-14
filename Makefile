public/favicon.ico:
	convert public/apple-touch-icon.png -define icon:auto-resize=64,48,32,16 public/favicon.ico

server:
	node_modules/nodemon/bin/nodemon.js index.js
