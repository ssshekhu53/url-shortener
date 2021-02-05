from flask import Flask, render_template, request, jsonify
import pyshorteners
app = Flask(__name__)


@app.route('/')
def home():
    return render_template('index.html')


@app.route('/tinyfy', methods=['POST'])
def tinyfy():
    shortener = pyshorteners.Shortener()
    response = {'short_url': shortener.tinyurl.short(request.form['urlField'])}
    print(response['short_url'])
    return jsonify(response)


if __name__ == '__main__':
    app.run(port=8080, debug=True)
