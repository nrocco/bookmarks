var express = require('express'),
    cookieParser = require('cookie-parser'),
    read = require("node-readability"),
    sanitizer = require("sanitizer"),
    pgp = require("pg-promise")();

var app = express(),
    db = pgp(process.env.DATABASE_URL),
    apiSecret = process.env.API_SECRET || 'foobar';

app.use(cookieParser())

app.use(function (request, response, next) {
    console.log(new Date().toISOString(), request.method, request.originalUrl);
    next();
});

app.use(function (request, response, next) {
    if (request.originalUrl !== '/' && (!request.cookies.secret || request.cookies.secret !== apiSecret)) {
        return response.sendStatus(401);
    }
    next();
});

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

app.set('port', (process.env.PORT || 5000));
app.use(express.static('public'));


// ============================================================


app.get('/bookmarks', function (request, response) {
    var params = {
        limit: request.query.limit || 25,
        offset: request.query.offset || 0,
    }

    if (request.query.q) {
        params.q = request.query.q
        var query = "SELECT id, created, name, url FROM bookmarks WHERE to_tsvector('english', content || ' ' || name || ' ' || url) @@ to_tsquery('english', ${q}) ORDER BY created DESC LIMIT ${limit} OFFSET ${offset}";
    } else {
        var query = "SELECT id, created, name, url FROM bookmarks ORDER BY created DESC LIMIT ${limit} OFFSET ${offset}";
    }

    db.manyOrNone(query, params).then(function (data) {
        response.status(200).json(data);
    }).catch(function (error) {
        response.status(500).json({status: 'error'});
    });
});

app.get('/bookmarks/add', function (request, response) {
    scraper(request.query.url, function (data) {
        db.none("INSERT INTO bookmarks (url, name, content) VALUES (${url}, ${title}, ${contents})", data).then(function () {
            response.redirect(request.query.url);
        }).catch(function (error) {
            response.status(500).json({status: 'error'});
        });
    });
});

app.get('/bookmarks/:id', function (request, response) {
    db.one("SELECT id, created, url, name, content FROM bookmarks WHERE id = $1", request.params.id).then(function (data) {
        response.status(200).json(data);
    }).catch(function (error) {
        response.status(404);
    });
});

app.delete('/bookmarks/:id', function (request, response) {
    db.none("DELETE FROM bookmarks WHERE id = $1 LIMIT 1", request.params.id).then(function (data) {
        response.status(204);
    }).catch(function (error) {
        response.status(404);
    });
});

app.listen(app.get('port'), function() {
  console.log('Node app is running on http://0.0.0.0:'+app.get('port'));
});
