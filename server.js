#!/usr/bin/env node
var express = require('express'),
    app = express(),
    read = require("node-readability"),
    sanitizer = require("sanitizer");

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


app.get('/', function (req, res) {
    scraper(req.query.uri, function (data) {
        res.send(data.contents);
    });
});

app.listen(3000, function () {
    console.log('Example app listening on port 3000!');
});
