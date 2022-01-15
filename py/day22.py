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


def main(infile: str):
    print('hi!')

    instructions = list()

    with open(infile, 'r') as f:
        for i, line in enumerate(f.readlines()):
            # on x=-3..43,y=-28..22,z=-6..38
            pattern = r'^(.*?) x=(.*?)\.\.(.*?),y=(.*?)\.\.(.*),z=(.*?)\.\.(.*?)$'
            m = re.match(pattern, line.strip())
            onoff, xmin, xmax, ymin, ymax, zmin, zmax = m.groups()
            instructions.append(
                (Box(int(xmin), int(xmax), int(ymin), int(ymax), int(zmin), int(zmax)),
                 1 if onoff == 'on' else 0)
            )

    t0 = time.time()
    print('start')
    area = np.zeros((101, 101, 101), dtype=int)
    for box, onoff in instructions:
        fifty = Box(-50, 50, -50, 50, -50, 50)
        ix = fifty.intersect(box)
        if ix:
            area[50+ix.xmin:50+ix.xmax+1,
                 50+ix.ymin:50+ix.ymax+1,
                 50+ix.zmin:50+ix.zmax+1] = onoff

    print('1:', area.sum(), time.time() - t0, 's')

    t1 = time.time()
    print('start1')
    ones = list()
    for box, onoff, in instructions:
        for one in ones:
            one.add(box)
        if onoff:
            ones.append(box)

    print('2:', sum(one.area() for one in ones), time.time() - t1, 's')


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
            return

        y0, y1 = max(o.ymin, self.ymin), min(o.ymax, self.ymax)
        if y0 > y1:
            return

        z0, z1 = max(o.zmin, self.zmin), min(o.zmax, self.zmax)
        if z0 > z1:
            return

        return Box(x0, x1, y0, y1, z0, z1)

    def add(self, o):
        ix = self.intersect(o)
        if ix:
            for c in self.children:
                c.add(ix)
            self.children.append(ix)

    def area(self) -> int:
        return ((self.xmax-self.xmin+1)*(self.ymax-self.ymin+1)*(self.zmax-self.zmin+1) -
                sum(c.area() for c in self.children))


if __name__ == '__main__':
    fire.Fire(main)
