import json
import math

from flask import Flask

from dist_util import dist
from id import valid_id, rem_id, add_id

app = Flask(__name__)

teams = {}
locs = {}


def error(msg):
    print(f'[ERROR]: {msg}')


@app.route('/del/<int:id>')
def del_user(id):
    if not valid_id(id):
        error(f'DEL: Invalid ID: {id}')
        return json.dumps({'error': f'Invalid ID {id}'})
    else:
        print(f'Removed ID {id}')
        rem_id(id)

    return json.dumps({'error': ''})


@app.route('/add/')
def add_user():
    new_id = add_id()
    print(f'Created user {new_id}')
    return json.dumps({'id': new_id, 'error': ''})


@app.route('/loc/<int:id>/<float:lat>/<float:long>/')
def loc(id, lat, long):
    if not valid_id(id):
        error(f'SET LOC: Invalid ID: {id}')
        return json.dumps({'error': f'Invalid ID {id}'})
    else:
        cloc = (lat, long)
        locs[id] = cloc

        cteam = teams[id]
        best = math.inf
        for oth_id, loc in locs:
            if teams[oth_id] != cteam:
                best = min(best, dist(cloc, loc))

        print(f'Changed loc of ID {id} to {cloc}')
        return json.dumps({'dist': best, 'error': ''})


@app.route('/team/<int:id>/<int:team>')
def set_team(id, team):
    if not valid_id(id):
        error(f'SET TEAM: Invalid ID: {id}')
        return json.dumps({'error': f'Invalid ID {id}'})
    else:
        print(f'Set team of ID: {id} to {team}')
        teams[id] = team
        return json.dumps({'error': ''})


if __name__ == '__main__':
    app.run()
