import math
from dataclasses import dataclass
from collections import deque, defaultdict
from typing import *

import numpy as np
import itertools
import re
import time

import fire


@dataclass
class Node:
    name: str

    def __post_init__(self):
        self.is_lower = (self.name == self.name.lower())
        self.edges: List['Node'] = list()


def main(infile: str):
    print('hi!')

    nodes = dict()
    with open(infile, 'r') as f:
        for line in f.readlines():
            l, r = line.strip().split('-')
            for x in [l, r]:
                if x not in nodes:
                    nodes[x] = Node(x)

            nodes[l].edges.append(nodes[r])
            nodes[r].edges.append(nodes[l])

    t0 = time.time()
    visited = set()
    def recurse(node):
        if node.name == 'end':
            return 1

        paths = 0
        if node.is_lower:
            visited.add(node.name)
        for edge in node.edges:
            if edge.name in visited:
                continue
            paths += recurse(edge)
        if node.is_lower:
            visited.remove(node.name)
        return paths

    paths = recurse(nodes['start'])
    print(f'1: {paths}', time.time() - t0, 's')

    t1 = time.time()
    vc = defaultdict(lambda: 0)
    vc['start'] = 2

    def recurse(node):
        if node.name == 'end':
            return 1

        if node.is_lower:
            vc[node.name] += 1

        paths = 0
        if any(x == 2 for x in vc.values()):
            max = 1
        else:
            max = 2
        for edge in node.edges:
            if vc[edge.name] >= max:
                continue
            paths += recurse(edge)

        if node.is_lower:
            vc[node.name] -= 1
        return paths

    paths = recurse(nodes['start'])
    print(f'2: {paths}', time.time() - t1, 's')


if __name__ == '__main__':
    fire.Fire(main)
