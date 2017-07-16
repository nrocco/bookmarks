bookmarks
=========

Personal zero-touch bookmarking app in the cloud, with full text search support.


Installation
------------

First clone this git repository:

    $ git clone https://github.com/nrocco/bookmarks.git

Then compile

    $ make

Now you can run bookmarks server:

    $ build/bookmarks -help
    Usage of build/bookmarks:
      -database string
            The connection string of the database server
      -http string
            Address to listen for HTTP requests on (default "0.0.0.0:8000")
      -secret string
            The secret hash to authenticate to the api


Usage
-----

Alernatively you can use the docker container:

    $ docker run -p 8000:8000 nrocco/bookmarks -database "postgres://xxxxxxx"


Bookmarklet
-----------

    javascript:(function()%7Blocation.href=%22http://0.0.0.0/bookmarks/add?url=%22+encodeURIComponent(location.href)+%22&title=%22+encodeURIComponent(document.title);%7D)()


Contributing
------------

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Make sure that tests pass (`make test`)
5. Push to the branch (`git push origin my-new-feature`)
6. Create new Pull Request


Contributors
------------

- Nico Di Rocco (https://github.com/nrocco)
