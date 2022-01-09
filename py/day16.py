import heapq
import math
from dataclasses import dataclass
from collections import deque, defaultdict, Counter
from typing import *
import time

import numpy as np
import itertools
import re

import fire


def main(infile: str):
    print('hi!')

    with open(infile, 'r') as f:
        for i, line in enumerate(f.readlines()):
            input = line.strip()

    t0 = time.time()
    input_bin = bin(int(input, 16))[2:]
    if len(input_bin) % 4 > 0:
        input_bin = '0'*(4-len(input_bin) % 4) + input_bin

    parsed, _ = parse(input_bin)
    print(f'1: {parsed.get_version_sum()}', time.time() - t0, 's')
    t1 = time.time()
    print(f'2: {parsed.get_value()}', time.time() - t1, 's')


@dataclass
class Literal:
    v: int
    t: int
    val: int

    def get_version_sum(self):
        return self.v

    def get_value(self):
        return self.val


@dataclass
class Operator:
    v: int
    t: int
    I: int
    L: int
    subpackets: list

    def get_version_sum(self):
        return self.v + sum(sp.get_version_sum()
                            for sp in self.subpackets)

    def get_value(self):
        if self.t == 0:
            return sum(sp.get_value()
                       for sp in self.subpackets)
        elif self.t == 1:
            return math.prod(sp.get_value()
                             for sp in self.subpackets)
        elif self.t == 2:
            return min(sp.get_value()
                       for sp in self.subpackets)
        elif self.t == 3:
            return max(sp.get_value()
                       for sp in self.subpackets)
        elif self.t == 5:
            v0 = self.subpackets[0].get_value()
            v1 = self.subpackets[1].get_value()
            return 1 if v0 > v1 else 0
        elif self.t == 6:
            v0 = self.subpackets[0].get_value()
            v1 = self.subpackets[1].get_value()
            return 1 if v0 < v1 else 0
        elif self.t == 7:
            v0 = self.subpackets[0].get_value()
            v1 = self.subpackets[1].get_value()
            return 1 if v0 == v1 else 0
        else:
            raise ValueError(self.t)


def parse(b, i=0):
    v = int(b[i:i+3], 2)
    i += 3
    t = int(b[i:i+3], 2)
    i += 3

    if t == 4:
        # Literal
        vals = list()
        while True:
            leading = int(b[i])
            vals.append(b[i+1:i+5])
            i += 5
            if leading == 0:
                break
        return Literal(v=v, t=t, val=int(''.join(vals), 2)), i
    else:
        # operator
        I = int(b[i], 2)
        i += 1
        subpackets = []
        if I == 0:
            L = int(b[i:i+15], 2)
            i += 15
            fin = i+L
            while i < fin:
                x, i = parse(b, i)
                subpackets.append(x)
        else:
            L = int(b[i:i+11], 2)
            i += 11
            for _ in range(L):
                x, i = parse(b, i)
                subpackets.append(x)

        return Operator(v=v, t=t, I=I, L=L, subpackets=subpackets), i


if __name__ == '__main__':
    fire.Fire(main)
