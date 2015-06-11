#coding:utf-8

import time
import json
import urllib
import urllib2

CLIENT_VERSION = '201506032051'

def do_auth(user, password, did):
	url = 'http://127.0.0.1/v1/auth'
	params = {
		'user': user,
		'password': password,
		'did': did,
		'cv': CLIENT_VERSION,
	}

	response = urllib2.urlopen(url, urllib.urlencode(params))
	data = json.loads(response.read())
	if data['r'] != 0:
		print 'auth failed'
	return data

def get_message_ids(token):
	url = 'http://127.0.0.1/v1/messages?token=%s&type=1&lastid=0&cv=%s' % (token, CLIENT_VERSION)
	response = urllib2.urlopen(url)
	data = json.loads(response.read())
	if data['r'] != 0:
		print 'messages failed'
	return data['ids']

def get_message_datas(token, ids):
	url = 'http://127.0.0.1/v1/message'
	params = {
		'token': token,
		'cv': CLIENT_VERSION,
		'id': ','.join(['%s' % id for id in ids]),
	}

	response = urllib2.urlopen(url, urllib.urlencode(params))
	data = json.loads(response.read())
	if data['r'] != 0:
		print 'messages failed'
	return data['msgs']

def do_sync(token, lastid):
	url = 'http://127.0.0.1/v1/syncid'
	params = {
		'token': token,
		'type': 1,
		'lastid': lastid,
		'cv': CLIENT_VERSION,
	}

	response = urllib2.urlopen(url, urllib.urlencode(params))
	data = json.loads(response.read())
	if data['r'] != 0:
		print 'messages failed'
	return data


if __name__ == "__main__":
	user = '18612024052'
	password = '12345678'
	did = 'androidid123456'

	auth_data = do_auth(user, password, did)
	print user, password, did, auth_data

	message_ids = get_message_ids(auth_data['token'])
	print message_ids

	print get_message_datas(auth_data['token'], message_ids)

	print do_sync(auth_data['token'], message_ids[-1])

	message_ids = get_message_ids(auth_data['token'])
	print message_ids