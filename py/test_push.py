#coding:utf-8

import urllib
import urllib2
import time
from binascii import crc32

def run():
	url = 'http://127.0.0.1:8081/v1/push'
	access_key = "bRVYoA5Y70m3dHYk"
	secret_key = "rP#S,9O4kl]GwjOD"
	key = '18612024052,15819459262'
	message = '{"t": 1, "id": %s, "data": {"a": "abc%s"}}' % (int(time.time()), int(time.time()))
	vc = crc32unsigned('%s%s%s' % (key, message, secret_key))
	params = {
		'ac' : access_key,
		'key': key,
		'message': message,
		'vc' : vc,
	}

	print url, params
	response = urllib2.urlopen(url, urllib.urlencode(params))
	text = response.read()
	print text

def crc32unsigned(b):
    value = crc32(b)
    if value < 0:
        value = 0xFFFFFFFF & value
    return value

run()