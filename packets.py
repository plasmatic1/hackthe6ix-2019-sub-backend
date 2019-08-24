import json


class Packet:
    def __str__(self):
        return json.dumps(self)

    def __repr__(self):
        return json.dumps(self, indent=4)


class OutIdPacket(Packet):
    def __init__(self, id):
        self.id = id


class OutDistPacket(Packet):
    def __init__(self, dist):
        self.dist = dist
