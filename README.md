bookmarks
=========

Personal zero-touch bookmarking app in the cloud, with full text search support.


installation
------------

1.  First get a Heroku account.
2.  Deploy a Node.js based app.
3.  Choose `GitHub` as your deployment method.
4.  Connect it to this repository or your fork.
5.  Install the Heroku Postgres add-on for this app.
6.  Under `settings` add a new config var `API_SECRET` and set this to
    something long and random. This is your password that clients need to
    manage bookmarks.
7.  Visit your newly created webapp at https://[name].herokuapp.com
8.  Fill in your `API_SECRET` in the search bar and hit `login`.
