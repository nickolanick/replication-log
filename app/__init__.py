import os

from flask import Flask


app = Flask(__name__)
import api.routes


if __name__ == '__main__':
        app.run(debug=True, use_reloader=True)
