import json
import geopy.distance as dist
from flask import Flask, render_template, redirect

from id import valid_id, rem_id, add_id

app = Flask(__name__)

teams = {}
locs = {}


def error(msg):
    print(f'[ERROR]: {msg}')


@app.route('/del/<int:id>/')
def del_user(id):
    if not valid_id(id):
        error(f'DEL: Invalid ID: {id}')
        return json.dumps({'error': f'Invalid ID {id}'})
    else:
        print(f'Removed ID {id}')
        del teams[id]
        del locs[id]
        rem_id(id)

    return json.dumps({'error': ''})


@app.route('/add/')
def add_user():
    new_id = add_id()

    # Defaults for loc and team
    teams[new_id] = 0
    locs[new_id] = (43.659743, -79.397698)  # (Lat, Long) for Bahen Centre

    print(f'Created user {new_id}')
    return json.dumps({'id': new_id, 'error': ''})


@app.route('/loc/<int:id>/<string:lat>/<string:long>/')
def loc(id, lat, long):
    try:
        lat = float(lat)
        long = float(long)
    except ValueError as e:
        return json.dumps({'error': f'Error while casting to float: {str(e)}'})

    if not valid_id(id):
        error(f'SET LOC: Invalid ID: {id}')
        return json.dumps({'error': f'Invalid ID {id}'})
    else:
        cloc = (lat, long)  # (Latitude, Longitude)
        locs[id] = cloc

        cteam = teams[id]
        best = 1e101  # We're going to define any value >1e100 as infinity
        for oth_id, loc in locs.items():
            if teams[oth_id] != cteam:
                best = min(best, float(dist.distance(cloc, loc).m))

        print(f'Changed loc of ID {id} to {cloc}')
        return json.dumps({'dist': best, 'error': ''})


@app.route('/team/<int:id>/<int:team>/')
def set_team(id, team):
    if not valid_id(id):
        error(f'SET TEAM: Invalid ID: {id}')
        return json.dumps({'error': f'Invalid ID {id}'})
    else:
        print(f'Set team of ID: {id} to {team}')
        teams[id] = team
        return json.dumps({'error': ''})


@app.route('/list/')
def list_info():
    users = []
    for k, v in teams.items():
        lat, long = locs[k]
        users.append((k, v, lat, long))

    return render_template('list.html', users=users)


@app.route('/del2/<int:id>/')
def del_user_redirect(id):
    del_user(id)
    return redirect('/list/')


@app.route('/add2/')
def add_user_redirect():
    add_user()
    return redirect('/list/')


@app.route('/team2/<int:id>/<int:team>/')
def set_team_redirect(id, team):
    set_team(id, team)
    return redirect('/list/')


if __name__ == '__main__':
    app.run()
