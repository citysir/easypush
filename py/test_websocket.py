import websocket
import thread
import time
import json
import urllib
import urllib2

def on_message(ws, message):
	print "on_message:", message
	reply = json.loads(message)
	reply["t"] = 2
	ws.send(json.dumps(reply))

def on_error(ws, error):
	print error

def on_close(ws):
	print "### closed ###"

def on_open(ws):
	def run(*args):
		while True:
			time.sleep(1)
	thread.start_new_thread(run, ())

def do_auth(user, password, did):
	url = 'http://127.0.0.1/v1/auth'
	params = {
		'user': user,
		'password': password,
		'did': did,
	}

	response = urllib2.urlopen(url, urllib.urlencode(params))
	data = json.loads(response.read())
	if data['r'] != 0:
		print 'auth failed'
	return data


if __name__ == "__main__":
	user = '18612024052'
	password = '12345678'
	did = 'androidid123456'

	auth_data = do_auth(user, password, did)

	websocket.enableTrace(True)

	url = "ws://%s:8080/v1/ws?token=%s&user=%s&did=%s" % (auth_data['node'], auth_data['token'], user, did)
	print url
	ws = websocket.WebSocketApp(url,
							  on_message = on_message,
							  on_error = on_error,
							  on_close = on_close)
	ws.on_open = on_open
	ws.run_forever()