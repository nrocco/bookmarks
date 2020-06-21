#!/usr/bin/env python3
import argparse
import json
import os
import subprocess
import urllib.request


def get_release(repository, tag):
    request = urllib.request.Request(method='GET', url='https://api.github.com/repos/{}/releases/tags/{}'.format(repository, tag))
    response = urllib.request.urlopen(request)
    if response.status != 200:
        return None
    return json.load(response)


def create_release(repository, tag, draft, prerelease):
    request = urllib.request.Request(method='POST', url='https://api.github.com/repos/{}/releases'.format(repository))
    request.add_header('Authorization', 'Bearer {}'.format(os.environ['GITHUB_TOKEN']))
    request.add_header('Content-Type', 'application/json')
    request.data = json.dumps({
        'name': 'Release {}'.format(tag),
        'tag_name': tag,
        'draft': draft,
        'prerelease': prerelease,
    }).encode()
    response = urllib.request.urlopen(request)
    if response.status != 201:
        raise Exception("Could not create release")
    return json.load(response)


def upload_asset(repository, release_id, asset):
    name = os.path.basename(asset)
    size = os.path.getsize(asset)
    content_type = subprocess.run(['file', '--brief', '--mime-type', asset], stdout=subprocess.PIPE).stdout.decode().strip() or 'application/octet-stream'
    file = open(asset, 'rb')

    request = urllib.request.Request(method='POST', url='https://uploads.github.com/repos/{}/releases/{}/assets?name={}'.format(repository, release_id, name), data=file)
    request.add_header('Authorization', 'Bearer {}'.format(os.environ['GITHUB_TOKEN']))
    request.add_header('Content-Length', size)
    request.add_header('Content-Type', content_type)

    response = urllib.request.urlopen(request)
    if response.status != 201:
        raise Exception("Could not upload asset {}".format(asset))

    return json.load(response)


def main():
    parser = argparse.ArgumentParser(description='Create a new github release and upload assets')
    parser.add_argument('--draft', action='store_true')
    parser.add_argument('--prerelease', action='store_true')
    parser.add_argument('repository')
    parser.add_argument('version')
    parser.add_argument('assets', nargs='+')
    args = parser.parse_args()

    release = create_release(args.repository, args.version, args.draft, args.prerelease)
    print("==> Created release {} {} with id {}".format(args.repository, args.version, release['id']))

    for asset in args.assets:
        print("==> Uploading {}".format(asset))
        result = upload_asset(args.repository, release['id'], asset)


if __name__ == '__main__':
    main()
