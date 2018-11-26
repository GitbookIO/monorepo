import os
import sys
import path
import json


def read_monofile(filename):
    """A monofile"""
    return json.load(open(filename))


class CLI(object):
    def list(self):

    def pull(self):

    def main(self, action, *args):
        if action == "list":
            return self.list()
        elif action == "pull":

        elif action == "push":
            return self.push()
        print("monorepo ")
        print
        return ""


if __name__ == "__main__":
    cli = CLI()
    cli.main(*os.argv[1:])
