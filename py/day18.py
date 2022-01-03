import time
import heapq
import math
from dataclasses import dataclass
from collections import deque, defaultdict, Counter
from typing import *

import numpy as np
import itertools
import re

import fire


def main(
        #infile: str='test_input.txt',
        infile: str = 'input.txt',
):
    print('hi!')

    lines = []
    with open(f'../inputs/day18/{infile}', 'r') as f:
        for i, line in enumerate(f.readlines()):
            lines.append(eval(line))

    print('=========')
    snums = [parse(l) for l in lines]
    t0 = time.time()
    x = snums[0]
    for i in range(1, len(snums)):
        x = Snum(x, snums[i])
        reduce(x)

    print('1:', magnitude(x), time.time() - t0, 's', x)

    t1 = time.time()
    max_mag = 0
    for i, j in itertools.product(range(len(snums)), range(len(snums))):
        x = Snum(parse(lines[i]), parse(lines[j]))
        reduce(x)
        mag = magnitude(x)
        if mag > max_mag:
            max_mag = mag

    print('2:', max_mag, time.time() - t1, 's')


@dataclass
class Snum:
    l: Optional['Snum'] = None
    r: Optional['Snum'] = None
    val: Optional[int] = None

    def is_leaf(self):
        return self.val is not None

    def add(self, val, dir):
        if self.is_leaf():
            self.val += val
            return

        if dir == 'l':
            self.l.add(val, dir)
        else:
            self.r.add(val, dir)

    def explode(self, depth=0):
        if self.is_leaf():
            return

        if depth == 4:
            lval = self.l.val
            rval = self.r.val
            self.l = None
            self.r = None
            self.val = 0
            return lval, rval
        x = self.l.explode(depth+1)
        if x:
            lval, rval = x
            if rval:
                self.r.add(rval, 'l')
            return lval, 0
        x = self.r.explode(depth+1)
        if x:
            lval, rval = x
            if lval:
                self.l.add(lval, 'r')
            return 0, rval

    def __str__(self):
        return str(self.val) if self.is_leaf() else f'[{str(self.l)},{str(self.r)}]'


def explode(x):
    depth = 0
    stack = list()
    cur = x
    prevprev, prev, rval = None, None, None
    exploded = False
    while cur or len(stack):
        if cur:
            stack.append((cur, depth))
            cur = cur.l
            depth += 1
        else:
            cur, depth = stack.pop()
            if cur.is_leaf():
                if rval is not None:
                    cur.val += rval
                    return True
                prevprev = prev
                prev = cur
            elif depth == 4:
                lval, rval = cur.l.val, cur.r.val
                cur.l, cur.r = None, None
                cur.val = 0
                exploded = True
                if prevprev:
                    prevprev.val += lval

            cur = cur.r
            depth += 1

    return exploded


def split(node):
    stack = list()
    cur = node
    while cur or len(stack):
        if cur:
            stack.append(cur)
            cur = cur.l
        else:
            cur = stack.pop()
            if cur.is_leaf() and cur.val > 9:
                cur.l = Snum(val=cur.val // 2)
                cur.r = Snum(val=(cur.val + 1) // 2)
                cur.val = None
                return True
            cur = cur.r

    return False


def magnitude(x):
    if x.is_leaf():
        return x.val
    else:
        return 3 * magnitude(x.l) + 2 * magnitude(x.r)


def reduce(x):
    while True:
        if explode(x):
            continue
        if split(x):
            continue
        break


def parse(x):
    if isinstance(x, list):
        return Snum(parse(x[0]), parse(x[1]))
    else:
        return Snum(val=x)


if __name__ == '__main__':
    fire.Fire(main)
