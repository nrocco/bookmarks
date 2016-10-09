var express = require('express'),
    cookieParser = require('cookie-parser'),
    read = require("node-readability"),
    sanitizer = require("sanitizer"),
    pgp = require("pg-promise")();

var app = express(),
    bookmarkRouter = express.Router();
    db = pgp(process.env.DATABASE_URL),
    listenPort = (process.env.PORT || 5000),
    apiSecret = process.env.API_SECRET || 'foobar';

function scraper(url, callback) {
    read(url, function(err, doc) {
        if (err) {
            throw err;
        }

        var obj = {
            "url": url,
            "title": doc.title.trim(),
            "contents": stripHTML(doc.content || "")
        };

        callback(obj);
    });
}

function stripHTML(html) {
    var clean = sanitizer.sanitize(html, function (str) {
        return str;
    });

    // Remove all remaining HTML tags.
    clean = clean.replace(/<(?:.|\n)*?>/gm, "");

    // RegEx to remove needless newlines and whitespace.
    // See: http://stackoverflow.com/questions/816085/removing-redundant-line-breaks-with-regular-expressions
    clean = clean.replace(/(?:(?:\r\n|\r|\n)\s*){2,}/ig, "\n");

    return clean.trim();
}


// ============================================================


bookmarkRouter.use(function (request, response, next) {
    if (!request.cookies.secret || request.cookies.secret !== apiSecret) {
        return response.sendStatus(401);
    }
    next();
});

bookmarkRouter.get('/bookmarks', function (request, response) {
    var params = {
        limit: request.query.limit || 25,
        offset: request.query.offset || 0,
        q: request.query.q
    }

    var query = [
        "SELECT id, created, name, url, substring(content, 0, 250) AS excerpt FROM bookmarks",
        (params.q) ? "WHERE fts @@ to_tsquery('english', ${q})" : "",
        "ORDER BY created DESC LIMIT ${limit} OFFSET ${offset}"
    ].join(" ");

    db.manyOrNone(query, params).then(function (data) {
        response.status(200).json(data);
    }).catch(function (error) {
        response.status(500).json({status: 'error'});
    });
});

bookmarkRouter.get('/bookmarks/add', function (request, response) {
    scraper(request.query.url, function (data) {
        db.none("INSERT INTO bookmarks (url, name, content, fts) VALUES (${url}, ${title}, ${contents}, setweight(to_tsvector('english', ${title}), 'A') || setweight(to_tsvector('english', replace(replace(replace(${url}, 'http', ''), '-', ' '), '/', ' ')), 'B') || setweight(to_tsvector('english', ${contents}), 'D'))", data).then(function () {
            response.redirect(request.query.url);
        }).catch(function (error) {
            response.status(500).json({status: 'error'});
        });
    });
});

bookmarkRouter.get('/bookmarks/:id', function (request, response) {
    db.one("SELECT id, created, url, name, content FROM bookmarks WHERE id = $1", request.params.id).then(function (data) {
        response.status(200).json(data);
    }).catch(function (error) {
        response.status(404);
    });
});

bookmarkRouter.delete('/bookmarks/:id', function (request, response) {
    db.none("DELETE FROM bookmarks WHERE id = $1 LIMIT 1", request.params.id).then(function (data) {
        response.status(204);
    }).catch(function (error) {
        response.status(404);
    });
});


// ============================================================


app.use(cookieParser());
app.use(express.static('public'));
app.use('/', bookmarkRouter);
app.listen(listenPort, function() {
  console.log('Node app is running on http://0.0.0.0:'+listenPort);
});
