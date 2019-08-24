import json
import math

from flask import Flask

from id import valid_id, rem_id, add_id
from packets import OutIdPacket, OutDistPacket

app = Flask(__name__)

teams = {}
locs = {}


def dist(a, b):
    """
    Dist between two points (a and b)
    :param a: Point A as a tuple (x, y)
    :param b: Point B as a tuple (x, y)
    :return: The distance (as a float)
    """
    da = a[0] - b[0]
    db = a[1] - b[1]
    return math.sqrt(da * da + db * db)


def error(msg):
    print(f'[ERROR]: {msg}')


# 261087034

@app.route('/del/<int:id>')
def del_user(id):
    if not valid_id(id):
        error(f'DEL: Invalid ID: {id}')
        return json.dumps({'error': f'Invalid ID {id}'})
    else:
        print(f'Removed ID {id}')
        rem_id(id)

    return ''


@app.route('/add/')
def add_user():
    new_id = add_id()
    print(f'Created user {new_id}')
    return json.dumps({'id': new_id})


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
        return json.dumps({'dist': best})


@app.route('/team/<int:id>/<int:team>')
def set_team(id, team):
    if not valid_id(id):
        error(f'SET TEAM: Invalid ID: {id}')
        return json.dumps({'error': f'Invalid ID {id}'})
    else:
        print(f'Set team of ID: {id} to {team}')
        teams[id] = team
        return ''


if __name__ == '__main__':
    app.run()
