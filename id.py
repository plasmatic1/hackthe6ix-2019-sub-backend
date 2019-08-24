from random import randint

maxv = (1 << 30)
maxv = maxv - 1 + maxv

ids = set()


def add_id():
    rn = randint(1, maxv)
    while rn in ids:
        rn = randint(1, maxv)
    ids.add(rn)
    return rn


def rem_id(id):
    ids.discard(id)


def valid_id(id):
    return id in ids
