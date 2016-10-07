var express = require('express'),
    read = require("node-readability"),
    sanitizer = require("sanitizer"),
    pgp = require("pg-promise")();

var app = express(),
    db = pgp(process.env.DATABASE_URL);

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

app.get('/', function (request, response) {
    scraper(request.query.uri, function (data) {
        db.none("INSERT INTO bookmarks (url, name, content) VALUES (${url}, ${title}, ${contents})", data).then(function () {
            response.status(204).json({status: 'ok'});
        }).catch(function (error) {
            response.status(500).json({status: 'error'});
        });
    });
});

app.get('/search', function (request, response) {
    db.manyOrNone("SELECT id, created, name, url FROM bookmarks WHERE to_tsvector('english', content) @@ to_tsquery('english', $1)", request.query.q).then(function (data) {
        response.status(200).json(data);
    }).catch(function (error) {
        response.status(500).json({status: 'error'});
    });
});

app.listen(app.get('port'), function() {
  console.log('Node app is running on port', app.get('port'));
});
