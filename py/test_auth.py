#coding:utf-8

import urllib2

CLIENT_VERSION = '201506031951'

def run():
	url = 'http://127.0.0.1/1/auth?user=zhouzhenhua&password=123456&did=01234567890&cv=%s' % CLIENT_VERSION
	response = urllib2.urlopen(url)
	text = response.read()
	print text

run()