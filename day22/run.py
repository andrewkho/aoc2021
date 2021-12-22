import time
import heapq
import math
from dataclasses import dataclass, field
from collections import deque, defaultdict, Counter
from typing import *

from scipy.signal import convolve2d

import numpy as np
import itertools
import re

import fire


def main(
        #infile: str='test_input.txt',
        infile: str = 'input.txt',
):
    print('hi!')

    instructions = list()

    with open(f'../inputs/day22/{infile}', 'r') as f:
        for i, line in enumerate(f.readlines()):
            # on x=-3..43,y=-28..22,z=-6..38
            pattern = r'^(.*?) x=(.*?)\.\.(.*?),y=(.*?)\.\.(.*),z=(.*?)\.\.(.*?)$'
            m = re.match(pattern, line.strip())
            onoff, xmin, xmax, ymin, ymax, zmin, zmax = m.groups()
            instructions.append(
                Instr(
                    1 if onoff == 'on' else 0,
                    Box(
                        int(xmin), int(xmax),
                        int(ymin), int(ymax),
                        int(zmin), int(zmax),
                    )
                )
            )

    #print(instructions)
    t0 = time.time()
    print('start')
    grid = defaultdict(lambda: 0)
    fifty = Box(-50, 50, -50, 50, -50, 50)
    for instr in instructions:
        ix = fifty.intersect(instr.box)
        if not ix:
            continue
        else:
            for i, j, k in itertools.product(range(ix.xmin, ix.xmax+1),
                                             range(ix.ymin, ix.ymax+1),
                                             range(ix.zmin, ix.zmax+1)):

                grid[i, j, k] = instr.val

    c = 0
    for i, j, k in itertools.product(range(-50, 51),
                                     range(-50, 51),
                                     range(-50, 51)):
        c += grid[i,j,k]

    print('1:', c, time.time() - t0, 's')

    t1 = time.time()
    print('start1')
    r = Region()
    for instr in instructions:
        r.add(instr.box, instr.val)

    print('2:', r.area(), time.time() - t1, 's')


@dataclass
class Box:
    # boundaries are INCLUSIVE
    xmin: int
    xmax: int
    ymin: int
    ymax: int
    zmin: int
    zmax: int
    children: List['Box'] = field(default_factory=list)

    def intersect(self, o: 'Box') -> Optional['Box']:
        x0, x1 = max(o.xmin, self.xmin), min(o.xmax, self.xmax)
        if x0 > x1:
            return None

        y0, y1 = max(o.ymin, self.ymin), min(o.ymax, self.ymax)
        if y0 > y1:
            return None

        z0, z1 = max(o.zmin, self.zmin), min(o.zmax, self.zmax)
        if z0 > z1:
            return None

        return Box(
            x0, x1,
            y0, y1,
            z0, z1
        )

    def area(self) -> int:
        return ((self.xmax-self.xmin+1)*(self.ymax-self.ymin+1)*(self.zmax-self.zmin+1) -
                sum(c.area() for c in self.children))


@dataclass
class Instr:
    val: int
    box: Box


@dataclass
class Region:
    ones: List[Box] = field(default_factory=list)

    def add(self, box: Box, on_or_off: int):
        new_minus = list()
        minus_check = list()
        for one in self.ones:
            ix = one.intersect(box)
            if ix:
                new_minus.append(ix)
                minus_check.extend(one.children)
                one.children.append(ix)
        for minus_one in minus_check:
            ix = minus_one.intersect(box)
            if ix:
                self.ones.append(ix)

        if on_or_off:
            self.ones.append(box)

    def area(self) -> int:
        return sum(one.area() for one in self.ones)


if __name__ == '__main__':
    fire.Fire(main)
