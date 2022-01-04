from dataclasses import dataclass
from typing import *
import numpy as np
import itertools
import re
import time

import fire


@dataclass
class Point:
    x: int
    y: int


@dataclass
class Line:
    l: Point
    r: Point


def main(
        #infile: str='test_input.txt',
        infile: str = 'input.txt',
):
    print('hi!')

    lines = []
    with open(infile, 'r') as f:
        for l in f.readlines():
            l, r = l.split('->')
            p1 = l.strip().split(',')
            p1 = Point(int(p1[0]), int(p1[1]))
            p2 = r.strip().split(',')
            p2 = Point(int(p2[0]), int(p2[1]))
            lines.append(Line(p1, p2))

    # print(lines)

    max_x = max(max(line.l.x, line.r.x) for line in lines) + 1
    max_y = max(max(line.l.y, line.r.y) for line in lines) + 1

    print(max_x, max_y)
    arr = np.zeros((max_y, max_x), dtype=np.int32)
    t0 = time.time()
    for line in lines:
        if line.l.x == line.r.x:
            y0 = min(line.l.y, line.r.y)
            y1 = max(line.l.y, line.r.y)
            for j in range((y1-y0)+1):
                arr[y0+j, line.l.x] += 1
        elif line.l.y == line.r.y:
            x0 = min(line.l.x, line.r.x)
            x1 = max(line.l.x, line.r.x)
            for j in range((x1-x0)+1):
                arr[line.l.y, x0+j] += 1
        else:
            # Diagonal
            x0, y0 = line.l.x, line.l.y
            x1, y1 = line.r.x, line.r.y

            if x0 < x1:
                dx = 1
            else:
                dx = -1
            if y0 < y1:
                dy = 1
            else:
                dy = -1

            for j in range(abs(x1-x0)+1):
                arr[y0+j*dy, x0+j*dx] += 1

    print(f'1: {(arr>1).sum()}', time.time() - t0, 's')


if __name__ == '__main__':
    fire.Fire(main)
