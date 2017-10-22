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

    $ build/bookmarks-darwin-amd64


Usage
-----

Alernatively you can use the docker container:

    $ docker run -p 3000:3000 nrocco/bookmarks


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
